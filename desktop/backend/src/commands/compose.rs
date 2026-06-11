use tauri::AppHandle;
use serde_json::Value;
use tauri_plugin_shell::ShellExt;
use crate::sidecar;

#[tauri::command]
pub async fn start_project(app: AppHandle, name: String) -> Result<(), String> {
    run_simple(&app, vec!["start", &name]).await
}

#[tauri::command]
pub async fn stop_project(app: AppHandle, name: String) -> Result<(), String> {
    run_simple(&app, vec!["stop", &name]).await
}

#[tauri::command]
pub async fn get_status(app: AppHandle, name: String) -> Result<Value, String> {
    sidecar::run_json(&app, vec!["status", &name]).await
}

#[tauri::command]
pub async fn reconfigure_project(app: AppHandle, name: String) -> Result<(), String> {
    run_simple(&app, vec!["reconfigure", &name]).await
}

async fn run_simple(app: &AppHandle, args: Vec<&str>) -> Result<(), String> {
    let out = app
        .shell()
        .sidecar("atempo")
        .map_err(|e| e.to_string())?
        .args(&args)
        .output()
        .await
        .map_err(|e| e.to_string())?;

    if out.status.success() { Ok(()) } else {
        Err(String::from_utf8_lossy(&out.stderr).to_string())
    }
}
