## 📌 Tipo de Alteração

Marque com um `x` a opção que melhor descreve sua mudança:

- [ ] ✨ `feat`: Nova funcionalidade ou melhoria de produto
- [ ] 🐛 `fix`: Correção de bug
- [ ] 📝 `docs`: Atualização ou adição de documentação
- [ ] ♻️ `refactor`: Refatoração de código sem alteração de comportamento externo
- [ ] ⚡ `perf`: Otimização de performance
- [ ] 🧪 `test`: Adição ou correção de testes automatizados
- [ ] 🔧 `chore`: Tarefas de manutenção, dependências ou infraestrutura
- [ ] 🛡️ `security`: Correção ou aprimoramento de segurança
- [ ] 👷 `ci`: Alterações nos fluxos de CI/CD

---

## 🎯 Contexto e Motivação

Descreva o motivo desta alteração e o contexto do problema/oportunidade. Se aplicável, vincule a issue correspondente para fechamento automático:

- **Issue vinculada**: Closes # <!-- ou Resolves # -->

---

## 🛠️ Resumo das Alterações Técnicas

Liste as principais mudanças técnicas realizadas neste PR:

- 
- 
- 

---

## 🔒 Considerações de Segurança e Integridade

- [ ] **Autenticação & Autorização**: Validação de escopo de usuário, permissões RBAC e integridade de cookies `HttpOnly`.
- [ ] **Sanitização & Validação**: Validação de inputs no backend, proteção contra XSS, NoSQL Injection e *Path Traversal*.
- [ ] **Exposição de Dados**: Nenhuma informação sensível, credencial ou segredo foi exposto.

---

## 🧪 Checklist de Validação Obrigatória

Antes de solicitar a revisão do PR, certifique-se de que todos os itens abaixo foram cumpridos:

- [ ] A branch foi criada e rebaseada a partir da `main` mais recente (`git pull origin main --ff-only` / `git rebase origin/main`).
- [ ] O comando de validação completa passou com 100% de sucesso (Backend, Frontend, Terraform e Docs):
  ```bash
  ./scripts/check.sh all
  ```
- [ ] A documentação Swagger foi atualizada caso novos endpoints ou DTOs tenham sido alterados:
  ```bash
  ./scripts/swagger.sh
  ```
- [ ] A documentação canônica em `docs/` e o log diário em `docs/logs/` (indexado em `docs/log.md`) foram atualizados.
- [ ] Os commits seguem o padrão [Conventional Commits](https://www.conventionalcommits.org/).
