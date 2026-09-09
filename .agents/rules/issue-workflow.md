---
description: Obrigatoriedade da ativação do cvmc-workflow e sincronização da main remota na resolução de issues
globs: "**/*"
---

# 📋 Regra de Resolução de Issues — CVMC

Esta diretriz é **mandatória** e deve ser seguida por qualquer agente ou desenvolvedor ao resolver, implementar ou corrigir issues (bugs, features ou melhorias) no projeto **CVMC (Como Vai Meu Carro)**.

---

## 🎯 1. Ativação Obrigatória do `cvmc-workflow`

Sempre que a tarefa envolver uma issue do GitHub ou resolução de funcionalidade/defeito:
1. **Ative e consulte a skill `cvmc-workflow`** como referência de procedimento.
2. O ciclo de vida completo deve ser respeitado:
   - Mapeamento e leitura da issue via GitHub MCP (`get_issue`).
   - Sincronização com `origin/main`.
   - Desenvolvimento alinhado aos padrões arquiteturais (`cvmc-dev`) e de segurança (`cvmc-security`).
   - Commits semânticos no padrão Conventional Commits referenciando o ID da issue (`#<id>`).
   - Validação com `./scripts/swagger.sh` (se aplicável) e `./scripts/check.sh all`.
   - Abertura de Pull Request direcionado para a branch `main` via GitHub MCP (`create_pull_request`).
   - Fechamento e comentário na issue via GitHub MCP (`add_issue_comment` e `update_issue`).

---

## 🔄 2. Sincronização da `main` antes de Worktrees ou Branches

- **Nunca** inicie alterações ou crie worktrees a partir de uma `main` local sem antes sincronizá-la com `origin/main`.
- Caso esteja executando dentro de uma worktree pré-criada, verifique imediatamente o commit base em relação a `origin/main`:
  ```bash
  git fetch origin main
  # Se a branch local estiver atrás de origin/main e sem commits de trabalho:
  git reset --hard origin/main
  ```
- Isso previne conflitos de merge, regressões ou código baseado em versões ultrapassadas da branch principal.
