resource "google_storage_bucket" "this" {
  name                        = "fleps-assets"
  location                    = "europe-west2"
  uniform_bucket_level_access = true
}
resource "google_storage_bucket_iam_member" "this" {
  bucket = google_storage_bucket.this.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}
resource "google_storage_bucket_object" "tos" {
  name   = "tos.html"
  bucket = google_storage_bucket.this.name
  source = "${path.module}/public/tos.html"
}
resource "google_storage_bucket_object" "privacy" {
  name   = "privacy.html"
  bucket = google_storage_bucket.this.name
  source = "${path.module}/public/privacy.html"
}
resource "google_storage_bucket_object" "support" {
  name   = "support.html"
  bucket = google_storage_bucket.this.name
  source = "${path.module}/public/support.html"
}
resource "google_storage_bucket_object" "logo" {
  name   = "logo640x640.png"
  bucket = google_storage_bucket.this.name
  source = "${path.module}/public/logo640x640.png"
}