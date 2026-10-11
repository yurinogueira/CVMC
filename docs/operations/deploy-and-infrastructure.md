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
timestamp: 2026-10-10
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

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Pipelines de CI/CD](ci-cd-pipelines.md)
- [Runbook: Release e Rollback](runbooks/release-and-rollback.md)
- [Runbook: Recuperação de Desastres](runbooks/disaster-recovery.md)
