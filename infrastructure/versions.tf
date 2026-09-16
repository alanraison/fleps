terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "8.3.0"
    }
  }

  backend "gcs" {
    bucket = "fleps-tofu-state.newredo.com"
    prefix = "state/"
  }
}

provider "google" {
  project = var.google_project_id
  region  = var.google_region
}
