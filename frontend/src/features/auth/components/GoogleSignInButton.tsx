import React, { useEffect, useRef, useState } from "react";
import { Box, Button, CircularProgress } from "@mui/material";
import { Capacitor } from "@capacitor/core";
import { SocialLogin } from "@capgo/capacitor-social-login";

export function GoogleIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg
      width="18"
      height="18"
      viewBox="0 0 18 18"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <path
        d="M17.64 9.2c0-.637-.057-1.251-.164-1.84H9v3.481h4.844c-.209 1.125-.843 2.078-1.796 2.717v2.258h2.908c1.702-1.567 2.684-3.874 2.684-6.616z"
        fill="#4285F4"
      />
      <path
        d="M9 18c2.43 0 4.467-.806 5.956-2.184l-2.908-2.258c-.806.54-1.837.86-3.048.86-2.344 0-4.328-1.584-5.036-3.711H.957v2.332C2.438 15.983 5.482 18 9 18z"
        fill="#34A853"
      />
      <path
        d="M3.964 10.707c-.18-.54-.282-1.117-.282-1.707s.102-1.167.282-1.707V4.961H.957C.347 6.175 0 7.55 0 9s.348 2.825.957 4.039l3.007-2.332z"
        fill="#FBBC05"
      />
      <path
        d="M9 3.58c1.321 0 2.508.454 3.44 1.345l2.582-2.58C13.463.891 11.426 0 9 0 5.482 0 2.438 2.017.957 4.961L3.964 7.293C4.672 5.166 6.656 3.58 9 3.58z"
        fill="#EA4335"
      />
    </svg>
  );
}

declare global {
  interface Window {
    google?: {
      accounts: {
        id: {
          initialize: (config: {
            client_id: string;
            callback: (response: { credential: string }) => void;
            auto_select?: boolean;
          }) => void;
          renderButton: (
            parent: HTMLElement,
            options: {
              theme?: string;
              size?: string;
              width?: string | number;
              text?: string;
              shape?: string;
              logo_alignment?: string;
              locale?: string;
            },
          ) => void;
          prompt: () => void;
        };
      };
    };
  }
}

interface GoogleSignInButtonProps {
  onSuccess: (idToken: string) => void;
  onError?: (error: Error) => void;
  text?: "signin_with" | "signup_with" | "continue_with";
  disabled?: boolean;
  label?: string;
  width?: number;
}

export function GoogleSignInButton({
  onSuccess,
  onError,
  text = "signin_with",
  disabled = false,
  label,
  width,
}: GoogleSignInButtonProps) {
  const buttonRef = useRef<HTMLDivElement>(null);
  const [isGsiRendered, setIsGsiRendered] = useState(false);
  const [loadingNative, setLoadingNative] = useState(false);
  const isNative = Capacitor.isNativePlatform();

  const clientId =
    import.meta.env.VITE_GOOGLE_CLIENT_ID ||
    "cvmc-google-client-id-placeholder";

  useEffect(() => {
    // On native mobile platforms (Android/iOS), use native SocialLogin rather than Web GIS
    if (isNative) {
      return;
    }

    // If window.google already exists, render the button
    const initializeGsi = () => {
      if (!window.google?.accounts?.id || !buttonRef.current) return;
      try {
        window.google.accounts.id.initialize({
          client_id: clientId,
          callback: (response) => {
            if (response.credential) {
              onSuccess(response.credential);
            }
          },
        });

        // Clamped width between Google GIS limits (200px to 400px)
        const effectiveWidth = Math.min(380, Math.max(200, width || 280));

        window.google.accounts.id.renderButton(buttonRef.current, {
          theme: "outline",
          size: "large",
          width: effectiveWidth,
          text,
          shape: "rectangular",
          logo_alignment: "center",
          locale: "pt-BR",
        });
        setIsGsiRendered(true);
      } catch {
        // Fallback to custom button if rendering fails
      }
    };

    if (window.google?.accounts?.id) {
      initializeGsi();
      return;
    }

    // Load GIS script dynamically on web
    const scriptId = "google-gsi-client";
    let script = document.getElementById(scriptId) as HTMLScriptElement | null;
    if (!script) {
      script = document.createElement("script");
      script.id = scriptId;
      script.src = "https://accounts.google.com/gsi/client";
      script.async = true;
      script.defer = true;
      script.onload = initializeGsi;
      document.body.appendChild(script);
    } else {
      script.addEventListener("load", initializeGsi);
    }
  }, [clientId, text, onSuccess, isNative, width]);

  const buttonText =
    label ||
    (text === "signup_with" ? "Cadastrar com Google" : "Entrar com Google");

  const handleNativeLogin = async () => {
    if (disabled || loadingNative) return;
    try {
      setLoadingNative(true);
      await SocialLogin.initialize({
        google: {
          webClientId: clientId,
        },
      });

      const res = await SocialLogin.login({
        provider: "google",
        options: {},
      });

      const result = res?.result as { idToken?: string | null } | undefined;
      if (result?.idToken) {
        onSuccess(result.idToken);
      } else {
        onError?.(new Error("Token do Google não retornado pelo dispositivo."));
      }
    } catch (err: unknown) {
      const error = err as { code?: string; message?: string };
      // User cancelled account selection; dismiss cleanly
      if (
        error?.code === "USER_CANCELLED" ||
        error?.message?.toLowerCase().includes("cancel")
      ) {
        return;
      }
      onError?.(
        err instanceof Error
          ? err
          : new Error(
              error?.message || "Erro na autenticação nativa com o Google.",
            ),
      );
    } finally {
      setLoadingNative(false);
    }
  };

  const handleFallbackClick = () => {
    if (isNative) {
      handleNativeLogin();
    } else if (window.google?.accounts?.id) {
      window.google.accounts.id.prompt();
    }
  };

  return (
    <Box
      sx={{
        width: "100%",
        maxWidth: "100%",
        display: "flex",
        justifyContent: "center",
        overflow: "hidden",
      }}
    >
      {/* Official Google Web Button container (only rendered on web) */}
      {!isNative && (
        <Box
          ref={buttonRef}
          sx={{
            width: "100%",
            maxWidth: "100%",
            display: isGsiRendered ? "flex" : "none",
            justifyContent: "center",
            overflow: "hidden",
            "& iframe": {
              margin: "0 auto",
              maxWidth: "100% !important",
            },
          }}
        />
      )}

      {/* Styled button for Native Mobile platforms and Web fallback */}
      {(isNative || !isGsiRendered) && (
        <Button
          fullWidth
          variant="outlined"
          startIcon={
            loadingNative ? (
              <CircularProgress size={18} color="inherit" />
            ) : (
              <GoogleIcon />
            )
          }
          disabled={disabled || loadingNative}
          onClick={handleFallbackClick}
          data-testid="google-signin-btn"
          sx={{
            color: "text.primary",
            borderColor: "#E2E8F0",
            backgroundColor: "#FFFFFF",
            textTransform: "none",
            fontWeight: 600,
            py: 1.1,
            maxWidth: "100%",
            "&:hover": {
              borderColor: "#CBD5E1",
              backgroundColor: "#F8FAFC",
            },
          }}
        >
          {loadingNative ? "Conectando ao Google..." : buttonText}
        </Button>
      )}
    </Box>
  );
}
