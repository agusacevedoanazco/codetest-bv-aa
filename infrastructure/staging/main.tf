resource "google_cloud_run_service" "staging-bv-ms" {
  name     = "staging-bv-ms"
  location = var.gcp-config.region
  project  = var.gcp-config.project

  template {
    spec {
      containers {
        image = var.bv-ms-img
        ports {
          container_port = 32080
        }
        resources {
          limits = {
            cpu    = "200m"
            memory = "200Mi"
          }
          requests = {
            cpu    = "100m"
            memory = "100Mi"
          }
        }
      }
      timeout_seconds = 90

    }
    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "2"
        "autoscaling.knative.dev/minScale" = "0"
        "run.googleapis.com/ingress"       = "all"
      }
    }
  }

  autogenerate_revision_name = true
}

data "google_iam_policy" "allowall" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "allowall" {
  location = google_cloud_run_service.staging-bv-ms.location
  project  = google_cloud_run_service.staging-bv-ms.project
  service  = google_cloud_run_service.staging-bv-ms.name

  policy_data = data.google_iam_policy.allowall
}
