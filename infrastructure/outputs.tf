output "added_to_space_url" {
  value = google_cloud_run_v2_service.added_to_space.uri
}

output "removed_from_space_url" {
  value = google_cloud_run_v2_service.removed_from_space.uri
}

output "command_url" {
  value = google_cloud_run_v2_service.command.uri
}

output "message_url" {
  value = google_cloud_run_v2_service.message.uri
}