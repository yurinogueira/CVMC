# Cloudflare DNS Record para o Frontend (GitHub Pages)
resource "cloudflare_dns_record" "frontend" {
  zone_id = var.cloudflare_zone_id
  name    = "cvmc"
  content = var.github_pages_target
  type    = "CNAME"
  proxied = true
  ttl     = 1
  comment = "Managed by Terraform - CVMC Frontend GitHub Pages"
}

# Cloudflare DNS Record para o Backend (Oracle OCI Reserved Public IP)
resource "cloudflare_dns_record" "backend" {
  zone_id = var.cloudflare_zone_id
  name    = "api-cvmc"
  content = oci_core_public_ip.server_reserved_ip.ip_address
  type    = "A"
  proxied = true
  ttl     = 1
  comment = "Managed by Terraform - CVMC Backend API OCI"
}

moved {
  from = cloudflare_record.frontend
  to   = cloudflare_dns_record.frontend
}

moved {
  from = cloudflare_record.backend
  to   = cloudflare_dns_record.backend
}
