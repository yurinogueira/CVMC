import { useEffect, useState, useCallback } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import {
  Box,
  Typography,
  Card,
  Button,
  Stack,
  Chip,
  Grid,
  Skeleton,
  Alert,
  Pagination,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  CircularProgress,
  Snackbar,
} from "@mui/material";
import BuildCircleRoundedIcon from "@mui/icons-material/BuildCircleRounded";
import AddRoundedIcon from "@mui/icons-material/AddRounded";
import CheckCircleRoundedIcon from "@mui/icons-material/CheckCircleRounded";
import WarningAmberRoundedIcon from "@mui/icons-material/WarningAmberRounded";
import ErrorOutlineRoundedIcon from "@mui/icons-material/ErrorOutlineRounded";
import DirectionsCarRoundedIcon from "@mui/icons-material/DirectionsCarRounded";
import TwoWheelerRoundedIcon from "@mui/icons-material/TwoWheelerRounded";
import LocalShippingRoundedIcon from "@mui/icons-material/LocalShippingRounded";
import HistoryRoundedIcon from "@mui/icons-material/HistoryRounded";
import EventRepeatRoundedIcon from "@mui/icons-material/EventRepeatRounded";
import TuneRoundedIcon from "@mui/icons-material/TuneRounded";
import { useDocumentTitle } from "../../shared";
import { brandColors } from "../../../styles/theme";
import { carService } from "../../cars/services/car.service";
import { Car } from "../../cars/types/car.types";
import { maintenanceService } from "../services/maintenance.service";
import {
  Maintenance,
  FleetUpcomingSummary,
  CarUpcomingMaintenances,
} from "../types/maintenance.types";
import { UpcomingMaintenanceCard } from "../components/UpcomingMaintenanceCard";
import { MaintenanceCard } from "../components/MaintenanceCard";
import { AddMaintenanceDialog } from "../components/AddMaintenanceDialog";
import { MaintenanceIntervalsDialog } from "../components/MaintenanceIntervalsDialog";

const ITEMS_PER_PAGE = 5;

function getVehicleIcon(vehicleType?: string) {
  switch (vehicleType) {
    case "motorcycles":
      return <TwoWheelerRoundedIcon fontSize="small" />;
    case "trucks":
      return <LocalShippingRoundedIcon fontSize="small" />;
    default:
      return <DirectionsCarRoundedIcon fontSize="small" />;
  }
}

export function MaintenancePage() {
  useDocumentTitle("Manutenções & Revisões");
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();
  const queryCarId = searchParams.get("carId");

  const [cars, setCars] = useState<Car[]>([]);
  const [selectedCarId, setSelectedCarId] = useState<string>("");
  const [fleetSummary, setFleetSummary] = useState<FleetUpcomingSummary | null>(
    null,
  );
  const [maintenances, setMaintenances] = useState<Maintenance[]>([]);
  const [loadingInitial, setLoadingInitial] = useState(true);
  const [loadingHistory, setLoadingHistory] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [page, setPage] = useState(1);

  // Quick register / New maintenance dialog state
  const [dialogOpen, setDialogOpen] = useState(false);
  const [preselectedTypes, setPreselectedTypes] = useState<string[]>([]);

  // Delete confirmation state
  const [maintenanceToDelete, setMaintenanceToDelete] =
    useState<Maintenance | null>(null);
  const [deleting, setDeleting] = useState(false);
  const [toastMessage, setToastMessage] = useState<string | null>(null);
  const [intervalsDialogOpen, setIntervalsDialogOpen] = useState(false);

  const refreshData = useCallback(
    async (targetCarId?: string) => {
      try {
        const [carsData, summaryData] = await Promise.all([
          carService.list(),
          maintenanceService.getUpcomingSummary().catch(() => null),
        ]);

        setCars(carsData);
        if (summaryData) {
          setFleetSummary(summaryData);
        }

        const effectiveCarId =
          targetCarId ||
          (queryCarId && carsData.some((c) => c.id === queryCarId)
            ? queryCarId
            : carsData[0]?.id || "");

        if (effectiveCarId) {
          setSelectedCarId(effectiveCarId);
          const list = await maintenanceService
            .listByCar(effectiveCarId)
            .catch(() => []);
          const sorted = [...list].sort(
            (a, b) => new Date(b.date).getTime() - new Date(a.date).getTime(),
          );
          setMaintenances(sorted);
        } else {
          setMaintenances([]);
        }
      } catch {
        setErrorMsg(
          "Não foi possível carregar os dados de veículos e revisões preventivas.",
        );
      }
    },
    [queryCarId],
  );

  useEffect(() => {
    let isMounted = true;
    Promise.all([
      carService.list(),
      maintenanceService.getUpcomingSummary().catch(() => null),
    ])
      .then(async ([carsData, summaryData]) => {
        if (!isMounted) return;
        setCars(carsData);
        if (summaryData) {
          setFleetSummary(summaryData);
        }
        const initialCarId =
          queryCarId && carsData.some((c) => c.id === queryCarId)
            ? queryCarId
            : carsData[0]?.id || "";
        if (initialCarId) {
          setSelectedCarId(initialCarId);
          const list = await maintenanceService
            .listByCar(initialCarId)
            .catch(() => []);
          if (!isMounted) return;
          const sorted = [...list].sort(
            (a, b) => new Date(b.date).getTime() - new Date(a.date).getTime(),
          );
          setMaintenances(sorted);
        }
      })
      .catch(() => {
        if (isMounted) {
          setErrorMsg(
            "Não foi possível carregar os dados de veículos e revisões preventivas.",
          );
        }
      })
      .finally(() => {
        if (isMounted) {
          setLoadingInitial(false);
        }
      });

    return () => {
      isMounted = false;
    };
  }, [queryCarId]);

  const selectedCar = cars.find((c) => c.id === selectedCarId) || null;
  const selectedUpcoming: CarUpcomingMaintenances | undefined =
    fleetSummary?.vehicles.find((v) => v.carId === selectedCarId);

  const handleSelectCar = async (carId: string) => {
    setSelectedCarId(carId);
    setSearchParams({ carId }, { replace: true });
    setLoadingHistory(true);
    try {
      const list = await maintenanceService.listByCar(carId);
      const sorted = [...list].sort(
        (a, b) => new Date(b.date).getTime() - new Date(a.date).getTime(),
      );
      setMaintenances(sorted);
      setPage(1);
    } catch {
      setMaintenances([]);
    } finally {
      setLoadingHistory(false);
    }
  };

  const handleOpenNewMaintenance = (initialType?: string) => {
    setPreselectedTypes(initialType ? [initialType] : []);
    setDialogOpen(true);
  };

  const handleMaintenanceCreated = async () => {
    setToastMessage("Manutenção registrada com sucesso!");
    await refreshData(selectedCarId);
  };

  const handleConfirmDelete = async () => {
    if (!maintenanceToDelete) return;
    try {
      setDeleting(true);
      await maintenanceService.delete(maintenanceToDelete.id);
      setMaintenanceToDelete(null);
      setToastMessage("Manutenção excluída com sucesso.");
      await refreshData(selectedCarId);
    } catch {
      setErrorMsg("Não foi possível excluir o registro de manutenção.");
    } finally {
      setDeleting(false);
    }
  };

  // Fleet Banner Dynamic Appearance
  const getFleetBannerConfig = () => {
    const overdueCount = fleetSummary?.overdueCount || 0;
    const upcomingCount = fleetSummary?.upcomingCount || 0;

    if (overdueCount > 0) {
      return {
        title: "Atenção: Manutenções Vencidas",
        chipLabel: `${overdueCount} ${overdueCount === 1 ? "item vencido" : "itens vencidos"}`,
        chipBg: "#FEE2E2",
        chipColor: "#991B1B",
        iconBg: "rgba(220, 38, 38, 0.12)",
        iconColor: "#DC2626",
        borderColor: "#FCA5A5",
        icon: <ErrorOutlineRoundedIcon sx={{ fontSize: 28 }} />,
        description:
          upcomingCount > 0
            ? `Você possui ${overdueCount} revisão(ões) com prazo ou quilometragem excedida e ${upcomingCount} próxima(s) do vencimento.`
            : `Você possui ${overdueCount} revisão(ões) preventiva(s) com prazo ou quilometragem excedida. Agende a manutenção o quanto antes.`,
      };
    }

    if (upcomingCount > 0) {
      return {
        title: "Revisões Preventivas Próximas",
        chipLabel: `${upcomingCount} ${upcomingCount === 1 ? "item próximo" : "itens próximos"}`,
        chipBg: "#FEF3C7",
        chipColor: "#92400E",
        iconBg: "rgba(217, 119, 6, 0.12)",
        iconColor: "#D97706",
        borderColor: "#FCD34D",
        icon: <WarningAmberRoundedIcon sx={{ fontSize: 28 }} />,
        description: `Existem ${upcomingCount} revisão(ões) periódica(s) próximas do limite de quilometragem ou tempo (12 meses).`,
      };
    }

    return {
      title: "Sua frota está em dia",
      chipLabel: "Tudo Certo",
      chipBg: brandColors.green,
      chipColor: "#064E3B",
      iconBg: "rgba(22, 163, 74, 0.12)",
      iconColor: "#16A34A",
      borderColor: "#E2E8F0",
      icon: <CheckCircleRoundedIcon sx={{ fontSize: 28 }} />,
      description:
        "Nenhum alerta de manutenção atrasada ou revisão pendente no momento.",
    };
  };

  const bannerCfg = getFleetBannerConfig();
  const totalPages = Math.max(
    1,
    Math.ceil(maintenances.length / ITEMS_PER_PAGE),
  );
  const paginatedMaintenances = maintenances.slice(
    (page - 1) * ITEMS_PER_PAGE,
    page * ITEMS_PER_PAGE,
  );

  return (
    <Box>
      {/* Header */}
      <Stack
        direction={{ xs: "column", sm: "row" }}
        spacing={2}
        sx={{
          alignItems: { xs: "flex-start", sm: "center" },
          justifyContent: "space-between",
          mb: 3,
        }}
      >
        <Box>
          <Typography
            component="h1"
            variant="h5"
            sx={{ fontWeight: 800, color: "text.primary" }}
          >
            Manutenções & Revisões
          </Typography>
          <Typography variant="body2" sx={{ color: "text.secondary" }}>
            Monitoramento inteligente de revisões periódicas e histórico
            completo por veículo.
          </Typography>
        </Box>

        <Button
          variant="contained"
          startIcon={<AddRoundedIcon />}
          disabled={!selectedCar}
          onClick={() => handleOpenNewMaintenance()}
          sx={{ px: 2.5, py: 1.1 }}
        >
          Registrar Manutenção
        </Button>
      </Stack>

      {errorMsg && (
        <Alert
          severity="error"
          onClose={() => setErrorMsg(null)}
          sx={{ mb: 3, borderRadius: 2 }}
        >
          {errorMsg}
        </Alert>
      )}

      {/* Dynamic Fleet Info Status Banner */}
      <Card
        elevation={0}
        sx={{
          mb: 4,
          p: 3,
          bgcolor: "background.paper",
          border: `1px solid ${bannerCfg.borderColor}`,
          display: "flex",
          alignItems: "center",
          gap: 2,
        }}
      >
        <Box
          sx={{
            width: 48,
            height: 48,
            borderRadius: 2,
            bgcolor: bannerCfg.iconBg,
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            color: bannerCfg.iconColor,
            flexShrink: 0,
          }}
        >
          {bannerCfg.icon}
        </Box>
        <Box sx={{ flex: 1 }}>
          <Stack
            direction="row"
            spacing={1}
            sx={{ alignItems: "center", mb: 0.5, flexWrap: "wrap", gap: 0.5 }}
          >
            <Typography
              component="h2"
              variant="subtitle1"
              sx={{ fontWeight: 700 }}
            >
              {bannerCfg.title}
            </Typography>
            <Chip
              label={bannerCfg.chipLabel}
              size="small"
              sx={{
                bgcolor: bannerCfg.chipBg,
                color: bannerCfg.chipColor,
                fontWeight: 700,
                borderRadius: 1,
              }}
            />
          </Stack>
          <Typography variant="body2" sx={{ color: "text.secondary" }}>
            {bannerCfg.description}
          </Typography>
        </Box>
      </Card>

      {loadingInitial ? (
        <Grid container spacing={2.5}>
          {[1, 2, 3].map((idx) => (
            <Grid size={{ xs: 12, sm: 6, md: 6, lg: 4 }} key={idx}>
              <Card
                elevation={0}
                sx={{ p: 2.5, border: "1px solid #E2E8F0", borderRadius: 2.5 }}
              >
                <Skeleton variant="text" width="60%" height={28} />
                <Skeleton
                  variant="text"
                  width="40%"
                  height={20}
                  sx={{ mb: 2 }}
                />
                <Skeleton
                  variant="rectangular"
                  height={90}
                  sx={{ borderRadius: 2 }}
                />
              </Card>
            </Grid>
          ))}
        </Grid>
      ) : cars.length === 0 ? (
        /* Empty Fleet State */
        <Card
          elevation={0}
          sx={{
            py: 8,
            px: 3,
            textAlign: "center",
            border: "1px dashed #CBD5E1",
            bgcolor: "background.paper",
          }}
        >
          <Box
            sx={{
              width: 56,
              height: 56,
              borderRadius: "50%",
              bgcolor: "rgba(2, 132, 199, 0.1)",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              mx: "auto",
              mb: 2,
              color: "primary.main",
            }}
          >
            <DirectionsCarRoundedIcon sx={{ fontSize: 32 }} />
          </Box>
          <Typography
            component="h2"
            variant="h6"
            sx={{ fontWeight: 700, mb: 0.5 }}
          >
            Nenhum veículo cadastrado
          </Typography>
          <Typography
            variant="body2"
            sx={{ color: "text.secondary", maxWidth: 460, mx: "auto", mb: 3 }}
          >
            Cadastre seu veículo primeiro para acompanhar as próximas revisões
            preventivas e registrar o histórico de manutenções.
          </Typography>
          <Button
            variant="contained"
            onClick={() => navigate("/vehicles")}
            startIcon={<AddRoundedIcon />}
          >
            Ir para Meus Veículos
          </Button>
        </Card>
      ) : (
        <>
          {/* Vehicle Selector */}
          <Box sx={{ mb: 3.5 }}>
            <Typography
              variant="caption"
              sx={{
                fontWeight: 700,
                color: "text.secondary",
                textTransform: "uppercase",
                letterSpacing: 0.5,
                display: "block",
                mb: 1,
              }}
            >
              Veículo Selecionado
            </Typography>
            <Stack
              direction="row"
              spacing={1}
              sx={{ flexWrap: "wrap", gap: 1 }}
            >
              {cars.map((car) => {
                const isSelected = car.id === selectedCarId;
                const carSummary = fleetSummary?.vehicles.find(
                  (v) => v.carId === car.id,
                );
                const alertCount =
                  (carSummary?.overdueCount || 0) +
                  (carSummary?.upcomingCount || 0);

                return (
                  <Chip
                    key={car.id}
                    icon={getVehicleIcon(car.vehicleType)}
                    label={`${car.name} (${(car.lastMileage || 0).toLocaleString("pt-BR")} km)${
                      alertCount > 0 ? ` • ${alertCount}` : ""
                    }`}
                    onClick={() => handleSelectCar(car.id)}
                    color={isSelected ? "primary" : "default"}
                    variant={isSelected ? "filled" : "outlined"}
                    sx={{
                      borderRadius: 2,
                      py: 2.2,
                      px: 1,
                      fontWeight: isSelected ? 700 : 600,
                      cursor: "pointer",
                      borderColor: isSelected ? "primary.main" : "#CBD5E1",
                    }}
                  />
                );
              })}
            </Stack>
          </Box>

          {/* Upcoming Periodic Maintenances Section */}
          {selectedUpcoming && selectedUpcoming.items.length > 0 && (
            <Box sx={{ mb: 5 }}>
              <Stack
                direction={{ xs: "column", sm: "row" }}
                spacing={1.5}
                sx={{
                  alignItems: { xs: "flex-start", sm: "center" },
                  justifyContent: "space-between",
                  mb: 2,
                }}
              >
                <Stack
                  direction="row"
                  spacing={1}
                  sx={{ alignItems: "center" }}
                >
                  <EventRepeatRoundedIcon color="primary" />
                  <Box>
                    <Typography
                      component="h2"
                      variant="h6"
                      sx={{ fontWeight: 800, lineHeight: 1.2 }}
                    >
                      Próximas Manutenções Preventivas
                    </Typography>
                    <Typography
                      variant="caption"
                      sx={{ color: "text.secondary" }}
                    >
                      Monitoramento automático por quilometragem e tempo para{" "}
                      <strong>{selectedUpcoming.carName}</strong>
                    </Typography>
                  </Box>
                </Stack>

                <Button
                  variant="outlined"
                  size="small"
                  startIcon={<TuneRoundedIcon />}
                  onClick={() => setIntervalsDialogOpen(true)}
                  disabled={!selectedCar}
                  sx={{
                    borderRadius: 2,
                    textTransform: "none",
                    fontWeight: 700,
                    borderColor: "#CBD5E1",
                    color: "text.primary",
                    "&:hover": {
                      borderColor: "primary.main",
                      bgcolor: "rgba(2, 132, 199, 0.04)",
                      color: "primary.main",
                    },
                  }}
                >
                  Ajustar Intervalos
                </Button>
              </Stack>

              <Grid container spacing={2.5}>
                {selectedUpcoming.items.map((item) => (
                  <Grid size={{ xs: 12, sm: 6, md: 6, lg: 4 }} key={item.type}>
                    <UpcomingMaintenanceCard
                      item={item}
                      onQuickRegister={(type) => handleOpenNewMaintenance(type)}
                    />
                  </Grid>
                ))}
              </Grid>
            </Box>
          )}

          {/* Chronological Maintenance History Section */}
          <Box>
            <Stack
              direction={{ xs: "column", sm: "row" }}
              spacing={1.5}
              sx={{
                alignItems: { xs: "flex-start", sm: "center" },
                justifyContent: "space-between",
                mb: 2,
              }}
            >
              <Stack direction="row" spacing={1} sx={{ alignItems: "center" }}>
                <HistoryRoundedIcon color="primary" />
                <Typography
                  component="h2"
                  variant="h6"
                  sx={{ fontWeight: 800 }}
                >
                  Histórico de Manutenções Realizadas
                </Typography>
                <Chip
                  label={maintenances.length}
                  size="small"
                  sx={{ fontWeight: 700, borderRadius: 1 }}
                />
              </Stack>
            </Stack>

            {loadingHistory ? (
              <Stack spacing={2}>
                {[1, 2].map((i) => (
                  <Skeleton
                    key={i}
                    variant="rectangular"
                    height={110}
                    sx={{ borderRadius: 2.5 }}
                  />
                ))}
              </Stack>
            ) : maintenances.length === 0 ? (
              <Card
                elevation={0}
                sx={{
                  py: 7,
                  px: 3,
                  textAlign: "center",
                  border: "1px dashed #CBD5E1",
                  bgcolor: "background.paper",
                }}
              >
                <Box
                  sx={{
                    width: 56,
                    height: 56,
                    borderRadius: "50%",
                    bgcolor: "rgba(2, 132, 199, 0.1)",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    mx: "auto",
                    mb: 2,
                    color: "primary.main",
                  }}
                >
                  <BuildCircleRoundedIcon sx={{ fontSize: 32 }} />
                </Box>
                <Typography
                  component="h3"
                  variant="h6"
                  sx={{ fontWeight: 700, mb: 0.5 }}
                >
                  Nenhuma manutenção registrada
                </Typography>
                <Typography
                  variant="body2"
                  sx={{
                    color: "text.secondary",
                    maxWidth: 460,
                    mx: "auto",
                    mb: 3,
                  }}
                >
                  Registre trocas de óleo, filtros, pastilhas de freio ou
                  revisões para atualizar os indicadores preventivos de{" "}
                  <strong>{selectedCar?.name}</strong>.
                </Typography>
                <Button
                  variant="contained"
                  startIcon={<AddRoundedIcon />}
                  onClick={() => handleOpenNewMaintenance()}
                >
                  Nova Manutenção
                </Button>
              </Card>
            ) : (
              <Stack spacing={2}>
                {paginatedMaintenances.map((maint) => (
                  <MaintenanceCard
                    key={maint.id}
                    maintenance={maint}
                    onEdit={(m) =>
                      navigate(`/vehicles/${m.carId}/maintenance/${m.id}/edit`)
                    }
                    onDelete={(maintId) => {
                      const target = maintenances.find((x) => x.id === maintId);
                      if (target) setMaintenanceToDelete(target);
                    }}
                  />
                ))}

                {totalPages > 1 && (
                  <Box
                    sx={{
                      display: "flex",
                      justifyContent: "center",
                      pt: 2,
                    }}
                  >
                    <Pagination
                      count={totalPages}
                      page={page}
                      onChange={(_, value) => setPage(value)}
                      color="primary"
                      shape="rounded"
                    />
                  </Box>
                )}
              </Stack>
            )}
          </Box>
        </>
      )}

      {/* Add Maintenance Dialog with Preselected Type & Current Mileage */}
      {selectedCar && (
        <AddMaintenanceDialog
          open={dialogOpen}
          carId={selectedCar.id}
          carName={selectedCar.name}
          lastMileage={selectedCar.lastMileage}
          initialTypes={preselectedTypes}
          onClose={() => setDialogOpen(false)}
          onMaintenanceCreated={() => {
            void handleMaintenanceCreated();
          }}
        />
      )}

      {/* Delete Confirmation Dialog */}
      <Dialog
        open={Boolean(maintenanceToDelete)}
        onClose={() => !deleting && setMaintenanceToDelete(null)}
        maxWidth="xs"
        fullWidth
      >
        <DialogTitle sx={{ fontWeight: 700 }}>Excluir Manutenção</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Tem certeza que deseja excluir o registro{" "}
            <strong>{maintenanceToDelete?.title}</strong>? Esta ação atualizará
            o cálculo das próximas revisões preventivas.
          </DialogContentText>
        </DialogContent>
        <DialogActions sx={{ p: 2 }}>
          <Button
            onClick={() => setMaintenanceToDelete(null)}
            disabled={deleting}
            color="inherit"
          >
            Cancelar
          </Button>
          <Button
            onClick={() => void handleConfirmDelete()}
            color="error"
            variant="contained"
            disabled={deleting}
            startIcon={
              deleting ? <CircularProgress size={16} color="inherit" /> : null
            }
          >
            Excluir
          </Button>
        </DialogActions>
      </Dialog>

      {/* Maintenance Intervals Dialog */}
      <MaintenanceIntervalsDialog
        open={intervalsDialogOpen}
        car={selectedCar}
        onClose={() => setIntervalsDialogOpen(false)}
        onSaved={async () => {
          await refreshData(selectedCarId);
          setToastMessage("Intervalos de revisão atualizados com sucesso!");
        }}
      />

      {/* Success Toast */}
      <Snackbar
        open={Boolean(toastMessage)}
        autoHideDuration={4000}
        onClose={() => setToastMessage(null)}
        anchorOrigin={{ vertical: "bottom", horizontal: "right" }}
      >
        <Alert
          onClose={() => setToastMessage(null)}
          severity="success"
          variant="filled"
          sx={{ width: "100%" }}
        >
          {toastMessage}
        </Alert>
      </Snackbar>
    </Box>
  );
}
