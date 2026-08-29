---
name: cvmc-dev
description: >-
  Guia e comandos otimizados para editar, validar e executar o projeto CVMC
  (Go backend, React frontend e Docker Compose), minimizando o consumo de tokens.
---

# Skill: Desenvolvimento e Execução do CVMC

Esta skill define as diretrizes de arquitetura, fluxos de edição, documentação OpenAPI/Swagger e o conjunto de ferramentas otimizadas para interagir com o projeto **CVMC (Como Vai Meu Carro)**, garantindo velocidade máxima e economia drástica de tokens na janela de contexto.

---

## ⚡ Regra de Ouro: Economia de Tokens

> [!IMPORTANT]
> **Nunca** execute comandos brutos e verbosos como `go test ./...` ou `npm run lint` diretamente, pois eles poluem o contexto com dezenas de linhas irrelevantes (`? [no test files]`, logs de build, etc.).
> **Sempre utilize os scripts auxiliares compactos em `scripts/`**:
> - `./scripts/check.sh all` ou `./scripts/check.sh backend|frontend|terraform`
> - `./scripts/fix.sh`
> - `./scripts/swagger.sh` (obrigatório ao alterar rotas/handlers do backend)
> - `./scripts/dev.sh start|stop|status|logs|restart`

### Comparativo de Comandos

| Tarefa | Comando Verboso (Evitar) | Comando Compacto (Usar) | Redução de Tokens |
| :--- | :--- | :--- | :--- |
| **Checagem Geral** | `go vet + go test + tsc + eslint + prettier + vitest + terraform fmt` | `./scripts/check.sh all` | **~85% menos tokens** |
| **Checar Backend** | `cd backend && go vet ./... && go test ./...` | `./scripts/check.sh backend` | Retorna 1 linha em sucesso |
| **Checar Frontend** | `cd frontend && npx tsc -b && npm run lint && ...` | `./scripts/check.sh frontend` | Retorna 1 linha em sucesso |
| **Checar Terraform**| `cd terraform && terraform fmt -check` | `./scripts/check.sh terraform` | Retorna 1 linha em sucesso |
| **Auto-Formatar** | `go fmt + prettier --write + eslint --fix + terraform fmt` | `./scripts/fix.sh` | Retorno limpo e direto |
| **Regenerar Swagger**| `cd backend && swag init ...` | `./scripts/swagger.sh` | Retorno direto e limpo |
| **Subir Stack** | `docker compose up -d` | `./scripts/dev.sh start` | Formata tabela limpa |
| **Status Containers** | `docker compose ps` | `./scripts/dev.sh status` | Resumo de portas e status |

---

## 🏗️ Visão Geral da Arquitetura

### 1. Backend (`backend/`)
- **Linguagem**: Go 1.25.
- **Padrão**: Clean Architecture + DDD.
- **Estrutura de Pastas**:
  - `internal/domain/`: Entidades puras (`car`, `user`, `maintenance`, `auth`), regras de negócio e validações de domínio. Todos os structs expostos na API contêm tags `json:"camelCase"`.
  - `internal/application/ports/`: Interfaces/contratos de repositórios, serviços de autenticação (`TokenService`, `PasswordHasher`), envio de e-mails (`EmailSender`) e storage (`StorageService`).
  - `internal/application/usecase/`:
    - `auth/`: Login, registro, validação de e-mail, reenvio de confirmação, esqueci minha senha e redefinição de senha com defesa anti-enumeração OWASP.
    - `car/`: Criação com validação de cota (`maxVehicles`), listagem, busca por ID, atualização (`UpdateCar`) e exclusão.
    - `user/`: Consulta de perfil público/privado, atualização de dados cadastrais e alteração segura de senha.
    - `maintenance/`: Histórico de manutenções e revisões preventivas.
    - `fipe/`: Integração resiliente com a API FIPE e cache multinível.
  - `internal/infrastructure/`: Implementações concretas (MongoDB com queries seguras, JWT, bcrypt com limite de 72 chars, Resend SMTP com codificação RFC 2045/MIME e Local Storage sanitizado).
  - `internal/interfaces/rest/`:
    - `handlers/`: Handlers HTTP que recebem DTOs, extraem a identidade via cookies `HttpOnly` (`cvmc_access_token`) ou Bearer token (sempre retornando `""` em caso de erro, forçando 401) e possuem anotações declarativas do Swaggo.
    - `router.go`: Registro de rotas HTTP, montagem da cadeia de middlewares de segurança (`SecurityHeaders`, `CORS` com whitelist, `RateLimiter` multicamada, `BodyLimit` de 1MB) e rota condicional `GET /swagger/` (apenas em `LOG_LEVEL=debug`).
  - `internal/shared/`: Middlewares (CORS, RateLimit, BodyLimit, SecurityHeaders, RequestID, Logging) e envelope HTTP padronizado (`httpx.Success`, `httpx.Created`, `httpx.Error`).
  - `docs/`: Documentação gerada automaticamente pelo Swaggo (`docs.go`, `swagger.json`, `swagger.yaml`).
  - `cmd/api/main.go`: Ponto de entrada do servidor HTTP com timeouts defensivos e graceful shutdown.

### 2. Frontend (`frontend/`)
- **Stack**: React 19, TypeScript, Vite, Material UI (MUI v6), `@mui/icons-material`, Zustand, React Router v7, Axios, Vitest.
- **Estrutura de Pastas**:
  - `src/features/auth/`: Módulo de autenticação com layout Split-Screen (`LoginPage.tsx`, `RegisterPage.tsx`, `ForgotPasswordPage.tsx`, `ResetPasswordPage.tsx`, `VerifyEmailPage.tsx`, `AuthHeroBanner.tsx`), store Zustand (`auth.store.ts` gerenciando apenas perfil público, sem tokens) e `auth.service.ts`.
  - `src/features/cars/`: Gestão completa de veículos (`VehiclesPage.tsx`, `VehicleCard.tsx`, `AddCarDialog.tsx`, `EditCarDialog.tsx`, `ImageUploadField.tsx`, `VehicleImage.tsx`, `car.service.ts`).
  - `src/features/profile/`: Gestão de conta e perfil (`ProfilePage.tsx`, `ProfileCard.tsx`, `ChangePasswordDialog.tsx`, `user.service.ts`).
  - `src/features/dashboard/`: Painel geral com KPIs, CTA em gradiente, resumo de manutenções e Empty State inteligente (`DashboardPage.tsx`).
  - `src/features/maintenance/`: Histórico e revisões de veículos (`MaintenancePage.tsx`).
  - `src/layouts/`: Shell SaaS persistente (`Sidebar.tsx`, `Topbar.tsx`, `AppLayout.tsx`, `BrandLogo.tsx`).
  - `src/routes/`: Definição de rotas públicas (`/login`, `/register`, `/forgot-password`, `/reset-password`, `/verify-email`) e privadas via `ProtectedRoute.tsx`.
  - `src/services/api/`: Cliente Axios (`client.ts`) configurado para `/api/v1` com `withCredentials: true` para transporte transparente e seguro de cookies `HttpOnly`.
  - `src/services/storage/`: Wrapper de acesso seguro ao `localStorage` (`storage.ts`) para suporte a testes e SSR (apenas para preferências de UI, nunca tokens).
  - `src/styles/`: Configuração de tema Material UI com harmonia análoga (`theme.ts`).

---

## 🛠️ Fluxo de Trabalho e Regras Mandatórias

### 1. Editando o Backend & Regra do Swagger
1. Ao adicionar/modificar entidades, atualize `internal/domain/` garantindo tags `json:"camelCase"`.
2. Defina os contratos em `internal/application/ports/`.
3. Implemente as regras nos use cases (`internal/application/usecase/`) e crie/atualize os testes unitários (`*_test.go`).
4. Implemente os adapters em `internal/infrastructure/` e handlers em `internal/interfaces/rest/handlers/`.
5. > [!IMPORTANT]
   > **Regra Obrigatória do Swagger**: Sempre que adicionar, alterar ou remover uma rota, parâmetro ou handler HTTP, adicione/atualize os comentários `@Summary`, `@Tags`, `@Param`, `@Success`, `@Router`, etc. no handler e **execute imediatamente**:
   > ```bash
   > ./scripts/swagger.sh
   > ```
6. Valide a integridade do backend com:
   ```bash
   ./scripts/check.sh backend
   ```

### 2. Editando o Frontend
1. Crie ou modifique componentes/telas em `src/features/<feature>/`.
2. Utilize ícones semânticos de `@mui/icons-material` e componentes padronizados do Design System.
3. Para acesso ao armazenamento local de preferências, use sempre o utilitário seguro `safeStorage` de `src/services/storage/storage.ts`.
4. Gerencie estado global via stores Zustand em `src/features/<feature>/state/`.
5. Desacople objetos reativos de dependências de `useEffect` (usando `useAuthStore.getState()` ou seletores estáveis) para prevenir loops de requisição.
6. Valide tipos, lint, formatação e testes unitários com:
   ```bash
   ./scripts/check.sh frontend
   ```
7. Se houver divergências de formatação ou lint simples:
   ```bash
   ./scripts/fix.sh
   ```

### 3. Rodando o Projeto com Docker Compose
Para inicializar e gerenciar a stack completa:
- **Iniciar serviços**: `./scripts/dev.sh start`
- **Verificar status**: `./scripts/dev.sh status`
- **Consultar logs**: `./scripts/dev.sh logs [backend|frontend|mongo]`
- **Reiniciar serviço**: `./scripts/dev.sh restart backend`
- **Parar serviços**: `./scripts/dev.sh stop`

Portas e URLs padrão:
- **Frontend App**: `http://localhost:5173`
- **Backend API**: `http://localhost:8080`
- **Documentação Swagger UI**: `http://localhost:8080/swagger/index.html`
- **MongoDB**: `localhost:27017`

---

## 🚨 Resolução Rápida de Erros (Runbook)

- **Falha no `check.sh backend`**:
  - O script isola e exibe apenas o pacote ou teste que falhou.
  - Para rodar um teste específico sem ruído:
    `cd backend && go test -v ./internal/application/usecase/auth`
- **Falha no `check.sh frontend`**:
  - Se for erro de tipos: verifique a interface TypeScript indicada pelo `tsc`.
  - Se for erro de formatação: rode `./scripts/fix.sh`.
- **Containers não sobem ou alterações no Go não refletem**:
  - Execute `docker compose build backend && docker compose up -d backend` para recompilar a imagem da API Go.
  - Execute `./scripts/dev.sh logs` para ver as últimas 30 linhas de log.
