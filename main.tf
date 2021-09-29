terraform {
  required_version = ">= 0.14"

  required_providers {
    # Cloud Run support was added on 3.3.0
    google = ">= 3.17"
  }
}

provider "google" {
  project = "${var.project_id}"
  region  = "${var.region}"
  credentials = "${var.key}"
}

resource "google_project_service" "run_api" {
  service = "run.googleapis.com"

  disable_on_destroy = true
}

resource "google_project_service" "firestore" {
  service = "firestore.googleapis.com"
  disable_on_destroy = false
}

resource "google_firestore_document" "initial-pet" {
  collection  = "Pet"
  document_id = "initial-pet"
  fields      = "{\"petname\":{\"stringValue\":\"Gopher\"}, \"owner\":{\"stringValue\":\"Fairy\"}, \"likes\":{\"integerValue\":\"20\"}, \"image\":{\"stringValue\":\"https://cdn.shopify.com/s/files/1/0541/3955/1941/files/gopher_480x480.png?v=1621428369\"}}"
}

resource "google_cloud_run_service" "run_service" {
  name = "app"
  location = "${var.region}"

  template {
    spec {
      containers {
        env {
          name = "GOOGLE_CLOUD_PROJECT"
          value = "${var.project_id}"
        }
        env {
          name = "GOOGLE_APPLICATION_CREDENTIALS"
          value = "${var.key}"
        }
        image = "gcr.io/${var.project_id}/go-pets:latest"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  # Waits for the Cloud Run API to be enabled
  depends_on = [google_project_service.run_api]
}

# Allow unauthenticated users to invoke the service
resource "google_cloud_run_service_iam_member" "run_all_users" {
  service  = google_cloud_run_service.run_service.name
  location = google_cloud_run_service.run_service.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# Display the service URL
output "service_url" {
  value = google_cloud_run_service.run_service.status[0].url
}