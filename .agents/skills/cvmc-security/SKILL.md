---
name: cvmc-security
description: >-
  Auditoria contínua de segurança, prevenção de vulnerabilidades (OWASP Top 10),
  validação de autenticação por cookies HttpOnly, CORS estrito, rate limiting e
  sanitização de infraestrutura para o projeto CVMC.
---

# Skill: Segurança e Prevenção de Vulnerabilidades — CVMC

Esta skill define o protocolo de auditoria contínua, checklist de segurança pré-commit/pré-deploy e as verificações mandatórias para garantir que nenhuma vulnerabilidade seja introduzida no **CVMC (Como Vai Meu Carro)**.

---

## 🛡️ Protocolo de Segurança Pré-Deploy / Pré-PR

Antes de finalizar qualquer alteração que toque em rotas, autenticação, armazenamento ou infraestrutura, execute o seguinte checklist:

### 1. Checklist de Autenticação & Sessões
- [ ] O handler extrai o usuário via `claims.UserID` a partir de `cvmc_access_token` ou Bearer JWT?
- [ ] Em caso de falha de parse ou token expirado, a função retorna estritamente vazio `""` (resultando em `401 Unauthorized`), **sem nunca fazer fallback para o token bruto**?
- [ ] Os tokens são gerados e trafegados exclusivamente em cookies `HttpOnly`, `Secure`, `SameSite=Lax` com domínio configurado?
- [ ] O frontend não manipula `accessToken` ou `refreshToken` no `localStorage` ou `sessionStorage`?
- [ ] Senhas no cadastro possuem validação de no mínimo 8 caracteres?

### 2. Checklist de Rede & Middlewares
- [ ] Novas rotas no backend passam pela cadeia de segurança:
  ```go
  middleware.SecurityHeaders(
      middleware.CORS(
          globalLimiter.Limit(
              middleware.BodyLimit(1<<20)(mux),
          ),
          cfg.AllowedOrigins,
      ),
  )
  ```
- [ ] O CORS utiliza estritamente `cfg.AllowedOrigins` (sem wildcard `*` com credentials)?
- [ ] Rotas sensíveis (login/register) possuem rate limit restrito (ex.: 10 req/min)?
- [ ] O endpoint `/swagger/` está estritamente condicionado a `cfg.LogLevel == "debug"`?

### 3. Checklist de Arquivos & Storage
- [ ] O provedor de storage valida caminhos absolutos e impede path traversal (`../`) antes de salvar ou deletar arquivos?

### 4. Checklist de Infraestrutura & Terraform
- [ ] O arquivo `docker-compose.yml` de produção não expõe portas de banco de dados (`27017:27017`) diretamente no host?
- [ ] Nenhum segredo ou chave privada RSA está em texto plano em `.tfvars` ou `.hcl`?
- [ ] A regra de SSH no OCI utiliza `var.admin_cidr` em vez de `0.0.0.0/0`?

---

## ⚡ Comandos Rápidos de Validação de Segurança

Execute sempre a validação compacta padrão:

```bash
# Validação geral de integridade, testes unitários, tipos e lints
./scripts/check.sh all
```

Ao adicionar ou modificar endpoints HTTP:
```bash
# Regenera a documentação da API Swagger mantendo os tipos consistentes
./scripts/swagger.sh
```

---

## 🚨 Padrões de Código Inseguros a Rejeitar Imediatamente

| Padrão Inseguro (PROIBIDO) | Padrão Seguro (OBRIGATÓRIO) |
| :--- | :--- |
| `return raw` quando `ParseAccessToken` falha | `return ""` (dispara 401) |
| `safeStorage.setItem("cvmc.accessToken", token)` | Cookies `HttpOnly` injetados pelo backend |
| `w.Header().Set("Access-Control-Allow-Origin", "*")` com credentials | Whitelist de origens via `ALLOWED_ORIGINS` |
| `filepath.Join(base, path)` sem verificação de prefixo | `filepath.Abs()` + checagem de prefixo do diretório base |
| `ports: - "27017:27017"` no MongoDB | Comunicação apenas via rede interna Docker (`mongodb://mongo:27017`) |
| `private_key = <<EOF ... EOF` em arquivos git | `private_key_path = "~/.oci/oci_api_key.pem"` |
