import React from "react";
import { Box, BoxProps } from "@mui/material";
import carImg from "../../../../assets/vehicles/car.webp";

interface VectorProps extends BoxProps {
  alt?: string;
}

export function CarSedanVector({ alt = "Carro Sedã", ...props }: VectorProps) {
  return (
    <Box
      sx={{
        width: "100%",
        height: "100%",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        p: { xs: 1.5, sm: 2 },
        boxSizing: "border-box",
        background:
          "radial-gradient(ellipse at center, rgba(56, 189, 248, 0.16) 0%, rgba(15, 23, 42, 0.95) 75%)",
        ...props.sx,
      }}
      {...props}
    >
      <Box
        component="img"
        src={carImg}
        alt={alt}
        loading="lazy"
        sx={{
          maxWidth: "94%",
          maxHeight: "90%",
          objectFit: "contain",
          filter: "drop-shadow(0 12px 24px rgba(0,0,0,0.55))",
          transition: "transform 0.3s ease",
          "&:hover": {
            transform: "scale(1.03)",
          },
        }}
      />
    </Box>
  );
}
