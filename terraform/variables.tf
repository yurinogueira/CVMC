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

variable "private_key" {
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
