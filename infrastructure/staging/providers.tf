terraform {
  backend "gcs" {
    bucket = "tfstate-codetest-bv-aa"
    prefix = "staging/microservice"
  }

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "5.32.0"
    }
  }
}
