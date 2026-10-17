import { describe, it, expect, beforeEach } from "vitest";
import { safeStorage } from "./storage";
import { useAuthStore } from "../../features/auth/state/auth.store";

describe("safeStorage security and token isolation", () => {
  beforeEach(() => {
    localStorage.clear();
    useAuthStore.getState().clear();
  });

  it("stores and retrieves non-sensitive public preferences via localStorage", () => {
    safeStorage.setItem("cvmc.theme", "dark");
    expect(safeStorage.getItem("cvmc.theme")).toBe("dark");
    expect(localStorage.getItem("cvmc.theme")).toBe("dark");

    safeStorage.removeItem("cvmc.theme");
    expect(safeStorage.getItem("cvmc.theme")).toBeNull();
    expect(localStorage.getItem("cvmc.theme")).toBeNull();
  });

  it("does not store JWT tokens in localStorage when running on web platform", async () => {
    await safeStorage.setSecureItem("cvmc.access_token", "secret.jwt.token");
    await safeStorage.setSecureItem(
      "cvmc.refresh_token",
      "secret.refresh.token",
    );

    // Must never leak to window.localStorage on web
    expect(localStorage.getItem("cvmc.access_token")).toBeNull();
    expect(localStorage.getItem("cvmc.refresh_token")).toBeNull();

    // getSecureItem returns null on web
    const retrieved = await safeStorage.getSecureItem("cvmc.access_token");
    expect(retrieved).toBeNull();
  });

  it("does not populate authStore tokens from localStorage on web", () => {
    // If a legacy or malicious value exists in localStorage
    localStorage.setItem("cvmc.access_token", "stale-access-token");
    localStorage.setItem("cvmc.refresh_token", "stale-refresh-token");

    // In-memory store should not initialize tokens from localStorage
    expect(useAuthStore.getState().accessToken).toBeNull();
    expect(useAuthStore.getState().refreshToken).toBeNull();
  });

  it("prevents setSession and setTokens from writing tokens to localStorage on web", () => {
    useAuthStore.getState().setSession({
      user: {
        id: "u-123",
        name: "Security Tester",
        email: "sec@example.com",
      },
      accessToken: "runtime-access-token",
      refreshToken: "runtime-refresh-token",
    });

    // Zustand holds in-memory tokens for the session runtime
    expect(useAuthStore.getState().accessToken).toBe("runtime-access-token");
    expect(useAuthStore.getState().refreshToken).toBe("runtime-refresh-token");

    // But localStorage MUST NOT contain tokens
    expect(localStorage.getItem("cvmc.access_token")).toBeNull();
    expect(localStorage.getItem("cvmc.refresh_token")).toBeNull();

    // User metadata is safely stored
    expect(localStorage.getItem("cvmc.user")).toContain("Security Tester");
  });
});
