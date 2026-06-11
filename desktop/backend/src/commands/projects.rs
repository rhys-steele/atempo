use tauri::AppHandle;
use serde_json::Value;
use tauri_plugin_shell::ShellExt;
use crate::sidecar;

#[tauri::command]
pub async fn list_projects(app: AppHandle) -> Result<Value, String> {
    sidecar::run_json(&app, vec!["projects"]).await
}

#[tauri::command]
pub async fn describe_project(app: AppHandle, name: String) -> Result<Value, String> {
    sidecar::run_json(&app, vec!["describe", &name]).await
}

#[tauri::command]
pub async fn remove_project(app: AppHandle, name: String) -> Result<(), String> {
    let args = vec!["remove", name.as_str()];
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

