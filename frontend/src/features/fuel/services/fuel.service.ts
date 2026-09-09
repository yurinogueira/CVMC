import { apiClient } from "../../../services/api/client";
import { ApiEnvelope } from "../../auth/types/auth.types";
import {
  Fueling,
  CreateFuelingInput,
  UpdateFuelingInput,
} from "../types/fuel.types";

export const fuelService = {
  async listByCar(carId: string): Promise<Fueling[]> {
    const response = await apiClient.get<ApiEnvelope<Fueling[]>>(
      `/cars/${carId}/fuelings`,
    );
    return response.data.data || [];
  },

  async get(fuelingId: string): Promise<Fueling> {
    const response = await apiClient.get<ApiEnvelope<Fueling>>(
      `/fuelings/${fuelingId}`,
    );
    return response.data.data;
  },

  async create(carId: string, payload: CreateFuelingInput): Promise<Fueling> {
    const response = await apiClient.post<ApiEnvelope<Fueling>>(
      `/cars/${carId}/fuelings`,
      payload,
    );
    return response.data.data;
  },

  async update(
    fuelingId: string,
    payload: UpdateFuelingInput,
  ): Promise<Fueling> {
    const response = await apiClient.put<ApiEnvelope<Fueling>>(
      `/fuelings/${fuelingId}`,
      payload,
    );
    return response.data.data;
  },

  async delete(fuelingId: string): Promise<void> {
    await apiClient.delete(`/fuelings/${fuelingId}`);
  },
};
