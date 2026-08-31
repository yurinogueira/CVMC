terraform {
  required_version = ">= 1.15.8"

  required_providers {
    oci = {
      source  = "oracle/oci"
      version = "~> 8.29.0"
    }
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 2.17.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.9.0"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 5.24.0"
    }
  }
}
