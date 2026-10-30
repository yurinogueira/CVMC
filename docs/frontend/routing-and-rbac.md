---
type: frontend
title: "Roteamento e Controle de Acesso (RBAC) — CVMC"
description: "Estrutura de rotas no React Router v7, proteção declarativa com ProtectedRoute e controle de papéis."
tags:
  - frontend
  - routing
  - react-router
  - rbac
timestamp: 2026-10-29
---

# 🚦 Roteamento e Controle de Acesso (RBAC) — CVMC

O sistema de rotas do **CVMC** é orquestrado via **React Router v7**, integrando lazy loading de telas e guardiões de acesso declarativos.

Para navegação geral, retorne ao [Catálogo Canônico](../index.md).

---

## 🗺️ Mapa de Rotas Principais

- **Públicas**:
  - `/login` — Autenticação de condutores.
  - `/register` — Criação de conta.
  - `/forgot-password` e `/reset-password` — Recuperação de senha.
  - `/verify-email` — Confirmação de e-mail.
  - `/privacy` — Página pública de Política de Privacidade (LGPD) e instruções de Exclusão de Conta e Dados (`#exclusao-de-conta`), pré-renderizada em `dist/privacy/index.html` pelo `spaPrerenderPlugin` no `vite.config.ts` para retornar HTTP 200 OK direto no GitHub Pages e atender ao verificador do Google Play Console sem exigir login.
- **Protegidas (`ProtectedRoute`)**:
  - `/dashboard` — Painel consolidado com KPIs, alertas de revisões e resumo de frotas.
  - `/vehicles` — Gestão, listagem e cadastro de veículos.
  - `/vehicles/:id` — Detalhes completos do automóvel, histórico de combustível e manutenções preventivas/corretivas.
  - `/vehicles/:id/maintenance/new` — Formulário de registro de nova manutenção e anexos.
  - `/vehicles/:id/maintenance/:maintenanceId/edit` — Edição e auditoria de ordem de manutenção.
  - `/maintenance` — Histórico global de manutenções e acompanhamento de revisões periódicas.
  - `/reports` — Relatórios e exportação de dossiê consolidado de manutenções em PDF assíncrono.
  - `/profile` — Edição de dados do usuário, limites de frota, quota de storage e segurança de conta.

---

## 🔝 Resolução de Títulos no Topbar (`getTopbarTitle`)

A barra superior (`Topbar.tsx`) sincroniza seu título de navegação dinamicamente de acordo com o `pathname` ativo:
1. **Normalização de Trailing Slashes**: Remove barras finais (ex: `/vehicles/` -> `/vehicles`), prevenindo fallbacks indesejados.
2. **Mapeamento Estático & Dinâmico**: Rotas exatas utilizam o dicionário `staticRouteTitles`, enquanto rotas parametrizadas (ex: `/vehicles/:id`) utilizam correspondência por expressões regulares em `dynamicRouteTitles`.
3. **Harmonia com Sidebar e PageHeader**: Os títulos no Topbar espelham a nomenclatura das abas do menu lateral e dos cabeçalhos `<PageHeader>`.
4. **Governança Automatizada**: Todas as rotas são validadas de forma contínua pelo script `scripts/check-typography.sh` e pela suíte de testes `Topbar.test.tsx`.

---

## 🛡️ Guardião de Rota (`ProtectedRoute`)

O componente intercepta tentativas de acesso não autenticadas, redirecionando o condutor para `/login` enquanto preserva a rota original (`from: location`) para redirecionamento após o login com sucesso.

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Componentes UI e Design System](ui-components.md)
- [Gerenciamento de Estado com Zustand](state-management.md)
- [Autenticação e Segurança em Camadas](../architecture/auth-and-security.md)
