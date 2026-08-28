import React from "react";
import { Box, Typography, Stack, SxProps, Theme } from "@mui/material";
import { brandColors } from "../../styles/theme";

export interface BrandLogoProps {
  variant?: "full" | "compact" | "icon";
  size?: "small" | "medium" | "large";
  colorMode?: "light" | "dark" | "gradient";
  showSubtitle?: boolean;
  sx?: SxProps<Theme>;
}

export function BrandLogoIcon({
  size = 40,
  glow = true,
}: {
  size?: number;
  glow?: boolean;
}) {
  return (
    <Box
      sx={{
        width: size,
        height: size,
        borderRadius: size > 32 ? 2.5 : 1.5,
        background: brandColors.gradient,
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        boxShadow: glow ? "0 4px 14px rgba(2, 132, 199, 0.35)" : "none",
        flexShrink: 0,
        position: "relative",
        overflow: "hidden",
      }}
    >
      <svg
        viewBox="0 0 48 48"
        width={size * 0.72}
        height={size * 0.72}
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        <defs>
          <linearGradient
            id="logoInnerGrad"
            x1="0%"
            y1="0%"
            x2="100%"
            y2="100%"
          >
            <stop offset="0%" stopColor="#FFFFFF" />
            <stop offset="100%" stopColor="#E0F2FE" />
          </linearGradient>
          <linearGradient id="speedArcGrad" x1="0%" y1="0%" x2="100%" y2="0%">
            <stop offset="0%" stopColor="#38BDF8" />
            <stop offset="100%" stopColor="#34D399" />
          </linearGradient>
        </defs>

        {/* Speed Arc */}
        <path
          d="M10 32 C10 18, 20 10, 36 12 C41 13, 44 16, 44 20"
          stroke="url(#speedArcGrad)"
          strokeWidth="3.5"
          strokeLinecap="round"
        />

        {/* Aerodynamic Car Body 3/4 */}
        <path
          d="M8 35 C9 32, 12 30, 16 29 L21 24 C23 22, 28 22, 33 24 L38 29 C40 30, 42 32, 42 35 L8 35 Z"
          fill="url(#logoInnerGrad)"
        />

        {/* Windshield */}
        <path
          d="M17 29 L21 25 C23 23.5, 27 23.5, 31 25 L35 29 Z"
          fill="#0284C7"
          opacity="0.9"
        />

        {/* Front Wheel */}
        <circle
          cx="34"
          cy="35"
          r="3.5"
          fill="#0F172A"
          stroke="#38BDF8"
          strokeWidth="1"
        />
        <circle cx="34" cy="35" r="1.5" fill="#FFFFFF" />

        {/* Rear Wheel */}
        <circle
          cx="15"
          cy="35"
          r="3.5"
          fill="#0F172A"
          stroke="#38BDF8"
          strokeWidth="1"
        />
        <circle cx="15" cy="35" r="1.5" fill="#FFFFFF" />
      </svg>
    </Box>
  );
}

export function BrandLogo({
  variant = "full",
  size = "medium",
  colorMode = "gradient",
  showSubtitle = true,
  sx,
}: BrandLogoProps) {
  const iconSize = size === "small" ? 32 : size === "large" ? 48 : 40;
  const titleSize =
    size === "small" ? "1.1rem" : size === "large" ? "1.6rem" : "1.35rem";
  const subtitleSize =
    size === "small" ? "0.68rem" : size === "large" ? "0.85rem" : "0.75rem";

  const textColor =
    colorMode === "light"
      ? "#FFFFFF"
      : colorMode === "dark"
        ? "#0F172A"
        : "text.primary";

  const subtitleColor =
    colorMode === "light" ? "rgba(255, 255, 255, 0.85)" : "text.secondary";

  if (variant === "icon") {
    return (
      <Box sx={{ display: "inline-flex", ...sx }}>
        <BrandLogoIcon size={iconSize} />
      </Box>
    );
  }

  return (
    <Stack
      direction="row"
      spacing={1.5}
      sx={{ alignItems: "center", textDecoration: "none", ...sx }}
    >
      <BrandLogoIcon size={iconSize} />

      <Box>
        <Typography
          variant="h6"
          component="span"
          sx={{
            fontWeight: 900,
            fontSize: titleSize,
            lineHeight: 1.1,
            letterSpacing: -0.5,
            color: textColor,
            display: "block",
          }}
        >
          CVMC
        </Typography>

        {variant === "full" && showSubtitle && (
          <Typography
            variant="caption"
            sx={{
              color: subtitleColor,
              fontWeight: 600,
              fontSize: subtitleSize,
              letterSpacing: 0.2,
              display: "block",
              lineHeight: 1.1,
            }}
          >
            Como Vai Meu Carro
          </Typography>
        )}
      </Box>
    </Stack>
  );
}
