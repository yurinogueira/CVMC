import { describe, it, expect, vi, beforeEach } from "vitest";
import {
  pushNotificationService,
  resolveNotificationDeepLink,
} from "./push.service";
import { apiClient } from "../api/client";
import { safeStorage } from "../storage/storage";

vi.mock("../api/client", () => ({
  apiClient: {
    post: vi.fn(),
    delete: vi.fn(),
  },
}));

describe("pushNotificationService", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    safeStorage.removeItem("cvmc.fcm_token");
  });

  describe("resolveNotificationDeepLink", () => {
    it("returns valid internal route from data.route", () => {
      expect(resolveNotificationDeepLink({ route: "/vehicles/abc-123" })).toBe(
        "/vehicles/abc-123",
      );
    });

    it("rejects external or protocol-relative URLs", () => {
      expect(
        resolveNotificationDeepLink({ route: "https://evil.com" }),
      ).toBeUndefined();
      expect(
        resolveNotificationDeepLink({ route: "//evil.com" }),
      ).toBeUndefined();
    });

    it("falls back to /vehicles/:carId when carId is present", () => {
      expect(resolveNotificationDeepLink({ carId: "car_42" })).toBe(
        "/vehicles/car_42",
      );
    });

    it("appends carId query param when route is /maintenance", () => {
      expect(
        resolveNotificationDeepLink({
          route: "/maintenance",
          carId: "car_99",
        }),
      ).toBe("/maintenance?carId=car_99");
    });
  });

  it("registers device token on backend and stores it locally", async () => {
    vi.mocked(apiClient.post).mockResolvedValueOnce({
      data: { success: true },
    });

    await pushNotificationService.registerTokenOnBackend(
      "fcm-token-999",
      "android",
    );

    expect(apiClient.post).toHaveBeenCalledWith("/devices/token", {
      token: "fcm-token-999",
      platform: "android",
    });
    expect(pushNotificationService.getStoredToken()).toBe("fcm-token-999");
  });

  it("unregisters stored token from backend on logout and clears local storage", async () => {
    safeStorage.setItem("cvmc.fcm_token", "fcm-token-999");
    vi.mocked(apiClient.delete).mockResolvedValueOnce({
      data: { success: true },
    });

    await pushNotificationService.unregisterTokenFromBackend();

    expect(apiClient.delete).toHaveBeenCalledWith("/devices/token", {
      data: { token: "fcm-token-999" },
    });
    expect(pushNotificationService.getStoredToken()).toBeNull();
  });

  it("dispatches test notification via /notifications/test", async () => {
    vi.mocked(apiClient.post).mockResolvedValueOnce({
      data: {
        success: true,
        data: {
          message: "Notificação push de teste enviada com sucesso!",
          devicesNotified: 1,
          maintenanceAlerts: 0,
        },
      },
    });

    const result = await pushNotificationService.sendTestNotification();

    expect(apiClient.post).toHaveBeenCalledWith("/notifications/test");
    expect(result.devicesNotified).toBe(1);
  });
});
