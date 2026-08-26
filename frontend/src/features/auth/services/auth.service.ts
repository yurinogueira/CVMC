import { apiClient } from "../../../services/api/client";
import {
  ApiEnvelope,
  AuthResponseData,
  ForgotPasswordPayload,
  LoginPayload,
  RegisterPayload,
  ResetPasswordPayload,
  User,
  VerifyEmailPayload,
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

  async forgotPassword(
    payload: ForgotPasswordPayload,
  ): Promise<{ message: string }> {
    const response = await apiClient.post<ApiEnvelope<{ message: string }>>(
      "/auth/forgot-password",
      payload,
    );
    return response.data.data;
  },

  async resetPassword(
    payload: ResetPasswordPayload,
  ): Promise<{ message: string }> {
    const response = await apiClient.post<ApiEnvelope<{ message: string }>>(
      "/auth/reset-password",
      payload,
    );
    return response.data.data;
  },

  async verifyEmail(payload: VerifyEmailPayload): Promise<{ message: string }> {
    const response = await apiClient.post<ApiEnvelope<{ message: string }>>(
      "/auth/verify-email",
      payload,
    );
    return response.data.data;
  },

  async resendVerification(): Promise<{ message: string }> {
    const response = await apiClient.post<ApiEnvelope<{ message: string }>>(
      "/auth/resend-verification",
    );
    return response.data.data;
  },
};
