## 📌 Descrição das Alterações

<!-- Forneça um resumo claro e conciso das alterações realizadas e o contexto técnico ou de negócio. -->

---

## 🏷️ Tipo de Alteração

Marque as opções que se aplicam a este Pull Request:

- [ ] ✨ `feat`: Nova funcionalidade ou recurso para o usuário
- [ ] 🐛 `fix`: Correção de bug ou comportamento inesperado
- [ ] ♻️ `refactor`: Refatoração de código sem alteração de comportamento externo
- [ ] ⚡ `perf`: Melhoria de performance ou carregamento (LCP/FCP)
- [ ] 📝 `docs`: Atualização ou adição de documentação
- [ ] 🔧 `chore`: Manutenção de dependências, builds ou arquivos auxiliares
- [ ] 👷 `ci`: Alterações em pipelines de CI/CD ou automações do GitHub Actions
- [ ] 🔒 `security`: Correções ou melhorias focadas em segurança defensiva

---

## 🔗 Issues Relacionadas

<!-- Vincule a issue correspondente para fechamento automático após o merge. Exemplos: Closes #123, Resolves #456 -->
- Closes #
- Resolves #

---

## 🧪 Checklist de Validação Técnica

Antes de submeter o PR, confirme se os seguintes passos foram executados com sucesso no seu ambiente local:

- [ ] **Sincronia com a Main**: A branch foi atualizada com o commit mais recente de `origin/main` (`git fetch origin main && git rebase origin/main`).
- [ ] **Regeneração da Documentação Swagger** (obrigatório caso tenha alterado rotas/handlers HTTP da API Go):
  ```bash
  ./scripts/swagger.sh
  ```
- [ ] **Validação Completa da Stack** (Backend Go vet + testes, Frontend typecheck + lint + format + vitest e Terraform fmt):
  ```bash
  ./scripts/check.sh all
  ```

---

## 🖼️ Demonstração Visual / Evidências (Obrigatório para Frontend)

<!-- Se este PR incluir alterações visuais, telas novas ou componentes do Design System, inclua capturas de tela ou GIFs antes/depois. -->

| Antes | Depois |
| :---: | :---: |
| _Insira imagem ou N/A_ | _Insira imagem ou N/A_ |

---

## 🛡️ Checklist de Segurança & Boas Práticas

- [ ] Tokens JWT trafegados estritamente através de cookies `HttpOnly`, `Secure`, `SameSite=Lax`.
- [ ] Nenhum token sensível armazenado em `localStorage` ou `sessionStorage`.
- [ ] Nenhuma credencial, segredo ou IP privado/público exposto em arquivos de código, configurações ou workflows.
- [ ] Entradas de usuário devidamente sanitizadas (sem injeção NoSQL, XSS ou Path Traversal).
