variable "google_project_id" {
  description = "The ID of the Google Cloud project."
  type        = string
}

variable "google_region" {
  description = "The region of the Google Cloud project."
  type        = string
}

variable "chat_service_account" {
  description = "The service account email for the Chat API."
  type        = string
}
variable "main_space" {
  description = "Chat space (e.g. spaces/AAAA) where submitted predictions are announced. Leave empty to disable."
  type        = string
  default     = ""
}
