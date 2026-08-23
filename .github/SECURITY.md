# Política de Segurança (Security Policy) — CVMC

A equipe do **CVMC (Como Vai Meu Carro)** leva a segurança de seus sistemas, dados e usuários muito a sério. Como a plataforma opera no modelo **SaaS com Entrega Contínua (Continuous Delivery)**, correções e atualizações são integradas e disponibilizadas diretamente em produção.

---

## 🌐 Escopo da Plataforma (SaaS & Continuous Delivery)

Por se tratar de um serviço web SaaS, não existem versões legadas ou pacotes offline para download: **o ambiente ativo em produção e a branch `main` representam sempre a versão oficial suportada**.

### ✅ Serviços em Escopo
- **Aplicação Web (Frontend)**: `cvmc.yurinogueira.dev.br`
- **API REST (Backend)**: `api-cvmc.yurinogueira.dev.br`
- **Código-fonte & Configurações**: Branch `main` deste repositório.

### 🎯 Vulnerabilidades Prioritárias em Escopo
- Quebra de autenticação ou vazamento de credenciais/sessões (tokens JWT via cookies `HttpOnly`).
- Falhas de autorização e controle de acesso a dados (BOLA / IDOR em veículos, manutenções ou perfis).
- Injeções (NoSQL Injection, Command Injection, etc.).
- Vulnerabilidades de manipulação de arquivos ou *Path Traversal* em rotas de upload.
- Cross-Site Scripting (XSS) com impacto demonstrável.
- Falhas de integridade em regras de negócio e permissões de usuários.

### ❌ Fora de Escopo
- Ataques volumétricos de Negação de Serviço (DDoS / DoS).
- Ataques de Engenharia Social, Phishing ou tentativas de força bruta física.
- Relatórios puramente teóricos originados de scanners automatizados sem uma Prova de Conceito (PoC) prática e funcional.
- Problemas relacionados a ambientes de desenvolvimento local sem réplica no ambiente de produção.

---

## 🚨 Como Reportar uma Vulnerabilidade (Responsible Disclosure)

> [!IMPORTANT]
> **Por favor, NÃO abra issues públicas no GitHub para relatar problemas de segurança ou potenciais vulnerabilidades.**

Para reportar uma vulnerabilidade de maneira segura e confidencial:

1. **GitHub Private Vulnerability Reporting (Recomendado)**:
   - Acesse a aba **Security** do repositório no GitHub.
   - Clique em **Advisories** > **Report a vulnerability** para abrir um relatório privado e criptografado diretamente com os mantenedores.

2. **Contato Direto**:
   - Caso o recurso de advisories não esteja habilitado, entre em contato de forma privada com o mantenedor responsável pelo repositório.

### O que incluir no relatório:
Para acelerar a triagem e a resolução, solicitamos:
- **Resumo do problema**: descrição clara da fragilidade identificada.
- **Componente / Endpoint afetado**: rota HTTP, página ou módulo da aplicação.
- **Passos para reprodução (PoC)**: comandos cURL, payloads de teste, capturas de tela ou scripts demonstrando o impacto de forma determinística.
- **Impacto estimado**: o que um agente malicioso conseguiria acessar, alterar ou comprometer.
- **Sugestão de mitigação** (opcional).

### Prazos e Processo de Resposta:
- **Confirmação inicial**: Acusamos o recebimento do relatório em até **48 horas**.
- **Triagem e Validação**: Avaliação de impacto em até **5 dias úteis**.
- **Correção e Deploy Contínuo**: A aplicação da correção é prioritária e publicada diretamente em produção no ciclo de CI/CD.
- **Reconhecimento**: Mediante consentimento do pesquisador, podemos atribuir os devidos créditos no GitHub Security Advisory da correção.

---

## 🔒 Postura e Controles de Segurança Implementados

O CVMC adota princípios de Defesa em Profundidade (*Defense in Depth*):

1. **Autenticação & Sessões**:
   - Tokens trafegados exclusivamente através de cookies `HttpOnly`, `Secure` e `SameSite=Lax`.
   - Isolamento total de tokens no Frontend (nunca persistidos em `localStorage` ou `sessionStorage`).
   - Rejeição estrita com `401 Unauthorized` em caso de tokens expirados ou adulterados.

2. **Rede e Proteção HTTP**:
   - Headers de segurança automáticos (`HSTS`, `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`, `Referrer-Policy`).
   - Política de CORS estrita baseada em whitelist (`ALLOWED_ORIGINS`), sem permissão de `*` para requisições autenticadas.
   - Rate limiting defensivo em nível de aplicação para mitigação de abusos.
   - Limite rígido no tamanho do corpo da requisição (`BodyLimit`).

3. **Armazenamento e Sanitização de Arquivos**:
   - Validação canônica absoluta de paths para neutralizar ataques de *Path Traversal* (`../`).

4. **Infraestrutura e Isolamento**:
   - Instâncias de banco de dados (MongoDB) isoladas na rede interna Docker, sem exposição direta de portas na internet pública.
   - Segredos e chaves gerenciados exclusivamente através de variáveis de ambiente seguras e GitHub Actions Secrets.
