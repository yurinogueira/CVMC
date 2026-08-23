output "compartment_id" {
  value = var.compartment_ocid
}

output "instance_id" {
  value = oci_core_instance.server.id
}

output "instance_public_ip" {
  value       = oci_core_public_ip.server_reserved_ip.ip_address
  description = "IP público estático reservado da instância backend"
}

output "instance_private_ip" {
  value = data.oci_core_private_ips.server_private_ips.private_ips[0].ip_address
}

output "object_storage_namespace" {
  value = data.oci_objectstorage_namespace.main.namespace
}

output "object_storage_bucket" {
  value = oci_objectstorage_bucket.files.name
}

output "ssh_command" {
  value = "ssh ubuntu@${oci_core_public_ip.server_reserved_ip.ip_address}"
}

output "mongodb_project_id" {
  value       = mongodbatlas_project.project.id
  description = "ID do projeto no MongoDB Atlas"
}

output "mongodb_cluster_connection_string" {
  value       = mongodbatlas_cluster.cluster.connection_strings[0].standard_srv
  description = "Standard SRV connection string do cluster MongoDB"
}

output "mongodb_app_username" {
  value       = mongodbatlas_database_user.app_user.username
  description = "Usuário da aplicação no MongoDB"
}

output "mongodb_app_password" {
  value       = random_password.mongodb_app_password.result
  description = "Senha do usuário da aplicação no MongoDB"
  sensitive   = true
}

output "mongodb_uri" {
  value       = "mongodb+srv://${mongodbatlas_database_user.app_user.username}:${random_password.mongodb_app_password.result}@${replace(mongodbatlas_cluster.cluster.connection_strings[0].standard_srv, "mongodb+srv://", "")}/?retryWrites=true&w=majority"
  description = "Connection URI completa para a aplicação (.env)"
  sensitive   = true
}

output "frontend_url" {
  value       = "https://${cloudflare_dns_record.frontend.name}.yurinogueira.dev.br"
  description = "URL de acesso ao Frontend"
}

output "backend_api_url" {
  value       = "https://${cloudflare_dns_record.backend.name}.yurinogueira.dev.br"
  description = "URL de acesso à API Backend"
}


