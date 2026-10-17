import { create } from "zustand";
import { User } from "../types/auth.types";
import { safeStorage } from "../../../services/storage/storage";

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
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
    accessToken: null,
    refreshToken: null,
    setUser: (user) => {
      safeStorage.setItem("cvmc.user", JSON.stringify(user));
      set({ user, isAuthenticated: true });
    },
    setSession: ({ user, accessToken, refreshToken }) => {
      safeStorage.setItem("cvmc.user", JSON.stringify(user));
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

export const initMobileAuthSession = async (): Promise<void> => {
  const accessToken = await safeStorage.getSecureItem("cvmc.access_token");
  const refreshToken = await safeStorage.getSecureItem("cvmc.refresh_token");
  if (accessToken) {
    useAuthStore.getState().setTokens({
      accessToken,
      refreshToken: refreshToken ?? undefined,
    });
  }
};
