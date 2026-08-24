variable "tenancy_ocid" {
  description = "OCID da tenancy"
  type        = string
}

variable "compartment_ocid" {
  description = "OCID do compartment CVMC onde os recursos serão gerenciados"
  type        = string
}

variable "user_ocid" {
  description = "OCID do usuário OCI usado pelo Terraform"
  type        = string
  sensitive   = true
}

variable "fingerprint" {
  description = "Fingerprint da chave de API OCI"
  type        = string
  sensitive   = true
}

variable "private_key_path" {
  description = "Conteúdo da chave privada de API OCI"
  type        = string
  sensitive   = true
}

variable "region" {
  description = "Região OCI"
  type        = string

  default = "sa-saopaulo-1"
}

variable "availability_domain" {
  description = "Availability Domain"
  type        = string
  default     = ""
}

variable "project_name" {
  description = "Nome do projeto"
  type        = string

  default = "project"
}

variable "ssh_public_key" {
  description = "Chave pública SSH"
  type        = string
}

variable "atlas_org_id" {
  description = "ID da Organização no MongoDB Atlas"
  type        = string
}

variable "atlas_public_key" {
  description = "Public API Key do MongoDB Atlas"
  type        = string
  sensitive   = true
}

variable "atlas_private_key" {
  description = "Private API Key do MongoDB Atlas"
  type        = string
  sensitive   = true
}

variable "mongodb_cluster_name" {
  description = "Nome do cluster MongoDB Atlas"
  type        = string
  default     = "cvmc-cluster"
}

variable "mongodb_app_username" {
  description = "Nome de usuário da aplicação no MongoDB Atlas"
  type        = string
  default     = "cvmc_app"
}

variable "cloudflare_api_token" {
  description = "API Token da Cloudflare com permissão Zone.DNS"
  type        = string
  sensitive   = true
}

variable "cloudflare_zone_id" {
  description = "Zone ID do domínio na Cloudflare (ex: yurinogueira.dev.br)"
  type        = string
}

variable "github_pages_target" {
  description = "Domínio de destino do GitHub Pages (ex: yurinogueira.github.io)"
  type        = string
  default     = "yurinogueira.github.io"
}

