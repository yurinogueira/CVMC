output "compartment_id" {
  value = var.compartment_ocid
}

output "instance_id" {
  value = oci_core_instance.server.id
}

output "instance_public_ip" {
  value = oci_core_instance.server.public_ip
}

output "instance_private_ip" {
  value = oci_core_instance.server.private_ip
}

output "data_volume_id" {
  value = oci_core_volume.data.id
}

output "object_storage_namespace" {
  value = data.oci_objectstorage_namespace.main.namespace
}

output "object_storage_bucket" {
  value = oci_objectstorage_bucket.files.name
}

output "ssh_command" {
  value = "ssh ubuntu@${oci_core_instance.server.public_ip}"
}
