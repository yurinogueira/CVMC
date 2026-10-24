import { useState, useMemo } from "react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Stack,
  Box,
  Typography,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Card,
  CardContent,
  Chip,
  Alert,
  CircularProgress,
  IconButton,
  Grid,
} from "@mui/material";
import TuneRoundedIcon from "@mui/icons-material/TuneRounded";
import RestartAltRoundedIcon from "@mui/icons-material/RestartAltRounded";
import CloseRoundedIcon from "@mui/icons-material/CloseRounded";
import OpacityRoundedIcon from "@mui/icons-material/OpacityRounded";
import FilterAltRoundedIcon from "@mui/icons-material/FilterAltRounded";
import AirRoundedIcon from "@mui/icons-material/AirRounded";
import AcUnitRoundedIcon from "@mui/icons-material/AcUnitRounded";
import LocalGasStationRoundedIcon from "@mui/icons-material/LocalGasStationRounded";
import BuildRoundedIcon from "@mui/icons-material/BuildRounded";

import { Car, MaintenanceIntervalConfig } from "../../cars/types/car.types";
import { carService } from "../../cars/services/car.service";
import {
  PeriodicItemDefinition,
  getPeriodicItemsForVehicle,
  getKmOptionsForVehicle,
  getMonthOptionsForVehicle,
  getDefaultIntervalsForVehicle,
  isMotorcycle,
} from "../constants/periodicIntervalPresets";

interface MaintenanceIntervalsDialogProps {
  open: boolean;
  car: Car | null;
  onClose: () => void;
  onSaved?: (updatedCar: Car) => void;
}

interface MaintenanceIntervalsFormProps {
  car: Car;
  onClose: () => void;
  onSaved?: (updatedCar: Car) => void;
}

function getItemIcon(type: string) {
  switch (type) {
    case "Óleo de Motor":
      return <OpacityRoundedIcon sx={{ color: "#0284C7" }} />;
    case "Filtro do Óleo de Motor":
      return <FilterAltRoundedIcon sx={{ color: "#0284C7" }} />;
    case "Filtro de Ar do Motor":
      return <AirRoundedIcon sx={{ color: "#0D9488" }} />;
    case "Filtro do Ar-Condicionado":
      return <AcUnitRoundedIcon sx={{ color: "#0284C7" }} />;
    case "Filtro de Combustível":
      return <LocalGasStationRoundedIcon sx={{ color: "#D97706" }} />;
    default:
      return <BuildRoundedIcon sx={{ color: "#64748B" }} />;
  }
}

function formatMonthLabel(months: number): string {
  if (months === 12) return "12 meses (1 ano)";
  if (months === 24) return "24 meses (2 anos)";
  if (months === 6) return "6 meses (semestral)";
  if (months === 3) return "3 meses (trimestral)";
  if (months === 18) return "18 meses (1 ano e meio)";
  return `${months} meses`;
}

function computeInitialIntervals(
  car: Car,
  items: PeriodicItemDefinition[],
  defaults: Record<string, { intervalKm: number; intervalMonths: number }>,
): Record<string, MaintenanceIntervalConfig> {
  const initial: Record<string, MaintenanceIntervalConfig> = {};
  for (const item of items) {
    const custom = car.maintenanceIntervals?.[item.type];
    const def = defaults[item.type];
    initial[item.type] = {
      intervalKm: custom?.intervalKm || def?.intervalKm || 10000,
      intervalMonths: custom?.intervalMonths || def?.intervalMonths || 12,
    };
  }
  return initial;
}

function MaintenanceIntervalsForm({
  car,
  onClose,
  onSaved,
}: MaintenanceIntervalsFormProps) {
  const vehicleType = car.vehicleType;
  const items = useMemo(
    () => getPeriodicItemsForVehicle(vehicleType),
    [vehicleType],
  );
  const kmOptions = useMemo(
    () => getKmOptionsForVehicle(vehicleType),
    [vehicleType],
  );
  const monthOptions = useMemo(() => getMonthOptionsForVehicle(), []);
  const defaultIntervals = useMemo(
    () => getDefaultIntervalsForVehicle(vehicleType),
    [vehicleType],
  );

  const [intervals, setIntervals] = useState<
    Record<string, MaintenanceIntervalConfig>
  >(() => computeInitialIntervals(car, items, defaultIntervals));
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const handleKmChange = (itemType: string, km: number) => {
    setIntervals((prev) => ({
      ...prev,
      [itemType]: {
        ...(prev[itemType] || { intervalMonths: 12 }),
        intervalKm: km,
      },
    }));
  };

  const handleMonthsChange = (itemType: string, months: number) => {
    setIntervals((prev) => ({
      ...prev,
      [itemType]: {
        ...(prev[itemType] || { intervalKm: 10000 }),
        intervalMonths: months,
      },
    }));
  };

  const handleResetToDefaults = () => {
    const resetValues: Record<string, MaintenanceIntervalConfig> = {};
    for (const item of items) {
      const def = defaultIntervals[item.type];
      if (def) {
        resetValues[item.type] = { ...def };
      }
    }
    setIntervals(resetValues);
  };

  const handleSave = async () => {
    try {
      setLoading(true);
      setErrorMsg(null);

      const updatedCar = await carService.update(car.id, {
        name: car.name,
        manufacturer: car.manufacturer,
        model: car.model,
        yearManufacture: car.yearManufacture,
        yearModel: car.yearModel,
        lastMileage: car.lastMileage,
        vehicleType: car.vehicleType,
        imageUrl: car.imageUrl,
        fipeCode: car.fipeCode,
        fipePrice: car.fipePrice,
        fuel: car.fuel,
        maintenanceIntervals: intervals,
      });

      onSaved?.(updatedCar);
      onClose();
    } catch (err: unknown) {
      const errorObj = err as { response?: { data?: { message?: string } } };
      setErrorMsg(
        errorObj.response?.data?.message ||
          "Não foi possível salvar os intervalos de revisão. Tente novamente.",
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <DialogTitle
        id="maintenance-intervals-dialog-title"
        sx={{
          fontWeight: 800,
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          pb: 1,
        }}
      >
        <Stack direction="row" spacing={1.5} sx={{ alignItems: "center" }}>
          <Box
            sx={{
              width: 38,
              height: 38,
              borderRadius: 2,
              bgcolor: "rgba(2, 132, 199, 0.1)",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              color: "primary.main",
            }}
          >
            <TuneRoundedIcon />
          </Box>
          <Box>
            <Typography variant="h6" sx={{ fontWeight: 800, lineHeight: 1.2 }}>
              Intervalos de Revisão
            </Typography>
            <Typography variant="caption" sx={{ color: "text.secondary" }}>
              Personalize a frequência para <strong>{car.name}</strong> (
              {isMotorcycle(car.vehicleType) ? "Moto" : "Carro"})
            </Typography>
          </Box>
        </Stack>

        <IconButton
          onClick={onClose}
          size="small"
          aria-label="fechar"
          sx={{ color: "text.secondary" }}
        >
          <CloseRoundedIcon />
        </IconButton>
      </DialogTitle>

      <DialogContent
        dividers
        sx={{ borderColor: "#F1F5F9", px: { xs: 2, sm: 3 }, py: 2.5 }}
      >
        {errorMsg && (
          <Alert
            severity="error"
            onClose={() => setErrorMsg(null)}
            sx={{ mb: 2, borderRadius: 2 }}
          >
            {errorMsg}
          </Alert>
        )}

        <Stack
          direction="row"
          sx={{
            justifyContent: "space-between",
            alignItems: "center",
            mb: 2.5,
            flexWrap: "wrap",
            gap: 1,
          }}
        >
          <Typography variant="body2" sx={{ color: "text.secondary" }}>
            Selecione as metas de quilometragem e tempo desejadas:
          </Typography>
          <Button
            size="small"
            variant="text"
            startIcon={<RestartAltRoundedIcon fontSize="small" />}
            onClick={handleResetToDefaults}
            sx={{
              fontWeight: 700,
              color: "text.secondary",
              textTransform: "none",
              borderRadius: 1.5,
              "&:hover": {
                color: "primary.main",
                bgcolor: "rgba(2, 132, 199, 0.08)",
              },
            }}
          >
            Restaurar Padrões Recomendados
          </Button>
        </Stack>

        <Stack spacing={2}>
          {items.map((item) => {
            const currentKm = intervals[item.type]?.intervalKm || 10000;
            const currentMonths = intervals[item.type]?.intervalMonths || 12;
            const def = defaultIntervals[item.type];
            const isCustom =
              def &&
              (currentKm !== def.intervalKm ||
                currentMonths !== def.intervalMonths);

            return (
              <Card
                key={item.type}
                elevation={0}
                sx={{
                  border: "1px solid",
                  borderColor: isCustom ? "primary.light" : "#E2E8F0",
                  borderRadius: 2.5,
                  bgcolor: isCustom
                    ? "rgba(2, 132, 199, 0.02)"
                    : "background.paper",
                  transition: "all 0.2s ease-in-out",
                }}
              >
                <CardContent
                  sx={{
                    p: { xs: 2, sm: 2.5 },
                    "&:last-child": { pb: { xs: 2, sm: 2.5 } },
                  }}
                >
                  <Stack
                    direction="row"
                    spacing={1.5}
                    sx={{
                      alignItems: "flex-start",
                      justifyContent: "space-between",
                      mb: 1.5,
                    }}
                  >
                    <Stack
                      direction="row"
                      spacing={1.5}
                      sx={{ alignItems: "center" }}
                    >
                      <Box
                        sx={{
                          width: 36,
                          height: 36,
                          borderRadius: 2,
                          bgcolor: "#F8FAFC",
                          display: "flex",
                          alignItems: "center",
                          justifyContent: "center",
                        }}
                      >
                        {getItemIcon(item.type)}
                      </Box>
                      <Box>
                        <Typography
                          variant="subtitle2"
                          sx={{ fontWeight: 800, color: "text.primary" }}
                        >
                          {item.label}
                        </Typography>
                        <Typography
                          variant="caption"
                          sx={{ color: "text.secondary" }}
                        >
                          {item.description}
                        </Typography>
                      </Box>
                    </Stack>

                    <Chip
                      size="small"
                      label={isCustom ? "Personalizado" : "Padrão"}
                      sx={{
                        fontWeight: 700,
                        fontSize: "0.7rem",
                        borderRadius: 1,
                        bgcolor: isCustom
                          ? "rgba(2, 132, 199, 0.12)"
                          : "#F1F5F9",
                        color: isCustom ? "primary.main" : "text.secondary",
                      }}
                    />
                  </Stack>

                  <Grid container spacing={1.5}>
                    <Grid size={{ xs: 12, sm: 6 }}>
                      <FormControl fullWidth size="small">
                        <InputLabel id={`km-label-${item.type}`}>
                          Quilometragem
                        </InputLabel>
                        <Select
                          labelId={`km-label-${item.type}`}
                          id={`km-select-${item.type}`}
                          value={currentKm}
                          label="Quilometragem"
                          onChange={(e) =>
                            handleKmChange(item.type, Number(e.target.value))
                          }
                          sx={{ borderRadius: 1.5, fontWeight: 600 }}
                        >
                          {kmOptions.map((km) => (
                            <MenuItem
                              key={km}
                              value={km}
                              sx={{ fontWeight: 600 }}
                            >
                              A cada {km.toLocaleString("pt-BR")} km
                            </MenuItem>
                          ))}
                        </Select>
                      </FormControl>
                    </Grid>

                    <Grid size={{ xs: 12, sm: 6 }}>
                      <FormControl fullWidth size="small">
                        <InputLabel id={`month-label-${item.type}`}>
                          Tempo (Meses)
                        </InputLabel>
                        <Select
                          labelId={`month-label-${item.type}`}
                          id={`month-select-${item.type}`}
                          value={currentMonths}
                          label="Tempo (Meses)"
                          onChange={(e) =>
                            handleMonthsChange(
                              item.type,
                              Number(e.target.value),
                            )
                          }
                          sx={{ borderRadius: 1.5, fontWeight: 600 }}
                        >
                          {monthOptions.map((months) => (
                            <MenuItem
                              key={months}
                              value={months}
                              sx={{ fontWeight: 600 }}
                            >
                              A cada {formatMonthLabel(months)}
                            </MenuItem>
                          ))}
                        </Select>
                      </FormControl>
                    </Grid>
                  </Grid>
                </CardContent>
              </Card>
            );
          })}
        </Stack>
      </DialogContent>

      <DialogActions sx={{ px: { xs: 2, sm: 3 }, py: 2 }}>
        <Button
          onClick={onClose}
          color="inherit"
          disabled={loading}
          sx={{ borderRadius: 2, fontWeight: 600, color: "text.secondary" }}
        >
          Cancelar
        </Button>
        <Button
          onClick={handleSave}
          variant="contained"
          disabled={loading}
          startIcon={
            loading ? <CircularProgress size={18} color="inherit" /> : null
          }
          sx={{ borderRadius: 2, px: 3, fontWeight: 700 }}
        >
          {loading ? "Salvando..." : "Salvar Intervalos"}
        </Button>
      </DialogActions>
    </>
  );
}

export function MaintenanceIntervalsDialog({
  open,
  car,
  onClose,
  onSaved,
}: MaintenanceIntervalsDialogProps) {
  if (!car) return null;

  return (
    <Dialog
      open={open}
      onClose={onClose}
      maxWidth="sm"
      fullWidth
      aria-labelledby="maintenance-intervals-dialog-title"
      slotProps={{
        paper: {
          sx: {
            borderRadius: { xs: 2, sm: 3 },
            p: { xs: 0.5, sm: 1 },
          },
        },
      }}
    >
      <MaintenanceIntervalsForm
        key={`${car.id}-${open}`}
        car={car}
        onClose={onClose}
        onSaved={onSaved}
      />
    </Dialog>
  );
}
