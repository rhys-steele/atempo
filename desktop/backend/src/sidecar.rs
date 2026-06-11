use serde_json::Value;
use tauri::{AppHandle, Emitter};
use tauri_plugin_shell::ShellExt;
use tauri_plugin_shell::process::CommandEvent;

/// Run an atempo command with --json and return parsed output.
pub async fn run_json(app: &AppHandle, args: Vec<&str>) -> Result<Value, String> {
    let mut full_args = args;
    full_args.push("--json");

    let sidecar = app
        .shell()
        .sidecar("atempo")
        .map_err(|e| e.to_string())?
        .args(&full_args);

    let output = sidecar
        .output()
        .await
        .map_err(|e| e.to_string())?;

    if !output.status.success() {
        let stderr = String::from_utf8_lossy(&output.stderr).to_string();
        return Err(stderr);
    }

    let stdout = String::from_utf8_lossy(&output.stdout).to_string();
    serde_json::from_str(&stdout).map_err(|e| format!("parse error: {e}\noutput: {stdout}"))
}

/// Spawn an atempo command with --stream-json, emitting each line as a Tauri event.
pub async fn stream_json(app: &AppHandle, args: Vec<&str>, event_name: &str) -> Result<(), String> {
    let mut full_args = args;
    full_args.push("--stream-json");

    let (mut rx, _child) = app
        .shell()
        .sidecar("atempo")
        .map_err(|e| e.to_string())?
        .args(&full_args)
        .spawn()
        .map_err(|e| e.to_string())?;

    let app = app.clone();
    let event_name = event_name.to_string();

    tauri::async_runtime::spawn(async move {
        while let Some(event) = rx.recv().await {
            match event {
                CommandEvent::Stdout(line) => {
                    let text = String::from_utf8_lossy(&line).to_string();
                    if let Ok(val) = serde_json::from_str::<Value>(&text) {
                        let _ = app.emit(&event_name, val);
                    }
                }
                CommandEvent::Stderr(line) => {
                    let text = String::from_utf8_lossy(&line).to_string();
                    let _ = app.emit(&event_name, serde_json::json!({
                        "type": "error",
                        "data": text
                    }));
                }
                CommandEvent::Terminated(_) => {
                    let _ = app.emit(&event_name, serde_json::json!({ "type": "done" }));
                    break;
                }
                _ => {}
            }
        }
    });

    Ok(())
}
