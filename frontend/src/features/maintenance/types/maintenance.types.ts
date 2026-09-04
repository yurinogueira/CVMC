export interface Maintenance {
  id: string;
  carId: string;
  title: string;
  description: string;
  date: string;
  mileage: number;
  createdAt?: string;
  updatedAt?: string;
}

export interface CreateMaintenancePayload {
  title: string;
  description?: string;
  date: string;
  mileage: number;
}

export interface UpdateMaintenancePayload {
  title: string;
  description?: string;
  date: string;
  mileage: number;
}
