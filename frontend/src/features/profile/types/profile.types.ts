import { User } from "../../auth/types/auth.types";

export interface ProfileData {
  user: User;
  vehiclesCount: number;
  maxVehicles: number;
  filesCount?: number;
  maxFiles?: number;
  storageUsedBytes?: number;
  maxStorageBytes?: number;
}

export interface UpdateProfilePayload {
  name: string;
}

export interface UpdatePasswordPayload {
  currentPassword: string;
  newPassword: string;
}

export interface DeleteAccountPayload {
  currentPassword?: string;
  confirmation?: string;
}
