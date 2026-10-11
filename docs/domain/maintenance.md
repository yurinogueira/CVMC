---
type: domain
title: "Subdomínio: Manutenções Preventivas e Corretivas (Maintenance)"
description: "Ordens de serviço, controle de custos mecânicos, oficinas, agendamento de revisões e notas fiscais."
tags:
  - domain
  - maintenance
  - workshop
  - costs
resource: backend/internal/domain/maintenance
timestamp: 2026-10-10
---

# 🔧 Subdomínio: Manutenções (Maintenance)

O subdomínio **Maintenance** rastreia o histórico mecânico, elétrico e preventivo de cada veículo, fornecendo métricas de custo total de manutenção.

Para navegação geral, retorne ao [Catálogo Canônico](../index.md).

---

## 🧩 Modelo de Domínio (`Maintenance`)

Estruturado em `backend/internal/domain/maintenance/maintenance.go`:

```go
type MaintenanceType string

const (
    MaintenancePreventive MaintenanceType = "preventive"
    MaintenanceCorrective MaintenanceType = "corrective"
)

type Maintenance struct {
    ID              string          `json:"id"`
    CarID           string          `json:"carId"`
    UserID          string          `json:"userId"`
    Type            MaintenanceType `json:"type"`
    Description     string          `json:"description"`
    Workshop        string          `json:"workshop,omitempty"`
    KmAtMaintenance int             `json:"kmAtMaintenance"`
    Cost            float64         `json:"cost"`
    Date            time.Time       `json:"date"`
    Attachments     []string        `json:"attachments,omitempty"`
    CreatedAt       time.Time       `json:"createdAt"`
}
```

---

## 📋 Regras de Negócio

1. **Vínculo Obrigatório ao Veículo**:
   - Toda manutenção deve estar vinculada a um `CarID` existente de propriedade do usuário autenticado.
2. **Atualização da Quilometragem do Carro**:
   - Ao registrar uma manutenção com `KmAtMaintenance > car.CurrentKm`, o sistema atualiza automaticamente a quilometragem atual do veículo.
3. **Anexos e Notas Fiscais**:
   - Cada manutenção pode conter múltiplos anexos (comprovantes e ordens de serviço) devidamente higienizados e validados quanto ao tipo MIME.

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Subdomínio de Veículos](car.md)
- [Subdomínio de Anexos](attachment.md)
- [Armazenamento e Gestão de Mídia](../architecture/storage-and-media.md)
