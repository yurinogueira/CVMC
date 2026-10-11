import { apiClient } from "../../../services/api/client";
import {
  ApiEnvelope,
  AuthProvidersData,
  User,
} from "../../auth/types/auth.types";
import {
  ProfileData,
  UpdatePasswordPayload,
  UpdateProfilePayload,
} from "../types/profile.types";

export const profileService = {
  async getProfile(): Promise<ProfileData> {
    const response =
      await apiClient.get<ApiEnvelope<ProfileData>>("/user/profile");
    return response.data.data;
  },

  async updateProfile(payload: UpdateProfilePayload): Promise<User> {
    const response = await apiClient.put<ApiEnvelope<{ user: User }>>(
      "/user/profile",
      payload,
    );
    return response.data.data.user;
  },

  async updatePassword(
    payload: UpdatePasswordPayload,
  ): Promise<{ message: string }> {
    const response = await apiClient.put<ApiEnvelope<{ message: string }>>(
      "/user/password",
      payload,
    );
    return response.data.data;
  },

  async getAuthProviders(): Promise<AuthProvidersData> {
    const response = await apiClient.get<ApiEnvelope<AuthProvidersData>>(
      "/user/auth-providers",
    );
    return response.data.data;
  },

  async linkGoogle(idToken: string): Promise<{ message: string }> {
    const response = await apiClient.post<ApiEnvelope<{ message: string }>>(
      "/user/auth-providers/google/link",
      { idToken },
    );
    return response.data.data;
  },

  async unlinkGoogle(): Promise<{ message: string }> {
    const response = await apiClient.post<ApiEnvelope<{ message: string }>>(
      "/user/auth-providers/google/unlink",
    );
    return response.data.data;
  },

  async setPassword(newPassword: string): Promise<{ message: string }> {
    const response = await apiClient.post<ApiEnvelope<{ message: string }>>(
      "/user/password/set",
      { newPassword },
    );
    return response.data.data;
  },
};
