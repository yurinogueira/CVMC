import {
  Card,
  CardContent,
  Typography,
  Box,
  Stack,
  Chip,
  Button,
  LinearProgress,
} from "@mui/material";
import CheckCircleRoundedIcon from "@mui/icons-material/CheckCircleRounded";
import WarningAmberRoundedIcon from "@mui/icons-material/WarningAmberRounded";
import ErrorOutlineRoundedIcon from "@mui/icons-material/ErrorOutlineRounded";
import OpacityRoundedIcon from "@mui/icons-material/OpacityRounded";
import FilterAltRoundedIcon from "@mui/icons-material/FilterAltRounded";
import AirRoundedIcon from "@mui/icons-material/AirRounded";
import AcUnitRoundedIcon from "@mui/icons-material/AcUnitRounded";
import LocalGasStationRoundedIcon from "@mui/icons-material/LocalGasStationRounded";
import AddCircleOutlineRoundedIcon from "@mui/icons-material/AddCircleOutlineRounded";
import SpeedRoundedIcon from "@mui/icons-material/SpeedRounded";
import CalendarMonthRoundedIcon from "@mui/icons-material/CalendarMonthRounded";
import {
  UpcomingMaintenanceItem,
  UpcomingMaintenanceStatus,
} from "../types/maintenance.types";

interface UpcomingMaintenanceCardProps {
  item: UpcomingMaintenanceItem;
  onQuickRegister?: (maintenanceType: string) => void;
}

function getStatusConfig(status: UpcomingMaintenanceStatus) {
  switch (status) {
    case "overdue":
      return {
        label: "Vencido",
        chipBg: "#FEE2E2",
        chipColor: "#991B1B",
        borderColor: "#FCA5A5",
        progressColor: "error" as const,
        icon: <ErrorOutlineRoundedIcon sx={{ fontSize: 16 }} />,
      };
    case "upcoming":
      return {
        label: "Próximo",
        chipBg: "#FEF3C7",
        chipColor: "#92400E",
        borderColor: "#FCD34D",
        progressColor: "warning" as const,
        icon: <WarningAmberRoundedIcon sx={{ fontSize: 16 }} />,
      };
    default:
      return {
        label: "Em dia",
        chipBg: "#DCFCE7",
        chipColor: "#166534",
        borderColor: "#E2E8F0",
        progressColor: "success" as const,
        icon: <CheckCircleRoundedIcon sx={{ fontSize: 16 }} />,
      };
  }
}

function getItemIcon(type: string) {
  switch (type) {
    case "Óleo de Motor":
      return <OpacityRoundedIcon sx={{ fontSize: 22 }} />;
    case "Filtro do Óleo de Motor":
      return <FilterAltRoundedIcon sx={{ fontSize: 22 }} />;
    case "Filtro de Ar do Motor":
      return <AirRoundedIcon sx={{ fontSize: 22 }} />;
    case "Filtro do Ar-Condicionado":
      return <AcUnitRoundedIcon sx={{ fontSize: 22 }} />;
    case "Filtro de Combustível":
      return <LocalGasStationRoundedIcon sx={{ fontSize: 22 }} />;
    default:
      return <FilterAltRoundedIcon sx={{ fontSize: 22 }} />;
  }
}

function formatDate(dateStr?: string): string {
  if (!dateStr) return "--";
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
}

export function UpcomingMaintenanceCard({
  item,
  onQuickRegister,
}: UpcomingMaintenanceCardProps) {
  const statusCfg = getStatusConfig(item.status);
  const clampedProgress = Math.min(100, Math.max(0, item.progressPercent || 0));

  const formatRemainingKm = () => {
    if (item.remainingKm < 0) {
      return `Excedido em ${Math.abs(item.remainingKm).toLocaleString("pt-BR")} km`;
    }
    if (item.remainingKm === 0) {
      return "Limite de km atingido";
    }
    return `Faltam ${item.remainingKm.toLocaleString("pt-BR")} km`;
  };

  const formatRemainingDays = () => {
    if (item.remainingDays < 0) {
      return `Vencido há ${Math.abs(item.remainingDays)} dias`;
    }
    if (item.remainingDays === 0) {
      return "Vence hoje";
    }
    return `em até ${item.remainingDays} dias`;
  };

  return (
    <Card
      elevation={0}
      sx={{
        height: "100%",
        display: "flex",
        flexDirection: "column",
        borderRadius: 2.5,
        border: `1px solid ${statusCfg.borderColor}`,
        bgcolor: "background.paper",
        transition: "transform 0.2s ease, box-shadow 0.2s ease",
        "&:hover": {
          transform: "translateY(-2px)",
          boxShadow: "0 8px 20px -4px rgba(15, 23, 42, 0.08)",
        },
      }}
    >
      <CardContent
        sx={{
          p: 2.5,
          flex: 1,
          display: "flex",
          flexDirection: "column",
          justifyContent: "space-between",
        }}
      >
        <Box>
          {/* Header: Icon + Label + Status Chip */}
          <Stack
            direction="row"
            spacing={1.5}
            sx={{
              alignItems: "flex-start",
              justifyContent: "space-between",
              mb: 2,
            }}
          >
            <Stack
              direction="row"
              spacing={1.5}
              sx={{ alignItems: "center", minWidth: 0, flex: 1, mr: 1 }}
            >
              <Box
                sx={{
                  width: 40,
                  height: 40,
                  borderRadius: 2,
                  bgcolor: "rgba(2, 132, 199, 0.1)",
                  color: "primary.main",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  flexShrink: 0,
                }}
              >
                {getItemIcon(item.type)}
              </Box>
              <Box sx={{ minWidth: 0 }}>
                <Typography
                  variant="subtitle1"
                  component="h3"
                  sx={{
                    fontWeight: 800,
                    color: "text.primary",
                    lineHeight: 1.25,
                    overflowWrap: "break-word",
                  }}
                >
                  {item.label}
                </Typography>
                <Typography
                  variant="caption"
                  sx={{ color: "text.secondary", display: "block" }}
                >
                  A cada {item.intervalKm.toLocaleString("pt-BR")} km ou{" "}
                  {item.intervalMonths} meses
                </Typography>
              </Box>
            </Stack>

            <Chip
              icon={statusCfg.icon}
              label={statusCfg.label}
              size="small"
              sx={{
                bgcolor: statusCfg.chipBg,
                color: statusCfg.chipColor,
                fontWeight: 700,
                borderRadius: 1.5,
                flexShrink: 0,
                "& .MuiChip-icon": {
                  color: statusCfg.chipColor,
                },
              }}
            />
          </Stack>

          {/* Progress Bar */}
          <Box sx={{ mb: 2 }}>
            <Stack
              direction="row"
              sx={{
                justifyContent: "space-between",
                alignItems: "center",
                mb: 0.75,
              }}
            >
              <Typography
                variant="caption"
                sx={{ fontWeight: 700, color: "text.primary" }}
              >
                {formatRemainingKm()}
              </Typography>
              <Typography
                variant="caption"
                sx={{ fontWeight: 600, color: "text.secondary" }}
              >
                {clampedProgress}% do ciclo
              </Typography>
            </Stack>
            <LinearProgress
              variant="determinate"
              value={clampedProgress}
              color={statusCfg.progressColor}
              sx={{
                height: 8,
                borderRadius: 4,
                bgcolor: "#F1F5F9",
              }}
            />
          </Box>

          {/* Metrics Details */}
          <Stack
            spacing={0.8}
            sx={{
              p: 1.5,
              borderRadius: 2,
              bgcolor: "#F8FAFC",
              border: "1px solid #F1F5F9",
              mb: 2,
            }}
          >
            <Stack
              direction="row"
              sx={{ justifyContent: "space-between", alignItems: "center" }}
            >
              <Stack
                direction="row"
                spacing={0.75}
                sx={{ alignItems: "center", color: "text.secondary" }}
              >
                <SpeedRoundedIcon sx={{ fontSize: 16 }} />
                <Typography variant="caption">Próxima troca (km):</Typography>
              </Stack>
              <Typography variant="caption" sx={{ fontWeight: 700 }}>
                {item.nextDueMileage.toLocaleString("pt-BR")} km
              </Typography>
            </Stack>

            <Stack
              direction="row"
              sx={{
                justifyContent: "space-between",
                alignItems: "center",
                gap: 1,
              }}
            >
              <Stack
                direction="row"
                spacing={0.75}
                sx={{
                  alignItems: "center",
                  color: "text.secondary",
                  flexShrink: 0,
                }}
              >
                <CalendarMonthRoundedIcon sx={{ fontSize: 16 }} />
                <Typography variant="caption">Data limite:</Typography>
              </Stack>
              <Typography
                variant="caption"
                sx={{ fontWeight: 700, textAlign: "right" }}
              >
                {formatDate(item.nextDueDate)} ({formatRemainingDays()})
              </Typography>
            </Stack>

            <Typography
              variant="caption"
              sx={{ color: "text.secondary", pt: 0.5, fontSize: "0.72rem" }}
            >
              {item.hasHistory && item.lastPerformedAt
                ? `Última troca: ${formatDate(item.lastPerformedAt)} aos ${item.lastMileage.toLocaleString("pt-BR")} km`
                : `Referência inicial: ${item.lastMileage.toLocaleString("pt-BR")} km (sem registro anterior)`}
            </Typography>
          </Stack>
        </Box>

        {/* Quick Action Button */}
        {onQuickRegister && (
          <Button
            variant={item.status === "ok" ? "outlined" : "contained"}
            color={
              item.status === "overdue"
                ? "error"
                : item.status === "upcoming"
                  ? "warning"
                  : "primary"
            }
            size="small"
            fullWidth
            startIcon={<AddCircleOutlineRoundedIcon />}
            onClick={() => onQuickRegister(item.type)}
            sx={{
              borderRadius: 2,
              textTransform: "none",
              fontWeight: 700,
              py: 0.8,
              whiteSpace: "normal",
              lineHeight: 1.3,
            }}
          >
            Registrar {item.label.toLowerCase()}
          </Button>
        )}
      </CardContent>
    </Card>
  );
}
