---
type: runbook
title: "Runbook: Release e Rollback em Produção — CVMC"
description: "Protocolo de implantação de versões e procedimentos de reversão imediata (rollback) em caso de falha."
tags:
  - runbook
  - operations
  - release
  - rollback
timestamp: 2026-10-10
---

# 🔄 Runbook: Release e Rollback em Produção — CVMC

Procedimentos padrão para deploy de novas versões e recuperação por reversão imediata no **CVMC**.

Para navegação geral, retorne ao [Catálogo Canônico](../../index.md).

---

## 🚀 1. Procedimento de Release

1. Certifique-se de que a validação passou 100%: `./scripts/check.sh all`.
2. O merge na branch `main` dispara o workflow do GitHub Actions `.github/workflows/backend.yml`.
3. Alternativamente, execute localmente: `./scripts/deploy-backend.sh`.

---

## ⏪ 2. Procedimento de Rollback Imediato

1. **Reversão via Git**:
   ```bash
   git revert HEAD -m 1
   git push origin main
   ```
2. **Reversão Manual no Servidor**:
   - Se o deploy falhou e a versão anterior do binário estiver preservada em `/opt/cvmc/cvmc-api.bak`:
     ```bash
     sudo cp /opt/cvmc/cvmc-api.bak /opt/cvmc/cvmc-api
     sudo systemctl restart cvmc-backend
     ```

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../../index.md)
- [Deploy e Infraestrutura](../deploy-and-infrastructure.md)
- [Pipelines de CI/CD](../ci-cd-pipelines.md)
