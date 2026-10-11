---
type: runbook
title: "Runbook: Recuperação de Desastres (Disaster Recovery) — CVMC"
description: "Restauração de serviços, recuperação de nó OCI e procedimentos de contingência operacional."
tags:
  - runbook
  - operations
  - disaster-recovery
  - failover
timestamp: 2026-10-10
---

# 🚨 Runbook: Recuperação de Desastres — CVMC

Procedimento para restauração emergencial da stack de serviços do **CVMC** em caso de indisponibilidade da instância OCI.

Para navegação geral, retorne ao [Catálogo Canônico](../../index.md).

---

## 🛠️ Procedimento Passo a Passo

1. **Verificação de Saúde**:
   - Teste conectividade SSH e status do serviço:
     ```bash
     ssh $OCI_USER@$OCI_HOST "systemctl status cvmc-backend"
     ```
2. **Reinício de Emergência do Serviço**:
   ```bash
   sudo systemctl restart cvmc-backend
   sudo journalctl -u cvmc-backend -n 50 --no-pager
   ```
3. **Provisionamento de Nova Instância**:
   - Caso a VM esteja corrompida, execute os scripts de provisionamento e provisionamento de DNS:
     ```bash
     ./scripts/provision-server.sh
     ./scripts/deploy-backend.sh
     ```

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../../index.md)
- [Deploy e Infraestrutura](../deploy-and-infrastructure.md)
- [Runbook: Backup de Banco de Dados](database-backup.md)
