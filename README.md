# 🚗 CVMC - Como Vai Meu Carro

<div align="center">

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![React](https://img.shields.io/badge/React-19-61DAFB?style=for-the-badge&logo=react&logoColor=black)
![TypeScript](https://img.shields.io/badge/TypeScript-5.x-3178C6?style=for-the-badge&logo=typescript&logoColor=white)
![Material UI](https://img.shields.io/badge/Material_UI-v6-007FFF?style=for-the-badge&logo=mui&logoColor=white)
![CI Backend](https://img.shields.io/github/actions/workflow/status/yurinogueira/CVMC/backend.yml?branch=main&label=CI%20Backend&style=for-the-badge)
![CI Frontend](https://img.shields.io/github/actions/workflow/status/yurinogueira/CVMC/frontend.yml?branch=main&label=CI%20Frontend&style=for-the-badge)
![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)

**Plataforma moderna para gestão, controle de custos e histórico de manutenção veicular.**

[Acessar Web App](https://cvmc.yurinogueira.dev.br)

</div>

---

## 📌 Links de Produção

| Serviço | URL | Descrição |
| :--- | :--- | :--- |
| **🌐 Web App** | [cvmc.yurinogueira.dev.br](https://cvmc.yurinogueira.dev.br) | Aplicação SPA em produção (GitHub Pages + Cloudflare) |

---

## 🌟 Funcionalidades

- 🔐 **Autenticação Segura & Controle de Acesso (RBAC)**: Sessões protegidas via cookies `HttpOnly` com isolamento de credenciais no frontend, suporte a refresh tokens e permissões por perfil.
- 🚘 **Gestão de Veículos**: Cadastro, busca e atualização de veículos com dados técnicos, identificação, quilometragem e especificações de combustível.
- 🛠️ **Histórico e Planejamento de Manutenções**: Registro detalhado de intervenções preventivas e corretivas, custos de peças e serviços, oficinas e acompanhamento por odômetro.
- 💰 **Controle Financeiro e Custos Operacionais**: Consolidação das despesas operacionais do veículo (manutenções, serviços e abastecimentos), oferecendo visão clara do custo total de propriedade.
- 📎 **Armazenamento de Anexos e Comprovantes**: Gestão de notas fiscais e relatórios técnicos com validação de formato e isolamento seguro de arquivos.
- 📊 **Painel Analítico e Indicadores**: Visualização de métricas de custos, próximos alertas de revisão e saúde veicular.
- 🛡️ **Segurança em Camadas**: Proteção contínua contra vetores de ataque comuns (CORS estrito, rate limiting por IP, validação de integridade de payloads e cabeçalhos de segurança HTTP).

---

## 🏗️ Arquitetura & Stack Tecnológica

### Backend (Go)
- **Linguagem & Padrões**: Go, adotando **Clean Architecture** e princípios de **Domain-Driven Design (DDD)** para desacoplamento de regras de negócio:
  - `internal/domain/`: Entidades puras, regras de negócio e tipos de domínio invariantes.
  - `internal/application/`: Casos de uso (*use cases*) e contratos de portas (*ports*) para serviços e repositórios.
  - `internal/infrastructure/`: Adaptadores técnicos de banco de dados, criptografia e integrações externas.
  - `internal/interfaces/`: Handlers HTTP/REST, middlewares de segurança e documentação OpenAPI/Swagger.
- **Banco de Dados**: MongoDB (armazenamento baseado em documentos para dados transacionais e analíticos).

### Frontend (React)
- **Ecossistema**: React, TypeScript, Vite e React Router.
- **Padrão de Arquitetura**: Arquitetura modular orientada a domínios (*Feature-Based Architecture*), na qual cada módulo de funcionalidade encapsula suas telas, componentes, estado e serviços específicos.
- **Interface do Usuário**: Material UI (MUI) com sistema de design responsivo e componentes acessíveis.
- **Gerenciamento de Estado**: Estado de UI e sessão desacoplados utilizando Zustand.
- **Comunicação HTTP**: Cliente Axios padronizado para transporte automático e seguro de credenciais via cookies `HttpOnly`.

### Infraestrutura, DevOps & CI/CD
- **Containerização**: Docker e Docker Compose para execução e paridade de ambiente local.
- **Nuvem & Rede**: Infraestrutura distribuída em nuvem com proxy reverso e proteção de borda.
- **Infraestrutura como Código (IaC)**: Terraform para provisionamento declarativo e reprodutível de recursos e DNS.
- **Integração & Entrega Contínua (CI/CD)**: GitHub Actions cobrindo análise estática, testes unitários, formatação e pipelines de deploy.

---

## 📂 Estrutura do Repositório

A organização do repositório é orientada à separação clara de responsabilidades entre as camadas de aplicação, infraestrutura e governança do código:

```text
CVMC/
├── .github/                  # Pipelines de CI/CD, templates de issues e automações
├── backend/                  # API REST em Go (Clean Architecture & DDD)
│   ├── cmd/                  # Pontos de entrada executáveis da aplicação
│   ├── docs/                 # Documentação Swagger / OpenAPI gerada
│   └── internal/             # Núcleo da API: domínio, casos de uso, infraestrutura e interfaces
├── frontend/                 # Aplicação Web SPA em React + TypeScript
│   └── src/
│       ├── features/         # Módulos de domínio da aplicação (feature-based)
│       ├── layouts/          # Casca estrutural da aplicação (Sidebar, Topbar, AppLayout)
│       ├── routes/           # Definição e proteção de rotas públicas e autenticadas
│       └── services/         # Clientes de integração de API e serviços globais
├── deploy/                   # Configurações de servidores web, proxies e templates de deploy
├── scripts/                  # Utilitários de desenvolvimento, validação e manutenção
├── terraform/                # Especificações de Infraestrutura como Código (IaC)
└── docker-compose.yml        # Orquestração do ambiente de desenvolvimento local
```

> **Design de Extensibilidade**: Novas funcionalidades de negócio são adicionadas de forma modular e isolada dentro de `frontend/src/features/<nome-da-feature>/` e orquestradas em `backend/internal/application/usecase/`, garantindo evolução contínua sem necessidade de reestruturação do projeto.

---

## 🚀 Como Executar Localmente

### Pré-requisitos
- [Docker](https://www.docker.com/) e Docker Compose instalados **OU**
- [Go](https://golang.org/) e [Node.js](https://nodejs.org/) instalados no ambiente local.

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
| `./scripts/check.sh all` | Executa validação estática completa: testes, linters, type-checking e formatação |
| `./scripts/check.sh backend` | Valida apenas o backend Go (`go vet` e `go test`) de forma concisa |
| `./scripts/check.sh frontend` | Valida apenas o frontend React (`tsc`, `eslint`, `vitest`) |
| `./scripts/fix.sh` | Formata automaticamente o código Go (`go fmt`) e Frontend (`prettier`, `eslint --fix`) |
| `./scripts/swagger.sh` | Regenera a documentação OpenAPI/Swagger a partir das anotações dos handlers |
| `./scripts/dev.sh start\|stop\|status` | Gerencia o ciclo de vida dos containers Docker locais |

---

## 🔒 Segurança

Diretrizes de segurança e canais para reporte responsável de vulnerabilidades estão documentados em [SECURITY.md](.github/SECURITY.md).

---

## 📄 Licença

Este projeto está sob a licença **MIT**. Consulte o arquivo [LICENSE](LICENSE) para mais detalhes.

---

<div align="center">
Desenvolvido por <a href="https://github.com/yurinogueira">Yuri Nogueira</a>
</div>
