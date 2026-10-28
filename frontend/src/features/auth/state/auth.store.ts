import { create } from "zustand";
import { Capacitor } from "@capacitor/core";
import { User } from "../types/auth.types";
import { safeStorage } from "../../../services/storage/storage";

export interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isInitializing: boolean;
  accessToken: string | null;
  refreshToken: string | null;
  setUser: (user: User) => void;
  setSession: (session: {
    user: User;
    accessToken?: string;
    refreshToken?: string;
  }) => void;
  setTokens: (tokens: { accessToken: string; refreshToken?: string }) => void;
  clear: () => void;
  setInitializing: (isInitializing: boolean) => void;
}

const getStoredUser = (): User | null => {
  const userStr = safeStorage.getItem("cvmc.user");
  if (userStr) {
    try {
      return JSON.parse(userStr);
    } catch {
      return null;
    }
  }
  return null;
};

export const useAuthStore = create<AuthState>((set) => {
  const user = getStoredUser();

  return {
    user,
    isAuthenticated: Boolean(user),
    isInitializing: true,
    accessToken: null,
    refreshToken: null,
    setInitializing: (isInitializing) => set({ isInitializing }),
    setUser: (user) => {
      safeStorage.setItem("cvmc.user", JSON.stringify(user));
      void safeStorage.setSecureItem("cvmc.user", JSON.stringify(user));
      set({ user, isAuthenticated: true });
    },
    setSession: ({ user, accessToken, refreshToken }) => {
      safeStorage.setItem("cvmc.user", JSON.stringify(user));
      void safeStorage.setSecureItem("cvmc.user", JSON.stringify(user));
      if (accessToken) {
        void safeStorage.setSecureItem("cvmc.access_token", accessToken);
      }
      if (refreshToken) {
        void safeStorage.setSecureItem("cvmc.refresh_token", refreshToken);
      }
      set((state) => ({
        user,
        isAuthenticated: true,
        accessToken: accessToken ?? state.accessToken,
        refreshToken: refreshToken ?? state.refreshToken,
      }));
    },
    setTokens: ({ accessToken, refreshToken }) => {
      void safeStorage.setSecureItem("cvmc.access_token", accessToken);
      if (refreshToken) {
        void safeStorage.setSecureItem("cvmc.refresh_token", refreshToken);
      }
      set((state) => ({
        accessToken,
        refreshToken: refreshToken ?? state.refreshToken,
      }));
    },
    clear: () => {
      safeStorage.removeItem("cvmc.user");
      void safeStorage.removeSecureItem("cvmc.user");
      void safeStorage.removeSecureItem("cvmc.access_token");
      void safeStorage.removeSecureItem("cvmc.refresh_token");
      set({
        user: null,
        isAuthenticated: false,
        accessToken: null,
        refreshToken: null,
      });
    },
  };
});

let initPromise: Promise<void> | null = null;

export const initAuthSession = async (): Promise<void> => {
  if (initPromise) {
    return initPromise;
  }

  initPromise = (async () => {
    try {
      if (Capacitor.isNativePlatform()) {
        const [storedUserStr, accessToken, refreshToken] = await Promise.all([
          safeStorage.getSecureItem("cvmc.user"),
          safeStorage.getSecureItem("cvmc.access_token"),
          safeStorage.getSecureItem("cvmc.refresh_token"),
        ]);

        let parsedUser: User | null = null;
        if (storedUserStr) {
          try {
            parsedUser = JSON.parse(storedUserStr);
          } catch {
            parsedUser = null;
          }
        }
        if (!parsedUser) {
          parsedUser = getStoredUser();
        }

        if (accessToken && parsedUser) {
          useAuthStore.setState({
            user: parsedUser,
            isAuthenticated: true,
            accessToken,
            refreshToken: refreshToken ?? null,
          });
        } else if (accessToken) {
          useAuthStore.setState({
            accessToken,
            refreshToken: refreshToken ?? null,
            isAuthenticated: true,
          });
        } else {
          useAuthStore.setState({
            user: null,
            isAuthenticated: false,
            accessToken: null,
            refreshToken: null,
          });
        }
      } else {
        const user = getStoredUser();
        if (user) {
          useAuthStore.setState({
            user,
            isAuthenticated: true,
          });
        } else {
          useAuthStore.setState({
            user: null,
            isAuthenticated: false,
          });
        }
      }
    } catch {
      // Graceful fallback on initialization error
    } finally {
      useAuthStore.setState({ isInitializing: false });
      initPromise = null;
    }
  })();

  return initPromise;
};

export const initMobileAuthSession = initAuthSession;
