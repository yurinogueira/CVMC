---
type: domain
title: "Subdomínio: Consulta e Cache da Tabela FIPE (Fipe)"
description: "Integração externa com a Tabela FIPE para obtenção de marcas, modelos, anos-modelo e histórico de preços."
tags:
  - domain
  - fipe
  - pricing
  - market-value
resource: backend/internal/domain/fipe
timestamp: 2026-10-10
---

# 📊 Subdomínio: Tabela FIPE (Fipe)

O subdomínio **Fipe** provê integração com o índice oficial de preços de veículos no Brasil, permitindo avaliar a depreciação e o valor médio de mercado dos carros cadastrados.

Para navegação geral, retorne ao [Catálogo Canônico](../index.md).

---

## 🧩 Modelo de Domínio (`FipePrice`)

Estruturado em `backend/internal/domain/fipe/fipe.go`:

```go
type Brand struct {
    Code string `json:"code"`
    Name string `json:"name"`
}

type Model struct {
    Code string `json:"code"`
    Name string `json:"name"`
}

type YearModel struct {
    Code string `json:"code"`
    Name string `json:"name"`
}

type FipePrice struct {
    FipeCode      string    `json:"fipeCode"`
    Brand         string    `json:"brand"`
    Model         string    `json:"model"`
    ModelYear     int       `json:"modelYear"`
    Fuel          string    `json:"fuel"`
    Price         float64   `json:"price"`
    ReferenceMonth string   `json:"referenceMonth"`
    UpdatedAt     time.Time `json:"updatedAt"`
}
```

---

## 📋 Resiliência e Estratégia de Cache

1. **Camada de Cache no Banco**:
   - Consultas de preços FIPE são armazenadas no MongoDB para evitar gargalos na API externa e garantir disponibilidade offline.
2. **Atualização Periódica**:
   - As cotações da Tabela FIPE mudam mensalmente. O cache possui TTL ou renovação controlada por mês de referência.

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Subdomínio de Veículos](car.md)
- [Isolamento de Dados e Persistência](../architecture/multitenancy-and-data.md)
