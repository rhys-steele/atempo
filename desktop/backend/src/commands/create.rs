use tauri::AppHandle;
use crate::sidecar;

#[tauri::command]
pub async fn create_project(
    app: AppHandle,
    framework: String,
    name: String,
    version: Option<String>,
    directory: String,
) -> Result<(), String> {
    let mut args = vec!["create", framework.as_str(), name.as_str()];
    if let Some(ref v) = version {
        args.extend_from_slice(&["--version", v.as_str()]);
    }
    args.extend_from_slice(&["--dir", directory.as_str()]);

    sidecar::stream_json(&app, args, "create:event").await
}
