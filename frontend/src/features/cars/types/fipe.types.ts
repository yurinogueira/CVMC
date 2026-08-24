export type VehicleType = "cars" | "motorcycles" | "trucks";

export interface FipeBrand {
  code: string;
  name: string;
  vehicleType?: string;
}

export interface FipeModel {
  code: string;
  name: string;
}

export interface FipeYear {
  code: string;
  name: string;
}

export interface FipeVehicleDetail {
  brand: string;
  codeFipe: string;
  fuel: string;
  fuelAcronym?: string;
  model: string;
  modelYear: number;
  price: string;
  priceValue?: number;
  referenceMonth: string;
  vehicleType: number;
}
