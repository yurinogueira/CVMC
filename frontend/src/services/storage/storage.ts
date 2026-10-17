import { Capacitor } from "@capacitor/core";
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

  // Asynchronous token storage using Capacitor Preferences strictly for native mobile persistence.
  // On web browsers, tokens are never stored in localStorage to preserve HttpOnly cookie isolation (ADR 0002).
  async getSecureItem(key: string): Promise<string | null> {
    if (!Capacitor.isNativePlatform()) {
      return null;
    }
    try {
      const { value } = await Preferences.get({ key });
      return value;
    } catch {
      return null;
    }
  },

  async setSecureItem(key: string, value: string): Promise<void> {
    if (!Capacitor.isNativePlatform()) {
      return;
    }
    try {
      await Preferences.set({ key, value });
    } catch {
      // ignore
    }
  },

  async removeSecureItem(key: string): Promise<void> {
    if (!Capacitor.isNativePlatform()) {
      this.removeItem(key);
      return;
    }
    try {
      await Preferences.remove({ key });
    } catch {
      // ignore
    }
  },
};
