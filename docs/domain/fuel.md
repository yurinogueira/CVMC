---
type: domain
title: "Subdomínio: Abastecimentos & Consumo de Combustível (Fuel)"
description: "Registro de abastecimentos, tipos de combustível, cálculo de consumo médio em km/l e gastos financeiros."
tags:
  - domain
  - fuel
  - efficiency
  - consumption
resource: backend/internal/domain/fuel
timestamp: 2026-10-10
---

# ⛽ Subdomínio: Abastecimentos (Fuel)

O subdomínio **Fuel** gerencia os abastecimentos realizados nos veículos, oferecendo cálculo estatístico de consumo médio e custo por quilômetro rodado.

Para navegação geral, retorne ao [Catálogo Canônico](../index.md).

---

## 🧩 Modelo de Domínio (`Fueling`)

Estruturado em `backend/internal/domain/fuel/fuel.go`:

```go
type FuelType string

const (
    FuelGasoline FuelType = "gasoline"
    FuelEthanol  FuelType = "ethanol"
    FuelDiesel   FuelType = "diesel"
    FuelCNG      FuelType = "cng"
    FuelElectric FuelType = "electric"
)

type Fueling struct {
    ID          string    `json:"id"`
    CarID       string    `json:"carId"`
    UserID      string    `json:"userId"`
    FuelType    FuelType  `json:"fuelType"`
    Liters      float64   `json:"liters"`
    PricePerLtr float64   `json:"pricePerLiter"`
    TotalCost   float64   `json:"totalCost"`
    CurrentKm   int       `json:"currentKm"`
    GasStation  string    `json:"gasStation,omitempty"`
    Date        time.Time `json:"date"`
    CreatedAt   time.Time `json:"createdAt"`
}
```

---

## 📋 Regras de Negócio e Cálculos

1. **Cálculo de Autonomia e Consumo Médio**:
   - Entre dois abastecimentos consecutivos com tanque cheio (*full tank*), o consumo médio é calculado pela fórmula:
     $$\text{Consumo (km/l)} = \frac{\text{Km Atual} - \text{Km Anterior}}{\text{Litros Abastecidos}}$$
2. **Consistência de Valores Financeiros**:
   - O `TotalCost` deve corresponder ao produto de `Liters * PricePerLtr` (com tolerância de arredondamento de centavos).

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Subdomínio de Veículos](car.md)
- [Isolamento de Dados e Persistência](../architecture/multitenancy-and-data.md)
