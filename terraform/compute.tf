data "oci_core_images" "ubuntu" {
  compartment_id = var.compartment_ocid

  operating_system         = "Canonical Ubuntu"
  operating_system_version = "24.04"
  shape                    = "VM.Standard.A1.Flex"

  sort_by    = "TIMECREATED"
  sort_order = "DESC"
}

resource "oci_core_instance" "server" {
  availability_domain = var.availability_domain
  compartment_id      = var.compartment_ocid

  display_name = "${var.project_name}-server"

  shape = "VM.Standard.A1.Flex"

  shape_config {
    ocpus         = 1
    memory_in_gbs = 6
  }

  create_vnic_details {
    subnet_id                 = oci_core_subnet.public.id
    assign_public_ip          = true
    assign_private_dns_record = true

    hostname_label = "server"
  }

  source_details {
    source_id   = data.oci_core_images.ubuntu.images[0].id
    source_type = "image"

    boot_volume_size_in_gbs = 50
  }

  metadata = {
    ssh_authorized_keys = var.ssh_public_key
  }

  preserve_boot_volume = true
}
