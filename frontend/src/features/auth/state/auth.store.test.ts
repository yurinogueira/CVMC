import { describe, it, expect, beforeEach, vi } from "vitest";
import { Capacitor } from "@capacitor/core";
import { Preferences } from "@capacitor/preferences";
import { useAuthStore, initAuthSession } from "./auth.store";

vi.mock("@capacitor/preferences", () => ({
  Preferences: {
    get: vi.fn(),
    set: vi.fn(),
    remove: vi.fn(),
  },
}));

describe("auth.store session persistence and initialization", () => {
  beforeEach(() => {
    localStorage.clear();
    useAuthStore.getState().clear();
    useAuthStore.setState({ isInitializing: true });
    vi.restoreAllMocks();
  });

  it("initializes with isInitializing = true", () => {
    expect(useAuthStore.getState().isInitializing).toBe(true);
  });

  it("restores user and tokens from Capacitor Preferences on native mobile", async () => {
    vi.spyOn(Capacitor, "isNativePlatform").mockReturnValue(true);

    const mockUser = {
      id: "u-mobile-1",
      name: "Motorista Mobile",
      email: "motorista@cvmc.com",
    };

    vi.mocked(Preferences.get).mockImplementation(async ({ key }) => {
      if (key === "cvmc.user") return { value: JSON.stringify(mockUser) };
      if (key === "cvmc.access_token") return { value: "native-access-token" };
      if (key === "cvmc.refresh_token")
        return { value: "native-refresh-token" };
      return { value: null };
    });

    await initAuthSession();

    const state = useAuthStore.getState();
    expect(state.isInitializing).toBe(false);
    expect(state.isAuthenticated).toBe(true);
    expect(state.user).toEqual(mockUser);
    expect(state.accessToken).toBe("native-access-token");
    expect(state.refreshToken).toBe("native-refresh-token");
  });

  it("handles empty Preferences on native mobile without authenticating", async () => {
    vi.spyOn(Capacitor, "isNativePlatform").mockReturnValue(true);
    vi.mocked(Preferences.get).mockResolvedValue({ value: null });

    await initAuthSession();

    const state = useAuthStore.getState();
    expect(state.isInitializing).toBe(false);
    expect(state.isAuthenticated).toBe(false);
    expect(state.user).toBeNull();
    expect(state.accessToken).toBeNull();
    expect(state.refreshToken).toBeNull();
  });

  it("restores user from localStorage on web platform and sets isInitializing = false", async () => {
    vi.spyOn(Capacitor, "isNativePlatform").mockReturnValue(false);

    const webUser = {
      id: "u-web-1",
      name: "Condutor Web",
      email: "web@cvmc.com",
    };
    localStorage.setItem("cvmc.user", JSON.stringify(webUser));

    await initAuthSession();

    const state = useAuthStore.getState();
    expect(state.isInitializing).toBe(false);
    expect(state.isAuthenticated).toBe(true);
    expect(state.user).toEqual(webUser);
    // On web, tokens are never in state on initial boot (they travel via HttpOnly cookies)
    expect(state.accessToken).toBeNull();
  });

  it("clears both localStorage and Preferences when clear() is invoked", async () => {
    vi.spyOn(Capacitor, "isNativePlatform").mockReturnValue(true);
    const removeSpy = vi.mocked(Preferences.remove).mockResolvedValue();

    useAuthStore.getState().setSession({
      user: { id: "u-1", name: "User 1", email: "u1@test.com" },
      accessToken: "token-1",
      refreshToken: "token-2",
    });

    useAuthStore.getState().clear();

    const state = useAuthStore.getState();
    expect(state.user).toBeNull();
    expect(state.isAuthenticated).toBe(false);
    expect(state.accessToken).toBeNull();
    expect(state.refreshToken).toBeNull();
    expect(localStorage.getItem("cvmc.user")).toBeNull();
    expect(removeSpy).toHaveBeenCalledWith({ key: "cvmc.user" });
    expect(removeSpy).toHaveBeenCalledWith({ key: "cvmc.access_token" });
    expect(removeSpy).toHaveBeenCalledWith({ key: "cvmc.refresh_token" });
  });
});
