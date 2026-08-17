resource "oci_objectstorage_bucket" "files" {
  compartment_id = var.compartment_ocid
  namespace      = data.oci_objectstorage_namespace.main.namespace

  name = "${var.project_name}-files"

  access_type = "NoPublicAccess"

  storage_tier = "Standard"

  versioning = "Enabled"
}

data "oci_objectstorage_namespace" "main" {
  compartment_id = var.compartment_ocid
}
