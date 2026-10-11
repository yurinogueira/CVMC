import React, { useEffect, useRef } from "react";
import { Box, Button } from "@mui/material";

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
}

export function GoogleSignInButton({
  onSuccess,
  text = "signin_with",
  disabled = false,
  label,
}: GoogleSignInButtonProps) {
  const buttonRef = useRef<HTMLDivElement>(null);
  const [isGsiRendered, setIsGsiRendered] = React.useState(false);
  const clientId =
    import.meta.env.VITE_GOOGLE_CLIENT_ID ||
    "cvmc-google-client-id-placeholder";

  useEffect(() => {
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
        window.google.accounts.id.renderButton(buttonRef.current, {
          theme: "outline",
          size: "large",
          width: 400,
          text,
          shape: "rectangular",
          logo_alignment: "left",
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

    // Load GIS script dynamically
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
  }, [clientId, text, onSuccess]);

  const buttonText =
    label ||
    (text === "signup_with" ? "Cadastrar com Google" : "Entrar com Google");

  return (
    <Box
      sx={{
        width: "100%",
        display: "flex",
        justifyContent: "center",
      }}
    >
      {/* Official Google Button container */}
      <Box
        ref={buttonRef}
        sx={{
          width: "100%",
          display: "flex",
          justifyContent: "center",
          "& > div": {
            width: "100% !important",
            display: "flex !important",
            justifyContent: "center !important",
          },
          "& iframe": {
            margin: "0 auto !important",
            width: "100% !important",
            minWidth: "100% !important",
          },
        }}
      />
      {/* Fallback button for test environments or when GIS is not rendered */}
      <Button
        fullWidth
        variant="outlined"
        startIcon={<GoogleIcon />}
        disabled={disabled}
        onClick={() => {
          if (window.google?.accounts?.id) {
            window.google.accounts.id.prompt();
          }
        }}
        data-testid="google-signin-btn"
        sx={{
          display: isGsiRendered ? "none" : "flex",
          color: "text.primary",
          borderColor: "#E2E8F0",
          backgroundColor: "#FFFFFF",
          textTransform: "none",
          fontWeight: 600,
          py: 1.1,
          "&:hover": {
            borderColor: "#CBD5E1",
            backgroundColor: "#F8FAFC",
          },
        }}
      >
        {buttonText}
      </Button>
    </Box>
  );
}
