---
name: cvmc-issues
description: >-
  Diretrizes e automações para criar issues de alta precisão (Bugs e Features)
  no projeto CVMC, formatadas de forma acionável para agentes de IA e desenvolvedores.
---

# Skill: Criação e Gestão de Issues no CVMC

Esta skill define os padrões e procedimentos obrigatórios para investigar, estruturar e publicar **Issues de Bug** e **Issues de Feature** no projeto **CVMC (Como Vai Meu Carro)**.

Issues bem escritas garantem que agentes de IA e desenvolvedores possam implementar correções e novas funcionalidades diretamente, sem ambiguidades e alinhados à arquitetura do projeto.

---

## 🎯 Princípios Fundamentais para Criação de Issues

1. **Intenção de Especificação (User Story / Feature Issue)**: Quando o usuário solicitar "Crie uma feature...", "Especifique a user story...", "Levante os requisitos...", o objetivo principal é formular e registrar uma **Issue de alta precisão técnica e funcional**, detalhando os requisitos e critérios de aceite antes da codificação. A implementação ocorrerá posteriormente no ciclo de desenvolvimento (`cvmc-workflow`).
2. **Investigação Prévia Obrigatória**: Antes de redigir a issue, explore a base de código para mapear os arquivos, rotas, contratos, workflows ou componentes envolvidos.
3. **Precisão de Localização**: Toda issue deve conter caminhos de arquivo exatos (`backend/...`, `frontend/...`, `.github/workflows/...`, `terraform/...`), nomes de structs, funções, endpoints REST (`/api/v1/...`) ou jobs de pipeline.
4. **Título Semântico**: Use a convenção do Conventional Commits para títulos de issue:
   - Bugs: `fix(<modulo>): <descrição sucinta em minúsculas>` (ex: `fix(auth): cookie de sessão não enviado na rota de perfil`)
   - Features/Melhorias: `feat(<modulo>): <descrição sucinta em minúsculas>` (ex: `feat(workflows): unificar e padronizar pipelines de ci e deploy`)
   - Manutenções/Infra: `chore(<modulo>): <descrição sucinta>` ou `ci(<modulo>): <descrição sucinta>`
5. **Alinhamento Arquitetural**: Respeite os padrões documentados em `cvmc-dev`, `cvmc-security` e as diretrizes de CI/CD do repositório.

---

## 🐛 Fluxo para Criação de Issues de Bug

Ao relatar um defeito, siga rigorosamente a seguinte estrutura:

### 1. Investigação do Bug
- Identifique a causa raiz ou o ponto de falha navegando no repositório.
- Colete payloads de requisição/resposta, logs do backend (`./scripts/dev.sh logs backend`) ou erros no console do frontend.
- Identifique se o problema ocorre no backend, frontend, banco MongoDB, infraestrutura Docker, workflows de CI/CD ou documentação Swagger.

### 2. Estrutura Padrão do Relato de Bug

```markdown
### Descrição do Problema
[Descrição clara e objetiva do erro observado]

### Componente / Camada Afetada
[Backend | Frontend | CI/CD & Workflows | Infraestrutura / Docker | Banco MongoDB | Autenticação | Swagger]

### Onde o Problema Ocorre (Localização Técnica)
- **Rota / Endpoint / Workflow**: `METODO /api/v1/...` ou `.github/workflows/arquivo.yml`
- **Arquivo(s) Backend**: `backend/internal/.../arquivo.go` (função/método `NomeFuncao`)
- **Componente(s) Frontend**: `frontend/src/features/.../Componente.tsx`
- **Store / Serviço**: `frontend/src/features/.../store.ts`

### Comportamento Atual vs. Comportamento Esperado
- **Comportamento Atual**: [O que acontece de errado hoje]
- **Comportamento Esperado**: [Como o sistema deve se comportar corretamente]

### Passos para Reproduzir
1. [Passo 1]
2. [Passo 2]
3. [Passo 3]

### Evidências e Contexto Técnico
\`\`\`json
{
  "status": 500,
  "message": "exemplo de erro retornado"
}
\`\`\`

### Gravidade / Severidade
[🔴 Bloqueante | 🟠 Alta | 🟡 Média | 🟢 Baixa]
```

---

## ✨ Fluxo para Criação de Issues de Feature / User Story

Ao propor uma nova funcionalidade, melhoria técnica ou automação de infraestrutura/CI:

### 1. Levantamento Arquitetural & Escopo
- Mapeie as camadas impactadas:
  - **Backend**: Entidades (`internal/domain/`), Contratos (`internal/application/ports/`), Use Cases (`internal/application/usecase/`), Adapters (`internal/infrastructure/`), REST/Swagger (`internal/interfaces/rest/`).
  - **Frontend**: Telas e componentes (`frontend/src/features/`), Stores Zustand (`features/<f>/state/`), Cliente API (`services/api/client.ts`).
  - **CI/CD & DevOps**: Workflows (`.github/workflows/`), scripts auxiliares (`scripts/`), Terraform (`terraform/`), Docker (`docker-compose.yml`, `deploy/`).

### 2. Estrutura Padrão da Proposta de Feature / User Story

```markdown
### Visão Geral e Contexto de Negócio / Técnico
**Como** [tipo de usuário/desenvolvedor/mantenedor],
**Quero** [ação, capacidade ou automação desejada],
**Para que** [benefício, redução de duplicidade, ganho de confiabilidade ou valor entregue].

### Requisitos Funcionais e Critérios de Aceite
- [ ] [Critério 1: Regra ou comportamento específico]
- [ ] [Critério 2: Validação de entradas / triggers / concorrência]
- [ ] [Critério 3: Resposta da API, comportamento da interface ou resultado de pipeline]

### Camadas / Componentes Técnicos Impactados
- [ ] Backend - Domínio e Contratos (`internal/domain/`, `internal/application/ports/`)
- [ ] Backend - Casos de Uso e Testes (`internal/application/usecase/`)
- [ ] Backend - Interfaces REST e Swagger (`internal/interfaces/rest/`, `./scripts/swagger.sh`)
- [ ] Frontend - Telas e Componentes MUI (`frontend/src/features/`)
- [ ] CI/CD & Workflows - GitHub Actions (`.github/workflows/`)
- [ ] Infraestrutura & IaC - Terraform / Docker (`terraform/`, `deploy/`)

### Detalhes Técnicos e Arquitetura Proposta
- **Arquivos Criados / Modificados / Deletados**: Caminhos exatos no repositório.
- **Triggers / Endpoints / Estruturas**: Detalhes de execução, parâmetros, contratos DTO ou eventos.
- **Estratégia de Execução & Artefatos**: Reaproveitamento de compilação, jobs condicionais, dependências (`needs`).

### Considerações de Segurança, Performance e Concorrência
- Princípio do menor privilégio em permissões (ex: `permissions: contents: read`, tokens restritos, cookies `HttpOnly`).
- Concorrência e idempotência (`concurrency` com `cancel-in-progress` em PRs, integridade de branches).
- Paridade com scripts locais de verificação (`./scripts/check.sh all`).

### Checklist de Implementação
- [ ] [Passo 1: Criação/ajuste de arquivos estruturais]
- [ ] [Passo 2: Implementação dos fluxos principais / jobs / use cases]
- [ ] [Passo 3: Remoção de arquivos obsoletos ou código redundante]
- [ ] [Passo 4: Validação completa via `./scripts/check.sh all`]
```

---

## 🚀 Métodos de Publicação de Issues

O agente pode publicar as issues de forma automatizada no GitHub utilizando as ferramentas integradas:

### Opção 1: Via GitHub MCP Tool (Preferencial em sessões assistidas)
Utilize a ferramenta `call_mcp_tool` chamando o servidor `github` e a tool `create_issue`:

```json
{
  "ServerName": "github",
  "ToolName": "create_issue",
  "Arguments": {
    "owner": "yurinogueira",
    "repo": "CVMC",
    "title": "feat(workflows): unificar e padronizar pipelines de ci e deploy",
    "body": "### Visão Geral...\n...",
    "labels": ["enhancement", "ci/cd"]
  }
}
```

### Opção 2: Via GitHub CLI (`gh issue create`)
Execute o comando via terminal (sandboxed):

```bash
gh issue create \
  --repo yurinogueira/CVMC \
  --title "feat(cars): adicionar exportação do histórico de manutenção em PDF" \
  --body-file /caminho/para/issue_body.md \
  --label "enhancement,feature"
```

### Opção 3: Apresentação Estruturada para Validação do Usuário
Apresente a User Story / Issue completa com todos os critérios de aceite e detalhes arquiteturais para revisão do usuário antes ou durante a criação.
