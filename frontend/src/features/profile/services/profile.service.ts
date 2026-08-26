import { apiClient } from "../../../services/api/client";
import { ApiEnvelope, User } from "../../auth/types/auth.types";
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
};
