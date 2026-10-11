---
type: operations
title: "Ambiente e Configuração — CVMC"
description: "Dicionário de variáveis de ambiente, gestão de segredos (.env.example) e perfis de execução."
tags:
  - operations
  - configuration
  - env
  - secrets
timestamp: 2026-10-10
---

# ⚙️ Ambiente e Configuração — CVMC

O gerenciamento de configurações no **CVMC** adota os princípios do *Twelve-Factor App*, lendo parâmetros a partir de variáveis de ambiente.

Para navegação geral, retorne ao [Catálogo Canônico](../index.md).

---

## 📋 Variáveis Principais do Backend

| Variável | Descrição | Padrão | Obrigatório em Prod |
| :--- | :--- | :--- | :--- |
| `PORT` | Porta de escuta da API REST | `8080` | Sim |
| `JWT_SECRET` | Chave secreta de assinatura JWT (mín. 32 chars) | - | Sim |
| `JWT_REFRESH_SECRET` | Chave de refresh JWT (mín. 32 chars, diferente de JWT_SECRET) | - | Sim |
| `MONGO_URI` | URI de conexão ao MongoDB | `mongodb://localhost:27017` | Sim |
| `MONGO_DATABASE` | Nome do banco no MongoDB | `cvmc` | Sim |
| `STORAGE_PROVIDER` | Provedor de storage (`local` ou `oci`) | `local` | Sim |
| `COOKIE_DOMAIN` | Domínio dos cookies de autenticação | - | Sim (`.cvmc.yurinogueira.dev.br`) |
| `COOKIE_SECURE` | Força flag Secure nos cookies | `true` | Sim |
| `SMTP_HOST` | Host do servidor SMTP para envio de e-mails | - | Sim |
| `ALLOWED_ORIGINS` | Lista de origens permitidas para CORS (web e mobile) | `http://localhost:5173,https://localhost,capacitor://localhost` | Sim (`https://cvmc.yurinogueira.dev.br,https://localhost,capacitor://localhost`) |
| `GOOGLE_CLIENT_ID` | OAuth 2.0 Client ID do Google Cloud para validação backend | - | Sim (em produção com login Google) |
| `VITE_GOOGLE_CLIENT_ID` | OAuth 2.0 Client ID do Google Cloud para GIS frontend | - | Sim (em produção com login Google) |


---

## 🔒 Regras de Segurança para Segredos

- **Proibido Comitar Arquivos `.env`**: Somente `.env.example` pode ser versionado.
- **Validação de Inicialização**: Chaves fracas (`change-me`) causam interrupção imediata na inicialização do backend.

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Docker e Ambiente Local](docker-and-local-dev.md)
- [Deploy e Infraestrutura](deploy-and-infrastructure.md)
