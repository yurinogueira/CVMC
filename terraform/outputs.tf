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

output "ssh_nsg_id" {
  value       = oci_core_network_security_group.ssh.id
  description = "OCID do NSG usado para regras dinâmicas de SSH"
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

output "smtp_host" {
  value       = "smtp.resend.com"
  description = "Host do servidor SMTP do Resend (Secret: SMTP_HOST)"
}

output "smtp_port" {
  value       = "587"
  description = "Porta STARTTLS do servidor SMTP (Secret: SMTP_PORT)"
}

output "smtp_user" {
  value       = "resend"
  description = "Usuário fixo do Resend SMTP (Secret: SMTP_USER)"
}

output "email_from" {
  value       = "no-reply@yurinogueira.dev.br"
  description = "Remetente de e-mail verificado no Resend (Secret: EMAIL_FROM)"
}

output "app_base_url" {
  value       = "https://${cloudflare_dns_record.frontend.name}.yurinogueira.dev.br"
  description = "URL base da aplicação frontend (Secret: APP_BASE_URL)"
}




