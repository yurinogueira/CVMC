terraform {
  required_version = ">= 1.15.8"

  required_providers {
    oci = {
      source  = "oracle/oci"
      version = "~> 8.28.0"
    }
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.28.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6.0"
    }
  }
}
