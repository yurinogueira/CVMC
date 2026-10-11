---
description: Diretrizes mandatórias de integridade e conformidade da documentação canônica (OKF) no CVMC
globs: "**/*"
---

# 📚 Regras de Documentação Canônica (OKF) — CVMC

Estas diretrizes são **mandatórias** e devem ser observadas por desenvolvedores e agentes autônomos de IA ao realizar qualquer modificação no projeto **CVMC (Como Vai Meu Carro)**.

---

## 🛑 1. Fonte Única da Verdade

- A base de documentação em `docs/` é estruturada segundo a especificação **Open Knowledge Format (OKF)** do Google Cloud e constitui a única fonte canônica da verdade sobre o sistema.
- Antes de iniciar o desenvolvimento ou propor refatorações de código, consulte o catálogo canônico em `docs/index.md` e a especificação de domínio correspondente em `docs/domain/` para evitar alucinações e premissas equivocadas.

---

## 🔄 2. Sincronização Universal e Mandatória ("Zero-Divergence & Continuous Docs")

> [!CAUTION]
> **Regra Universal Inegociável**: **TODA E QUALQUER** modificação realizada no repositório — sem exceção alguma — exige atualização da documentação canônica e registro histórico datado no mesmo commit ou Pull Request.
> Isso se aplica estritamente a:
> 1. **Novas Funcionalidades (`feat`)**: novas rotas, componentes de UI, fluxos de navegação, contratos de API ou entidades DDD.
> 2. **Correção de Defeitos (`fix`)**: correções de regras de negócio, ajustes de validação, correções de segurança ou comportamentos corrigidos.
> 3. **Atualização de Versões e Dependências (`chore`/`deps`/`ci`)**: bumps de runtime (ex: Go, Node), bibliotecas (Vite, React, MUI, drivers), imagens Docker (ex: MongoDB) ou providers Terraform.
> 4. **Refatorações e Arquitetura (`refactor`/`perf`)**: reestruturações de pastas, melhorias de desempenho, alteração de contratos ou novos ADRs.

### Deveres Inegociáveis do Agente em Toda Tarefa:
1. **Atualizar os Documentos Canônicos Correspondentes em `docs/`**:
   - Alterou entidades ou use cases? Atualize `docs/domain/`.
   - Alterou autenticação, permissões ou segurança? Atualize `docs/architecture/auth-and-security.md`.
   - Alterou telas, componentes, layouts ou rotas? Atualize `docs/frontend/`.
   - Alterou Docker, dependências, pipelines ou env vars? Atualize `docs/operations/`.
2. **Proibido Documento Órfão**: Todo novo arquivo `.md` criado sob `docs/` deve conter cabeçalho YAML frontmatter com campo `type:` válido e ser indexado em `docs/index.md`.
3. **Proibido Link Quebrado**: Todas as referências cruzadas devem utilizar links Markdown relativos válidos.
4. **Registro Histórico Particionado por Dia (`docs/logs/AAAA-MM-DD.md` e `docs/log.md`)**:
   - Cada entrega **DEVE** ser registrada no arquivo diário correspondente à data atual (`docs/logs/AAAA-MM-DD.md`):
     - **Se o arquivo do dia já existir**: Adicione uma nova seção de nível 2:
       `## 🎯 <Título Semântico> (#<id_issue_ou_pr>)`
       com o resumo, modificações por camada e documentos canônicos atualizados.
     - **Se for o primeiro registro do dia**: Crie o arquivo `docs/logs/AAAA-MM-DD.md` com frontmatter YAML (`type: log`, `timestamp: AAAA-MM-DD`, etc.), adicione a seção da intervenção e registre o novo dia no catálogo central [docs/log.md](log.md).
   - O catálogo central [docs/log.md](log.md) atua como índice enxuto, mantendo a listagem das datas com links para os arquivos diários e sumário em bullet points, economizando tokens e evitando arquivos gigantes.
   - O campo `timestamp:` no frontmatter de `docs/log.md` e do arquivo diário deve refletir a data da modificação (ISO `AAAA-MM-DD`).

---

## 🧪 3. Validação Mandatória Pré-Commit

- Antes de abrir Pull Requests ou concluir tarefas, execute a validação automatizada:
  ```bash
  ./scripts/check.sh docs
  ```
- O pipeline `.github/workflows/docs.yml` e o checklist de PR bloquearão entregas que não tenham documentação ou violem a conformidade OKF.
