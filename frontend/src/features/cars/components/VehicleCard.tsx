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
import DeleteOutlineRoundedIcon from "@mui/icons-material/DeleteOutlineRounded";
import { Car } from "../types/car.types";
import { brandColors } from "../../../styles/theme";

interface VehicleCardProps {
  car: Car;
  onDelete?: (id: string) => void;
}

export function VehicleCard({ car, onDelete }: VehicleCardProps) {
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
          alignItems="flex-start"
          justifyContent="space-between"
          sx={{ mb: 2 }}
        >
          <Stack direction="row" spacing={1.5} alignItems="center">
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

          {onDelete && (
            <Tooltip title="Remover veículo">
              <IconButton
                size="small"
                onClick={() => onDelete(car.id)}
                sx={{
                  color: "text.secondary",
                  "&:hover": { color: "#EF4444" },
                }}
              >
                <DeleteOutlineRoundedIcon fontSize="small" />
              </IconButton>
            </Tooltip>
          )}
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
        </Stack>
      </CardContent>
    </Card>
  );
}
