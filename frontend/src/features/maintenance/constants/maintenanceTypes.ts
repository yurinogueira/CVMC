export const OTHER_MAINTENANCE_TYPE = "Outro tipo de manutenção";

export const MAINTENANCE_TYPES = [
  "Óleo de Motor",
  "Filtro do Óleo de Motor",
  "Filtro de Ar do Motor",
  "Filtro de Combustível",
  "Filtro do Ar-Condicionado",
  "Fluido de Arrefecimento",
  "Pastilhas de Freio",
  "Fluido de Freio",
  "Fluido da Embreagem",
  "Fluido da Direção Hidráulica",
  "Fluido da Transmissão",
  "Velas de Ignição",
  "Pneus",
  "Alinhamento e Balanceamento",
  "Bateria",
  "Suspensão e Amortecedores",
  "Palhetas do Limpador",
  "Revisão Geral / Preventiva",
  OTHER_MAINTENANCE_TYPE,
] as const;

export type MaintenanceTypeOption = (typeof MAINTENANCE_TYPES)[number];
