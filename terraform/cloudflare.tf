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

# Cloudflare Cache Rules para Frontend (Assets imutáveis e HTML revalidado)
resource "cloudflare_ruleset" "frontend_cache_rules" {
  zone_id     = var.cloudflare_zone_id
  name        = "Frontend Static Cache Rules"
  description = "Aggressive caching for static hashed assets and bypass/revalidate for HTML"
  kind        = "zone"
  phase       = "http_request_cache_settings"

  rules = [
    {
      action      = "set_cache_settings"
      description = "Cache static assets (Vite hashed bundles, images, fonts) for 1 year"
      expression  = "(http.host eq \"cvmc.yurinogueira.dev.br\" and (starts_with(http.request.uri.path, \"/assets/\") or ends_with(http.request.uri.path, \".woff2\") or ends_with(http.request.uri.path, \".png\") or ends_with(http.request.uri.path, \".svg\") or ends_with(http.request.uri.path, \".ico\")))"
      action_parameters = {
        cache = true
        edge_ttl = {
          mode    = "override_origin"
          default = 31536000
        }
        browser_ttl = {
          mode    = "override_origin"
          default = 31536000
        }
        serve_stale = {
          disable_stale_while_updating = false
        }
      }
      enabled = true
    },
    {
      action      = "set_cache_settings"
      description = "Revalidate HTML documents and root"
      expression  = "(http.host eq \"cvmc.yurinogueira.dev.br\" and (http.request.uri.path eq \"/\" or ends_with(http.request.uri.path, \".html\")))"
      action_parameters = {
        cache = true
        edge_ttl = {
          mode    = "override_origin"
          default = 0
        }
        browser_ttl = {
          mode    = "override_origin"
          default = 0
        }
      }
      enabled = true
    }
  ]
}

# Cloudflare Response Headers Transform Rules para Frontend
resource "cloudflare_ruleset" "frontend_response_headers" {
  zone_id     = var.cloudflare_zone_id
  name        = "Frontend Security and Cache Response Headers"
  description = "Set immutable cache headers for hashed static assets"
  kind        = "zone"
  phase       = "http_response_headers_transform"

  rules = [
    {
      action      = "rewrite"
      description = "Set Cache-Control immutable for static assets"
      expression  = "(http.host eq \"cvmc.yurinogueira.dev.br\" and starts_with(http.request.uri.path, \"/assets/\"))"
      action_parameters = {
        headers = {
          "cache-control" = {
            operation = "set"
            value     = "public, max-age=31536000, immutable"
          }
        }
      }
      enabled = true
    },
    {
      action      = "rewrite"
      description = "Set Cache-Control no-cache for HTML pages"
      expression  = "(http.host eq \"cvmc.yurinogueira.dev.br\" and (http.request.uri.path eq \"/\" or ends_with(http.request.uri.path, \".html\")))"
      action_parameters = {
        headers = {
          "cache-control" = {
            operation = "set"
            value     = "public, max-age=0, must-revalidate"
          }
        }
      }
      enabled = true
    }
  ]
}

moved {
  from = cloudflare_record.frontend
  to   = cloudflare_dns_record.frontend
}

moved {
  from = cloudflare_record.backend
  to   = cloudflare_dns_record.backend
}

