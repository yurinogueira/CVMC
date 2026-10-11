---
type: frontend
title: "Componentes UI e Design System — CVMC"
description: "Padrões visuais do Material UI v6, paleta análoga de cores, tipografia, regras de padding e acessibilidade (a11y)."
tags:
  - frontend
  - ui
  - mui
  - design-system
  - accessibility
timestamp: 2026-10-10
---

# 🎨 Componentes UI e Design System — CVMC

O frontend do **CVMC** é desenvolvido sobre **Material UI v6** e **React 19**, priorizando consistência estética, alto contraste (WCAG AA/AAA) e alinhamento de layout.

Para navegação geral, retorne ao [Catálogo Canônico](../index.md).

---

## 🎨 Paleta Análoga e Tokens Visuais

- **Cores Principais**:
  - `primary.main`: `#0284C7` (Sky Blue)
  - `secondary.main`: `#0D9488` (Teal)
  - `background.default`: `#F8FAFC` (Slate Claro)
  - `surface.card`: `#FFFFFF` com bordas sutis `#E2E8F0`

---

## 📏 Padrão de Margens e Padding de Página (`AppLayout.tsx`)

Para garantir uniformidade visual entre todas as telas:

1. **Centralização no AppLayout**:
   - O componente `AppLayout.tsx` é a única fonte da verdade para o padding externo:
     ```tsx
     <Box sx={{ flex: 1, p: { xs: 1.5, sm: 2.5, md: 3 }, minWidth: 0, maxWidth: "100%" }}>
       <Box sx={{ maxWidth: 1600, mx: "auto", width: "100%" }}>
         <Outlet />
       </Box>
     </Box>
     ```
2. **Proibição de Padding Duplo**:
   - Os componentes de tela (`*Page.tsx`) **nunca** devem aplicar padding (`p: ...`) em sua tag raiz `<Box>`, prevenindo espaçamento descompensado.

---

## ♿ Acessibilidade e Semântica de Cabeçalhos

- **Hierarquia de Headings**:
  - Cada página contém exatamente um `<h1>` semântico (`component="h1"` com `variant="h4"`).
  - Títulos de seções utilizam `<h2>` e cards utilizam `<h3>`.
- **Prevenção de Erros de Hidratação em Diálogos**:
  - `<DialogTitle>` renderiza uma tag `<h2>` por padrão. Para aninhar `<Typography>` em modais, use explicitamente `component="div"` ou `component="span"`.

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Gerenciamento de Estado com Zustand](state-management.md)
- [Roteamento e RBAC](routing-and-rbac.md)
