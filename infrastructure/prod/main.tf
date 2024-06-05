resource "google_cloud_run_v2_service" "prod-bv-ms" {
  name     = "prod-bv-ms"
  location = "us-central1"
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    scaling {
      max_instance_count = 100
      min_instance_count = 0
    }
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
        startup_cpu_boost = false
        cpu_idle          = true
      }
      env {
        name  = "ENDPOINT_URL"
        value = var.endpoint-url
      }
    }
  }
}

resource "google_cloud_run_service_iam_binding" "svc-iam-binding" {
  location = google_cloud_run_v2_service.prod-bv-ms.location
  service  = google_cloud_run_v2_service.prod-bv-ms.name
  role     = "roles/run.invoker"
  members = [
    "allUsers"
  ]
}
