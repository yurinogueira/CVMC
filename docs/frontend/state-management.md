---
type: frontend
title: "Gerenciamento de Estado com Zustand — CVMC"
description: "Stores globais com Zustand, sincronização de perfil e proteção contra persistência de tokens sensíveis no cliente."
tags:
  - frontend
  - state-management
  - zustand
  - store
timestamp: 2026-10-10
---

# 📦 Gerenciamento de Estado com Zustand — CVMC

A gestão de estado global no **CVMC** utiliza **Zustand**, oferecendo reatividade sem a sobrecarga e verbosidade de bibliotecas legadas.

Para navegação geral, retorne ao [Catálogo Canônico](../index.md).

---

## 🏛️ Stores Globais

1. **`authStore` (`frontend/src/features/auth/state/auth.store.ts`)**:
   - Mantém o estado de login (`isAuthenticated`), prontidão de inicialização (`isInitializing`), usuário ativo (`user`) e status de carregamento.
   - **Zero Tokens no Web**: No navegador web, nunca guarda tokens JWT em `localStorage`. As credenciais trafegam por cookies `HttpOnly` gerenciados pelo navegador.
   - **Persistência Mobile Nativa**: No Android/Capacitor, `initAuthSession` restaura de forma assíncrona o usuário e os tokens (`accessToken` e `refreshToken`) persistidos no `@capacitor/preferences`. Enquanto `isInitializing` for verdadeiro, o `ProtectedRoute` exibe o fallback de loading e a SplashScreen nativa permanece ativa para evitar redirecionamentos indesejados.

---

## 🔄 Boas Práticas

- **Seletores Granulares**: Sempre consuma propriedades da store utilizando seletores específicos (ex: `useAuthStore(state => state.user)`) para evitar re-renderizações desnecessárias.
- **Ações Desacopladas**: Modificações de estado residem como métodos dentro da própria store.

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Componentes UI e Design System](ui-components.md)
- [Roteamento e RBAC](routing-and-rbac.md)
- [ADR 0003: Zustand para Estado Global](../adrs/0003-zustand-state-management.md)
