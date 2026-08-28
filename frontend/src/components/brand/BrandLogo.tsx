import React from "react";
import { Box, Typography, Stack, SxProps, Theme } from "@mui/material";

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
        borderRadius: "50%",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        boxShadow: glow ? "0 4px 14px rgba(2, 132, 199, 0.35)" : "none",
        flexShrink: 0,
        position: "relative",
        overflow: "hidden",
      }}
    >
      <Box
        component="img"
        src="/favicon.svg"
        alt="CVMC Logo"
        sx={{
          width: "100%",
          height: "100%",
          display: "block",
          objectFit: "contain",
        }}
      />
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
  const iconSize = size === "small" ? 34 : size === "large" ? 48 : 40;
  const titleSize =
    size === "small" ? "1.15rem" : size === "large" ? "1.6rem" : "1.35rem";
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
