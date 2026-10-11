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
type Car struct {
    ID        string    `json:"id"`
    UserID    string    `json:"userId"`
    Plate     string    `json:"plate"`
    Brand     string    `json:"brand"`
    Model     string    `json:"model"`
    Year      int       `json:"year"`
    Color     string    `json:"color"`
    CurrentKm int       `json:"currentKm"`
    FipeCode  string    `json:"fipeCode,omitempty"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt,omitempty"`
}
```

---

## 📋 Regras de Negócio e Validações

1. **Formato e Unicidade da Placa**:
   - A placa deve ser normalizada em letras maiúsculas e sem pontuação (suporta tanto o padrão antigo brasileiro `ABC1234` quanto o padrão Mercosul `ABC1D23`).
   - Um usuário não pode cadastrar dois veículos com a mesma placa.
2. **Atualização Monotônica da Quilometragem**:
   - Abastecimentos e manutenções podem atualizar o `CurrentKm` do veículo caso a quilometragem informada seja superior à registrada atualmente.
   - Não é permitido registrar quilometragem inferior à última registrada sem justificativa de troca de painel/odômetro.
3. **Associação FIPE**:
   - Cada veículo pode ser vinculado opcionalmente a um código FIPE (`FipeCode`) para cálculo automático do valor de mercado e histórico de depreciação.

---

## 🔗 Referências Cruzadas
- [Catálogo Canônico](../index.md)
- [Isolamento de Dados e Persistência](../architecture/multitenancy-and-data.md)
- [Subdomínio de Manutenções](maintenance.md)
- [Subdomínio de Abastecimentos](fuel.md)
- [Subdomínio da Tabela FIPE](fipe.md)
