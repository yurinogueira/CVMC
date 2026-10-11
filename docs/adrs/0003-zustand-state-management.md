---
type: adr
title: "ADR 0003: Adoção do Zustand para Gerenciamento de Estado no Frontend"
description: "Decisão de utilizar Zustand para estado reativo global no cliente, eliminando boilerplate de Redux."
tags:
  - adr
  - frontend
  - react
  - zustand
timestamp: 2026-10-10
---

# 📜 ADR 0003: Adoção do Zustand para Gerenciamento de Estado

- **Status**: Aceito
- **Data**: 2026-10-10
- **Decisores**: Equipe de Engenharia CVMC

Para navegação geral, retorne ao [Catálogo Canônico](../index.md).

---

## 📌 Contexto
Precisávamos de um gerenciador de estado leve, com baixo consumo de memória e simples de tipar em TypeScript para gerenciar o perfil do usuário e preferências de tela.

---

## 💡 Decisão
Adotamos **Zustand** como padrão único de gerenciamento de estado global no frontend:

1. Elimina reducers complexos, dispatchers e actions verbosas do Redux.
2. Suporta hooks nativos e seletores granulares, minimizando re-renderizações no React.
3. Não requer providers envolvendo a árvore de componentes raiz da aplicação.

---

## ⚖️ Consequências

### Positivas:
- Menor tamanho de bundle no Vite e inicialização instantânea da aplicação.
- Código direto, intuitivo e de fácil manutenção para desenvolvedores e agentes autônomos.

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Gerenciamento de Estado com Zustand](../frontend/state-management.md)
- [Componentes UI e Design System](../frontend/ui-components.md)
