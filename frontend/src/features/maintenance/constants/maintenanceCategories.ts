import { OTHER_MAINTENANCE_TYPE } from "./maintenanceTypes";

export interface MaintenanceCategoryGroup {
  id: string;
  name: string;
  items: string[];
}

export const MAINTENANCE_CATEGORIES: MaintenanceCategoryGroup[] = [
  {
    id: "engine_filters",
    name: "Motor & Filtros",
    items: [
      "Óleo de Motor",
      "Filtro do Óleo de Motor",
      "Filtro de Ar do Motor",
      "Filtro de Combustível",
      "Velas de Ignição",
    ],
  },
  {
    id: "fluids",
    name: "Fluidos & Arrefecimento",
    items: [
      "Fluido de Arrefecimento",
      "Fluido de Freio",
      "Fluido da Embreagem",
      "Fluido da Direção Hidráulica",
      "Fluido da Transmissão",
    ],
  },
  {
    id: "brakes_suspension_tires",
    name: "Freios, Suspensão & Pneus",
    items: [
      "Pastilhas de Freio",
      "Pneus",
      "Alinhamento e Balanceamento",
      "Suspensão e Amortecedores",
    ],
  },
  {
    id: "electrical_comfort_others",
    name: "Elétrica, Conforto & Outros",
    items: [
      "Bateria",
      "Filtro do Ar-Condicionado",
      "Palhetas do Limpador",
      "Revisão Geral / Preventiva",
      OTHER_MAINTENANCE_TYPE,
    ],
  },
];
