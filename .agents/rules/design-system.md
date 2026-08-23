---
description: Diretrizes de UI/UX, paleta de cores análogas, tipografia, componentes, layouts e regras de design para o frontend do CVMC
globs: frontend/**
---

# Regras de Design e UI/UX - CVMC (Como Vai Meu Carro)

Estas regras governam o desenvolvimento visual, componentes de interface, acessibilidade (WCAG AA/AAA), arquitetura de layouts, hierarquia de headings semânticos, performance de carregamento crítico (LCP/FCP) e padrões estéticos do frontend em React + Material UI v6.

---

## 🎨 1. Paleta de Cores, Harmonia Análoga e Alto Contraste (WCAG AA)

O projeto utiliza uma harmonia análoga moderna em tons de azul, ciano e verde, calibrada para cumprir os requisitos de contraste WCAG AA/AAA em todos os estados:

| Cor Hex | Nome / Tom | Papel Semântico | Uso no MUI / UI | Contraste Mínimo |
| :--- | :--- | :--- | :--- | :--- |
| **`#0284C7`** | Sky 600 | **Primary Main** | Botões primários, links principais, itens de destaque, ícones de foco. | **4.7:1** (c/ `#FFFFFF`) |
| **`#0369A1`** | Sky 700 | **Primary Dark** | Hover de botões, texto em botões brancos sobre gradiente, links em destaque. | **6.5:1** (c/ `#FFFFFF`) |
| **`#0EA5E9`** | Sky 500 | **Secondary** | Ações secundárias, abas ativas, sub-elementos. | **3.8:1** (elementos grandes) |
| **`#06B6D4`** | Cyan 500 | **Accent / Highlight** | Início/fim de gradientes sutis, bordas de cards e acentos visuais. | Acento visual |
| **`#16A34A`** | Green 600 | **Success** | Indicadores de sucesso, status positivo ("Em dia"), confirmações. | **4.5:1** (c/ `#FFFFFF`) |
| **`#0D9488`** | Teal 600 | **Info** | Badges e detalhes informativos. | **4.5:1** (c/ `#FFFFFF`) |

### Gradiente da Marca
- **Brand Gradient**: `linear-gradient(135deg, #0284C7 0%, #0369A1 50%, #0F766E 100%)`
- **Uso do Gradiente**:
  - Banner visual nas telas de autenticação (`AuthHeroBanner.tsx`).
  - Cards de CTA e boas-vindas no Dashboard.
  - Ícone/logo da marca no topo da Sidebar e tela de login.

---

## 🌗 2. Superfícies, Textos e Regras de Contraste

Para garantir legibilidade absoluta em qualquer dispositivo e ambiente de iluminação:
- **Background Principal**: `#F8FAFC` (Slate 50 - neutro muito suave).
- **Cards & Superfícies (Paper)**: `#FFFFFF` com bordas sutis em `#E2E8F0` (Slate 200).
- **Texto Principal**: `#0F172A` (Slate 900 - contraste > 15:1 contra o fundo branco).
- **Texto Secundário / Apoio**: `#475569` (Slate 600 - contraste **5.6:1** contra fundo branco, superando os 4.5:1 do WCAG AA).
- **Botões sobre Fundo Branco (ex.: dentro de CTA com gradiente)**:
  - Fundo: `#FFFFFF`
  - Texto: `#0369A1` (Sky 700) com peso `700`, garantindo **6.5:1** de contraste.
- **Texto em Badges Claros (ex.: `#A7F3D0`, `#86EFAC`)**:
  - Usar texto escuro `#064E3B` ou `#0F172A` (contraste > 7:1).

---

## 🧭 3. Hierarquia Semântica de Headings (A11y & SEO)

Os cabeçalhos devem ser organizados de forma sequencial e descendente, sem pular níveis de hierarquia:

1. **`h1` (Título Principal da Página)**:
   - Toda página deve possuir **exatamente um** elemento semântico `h1`.
   - Use `component="h1"` com `variant="h4"` ou `variant="h5"` (ex.: *"Olá, [Nome] 👋"*, *"Meus Veículos"*, *"Bem-vindo de volta"*).
2. **`h2` (Seções e Banners Principais)**:
   - Use `component="h2"` para banners de ação rápida (CTA), títulos de seções (ex.: *"Meus Veículos"*) e cartões de status da frota.
3. **`h3` (Itens de Grade, KPIs e Estados)**:
   - Use `component="h3"` para títulos de KPI cards (*"Veículos Cadastrados"*), títulos de cartões de itens individuais (`VehicleCard`), títulos de modais e títulos de Empty States.
4. **Banners Decorativos / Ilustrativos**:
   - Elementos em colunas puramente visuais (como `AuthHeroBanner`) devem utilizar `component="div"` ou `component="p"` para não poluir a árvore semântica do leitor de tela nem colidir com o `h1` do formulário.

---

## ⚡ 4. Performance, LCP e Carregamento Crítico

1. **Skeleton Loaders Imediatos (Zero Layout Shift)**:
   - **Nunca** utilize spinners genéricos de tela inteira (`CircularProgress`) que ocultam a estrutura da página enquanto a API responde.
   - Renderize o cabeçalho, banner e a grade de KPIs instantaneamente no primeiro frame, preenchendo valores assíncronos com `<Skeleton width={...} height={...} />`.
   - Isso antecipa o **Largest Contentful Paint (LCP)** para menos de 1 segundo e zera o **Cumulative Layout Shift (CLS)**.
2. **Lazy Loading de Modais e Diálogos Complexos**:
   - Componentes pesados (formulários de cadastro, modais de confirmação, seletores de data) devem ser importados com `lazy(() => import(...))` e renderizados apenas quando abertos (`open && <Suspense fallback={null}><Modal /></Suspense>`).

---

## 🔍 5. Metadados de SEO, Canonical e Indexação

1. **Sincronização de Título e Canonical URL**:
   - Utilize o hook `useDocumentTitle("Nome da Página")` em todas as rotas para atualizar o `document.title` e injetar a tag `<link rel="canonical" href="https://cvmc.yurinogueira.dev.br/rota/" />` dinamicamente.
2. **Pré-renderização Estática de Rotas no Vite**:
   - O plugin `spaPrerenderPlugin` em `vite.config.ts` deve gerar subdiretórios com `index.html` e a tag `canonical` exata para cada rota (`/dashboard/`, `/vehicles/`, `/maintenance/`, `/login/`, `/register/`), garantindo resposta nativa `HTTP 200 OK` em hosts estáticos (GitHub Pages, Cloudflare Pages).
3. **Políticas de Crawlers (`robots.txt`)**:
   - Mantenha `Allow: /` para todas as páginas de usuário para permitir auditorias de SEO completas (Lighthouse) e indexação por mecanismos de busca, mantendo bloqueados apenas endpoints privados de backend (`/api/`).

---

## 📐 6. Tipografia, Espaçamento e Formas

- **Fonte**: `'Inter', system-ui, -apple-system, sans-serif`.
- **Pesos**:
  - `400` (Regular) para textos corridos e inputs.
  - `500` (Medium) para labels, navegação e chips.
  - `600` / `700` / `800` (SemiBold / Bold / ExtraBold) para títulos (`h1`-`h6`), números de KPI e chamadas de ação.
- **Arredondamento de Bordas (Border Radius)**:
  - **Inputs, Botões, Badges, Chips**: `8px` (`borderRadius: 1` ou `8px`).
  - **Cards, Modais, Diálogos, Drawers**: `16px` a `24px` (`borderRadius: 2` a `3`).
- **Sombras (Elevation)**:
  - Cards padrão: `elevation={0}` com `border: 1px solid #E2E8F0` e sombra leve `0 4px 20px -2px rgba(2, 132, 199, 0.06), 0 2px 6px -1px rgba(0, 0, 0, 0.03)`.
  - Hover de Cards interativos: `transform: translateY(-3px)` e sombra `0 10px 25px -3px rgba(2, 132, 199, 0.12)`.

---

## 🚫 7. O que NÃO Fazer

- ❌ Não usar texto branco sobre cores claras como `#4CFCF7` ou `#A7F3D0` (exigem texto escuro `#064E3B` ou `#0F172A`).
- ❌ Não usar texto azul claro sobre fundo branco em botões (usar `#0369A1` ou `#0F172A` para manter contraste > 4.5:1).
- ❌ Não pular níveis de heading (ex.: ir direto de `h1` para `h6` ou `h4` para `h6`).
- ❌ Não bloquear a renderização inicial com loaders centrados que adiam o LCP.
- ❌ Não colocar diretivas `Disallow` no `robots.txt` para páginas que precisam ser avaliadas ou indexadas pelo Lighthouse/Googlebot.
- ❌ Não acessar `localStorage` diretamente sem o utilitário `safeStorage`.
