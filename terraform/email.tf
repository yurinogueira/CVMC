# 1. Domínio de envio de e-mail na Oracle OCI
resource "oci_email_email_domain" "cvmc" {
  compartment_id = var.compartment_ocid
  name           = var.email_domain_name
  description    = "Email domain for CVMC notification and auth service"
}

# 2. Chave DKIM para autenticação e alta entregabilidade dos e-mails
resource "oci_email_dkim" "cvmc" {
  email_domain_id = oci_email_email_domain.cvmc.id
  name            = "cvmc"
  description     = "DKIM key for CVMC email domain"
}

# 3. Remetente aprovado (Approved Sender) no OCI Email Delivery
resource "oci_email_sender" "cvmc" {
  compartment_id  = var.compartment_ocid
  email_address   = var.email_sender_address
  email_domain_id = oci_email_email_domain.cvmc.id
}

# 4. Credencial SMTP gerada no IAM da OCI para autenticação segura
resource "oci_identity_smtp_credential" "cvmc" {
  user_id     = var.user_ocid
  description = "SMTP credentials for CVMC backend email service"
}

# 5. Cloudflare DNS Record para DKIM (CNAME sem proxy)
resource "cloudflare_dns_record" "email_dkim" {
  zone_id = var.cloudflare_zone_id
  name    = "${oci_email_dkim.cvmc.dns_subdomain_name}._domainkey"
  content = oci_email_dkim.cvmc.cname_record_value
  type    = "CNAME"
  proxied = false
  ttl     = 1
  comment = "Managed by Terraform - OCI Email Delivery DKIM"
}

# 6. Cloudflare DNS Record para SPF (TXT sem proxy)
resource "cloudflare_dns_record" "email_spf" {
  zone_id = var.cloudflare_zone_id
  name    = "@"
  content = "v=spf1 include:recipientcloud.com ~all"
  type    = "TXT"
  proxied = false
  ttl     = 1
  comment = "Managed by Terraform - OCI Email Delivery SPF"
}
