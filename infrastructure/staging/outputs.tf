output "cloud_run_svc_url" {
  value = google_cloud_run_v2_service.staging-bv-ms.uri
}
