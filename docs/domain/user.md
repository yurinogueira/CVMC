---
type: domain
title: "Subdomínio: Usuários, Credenciais e Perfil (User)"
description: "Modelos de contas de usuário, papéis de acesso (RBAC), tokens de recuperação de senha e verificação de e-mail."
tags:
  - domain
  - user
  - auth
  - rbac
resource: backend/internal/domain/user
timestamp: 2026-10-10
---

# 👤 Subdomínio: Usuários & Perfil (User)

O subdomínio **User** gerencia a identidade dos proprietários e condutores na plataforma **CVMC**.

Para navegação geral, retorne ao [Catálogo Canônico](../index.md).

---

## 🧩 Modelo de Domínio (`User`)

Estruturado em `backend/internal/domain/user/user.go`:

```go
type Role string

const (
    RoleAdmin Role = "admin"
    RoleUser  Role = "user"
)

type User struct {
    ID                         string     `json:"id"`
    Name                       string     `json:"name"`
    Email                      string     `json:"email"`
    PasswordHash               string     `json:"-"`
    GoogleID                   string     `json:"googleId,omitempty"`
    Role                       Role       `json:"role"`
    EmailVerified              bool       `json:"emailVerified"`
    EmailVerifiedAt            *time.Time `json:"emailVerifiedAt,omitempty"`
    EmailVerificationTokenHash string     `json:"-"`
    EmailVerificationExpiresAt *time.Time `json:"-"`
    PasswordResetTokenHash     string     `json:"-"`
    PasswordResetExpiresAt     *time.Time `json:"-"`
    CreatedAt                  time.Time  `json:"createdAt"`
    UpdatedAt                  time.Time  `json:"updatedAt,omitempty"`
}
```

O método auxiliar `user.HasPassword()` verifica se o usuário possui senha local cadastrada (`PasswordHash != ""`).

---

## 🔐 Autenticação Social & Vinculação de Contas (Google)

O CVMC oferece autenticação via Google Cloud sem uso de intermediários (Google Identity Services / GIS no frontend e `google.golang.org/api/idtoken` no backend Go).

### Endpoints
- `POST /api/v1/auth/google`: Login e auto-cadastro via Google ID token. Usuários autenticados pelo Google têm `EmailVerified = true` garantido automaticamente. Se o e-mail já existir, a conta Google é vinculada de forma transparente.
- `GET /api/v1/user/auth-providers`: Retorna a lista de provedores vinculados (`google` e `password`) e indica se o usuário possui senha local (`hasPassword`).
- `POST /api/v1/user/auth-providers/google/link`: Vincula a conta Google ao usuário logado via ID token.
- `POST /api/v1/user/auth-providers/google/unlink`: Desvincula a conta Google.
- `POST /api/v1/user/password/set`: Permite que usuários cadastrados exclusivamente via Google definam uma senha local inicial para login convencional.

---

## 📋 Regras de Segurança e Privacidade

1. **Ocultação de Dados Críticos (`json:"-"`)**:
   - Campos como `PasswordHash` e hashes de tokens nunca são serializados em respostas JSON.
2. **Prevenção de Lockout na Desvinculação**:
   - Uma conta Google só pode ser desvinculada se o usuário possuir senha definida (`user.HasPassword()`). Caso contrário, a requisição é rejeitada para impedir que o usuário perca o acesso.
3. **Ciclo de Verificação e Redefinição**:
   - Tokens de verificação de e-mail possuem validade de 24 horas.
   - Tokens de redefinição de senha possuem validade de 30 minutos e são invalidados imediatamente após uso com sucesso.

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Autenticação e Segurança em Camadas](../architecture/auth-and-security.md)
- [Isolamento de Dados e Persistência](../architecture/multitenancy-and-data.md)
- [Subdomínio de Veículos](car.md)

