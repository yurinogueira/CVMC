---
type: architecture
title: "Autenticação e Segurança em Camadas — CVMC"
description: "Padrões defensivos de autenticação via cookies HttpOnly, isolamento de sessões, mitigação OWASP e rate limiting."
tags:
  - architecture
  - security
  - auth
  - httponly
  - owasp
timestamp: 2026-10-10
---

# 🔒 Autenticação e Segurança em Camadas — CVMC

A segurança no **CVMC** segue o princípio da **Defesa em Profundidade (*Defense in Depth*)**, implementando proteções rigorosas contra vulnerabilidades comuns (OWASP Top 10).

Para navegação geral, retorne ao [Catálogo Canônico](../index.md).

---

## 🍪 1. Transporte de Sessão via Cookies `HttpOnly`

- **Proteção Total contra Roubo de Tokens por XSS**:
  - Tokens de acesso e atualização (`cvmc_access_token` e `cvmc_refresh_token`) trafegam exclusivamente via cookies HTTP com flags:
    - `HttpOnly: true` (inacessíveis via scripts JavaScript / `document.cookie`).
    - `Secure: true` (obrigatório em produção / conexões HTTPS).
    - `SameSite: Lax` (proteção nativa contra CSRF em requisições de terceiros).
    - `Domain`: configurado para permitir navegação transparente em subdomínios em produção (`.cvmc.yurinogueira.dev.br`).
- **Ciclo de Vida e Expiração dos Tokens (TTL)**:
  - `Access Token`: validade de 24 horas (`24*time.Hour`, `MaxAge: 24*3600`).
  - `Refresh Token`: validade estendida de 30 dias (`30*24*time.Hour`, `MaxAge: 30*24*3600`) com renovação contínua (*rolling refresh*) a cada requisição em `/api/v1/auth/refresh`, garantindo longevidade da sessão móvel e persistência do vínculo de notificações push FCM.
- **Zero Tokens no LocalStorage (Web) & Persistência Nativa Segura (Mobile)**:
  - No navegador web, o frontend React nunca armazena tokens JWT em `localStorage` ou `sessionStorage`, prevenindo roubo de credenciais via ataques XSS.
  - Apenas metadados não sensíveis de perfil (como nome e preferências visuais) são cacheados na store Zustand.
  - No ambiente móvel nativo (Android via Capacitor), os tokens de sessão (`accessToken`, `refreshToken`) e o perfil do usuário são mantidos com isolamento em `@capacitor/preferences` (sandbox Android `SharedPreferences`), inicializados assincronamente no bootstrap do app via `initAuthSession` com controle de prontidão (`isInitializing`) para evitar condições de corrida em rotas protegidas.

---

## 🛡️ 2. Proteção Criptográfica & Limites Contra DoS

- **Hashing de Senhas com BCrypt**:
  - Senhas são armazenadas utilizando `bcrypt` com custo padrão.
  - Para evitar ataques de exaustão de CPU (DoS em bcrypt), senhas são validadas estritamente no intervalo `8 <= len(password) <= 72`.
- **Validação de Segredos JWT no Startup**:
  - O método `Config.Validate()` aborta a inicialização da aplicação (`log.Fatalf`) caso `JWT_SECRET` ou `JWT_REFRESH_SECRET` estejam vazios, utilizem padrões inseguros (`change-me`), possuam menos de 32 caracteres ou sejam idênticos.
- **Hashes Criptográficos para Tokens Temporários**:
  - Tokens de redefinição de senha e confirmação de e-mail nunca são salvos em texto plano; o banco armazena o hash `SHA-256` do token gerado via `crypto/rand`.
- **Validação de ID Tokens do Google sem Intermediários**:
  - A autenticação social com Google utiliza validação direta de ID tokens criptográficos via biblioteca oficial `google.golang.org/api/idtoken`, sem dependência do Firebase.
  - O backend valida a assinatura, a expiração e confere o claim `aud` contra a variável `GOOGLE_CLIENT_ID`, rejeitando incondicionalmente a validação de tokens quando `GOOGLE_CLIENT_ID` não estiver configurado para mitigar ataques de Confused Deputy / Audience Bypass.
  - Usuários autenticados pelo Google têm seus e-mails considerados verificados pelo provedor de identidade, e contas são associadas de forma segura prevenindo lockout.

---

## 🚦 3. Rate Limiting e Sanitização

- **Rate Limiters Diferenciados**:
  - **Global**: 100 requisições por minuto com burst de 100 por IP.
  - **Auth**: 10 requisições por minuto com burst de 10 nos endpoints `/api/v1/auth/*` para mitigar ataques de força bruta.
  - **Sensível/Strict**: 5 requisições por minuto com burst de 5 para redefinição de senha e envio de e-mails.
- **Cabeçalhos de Segurança HTTP**:
  - Injeção obrigatória via middleware de `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`, `Referrer-Policy: strict-origin-when-cross-origin` e `Content-Security-Policy`.

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Visão Geral da Arquitetura](overview.md)
- [Subdomínio de Usuários e Autenticação](../domain/user.md)
- [Subdomínio de Anexos & Uploads](../domain/attachment.md)
- [ADR 0002: Cookies HttpOnly e Sessões Stateless](../adrs/0002-httponly-cookie-sessions.md)
