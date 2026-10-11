---
type: log
title: "Diário de Bordo & Trilha de Mudanças — CVMC"
description: "Índice cronológico de intervenções de engenharia, refatorações, novas features e manutenções do CVMC."
tags:
  - log
  - audit
  - history
timestamp: 2026-10-10
---

# 🕒 Diário de Bordo & Histórico Cronológico — CVMC

Este documento funciona como o índice mestre de todas as mudanças de engenharia aplicadas ao projeto **CVMC (Como Vai Meu Carro)**. Cada entrada detalhada reside em um arquivo diário específico na pasta [docs/logs/2026-10-10.md](logs/2026-10-10.md), garantindo economia de tokens e histórico auditável.

Para navegação geral, retorne ao [Catálogo Canônico](index.md).

---

## 📅 [2026-10-10](logs/2026-10-10.md) — Correção de Layout Mobile e Responsividade Global
- **Frontend & Layout**: Resolução do overflow horizontal na tela de detalhes do veículo (`VehicleDetailsPage.tsx`) tornando abas roláveis com `variant="scrollable"`, aplicando confinamento flexbox (`minWidth: 0`, `maxWidth: "100%"`, `overflowX: "hidden"`) em `AppLayout.tsx` e resiliência na `Topbar.tsx`.
- **Formulários & Modais**: Empilhamento responsivo de botões de ação em `RegisterMaintenancePage.tsx` e margens seguras para dispositivos móveis em `AddFuelingDialog.tsx` e modais de exclusão.
- **Testes & Documentação**: Teste unitário para abas scrollable com MUI v6 e atualização da especificação de responsividade em `docs/frontend/ui-components.md`.

## 📅 [2026-10-10](logs/2026-10-10.md) — Sincronização Canônica e Modernização (PS -> CVMC)
- **Documentação**: Criação da árvore canônica de documentação Open Knowledge Format (OKF) em `docs/` com mais de 25 documentos abrangendo arquitetura, domínio, frontend, operações, runbooks e ADRs.
- **Governança & Automação**: Criação de `scripts/check-docs.sh`, inclusão do alvo `docs` em `scripts/check.sh`, novo workflow `.github/workflows/docs.yml` e modernização do template de PR.
- **Skills & Regras**: Criação da skill `cvmc-docs`, inclusão da regra `docs.md`, atualização de `design-system.md` (margens e modais) e regras defensivas em `security.md`.
- **Backend Defensivo**: Implementação do método `Config.Validate()` com erros tipados para validação rigorosa de entropia e não repetição de segredos JWT.
