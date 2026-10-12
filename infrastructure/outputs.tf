output "space_url" {
  value = google_cloud_run_v2_service.space.uri
}

output "command_url" {
  value = google_cloud_run_v2_service.command.uri
}

output "message_url" {
  value = google_cloud_run_v2_service.message.uri
}