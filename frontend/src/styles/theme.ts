import { createTheme } from "@mui/material/styles";

declare module "@mui/material/styles" {
  interface Palette {
    brand: {
      blue: string;
      sky: string;
      cyan: string;
      mint: string;
      green: string;
      gradient: string;
    };
  }
  interface PaletteOptions {
    brand?: {
      blue: string;
      sky: string;
      cyan: string;
      mint: string;
      green: string;
      gradient: string;
    };
  }
}

export const brandColors = {
  blue: "#0284C7",
  sky: "#38BDF8",
  cyan: "#06B6D4",
  mint: "#A7F3D0",
  green: "#86EFAC",
  gradient: "linear-gradient(135deg, #0284C7 0%, #0369A1 50%, #0F766E 100%)",
};

export const appTheme = createTheme({
  palette: {
    mode: "light",
    primary: {
      main: "#0284C7",
      light: "#38BDF8",
      dark: "#0369A1",
      contrastText: "#FFFFFF",
    },
    secondary: {
      main: "#0EA5E9",
      light: "#7DD3FC",
      dark: "#0284C7",
      contrastText: "#0F172A",
    },
    success: {
      main: "#16A34A",
      light: "#86EFAC",
      dark: "#15803D",
      contrastText: "#FFFFFF",
    },
    info: {
      main: "#0D9488",
      light: "#99F6E4",
      dark: "#0F766E",
      contrastText: "#FFFFFF",
    },
    background: {
      default: "#F8FAFC",
      paper: "#FFFFFF",
    },
    text: {
      primary: "#0F172A",
      secondary: "#475569",
    },
    divider: "#E2E8F0",
    brand: brandColors,
  },
  shape: {
    borderRadius: 8,
  },
  typography: {
    fontFamily: ["Inter", "system-ui", "-apple-system", "sans-serif"].join(","),
    button: {
      textTransform: "none",
      fontWeight: 600,
    },
    h4: {
      fontWeight: 700,
      color: "#0F172A",
    },
    h5: {
      fontWeight: 600,
      color: "#0F172A",
    },
    h6: {
      fontWeight: 600,
      color: "#0F172A",
    },
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          boxShadow: "none",
          "&:hover": {
            boxShadow: "0 4px 12px rgba(76, 146, 252, 0.25)",
          },
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          borderRadius: 16,
          border: "1px solid #E2E8F0",
          boxShadow:
            "0 4px 20px -2px rgba(76, 146, 252, 0.06), 0 2px 6px -1px rgba(0, 0, 0, 0.03)",
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        rounded: {
          borderRadius: 16,
        },
      },
    },
  },
});
