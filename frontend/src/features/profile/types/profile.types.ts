import { User } from "../../auth/types/auth.types";

export interface ProfileData {
  user: User;
  vehiclesCount: number;
  maxVehicles: number;
}

export interface UpdateProfilePayload {
  name: string;
}

export interface UpdatePasswordPayload {
  currentPassword: string;
  newPassword: string;
}
