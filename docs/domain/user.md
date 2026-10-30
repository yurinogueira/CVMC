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
    AvatarURL                  string     `json:"avatarUrl,omitempty"`
    Role                       Role       `json:"role"`
    EmailVerified              bool       `json:"emailVerified"`
    EmailVerifiedAt            *time.Time `json:"emailVerifiedAt,omitempty"`
    EmailVerificationTokenHash string     `json:"-"`
    EmailVerificationExpiresAt *time.Time `json:"-"`
    PasswordResetTokenHash     string     `json:"-"`
    PasswordResetExpiresAt     *time.Time `json:"-"`
    MaxVehicles                int        `json:"maxVehicles"`
    MaxFiles                   int        `json:"maxFiles"`
    MaxStorageBytes            int64      `json:"maxStorageBytes"`
    FilesCount                 int        `json:"filesCount"`
    StorageUsedBytes           int64      `json:"storageUsedBytes"`
    StorageQuotaSynced         bool       `json:"-"`
    CreatedAt                  time.Time  `json:"createdAt"`
    UpdatedAt                  time.Time  `json:"updatedAt,omitempty"`
}
```

O método auxiliar `user.HasPassword()` verifica se o usuário possui senha local cadastrada (`PasswordHash != ""`), enquanto `user.EffectiveMaxFiles()` e `user.EffectiveMaxStorageBytes()` retornam os limites da cota de armazenamento do usuário (padrão de **100 arquivos** e **100 MB**). Os campos `FilesCount` e `StorageUsedBytes` persistem o consumo acumulado no próprio documento do usuário para leitura em $O(1)$, enquanto `StorageQuotaSynced` controla a migração única (backfill e criptografia in-place) de arquivos legados pré-existentes.

---

## 🔐 Autenticação Social & Vinculação de Contas (Google)

O CVMC oferece autenticação via Google Cloud sem uso de intermediários (Google Identity Services / GIS no frontend e `google.golang.org/api/idtoken` no backend Go).

### Endpoints
- `POST /api/v1/auth/google`: Login e auto-cadastro via Google ID token. Usuários autenticados pelo Google têm `EmailVerified = true` garantido automaticamente. Se o e-mail já existir, a conta Google é vinculada de forma transparente.
- `GET /api/v1/user/auth-providers`: Retorna a lista de provedores vinculados (`google` e `password`) e indica se o usuário possui senha local (`hasPassword`).
- `POST /api/v1/user/auth-providers/google/link`: Vincula a conta Google ao usuário logado via ID token.
- `POST /api/v1/user/auth-providers/google/unlink`: Desvincula a conta Google.
- `POST /api/v1/user/password/set`: Permite que usuários cadastrados exclusivamente via Google definam uma senha local inicial para login convencional.
- `DELETE /api/v1/user/account`: Exclui permanentemente a conta do usuário autenticado e executa a remoção em cascata de todos os seus dados (veículos próprios, manutenções, abastecimentos, vínculos de compartilhamento em veículos de terceiros e tokens FCM de dispositivos móveis), limpando também os cookies `HttpOnly` de sessão. Exige confirmação da senha atual (`currentPassword`) para contas com senha local ou a palavra-chave `"EXCLUIR"` (`confirmation`) para contas exclusivamente sociais.

### Sincronização e Exibição de Foto de Perfil (AvatarURL)
- **Captura no Token Google**: O claim `picture` do Google ID Token é capturado no payload e sanitizado via `domainuser.SanitizeAvatarURL()`.
- **Sanitização Estrita de Esquema**: Rejeita esquemas inseguros (`http://`, `javascript:`, URIs de dados base64), exigindo estritamente esquema `https://` com host válido para prevenir SSRF ou injeções arbitrárias.
- **Sincronização Contínua**:
  - No registro de novo usuário via Google (`POST /api/v1/auth/google`), a foto é persistida em `AvatarURL`.
  - No login recorrente ou auto-linking por e-mail (`POST /api/v1/auth/google`), `AvatarURL` é atualizada automaticamente caso o Google forneça uma nova URL de imagem.
  - Na vinculação de conta Google (`POST /api/v1/user/auth-providers/google/link`), o avatar é preenchido caso o usuário ainda não possua foto cadastrada.
  - Na desvinculação (`POST /api/v1/user/auth-providers/google/unlink`), a foto atual é mantida garantindo a integridade visual do perfil.
- **Renderização e Proteção no Frontend**:
  - Os componentes `<Avatar>` na `Topbar`, `Sidebar` e `ProfilePage` consomem `user.avatarUrl`.
  - Configuração obrigatória de `slotProps={{ img: { referrerPolicy: "no-referrer" } }}` para evitar bloqueios de requisição com erro `403 Forbidden` pelo CDN `*.googleusercontent.com`.
  - Fallback automático para a inicial maiúscula do nome do usuário caso a imagem falhe ao carregar ou o campo seja nulo.

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

