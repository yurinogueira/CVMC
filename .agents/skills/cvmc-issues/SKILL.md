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

1. **Investigação Prévia Obrigatória**: Antes de redigir a issue, explore a base de código para mapear os arquivos, rotas, contratos e componentes envolvidos.
2. **Precisão de Localização**: Toda issue deve conter caminhos de arquivo exatos (`backend/...`, `frontend/...`), nomes de structs, funções, endpoints REST (`/api/v1/...`) e stores Zustand.
3. **Título Semântico**: Use a convenção do Conventional Commits para títulos de issue:
   - Bugs: `fix(<modulo>): <descrição sucinta em minúsculas>` (ex: `fix(auth): cookie de sessão não enviado na rota de perfil`)
   - Features: `feat(<modulo>): <descrição sucinta em minúsculas>` (ex: `feat(maintenance): alerta de vencimento de revisão por quilometragem`)
4. **Alinhamento Arquitetural**: Respeite os padrões documentados em `cvmc-dev` (Clean Arch + DDD no Go 1.25 e React 19 + MUI v6 no Frontend) e `cvmc-security` (Cookies HttpOnly, validação rigorosa de ownership, sanitização).

---

## 🐛 Fluxo para Criação de Issues de Bug

Ao relatar um defeito, siga rigorosamente a seguinte estrutura:

### 1. Investigação do Bug
- Identifique a causa raiz ou o ponto de falha navegando no repositório.
- Colete payloads de requisição/resposta, logs do backend (`./scripts/dev.sh logs backend`) ou erros no console do frontend.
- Identifique se o problema ocorre no backend, frontend, banco MongoDB, infraestrutura Docker ou documentação Swagger.

### 2. Estrutura Padrão do Relato de Bug

```markdown
### Descrição do Problema
[Descrição clara e objetiva do erro observado]

### Componente / Camada Afetada
[Backend | Frontend | Infraestrutura / Docker | Banco MongoDB | Autenticação | Swagger]

### Onde o Problema Ocorre (Localização Técnica)
- **Rota / Endpoint**: `METODO /api/v1/...`
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

## ✨ Fluxo para Criação de Issues de Feature

Ao propor uma nova funcionalidade ou melhoria técnica, detalhe a especificação funcional e arquitetural:

### 1. Levantamento Arquitetural
- Mapeie as camadas que serão impactadas no backend:
  - Entidades e regras puras (`backend/internal/domain/`)
  - Interfaces de repositório e serviços (`backend/internal/application/ports/`)
  - Use cases e testes unitários (`backend/internal/application/usecase/`)
  - Adapters de infraestrutura (`backend/internal/infrastructure/`)
  - Handlers REST, rotas e Swagger (`backend/internal/interfaces/rest/`)
- Mapeie o fluxo no frontend:
  - Telas e componentes visuais (`frontend/src/features/<feature>/`)
  - Estado global e stores Zustand (`frontend/src/features/<feature>/state/`)
  - Integração com cliente Axios (`frontend/src/services/api/client.ts`)

### 2. Estrutura Padrão da Proposta de Feature

```markdown
### Visão Geral e Contexto de Negócio
**Como** [tipo de usuário/ator],
**Quero** [ação ou capacidade desejada],
**Para que** [benefício ou valor entregue].

### Requisitos Funcionais e Critérios de Aceite
- [ ] [Critério 1: Regra de negócio específica]
- [ ] [Critério 2: Validação de dados de entrada]
- [ ] [Critério 3: Resposta da API ou comportamento da interface]

### Camadas Técnicas Impactadas
- [x] Backend - Domínio e Contratos (`internal/domain/`, `internal/application/ports/`)
- [x] Backend - Casos de Uso e Testes (`internal/application/usecase/`)
- [x] Backend - Interfaces REST e Swagger (`internal/interfaces/rest/`, `./scripts/swagger.sh`)
- [ ] Frontend - Telas e Componentes MUI (`frontend/src/features/`)

### Detalhes Técnicos e Arquitetura Proposta
- **Endpoints REST**: `POST /api/v1/recurso`
- **DTOs / Estruturas**: Structs em Go com tags `json:"camelCase"`
- **Componentes UI**: Componentes do Design System MUI e ícones de `@mui/icons-material`

### Considerações de Segurança e Performance
- Autenticação e autorização via cookies `HttpOnly` com validação de ownership (`userID`).
- Defesa em profundidade (sanitização de inputs, rate limiting, CORS restrito).
- Execução de testes compactos sem poluição de tokens (`./scripts/check.sh all`).

### Checklist de Implementação
- [ ] Criar/atualizar entidades e interfaces no backend
- [ ] Implementar caso de uso e testes unitários (`*_test.go`)
- [ ] Adicionar handlers REST e anotações Swagger
- [ ] Executar `./scripts/swagger.sh` e `./scripts/check.sh backend`
- [ ] Implementar componentes e integração no frontend
- [ ] Validar com `./scripts/check.sh frontend`
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
    "title": "fix(auth): validar expiracao do cookie de sessao no middleware",
    "body": "### Descrição do Problema\n...",
    "labels": ["bug"]
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
Se o usuário solicitar apenas o rascunho ou a estruturação antes do envio, gere o Markdown completo dentro do chat ou artefato estruturado para revisão.
