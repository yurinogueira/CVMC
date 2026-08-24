data "oci_identity_availability_domains" "ads" {
  compartment_id = var.tenancy_ocid
}

data "oci_core_images" "ubuntu" {
  compartment_id = var.compartment_ocid

  operating_system         = "Canonical Ubuntu"
  operating_system_version = "24.04"
  shape                    = "VM.Standard.E2.1.Micro"

  sort_by    = "TIMECREATED"
  sort_order = "DESC"
}

resource "oci_core_instance" "server" {
  availability_domain = data.oci_identity_availability_domains.ads.availability_domains[0].name
  compartment_id      = var.compartment_ocid

  display_name = "${var.project_name}-server"

  shape = "VM.Standard.E2.1.Micro"

  create_vnic_details {
    subnet_id                 = oci_core_subnet.public.id
    assign_public_ip          = false
    assign_private_dns_record = true
    nsg_ids                   = [oci_core_network_security_group.ssh.id]

    hostname_label = "server"
  }

  source_details {
    source_id   = data.oci_core_images.ubuntu.images[0].id
    source_type = "image"

    boot_volume_size_in_gbs = 50
  }

  metadata = {
    ssh_authorized_keys = var.ssh_public_key
    user_data           = base64encode(file("${path.module}/cloud-init.sh"))
  }

  preserve_boot_volume = true

  lifecycle {
    ignore_changes = [
      metadata,
      source_details,
    ]
    prevent_destroy = true
  }
}

data "oci_core_vnic_attachments" "server_vnics" {
  compartment_id = var.compartment_ocid
  instance_id    = oci_core_instance.server.id
}

data "oci_core_private_ips" "server_private_ips" {
  vnic_id = data.oci_core_vnic_attachments.server_vnics.vnic_attachments[0].vnic_id
}

resource "oci_core_public_ip" "server_reserved_ip" {
  compartment_id = var.compartment_ocid
  display_name   = "${var.project_name}-backend-reserved-ip"
  lifetime       = "RESERVED"
  private_ip_id  = data.oci_core_private_ips.server_private_ips.private_ips[0].id
}

