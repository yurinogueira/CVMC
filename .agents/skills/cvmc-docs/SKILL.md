---
name: cvmc-docs
description: >-
  Diretrizes e automações para leitura, criação e melhorias contínuas da documentação
  canônica (OKF) no CVMC, orientando desenvolvedores e agentes autônomos a consultar
  o catálogo docs/index.md, manter a sincronia entre código e docs e validar via scripts.
---

# Skill: Gestão e Governança da Documentação Canônica (OKF) — CVMC

Esta skill estabelece os padrões e procedimentos obrigatórios para **ler**, **criar** e **evoluir** a documentação canônica do projeto **CVMC (Como Vai Meu Carro)**, estruturada sob o padrão **Open Knowledge Format (OKF)** do Google Cloud como um *LLM Wiki* em `docs/`.

---

## 🎯 Por Que a Documentação Canônica é Obrigatória?

No ecossistema do CVMC, a base em `docs/` é a **Única Fonte Canônica da Verdade**. Ela elimina alucinações de modelos de IA, reduz drasticamente o consumo de tokens na janela de contexto e garante que novos desenvolvedores e agentes compreendam as invariantes arquiteturais, regras de domínio DDD e padrões de segurança sem precisar garimpar o código-fonte.

---

## 📖 1. Protocolo de Leitura (Divulgação Progressiva & Economia de Tokens)

> [!IMPORTANT]
> **Regra de Ouro**: **NUNCA** leia arquivos de documentação aleatoriamente nem carregue múltiplos documentos de uma vez.
> Pratique sempre a **Divulgação Progressiva (*Progressive Disclosure*)**:

1. **Passo 1 (Localização no Catálogo)**:
   - Abra o catálogo mestre `docs/index.md` utilizando `view_file`.
   - Localize o tópico desejado pelo mapa conceitual ou pelo índice estruturado (`architecture/`, `domain/`, `frontend/`, `operations/`, `adrs/`).
2. **Passo 2 (Leitura Cirúrgica)**:
   - Abra diretamente o documento mapeado (ex: `docs/domain/car.md` ou `docs/architecture/auth-and-security.md`).
   - Leia as regras de negócio, interfaces e referências cruzadas específicas da tarefa em andamento.
3. **Passo 3 (Navegação via Grafo)**:
   - Siga os links relativos citados no documento para aprofundar em conceitos correlacionados, evitando buscas cegas no repositório.

---

## ✍️ 2. Protocolo para Criação de Novos Documentos

Ao especificar um novo subdomínio, componente arquitetural, manual de infraestrutura ou Registro de Decisão Arquitetural (ADR):

### 1. Nomenclatura e Localização
- Salve o arquivo na subpasta semântica correta sob `docs/`:
  - `docs/architecture/<topico>.md`: Decisões globais de arquitetura, segurança e armazenamento.
  - `docs/domain/<subdominio>.md`: Regras de negócio, entidades e agregados DDD do backend.
  - `docs/frontend/<tema>.md`: Padrões de tela, estado Zustand, rotas e componentes MUI v6.
  - `docs/operations/<assunto>.md` ou `docs/operations/runbooks/<procedimento>.md`: Manuais de infraestrutura, Docker, CI/CD e runbooks de contingência.
  - `docs/adrs/<numero-sequencial>-<titulo-curto>.md`: Registros de decisão arquitetural imutáveis.
- Utilize nomes em minúsculas com hífen (`kebab-case`).

### 2. Cabeçalho YAML Frontmatter Obrigatório
Todo arquivo `.md` em `docs/` deve iniciar na linha 1 com delimitadores `---` e conter:

```yaml
---
type: domain # [index | log | architecture | domain | frontend | operations | adr | runbook | concept]
title: "Título Semântico e Claro"
description: "Resumo explicativo do propósito e escopo deste documento."
tags:
  - domain
  - car
  - ddd
resource: backend/internal/domain/car # Opcional: caminho do código-fonte correspondente
timestamp: 2026-10-10
---
```

#### Tipos Permitidos (`type`):
- `index`: Reservado para `docs/index.md`.
- `log`: Reservado para `docs/log.md` e arquivos diários sob `docs/logs/*.md`.
- `architecture`: Visão sistêmica, camadas, segurança e storage.
- `domain`: Regras de negócio, entidades e contratos DDD.
- `frontend`: Diretrizes de estado, componentes UI e rotas SPA.
- `operations`: Manuais de ambiente, deploy, Docker e pipelines.
- `runbook`: Procedimentos passo a passo de recuperação ou manutenção.
- `adr`: Architecture Decision Records.

### 3. Grafo de Links e Conexões Bidirecionais
- Todas as referências cruzadas devem utilizar caminhos relativos válidos do Markdown (ex: `[Visão Geral](../architecture/overview.md)` ou `[ADR 0001](../adrs/0001-clean-architecture-go.md)`).
- Todo documento criado deve conter um link de retorno para o catálogo mestre: `[Catálogo Canônico](../index.md)`.

### 4. Indexação Mandatória em `docs/index.md` e `docs/log.md`
- **Documentos conceituais e arquiteturais**: Todo novo documento sob `docs/` deve ter seu link e descrição inseridos em `docs/index.md`.
- **Arquivos diários de log (`docs/logs/AAAA-MM-DD.md`)**: Devem ser referenciados e indexados no catálogo central `docs/log.md`.

### 5. Registro Cronológico Particionado em `docs/logs/` e `docs/log.md`
- Toda adição, evolução, correção ou atualização de dependência **DEVE** ser registrada no arquivo diário correspondente à data atual (`docs/logs/AAAA-MM-DD.md`).
- **Se o arquivo do dia já existir**: Adicione uma nova seção de nível 2 (`## 🎯 <Título> (#PR/Issue)`) e insira um novo bullet no dia correspondente em `docs/log.md`.
- **Se for o primeiro registro do dia**: Crie `docs/logs/AAAA-MM-DD.md` com YAML frontmatter (`type: log`, `timestamp: AAAA-MM-DD`, etc.) e registre o novo dia no índice central `docs/log.md`.
- Formato obrigatório da seção no log diário:
  ```markdown
  ## 🎯 <Título Semântico da Mudança> (#<issue_ou_pr>)

  ### 📋 Resumo da Intervenção
  [Descrição sucinta da motivação, contexto de negócio/técnico e impacto gerado]

  ### 🛠️ Modificações por Camada
  - **Backend / Domínio**: [Detalhamento de structs, endpoints, use cases ou contratos]
  - **Frontend / UI**: [Detalhamento de telas, modais, componentes ou stores Zustand]
  - **Operações / Infra / CI**: [Detalhamento de dependências, Docker, workflows ou terraform]
  - **Documentos Canônicos Atualizados**:
    - [docs/domain/car.md](../domain/car.md)
    - [docs/architecture/auth-and-security.md](../architecture/auth-and-security.md)
  ```

---

## 🔄 3. Protocolo Universal de Manutenção Contínua ("Zero-Divergence OKF")

> [!CAUTION]
> **Regra Universal Inegociável**: **QUALQUER TAREFA** executada por um desenvolvedor ou agente autônomo — seja criar novas telas/regras, corrigir um bug sutil, alterar uma biblioteca, atualizar uma versão de runtime ou refatorar código — **EXIGE A ATUALIZAÇÃO DA DOCUMENTAÇÃO CANÔNICA E DO LOG DIÁRIO NO MESMO PR**.
> Nunca encerre uma tarefa sem atualizar `docs/` e `docs/logs/AAAA-MM-DD.md` (indexado em `docs/log.md`).

### Guia de Correspondência de Alterações:

1. **Novas Funcionalidades (`feat`)**:
   - Novas rotas, telas ou componentes: `docs/frontend/ui-components.md` e `docs/frontend/routing-and-rbac.md`.
   - Novas regras de negócio ou entidades: `docs/domain/<subdominio>.md`.
   - Novos fluxos de dados ou permissões: `docs/architecture/auth-and-security.md` ou `multitenancy-and-data.md`.
   - Registro datado: `docs/logs/AAAA-MM-DD.md` e sumário em `docs/log.md`.

2. **Correções de Defeitos e Segurança (`fix`)**:
   - Correção em autenticação, cookies, sessões ou roles: `docs/architecture/auth-and-security.md` e `docs/domain/user.md`.
   - Correção em cálculos, veículos ou manutenções: `docs/domain/<subdominio>.md` (ex: `docs/domain/car.md`, `maintenance.md`, `fuel.md`).
   - Correção de layout, responsividade ou acessibilidade: `docs/frontend/ui-components.md`.
   - Registro datado: `docs/logs/AAAA-MM-DD.md` e sumário em `docs/log.md` (descrever a causa raiz e a correção aplicada).

3. **Atualização de Versões e Dependências (`chore`/`deps`/`ci`)**:
   - Bumps de runtime (Go, Node): `docs/operations/docker-and-local-dev.md`, `docs/operations/ci-cd-pipelines.md`, `docs/index.md` e `.agents/skills/cvmc-dev/SKILL.md`.
   - Bumps de banco de dados (MongoDB): `docs/operations/docker-and-local-dev.md` e `docs/index.md`.
   - Bumps de bibliotecas frontend (React, Vite, MUI): `docs/frontend/ui-components.md` e `docs/operations/docker-and-local-dev.md`.
   - Bumps de Terraform ou nuvem: `docs/operations/deploy-and-infrastructure.md`.
   - Registro datado: `docs/logs/AAAA-MM-DD.md` e sumário em `docs/log.md` (especificar quais pacotes foram atualizados e motivação).

4. **Decisões Arquiteturais e Refatorações (`refactor`/`adr`)**:
   - Nova decisão técnica estrutural: novo arquivo em `docs/adrs/XXXX-<nome>.md`, indexado em `docs/index.md`.
   - Mudança estrutural de serviços ou storage: `docs/architecture/overview.md` ou `storage-and-media.md`.
   - Registro datado: `docs/logs/AAAA-MM-DD.md` e sumário em `docs/log.md`.

---

## 🧪 4. Validação Automatizada Pré-Commit

Antes de submeter alterações ou abrir Pull Requests, execute sempre o validador automatizado:

```bash
# Validação específica da documentação (execução em milissegundos):
./scripts/check-docs.sh

# Ou via o script unificado de checagem do CVMC:
./scripts/check.sh docs

# Validação completa de todo o repositório (Go + Frontend + Terraform + Docs):
./scripts/check.sh all
```

O validador verifica:
1. Conformidade estrita da sintaxe YAML frontmatter.
2. Presença dos campos obrigatórios (`type`, `title`, `description`).
3. Uso de valores permitidos no campo `type`.
4. Resolução de 100% dos links Markdown locais (zero links quebrados).
5. Indexação de todos os arquivos no catálogo `docs/index.md`.
