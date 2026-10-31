export interface MaintenanceIntervalConfig {
  intervalKm: number;
  intervalMonths: number;
}

export interface Car {
  id: string;
  name: string;
  manufacturer: string;
  model: string;
  yearManufacture: number;
  yearModel: number;
  lastMileage: number;
  vehicleType?: string;
  imageUrl?: string;
  fipeCode?: string;
  fipePrice?: string;
  fuel?: string;
  ownerId: string;
  sharedWith?: string[];
  maintenanceIntervals?: Record<string, MaintenanceIntervalConfig>;
  mileageUpdatedAt?: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface CreateCarPayload {
  name: string;
  manufacturer: string;
  model: string;
  yearManufacture: number;
  yearModel: number;
  lastMileage: number;
  vehicleType?: string;
  imageUrl?: string;
  fipeCode?: string;
  fipePrice?: string;
  fuel?: string;
  maintenanceIntervals?: Record<string, MaintenanceIntervalConfig>;
}

export interface UpdateCarPayload {
  name: string;
  manufacturer: string;
  model: string;
  yearManufacture: number;
  yearModel: number;
  lastMileage: number;
  vehicleType?: string;
  imageUrl?: string;
  fipeCode?: string;
  fipePrice?: string;
  fuel?: string;
  maintenanceIntervals?: Record<string, MaintenanceIntervalConfig>;
}
