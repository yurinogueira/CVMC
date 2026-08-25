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
import DirectionsCarFilledRoundedIcon from "@mui/icons-material/DirectionsCarFilledRounded";
import SpeedRoundedIcon from "@mui/icons-material/SpeedRounded";
import CalendarMonthRoundedIcon from "@mui/icons-material/CalendarMonthRounded";
import EditRoundedIcon from "@mui/icons-material/EditRounded";
import DeleteOutlineRoundedIcon from "@mui/icons-material/DeleteOutlineRounded";
import { Car } from "../types/car.types";
import { brandColors } from "../../../styles/theme";

interface VehicleCardProps {
  car: Car;
  onEdit?: (car: Car) => void;
  onDelete?: (id: string) => void;
}

export function VehicleCard({ car, onEdit, onDelete }: VehicleCardProps) {
  return (
    <Card
      elevation={0}
      sx={{
        height: "100%",
        display: "flex",
        flexDirection: "column",
        justifyContent: "space-between",
        transition: "transform 0.2s ease, box-shadow 0.2s ease",
        "&:hover": {
          transform: "translateY(-3px)",
          boxShadow: "0 10px 25px -3px rgba(76, 146, 252, 0.12)",
        },
      }}
    >
      <CardContent sx={{ p: 2.5 }}>
        {/* Card Header */}
        <Stack
          direction="row"
          sx={{
            alignItems: "flex-start",
            justifyContent: "space-between",
            mb: 2,
          }}
        >
          <Stack direction="row" spacing={1.5} sx={{ alignItems: "center" }}>
            <Box
              sx={{
                width: 44,
                height: 44,
                borderRadius: 2,
                bgcolor: "rgba(76, 146, 252, 0.1)",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                color: "primary.main",
              }}
            >
              <DirectionsCarFilledRoundedIcon fontSize="medium" />
            </Box>
            <Box>
              <Typography
                variant="subtitle1"
                component="h3"
                sx={{ fontWeight: 700, lineHeight: 1.2 }}
              >
                {car.name}
              </Typography>
              <Typography
                variant="caption"
                sx={{ color: "text.secondary", fontWeight: 500 }}
              >
                {car.manufacturer} • {car.model}
              </Typography>
            </Box>
          </Stack>

          <Stack direction="row" spacing={0.5}>
            {onEdit && (
              <Tooltip title="Editar veículo">
                <IconButton
                  size="small"
                  onClick={() => onEdit(car)}
                  sx={{
                    color: "text.secondary",
                    "&:hover": { color: "primary.main" },
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
                  onClick={() => onDelete(car.id)}
                  sx={{
                    color: "text.secondary",
                    "&:hover": { color: "#EF4444" },
                  }}
                  aria-label="Remover veículo"
                >
                  <DeleteOutlineRoundedIcon fontSize="small" />
                </IconButton>
              </Tooltip>
            )}
          </Stack>
        </Stack>

        {/* Badges / Details */}
        <Stack
          direction="row"
          spacing={1}
          sx={{ my: 2, flexWrap: "wrap", gap: 1 }}
        >
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
