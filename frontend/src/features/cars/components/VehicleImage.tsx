import React, { useState } from "react";
import { Box, SxProps, Theme, Skeleton } from "@mui/material";
import { CarSedanVector } from "./vectors/CarSedanVector";
import { MotorcycleVector } from "./vectors/MotorcycleVector";
import { TruckVector } from "./vectors/TruckVector";

export interface VehicleImageProps {
  imageUrl?: string | null;
  vehicleType?: "cars" | "motorcycles" | "trucks" | string;
  alt: string;
  className?: string;
  sx?: SxProps<Theme>;
  aspectRatio?: string | number;
  height?: number | string;
  width?: number | string;
  fit?: "cover" | "contain";
  borderRadius?: number | string;
}

export function VehicleImage({
  imageUrl,
  vehicleType = "cars",
  alt,
  className,
  sx,
  aspectRatio = "16 / 9",
  height,
  width = "100%",
  fit = "cover",
  borderRadius = 2,
}: VehicleImageProps) {
  const [imgError, setImgError] = useState(false);
  const [imgLoaded, setImgLoaded] = useState(false);

  const normalizedType = (vehicleType || "cars").toLowerCase();

  const renderFallbackVector = () => {
    switch (normalizedType) {
      case "motorcycles":
      case "motos":
      case "moto":
        return <MotorcycleVector aria-label={alt} />;
      case "trucks":
      case "caminhoes":
      case "caminhao":
      case "picape":
        return <TruckVector aria-label={alt} />;
      case "cars":
      case "carros":
      case "carro":
      default:
        return <CarSedanVector aria-label={alt} />;
    }
  };

  const hasCustomImage = Boolean(imageUrl && !imgError);

  return (
    <Box
      className={className}
      sx={{
        width,
        height,
        aspectRatio,
        position: "relative",
        overflow: "hidden",
        borderRadius,
        bgcolor: "#0F172A",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        border: "1px solid rgba(226, 232, 240, 0.8)",
        boxSizing: "border-box",
        ...sx,
      }}
    >
      {/* Background vector fallback when no custom image or when custom image errored */}
      {!hasCustomImage && (
        <Box
          sx={{
            width: "100%",
            height: "100%",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
          }}
        >
          {renderFallbackVector()}
        </Box>
      )}

      {/* Real image if provided */}
      {imageUrl && !imgError && (
        <>
          {!imgLoaded && (
            <Skeleton
              variant="rectangular"
              width="100%"
              height="100%"
              animation="wave"
              sx={{
                position: "absolute",
                top: 0,
                left: 0,
                bgcolor: "#1E293B",
              }}
            />
          )}

          <Box
            component="img"
            src={imageUrl}
            alt={alt}
            loading="lazy"
            onLoad={() => setImgLoaded(true)}
            onError={() => setImgError(true)}
            sx={{
              width: "100%",
              height: "100%",
              objectFit: fit,
              objectPosition: "center",
              display: "block",
              opacity: imgLoaded ? 1 : 0,
              transition: "opacity 0.3s ease-in-out",
            }}
          />
        </>
      )}
    </Box>
  );
}
