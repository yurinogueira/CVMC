import { create } from "zustand";
import { User } from "../types/auth.types";
import { safeStorage } from "../../../services/storage/storage";

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  setUser: (user: User) => void;
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
    setUser: (user) => {
      safeStorage.setItem("cvmc.user", JSON.stringify(user));
      set({ user, isAuthenticated: true });
    },
    clear: () => {
      safeStorage.removeItem("cvmc.user");
      set({ user: null, isAuthenticated: false });
    },
  };
});
