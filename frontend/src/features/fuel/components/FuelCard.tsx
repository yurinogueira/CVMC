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
import LocalGasStationRoundedIcon from "@mui/icons-material/LocalGasStationRounded";
import CalendarMonthRoundedIcon from "@mui/icons-material/CalendarMonthRounded";
import AttachMoneyRoundedIcon from "@mui/icons-material/AttachMoneyRounded";
import PlaceRoundedIcon from "@mui/icons-material/PlaceRounded";
import EditRoundedIcon from "@mui/icons-material/EditRounded";
import DeleteOutlineRoundedIcon from "@mui/icons-material/DeleteOutlineRounded";
import CheckCircleRoundedIcon from "@mui/icons-material/CheckCircleRounded";
import OpacityRoundedIcon from "@mui/icons-material/OpacityRounded";

import { Fueling } from "../types/fuel.types";
import { brandColors } from "../../../styles/theme";

interface FuelCardProps {
  fueling: Fueling;
  onEdit?: (fueling: Fueling) => void;
  onDelete?: (fuelingId: string) => void;
}

export function FuelCard({ fueling, onEdit, onDelete }: FuelCardProps) {
  const formatDate = (dateStr: string): string => {
    try {
      const datePart = dateStr.split("T")[0];
      const parts = datePart.split("-");
      if (parts.length === 3) {
        return `${parts[2]}/${parts[1]}/${parts[0]}`;
      }
      return new Date(dateStr).toLocaleDateString("pt-BR");
    } catch {
      return dateStr;
    }
  };

  const formatCost = (val: number): string => {
    return new Intl.NumberFormat("pt-BR", {
      style: "currency",
      currency: "BRL",
    }).format(val);
  };

  const formatPricePerLiter = (val: number): string => {
    return new Intl.NumberFormat("pt-BR", {
      style: "currency",
      currency: "BRL",
      minimumFractionDigits: 2,
      maximumFractionDigits: 3,
    }).format(val);
  };

  return (
    <Card
      elevation={0}
      sx={{
        border: "1px solid #E2E8F0",
        borderRadius: 2.5,
        bgcolor: "background.paper",
        transition: "all 0.2s ease",
        "&:hover": {
          borderColor: "#CBD5E1",
          boxShadow: "0 4px 16px -2px rgba(2, 132, 199, 0.08)",
        },
      }}
    >
      <CardContent sx={{ p: { xs: 2, sm: 2.5 } }}>
        <Stack spacing={2}>
          {/* Top Line: Header with Fuel Type & Action Buttons */}
          <Stack
            direction="row"
            sx={{
              alignItems: "flex-start",
              justifyContent: "space-between",
              gap: 1.5,
            }}
          >
            <Stack direction="row" spacing={1.5} sx={{ alignItems: "center" }}>
              <Box
                sx={{
                  width: 42,
                  height: 42,
                  borderRadius: 2,
                  bgcolor: "rgba(2, 132, 199, 0.1)",
                  color: "primary.main",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  flexShrink: 0,
                }}
              >
                <LocalGasStationRoundedIcon sx={{ fontSize: 24 }} />
              </Box>
              <Box>
                <Stack
                  direction="row"
                  spacing={1}
                  sx={{ alignItems: "center", flexWrap: "wrap", gap: 0.5 }}
                >
                  <Typography
                    variant="subtitle1"
                    sx={{ fontWeight: 800, color: "text.primary" }}
                  >
                    {fueling.fuelType}
                  </Typography>
                  {fueling.isFullTank && (
                    <Chip
                      icon={
                        <CheckCircleRoundedIcon
                          sx={{ "&&": { fontSize: 14 } }}
                        />
                      }
                      label="Tanque Cheio"
                      size="small"
                      sx={{
                        bgcolor: brandColors.mint,
                        color: "#064E3B",
                        fontWeight: 700,
                        fontSize: "0.75rem",
                        height: 22,
                        borderRadius: 1,
                      }}
                    />
                  )}
                </Stack>

                {fueling.gasStation && (
                  <Stack
                    direction="row"
                    spacing={0.5}
                    sx={{ alignItems: "center", mt: 0.25 }}
                  >
                    <PlaceRoundedIcon
                      sx={{ fontSize: 14, color: "text.secondary" }}
                    />
                    <Typography
                      variant="caption"
                      sx={{ color: "text.secondary", fontWeight: 500 }}
                    >
                      {fueling.gasStation}
                    </Typography>
                  </Stack>
                )}
              </Box>
            </Stack>

            {/* Action Buttons: Edit / Delete */}
            <Stack direction="row" spacing={0.5} sx={{ flexShrink: 0 }}>
              {onEdit && (
                <Tooltip title="Editar abastecimento" arrow>
                  <IconButton
                    size="small"
                    aria-label="Editar abastecimento"
                    onClick={() => onEdit(fueling)}
                    sx={{
                      color: "text.secondary",
                      "&:hover": {
                        color: "primary.main",
                        bgcolor: "rgba(2, 132, 199, 0.08)",
                      },
                    }}
                  >
                    <EditRoundedIcon fontSize="small" />
                  </IconButton>
                </Tooltip>
              )}
              {onDelete && (
                <Tooltip title="Excluir abastecimento" arrow>
                  <IconButton
                    size="small"
                    aria-label="Excluir abastecimento"
                    onClick={() => onDelete(fueling.id)}
                    sx={{
                      color: "text.secondary",
                      "&:hover": {
                        color: "error.main",
                        bgcolor: "rgba(239, 68, 68, 0.08)",
                      },
                    }}
                  >
                    <DeleteOutlineRoundedIcon fontSize="small" />
                  </IconButton>
                </Tooltip>
              )}
            </Stack>
          </Stack>

          {/* Notes description if any */}
          {fueling.notes && (
            <Typography
              variant="body2"
              sx={{ color: "text.secondary", lineHeight: 1.5 }}
            >
              {fueling.notes}
            </Typography>
          )}

          {/* Badges / Metrics Row: Date, Liters, Price/L, Total Cost */}
          <Stack
            direction="row"
            spacing={1.5}
            sx={{ flexWrap: "wrap", gap: 1.5, alignItems: "center" }}
          >
            <Chip
              icon={
                <CalendarMonthRoundedIcon sx={{ "&&": { fontSize: 16 } }} />
              }
              label={formatDate(fueling.date)}
              size="small"
              sx={{
                bgcolor: "#F1F5F9",
                color: "text.secondary",
                fontWeight: 600,
                borderRadius: 1.5,
              }}
            />

            <Chip
              icon={
                <OpacityRoundedIcon
                  sx={{ "&&": { fontSize: 16, color: "primary.dark" } }}
                />
              }
              label={`${fueling.liters.toLocaleString("pt-BR", { minimumFractionDigits: 1, maximumFractionDigits: 2 })} Litros`}
              size="small"
              sx={{
                bgcolor: "#E0F2FE",
                color: "primary.dark",
                fontWeight: 700,
                borderRadius: 1.5,
              }}
            />

            <Chip
              label={`${formatPricePerLiter(fueling.pricePerLiter)} / L`}
              size="small"
              sx={{
                bgcolor: "#F8FAFC",
                border: "1px solid #E2E8F0",
                color: "text.primary",
                fontWeight: 600,
                borderRadius: 1.5,
              }}
            />

            <Chip
              icon={
                <AttachMoneyRoundedIcon
                  sx={{ "&&": { fontSize: 16, color: "#15803D" } }}
                />
              }
              label={formatCost(fueling.totalCost)}
              size="small"
              sx={{
                bgcolor: "#DCFCE7",
                color: "#15803D",
                fontWeight: 800,
                borderRadius: 1.5,
              }}
            />
          </Stack>
        </Stack>
      </CardContent>
    </Card>
  );
}
