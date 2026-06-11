use tauri::AppHandle;
use serde_json::Value;
use crate::sidecar;

#[tauri::command]
pub async fn run_doctor(app: AppHandle) -> Result<Value, String> {
    sidecar::run_json(&app, vec!["doctor"]).await
}
