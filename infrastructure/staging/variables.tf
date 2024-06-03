variable "gcp-config" {
  type = object({
    project = string
    region  = string
  })
}

variable "bv-ms-img" {
  type = string
}
