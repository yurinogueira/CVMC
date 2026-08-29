# Guia de Contribuição — CVMC (Como Vai Meu Carro) 🚗

Obrigado pelo seu interesse em contribuir com o **CVMC**! Este projeto é uma plataforma moderna e completa para gestão, histórico e controle preventivo de manutenção veicular.

Este guia estabelece os padrões técnicos, fluxos de versionamento, boas práticas de segurança e diretrizes de desenvolvimento para garantir uma base de código sólida, segura e agradável de manter.

---

## 🧭 Índice

1. [Código de Conduta](#-código-de-conduta)
2. [Arquitetura & Stack Tecnológica](#-arquitetura--stack-tecnológica)
3. [Configuração do Ambiente Local](#-configuração-do-ambiente-local)
4. [Scripts Auxiliares de Produtividade](#-scripts-auxiliares-de-produtividade)
5. [Fluxo de Trabalho Git (Branching & Commits)](#-fluxo-de-trabalho-git)
6. [Diretrizes de Desenvolvimento](#-diretrizes-de-desenvolvimento)
   - [Backend (Go)](#backend-go)
   - [Frontend (React + MUI)](#frontend-react--mui)
   - [Segurança & Defesa em Profundidade](#segurança--defesa-em-profundidade)
7. [Submissão de Pull Requests](#-submissão-de-pull-requests)
8. [Reportando Problemas & Sugerindo Funcionalidades](#-reportando-problemas--sugerindo-funcionalidades)

---

## 📜 Código de Conduta

Ao participar deste projeto, você concorda em cumprir e zelar pelos princípios estabelecidos em nosso [Código de Conduta](CODE_OF_CONDUCT.md) (Contributor Covenant 2.1). Esperamos um ambiente acolhedor, respeitoso e livre de assédio para todos.

---

## 🏗️ Arquitetura & Stack Tecnológica

O CVMC adota uma arquitetura desacoplada e modular:

- **Backend (`/backend`)**:
  - **Linguagem**: Go 1.25.
  - **Arquitetura**: Clean Architecture + Domain-Driven Design (DDD) (`domain`, `ports`, `usecase`, `infrastructure`, `interfaces/rest`).
  - **Banco de Dados**: MongoDB (Driver Oficial Go v2).
  - **Autenticação**: JWT com transporte exclusivo em cookies `HttpOnly`, `Secure` e `SameSite=Lax`.
  - **Documentação da API**: OpenAPI 2.0 / 3.0 via Swaggo (`./scripts/swagger.sh`).
- **Frontend (`/frontend`)**:
  - **Framework**: React 19 + TypeScript + Vite.
  - **UI / Design System**: Material UI v6 (`@mui/material`, `@mui/icons-material`) com paleta análoga e acessibilidade WCAG AA.
  - **Gerenciamento de Estado**: Zustand (gerenciando apenas perfis públicos e estado de UI, sem tokens).
  - **Roteamento & Requisições**: React Router v7 e Axios (`withCredentials: true`).
- **Infraestrutura & DevOps**:
  - **Ambiente Local**: Docker Compose (`docker-compose.yml`).
  - **IaC**: Terraform para provisionamento em nuvem (OCI & Cloudflare).
  - **CI/CD**: GitHub Actions para testes, linters, segurança e deploy contínuo.

---

## 💻 Configuração do Ambiente Local

### Pré-requisitos
- **Git**: `>= 2.30`
- **Docker & Docker Compose**: Docker 24+ com Compose v2.
- **Go**: `>= 1.25` (opcional caso utilize apenas os containers).
- **Node.js**: `>= 20 LTS` e **npm** `>= 10` (opcional caso utilize apenas os containers).

### Inicializando o Projeto

1. Clone o repositório:
   ```bash
   git clone https://github.com/yurinogueira/CVMC.git
   cd CVMC
   ```

2. Inicialize a stack completa com o script de desenvolvimento:
   ```bash
   ./scripts/dev.sh start
   ```

3. URLs e Portas de Acesso:
   - **Aplicação Frontend**: [http://localhost:5173](http://localhost:5173)
   - **API REST Backend**: [http://localhost:8080](http://localhost:8080)
   - **Documentação Swagger UI**: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
   - **MongoDB**: `localhost:27017`

---

## ⚡ Scripts Auxiliares de Produtividade

O projeto conta com scripts otimizados no diretório `scripts/` para acelerar verificações e economizar recursos:

| Comando | Descrição |
| :--- | :--- |
| `./scripts/check.sh all` | Executa validação completa (Backend vet/testes, Frontend tsc/lint/testes e Terraform fmt) |
| `./scripts/check.sh backend` | Valida formatação, `go vet` e todas as suítes de testes unitários em Go |
| `./scripts/check.sh frontend` | Executa `tsc -b`, ESLint, Prettier e Vitest no frontend |
| `./scripts/check.sh terraform`| Verifica a formatação dos arquivos `.tf` |
| `./scripts/fix.sh` | Corrige automaticamente formatação em Go, Prettier, ESLint e Terraform |
| `./scripts/swagger.sh` | **Obrigatório**: Regenera os artefatos OpenAPI/Swagger após alterações em rotas/handlers |
| `./scripts/dev.sh start\|stop\|status\|logs` | Gerencia o ciclo de vida dos containers Docker |

---

## 🔄 Fluxo de Trabalho Git

Seguimos um fluxo estrito para manter o histórico linear, rastreável e sempre pronto para deploy contínuo:

### 1. Sincronização Obrigatória com a `main`
Antes de criar qualquer branch ou iniciar modificações, **sempre sincronize com a versão mais recente da `main` remota**:

```bash
# 1. Buscar referências remotas
git fetch origin main

# 2. Criar a nova branch a partir de origin/main
git checkout -b <tipo>/<nome-da-branch> origin/main
```

### 2. Padrão de Nomenclatura de Branches
- **Com Issue vinculada**: `<tipo>/<id_da_issue>-<descricao-curta>`
  - `feat/24-user-profile-and-auth-limits`
  - `fix/28-remove-hardcoded-deploy-ip`
  - `feat/42-edit-vehicle-and-year`
- **Sem Issue vinculada (docs/skills/infra)**: `<tipo>/<descricao-curta>`
  - `chore/skills-enhancement`
  - `docs/community-health-files`
  - `feat/dependency-updates`

### 3. Commits Semânticos (Conventional Commits)
Organize seus commits de forma atômica seguindo o padrão:
- **Com Issue**: `<tipo>(<escopo>): <descrição clara no imperativo> (#<id_da_issue>)`
  - Exemplo: `feat(cars): permitir edicao de veiculo e incluir campo de ano (#42)`
  - Exemplo: `fix(profile): desacoplar user do useEffect para prevenir loop (#47)`
- **Sem Issue**: `<tipo>(<escopo>): <descrição clara no imperativo>`
  - Exemplo: `docs(community): adicionar pull request template, contributing e code of conduct`
  - Exemplo: `chore(deps): atualizar dependencias do frontend e driver mongodb`

---

## 📐 Diretrizes de Desenvolvimento

### Backend (Go)
1. **Clean Architecture**: Regras de negócio puras em `domain/`, contratos em `ports/`, orquestração em `usecase/`, implementações em `infrastructure/` e HTTP em `interfaces/rest/`.
2. **Tags JSON CamelCase**: Todas as structs expostas na API devem utilizar `json:"camelCase"`.
3. **Regra do Swagger**: Sempre que adicionar ou alterar handlers HTTP, adicione anotações Swaggo (`@Summary`, `@Tags`, `@Param`, `@Success`, `@Router`) e execute:
   ```bash
   ./scripts/swagger.sh
   ```
4. **Testes**: Crie testes unitários para use cases (`*_test.go`) cobrindo cenários de sucesso e caminhos de erro.

### Frontend (React + MUI)
1. **Design System & Acessibilidade**:
   - Utilize a paleta análoga do projeto (`#0284C7`, `#0369A1`, `#0EA5E9`, `#16A34A`).
   - Garanta contraste mínimo WCAG AA (>= 4.5:1 para textos regulares).
   - Use hierarquia semântica correta (`h1` único por página, seguido por `h2` e `h3`).
2. **Gerenciamento de Estado & Efeitos**:
   - Stores Zustand apenas para estado de UI e perfil público.
   - Desacople objetos instáveis de arrays de dependência de `useEffect` (use seletores de store ou `getState()`) para prevenir loops de requisição.
3. **Armazenamento Seguro**: Use sempre `safeStorage` (`src/services/storage/storage.ts`) para preferências de interface. **Nunca salve tokens JWT no storage do navegador**.

### Segurança & Defesa em Profundidade
- **Cookies Seguros**: Sessões operam exclusivamente com cookies `HttpOnly`, `Secure` e `SameSite=Lax`.
- **Prevenção a DoS no Bcrypt**: Valide tamanho de senha com limite máximo estrito: `8 <= len(password) <= 72`.
- **Anti-Enumeração**: Endpoints de recuperação de senha (`/forgot-password`) respondem com `200 OK` genérico e tempo constante.
- **Hashing de Tokens**: Tokens de reset de senha e verificação de e-mail devem ser armazenados exclusivamente como hash SHA-256 no banco de dados.
- **Sanitização**: Sanitização rigorosa de cabeçalhos SMTP contra CRLF Injection (CWE-93) e validação de paths contra Path Traversal (`../`).
- **Segredos & Infraestrutura**: Nunca exponha IPs de servidores, chaves RSA ou credenciais no repositório.

---

## 🚀 Submissão de Pull Requests

1. Certifique-se de que a branch está perfeitamente rebaseada com a `main` remota:
   ```bash
   git fetch origin main
   git rebase origin/main
   ```
2. Execute a verificação completa e garanta 100% de sucesso:
   ```bash
   ./scripts/check.sh all
   ```
3. Faça o push para sua branch remota:
   ```bash
   git push -u origin <nome-da-sua-branch>
   ```
4. Abra o Pull Request apontando para a base `main`.
5. Preencha detalhadamente o [Pull Request Template](.github/pull_request_template.md), incluindo o resumo, vínculo de issues e evidências visuais (se houver alteração de interface).

---

## 💡 Reportando Problemas & Sugerindo Funcionalidades

- **Problemas de Segurança**: Consulte a nossa [Política de Segurança](.github/SECURITY.md) para reporte responsável e confidencial (via GitHub Security Advisories privados). **Não abra issues públicas para vulnerabilidades de segurança**.
- **Relato de Bugs**: Utilize o template de [Bug Report](.github/ISSUE_TEMPLATE/bug_report.yml) fornecendo contexto, passos para reprodução e evidências.
- **Sugestão de Features**: Utilize o template de [Feature Request](.github/ISSUE_TEMPLATE/feature_request.yml) estruturado no formato User Story com critérios de aceite e camadas impactadas.

---

Agradecemos sua colaboração para tornar o CVMC cada vez melhor! 🚀
