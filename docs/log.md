---
type: log
title: "Diário de Bordo & Trilha de Mudanças — CVMC"
description: "Índice cronológico de intervenções de engenharia, refatorações, novas features e manutenções do CVMC."
tags:
  - log
  - audit
  - history
timestamp: 2026-10-24
---

# 🕒 Diário de Bordo & Histórico Cronológico — CVMC

Este documento funciona como o índice mestre de todas as mudanças de engenharia aplicadas ao projeto **CVMC (Como Vai Meu Carro)**. Cada entrada detalhada reside em um arquivo diário específico na pasta [docs/logs/2026-10-24.md](logs/2026-10-24.md), garantindo economia de tokens e histórico auditável.

Para navegação geral, retorne ao [Catálogo Canônico](index.md).

## 📅 [2026-10-24](logs/2026-10-24.md) — Monitoramento de Revisões Periódicas Próximas (#120), Ícones da Play Store (#122) & Migração para Oxlint + TypeScript 7
- **Migração de Linter para Oxlint + Atualização para TypeScript 7.0.2**: Substituição de `eslint` e `typescript-eslint` pelo `oxlint` (`^1.89.0`) com plugins `typescript`, `unicorn`, `oxc` e `react` em modo estrito (`--deny-warnings`), desbloqueando a atualização do compilador para `typescript@7.0.2` e aplicando todas as atualizações de bibliotecas do frontend (`@mui/material@9.5.0`, `react@19.3.1`, `vitest@5.0.5`, `zustand@5.0.16`, `@capgo/capacitor-social-login@8.7.5`).
- **Ícones Adaptativos Android, Splash Screens e Assets da Google Play Store (#122)**: Substituição de todos os ícones `mipmap-*` (`ic_launcher`, `ic_launcher_round`, `ic_launcher_foreground`), suporte a Adaptive Icons e `<monochrome>` (Android 13+), criação do ícone de notificação push FCM (`ic_stat_cvmc.xml`), geração das 11 telas de Splash Screen (`#0284C7`), criação dos ativos do Google Play Console (`playstore-icon-512.png` e `feature-graphic-1024x500.png`) e ícones PWA/Maskable/Open Graph (`pwa-192x192.png`, `pwa-512x512.png`, `maskable-icon-512x512.png`, `og-image.png`).
- **Motor de Revisões Periódicas por Categoria (`periodic.go`)**: Regras de manutenção preventiva por quilometragem ou 12 meses para carros/caminhonetes (10.000 km / 20.000 km) e motos (2.500 km / 5.000 km), com marco zero referencial (`InitialMileage`) e endpoints `GET /api/v1/cars/{id}/upcoming-maintenances` e `GET /api/v1/maintenances/upcoming`.
- **Cronjob & Consolidação Anti-Spam de Push**: Scheduler em background (`maintenance_cron.go`, `MAINTENANCE_CRON_INTERVAL`) que agrupa múltiplos itens pendentes do mesmo veículo em uma única notificação push com *cooldown* de 7 dias.
- **Tela `/maintenance` e Cards Preventivos no Frontend**: Novo componente `UpcomingMaintenanceCard.tsx` com barra de progresso e registro rápido pré-preenchido, seletor de veículos com *deep-link* (`?carId=`), banner dinâmico de status da frota e histórico paginado.

## 📅 [2026-10-17](logs/2026-10-17.md) — Infraestrutura de Notificações Push com FCM no Backend e App Mobile (#93)
- **Domínio, Repositórios e FCM HTTP v1 (Go)**: Criação da entidade `DeviceToken`, coleção `device_tokens` no MongoDB com índices otimizados, cliente FCM HTTP v1 com limpeza automática de tokens `UNREGISTERED` e endpoints `POST /api/v1/devices/token`, `DELETE /api/v1/devices/token` e `POST /api/v1/notifications/test`.
- **Alertas de Manutenção Preventiva**: Implementação de verificação automática de revisões próximas ou vencidas por quilometragem/data com payload de *deep-link* (`/vehicles/:id`).
- **App Mobile Android & Capacitor Push**: Integração de `@capacitor/push-notifications`, permissão `POST_NOTIFICATIONS` para Android 13+, canal `cvmc_alerts`, toast em primeiro plano, *deep-linking*, desvínculo seguro no logout e card de teste em `ProfilePage.tsx`.
- **CI/CD & Deploy Automatizado de Segredos FCM**: Injeção automática de `GOOGLE_SERVICES_JSON` no build do APK Android (`frontend.yml`) e provisionamento de `FIREBASE_CREDENTIALS_JSON` no servidor OCI (`backend.yml`).

## 📅 [2026-10-16](logs/2026-10-16.md) — Prevenção de Vazamento de Tokens em LocalStorage
- **Isolamento de Tokens no Frontend**: Refatoração do `safeStorage.setSecureItem` e `useAuthStore` para impedir terminantemente o salvamento ou leitura de tokens JWT (`cvmc.access_token` e `cvmc.refresh_token`) no `localStorage` de navegadores web.
- **Isolamento de Runtime Mobile**: Restrição de persistência segura via `@capacitor/preferences` estritamente ao ambiente nativo (`Capacitor.isNativePlatform()`), assegurando que na web a aplicação opera exclusivamente com cookies `HttpOnly` conforme a [ADR 0002](adrs/0002-httponly-cookie-sessions.md).
- **Testes de Segurança**: Adicionada suíte de testes em `src/services/storage/storage.test.ts` validando isolamento de armazenamento local na web.

## 📅 [2026-10-11](logs/2026-10-11.md) — Empacotamento Android com Capacitor e Configuração de Build (#92)
- **Capacitor 8 no Frontend**: Configuração do `capacitor.config.ts` com `appId: com.cvmc.app`, inicialização do projeto `frontend/android` e scripts `cap:sync`, `cap:open` e `cap:build`.
- **Comunicação de API & Sessão Mobile**: BaseURL dinâmica para ambiente nativo, injeção de `Authorization: Bearer <token>`, persistência segura com `@capacitor/preferences` e retorno de tokens em JSON nos endpoints de autenticação em Go.
- **UX Mobile Nativa**: Controle da cor da barra de status (`#0284C7`), splash screen, hardware back button do Android e suporte a Safe Area Insets.
- **Correção de CORS e Tratamento de Erros**: Liberação das origens `https://localhost` e `capacitor://localhost` no backend Go e no deploy OCI (`backend.yml`), e refinamento no `LoginPage.tsx` para distinguir erro 401 de falha de conexão.
- **Status Bar Fixa e Preta**: Desativação de sobreposição de tela no Android (`overlay: false`) e definição de fundo preto sólido (`#000000`) com ícones claros para legibilidade total.
- **Sobreposição de Sidebar e Contraste Mobile (#104)**: Aplicação de safe area insets (`env(safe-area-inset-top)` e `env(safe-area-inset-bottom)`) no Drawer da `Sidebar.tsx`, eliminação de colisão com a barra de status do celular e configuração de alto contraste (`Style.Light`, `android:windowLightStatusBar="true"`, `#FFFFFF`) para o topo da aplicação.
- **Automação de Build do APK**: Novo job `build-android-apk` no workflow `frontend.yml` com injeção automática de `VITE_GOOGLE_CLIENT_ID`, build Gradle e disponibilização do artefato `app-debug-apk` no GitHub Actions.
- **Documentação de Build**: Criação do runbook operacional [docs/operations/mobile-build.md](operations/mobile-build.md) detalhando compilação do APK de debug, geração de bundle AAB e assinatura com keystore para publicação na Google Play.
- **Login Google Nativo e Layout Mobile (#105)**: Integração do plugin `@capgo/capacitor-social-login` com Android Credential Manager para contornar bloqueio de OAuth em WebViews, layout 100% responsivo no card de perfil sem overflow horizontal e feedback de erro aprimorado.
- **Validação Estrita de Client ID e Audience no Provedor Google**: Mitigação de vulnerabilidade de Confused Deputy / Audience Bypass em `validator.go`, abortando validação caso `GOOGLE_CLIENT_ID` não esteja configurado e validando o claim `aud` incondicionalmente.

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
