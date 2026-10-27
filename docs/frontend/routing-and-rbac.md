---
type: frontend
title: "Roteamento e Controle de Acesso (RBAC) — CVMC"
description: "Estrutura de rotas no React Router v7, proteção declarativa com ProtectedRoute e controle de papéis."
tags:
  - frontend
  - routing
  - react-router
  - rbac
timestamp: 2026-10-10
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
  - `/dashboard` — Painel geral com KPIs e visão consolidada.
  - `/cars` — Listagem, cadastro e detalhes dos veículos.
  - `/maintenance` — Histórico e agendamento de manutenções.
  - `/fuel` — Histórico e registro de abastecimentos.
  - `/profile` — Edição de dados pessoais e alteração de senha.

---

## 🛡️ Guardião de Rota (`ProtectedRoute`)

O componente intercepta tentativas de acesso não autenticadas, redirecionando o condutor para `/login` enquanto preserva a rota original (`from: location`) para redirecionamento após o login com sucesso.

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Componentes UI e Design System](ui-components.md)
- [Gerenciamento de Estado com Zustand](state-management.md)
- [Autenticação e Segurança em Camadas](../architecture/auth-and-security.md)
