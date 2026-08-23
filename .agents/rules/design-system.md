---
description: Diretrizes de UI/UX, paleta de cores análogas, tipografia, componentes, layouts e regras de design para o frontend do CVMC
globs: frontend/**
---

# Regras de Design e UI/UX - CVMC (Como Vai Meu Carro)

Estas regras governam o desenvolvimento visual, componentes de interface, acessibilidade, arquitetura de layouts e padrões estéticos do frontend em React + Material UI v6.

---

## 🎨 1. Paleta de Cores e Harmonia Análoga

O projeto utiliza uma harmonia análoga moderna baseada em azuis, cianos e verdes vibrantes:

| Cor Hex | Nome / Tom | Papel Semântico | Uso no MUI / UI |
| :--- | :--- | :--- | :--- |
| **`#4C92FC`** | Cobalt / Vibrant Blue | **Primary** | Botões de ação principal, links ativos, cabeçalhos de destaque, estados focados. |
| **`#4CCAFC`** | Sky / Cerulean Blue | **Secondary** | Ações secundárias, abas ativas, ícones de destaque, sub-elementos. |
| **`#4CFCF7`** | Electric Cyan | **Accent / Highlight** | Início/fim de gradientes, pontos de atenção sutis, bordas de cards selecionados. |
| **`#4CFCBB`** | Mint Turquoise | **Info / Badge Accent** | Badges de informação, tags de quilometragem/status, detalhes visuais suaves. |
| **`#4CFC7F`** | Spring Green | **Success** | Indicadores de sucesso, status positivo de veículos ("Em dia"), confirmações. |

### Gradiente da Marca
- **Brand Gradient**: `linear-gradient(135deg, #4C92FC 0%, #4CCAFC 50%, #4CFCF7 100%)`
- **Uso do Gradiente**:
  - Banner visual nas telas de autenticação (`AuthHeroBanner.tsx`).
  - Cards de CTA e boas-vindas no Dashboard.
  - Ícone/logo da marca no topo da Sidebar e tela de login.

---

## 🌗 2. Superfícies, Textos e Acessibilidade (WCAG AA)

Para garantir legibilidade perfeita sem ofuscar a visão:
- **Background Principal**: `#F8FAFC` (Slate 50 - neutro muito suave).
- **Cards & Superfícies (Paper)**: `#FFFFFF` com bordas sutis em `#E2E8F0` (Slate 200).
- **Texto Principal**: `#0F172A` (Slate 900 - contraste > 7:1 contra o fundo branco).
- **Texto Secundário / Apoio**: `#64748B` (Slate 500).
- **Texto em Botões Primários (`#4C92FC`)**: `#FFFFFF` (para alto contraste).
- **Texto em Badges Claros (ex.: `#4CFCF7`, `#4CFCBB`, `#4CFC7F`)**: Usar texto escuro `#0F172A` ou `#064E3B` para nunca perder legibilidade.

---

## 📐 3. Tipografia, Espaçamento e Formas

- **Fonte**: `'Inter', system-ui, -apple-system, sans-serif`.
- **Pesos**:
  - `400` (Regular) para textos corridos e inputs.
  - `500` (Medium) para labels, navegação e chips.
  - `600` / `700` / `800` (SemiBold / Bold / ExtraBold) para títulos (`h1`-`h6`), números de KPI e chamadas.
- **Arredondamento de Bordas (Border Radius)**:
  - **Inputs, Botões, Badges, Chips**: `8px` (`borderRadius: 1` ou `8px`).
  - **Cards, Modais, Diálogos, Drawers**: `16px` a `24px` (`borderRadius: 2` a `3`).
- **Sombras (Elevation)**:
  - Cards padrão: `elevation={0}` com `border: 1px solid #E2E8F0` e sombra leve `0 4px 20px -2px rgba(76, 146, 252, 0.06), 0 2px 6px -1px rgba(0, 0, 0, 0.03)`.
  - Hover de Cards interativos: `transform: translateY(-3px)` e sombra `0 10px 25px -3px rgba(76, 146, 252, 0.12)`.

---

## 🏛️ 4. Padrões Estruturais de Layout

### A. Layout de Autenticação (Split-Screen)
- **Desktop (md+)**: Tela dividida em 50%/50%.
  - Lado esquerdo: `AuthHeroBanner` com gradiente da marca, logotipo, título de impacto e 3 pilares de valor com ícones.
  - Lado direito: Formulário centralizado em Card limpo com alternância entre Entrar e Criar Conta.
- **Mobile (xs-sm)**: Banner esquerdo oculto, formulário ocupa 100% com logo compacto no topo.

### B. Shell da Aplicação Autenticada (`AppLayout`)
- **Sidebar (Painel SaaS)**: Largura fixa de `260px` no desktop; gaveta retrátil temporária no mobile.
  - Destaque ativo para o item selecionado (`rgba(76, 146, 252, 0.1)` com indicador lateral).
  - Rodapé da sidebar com avatar, nome, e-mail do usuário e botão de logout.
- **Topbar**: Fixa no topo, com botão de menu no mobile, título dinâmico da rota atual e menu dropdown de perfil com avatar.
- **Área de Conteúdo**: Fundo `#F8FAFC`, padding responsivo (`p: { xs: 2, sm: 3, md: 4 }`).

### C. Telas de Módulos & Dashboard
- **KPI Cards**: Cards métricos no topo com valor numérico em destaque (peso 800), rótulo em caixa alta e ícone semântico em container colorido suave.
- **Empty States**: Quando não há registros (carros ou manutenções), exibir card com borda tracejada (`1px dashed #CBD5E1`), ícone circular centralizado, mensagem explicativa e botão de CTA primário.
- **Modais e Diálogos (`Dialog`)**: Diálogos com `maxWidth="sm"`, `borderRadius: 3`, cabeçalho sem divisórias pesadas e botões de ação bem espaçados no rodapé.

---

## 🧩 5. Padrões de Código e Componentes

### Ícones
- Sempre importar ícones oficiais de `@mui/icons-material` com visual moderno (preferência para sufixos `Rounded` ou `Outlined`, ex: `DirectionsCarRoundedIcon`, `BuildCircleRoundedIcon`).

### Armazenamento e Estado
- **LocalStorage Seguro**: Nunca acesse `window.localStorage` diretamente sem tratamento de erro. Utilize sempre o utilitário `safeStorage` em `src/services/storage/storage.ts` para compatibilidade com SSR e testes automatizados.
- **Zustand Stores**: Estado de autenticação centralizado em `useAuthStore` com reidratação lazy.

---

## 🚫 6. O que NÃO Fazer

- ❌ Não usar texto branco sobre `#4CFCF7` ou `#4CFCBB` (cores claras exigem texto escuro para contraste WCAG).
- ❌ Não utilizar cores fora da harmonia definida (evitar laranjas, roxos ou vermelhos que não sejam para erros críticos `#EF4444`).
- ❌ Não utilizar cantos pontiagudos (0px) ou arredondamento excessivo não padronizado (pill buttons de 50px em tudo).
- ❌ Não colocar sombras pesadas e escuras (`rgba(0,0,0,0.5)`).
- ❌ Não acessar `localStorage` diretamente no corpo de módulos sem o utilitário `safeStorage`.
