import { useEffect, useState } from "react";
import { BrowserRouter, useNavigate, useLocation } from "react-router-dom";
import { Alert, Button, Snackbar } from "@mui/material";
import { Capacitor } from "@capacitor/core";
import { App as CapApp } from "@capacitor/app";
import { StatusBar, Style } from "@capacitor/status-bar";
import { SplashScreen } from "@capacitor/splash-screen";
import { AppRoutes } from "./routes";
import {
  useAuthStore,
  initAuthSession,
} from "./features/auth/state/auth.store";
import {
  pushNotificationService,
  ForegroundNotificationPayload,
} from "./services/notifications/push.service";

function MobileAppBridge() {
  const navigate = useNavigate();
  const location = useLocation();
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const isInitializing = useAuthStore((state) => state.isInitializing);
  const [foregroundPush, setForegroundPush] =
    useState<ForegroundNotificationPayload | null>(null);

  useEffect(() => {
    void initAuthSession();
  }, []);

  useEffect(() => {
    if (!Capacitor.isNativePlatform() || isInitializing) {
      return;
    }

    // Hide Splash Screen after session initialization is ready
    void (async () => {
      try {
        await SplashScreen.hide();
      } catch {
        // Plugin error fallback
      }
    })();
  }, [isInitializing]);

  useEffect(() => {
    if (!Capacitor.isNativePlatform()) {
      return;
    }

    // Configure Status Bar (Light background with dark icons for optimal contrast)
    void (async () => {
      try {
        await StatusBar.setBackgroundColor({ color: "#FFFFFF" });
        await StatusBar.setStyle({ style: Style.Light });
      } catch {
        // Plugin error fallback
      }
    })();

    // Hardware Back Button listener
    const backButtonHandler = CapApp.addListener(
      "backButton",
      ({ canGoBack }) => {
        if (
          location.pathname === "/dashboard" ||
          location.pathname === "/login" ||
          location.pathname === "/"
        ) {
          void CapApp.exitApp();
        } else if (canGoBack) {
          navigate(-1);
        } else {
          navigate("/dashboard", { replace: true });
        }
      },
    );

    return () => {
      void backButtonHandler.then((handle) => handle.remove());
    };
  }, [location.pathname, navigate]);

  useEffect(() => {
    if (!Capacitor.isNativePlatform() || !isAuthenticated || isInitializing) {
      return;
    }

    void pushNotificationService.init({
      onForegroundNotification: (payload) => {
        setForegroundPush(payload);
      },
      onNavigate: (targetRoute) => {
        navigate(targetRoute);
      },
    });
  }, [isAuthenticated, isInitializing, navigate]);

  return (
    <>
      <AppRoutes />
      <Snackbar
        open={Boolean(foregroundPush)}
        autoHideDuration={6000}
        onClose={() => setForegroundPush(null)}
        anchorOrigin={{ vertical: "top", horizontal: "center" }}
      >
        <Alert
          onClose={() => setForegroundPush(null)}
          severity="info"
          variant="filled"
          sx={{ width: "100%", alignItems: "center" }}
          action={
            foregroundPush?.route ? (
              <Button
                color="inherit"
                size="small"
                onClick={() => {
                  const route = foregroundPush.route;
                  setForegroundPush(null);
                  if (route) {
                    navigate(route);
                  }
                }}
              >
                Abrir
              </Button>
            ) : undefined
          }
        >
          <strong>{foregroundPush?.title}</strong>
          {foregroundPush?.body ? ` — ${foregroundPush.body}` : ""}
        </Alert>
      </Snackbar>
    </>
  );
}

export function App() {
  return (
    <BrowserRouter>
      <MobileAppBridge />
    </BrowserRouter>
  );
}
