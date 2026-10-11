---
type: adr
title: "ADR 0002: Autenticação Exclusiva via Cookies HttpOnly e Sessões Stateless"
description: "Decisão de trafegar tokens de autenticação exclusivamente via cookies HttpOnly, abolindo armazenamento no localStorage."
tags:
  - adr
  - security
  - auth
  - cookies
timestamp: 2026-10-10
---

# 📜 ADR 0002: Autenticação Exclusiva via Cookies HttpOnly

- **Status**: Aceito
- **Data**: 2026-10-10
- **Decisores**: Equipe de Engenharia CVMC

Para navegação geral, retorne ao [Catálogo Canônico](../index.md).

---

## 📌 Contexto
O armazenamento de tokens JWT em `localStorage` expunha a aplicação a riscos de roubo de credenciais em caso de exploração de vulnerabilidades Cross-Site Scripting (XSS).

---

## 💡 Decisão
Adotamos o transporte exclusivo de tokens de autenticação via cookies com flags `HttpOnly`, `Secure` e `SameSite=Lax`:

1. O corpo das respostas de login e refresh não contém tokens brutos.
2. O frontend React permanece completamente agnóstico ao token físico, delegando o envio de credenciais ao próprio navegador.
3. Sessões permanecem stateless no backend, validadas através da assinatura criptográfica das chaves JWT.

---

## ⚖️ Consequências

### Positivas:
- **Proteção Completa contra XSS**: Tokens tornam-se inacessíveis para qualquer código JavaScript malicioso em execução no cliente.
- **Sincronia Automática**: Requisições para o backend enviam credenciais automaticamente via cookies.

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Autenticação e Segurança em Camadas](../architecture/auth-and-security.md)
- [Subdomínio de Usuários](../domain/user.md)
