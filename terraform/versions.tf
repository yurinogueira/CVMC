terraform {
  required_version = ">= 1.15.8"

  required_providers {
    oci = {
      source  = "oracle/oci"
      version = "~> 9.2.0"
    }
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 2.18.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.9.0"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 5.25.0"
    }
  }
}
