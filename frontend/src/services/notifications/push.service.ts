import { Capacitor, PluginListenerHandle } from "@capacitor/core";
import {
  PushNotifications,
  PushNotificationSchema,
  ActionPerformed,
  Token,
} from "@capacitor/push-notifications";
import { apiClient } from "../api/client";
import { safeStorage } from "../storage/storage";

const FCM_TOKEN_STORAGE_KEY = "cvmc.fcm_token";
export const DEFAULT_ANDROID_CHANNEL_ID = "cvmc_alerts";

export interface ForegroundNotificationPayload {
  title: string;
  body: string;
  route?: string;
}

export interface PushInitCallbacks {
  onForegroundNotification?: (payload: ForegroundNotificationPayload) => void;
  onNavigate?: (route: string) => void;
}

export interface TestNotificationResponse {
  message: string;
  devicesNotified: number;
  maintenanceAlerts: number;
}

interface ApiEnvelope<T> {
  success: boolean;
  data: T;
  message?: string;
}

export function resolveNotificationDeepLink(
  data?: Record<string, unknown>,
): string | undefined {
  if (!data) {
    return undefined;
  }
  const rawRoute = typeof data.route === "string" ? data.route.trim() : "";
  const carId = typeof data.carId === "string" ? data.carId.trim() : "";
  const validCarId = carId && /^[a-zA-Z0-9_-]+$/.test(carId) ? carId : "";

  if (rawRoute.startsWith("/") && !rawRoute.startsWith("//")) {
    if (rawRoute === "/maintenance" && validCarId) {
      return `/maintenance?carId=${validCarId}`;
    }
    return rawRoute;
  }
  if (validCarId) {
    return `/vehicles/${validCarId}`;
  }
  return undefined;
}

let activeListeners: PluginListenerHandle[] = [];
let initialized = false;

export const pushNotificationService = {
  getStoredToken(): string | null {
    return safeStorage.getItem(FCM_TOKEN_STORAGE_KEY);
  },

  async registerTokenOnBackend(
    token: string,
    platform: string = Capacitor.getPlatform(),
  ): Promise<void> {
    const cleanToken = token.trim();
    if (!cleanToken) {
      return;
    }
    await apiClient.post("/devices/token", {
      token: cleanToken,
      platform: platform || "android",
    });
    safeStorage.setItem(FCM_TOKEN_STORAGE_KEY, cleanToken);
  },

  async unregisterTokenFromBackend(token?: string | null): Promise<void> {
    const targetToken = (token ?? this.getStoredToken() ?? "").trim();
    safeStorage.removeItem(FCM_TOKEN_STORAGE_KEY);
    if (!targetToken) {
      return;
    }
    try {
      await apiClient.delete("/devices/token", {
        data: { token: targetToken },
      });
    } catch {
      // Ignore network errors during logout cleanup
    }
  },

  async sendTestNotification(): Promise<TestNotificationResponse> {
    const response = await apiClient.post<
      ApiEnvelope<TestNotificationResponse>
    >("/notifications/test");
    return response.data.data;
  },

  async init(callbacks: PushInitCallbacks = {}): Promise<void> {
    if (!Capacitor.isNativePlatform() || initialized) {
      return;
    }
    initialized = true;

    try {
      // 1. Configure Android notification channel with sound and vibration
      if (Capacitor.getPlatform() === "android") {
        await PushNotifications.createChannel({
          id: DEFAULT_ANDROID_CHANNEL_ID,
          name: "Alertas e Manutenções CVMC",
          description:
            "Avisos de revisões preventivas, quilometragem e alertas da sua frota",
          importance: 5,
          visibility: 1,
          sound: "default",
          vibration: true,
        }).catch(() => {
          // Ignore channel creation error on unsupported versions
        });
      }

      // 2. Check and request runtime permission (Android 13+ POST_NOTIFICATIONS)
      let permStatus = await PushNotifications.checkPermissions();
      if (permStatus.receive === "prompt") {
        permStatus = await PushNotifications.requestPermissions();
      }
      if (permStatus.receive !== "granted") {
        initialized = false;
        return;
      }

      // 3. Register listeners
      const regHandle = await PushNotifications.addListener(
        "registration",
        (token: Token) => {
          void this.registerTokenOnBackend(
            token.value,
            Capacitor.getPlatform(),
          ).catch(() => {
            // Silently ignore registration sync error
          });
        },
      );

      const fgHandle = await PushNotifications.addListener(
        "pushNotificationReceived",
        (notification: PushNotificationSchema) => {
          const data = (notification.data ?? {}) as Record<string, unknown>;
          callbacks.onForegroundNotification?.({
            title: notification.title || "Notificação CVMC",
            body: notification.body || "",
            route: resolveNotificationDeepLink(data),
          });
        },
      );

      const actionHandle = await PushNotifications.addListener(
        "pushNotificationActionPerformed",
        (action: ActionPerformed) => {
          const data = (action.notification?.data ?? {}) as Record<
            string,
            unknown
          >;
          const targetRoute = resolveNotificationDeepLink(data);
          if (targetRoute) {
            callbacks.onNavigate?.(targetRoute);
          }
        },
      );

      activeListeners = [regHandle, fgHandle, actionHandle];

      // 4. Trigger FCM registration
      await PushNotifications.register();
    } catch {
      initialized = false;
    }
  },

  async cleanup(): Promise<void> {
    for (const handle of activeListeners) {
      await handle.remove().catch(() => {});
    }
    activeListeners = [];
    initialized = false;
  },
};
