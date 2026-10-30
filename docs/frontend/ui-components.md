---
type: frontend
title: "Componentes UI e Design System — CVMC"
description: "Padrões visuais do Material UI v9, paleta análoga de cores, tipografia, regras de padding e acessibilidade (a11y)."
tags:
  - frontend
  - ui
  - mui
  - design-system
  - accessibility
timestamp: 2026-10-29
---

# 🎨 Componentes UI e Design System — CVMC

O frontend do **CVMC** é desenvolvido sobre **Material UI v9** e **React 19**, priorizando consistência estética, alto contraste (WCAG AA/AAA) e alinhamento de layout.

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

## 📱 Responsividade Mobile e Confinamento Flexbox

Para prevenir estouro horizontal (*horizontal overflow*) e garantir adaptação impecável em telas móveis:

1. **Confinamento em Itens Flexíveis**:
   - Elementos principais de layout (`Box component="main"`) devem declarar explicitamente `minWidth: 0`, `maxWidth: "100%"` e `overflowX: "hidden"`. Por padrão da especificação CSS Flexbox, itens possuem `min-width: auto`, o que permite que elementos filhos com conteúdo amplo estiquem o layout além de 100vw.
2. **Abas Roláveis no Mobile (`<Tabs>`)**:
   - Em abas de navegação internas com textos longos ou chips de contagem (como em `VehicleDetailsPage.tsx`), deve-se utilizar obrigatoriamente:
     ```tsx
     <Tabs
       variant="scrollable"
       scrollButtons="auto"
       allowScrollButtonsMobile
       sx={{
         "& .MuiTabs-scroller": { overflowX: "auto" },
         "& .MuiTab-root": { minWidth: { xs: "auto", sm: 160 }, px: { xs: 1.5, sm: 2 } },
       }}
     >
     ```
3. **Barra Superior Resiliente (`Topbar.tsx`)**:
   - O título deve possuir `noWrap` com `textOverflow: "ellipsis"`, e o bloco de avatar/ações deve declarar `flexShrink: 0`, impedindo que títulos longos empurrem o perfil para fora da tela.
4. **Empilhamento de Ações em Telas Estreitas**:
   - Cabeçalhos de seção e formulários devem utilizar `Stack direction={{ xs: "column", sm: "row" }}` com botões de largura total (`width: { xs: "100%", sm: "auto" }`) em resoluções móveis.
5. **Alinhamento e Centralização de Botões Sociais (`GoogleSignInButton.tsx`)**:
   - Botões de provedores externos renderizados via SDKs de terceiros (como o Google Identity Services / GIS) devem utilizar `logo_alignment: "center"` e container flexível centralizado (`justifyContent: "center"`), garantindo agrupamento simétrico do logotipo e texto com centralização natural e harmônica com os formulários de autenticação.
6. **Safe Area Insets no Menu Lateral Móvel (`Sidebar.tsx`)**:
   - O Drawer temporário para dispositivos móveis deve obrigatoriamente declarar `pt: "env(safe-area-inset-top, 0px)"` e `pb: "env(safe-area-inset-bottom, 0px)"` em seu container (`& .MuiDrawer-paper`). Como o Drawer é renderizado em um portal com `position: fixed; top: 0`, ele não herda o padding do `body`. A aplicação de safe area insets previne que o logotipo e cabeçalho colidam com a barra de status do sistema (relógio, bateria, etc.) e que a seção de usuário encavale na barra de gestos do Android/iOS.
7. **Alto Contraste da Barra de Status Nativa**:
   - Em aplicações com tema visual claro (`mode: "light"`), a barra de status do celular deve operar com ícones e textos escuros (`Style.Light` no Capacitor, `android:windowLightStatusBar="true"` no Android nativo e `<meta name="theme-color" content="#FFFFFF">` no HTML), assegurando contraste nítido e legibilidade perfeita do relógio e indicadores sobre os cabeçalhos claros do app.
8. **Grid Responsivo em Resoluções Intermediárias e Monitores Verticais (`MaintenancePage.tsx`)**:
   - Em grids com alta densidade de métricas e ações (como os cards de `UpcomingMaintenanceCard.tsx`), a presença do drawer lateral persistente de 260px (`DRAWER_WIDTH`) a partir de `md` (900px) reduz a largura útil do container para ~690px em resoluções como 1000x1800.
   - Para evitar achatamento visual, truncamento de badges de status e quebras descompensadas de texto, o grid preventivo adota a proporção adaptativa `size={{ xs: 12, sm: 6, md: 6, lg: 4 }}` (renderizando 2 colunas no breakpoint `md` e reservando 3 colunas para larguras amplas `lg` >= 1200px), associado a confinamento interno com `minWidth: 0` e `overflowWrap: "break-word"`.

---

## ♿ Acessibilidade e Semântica de Cabeçalhos

- **Hierarquia de Headings e Contrato Mandatório de Cabeçalho (`PageHeader`)**:
  - Toda página autenticada da aplicação deve obrigatoriamente utilizar o componente padronizado `<PageHeader>` (`frontend/src/features/shared/components/PageHeader.tsx`).
  - **Tag `<h1>` Canônica Única**: Cada tela contém exatamente uma tag `<h1>` semântica para conformidade com WCAG AA/AAA e acessibilidade por leitores de tela. O cabeçalho é renderizado com `component="h1"`, `variant="h5"`, `fontWeight: 800` e tamanho de fonte responsivo (`fontSize: { xs: "1.25rem", sm: "1.5rem" }`).
  - **Proibição de Declarações Ad-Hoc**: É expressamente proibido declarar `<Typography component="h1">` manual ou títulos de tela com `variant="h4"` / pesos arbitrários diretamente em arquivos de página (`*Page.tsx`).
  - **Subtítulo Padronizado**: Subtítulos de contextualização utilizam `variant="body2"` e cor `text.secondary` com espaçamento uniforme (`mt: 0.5`).
  - **Ações e Botões de Voltar Responsivos**: O slot `action` alinha botões à direita no desktop e os empilha com largura adaptativa no mobile. O botão voltar é padronizado pelo objeto `backButton={{ onClick, label }}` com bordas sutis e contraste acessível.
  - **Estado de Carregamento sem Layout Shift**: Ao carregar entidades (como dados do veículo), deve-se passar `loading={true}` para o `PageHeader`, renderizando o Skeleton correspondente sem quebra de alinhamento visual.
  - **Validação Automatizada Contínua**: O cumprimento dessas diretrizes visuais é verificado de forma automatizada pelo script `scripts/check-typography.sh` e pela suíte arquitetural Vitest (`TypographyHierarchy.test.tsx`), integrados ao comando `./scripts/check.sh frontend` e às esteiras de CI. Qualquer cabeçalho fora do padrão interrompe o build.
- **Prevenção de Erros de Hidratação em Diálogos**:
  - `<DialogTitle>` renderiza uma tag `<h2>` por padrão. Para aninhar `<Typography>` em modais, use explicitamente `component="div"` ou `component="span"`.
- **Cards de Próximas Manutenções (`UpcomingMaintenanceCard.tsx` e `MaintenancePage.tsx`)**:
  - Exibem indicadores semânticos de status (`Em dia` em verde, `Próximo` em âmbar e `Vencido` em vermelho), barra de progresso (`LinearProgress`) proporcional ao desgaste por quilometragem/tempo e botão de ação rápida que abre o `AddMaintenanceDialog` com o tipo de manutenção e a quilometragem atual do veículo pré-selecionados.

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Gerenciamento de Estado com Zustand](state-management.md)
- [Roteamento e RBAC](routing-and-rbac.md)
