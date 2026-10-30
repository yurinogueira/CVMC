---
type: operations
title: "Deploy e Topologia de Nuvem — CVMC"
description: "Infraestrutura na Oracle Cloud (OCI Always Free), proxy reverso Caddy, Systemd e Cloudflare."
tags:
  - operations
  - deploy
  - oci
  - caddy
  - cloudflare
timestamp: 2026-10-30
---

# 🚀 Deploy e Topologia de Nuvem — CVMC

A infraestrutura produtiva do **CVMC** é hospedada no nível Always Free da **Oracle Cloud Infrastructure (OCI)**.

Para navegação geral, retorne ao [Catálogo Canônico](../index.md).

---

## 🗺️ Topologia e Fluxo de Rede

```mermaid
flowchart LR
    User([Usuário / Navegador]) --> CF[Cloudflare DNS & WAF]
    CF --> Caddy[Caddy Reverse Proxy (Porta 443)]
    subgraph OCI["Instância OCI Compute (Ubuntu ARM/AMD)"]
        Caddy -->|Proxy HTTP :8080| Backend[cvmc-backend.service]
        Backend --> Storage[OCI Object Storage]
    end
    Backend --> MongoAtlas[(MongoDB Atlas)]
```

---

## ⚙️ Componentes de Produção

1. **Serviço Systemd (`cvmc-backend.service`)**:
   - Binário compilado estaticamente em Go (`CGO_ENABLED=0 GOOS=linux GOARCH=amd64`).
   - Gerenciado como serviço de sistema com reinício automático (`Restart=always`).
2. **Caddy Reverse Proxy**:
   - Gerencia terminação TLS automática e roteamento para o backend local.
3. **Deploy Automatizado**:
   - Executado via script `scripts/deploy-backend.sh` e esteira GitHub Actions.

---

## 🧱 Infraestrutura como Código (Terraform)

### Provedores homologados

| Provedor | Versão travada (`.terraform.lock.hcl`) | Constraint (`versions.tf`) |
| :--- | :--- | :--- |
| `oracle/oci` | `9.11.0` | `~> 9.11.0` |
| `mongodb/mongodbatlas` | `2.20.1` | `~> 2.20.0` |
| `cloudflare/cloudflare` | `5.28.1` | `~> 5.28.0` |
| `hashicorp/random` | `3.9.1` | `~> 3.9.0` |

### Protocolo de homologação de upgrades (Dependabot)
1. PRs do grupo `terraform` **não têm auto-merge**; exigem revisão manual.
2. Validar localmente: `terraform -chdir=terraform init -backend=false && terraform -chdir=terraform validate`.
3. Executar **plan somente leitura** com credenciais reais via `workflow_dispatch` do workflow `Terraform` com `apply=false` (o merge feito por bot não dispara o workflow de `push`).
4. Qualquer `forces replacement` deve ser isolado com um plan de controle na versão anterior do provedor para distinguir regressão de drift pré-existente.
5. `apply` apenas via `workflow_dispatch` com `apply=true`, após plan revisado.

### Drifts conhecidos e mitigados
- **`oci_objectstorage_preauthrequest.files_par`**: a API da OCI não retorna `bucket_listing_action` na leitura, gerando `forces replacement` a cada plan (independe da versão do provedor — reproduzido em 9.10.0 e 9.11.0). Mitigado com `lifecycle { ignore_changes = [bucket_listing_action] }` para evitar a rotação da URL do PAR consumida pelo backend (issue #147).
- **`mongodbatlas_cluster`**: recurso depreciado (aviso no `validate`); migração para `mongodbatlas_advanced_cluster` deve ocorrer antes do próximo major do provedor.

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Pipelines de CI/CD](ci-cd-pipelines.md)
- [Runbook: Release e Rollback](runbooks/release-and-rollback.md)
- [Runbook: Recuperação de Desastres](runbooks/disaster-recovery.md)
