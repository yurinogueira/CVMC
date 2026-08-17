import { create } from 'zustand';

type AuthState = {
  accessToken: string | null;
  refreshToken: string | null;
  setTokens: (tokens: { accessToken: string; refreshToken: string }) => void;
  clear: () => void;
};

export const useAuthStore = create<AuthState>((set) => ({
  accessToken: null,
  refreshToken: null,
  setTokens: ({ accessToken, refreshToken }) => {
    localStorage.setItem('cvmc.accessToken', accessToken);
    localStorage.setItem('cvmc.refreshToken', refreshToken);
    set({ accessToken, refreshToken });
  },
  clear: () => {
    localStorage.removeItem('cvmc.accessToken');
    localStorage.removeItem('cvmc.refreshToken');
    set({ accessToken: null, refreshToken: null });
  },
}));
