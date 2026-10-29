resource "oci_objectstorage_bucket" "files" {
  compartment_id = var.compartment_ocid
  namespace      = data.oci_objectstorage_namespace.main.namespace

  name = "${var.project_name}-files"

  access_type = "NoPublicAccess"

  storage_tier = "Standard"

  versioning = "Enabled"
}

resource "oci_objectstorage_preauthrequest" "files_par" {
  namespace             = data.oci_objectstorage_namespace.main.namespace
  bucket                = oci_objectstorage_bucket.files.name
  name                  = "${var.project_name}-files-par"
  access_type           = "AnyObjectReadWrite"
  bucket_listing_action = "ListObjects"
  time_expires          = "2035-12-31T23:59:59Z"
}

data "oci_objectstorage_namespace" "main" {
  compartment_id = var.compartment_ocid
}
