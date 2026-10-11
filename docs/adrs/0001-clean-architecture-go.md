---
type: adr
title: "ADR 0001: Adoção de Clean Architecture e DDD no Backend Go"
description: "Decisão de estruturar o backend Go em camadas concêntricas (Clean Architecture) com princípios de Domain-Driven Design."
tags:
  - adr
  - architecture
  - go
  - clean-architecture
  - ddd
timestamp: 2026-10-10
---

# 📜 ADR 0001: Adoção de Clean Architecture e DDD no Backend Go

- **Status**: Aceito
- **Data**: 2026-10-10
- **Decisores**: Equipe de Engenharia CVMC

Para navegação geral, retorne ao [Catálogo Canônico](../index.md).

---

## 📌 Contexto
O **CVMC (Como Vai Meu Carro)** começou como uma API REST para gestão de automóveis e despesas veiculares. Com o crescimento das regras de negócio (cálculo de consumo médio, validações FIPE e ordens de manutenção), manter regras acopladas a frameworks web e coleções de banco de dados gerava risco de regressões e dificuldades para testes unitários.

---

## 💡 Decisão
Decidimos adotar a **Clean Architecture** combinada com princípios táticos de **Domain-Driven Design (DDD)**:

1. A camada central `internal/domain/` contém entidades puras e regras invariantes sem qualquer importação de frameworks ou drivers externos.
2. A camada `internal/application/ports/` define os contratos de abstração (interfaces de repositórios e serviços).
3. A camada `internal/application/usecase/` orquestra os fluxos de trabalho do sistema consumindo apenas interfaces.
4. Frameworks HTTP (`net/http`, routers) e drivers de banco residem exclusivamente nas camadas periféricas (`interfaces/rest` e `infrastructure`).

---

## ⚖️ Consequências

### Positivas:
- **Testabilidade Superior**: Use cases podem ser testados unitariamente de forma isolada com mocks em milissegundos sem depender do MongoDB.
- **Independência de Provedores**: O provedor de storage ou banco de dados pode ser substituído sem alterar regras de negócio.
- **Clareza de Limites**: Desenvolvedores e agentes de IA têm limites óbvios para cada tipo de modificação.

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Visão Geral da Arquitetura](../architecture/overview.md)
- [Isolamento de Dados e Persistência](../architecture/multitenancy-and-data.md)
