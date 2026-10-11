---
type: runbook
title: "Runbook: Backup e Restauração de Banco de Dados — CVMC"
description: "Procedimentos operacionais de cópia de segurança e restauração no MongoDB local e Atlas."
tags:
  - runbook
  - operations
  - backup
  - restore
  - mongodb
timestamp: 2026-10-10
---

# 💾 Runbook: Backup e Restauração de Banco de Dados — CVMC

Procedimento para cópias de segurança e restauração de dados no **CVMC**.

Para navegação geral, retorne ao [Catálogo Canônico](../../index.md).

---

## 📦 1. Geração de Backup (`mongodump`)

```bash
# Backup local:
mongodump --uri="mongodb://localhost:27017/cvmc" --out=/backup/cvmc-$(date +%F)

# Backup comprimido em arquivo único:
mongodump --uri="$MONGO_URI" --archive=/backup/cvmc-backup.gz --gzip
```

---

## ♻️ 2. Restauração (`mongorestore`)

```bash
# Restaurar a partir de arquivo comprimido:
mongorestore --uri="$MONGO_URI" --archive=/backup/cvmc-backup.gz --gzip --drop
```

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../../index.md)
- [Runbook: Recuperação de Desastres](disaster-recovery.md)
- [Isolamento de Dados e Persistência](../../architecture/multitenancy-and-data.md)
