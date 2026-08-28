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
}
