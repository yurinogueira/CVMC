---
type: log
title: "Diário de Bordo & Trilha de Mudanças — CVMC"
description: "Índice cronológico de intervenções de engenharia, refatorações, novas features e manutenções do CVMC."
tags:
  - log
  - audit
  - history
timestamp: 2026-10-11
---

# 🕒 Diário de Bordo & Histórico Cronológico — CVMC

Este documento funciona como o índice mestre de todas as mudanças de engenharia aplicadas ao projeto **CVMC (Como Vai Meu Carro)**. Cada entrada detalhada reside em um arquivo diário específico na pasta [docs/logs/2026-10-10.md](logs/2026-10-10.md), garantindo economia de tokens e histórico auditável.

Para navegação geral, retorne ao [Catálogo Canônico](index.md).

## 📅 [2026-10-11](logs/2026-10-11.md) — Empacotamento Android com Capacitor e Configuração de Build (#92)
- **Capacitor 8 no Frontend**: Configuração do `capacitor.config.ts` com `appId: com.cvmc.app`, inicialização do projeto `frontend/android` e scripts `cap:sync`, `cap:open` e `cap:build`.
- **Comunicação de API & Sessão Mobile**: BaseURL dinâmica para ambiente nativo, injeção de `Authorization: Bearer <token>`, persistência segura com `@capacitor/preferences` e retorno de tokens em JSON nos endpoints de autenticação em Go.
- **UX Mobile Nativa**: Controle da cor da barra de status (`#0284C7`), splash screen, hardware back button do Android e suporte a Safe Area Insets.
- **Correção de CORS e Tratamento de Erros**: Liberação das origens `https://localhost` e `capacitor://localhost` no backend Go e no deploy OCI (`backend.yml`), e refinamento no `LoginPage.tsx` para distinguir erro 401 de falha de conexão.
- **Status Bar Fixa e Preta**: Desativação de sobreposição de tela no Android (`overlay: false`) e definição de fundo preto sólido (`#000000`) com ícones claros para legibilidade total.
- **Sobreposição de Sidebar e Contraste Mobile (#104)**: Aplicação de safe area insets (`env(safe-area-inset-top)` e `env(safe-area-inset-bottom)`) no Drawer da `Sidebar.tsx`, eliminação de colisão com a barra de status do celular e configuração de alto contraste (`Style.Light`, `android:windowLightStatusBar="true"`, `#FFFFFF`) para o topo da aplicação.
- **Automação de Build do APK**: Novo job `build-android-apk` no workflow `frontend.yml` com injeção automática de `VITE_GOOGLE_CLIENT_ID`, build Gradle e disponibilização do artefato `app-debug-apk` no GitHub Actions.
- **Documentação de Build**: Criação do runbook operacional [docs/operations/mobile-build.md](operations/mobile-build.md) detalhando compilação do APK de debug, geração de bundle AAB e assinatura com keystore para publicação na Google Play.

## 📅 [2026-10-10](logs/2026-10-10.md) — Suporte a Login Social Google e Vínculo de Contas Existentes (#91)
- **Autenticação Direta Google Cloud**: Implementação de autenticação social com Google Identity Services (GIS) no frontend e validação nativa de ID tokens com `google.golang.org/api/idtoken` no backend Go, dispensando intermediários (sem Firebase).
- **Auto-linking e Gestão de Provedores**: Associação automática de contas existentes por e-mail, prevenção de lockout (desvinculação proibida sem senha cadastrada) e suporte a definição de senha inicial via `POST /api/v1/user/password/set`.
- **Interface e Experiência**: Botão GIS oficial em login/cadastro, card de contas conectadas na página de perfil e alternância dinâmica entre definição e troca de senha.

## 📅 [2026-10-10](logs/2026-10-10.md) — Correção de Layout Mobile e Responsividade Global
- **Frontend & Layout**: Resolução do overflow horizontal na tela de detalhes do veículo (`VehicleDetailsPage.tsx`) tornando abas roláveis com `variant="scrollable"`, aplicando confinamento flexbox (`minWidth: 0`, `maxWidth: "100%"`, `overflowX: "hidden"`) em `AppLayout.tsx` e resiliência na `Topbar.tsx`.
- **Formulários & Modais**: Empilhamento responsivo de botões de ação em `RegisterMaintenancePage.tsx` e margens seguras para dispositivos móveis em `AddFuelingDialog.tsx` e modais de exclusão.
- **Testes & Documentação**: Teste unitário para abas scrollable com MUI v6 e atualização da especificação de responsividade em `docs/frontend/ui-components.md`.

## 📅 [2026-10-10](logs/2026-10-10.md) — Sincronização Canônica e Modernização (PS -> CVMC)
- **Documentação**: Criação da árvore canônica de documentação Open Knowledge Format (OKF) em `docs/` com mais de 25 documentos abrangendo arquitetura, domínio, frontend, operações, runbooks e ADRs.
- **Governança & Automação**: Criação de `scripts/check-docs.sh`, inclusão do alvo `docs` em `scripts/check.sh`, novo workflow `.github/workflows/docs.yml` e modernização do template de PR.
- **Skills & Regras**: Criação da skill `cvmc-docs`, inclusão da regra `docs.md`, atualização de `design-system.md` (margens e modais) e regras defensivas em `security.md`.
- **Backend Defensivo**: Implementação do método `Config.Validate()` com erros tipados para validação rigorosa de entropia e não repetição de segredos JWT.
