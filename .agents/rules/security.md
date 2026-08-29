---
description: Diretrizes mandatórias de segurança, prevenção de vulnerabilidades, autenticação por cookies HttpOnly, CORS estrito, sanitização e infraestrutura segura no CVMC
globs: "**/*"
---

# 🔒 Regras de Segurança Mandatórias — CVMC

Estas diretrizes são **inegociáveis** e devem ser seguidas em qualquer alteração de código ou infraestrutura no projeto **CVMC (Como Vai Meu Carro)** para impedir a introdução de vulnerabilidades de segurança.

---

## 🛑 1. Autenticação, Autorização e Sessões

### Backend (Go)
- **Cookies `HttpOnly` Obrigatórios**:
  - Tokens de acesso e atualização (`cvmc_access_token` e `cvmc_refresh_token`) **NUNCA** devem ser retornados no corpo da resposta JSON.
  - Devem ser setados exclusivamente em cookies HTTP com:
    - `HttpOnly: true` (proteção total contra roubo via XSS)
    - `Secure: true` (em produção / HTTPS)
    - `SameSite: http.SameSiteLaxMode` (necessário para subdomínios `cvmc.` e `api.cvmc.`)
    - `Domain: .cvmc.yurinogueira.dev.br` (compartilhamento seguro entre subdomínios em produção)
- **Proibido Fallback de Token Inválido**:
  - Em funções de extração de identidade (`extractUserID`), se a validação/parse do token falhar ou for nulo, **retorne imediatamente `""`** (vazio) para forçar `401 Unauthorized`.
  - **NUNCA** retorne o token bruto, string recebida no header ou fallback para IDs padrão.
- **Validação de Senhas**:
  - Toda senha de usuário deve ser validada com tamanho mínimo de **8 caracteres** no registro.
  - Hashing sempre via `bcrypt` com custo mínimo `DefaultCost` (10).

### Frontend (React)
- **Proibido Armazenar Tokens no Client**:
  - **NUNCA** armazene `accessToken` ou `refreshToken` no `localStorage`, `sessionStorage` ou variáveis globais `window`.
  - O Axios deve estar configurado com `withCredentials: true` para enviar cookies automaticamente.
  - O `auth.store.ts` deve armazenar apenas informações públicas do perfil do usuário (`id`, `name`, `email`).

---

## 🌐 2. Proteção de Rede, CORS & Middlewares

- **CORS Estrito com Whitelist**:
  - **NUNCA** utilize wildcard `*` ou reflita cabeçalhos `Origin` arbitrariamente com `Access-Control-Allow-Credentials: true`.
  - As requisições devem ser validadas estritamente contra a lista de origens configurada (`ALLOWED_ORIGINS`).
- **Rate Limiting Ativo**:
  - Todas as rotas públicas devem passar pelo rate limiter global (100 req/min).
  - Rotas críticas de autenticação (`/auth/login`, `/auth/register`) devem ter limite restrito (10 req/min) contra força bruta.
- **Limite de Payload (Body Limit)**:
  - Todas as requisições HTTP devem ser limitadas a no máximo **1MB** via middleware `http.MaxBytesReader`.
- **Security Headers Padronizados**:
  - Toda resposta HTTP deve conter:
    - `X-Content-Type-Options: nosniff`
    - `X-Frame-Options: DENY`
    - `Referrer-Policy: strict-origin-when-cross-origin`
    - `Permissions-Policy: camera=(), microphone=(), geolocation=()`

---

## 🗄️ 3. Acesso ao Sistema de Arquivos & Storage

- **Prevenção de Path Traversal (LFI)**:
  - Nunca utilize `filepath.Join(basePath, userInput)` diretamente sem validar se o caminho final pertence ao diretório raiz:
    ```go
    absBase, _ := filepath.Abs(basePath)
    absFull, _ := filepath.Abs(fullPath)
    if !strings.HasPrefix(absFull, absBase+string(filepath.Separator)) && absFull != absBase {
        return errors.New("invalid path: traversal detected")
    }
    ```

---

## 🔑 4. Gestão de Segredos & Configuração

- **Rejeição Incondicional de Segredos JWT Padrão/Fracos**:
  - Em qualquer ambiente (independentemente de `LOG_LEVEL`), a aplicação **deve abortar a inicialização** (`log.Fatalf`) se `JWT_SECRET` ou `JWT_REFRESH_SECRET` contiverem valores default como `"change-me"`, `"change-me-too"` ou possuírem comprimento inferior a **32 caracteres** (`len < 32`).
- **Documentação de API Segura**:
  - O endpoint de Swagger UI (`/swagger/`) só pode ser registrado e acessível quando `LOG_LEVEL == "debug"`. Em produção, o endpoint não deve existir.

---

## 🐳 5. Infraestrutura, Docker & Terraform

- **Isolamento de Banco de Dados**:
  - Portas internas de bancos de dados (ex.: MongoDB `27017:27017`) **NUNCA** devem ser mapeadas para o host no `docker-compose.yml`. O acesso deve ocorrer exclusivamente pela rede interna dos containers.
- **Chaves de API do Terraform**:
  - **NUNCA** coloque chaves privadas RSA em texto plano em arquivos `.tfvars`, `.hcl` ou `.tf`.
  - Use sempre o caminho do arquivo de chave (`private_key_path = "~/.oci/oci_api_key.pem"`).
- **Acesso SSH Restrito**:
  - Regras de ingress de SSH (porta 22) no Terraform não devem usar `0.0.0.0/0` irrestrito; devem usar `var.admin_cidr`.
