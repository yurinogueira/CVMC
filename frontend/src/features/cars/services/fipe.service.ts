import { apiClient } from "../../../services/api/client";
import { ApiEnvelope } from "../../auth/types/auth.types";
import {
  FipeBrand,
  FipeModel,
  FipeVehicleDetail,
  FipeYear,
  VehicleType,
} from "../types/fipe.types";

export const fipeService = {
  async getBrands(vehicleType: VehicleType = "cars"): Promise<FipeBrand[]> {
    const response = await apiClient.get<ApiEnvelope<FipeBrand[]>>(
      `/fipe/${vehicleType}/brands`,
    );
    return response.data.data || [];
  },

  async getModels(
    vehicleType: VehicleType,
    brandId: string,
  ): Promise<FipeModel[]> {
    const response = await apiClient.get<ApiEnvelope<FipeModel[]>>(
      `/fipe/${vehicleType}/brands/${brandId}/models`,
    );
    return response.data.data || [];
  },

  async getYears(
    vehicleType: VehicleType,
    brandId: string,
    modelId: string,
  ): Promise<FipeYear[]> {
    const response = await apiClient.get<ApiEnvelope<FipeYear[]>>(
      `/fipe/${vehicleType}/brands/${brandId}/models/${modelId}/years`,
    );
    return response.data.data || [];
  },

  async getVehicleDetail(
    vehicleType: VehicleType,
    brandId: string,
    modelId: string,
    yearId: string,
  ): Promise<FipeVehicleDetail> {
    const response = await apiClient.get<ApiEnvelope<FipeVehicleDetail>>(
      `/fipe/${vehicleType}/brands/${brandId}/models/${modelId}/years/${yearId}`,
    );
    return response.data.data;
  },
};
