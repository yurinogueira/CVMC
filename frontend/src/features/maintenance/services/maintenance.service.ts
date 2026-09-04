import { apiClient } from "../../../services/api/client";
import { ApiEnvelope } from "../../auth/types/auth.types";
import {
  Maintenance,
  CreateMaintenancePayload,
  UpdateMaintenancePayload,
} from "../types/maintenance.types";

export const maintenanceService = {
  async listByCar(carId: string): Promise<Maintenance[]> {
    const response = await apiClient.get<ApiEnvelope<Maintenance[]>>(
      `/cars/${carId}/maintenances`,
    );
    return response.data.data || [];
  },

  async create(
    carId: string,
    payload: CreateMaintenancePayload,
  ): Promise<Maintenance> {
    const response = await apiClient.post<ApiEnvelope<Maintenance>>(
      `/cars/${carId}/maintenances`,
      payload,
    );
    return response.data.data;
  },

  async update(
    maintenanceId: string,
    payload: UpdateMaintenancePayload,
  ): Promise<Maintenance> {
    const response = await apiClient.put<ApiEnvelope<Maintenance>>(
      `/maintenances/${maintenanceId}`,
      payload,
    );
    return response.data.data;
  },

  async delete(maintenanceId: string): Promise<void> {
    await apiClient.delete(`/maintenances/${maintenanceId}`);
  },
};
