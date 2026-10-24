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
4. **Monitoramento de Revisões Periódicas (`periodic.go`)**:
   - O sistema monitora revisões preventivas por quilometragem ou tempo a partir da última troca registrada de cada item ou da quilometragem inicial de cadastro do veículo (`InitialMileage` / marco zero).
   - **Intervalos Padrão de Fábrica (Fallback)**:
     - **Carros e Caminhonetes (`cars`, `trucks`)**:
       - `Óleo de Motor` (Troca de óleo): 10.000 km ou 12 meses.
       - `Filtro do Óleo de Motor` (Troca de filtro de óleo): 10.000 km ou 12 meses.
       - `Filtro de Ar do Motor` (Troca de filtro de ar): 10.000 km ou 12 meses.
       - `Filtro do Ar-Condicionado` (Troca do filtro de ar-condicionado): 20.000 km ou 12 meses.
       - `Filtro de Combustível` (Troca do filtro de combustível): 20.000 km ou 12 meses.
     - **Motos (`motorcycles`)**:
       - `Óleo de Motor` (Troca de óleo): 2.500 km ou 12 meses.
       - `Filtro do Óleo de Motor` (Troca de filtro de óleo): 2.500 km ou 12 meses.
       - `Filtro de Ar do Motor` (Troca de filtro de ar): 5.000 km ou 12 meses.
   - **Intervalos Personalizáveis por Veículo (`car.MaintenanceIntervals`)**:
     - Cada proprietário pode configurar individualmente os intervalos desejados de cada item preventivo diretamente pelo painel do veículo ou pela tela de manutenções.
     - As opções são estritamente pré-definidas para garantir integridade:
       - Carros: 5.000, 7.500, 10.000, 15.000, 20.000, 30.000 km; 3, 6, 12, 18, 24 meses.
       - Motos: 1.000, 2.500, 5.000, 10.000 km; 3, 6, 12, 18, 24 meses.
     - A arquitetura orientada a dicionário (`map[string]MaintenanceIntervalConfig`) é aberta e preparada para suportar a adição de novos itens preventivos no futuro sem alterações estruturais no banco de dados.
     - Um botão de restauração permite reverter todos os itens para os padrões recomendados a qualquer momento.
   - **Classificação de Status**:
     - `ok` (`Em dia`): mais de 1.000 km e mais de 30 dias restantes.
     - `upcoming` (`Próximo`): faltam ≤ 1.000 km ou ≤ 30 dias.
     - `overdue` (`Vencido`): quilometragem ou data limite atingida/excedida (`<= 0`).
   - **Endpoints de Consulta**:
     - `GET /api/v1/cars/{id}/upcoming-maintenances`: lista detalhada das próximas revisões preventivas do veículo.
     - `GET /api/v1/maintenances/upcoming`: resumo consolidado de toda a frota do usuário autenticado.
5. **Cronjob de Alertas e Consolidação Anti-Spam**:
   - O scheduler em background (`internal/infrastructure/scheduler/maintenance_cron.go`, configurável via `MAINTENANCE_CRON_INTERVAL`) varre periodicamente os veículos de usuários com tokens FCM ativos.
   - Quando um ou mais itens atingem o limiar de alerta (faltando ≤ 500 km ou ≤ 15 dias, ou vencidos), o serviço de notificações consolida **todos os itens pendentes do veículo em uma única notificação push** (`Revisão Preventiva: {Nome do Veículo}`) apontando para `/maintenance?carId={carId}`, respeitando um *cooldown* de 7 dias (exceto se um novo item se tornar crítico).

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Subdomínio de Veículos](car.md)
- [Subdomínio de Anexos](attachment.md)
- [Armazenamento e Gestão de Mídia](../architecture/storage-and-media.md)
