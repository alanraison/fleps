data "google_project" "current" {
}

resource "google_project_service" "sqladmin" {
  service = "sqladmin.googleapis.com"
}

resource "google_sql_database_instance" "this" {
  name             = "football-predictions"
  database_version = "POSTGRES_18"

  settings {
    tier                             = "db-custom-N4-2-8192"
    availability_type                = "ZONAL"
    disk_autoresize                  = true
    disk_size                        = 20
    edition                          = "ENTERPRISE"
    data_disk_provisioned_iops       = "3000"
    data_disk_provisioned_throughput = "140"

    backup_configuration {
      enabled = true
      backup_retention_settings {
        retained_backups = 7
        retention_unit   = "COUNT"
      }
    }
    database_flags {
      name  = "cloudsql.iam_authentication"
      value = "on"
    }
  }
}

resource "google_sql_database" "this" {
  name      = "fleps"
  instance  = google_sql_database_instance.this.name
  charset   = "UTF8"
  collation = "en_US.UTF8"
}

resource "google_service_account" "this" {
  account_id   = "fleps-sa"
  display_name = "Fleps Service Account"
}

resource "google_sql_user" "this" {
  name     = trimsuffix(google_service_account.this.email, ".gserviceaccount.com")
  instance = google_sql_database_instance.this.name
  type     = "CLOUD_IAM_SERVICE_ACCOUNT"
}

resource "google_project_iam_member" "cloud_sql_instance_user" {
  project = data.google_project.current.project_id
  role    = "roles/cloudsql.instanceUser"
  member  = "serviceAccount:${google_service_account.this.email}"
}

resource "google_project_iam_member" "cloud_sql_client" {
  project = data.google_project.current.project_id
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.this.email}"
}
