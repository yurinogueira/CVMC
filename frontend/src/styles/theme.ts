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
  blue: "#4C92FC",
  sky: "#4CCAFC",
  cyan: "#4CFCF7",
  mint: "#4CFCBB",
  green: "#4CFC7F",
  gradient: "linear-gradient(135deg, #4C92FC 0%, #4CCAFC 50%, #4CFCF7 100%)",
};

export const appTheme = createTheme({
  palette: {
    mode: "light",
    primary: {
      main: brandColors.blue,
      light: "#78ACFD",
      dark: "#2563EB",
      contrastText: "#FFFFFF",
    },
    secondary: {
      main: brandColors.sky,
      light: "#8AE0FE",
      dark: "#0284C7",
      contrastText: "#0F172A",
    },
    success: {
      main: brandColors.green,
      light: "#85FDAE",
      dark: "#16A34A",
      contrastText: "#064E3B",
    },
    info: {
      main: brandColors.mint,
      light: "#8EFDD3",
      dark: "#0D9488",
      contrastText: "#064E3B",
    },
    background: {
      default: "#F8FAFC",
      paper: "#FFFFFF",
    },
    text: {
      primary: "#0F172A",
      secondary: "#64748B",
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
