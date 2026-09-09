import { useEffect, useState } from "react";
import { useParams, useNavigate, useLocation } from "react-router-dom";
import {
  Box,
  Typography,
  Grid,
  Card,
  CardContent,
  Button,
  Stack,
  Chip,
  Skeleton,
  Alert,
  Snackbar,
  Divider,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  CircularProgress,
  Tabs,
  Tab,
} from "@mui/material";
import ArrowBackRoundedIcon from "@mui/icons-material/ArrowBackRounded";
import AddRoundedIcon from "@mui/icons-material/AddRounded";
import BuildCircleRoundedIcon from "@mui/icons-material/BuildCircleRounded";
import SpeedRoundedIcon from "@mui/icons-material/SpeedRounded";
import CalendarMonthRoundedIcon from "@mui/icons-material/CalendarMonthRounded";
import LocalGasStationRoundedIcon from "@mui/icons-material/LocalGasStationRounded";
import TagRoundedIcon from "@mui/icons-material/TagRounded";
import MonetizationOnRoundedIcon from "@mui/icons-material/MonetizationOnRounded";
import DirectionsCarRoundedIcon from "@mui/icons-material/DirectionsCarRounded";
import TwoWheelerRoundedIcon from "@mui/icons-material/TwoWheelerRounded";
import LocalShippingRoundedIcon from "@mui/icons-material/LocalShippingRounded";
import DeleteOutlineRoundedIcon from "@mui/icons-material/DeleteOutlineRounded";

import { carService } from "../services/car.service";
import { Car } from "../types/car.types";
import { VehicleImage } from "../components/VehicleImage";
import { maintenanceService } from "../../maintenance/services/maintenance.service";
import { Maintenance } from "../../maintenance/types/maintenance.types";
import { MaintenanceCard } from "../../maintenance/components/MaintenanceCard";
import { fuelService } from "../../fuel/services/fuel.service";
import { Fueling } from "../../fuel/types/fuel.types";
import { FuelCard } from "../../fuel/components/FuelCard";
import { AddFuelingDialog } from "../../fuel/components/AddFuelingDialog";
import { useDocumentTitle } from "../../shared";
import { brandColors } from "../../../styles/theme";

export function VehicleDetailsPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [car, setCar] = useState<Car | null>(null);
  const [maintenances, setMaintenances] = useState<Maintenance[]>([]);
  const [fuelings, setFuelings] = useState<Fueling[]>([]);
  const [loadingCar, setLoadingCar] = useState(true);
  const [loadingMaintenances, setLoadingMaintenances] = useState(true);
  const [loadingFuelings, setLoadingFuelings] = useState(true);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const location = useLocation();
  const [toastMessage, setToastMessage] = useState<string | null>(
    () => (location.state as { message?: string } | null)?.message || null,
  );

  const [activeTab, setActiveTab] = useState<number>(0);
  const [fuelDialogOpen, setFuelDialogOpen] = useState<boolean>(false);
  const [selectedFueling, setSelectedFueling] = useState<Fueling | null>(null);

  const [fuelingToDelete, setFuelingToDelete] = useState<Fueling | null>(null);
  const [deletingFueling, setDeletingFueling] = useState(false);
  const [deleteFuelingError, setDeleteFuelingError] = useState<string | null>(
    null,
  );

  const [maintenanceToDelete, setMaintenanceToDelete] =
    useState<Maintenance | null>(null);
  const [deleting, setDeleting] = useState(false);
  const [deleteError, setDeleteError] = useState<string | null>(null);

  useDocumentTitle(car ? `${car.name} - Detalhes` : "Detalhes do Veículo");

  useEffect(() => {
    if (location.state?.message) {
      window.history.replaceState({}, document.title);
    }
  }, [location.state]);

  useEffect(() => {
    if (!id) return;

    let isMounted = true;

    // Fetch Car details
    carService
      .get(id)
      .then((carData) => {
        if (isMounted) {
          setCar(carData);
          setErrorMsg(null);
        }
      })
      .catch(() => {
        if (isMounted) {
          setErrorMsg(
            "Veículo não encontrado ou você não possui permissão para visualizá-lo.",
          );
        }
      })
      .finally(() => {
        if (isMounted) setLoadingCar(false);
      });

    // Fetch Maintenances
    maintenanceService
      .listByCar(id)
      .then((data) => {
        if (isMounted) {
          // Sort descending by date
          const sorted = [...data].sort(
            (a, b) => new Date(b.date).getTime() - new Date(a.date).getTime(),
          );
          setMaintenances(sorted);
        }
      })
      .catch(() => {
        if (isMounted) {
          // non-fatal, fallback to empty list
          setMaintenances([]);
        }
      })
      .finally(() => {
        if (isMounted) setLoadingMaintenances(false);
      });

    // Fetch Fuelings
    fuelService
      .listByCar(id)
      .then((data) => {
        if (isMounted) {
          const sorted = [...data].sort(
            (a, b) => new Date(b.date).getTime() - new Date(a.date).getTime(),
          );
          setFuelings(sorted);
        }
      })
      .catch(() => {
        if (isMounted) {
          setFuelings([]);
        }
      })
      .finally(() => {
        if (isMounted) setLoadingFuelings(false);
      });

    return () => {
      isMounted = false;
    };
  }, [id]);

  const formatBRL = (val: number): string => {
    return new Intl.NumberFormat("pt-BR", {
      style: "currency",
      currency: "BRL",
    }).format(val);
  };

  const totalFuelCost = fuelings.reduce(
    (sum, f) => sum + (f.totalCost || 0),
    0,
  );
  const totalMaintenanceCost = maintenances.reduce(
    (sum, m) => sum + (m.cost || 0),
    0,
  );
  const totalOperationalCost = totalFuelCost + totalMaintenanceCost;

  const getVehicleTypeIcon = (type?: string) => {
    switch (type) {
      case "motorcycles":
        return <TwoWheelerRoundedIcon fontSize="small" />;
      case "trucks":
        return <LocalShippingRoundedIcon fontSize="small" />;
      default:
        return <DirectionsCarRoundedIcon fontSize="small" />;
    }
  };

  const getVehicleTypeLabel = (type?: string) => {
    switch (type) {
      case "motorcycles":
        return "Moto";
      case "trucks":
        return "Caminhão";
      default:
        return "Carro";
    }
  };

  const formatMaintenanceDate = (dateStr?: string): string => {
    if (!dateStr) return "";
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

  const handleConfirmDelete = async () => {
    if (!maintenanceToDelete) return;
    try {
      setDeleting(true);
      setDeleteError(null);
      await maintenanceService.delete(maintenanceToDelete.id);
      setMaintenances((prev) =>
        prev.filter((m) => m.id !== maintenanceToDelete.id),
      );
      setMaintenanceToDelete(null);
      setToastMessage("Manutenção removida com sucesso");
    } catch (err: unknown) {
      const errorObj = err as { response?: { data?: { message?: string } } };
      setDeleteError(
        errorObj.response?.data?.message ||
          "Não foi possível excluir a manutenção. Tente novamente.",
      );
    } finally {
      setDeleting(false);
    }
  };

  const handleConfirmDeleteFueling = async () => {
    if (!fuelingToDelete) return;
    try {
      setDeletingFueling(true);
      setDeleteFuelingError(null);
      await fuelService.delete(fuelingToDelete.id);
      setFuelings((prev) => prev.filter((f) => f.id !== fuelingToDelete.id));
      setFuelingToDelete(null);
      setToastMessage("Abastecimento removido com sucesso");
    } catch (err: unknown) {
      const errorObj = err as { response?: { data?: { message?: string } } };
      setDeleteFuelingError(
        errorObj.response?.data?.message ||
          "Não foi possível excluir o abastecimento. Tente novamente.",
      );
    } finally {
      setDeletingFueling(false);
    }
  };

  const displayError = !id
    ? "Identificador do veículo não fornecido."
    : errorMsg;

  if (displayError) {
    return (
      <Box sx={{ py: 4 }}>
        <Button
          startIcon={<ArrowBackRoundedIcon />}
          onClick={() => navigate("/vehicles")}
          sx={{ mb: 3 }}
        >
          Voltar para Meus Veículos
        </Button>
        <Alert
          severity="error"
          sx={{ borderRadius: 2 }}
          action={
            <Button
              color="inherit"
              size="small"
              onClick={() => navigate("/vehicles")}
            >
              Ir para Veículos
            </Button>
          }
        >
          {displayError}
        </Alert>
      </Box>
    );
  }

  return (
    <Box>
      {/* Top Header Navigation */}
      <Stack
        direction={{ xs: "column", sm: "row" }}
        spacing={2}
        sx={{
          alignItems: { xs: "flex-start", sm: "center" },
          justifyContent: "space-between",
          mb: 3,
        }}
      >
        <Stack direction="row" spacing={1.5} sx={{ alignItems: "center" }}>
          <Button
            variant="outlined"
            size="small"
            startIcon={<ArrowBackRoundedIcon />}
            onClick={() => navigate("/vehicles")}
            sx={{
              borderRadius: 2,
              color: "text.secondary",
              borderColor: "#E2E8F0",
              "&:hover": {
                borderColor: "#CBD5E1",
                bgcolor: "#F8FAFC",
              },
            }}
          >
            Voltar
          </Button>
          {loadingCar ? (
            <Skeleton width={200} height={36} />
          ) : (
            <Typography
              component="h1"
              variant="h5"
              sx={{ fontWeight: 800, color: "text.primary" }}
            >
              {car?.name}
            </Typography>
          )}
        </Stack>

        {activeTab === 0 ? (
          <Button
            variant="contained"
            startIcon={<AddRoundedIcon />}
            onClick={() => {
              setSelectedFueling(null);
              setFuelDialogOpen(true);
            }}
            disabled={loadingCar || !car}
            sx={{ px: 2.5, py: 1.1 }}
          >
            Registrar Abastecimento
          </Button>
        ) : (
          <Button
            variant="contained"
            startIcon={<AddRoundedIcon />}
            onClick={() =>
              navigate(`/vehicles/${car?.id || id}/maintenance/new`)
            }
            disabled={loadingCar || !car}
            sx={{ px: 2.5, py: 1.1 }}
          >
            Registrar Manutenção
          </Button>
        )}
      </Stack>

      {/* Vehicle Overview Card */}
      {loadingCar ? (
        <Card
          elevation={0}
          sx={{
            p: 3,
            mb: 4,
            border: "1px solid #E2E8F0",
            borderRadius: 2.5,
            bgcolor: "background.paper",
          }}
        >
          <Grid container spacing={3}>
            <Grid size={{ xs: 12, md: 4 }}>
              <Skeleton
                variant="rectangular"
                height={220}
                sx={{ borderRadius: 2 }}
              />
            </Grid>
            <Grid size={{ xs: 12, md: 8 }}>
              <Skeleton width="40%" height={32} sx={{ mb: 1 }} />
              <Skeleton width="60%" height={24} sx={{ mb: 2 }} />
              <Skeleton width="100%" height={80} sx={{ borderRadius: 1.5 }} />
            </Grid>
          </Grid>
        </Card>
      ) : car ? (
        <Card
          elevation={0}
          sx={{
            mb: 4,
            border: "1px solid #E2E8F0",
            borderRadius: 2.5,
            bgcolor: "background.paper",
            overflow: "hidden",
          }}
        >
          <Grid container>
            {/* Vehicle Image / Vector Box */}
            <Grid size={{ xs: 12, md: 4 }}>
              <Box
                sx={{
                  position: "relative",
                  width: "100%",
                  height: "100%",
                  minHeight: 220,
                  bgcolor: "#F8FAFC",
                  borderRight: { md: "1px solid #E2E8F0" },
                  borderBottom: { xs: "1px solid #E2E8F0", md: "none" },
                }}
              >
                <VehicleImage
                  imageUrl={car.imageUrl}
                  vehicleType={car.vehicleType}
                  alt={`${car.manufacturer} ${car.model}`}
                  aspectRatio="16 / 10"
                  borderRadius={0}
                />
              </Box>
            </Grid>

            {/* Vehicle Details */}
            <Grid size={{ xs: 12, md: 8 }}>
              <CardContent sx={{ p: { xs: 2.5, sm: 3.5 } }}>
                <Box sx={{ mb: 2 }}>
                  <Stack
                    direction="row"
                    spacing={1}
                    sx={{ alignItems: "center", mb: 0.5 }}
                  >
                    <Typography
                      variant="h5"
                      component="h2"
                      sx={{ fontWeight: 800, color: "text.primary" }}
                    >
                      {car.name}
                    </Typography>
                    <Chip
                      icon={getVehicleTypeIcon(car.vehicleType)}
                      label={getVehicleTypeLabel(car.vehicleType)}
                      size="small"
                      sx={{
                        bgcolor: "#F1F5F9",
                        color: "text.secondary",
                        fontWeight: 600,
                        borderRadius: 1,
                      }}
                    />
                  </Stack>
                  <Typography
                    variant="subtitle1"
                    sx={{ color: "text.secondary", fontWeight: 600 }}
                  >
                    {car.manufacturer} • {car.model}
                  </Typography>
                </Box>

                {/* Specs Badges */}
                <Stack
                  direction="row"
                  spacing={1.5}
                  sx={{ flexWrap: "wrap", gap: 1.5, mb: 3 }}
                >
                  <Chip
                    icon={
                      <CalendarMonthRoundedIcon
                        sx={{ "&&": { fontSize: 18 } }}
                      />
                    }
                    label={`Ano: ${car.yearManufacture}/${car.yearModel}`}
                    sx={{
                      bgcolor: "#F1F5F9",
                      color: "text.primary",
                      fontWeight: 600,
                      borderRadius: 1.5,
                      py: 0.5,
                    }}
                  />
                  <Chip
                    icon={
                      <SpeedRoundedIcon
                        sx={{ "&&": { fontSize: 18, color: "#064E3B" } }}
                      />
                    }
                    label={`${(car.lastMileage || 0).toLocaleString("pt-BR")} km`}
                    sx={{
                      bgcolor: brandColors.mint,
                      color: "#064E3B",
                      fontWeight: 700,
                      borderRadius: 1.5,
                      py: 0.5,
                    }}
                  />
                </Stack>

                <Divider sx={{ my: 2 }} />

                {/* FIPE Section */}
                <Box>
                  <Typography
                    variant="caption"
                    sx={{
                      fontWeight: 700,
                      color: "text.secondary",
                      textTransform: "uppercase",
                      letterSpacing: 0.5,
                      display: "block",
                      mb: 1.5,
                    }}
                  >
                    Informações FIPE & Mercado
                  </Typography>

                  <Grid container spacing={2}>
                    <Grid size={{ xs: 12, sm: 4 }}>
                      <Stack
                        direction="row"
                        spacing={1}
                        sx={{ alignItems: "center" }}
                      >
                        <TagRoundedIcon
                          fontSize="small"
                          sx={{ color: "text.secondary" }}
                        />
                        <Box>
                          <Typography
                            variant="caption"
                            color="text.secondary"
                            sx={{ display: "block" }}
                          >
                            Código FIPE
                          </Typography>
                          <Typography variant="body2" sx={{ fontWeight: 600 }}>
                            {car.fipeCode || "Não vinculado"}
                          </Typography>
                        </Box>
                      </Stack>
                    </Grid>

                    <Grid size={{ xs: 12, sm: 4 }}>
                      <Stack
                        direction="row"
                        spacing={1}
                        sx={{ alignItems: "center" }}
                      >
                        <MonetizationOnRoundedIcon
                          fontSize="small"
                          sx={{ color: "success.main" }}
                        />
                        <Box>
                          <Typography
                            variant="caption"
                            color="text.secondary"
                            sx={{ display: "block" }}
                          >
                            Valor de Referência
                          </Typography>
                          <Typography
                            variant="body2"
                            sx={{ fontWeight: 700, color: "success.main" }}
                          >
                            {car.fipePrice || "Não consultado"}
                          </Typography>
                        </Box>
                      </Stack>
                    </Grid>

                    <Grid size={{ xs: 12, sm: 4 }}>
                      <Stack
                        direction="row"
                        spacing={1}
                        sx={{ alignItems: "center" }}
                      >
                        <LocalGasStationRoundedIcon
                          fontSize="small"
                          sx={{ color: "text.secondary" }}
                        />
                        <Box>
                          <Typography
                            variant="caption"
                            color="text.secondary"
                            sx={{ display: "block" }}
                          >
                            Combustível
                          </Typography>
                          <Typography variant="body2" sx={{ fontWeight: 600 }}>
                            {car.fuel || "Flex"}
                          </Typography>
                        </Box>
                      </Stack>
                    </Grid>
                  </Grid>
                </Box>
              </CardContent>
            </Grid>
          </Grid>
        </Card>
      ) : null}

      {/* Financial KPIs Summary Grid */}
      <Grid container spacing={2} sx={{ mb: 3 }}>
        {/* Card 1: Gastos com Combustível */}
        <Grid size={{ xs: 12, md: 4 }}>
          <Card
            elevation={0}
            sx={{
              border: "1px solid #E2E8F0",
              borderRadius: 2.5,
              p: 2.5,
              bgcolor: "background.paper",
            }}
          >
            <Stack direction="row" spacing={2} sx={{ alignItems: "center" }}>
              <Box
                sx={{
                  width: 48,
                  height: 48,
                  borderRadius: 2,
                  bgcolor: "rgba(2, 132, 199, 0.1)",
                  color: "primary.main",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                }}
              >
                <LocalGasStationRoundedIcon sx={{ fontSize: 28 }} />
              </Box>
              <Box>
                <Typography
                  variant="caption"
                  sx={{
                    color: "text.secondary",
                    fontWeight: 600,
                    textTransform: "uppercase",
                    letterSpacing: 0.5,
                  }}
                >
                  Gastos com Combustível
                </Typography>
                <Typography
                  variant="h6"
                  sx={{ fontWeight: 800, color: "text.primary" }}
                >
                  {formatBRL(totalFuelCost)}
                </Typography>
                <Typography variant="caption" sx={{ color: "text.secondary" }}>
                  {fuelings.length}{" "}
                  {fuelings.length === 1 ? "abastecimento" : "abastecimentos"}
                </Typography>
              </Box>
            </Stack>
          </Card>
        </Grid>

        {/* Card 2: Gastos com Manutenção */}
        <Grid size={{ xs: 12, md: 4 }}>
          <Card
            elevation={0}
            sx={{
              border: "1px solid #E2E8F0",
              borderRadius: 2.5,
              p: 2.5,
              bgcolor: "background.paper",
            }}
          >
            <Stack direction="row" spacing={2} sx={{ alignItems: "center" }}>
              <Box
                sx={{
                  width: 48,
                  height: 48,
                  borderRadius: 2,
                  bgcolor: "rgba(16, 185, 129, 0.1)",
                  color: "#059669",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                }}
              >
                <BuildCircleRoundedIcon sx={{ fontSize: 28 }} />
              </Box>
              <Box>
                <Typography
                  variant="caption"
                  sx={{
                    color: "text.secondary",
                    fontWeight: 600,
                    textTransform: "uppercase",
                    letterSpacing: 0.5,
                  }}
                >
                  Gastos com Manutenção
                </Typography>
                <Typography
                  variant="h6"
                  sx={{ fontWeight: 800, color: "text.primary" }}
                >
                  {formatBRL(totalMaintenanceCost)}
                </Typography>
                <Typography variant="caption" sx={{ color: "text.secondary" }}>
                  {maintenances.length}{" "}
                  {maintenances.length === 1 ? "serviço" : "serviços"}
                </Typography>
              </Box>
            </Stack>
          </Card>
        </Grid>

        {/* Card 3: Custo Operacional Total */}
        <Grid size={{ xs: 12, md: 4 }}>
          <Card
            elevation={0}
            sx={{
              border: "1px solid #0284C7",
              borderRadius: 2.5,
              p: 2.5,
              bgcolor: "rgba(2, 132, 199, 0.04)",
            }}
          >
            <Stack direction="row" spacing={2} sx={{ alignItems: "center" }}>
              <Box
                sx={{
                  width: 48,
                  height: 48,
                  borderRadius: 2,
                  bgcolor: "primary.main",
                  color: "#FFFFFF",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                }}
              >
                <MonetizationOnRoundedIcon sx={{ fontSize: 28 }} />
              </Box>
              <Box>
                <Typography
                  variant="caption"
                  sx={{
                    color: "primary.main",
                    fontWeight: 700,
                    textTransform: "uppercase",
                    letterSpacing: 0.5,
                  }}
                >
                  Custo Operacional Total
                </Typography>
                <Typography
                  variant="h6"
                  sx={{ fontWeight: 800, color: "primary.main" }}
                >
                  {formatBRL(totalOperationalCost)}
                </Typography>
                <Typography variant="caption" sx={{ color: "text.secondary" }}>
                  Combustível + Manutenções
                </Typography>
              </Box>
            </Stack>
          </Card>
        </Grid>
      </Grid>

      {/* Navigation Tabs */}
      <Box sx={{ borderBottom: 1, borderColor: "divider", mb: 3 }}>
        <Tabs
          value={activeTab}
          onChange={(_, val) => setActiveTab(val)}
          textColor="primary"
          indicatorColor="primary"
        >
          <Tab
            icon={<LocalGasStationRoundedIcon sx={{ fontSize: 20 }} />}
            iconPosition="start"
            label={
              <Stack direction="row" spacing={1} sx={{ alignItems: "center" }}>
                <span>Abastecimentos & Combustível</span>
                {!loadingFuelings && (
                  <Chip
                    label={fuelings.length}
                    size="small"
                    sx={{
                      height: 20,
                      fontSize: "0.75rem",
                      fontWeight: 700,
                      bgcolor: activeTab === 0 ? "primary.main" : "#E2E8F0",
                      color: activeTab === 0 ? "#FFFFFF" : "text.secondary",
                    }}
                  />
                )}
              </Stack>
            }
            sx={{ fontWeight: 700, textTransform: "none" }}
          />
          <Tab
            icon={<BuildCircleRoundedIcon sx={{ fontSize: 20 }} />}
            iconPosition="start"
            label={
              <Stack direction="row" spacing={1} sx={{ alignItems: "center" }}>
                <span>Manutenções & Revisões</span>
                {!loadingMaintenances && (
                  <Chip
                    label={maintenances.length}
                    size="small"
                    sx={{
                      height: 20,
                      fontSize: "0.75rem",
                      fontWeight: 700,
                      bgcolor: activeTab === 1 ? "primary.main" : "#E2E8F0",
                      color: activeTab === 1 ? "#FFFFFF" : "text.secondary",
                    }}
                  />
                )}
              </Stack>
            }
            sx={{ fontWeight: 700, textTransform: "none" }}
          />
        </Tabs>
      </Box>

      {/* Tab 0: Abastecimentos & Combustível */}
      {activeTab === 0 && (
        <Box>
          <Stack
            direction="row"
            sx={{
              alignItems: "center",
              justifyContent: "space-between",
              mb: 2.5,
            }}
          >
            <Stack direction="row" spacing={1} sx={{ alignItems: "center" }}>
              <Typography
                component="h2"
                variant="h6"
                sx={{ fontWeight: 800, color: "text.primary" }}
              >
                Histórico de Abastecimentos
              </Typography>
              {!loadingFuelings && (
                <Chip
                  label={fuelings.length}
                  size="small"
                  sx={{
                    bgcolor: "primary.main",
                    color: "#FFFFFF",
                    fontWeight: 700,
                    borderRadius: 1,
                  }}
                />
              )}
            </Stack>
            <Button
              variant="contained"
              size="small"
              startIcon={<AddRoundedIcon />}
              onClick={() => {
                setSelectedFueling(null);
                setFuelDialogOpen(true);
              }}
              disabled={!car}
              sx={{ borderRadius: 2 }}
            >
              Registrar Abastecimento
            </Button>
          </Stack>

          {loadingFuelings ? (
            <Stack spacing={2}>
              {[1, 2].map((i) => (
                <Card
                  key={i}
                  elevation={0}
                  sx={{
                    p: 2.5,
                    border: "1px solid #E2E8F0",
                    borderRadius: 2.5,
                    bgcolor: "background.paper",
                  }}
                >
                  <Skeleton width="40%" height={28} sx={{ mb: 1 }} />
                  <Skeleton width="20%" height={20} />
                </Card>
              ))}
            </Stack>
          ) : fuelings.length === 0 ? (
            <Card
              elevation={0}
              sx={{
                py: 7,
                px: 3,
                textAlign: "center",
                border: "1px dashed #CBD5E1",
                bgcolor: "background.paper",
                borderRadius: 2.5,
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
                <LocalGasStationRoundedIcon sx={{ fontSize: 32 }} />
              </Box>
              <Typography
                component="h3"
                variant="h6"
                sx={{ fontWeight: 700, mb: 0.5 }}
              >
                Nenhum abastecimento registrado para este veículo
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
                Registre abastecimentos de combustível para acompanhar o
                consumo, valor pago e controlar seus custos operacionais.
              </Typography>
              <Button
                variant="contained"
                startIcon={<AddRoundedIcon />}
                onClick={() => {
                  setSelectedFueling(null);
                  setFuelDialogOpen(true);
                }}
                disabled={!car}
              >
                Registrar Primeiro Abastecimento
              </Button>
            </Card>
          ) : (
            <Stack spacing={2}>
              {fuelings.map((f) => (
                <FuelCard
                  key={f.id}
                  fueling={f}
                  onEdit={(fuel) => {
                    setSelectedFueling(fuel);
                    setFuelDialogOpen(true);
                  }}
                  onDelete={(fuelId) => {
                    const target = fuelings.find((item) => item.id === fuelId);
                    if (target) {
                      setFuelingToDelete(target);
                      setDeleteFuelingError(null);
                    }
                  }}
                />
              ))}
            </Stack>
          )}
        </Box>
      )}

      {/* Tab 1: Manutenções & Revisões */}
      {activeTab === 1 && (
        <Box>
          <Stack
            direction="row"
            sx={{
              alignItems: "center",
              justifyContent: "space-between",
              mb: 2.5,
            }}
          >
            <Stack direction="row" spacing={1} sx={{ alignItems: "center" }}>
              <Typography
                component="h2"
                variant="h6"
                sx={{ fontWeight: 800, color: "text.primary" }}
              >
                Histórico de Manutenções
              </Typography>
              {!loadingMaintenances && (
                <Chip
                  label={maintenances.length}
                  size="small"
                  sx={{
                    bgcolor: "primary.main",
                    color: "#FFFFFF",
                    fontWeight: 700,
                    borderRadius: 1,
                  }}
                />
              )}
            </Stack>
            <Button
              variant="contained"
              size="small"
              startIcon={<AddRoundedIcon />}
              onClick={() =>
                navigate(`/vehicles/${car?.id || id}/maintenance/new`)
              }
              disabled={!car}
              sx={{ borderRadius: 2 }}
            >
              Registrar Manutenção
            </Button>
          </Stack>

          {loadingMaintenances ? (
            <Stack spacing={2}>
              {[1, 2].map((i) => (
                <Card
                  key={i}
                  elevation={0}
                  sx={{
                    p: 2.5,
                    border: "1px solid #E2E8F0",
                    borderRadius: 2.5,
                    bgcolor: "background.paper",
                  }}
                >
                  <Skeleton width="40%" height={28} sx={{ mb: 1 }} />
                  <Skeleton width="20%" height={20} />
                </Card>
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
                borderRadius: 2.5,
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
                Nenhuma manutenção registrada para este veículo
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
                Registre revisões, trocas de óleo, pastilhas de freio ou
                substituição de pneus para manter o histórico do seu veículo em
                dia.
              </Typography>
              <Button
                variant="contained"
                startIcon={<AddRoundedIcon />}
                onClick={() =>
                  navigate(`/vehicles/${car?.id || id}/maintenance/new`)
                }
                disabled={!car}
              >
                Registrar Primeira Manutenção
              </Button>
            </Card>
          ) : (
            <Stack spacing={2}>
              {maintenances.map((maint) => (
                <MaintenanceCard
                  key={maint.id}
                  maintenance={maint}
                  onEdit={(m) =>
                    navigate(
                      `/vehicles/${car?.id || id}/maintenance/${m.id}/edit`,
                    )
                  }
                  onDelete={(maintId) => {
                    const target = maintenances.find((m) => m.id === maintId);
                    if (target) {
                      setMaintenanceToDelete(target);
                      setDeleteError(null);
                    }
                  }}
                />
              ))}
            </Stack>
          )}
        </Box>
      )}

      {/* Add / Edit Fueling Dialog */}
      <AddFuelingDialog
        open={fuelDialogOpen}
        onClose={() => {
          setFuelDialogOpen(false);
          setSelectedFueling(null);
        }}
        initialData={selectedFueling}
        carId={car?.id || id || ""}
        onSuccess={(fueling, isEdit) => {
          if (isEdit) {
            setFuelings((prev) =>
              prev.map((f) => (f.id === fueling.id ? fueling : f)),
            );
            setToastMessage("Abastecimento atualizado com sucesso");
          } else {
            setFuelings((prev) => [fueling, ...prev]);
            setToastMessage("Abastecimento registrado com sucesso");
          }
        }}
      />

      {/* Delete Fueling Confirmation Dialog */}
      <Dialog
        open={Boolean(fuelingToDelete)}
        onClose={() => !deletingFueling && setFuelingToDelete(null)}
        maxWidth="xs"
        fullWidth
        slotProps={{
          paper: {
            sx: { borderRadius: 3, p: 1 },
          },
        }}
      >
        <DialogTitle sx={{ fontWeight: 800, pb: 1 }}>
          Excluir Abastecimento
        </DialogTitle>
        <DialogContent>
          <DialogContentText sx={{ color: "text.primary", mb: 2 }}>
            Tem certeza que deseja excluir o abastecimento de{" "}
            <strong>"{fuelingToDelete?.fuelType}"</strong> no valor de{" "}
            <strong>{formatBRL(fuelingToDelete?.totalCost || 0)}</strong>? Esta
            ação é irreversível e removerá o registro do histórico do veículo.
          </DialogContentText>
          {deleteFuelingError && (
            <Alert severity="error" sx={{ borderRadius: 2 }}>
              {deleteFuelingError}
            </Alert>
          )}
        </DialogContent>
        <DialogActions sx={{ px: 3, pb: 2 }}>
          <Button
            variant="outlined"
            color="inherit"
            onClick={() => setFuelingToDelete(null)}
            disabled={deletingFueling}
            sx={{ borderRadius: 2 }}
          >
            Cancelar
          </Button>
          <Button
            variant="contained"
            color="error"
            onClick={handleConfirmDeleteFueling}
            disabled={deletingFueling}
            startIcon={
              deletingFueling ? (
                <CircularProgress size={18} color="inherit" />
              ) : (
                <DeleteOutlineRoundedIcon />
              )
            }
            sx={{ borderRadius: 2 }}
          >
            Excluir
          </Button>
        </DialogActions>
      </Dialog>

      {/* Delete Maintenance Confirmation Dialog */}
      <Dialog
        open={Boolean(maintenanceToDelete)}
        onClose={() => !deleting && setMaintenanceToDelete(null)}
        maxWidth="xs"
        fullWidth
        slotProps={{
          paper: {
            sx: { borderRadius: 3, p: 1 },
          },
        }}
      >
        <DialogTitle sx={{ fontWeight: 800, pb: 1 }}>
          Excluir Manutenção
        </DialogTitle>
        <DialogContent>
          <DialogContentText sx={{ color: "text.primary", mb: 2 }}>
            Tem certeza que deseja excluir a manutenção{" "}
            <strong>"{maintenanceToDelete?.title}"</strong> realizada em{" "}
            <strong>{formatMaintenanceDate(maintenanceToDelete?.date)}</strong>?
            Esta ação é irreversível e removerá o registro do histórico do
            veículo.
          </DialogContentText>
          {deleteError && (
            <Alert severity="error" sx={{ borderRadius: 2 }}>
              {deleteError}
            </Alert>
          )}
        </DialogContent>
        <DialogActions sx={{ px: 3, pb: 2 }}>
          <Button
            variant="outlined"
            color="inherit"
            onClick={() => setMaintenanceToDelete(null)}
            disabled={deleting}
            sx={{ borderRadius: 2 }}
          >
            Cancelar
          </Button>
          <Button
            variant="contained"
            color="error"
            onClick={handleConfirmDelete}
            disabled={deleting}
            startIcon={
              deleting ? (
                <CircularProgress size={18} color="inherit" />
              ) : (
                <DeleteOutlineRoundedIcon />
              )
            }
            sx={{ borderRadius: 2 }}
          >
            Excluir
          </Button>
        </DialogActions>
      </Dialog>

      {/* Success Toast */}
      <Snackbar
        open={Boolean(toastMessage)}
        autoHideDuration={4000}
        onClose={() => setToastMessage(null)}
        anchorOrigin={{ vertical: "bottom", horizontal: "center" }}
      >
        <Alert
          onClose={() => setToastMessage(null)}
          severity="success"
          variant="filled"
          sx={{ width: "100%", borderRadius: 2 }}
        >
          {toastMessage}
        </Alert>
      </Snackbar>
    </Box>
  );
}
