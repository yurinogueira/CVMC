---
type: adr
title: "ADR 0005: Política de Versões do MongoDB e Paridade com MongoDB Atlas M0"
description: "Decisão de fixar a versão do MongoDB na linha 8 LTS no Docker Compose, garantindo paridade operacional com o cluster gerenciado MongoDB Atlas M0 e preservando volumes de desenvolvimento."
tags:
  - adr
  - architecture
  - database
  - mongodb
  - docker
  - parity
timestamp: 2026-10-30
---

# 📜 ADR 0005: Política de Versões do MongoDB e Paridade com MongoDB Atlas M0

- **Status**: Aceito
- **Data**: 2026-10-30
- **Decisores**: Equipe de Engenharia CVMC
- **Issues / PRs Relacionados**: Issue [#148](https://github.com/yurinogueira/CVMC/issues/148), PR [#132](https://github.com/yurinogueira/CVMC/pull/132)

Para navegação geral, retorne ao [Catálogo Canônico](../index.md).

---

## 📌 Contexto

O Dependabot propôs, através do PR [#132](https://github.com/yurinogueira/CVMC/pull/132), a atualização da imagem do banco de dados no `docker-compose.yml` de `mongo:8` para `mongo:9`.

Em conformidade com as diretrizes de governança da plataforma CVMC e a análise formal detalhada na Issue [#148](https://github.com/yurinogueira/CVMC/issues/148), upgrades de *major version* em bancos de dados relacionais ou documentais demandam auditoria de três pilares fundamentais:

1. **Paridade Operacional Dev vs. Prod (Zero-Divergence)**:
   - Em produção, o CVMC provisiona um cluster gerenciado MongoDB Atlas no tier gratuito compartilhado (`M0`), conforme declarado no código de infraestrutura Terraform (`terraform/mongodb.tf`):
     ```hcl
     resource "mongodbatlas_cluster" "cluster" {
       project_id                  = mongodbatlas_project.project.id
       name                        = var.mongodb_cluster_name
       provider_name               = "TENANT"
       backing_provider_name       = "AWS"
       provider_region_name        = "SA_EAST_1"
       provider_instance_size_name = "M0"
     }
     ```
   - O MongoDB Atlas M0 é mantido e atualizado automaticamente pelo provedor na linha **MongoDB 8.0 LTS**, não suportando versões da linha 9.
   - Adotar o MongoDB 9 localmente criaria assimetria de ambientes, onde operadores de agregação, tipos de índices ou comportamentos sintáticos presentes no desenvolvimento local poderiam falhar silenciosamente ao atingir a esteira de homologação ou produção.

2. **Integridade de Dados e Compatibilidade de Volumes Locais (`WiredTiger` e `fCV`)**:
   - Baterias de testes empíricos com containers demonstraram que a substituição direta da imagem `mongo:8` por `mongo:9` sobre volumes de dados pré-existentes (`mongo_data`) provoca encerramento fatal imediato do serviço com código de saída 62 (`exitCode: 62`):
     ```text
     Wrong mongod version: Found an invalid featureCompatibilityVersion document
     (ERROR: Location4926900: Invalid featureCompatibilityVersion document in admin.system.version:
     { _id: "featureCompatibilityVersion", version: "8.2" }.
     Expected one of the following versions: '9.0', '8.3', '8.0'.)
     ```
   - A imagem `mongo:8` opera com MongoDB 8.2 e fCV 8.2. A migração transparente in-place não é permitida pelo motor da versão 9.0 sem intervenção prévia manual na versão de compatibilidade, quebrando a experiência de novos e atuais desenvolvedores da plataforma.

3. **Compatibilidade do Driver Go**:
   - O backend utiliza o driver oficial `go.mongodb.org/mongo-driver/v2 v2.9.2` (`backend/go.mod`), cujo suporte e estabilidade de transações, BSON e índices estão plenamente testados e homologados para a linha 8 LTS.

---

## 💡 Decisão

1. **Fixação no MongoDB 8 LTS**: Manter a imagem `mongo:8` no `docker-compose.yml`, assegurando estabilidade, paridade estrita com o MongoDB Atlas M0 e preservação dos dados locais de desenvolvimento.
2. **Governança no Dependabot**: Configurar regra explícita de `ignore` para *major version updates* da dependência `mongo` em `.github/dependabot.yml`:
   ```yaml
   ignore:
     - dependency-name: "mongo"
       update-types: ["version-update:semver-major"]
   ```
3. **Desfecho do PR #132**: Fechar o PR [#132](https://github.com/yurinogueira/CVMC/pull/132) sem merge, documentando a justificativa técnica vinculada a este ADR e à Issue [#148](https://github.com/yurinogueira/CVMC/issues/148).
4. **Gatilho de Reavaliação Futura**: Qualquer migração futura para o MongoDB 9 ou superior fica condicionada à disponibilização prévia e homologação oficial da versão no cluster gerenciado MongoDB Atlas M0.

---

## ⚖️ Consequências

### Positivas:
- **Zero-Divergence Garantido**: Paridade fidedigna entre o ambiente de testes/desenvolvimento local e o ambiente de produção na nuvem.
- **Resiliência e Continuidade Local**: Prevenção contra quebras catastróficas de volumes locais (`exitCode: 62`) e necessidade de rotinas manuais de dump/restore pelos desenvolvedores.
- **Automação Limpa**: Dependabot não gerará ruído nem reabrirá PRs automáticos de migração de versão maior enquanto o bloqueio estiver ativo.

### Mitigações / Trade-offs:
- Recursos específicos introduzidos no MongoDB 9 (tais como melhorias avançadas em Queryable Encryption e Intelligent Workload Management) não estarão disponíveis localmente, o que não impacta a operação do CVMC neste momento.

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Isolamento de Dados e Persistência](../architecture/multitenancy-and-data.md)
- [Docker e Ambiente de Desenvolvimento Local](../operations/docker-and-local-dev.md)
- [Runbook: Backup e Restauração de Banco de Dados](../operations/runbooks/database-backup.md)
- [Infraestrutura e Deploy](../operations/deploy-and-infrastructure.md)
