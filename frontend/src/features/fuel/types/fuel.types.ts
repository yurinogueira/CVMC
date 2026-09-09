export interface Fueling {
  id: string;
  carId: string;
  date: string;
  fuelType: string;
  liters: number;
  pricePerLiter: number;
  totalCost: number;
  isFullTank: boolean;
  gasStation?: string;
  notes?: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface CreateFuelingInput {
  date?: string;
  fuelType: string;
  liters: number;
  pricePerLiter?: number;
  totalCost?: number;
  isFullTank?: boolean;
  gasStation?: string;
  notes?: string;
}

export interface UpdateFuelingInput {
  date?: string;
  fuelType: string;
  liters: number;
  pricePerLiter?: number;
  totalCost?: number;
  isFullTank?: boolean;
  gasStation?: string;
  notes?: string;
}

export const FUEL_TYPES = [
  "Gasolina Comum",
  "Gasolina Aditivada",
  "Etanol",
  "Diesel S10",
  "Diesel Comum",
  "GNV",
  "Elétrico",
] as const;
