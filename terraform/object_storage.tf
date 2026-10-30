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

  # A API da OCI não retorna `bucket_listing_action` na leitura do PAR; sem este
  # ignore, todo refresh gera drift e força a recriação do PAR (rotacionando a URL
  # consumida pelo backend). Ver issue #147.
  lifecycle {
    ignore_changes = [bucket_listing_action]
  }
}

data "oci_objectstorage_namespace" "main" {
  compartment_id = var.compartment_ocid
}
