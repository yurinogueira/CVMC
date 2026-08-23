resource "mongodbatlas_project" "project" {
  name   = "${var.project_name}-atlas-project"
  org_id = var.atlas_org_id

  lifecycle {
    prevent_destroy = true
  }
}

resource "mongodbatlas_cluster" "cluster" {
  project_id                     = mongodbatlas_project.project.id
  name                           = var.mongodb_cluster_name
  provider_name                  = "TENANT"
  backing_provider_name          = "AWS"
  provider_region_name           = "SA_EAST_1"
  provider_instance_size_name    = "M0"
  termination_protection_enabled = true

  lifecycle {
    prevent_destroy = true
  }
}

resource "random_password" "mongodb_app_password" {
  length           = 24
  special          = true
  override_special = "!-_=+"
}

resource "mongodbatlas_database_user" "app_user" {
  username           = var.mongodb_app_username
  password           = random_password.mongodb_app_password.result
  project_id         = mongodbatlas_project.project.id
  auth_database_name = "admin"

  roles {
    role_name     = "readWriteAnyDatabase"
    database_name = "admin"
  }
}

resource "mongodbatlas_project_ip_access_list" "backend_server" {
  project_id = mongodbatlas_project.project.id
  ip_address = oci_core_public_ip.server_reserved_ip.ip_address
  comment    = "OCI VM backend static reserved public IP"
}
