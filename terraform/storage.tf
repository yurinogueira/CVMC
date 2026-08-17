resource "oci_core_volume" "data" {
  compartment_id      = var.compartment_ocid
  availability_domain = var.availability_domain

  display_name = "${var.project_name}-data"

  size_in_gbs = 8
  vpus_per_gb = 10
}

resource "oci_core_volume_attachment" "data" {
  attachment_type = "paravirtualized"

  instance_id = oci_core_instance.server.id
  volume_id   = oci_core_volume.data.id

  display_name = "${var.project_name}-data-attachment"
}
