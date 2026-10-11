---
type: operations
title: "Pipelines de CI/CD — CVMC"
description: "Workflows do GitHub Actions para validação do backend, frontend, terraform e governança de documentação."
tags:
  - operations
  - ci-cd
  - github-actions
  - automation
timestamp: 2026-10-10
---

# 🤖 Pipelines de CI/CD — CVMC

A esteira de integração contínua do **CVMC** garante qualidade de código, segurança e conformidade arquitetural em cada Pull Request direcionado à branch `main`.

Para navegação geral, retorne ao [Catálogo Canônico](../index.md).

---

## 🏗️ Workflows Ativos

1. **`backend.yml`**:
   - Executa `go vet` e testes automatizados (`go test ./...`).
   - Compila o binário estático e opcionalmente dispara deploy no OCI.
2. **`frontend.yml`**:
   - Executa checagem de tipos (`tsc`), linter ESLint, formatação Prettier e testes Vitest.
   - Gera build de produção e publica no GitHub Pages / host estático.
3. **`terraform.yml`**:
   - Valida a formatação de infraestrutura (`terraform fmt -check`) e sintaxe HCL.
4. **`docs.yml` (OKF Quality Gate)**:
   - Valida frontmatter YAML, links locais quebrados e conformidade Zero-Divergence via `scripts/check-docs.sh`.

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Deploy e Infraestrutura](deploy-and-infrastructure.md)
- [Docker e Ambiente Local](docker-and-local-dev.md)
