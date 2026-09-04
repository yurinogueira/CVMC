import {
  Card,
  CardContent,
  Typography,
  Box,
  Stack,
  Chip,
  IconButton,
  Tooltip,
} from "@mui/material";
import SpeedRoundedIcon from "@mui/icons-material/SpeedRounded";
import CalendarMonthRoundedIcon from "@mui/icons-material/CalendarMonthRounded";
import EditRoundedIcon from "@mui/icons-material/EditRounded";
import DeleteOutlineRoundedIcon from "@mui/icons-material/DeleteOutlineRounded";
import { Car } from "../types/car.types";
import { brandColors } from "../../../styles/theme";
import { VehicleImage } from "./VehicleImage";

interface VehicleCardProps {
  car: Car;
  onEdit?: (car: Car) => void;
  onDelete?: (id: string) => void;
  onClick?: (car: Car) => void;
}

export function VehicleCard({
  car,
  onEdit,
  onDelete,
  onClick,
}: VehicleCardProps) {
  const isClickable = Boolean(onClick);

  return (
    <Card
      elevation={0}
      role={isClickable ? "button" : undefined}
      tabIndex={isClickable ? 0 : undefined}
      onClick={() => {
        if (onClick) onClick(car);
      }}
      onKeyDown={(e) => {
        if (onClick && (e.key === "Enter" || e.key === " ")) {
          e.preventDefault();
          onClick(car);
        }
      }}
      sx={{
        height: "100%",
        display: "flex",
        flexDirection: "column",
        justifyContent: "space-between",
        borderRadius: 2.5,
        border: "1px solid #E2E8F0",
        overflow: "hidden",
        bgcolor: "background.paper",
        cursor: isClickable ? "pointer" : "default",
        transition: "transform 0.2s ease, box-shadow 0.2s ease",
        "&:hover": {
          transform: "translateY(-4px)",
          boxShadow: "0 12px 28px -4px rgba(2, 132, 199, 0.12)",
        },
        "&:focus-visible": isClickable
          ? {
              outline: "2px solid #0284C7",
              outlineOffset: "2px",
            }
          : undefined,
      }}
    >
      {/* Vehicle Media Header */}
      <Box sx={{ position: "relative", width: "100%" }}>
        <VehicleImage
          imageUrl={car.imageUrl}
          vehicleType={car.vehicleType}
          alt={`${car.manufacturer} ${car.model}`}
          aspectRatio="16 / 9"
          borderRadius={0}
        />

        {/* Floating Action Buttons */}
        <Box
          onClick={(e) => e.stopPropagation()}
          onKeyDown={(e) => e.stopPropagation()}
          sx={{
            position: "absolute",
            top: 10,
            right: 10,
            display: "flex",
            gap: 0.5,
            bgcolor: "rgba(15, 23, 42, 0.65)",
            backdropFilter: "blur(6px)",
            borderRadius: 2,
            p: 0.5,
          }}
        >
          {onEdit && (
            <Tooltip title="Editar veículo">
              <IconButton
                size="small"
                onClick={(e) => {
                  e.stopPropagation();
                  onEdit(car);
                }}
                sx={{
                  color: "#FFFFFF",
                  "&:hover": {
                    color: "#38BDF8",
                    bgcolor: "rgba(255,255,255,0.1)",
                  },
                }}
                aria-label="Editar veículo"
              >
                <EditRoundedIcon fontSize="small" />
              </IconButton>
            </Tooltip>
          )}

          {onDelete && (
            <Tooltip title="Remover veículo">
              <IconButton
                size="small"
                onClick={(e) => {
                  e.stopPropagation();
                  onDelete(car.id);
                }}
                sx={{
                  color: "#FFFFFF",
                  "&:hover": {
                    color: "#F87171",
                    bgcolor: "rgba(255,255,255,0.1)",
                  },
                }}
                aria-label="Remover veículo"
              >
                <DeleteOutlineRoundedIcon fontSize="small" />
              </IconButton>
            </Tooltip>
          )}
        </Box>
      </Box>

      <CardContent
        sx={{
          p: 2.5,
          flex: 1,
          display: "flex",
          flexDirection: "column",
          justifyContent: "space-between",
        }}
      >
        {/* Card Titles */}
        <Box sx={{ mb: 2 }}>
          <Typography
            variant="subtitle1"
            component="h3"
            sx={{
              fontWeight: 800,
              lineHeight: 1.2,
              mb: 0.5,
              color: "text.primary",
            }}
          >
            {car.name}
          </Typography>
          <Typography
            variant="caption"
            sx={{
              color: "text.secondary",
              fontWeight: 600,
              fontSize: "0.85rem",
              display: "block",
            }}
          >
            {car.manufacturer} • {car.model}
          </Typography>
        </Box>

        {/* Badges / Details */}
        <Stack direction="row" spacing={1} sx={{ flexWrap: "wrap", gap: 1 }}>
          <Chip
            icon={<CalendarMonthRoundedIcon sx={{ "&&": { fontSize: 16 } }} />}
            label={`${car.yearManufacture}/${car.yearModel}`}
            size="small"
            sx={{
              bgcolor: "#F1F5F9",
              color: "text.primary",
              fontWeight: 500,
              borderRadius: 1,
            }}
          />
          <Chip
            icon={
              <SpeedRoundedIcon
                sx={{ "&&": { fontSize: 16, color: "#064E3B" } }}
              />
            }
            label={`${(car.lastMileage || 0).toLocaleString("pt-BR")} km`}
            size="small"
            sx={{
              bgcolor: brandColors.mint,
              color: "#064E3B",
              fontWeight: 600,
              borderRadius: 1,
            }}
          />
          {car.fipePrice && (
            <Chip
              label={`FIPE: ${car.fipePrice}`}
              size="small"
              sx={{
                bgcolor: "rgba(76, 146, 252, 0.1)",
                color: "primary.main",
                fontWeight: 600,
                borderRadius: 1,
              }}
            />
          )}
        </Stack>
      </CardContent>
    </Card>
  );
}
