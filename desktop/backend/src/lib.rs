mod commands;
mod sidecar;

use commands::{compose, create, doctor, logs, projects};

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .plugin(tauri_plugin_dialog::init())
        .plugin(tauri_plugin_opener::init())
        .plugin(tauri_plugin_fs::init())
        .invoke_handler(tauri::generate_handler![
            projects::list_projects,
            projects::describe_project,
            projects::remove_project,
            compose::start_project,
            compose::stop_project,
            compose::get_status,
            compose::reconfigure_project,
            create::create_project,
            logs::stream_logs,
            doctor::run_doctor,
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
