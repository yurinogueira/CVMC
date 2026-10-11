import { Preferences } from "@capacitor/preferences";

export const safeStorage = {
  getItem(key: string): string | null {
    if (typeof window === "undefined" || !window.localStorage) return null;
    try {
      return window.localStorage.getItem(key);
    } catch {
      return null;
    }
  },
  setItem(key: string, value: string): void {
    if (typeof window === "undefined" || !window.localStorage) return;
    try {
      window.localStorage.setItem(key, value);
    } catch {
      // ignore
    }
  },
  removeItem(key: string): void {
    if (typeof window === "undefined" || !window.localStorage) return;
    try {
      window.localStorage.removeItem(key);
    } catch {
      // ignore
    }
  },

  // Asynchronous token storage using Capacitor Preferences for native mobile persistence
  async getSecureItem(key: string): Promise<string | null> {
    try {
      const { value } = await Preferences.get({ key });
      if (value !== null) return value;
    } catch {
      // fallback to safeStorage.getItem
    }
    return this.getItem(key);
  },

  async setSecureItem(key: string, value: string): Promise<void> {
    this.setItem(key, value);
    try {
      await Preferences.set({ key, value });
    } catch {
      // ignore
    }
  },

  async removeSecureItem(key: string): Promise<void> {
    this.removeItem(key);
    try {
      await Preferences.remove({ key });
    } catch {
      // ignore
    }
  },
};
