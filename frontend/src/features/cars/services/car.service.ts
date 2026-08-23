import { apiClient } from "../../../services/api/client";
import { ApiEnvelope } from "../../auth/types/auth.types";
import { Car, CreateCarPayload, UpdateCarPayload } from "../types/car.types";

export const carService = {
  async list(): Promise<Car[]> {
    const response = await apiClient.get<ApiEnvelope<Car[]>>("/cars");
    return response.data.data || [];
  },

  async get(id: string): Promise<Car> {
    const response = await apiClient.get<ApiEnvelope<Car>>(`/cars/${id}`);
    return response.data.data;
  },

  async create(payload: CreateCarPayload): Promise<Car> {
    const response = await apiClient.post<ApiEnvelope<Car>>("/cars", payload);
    return response.data.data;
  },

  async update(id: string, payload: UpdateCarPayload): Promise<Car> {
    const response = await apiClient.put<ApiEnvelope<Car>>(
      `/cars/${id}`,
      payload,
    );
    return response.data.data;
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/cars/${id}`);
  },
};
