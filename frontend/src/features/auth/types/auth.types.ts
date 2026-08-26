export interface User {
  id: string;
  name: string;
  email: string;
  emailVerified?: boolean;
  emailVerifiedAt?: string;
  maxVehicles?: number;
  createdAt?: string;
  updatedAt?: string;
}

export interface AuthTokens {
  accessToken: string;
  refreshToken: string;
}

export interface AuthResponseData {
  user: User;
}

export interface ApiEnvelope<T> {
  success: boolean;
  data: T;
  message?: string;
  errors?: unknown;
}

export interface LoginPayload {
  email: string;
  password: string;
}

export interface RegisterPayload {
  name: string;
  email: string;
  password: string;
}

export interface ForgotPasswordPayload {
  email: string;
}

export interface ResetPasswordPayload {
  token: string;
  newPassword: string;
}

export interface VerifyEmailPayload {
  token: string;
}
