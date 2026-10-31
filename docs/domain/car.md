---
type: domain
title: "Subdomínio: Veículos & Gestão de Frota (Car)"
description: "Modelagem de dados do veículo, placa, quilometragem atual, histórico, validação e regras de negócio."
tags:
  - domain
  - car
  - vehicle
  - fleet
resource: backend/internal/domain/car
timestamp: 2026-10-10
---

# 🚗 Subdomínio: Veículos & Gestão de Frota (Car)

O subdomínio **Car** representa a entidade central da plataforma **CVMC**, encapsulando as informações cadastrais e métricas operacionais de cada automóvel.

Para navegação geral, retorne ao [Catálogo Canônico](../index.md).

---

## 🧩 Modelo de Domínio (`Car`)

Estruturado em `backend/internal/domain/car/car.go`:

```go
type MaintenanceIntervalConfig struct {
    IntervalKm     int `json:"intervalKm" bson:"intervalKm"`
    IntervalMonths int `json:"intervalMonths" bson:"intervalMonths"`
}

type Car struct {
    ID                   string                               `json:"id"`
    OwnerID              string                               `json:"ownerId"`
    Name                 string                               `json:"name"`
    Manufacturer         string                               `json:"manufacturer"`
    Model                string                               `json:"model"`
    YearManufacture      int                                  `json:"yearManufacture"`
    YearModel            int                                  `json:"yearModel"`
    InitialMileage       int                                  `json:"initialMileage,omitempty"`
    LastMileage          int                                  `json:"lastMileage"`
    VehicleType          string                               `json:"vehicleType,omitempty"`
    ImageUrl             string                               `json:"imageUrl,omitempty"`
    FIPECode             string                               `json:"fipeCode,omitempty"`
    FIPEPrice            string                               `json:"fipePrice,omitempty"`
    Fuel                 string                               `json:"fuel,omitempty"`
    SharedWith           []string                             `json:"sharedWith,omitempty"`
    MaintenanceIntervals map[string]MaintenanceIntervalConfig `json:"maintenanceIntervals,omitempty"`
    MileageUpdatedAt     time.Time                            `json:"mileageUpdatedAt,omitempty"`
    CreatedAt            time.Time                            `json:"createdAt"`
    UpdatedAt            time.Time                            `json:"updatedAt"`
}
```

---

## 📋 Regras de Negócio e Validações

1. **Atualização Monotônica da Quilometragem**:
   - Abastecimentos e manutenções podem atualizar o `LastMileage` do veículo caso a quilometragem informada seja superior à registrada atualmente.
   - Não é permitido registrar quilometragem inferior à última registrada sem justificativa.
2. **Associação FIPE**:
   - Cada veículo pode ser vinculado opcionalmente a um código FIPE (`FipeCode`) para cálculo automático do valor de referência de mercado.
3. **Intervalos Personalizados de Revisão Preventiva**:
   - Cada veículo armazena em `MaintenanceIntervals` suas metas específicas de quilometragem e meses para cada item preventivo (`Óleo de Motor`, filtros, etc.).
   - Quando não configurado, o sistema adota os valores recomendados de fábrica para a categoria (`cars`, `trucks`, `motorcycles`).
   - A estrutura é aberta e suporta novos componentes no futuro sem migrações de schema.
4. **Rastreamento de Atualização de Quilometragem e Lembrete Semanal (`MileageUpdatedAt`)**:
   - O campo `MileageUpdatedAt` registra o timestamp UTC exato da última alteração de `LastMileage` (inicializado em `CreateCar`, atualizado em `UpdateCar` quando a quilometragem muda, e avançado em `maintenance.Service` quando uma manutenção supera a quilometragem atual).
   - Para veículos legados sem `mileageUpdatedAt`, o método `EffectiveMileageUpdatedAt()` realiza fallback automático para `UpdatedAt` e `CreatedAt`.
   - O cronjob `MileageReminderCron` (`MILEAGE_REMINDER_CRON_INTERVAL`, padrão `168h` / 7 dias) avalia os veículos de usuários com tokens FCM ativos e envia uma notificação push (`Lembrete de Quilometragem 🚗`, `type: mileage_reminder`, rota `/cars/{id}`) caso o veículo esteja há mais de 3 dias sem atualização de quilometragem, respeitando cooldown de 7 dias por veículo/usuário.

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Isolamento de Dados e Persistência](../architecture/multitenancy-and-data.md)
- [Subdomínio de Manutenções](maintenance.md)
- [Subdomínio de Abastecimentos](fuel.md)
- [Subdomínio da Tabela FIPE](fipe.md)
