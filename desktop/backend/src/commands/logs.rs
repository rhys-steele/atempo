use tauri::AppHandle;
use crate::sidecar;

#[tauri::command]
pub async fn stream_logs(app: AppHandle, name: String) -> Result<(), String> {
    sidecar::stream_json(&app, vec!["logs", name.as_str()], "logs:event").await
}
