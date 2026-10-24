import { apiClient } from "../../../services/api/client";
import { ApiEnvelope } from "../../auth/types/auth.types";
import {
  Maintenance,
  CreateMaintenancePayload,
  UpdateMaintenancePayload,
  CarUpcomingMaintenances,
  FleetUpcomingSummary,
} from "../types/maintenance.types";

export const maintenanceService = {
  async listByCar(carId: string): Promise<Maintenance[]> {
    const response = await apiClient.get<ApiEnvelope<Maintenance[]>>(
      `/cars/${carId}/maintenances`,
    );
    return response.data.data || [];
  },

  async getUpcomingByCar(carId: string): Promise<CarUpcomingMaintenances> {
    const response = await apiClient.get<ApiEnvelope<CarUpcomingMaintenances>>(
      `/cars/${carId}/upcoming-maintenances`,
    );
    return response.data.data;
  },

  async getUpcomingSummary(): Promise<FleetUpcomingSummary> {
    const response = await apiClient.get<ApiEnvelope<FleetUpcomingSummary>>(
      `/maintenances/upcoming`,
    );
    return response.data.data;
  },

  async get(maintenanceId: string): Promise<Maintenance> {
    const response = await apiClient.get<ApiEnvelope<Maintenance>>(
      `/maintenances/${maintenanceId}`,
    );
    return response.data.data;
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
