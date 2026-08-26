---
name: cvmc-security
description: >-
  Auditoria contínua de segurança, prevenção de vulnerabilidades (OWASP Top 10),
  validação de autenticação por cookies HttpOnly, proteção contra enumeração,
  hashing seguro de tokens, rate limiting e sanitização de infraestrutura para o CVMC.
---

# Skill: Segurança, Defesa em Profundidade e Prevenção de Vulnerabilidades — CVMC

Esta skill define o protocolo mandatório de auditoria de segurança, checklist pré-commit/pré-deploy e as diretrizes de defesa em profundidade aplicadas a todas as camadas do **CVMC (Como Vai Meu Carro)**.

---

## 🛡️ Protocolo de Segurança Mandatório

Antes de submeter qualquer modificação que envolva autenticação, usuários, rotas, persistência ou infraestrutura, execute e valide cada um dos checklists abaixo:

### 1. Checklist de Autenticação, Tokens & Sessões
- [ ] **Transporte de Sessão**: Tokens JWT (`accessToken` e `refreshToken`) são trafegados exclusivamente em cookies `HttpOnly`, `Secure`, `SameSite=Lax` com domínio configurado.
- [ ] **Isolamento no Frontend**: O frontend **nunca** armazena tokens em `localStorage` ou `sessionStorage` (apenas dados públicos de perfil não sensíveis).
- [ ] **Falha Segura de Identidade**: Ao extrair o usuário via `claims.UserID`, em caso de falha de parse ou token expirado, a função retorna estritamente `""` (forçando `401 Unauthorized`), **sem nunca fazer fallback para o token bruto**.
- [ ] **Limites de Senha (Prevenção a DoS no Bcrypt)**:
  - Senhas no cadastro, redefinição e alteração validam estritamente o tamanho: `8 <= len(password) <= 72`.
  - O limite superior de 72 caracteres é obrigatório para evitar ataques de DoS por exaustão de CPU no algoritmo bcrypt.

### 2. Checklist de Validação de E-mail & Recuperação de Senha (OWASP Defense)
- [ ] **Prevenção de Enumeração de Usuários**:
  - O endpoint `POST /api/v1/auth/forgot-password` sempre responde com status `200 OK` e mensagem opaca genérica, independentemente de o e-mail existir ou não na base.
  - O timing de resposta não deve revelar a existência do e-mail.
- [ ] **Armazenamento de Tokens de Uso Único**:
  - Tokens para verificação de e-mail e recuperação de senha são gerados com 32 bytes aleatórios criptográficos (`crypto/rand`).
  - **Nunca armazene tokens em texto plano no banco**. Persista apenas o **hash SHA-256** do token.
  - Validade estrita com timestamp de expiração (máx. 30 min para reset de senha; 24h para validação de e-mail).
  - Invalidação atômica e imediata do token após o uso bem-sucedido.
- [ ] **Controle de Cota e Permissão por Status de Conta**:
  - Usuários não verificados (`emailVerified == false`) têm o cadastro de recursos sensíveis (ex.: veículos) bloqueado (`403 Forbidden`).
  - Verificação atômica de cota máxima de recursos (`maxVehicles`) antes da criação para evitar abusos.

### 3. Checklist de Rede, Middlewares & Rate Limiting
- [ ] **Cadeia Global de Segurança**:
  - Toda rota pública ou protegida passa pela cadeia de segurança:
    ```go
    middleware.SecurityHeaders(
        middleware.CORS(
            globalLimiter.Limit(
                middleware.BodyLimit(1<<20)(mux), // 1MB body limit
            ),
            cfg.AllowedOrigins,
        ),
    )
    ```
- [ ] **CORS Estrito**: Whitelist estrita via `cfg.AllowedOrigins` (nunca utilizar wildcard `*` em conjunto com `AllowCredentials`).
- [ ] **Rate Limiting Multicamada**:
  - Global: 100 req/min.
  - Auth padrão (login/register): 10 req/min.
  - Strict (forgot-password, reset-password, resend-verification): 5 req/min.
- [ ] **Condicionamento de Debug**: A rota `/swagger/` é exposta exclusivamente quando `cfg.LogLevel == "debug"`.

### 4. Checklist de Persistência & Sanitização
- [ ] **Sanitização NoSQL**: IDs e e-mails passam por sanitização (`SanitizeID`, `SanitizeEmail`) antes de consultas no MongoDB para mitigar injeção de operadores BSON.
- [ ] **Storage e Path Traversal**: Provedores de armazenamento validam caminhos absolutos e impedem path traversal (`../`) antes de qualquer operação em disco.

### 5. Checklist de Infraestrutura & Terraform
- [ ] O arquivo `docker-compose.yml` de produção não expõe portas de banco de dados (`27017:27017`) diretamente no host.
- [ ] Nenhum segredo ou chave privada RSA está em texto plano em `.tfvars` ou `.hcl`.
- [ ] A regra de SSH no OCI utiliza `var.admin_cidr` em vez de `0.0.0.0/0`.

---

## 🚨 Matriz de Padrões Inseguros vs. Padrões Seguros

| Padrão Inseguro (PROIBIDO) | Padrão Seguro (OBRIGATÓRIO) | Vulnerabilidade Evitada |
| :--- | :--- | :--- |
| `return raw` quando `ParseAccessToken` falha | `return ""` (dispara 401) | Bypass de autenticação |
| `safeStorage.setItem("cvmc.accessToken", token)` | Cookies `HttpOnly`, `Secure`, `SameSite=Lax` | Roubo de sessão via XSS |
| Vazar `ErrUserNotFound` em `/forgot-password` | Retornar `200 OK` genérico | Enumeração de usuários |
| Salvar token de reset em texto claro no banco | Salvar hash SHA-256 do token | Comprometimento por vazamento de banco |
| Senha sem limite máximo (`len > 72`) | Validar `8 <= len(password) <= 72` | DoS no Bcrypt (CPU Exhaustion) |
| Permitir cadastro de carro com e-mail não validado | Validar `user.EmailVerified == true` | Criação de spam e contas falsas |
| `w.Header().Set("Access-Control-Allow-Origin", "*")` com credenciais | Whitelist de origens via `ALLOWED_ORIGINS` | Vazamento de dados cross-origin |
| `filepath.Join(base, path)` sem checagem de prefixo | `filepath.Abs()` + checagem de prefixo base | Arbitrary File Read / Path Traversal |
| `ports: - "27017:27017"` no MongoDB em produção | Rede interna Docker (`mongodb://mongo:27017`) | Exposição pública do banco de dados |

---

## ⚡ Comandos Rápidos de Validação de Segurança

```bash
# Validação geral de integridade, testes unitários, tipos e lints
./scripts/check.sh all

# Regenera a documentação OpenAPI/Swagger garantindo tipagem precisa
./scripts/swagger.sh
```
