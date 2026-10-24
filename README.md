# 🚗 CVMC - Como Vai Meu Carro

<div align="center">

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![React](https://img.shields.io/badge/React-19-61DAFB?style=for-the-badge&logo=react&logoColor=black)
![TypeScript](https://img.shields.io/badge/TypeScript-7.x-3178C6?style=for-the-badge&logo=typescript&logoColor=white)
![Material UI](https://img.shields.io/badge/Material_UI-v9-007FFF?style=for-the-badge&logo=mui&logoColor=white)
![Capacitor](https://img.shields.io/badge/Capacitor-v8-119EFF?style=for-the-badge&logo=capacitor&logoColor=white)
![CI Backend](https://img.shields.io/github/actions/workflow/status/yurinogueira/CVMC/backend.yml?branch=main&label=CI%20Backend&style=for-the-badge)
![CI Frontend](https://img.shields.io/github/actions/workflow/status/yurinogueira/CVMC/frontend.yml?branch=main&label=CI%20Frontend&style=for-the-badge)
![CI Docs](https://img.shields.io/github/actions/workflow/status/yurinogueira/CVMC/docs.yml?branch=main&label=CI%20Docs&style=for-the-badge)
![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)

**Plataforma moderna para gestão de frotas e veículos particulares, controle financeiro, histórico de manutenções, abastecimentos e monitoramento preventivo.**

[🌐 Acessar Web App](https://cvmc.yurinogueira.dev.br) • [📚 Documentação Técnica (OKF)](docs/index.md)

</div>

---

## 📌 Links de Produção & Artefatos

| Serviço / Recurso | Destino | Descrição |
| :--- | :--- | :--- |
| **🌐 Web App** | [cvmc.yurinogueira.dev.br](https://cvmc.yurinogueira.dev.br) | Aplicação SPA em produção (Cloudflare DNS/WAF + Caddy + OCI) |
| **📚 Base de Conhecimento** | [Catálogo Canônico (OKF)](docs/index.md) | Documentação de arquitetura, subdomínios DDD, runbooks e ADRs |
| **📱 Aplicativo Android** | [release-apk/](release-apk/) | Binários APK de release e debug gerados via Capacitor |
| **📖 Documentação Swagger** | `/swagger/index.html` | Especificação interativa OpenAPI dos endpoints REST da API |

---

## 🌟 Principais Funcionalidades

- 🔐 **Autenticação Segura & RBAC**:
  - Sessões protegidas via cookies `HttpOnly` stateless com suporte a refresh token e isolamento total de credenciais no frontend.
  - Suporte a **Login com Google** (OAuth2 web e autenticação nativa Android via Capacitor) com vinculação automática de contas.
  - Gestão de perfil do condutor e controle de limites de veículos cadastrados.
- 🚘 **Gestão de Veículos & Identidade Visual**:
  - Cadastro detalhado de veículos: placa, marca, modelo, ano de fabricação, combustível e quilometragem atual.
  - Upload seguro de fotos personalizadas com isolamento de storage.
  - Renders visuais vetoriais (SVG) dinâmicos adaptados por categoria/tipo de carroceria.
- 📊 **Integração com a Tabela FIPE**:
  - Consulta oficial em tempo real de marcas, modelos, anos de fabricação e valor de mercado atualizado.
  - Camada de **cache persistido no MongoDB** com indexação automática para resposta instantânea e resiliência offline.
- 🛠️ **Histórico & Ordens de Manutenção**:
  - Registro completo de manutenções preventivas e corretivas com discriminação de peças, serviços, oficinas e valores.
  - Acompanhamento cronológico e por odômetro da vida útil do automóvel.
- ⏰ **Monitoramento & Alertas Preventivos**:
  - Acompanhamento de revisões futuras e alertas com base na quilometragem e periodicidade temporal.
  - Infraestrutura de **Notificações Push (FCM - Firebase Cloud Messaging)** para alertas no dispositivo do usuário.
- ⛽ **Gestão de Combustível & Custos Operacionais**:
  - Registro de abastecimentos com cálculo automático de eficiência e consumo médio (`km/l`).
  - Consolidação de custos operacionais totais (manutenções + serviços + abastecimentos), fornecendo visão real do custo total de propriedade (TCO).
- 📎 **Armazenamento de Anexos e Comprovantes**:
  - Upload e gestão de notas fiscais, ordens de serviço e laudos técnicos.
  - Validação rigorosa contra MIME types maliciosos, controle de tamanho e prevenção contra Path Traversal e Stored XSS.
- 📱 **Multiplataforma (Web SPA + Android Nativo)**:
  - Aplicação Web responsiva acessível via desktop e smartphones.
  - Aplicativo nativo Android compilado via **Capacitor 8** com suporte a splash screen, safe area insets, status bar estilizada e biometria/social login.
- 🛡️ **Segurança Defensiva em Camadas**:
  - CORS restrito, rate limiting por IP confiável, mitigação de timing attacks em senhas e headers de proteção HTTP.

---

## 🏗️ Arquitetura & Stack Tecnológica

O CVMC segue uma arquitetura moderna orientada a microsserviços conceituais e desacoplamento rigoroso:

### Backend (Go 1.26)
- **Linguagem & Filosofia**: Go 1.26 estruturado sob os princípios de **Clean Architecture** e **Domain-Driven Design (DDD)**:
  - `internal/domain/`: Entidades puras, regras de negócio invariantes e erros de domínio.
  - `internal/application/`: Casos de uso (*use cases*) e interfaces/portas (*ports*) para repositórios e serviços.
  - `internal/infrastructure/`: Adaptadores técnicos (MongoDB driver v2, Object Storage, integração FIPE, hash bcrypt, Firebase).
  - `internal/interfaces/`: Handlers HTTP/REST, middlewares de segurança e especificações Swagger.
- **Banco de Dados**: MongoDB (Atlas e local) utilizando o driver oficial `go.mongodb.org/mongo-driver/v2` com segregação lógica por usuário (`user_id`) e índices compostos.
- **Documentação de API**: OpenAPI/Swagger interativo gerado automaticamente via `swag`.

### Frontend (React 19 + TypeScript 7)
- **Ecossistema**: React 19, TypeScript 7.x, Vite 8 e React Router v7.
- **Padrão Arquitetural**: *Feature-Based Architecture*, com cada módulo de negócio encapsulando suas próprias telas, componentes, estado e serviços.
- **Design System & Acessibilidade**: Material UI v9 (MUI) com paleta balanceada (azul, ciano e esmeralda) e conformidade com WCAG AA.
- **Gerenciamento de Estado**: Zustand v5 para gerenciamento reativo e descentralizado do estado global da aplicação.
- **Qualidade & Linters**: **Oxlint** (linter em Rust ultrarrápido) e Prettier para padronização de código, com testes unitários via Vitest e React Testing Library.

### Mobile (Capacitor 8 / Android)
- **Runtime**: `@capacitor/core` e `@capacitor/android` v8 integrando a SPA como aplicação nativa Android.
- **Plugins Nativos**:
  - `@capacitor/push-notifications`: Recepção de alertas e mensagens push via Firebase Cloud Messaging.
  - `@capgo/capacitor-social-login`: Autenticação Google nativa com credenciais do dispositivo.
  - `@capacitor/status-bar` & `@capacitor/splash-screen`: Controle visual e retenção de splash nativa.

### Infraestrutura, DevOps & CI/CD
- **Containerização**: Docker e Docker Compose para execução e paridade total de ambiente de desenvolvimento.
- **Hospedagem & Nuvem**: Oracle Cloud Infrastructure (OCI Compute Always Free) com binário estático Go gerenciado via `systemd`.
- **Borda & Proxy Reverso**: Caddy com renovação automática de certificados TLS/HTTPS, protegido por DNS e WAF da Cloudflare.
- **Infraestrutura como Código (IaC)**: Terraform gerenciando recursos declarativos na nuvem (`terraform/`).
- **Automação de CI/CD**: GitHub Actions cobrindo compilação, linters, testes unitários, validação da base de conhecimento (OKF) e pipelines de deploy.

---

## 📂 Estrutura do Repositório

```text
CVMC/
├── .agents/                  # Skills, padrões e diretrizes especializadas para agentes de IA
├── .github/                  # Pipelines de CI/CD (Actions), templates de issue e PR
├── backend/                  # API REST em Go 1.26 (Clean Architecture & DDD)
│   ├── cmd/                  # Pontos de entrada executáveis (servidor HTTP)
│   ├── docs/                 # Documentação OpenAPI/Swagger gerada
│   └── internal/             # Núcleo da API: domínio, casos de uso, portas, infra e REST
├── docs/                     # Catálogo Canônico OKF (Open Knowledge Format & LLM Wiki)
│   ├── adrs/                 # Registros de Decisões de Arquitetura (ADRs)
│   ├── architecture/         # Especificações de arquitetura, segurança e persistência
│   ├── domain/               # Subdomínios DDD (Car, Maintenance, Fuel, Fipe, User, etc.)
│   ├── frontend/             # Padrões visuais, roteamento e stores Zustand
│   ├── operations/           # Manuais de deploy, infraestrutura, CI/CD, mobile e Runbooks
│   ├── logs/                 # Trilha de auditoria datada (AAAA-MM-DD.md)
│   ├── index.md              # Ponto de entrada central do conhecimento (Catálogo Canônico)
│   └── log.md                # Linha do tempo consolidada de intervenções
├── frontend/                 # Aplicação Web SPA e Mobile em React 19 + TypeScript 7
│   ├── android/              # Projeto nativo Android integrado via Capacitor 8
│   └── src/
│       ├── features/         # Módulos de domínio (cars, maintenances, fuel, fipe, auth, profile)
│       ├── layouts/          # Casca estrutural (Sidebar, Topbar, AppLayout)
│       ├── routes/           # Definição e proteção de rotas públicas e autenticadas
│       └── services/         # Clientes de API Axios e integrações
├── deploy/                   # Configurações de serviços de produção, Caddy e systemd
├── release-apk/              # Binários APK de produção e depuração do app Android
├── scripts/                  # Scripts utilitários de build, lint, testes e manutenção
├── terraform/                # Especificações de Infraestrutura como Código (IaC)
└── docker-compose.yml        # Orquestração do ambiente local (MongoDB + API + Frontend)
```

---

## 🚀 Como Executar Localmente

### Pré-requisitos
- [Docker](https://www.docker.com/) e Docker Compose instalados **OU**
- [Go 1.26+](https://golang.org/) e [Node.js 20+](https://nodejs.org/) instalados no ambiente local.

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

Inicie todos os serviços locais (MongoDB, Backend Go e Frontend React):

```bash
./scripts/dev.sh start
```

Ou diretamente pelo Docker Compose:

```bash
docker compose up -d
```

- **Frontend Web**: `http://localhost:5173`
- **Backend API**: `http://localhost:8080`
- **Swagger UI**: `http://localhost:8080/swagger/index.html`

Para monitorar logs ou verificar status:

```bash
./scripts/dev.sh status
./scripts/dev.sh logs
```

Para encerrar os serviços:

```bash
./scripts/dev.sh stop
```

---

## 📱 Executando e Compilando o App Mobile (Android)

O projeto mobile utiliza o Capacitor para compilar a SPA React como aplicativo nativo Android:

```bash
# 1. Navegue para o diretório frontend
cd frontend

# 2. Instale as dependências e compile a aplicação
npm ci
npm run build

# 3. Sincronize os arquivos web com o projeto Android nativo
npm run cap:sync

# 4. Abra o projeto no Android Studio para executar em emulador ou dispositivo físico
npm run cap:open
```

> Para mais detalhes operacionais de compilação e geração de APK/AAB assinado, consulte o [Guia de Build Mobile](docs/operations/mobile-build.md).

---

## 💻 Desenvolvimento & Scripts Úteis

O repositório disponibiliza comandos unificados na pasta `scripts/` para garantir velocidade, integridade e conformidade de código:

| Comando | Finalidade |
| :--- | :--- |
| `./scripts/check.sh all` | Executa verificação total: validação de documentação (OKF), backend, frontend e Terraform |
| `./scripts/check.sh backend` | Valida apenas o backend Go (`go vet` e testes automatizados `go test`) |
| `./scripts/check.sh frontend` | Valida o frontend React (`tsc -b`, `oxlint`, formatação Prettier e testes Vitest) |
| `./scripts/check.sh terraform` | Valida a formatação de Infraestrutura como Código (`terraform fmt -check`) |
| `./scripts/check.sh docs` | Valida integridade referencial de links e sincronismo da documentação canônica (OKF) |
| `./scripts/fix.sh` | Formata automaticamente o código Go (`go fmt`), Frontend (`prettier`) e Terraform |
| `./scripts/swagger.sh` | Regenera a documentação OpenAPI/Swagger a partir das anotações dos handlers HTTP |
| `./scripts/dev.sh <start\|stop\|status\|logs>` | Gerencia o ciclo de vida dos containers locais do Docker Compose |

---

## 📚 Base de Conhecimento Canônica (OKF)

O CVMC adota o padrão **Open Knowledge Format (OKF)** para governança de documentação técnica, servindo como a fonte única da verdade (*Single Source of Truth*) para engenheiros e agentes autônomos de IA.

A base é estruturada sob o princípio da **Divulgação Progressiva** (*Progressive Disclosure*):
- 🧭 **Ponto de Partida**: Consulte o [Catálogo Canônico & LLM Wiki](docs/index.md) para navegação estruturada.
- 🏛️ **Subdomínios DDD**: Especificações funcionais em [docs/domain/](docs/domain/).
- 🏗️ **Arquitetura & Segurança**: Decisões técnicas e mitigação de vulnerabilidades em [docs/architecture/](docs/architecture/).
- 📋 **Runbooks**: Guias práticos de contingência operacional e procedimentos de release em [docs/operations/](docs/operations/).

---

## 🔒 Segurança

Diretrizes de segurança, política de suporte e canais para reporte responsável de vulnerabilidades estão documentados em [SECURITY.md](.github/SECURITY.md).

---

## 🤝 Contribuição

Contribuições são bem-vindas! Consulte o guia de contribuição em [CONTRIBUTING.md](CONTRIBUTING.md) e nosso [Código de Conduta](CODE_OF_CONDUCT.md).

---

## 📄 Licença

Este projeto está sob a licença **MIT**. Consulte o arquivo [LICENSE](LICENSE) para mais detalhes.

---

<div align="center">
Desenvolvido por <a href="https://github.com/yurinogueira">Yuri Nogueira</a>
</div>
