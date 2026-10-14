resource "google_project_service" "chat" {
  service = "chat.googleapis.com"
}

resource "google_project_service" "artifact_registry" {
  service = "artifactregistry.googleapis.com"
}

resource "google_project_service" "cloud_run" {
  service = "run.googleapis.com"
}

resource "google_project_service" "marketplace" {
  service = "appsmarket-component.googleapis.com"
}

resource "google_project_service" "gsuiteaddons" {
  service = "gsuiteaddons.googleapis.com"
}

resource "google_artifact_registry_repository" "this" {
  format        = "DOCKER"
  repository_id = "fleps"

  depends_on = [google_project_service.artifact_registry]
}

resource "google_cloud_run_v2_service" "space" {
  name                = "fleps-space"
  location            = "europe-west2"
  description         = "Fleps Space handler"
  deletion_protection = false

  template {
    containers {
      image = "${google_artifact_registry_repository.this.registry_uri}/space:latest"
      env {
        name  = "PGHOST"
        value = google_sql_database_instance.this.connection_name
      }
      env {
        name  = "PGUSER"
        value = trimsuffix(google_service_account.this.email, ".gserviceaccount.com")
      }
      env {
        name  = "PGDATABASE"
        value = "fleps"
      }
      env {
        name  = "PGSSLMODE"
        value = "disable"
      }
    }

    service_account = google_service_account.this.email
  }

  depends_on = [google_project_service.cloud_run]
}

resource "google_cloud_run_v2_service_iam_member" "space" {
  name     = google_cloud_run_v2_service.space.name
  location = google_cloud_run_v2_service.space.location
  role     = "roles/run.invoker"
  member   = "serviceAccount:${var.chat_service_account}"
}

resource "google_cloud_run_v2_service" "command" {
  name                = "fleps-command"
  location            = "europe-west2"
  description         = "Fleps Command handler"
  deletion_protection = false

  template {
    containers {
      image = "${google_artifact_registry_repository.this.registry_uri}/command:latest"
      env {
        name  = "PGHOST"
        value = google_sql_database_instance.this.connection_name
      }
      env {
        name  = "PGUSER"
        value = trimsuffix(google_service_account.this.email, ".gserviceaccount.com")
      }
      env {
        name  = "PGDATABASE"
        value = "fleps"
      }
      env {
        name  = "PGSSLMODE"
        value = "disable"
      }
      env {
        name  = "ACTION_URL"
        value = google_cloud_run_v2_service.action.uri
      }
      env {
        name  = "ASSETS_BUCKET"
        value = google_storage_bucket.this.name
      }
    }
    service_account = google_service_account.this.email
  }
}

resource "google_storage_bucket_iam_member" "command_writer" {
  bucket = google_storage_bucket.this.name
  role   = "roles/storage.objectCreator"
  member = "serviceAccount:${google_service_account.this.email}"
}

resource "google_cloud_run_v2_service_iam_member" "command" {
  name     = google_cloud_run_v2_service.command.name
  location = google_cloud_run_v2_service.command.location
  role     = "roles/run.invoker"
  member   = "serviceAccount:${var.chat_service_account}"
}



resource "google_cloud_run_v2_service" "message" {
  name                = "fleps-message"
  location            = "europe-west2"
  description         = "Fleps Message handler"
  deletion_protection = false

  template {
    containers {
      image = "${google_artifact_registry_repository.this.registry_uri}/message:latest"
      env {
        name  = "PGHOST"
        value = google_sql_database_instance.this.connection_name
      }
      env {
        name  = "PGUSER"
        value = trimsuffix(google_service_account.this.email, ".gserviceaccount.com")
      }
      env {
        name  = "PGDATABASE"
        value = "fleps"
      }
      env {
        name  = "PGSSLMODE"
        value = "disable"
      }
    }
    service_account = google_service_account.this.email
  }
}

resource "google_cloud_run_v2_service_iam_member" "message" {
  name     = google_cloud_run_v2_service.message.name
  location = google_cloud_run_v2_service.message.location
  role     = "roles/run.invoker"
  member   = "serviceAccount:${var.chat_service_account}"
}

resource "google_cloud_run_v2_service" "action" {
  name                = "fleps-action"
  location            = "europe-west2"
  description         = "Fleps Action handler"
  deletion_protection = false

  template {
    containers {
      image = "${google_artifact_registry_repository.this.registry_uri}/action:latest"
      env {
        name  = "PGHOST"
        value = google_sql_database_instance.this.connection_name
      }
      env {
        name  = "PGUSER"
        value = trimsuffix(google_service_account.this.email, ".gserviceaccount.com")
      }
      env {
        name  = "PGDATABASE"
        value = "fleps"
      }
      env {
        name  = "PGSSLMODE"
        value = "disable"
      }
    }
    service_account = google_service_account.this.email
  }
}

resource "google_cloud_run_v2_service_iam_member" "action" {
  name     = google_cloud_run_v2_service.action.name
  location = google_cloud_run_v2_service.action.location
  role     = "roles/run.invoker"
  member   = "serviceAccount:${var.chat_service_account}"
}
