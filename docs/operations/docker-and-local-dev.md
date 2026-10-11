---
type: operations
title: "Docker e Ambiente de Desenvolvimento Local — CVMC"
description: "Composição de containers com Docker Compose, variáveis locais e scripts de inicialização rápida."
tags:
  - operations
  - docker
  - docker-compose
  - local-dev
timestamp: 2026-10-10
---

# 🐳 Docker e Ambiente de Desenvolvimento Local — CVMC

O ambiente local do **CVMC** é totalmente reproduzível através do Docker Compose, isolando dependências de banco de dados e simuladores.

Para navegação geral, retorne ao [Catálogo Canônico](../index.md).

---

## 🚀 Inicialização Rápida

O script `./scripts/dev.sh` facilita as rotinas diárias:

```bash
# Iniciar todos os serviços em segundo plano:
./scripts/dev.sh start

# Verificar status dos containers e portas:
./scripts/dev.sh status

# Visualizar logs consolidados:
./scripts/dev.sh logs

# Parar a stack:
./scripts/dev.sh stop
```

---

## 🧱 Serviços da Stack Local

1. **`backend`**: API REST em Go (porta interna 8080).
2. **`frontend`**: Aplicação React com Vite Hot Reload (porta 5173).
3. **`mongodb`**: Instância local de desenvolvimento do MongoDB (rede interna de containers sem exposição de portas inseguras ao host).

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Ambiente e Configuração](environment-and-config.md)
- [Pipelines de CI/CD](ci-cd-pipelines.md)
