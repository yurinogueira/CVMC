import { useEffect } from "react";
import { BrowserRouter, useNavigate, useLocation } from "react-router-dom";
import { Capacitor } from "@capacitor/core";
import { App as CapApp } from "@capacitor/app";
import { StatusBar, Style } from "@capacitor/status-bar";
import { SplashScreen } from "@capacitor/splash-screen";
import { AppRoutes } from "./routes";

function MobileAppBridge() {
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    if (!Capacitor.isNativePlatform()) {
      return;
    }

    // Configure Status Bar
    void (async () => {
      try {
        await StatusBar.setBackgroundColor({ color: "#0284C7" });
        await StatusBar.setStyle({ style: Style.Dark });
      } catch {
        // Plugin error fallback
      }
    })();

    // Hide Splash Screen after app mounted
    void (async () => {
      try {
        await SplashScreen.hide();
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

  return <AppRoutes />;
}

export function App() {
  return (
    <BrowserRouter>
      <MobileAppBridge />
    </BrowserRouter>
  );
}
