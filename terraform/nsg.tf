resource "oci_core_network_security_group" "ssh" {
  compartment_id = var.compartment_ocid
  vcn_id         = oci_core_vcn.main.id
  display_name   = "${var.project_name}-ssh-nsg"
}
