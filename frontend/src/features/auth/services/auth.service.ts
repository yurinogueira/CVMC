import { apiClient } from "../../../services/api/client";
import {
  ApiEnvelope,
  AuthResponseData,
  LoginPayload,
  RegisterPayload,
  User,
} from "../types/auth.types";

export const authService = {
  async login(payload: LoginPayload): Promise<AuthResponseData> {
    const response = await apiClient.post<ApiEnvelope<AuthResponseData>>(
      "/auth/login",
      payload,
    );
    return response.data.data;
  },

  async register(payload: RegisterPayload): Promise<AuthResponseData> {
    const response = await apiClient.post<ApiEnvelope<AuthResponseData>>(
      "/auth/register",
      payload,
    );
    return response.data.data;
  },

  async getMe(): Promise<User> {
    const response = await apiClient.get<ApiEnvelope<User>>("/auth/me");
    return response.data.data;
  },

  async logout(): Promise<void> {
    await apiClient.post("/auth/logout");
  },

  async refresh(): Promise<AuthResponseData> {
    const response =
      await apiClient.post<ApiEnvelope<AuthResponseData>>("/auth/refresh");
    return response.data.data;
  },
};
