export interface PeriodicItemDefinition {
  type: string;
  label: string;
  description: string;
  supportedVehicleTypes: ("cars" | "trucks" | "motorcycles")[];
  defaultInterval: {
    cars: { intervalKm: number; intervalMonths: number };
    trucks: { intervalKm: number; intervalMonths: number };
    motorcycles?: { intervalKm: number; intervalMonths: number };
  };
}

export const CAR_KM_OPTIONS = [5000, 7500, 10000, 15000, 20000, 30000] as const;
export const MOTO_KM_OPTIONS = [1000, 2500, 5000, 10000] as const;
export const MONTH_OPTIONS = [3, 6, 12, 18, 24] as const;

export const PERIODIC_ITEMS_CATALOG: PeriodicItemDefinition[] = [
  {
    type: "Óleo de Motor",
    label: "Troca de óleo",
    description: "Lubrificação essencial do motor contra atrito e desgaste.",
    supportedVehicleTypes: ["cars", "trucks", "motorcycles"],
    defaultInterval: {
      cars: { intervalKm: 10000, intervalMonths: 12 },
      trucks: { intervalKm: 10000, intervalMonths: 12 },
      motorcycles: { intervalKm: 2500, intervalMonths: 12 },
    },
  },
  {
    type: "Filtro do Óleo de Motor",
    label: "Troca de filtro de óleo",
    description: "Retenção de impurezas e partículas metálicas no óleo.",
    supportedVehicleTypes: ["cars", "trucks", "motorcycles"],
    defaultInterval: {
      cars: { intervalKm: 10000, intervalMonths: 12 },
      trucks: { intervalKm: 10000, intervalMonths: 12 },
      motorcycles: { intervalKm: 2500, intervalMonths: 12 },
    },
  },
  {
    type: "Filtro de Ar do Motor",
    label: "Troca de filtro de ar",
    description: "Garante ar limpo na câmara de combustão e eficiência.",
    supportedVehicleTypes: ["cars", "trucks", "motorcycles"],
    defaultInterval: {
      cars: { intervalKm: 10000, intervalMonths: 12 },
      trucks: { intervalKm: 10000, intervalMonths: 12 },
      motorcycles: { intervalKm: 5000, intervalMonths: 12 },
    },
  },
  {
    type: "Filtro do Ar-Condicionado",
    label: "Troca do filtro de ar-condicionado",
    description: "Filtro de cabine para retenção de ácaros, poeira e odores.",
    supportedVehicleTypes: ["cars", "trucks"],
    defaultInterval: {
      cars: { intervalKm: 20000, intervalMonths: 12 },
      trucks: { intervalKm: 20000, intervalMonths: 12 },
    },
  },
  {
    type: "Filtro de Combustível",
    label: "Troca do filtro de combustível",
    description: "Proteção dos bicos injetores e bomba contra resíduos.",
    supportedVehicleTypes: ["cars", "trucks"],
    defaultInterval: {
      cars: { intervalKm: 20000, intervalMonths: 12 },
      trucks: { intervalKm: 20000, intervalMonths: 12 },
    },
  },
];

export function isMotorcycle(vehicleType?: string): boolean {
  const vt = (vehicleType || "").toLowerCase().trim();
  return (
    vt === "motorcycles" ||
    vt === "motorcycle" ||
    vt === "moto" ||
    vt === "motos"
  );
}

export function getKmOptionsForVehicle(
  vehicleType?: string,
): readonly number[] {
  return isMotorcycle(vehicleType) ? MOTO_KM_OPTIONS : CAR_KM_OPTIONS;
}

export function getMonthOptionsForVehicle(): readonly number[] {
  return MONTH_OPTIONS;
}

export function getPeriodicItemsForVehicle(
  vehicleType?: string,
): PeriodicItemDefinition[] {
  const isMoto = isMotorcycle(vehicleType);
  const targetCategory = isMoto
    ? "motorcycles"
    : vehicleType === "trucks"
      ? "trucks"
      : "cars";
  return PERIODIC_ITEMS_CATALOG.filter((item) =>
    item.supportedVehicleTypes.includes(targetCategory),
  );
}

export function getDefaultIntervalsForVehicle(
  vehicleType?: string,
): Record<string, { intervalKm: number; intervalMonths: number }> {
  const isMoto = isMotorcycle(vehicleType);
  const items = getPeriodicItemsForVehicle(vehicleType);
  const defaults: Record<
    string,
    { intervalKm: number; intervalMonths: number }
  > = {};

  for (const item of items) {
    if (isMoto && item.defaultInterval.motorcycles) {
      defaults[item.type] = { ...item.defaultInterval.motorcycles };
    } else if (vehicleType === "trucks") {
      defaults[item.type] = { ...item.defaultInterval.trucks };
    } else {
      defaults[item.type] = { ...item.defaultInterval.cars };
    }
  }

  return defaults;
}
