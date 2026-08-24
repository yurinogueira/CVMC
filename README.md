# 🚗 CVMC - Como Vai Meu Carro

<div align="center">

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![React](https://img.shields.io/badge/React-19-61DAFB?style=for-the-badge&logo=react&logoColor=black)
![TypeScript](https://img.shields.io/badge/TypeScript-5.x-3178C6?style=for-the-badge&logo=typescript&logoColor=white)
![Material UI](https://img.shields.io/badge/Material_UI-v6-007FFF?style=for-the-badge&logo=mui&logoColor=white)
![CI Backend](https://img.shields.io/github/actions/workflow/status/yurinogueira/CVMC/ci-backend.yml?branch=main&label=CI%20Backend&style=for-the-badge)
![CI Frontend](https://img.shields.io/github/actions/workflow/status/yurinogueira/CVMC/ci-frontend.yml?branch=main&label=CI%20Frontend&style=for-the-badge)
![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)

**Plataforma moderna para gestão, controle de custos e histórico de manutenção veicular.**

[Acessar Web App](https://cvmc.yurinogueira.dev.br)

</div>

---

## 📌 Links de Produção

| Serviço | URL | Descrição |
| :--- | :--- | :--- |
| **🌐 Web App (Frontend)** | [cvmc.yurinogueira.dev.br](https://cvmc.yurinogueira.dev.br) | Aplicação SPA em produção (GitHub Pages + Cloudflare) |
| **🔌 API Backend** | `api-cvmc.yurinogueira.dev.br` | API RESTful Go hospedada em Oracle Cloud (OCI) |
| **📖 Swagger Docs** | `api-cvmc.yurinogueira.dev.br/swagger/index.html` | Documentação OpenAPI interativa dos endpoints |

---

## 🌟 Funcionalidades

- 🔐 **Autenticação Segura & RBAC**: Sessões gerenciadas via cookies `HttpOnly` com suporte a refresh token, sem expor tokens no frontend e com controle de acesso baseado em papéis (Admin/User).
- 🚘 **Gestão Completa de Veículos**: Cadastro, listagem, busca e edição de automóveis com placa, marca, modelo, ano de fabricação/modelo, quilometragem atual e tipo de combustível.
- 🛠️ **Histórico de Manutenções**: Registro minucioso de revisões preventivas e corretivas, datas, oficinas, custos de peças e mão de obra, além de alertas por quilometragem.
- 📎 **Upload de Comprovantes**: Anexação de notas fiscais, relatórios e recibos de serviços com validação de tipos MIME e isolamento de storage.
- 📊 **Dashboard Analítico**: Indicadores de saúde veicular, métricas de custos consolidados, próximos alertas de revisão e atalhos rápidos.
- 🛡️ **Segurança em Camadas**: CORS estrito, rate limiting por IP, validação de payloads, cabeçalhos de segurança HTTP (HSTS, CSP, X-Frame-Options) e sanitização contra Path Traversal.

---

## 🏗️ Arquitetura & Stack Tecnológica

### Backend (Go)
- **Linguagem**: Go 1.25
- **Arquitetura**: Clean Architecture + Domain-Driven Design (DDD)
  - `internal/domain/`: Entidades de negócio puras e regras de domínio.
  - `internal/application/ports/`: Contratos e interfaces de repositórios/serviços.
  - `internal/application/usecase/`: Casos de uso e orquestração de lógica de negócio.
  - `internal/infrastructure/`: Implementações de banco de dados (MongoDB), JWT, bcrypt e storage.
  - `internal/interfaces/rest/`: Handlers HTTP REST, roteamento e documentação Swagger com Swaggo.
- **Banco de Dados**: MongoDB 8 (Local) e MongoDB Atlas (Produção).

### Frontend (React)
- **Framework & Ferramentas**: React 19, TypeScript, Vite, React Router v7.
- **UI & Estilização**: Material UI (MUI v6) com `@mui/icons-material` e paleta de cores personalizada.
- **Gerenciamento de Estado**: Zustand (gerenciamento desacoplado de estado da UI e perfil).
- **Comunicação HTTP**: Axios configurado com `withCredentials: true` para transporte automático de cookies `HttpOnly`.

### Infraestrutura & DevOps
- **Containerização**: Docker e Docker Compose para ambiente de desenvolvimento reprodutível.
- **Nuvem & Hospedagem**: Oracle Cloud Infrastructure (OCI Compute Instance, VCN, Reserved IP).
- **Edge & DNS**: Cloudflare com SSL/TLS Full, proteção DDoS e proxy DNS.
- **Infraestrutura como Código**: Terraform para provisionamento automatizado de recursos OCI, MongoDB Atlas e Cloudflare DNS.
- **CI/CD**: GitHub Actions para validação de testes, linters, build e deploy automático.

---

## 📂 Estrutura do Repositório

```text
CVMC/
├── backend/                  # API REST em Go (Clean Architecture)
│   ├── cmd/api/              # Ponto de entrada (main.go)
│   ├── docs/                 # Documentação Swagger gerada
│   ├── internal/             # Domain, Use Cases, Infrastructure, Handlers
│   └── Dockerfile            # Container de produção Go
├── frontend/                 # Aplicação SPA em React 19 + Vite
│   ├── src/
│   │   ├── features/         # Módulos: auth, cars, dashboard, maintenance
│   │   ├── layouts/          # Shell da aplicação (Sidebar, Topbar, AppLayout)
│   │   ├── routes/           # Rotas públicas e rotas protegidas
│   │   └── services/         # Clientes de API e storage seguro
├── deploy/                   # Configurações de Caddy, systemd e exemplos de env
├── scripts/                  # Scripts otimizados de checagem, build e dev
├── terraform/                # Definições IaC (OCI, Cloudflare, Mongo Atlas)
└── docker-compose.yml        # Stack completa para desenvolvimento local
```

---

## 🚀 Como Executar Localmente

### Pré-requisitos
- Docker e Docker Compose instalados **OU**
- Go 1.25+ e Node.js 22+

### 1. Clonar o Repositório

```bash
git clone https://github.com/yurinogueira/CVMC.git
cd CVMC
```

### 2. Configurar Variáveis de Ambiente

Copie o arquivo de exemplo para `.env`:

```bash
cp .env.example .env
```

### 3. Executar com Docker Compose (Recomendado)

Inicie todos os serviços (MongoDB, Backend Go e Frontend React):

```bash
./scripts/dev.sh start
```

Ou diretamente pelo Docker Compose:

```bash
docker compose up -d
```

- **Frontend**: `http://localhost:5173`
- **Backend API**: `http://localhost:8080`
- **Swagger UI**: `http://localhost:8080/swagger/index.html`

Para verificar o status ou logs dos containers:

```bash
./scripts/dev.sh status
./scripts/dev.sh logs
```

Para encerrar os serviços:

```bash
./scripts/dev.sh stop
```

---

## 💻 Desenvolvimento & Scripts Úteis

O repositório inclui utilitários em `scripts/` para desenvolvimento ágil e validação de código:

| Script | Finalidade |
| :--- | :--- |
| `./scripts/check.sh all` | Executa testes unitários, type-checking, linters e formatação em todo o projeto |
| `./scripts/check.sh backend` | Valida apenas o backend Go (`go vet` e `go test`) de forma concisa |
| `./scripts/check.sh frontend` | Valida apenas o frontend React (`tsc`, `eslint`, `vitest`) |
| `./scripts/fix.sh` | Formata automaticamente o código Go (`go fmt`) e Frontend (`prettier`, `eslint --fix`) |
| `./scripts/swagger.sh` | Regenera a documentação OpenAPI/Swagger a partir das anotações dos handlers |
| `./scripts/dev.sh start\|stop\|status` | Gerencia o ciclo de vida dos containers Docker |

---

## 🔒 Segurança

Diretrizes de segurança e canais para reporte responsável de vulnerabilidades estão documentados em `SECURITY.md`.

---

## 📄 Licença

Este projeto está sob a licença **MIT**. Consulte o arquivo `LICENSE` para mais detalhes.

---

<div align="center">
Desenvolvido por Yuri Nogueira
</div>
