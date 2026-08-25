import { describe, it, expect, vi, beforeEach } from "vitest";
import { isAuthEndpoint, handleResponseRejected, apiClient } from "./client";
import { useAuthStore } from "../../features/auth/state/auth.store";
import { AxiosError, InternalAxiosRequestConfig } from "axios";

describe("apiClient and silent refresh interceptor", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    useAuthStore.getState().clear();

    Object.defineProperty(window, "location", {
      value: {
        pathname: "/vehicles",
        href: "http://localhost:5173/vehicles",
      },
      writable: true,
    });
  });

  it("identifies authentication endpoints correctly", () => {
    expect(isAuthEndpoint("/auth/login")).toBe(true);
    expect(isAuthEndpoint("/api/v1/auth/login")).toBe(true);
    expect(isAuthEndpoint("/auth/register")).toBe(true);
    expect(isAuthEndpoint("/auth/refresh")).toBe(true);
    expect(isAuthEndpoint("/auth/logout")).toBe(true);
    expect(isAuthEndpoint("/cars")).toBe(false);
    expect(isAuthEndpoint("/maintenance")).toBe(false);
    expect(isAuthEndpoint(undefined)).toBe(false);
  });

  it("does not attempt refresh on auth login endpoint when receiving 401", async () => {
    const postSpy = vi.spyOn(apiClient, "post");

    const error = {
      response: { status: 401 },
      config: { url: "/auth/login", headers: {} },
    } as unknown as AxiosError;

    await expect(handleResponseRejected(error)).rejects.toBeDefined();
    expect(postSpy).not.toHaveBeenCalledWith("/auth/refresh");
  });

  it("clears user and redirects to /login when /auth/refresh fails with 401", async () => {
    useAuthStore.getState().setUser({
      id: "u-1",
      name: "Test User",
      email: "test@example.com",
    });
    expect(useAuthStore.getState().isAuthenticated).toBe(true);

    const error = {
      response: { status: 401 },
      config: { url: "/auth/refresh", headers: {} },
    } as unknown as AxiosError;

    await expect(handleResponseRejected(error)).rejects.toBeDefined();

    expect(useAuthStore.getState().isAuthenticated).toBe(false);
    expect(window.location.href).toBe("/login");
  });

  it("performs silent refresh, updates store and retries original request upon 401", async () => {
    const refreshData = {
      success: true,
      data: {
        user: {
          id: "u-2",
          name: "Refreshed User",
          email: "refresh@example.com",
        },
      },
    };

    apiClient.defaults.adapter = async (config: InternalAxiosRequestConfig) => {
      if (config.url?.includes("/auth/refresh")) {
        return {
          data: refreshData,
          status: 200,
          statusText: "OK",
          headers: {},
          config,
        };
      }
      return {
        data: { success: true, data: [{ id: "car-1", brand: "Honda" }] },
        status: 200,
        statusText: "OK",
        headers: {},
        config,
      };
    };

    const error = {
      response: { status: 401 },
      config: { url: "/cars", method: "get", headers: {} },
    } as unknown as AxiosError;

    const response = (await handleResponseRejected(error)) as {
      data: { success: boolean };
    };

    expect(response.data.success).toBe(true);
    expect(useAuthStore.getState().user?.name).toBe("Refreshed User");
    expect(useAuthStore.getState().isAuthenticated).toBe(true);
  });

  it("redirects and clears session if retried request fails again with 401", async () => {
    useAuthStore.getState().setUser({
      id: "u-3",
      name: "Retry User",
      email: "retry@example.com",
    });

    const error = {
      response: { status: 401 },
      config: { url: "/cars", _retry: true, headers: {} },
    } as unknown as AxiosError;

    await expect(handleResponseRejected(error)).rejects.toBeDefined();

    expect(useAuthStore.getState().isAuthenticated).toBe(false);
    expect(window.location.href).toBe("/login");
  });
});
