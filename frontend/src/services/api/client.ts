import axios, { AxiosError, InternalAxiosRequestConfig } from "axios";
import { Capacitor } from "@capacitor/core";
import { useAuthStore } from "../../features/auth/state/auth.store";
import {
  ApiEnvelope,
  AuthResponseData,
} from "../../features/auth/types/auth.types";

const resolveBaseURL = (): string => {
  if (import.meta.env.VITE_API_BASE_URL) {
    return `${import.meta.env.VITE_API_BASE_URL}/api/v1`;
  }
  // In native mobile app (Capacitor), default to production API if not explicitly overridden
  if (Capacitor.isNativePlatform()) {
    return "https://api-cvmc.yurinogueira.dev.br/api/v1";
  }
  return "/api/v1";
};

export const apiClient = axios.create({
  baseURL: resolveBaseURL(),
  timeout: 15000,
  withCredentials: true,
});

apiClient.interceptors.request.use((config) => {
  const token = useAuthStore.getState().accessToken;
  if (token && !config.headers.Authorization) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

interface CustomAxiosRequestConfig extends InternalAxiosRequestConfig {
  _retry?: boolean;
}

interface QueuedPromise {
  resolve: (value?: unknown) => void;
  reject: (reason?: unknown) => void;
}

let isRefreshing = false;
let failedQueue: QueuedPromise[] = [];

const processQueue = (error: unknown = null) => {
  failedQueue.forEach((promise) => {
    if (error) {
      promise.reject(error);
    } else {
      promise.resolve();
    }
  });
  failedQueue = [];
};

export const isAuthEndpoint = (url?: string): boolean => {
  if (!url) return false;
  return (
    url.includes("/auth/login") ||
    url.includes("/auth/register") ||
    url.includes("/auth/refresh") ||
    url.includes("/auth/logout")
  );
};

export const handleResponseRejected = async (error: AxiosError) => {
  const originalRequest = error.config as CustomAxiosRequestConfig | undefined;
  const status = error?.response?.status;

  if (status !== 401 || !originalRequest) {
    return Promise.reject(error);
  }

  // Skip silent refresh for authentication endpoints (login, register, refresh, logout)
  if (isAuthEndpoint(originalRequest.url)) {
    if (originalRequest.url?.includes("/auth/refresh")) {
      useAuthStore.getState().clear();
      const currentPath =
        typeof window !== "undefined" ? window.location.pathname : "";
      if (
        currentPath &&
        currentPath !== "/login" &&
        currentPath !== "/register"
      ) {
        window.location.href = "/login";
      }
    }
    return Promise.reject(error);
  }

  // If the request was already retried and failed again with 401
  if (originalRequest._retry) {
    useAuthStore.getState().clear();
    const currentPath =
      typeof window !== "undefined" ? window.location.pathname : "";
    if (
      currentPath &&
      currentPath !== "/login" &&
      currentPath !== "/register"
    ) {
      window.location.href = "/login";
    }
    return Promise.reject(error);
  }

  // If a refresh operation is currently in progress, queue this request
  if (isRefreshing) {
    return new Promise((resolve, reject) => {
      failedQueue.push({ resolve, reject });
    })
      .then(() => {
        originalRequest._retry = true;
        return apiClient(originalRequest);
      })
      .catch((err) => Promise.reject(err));
  }

  originalRequest._retry = true;
  isRefreshing = true;

  try {
    const currentRefreshToken = useAuthStore.getState().refreshToken;
    const refreshHeaders: Record<string, string> = {};
    if (currentRefreshToken) {
      refreshHeaders.Authorization = `Bearer ${currentRefreshToken}`;
    }

    const refreshResponse = await apiClient.post<ApiEnvelope<AuthResponseData>>(
      "/auth/refresh",
      {},
      { headers: refreshHeaders },
    );

    if (refreshResponse.data?.data) {
      const { user, accessToken, refreshToken } = refreshResponse.data.data;
      if (accessToken) {
        useAuthStore.getState().setTokens({ accessToken, refreshToken });
      }
      if (user) {
        useAuthStore.getState().setUser(user);
      }
    }

    processQueue(null);
    return apiClient(originalRequest);
  } catch (refreshErr) {
    processQueue(refreshErr);
    useAuthStore.getState().clear();
    const currentPath =
      typeof window !== "undefined" ? window.location.pathname : "";
    if (
      currentPath &&
      currentPath !== "/login" &&
      currentPath !== "/register"
    ) {
      window.location.href = "/login";
    }
    return Promise.reject(refreshErr);
  } finally {
    isRefreshing = false;
  }
};

apiClient.interceptors.response.use(
  (response) => response,
  handleResponseRejected,
);
