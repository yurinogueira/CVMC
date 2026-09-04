import { useEffect, useState, lazy, Suspense } from "react";
import { Link as RouterLink, useNavigate } from "react-router-dom";
import {
  Box,
  Typography,
  Grid,
  Card,
  CardContent,
  Button,
  Stack,
  Skeleton,
  Alert,
} from "@mui/material";
import DirectionsCarFilledRoundedIcon from "@mui/icons-material/DirectionsCarFilledRounded";
import BuildCircleRoundedIcon from "@mui/icons-material/BuildCircleRounded";
import SpeedRoundedIcon from "@mui/icons-material/SpeedRounded";
import AddRoundedIcon from "@mui/icons-material/AddRounded";
import ArrowForwardRoundedIcon from "@mui/icons-material/ArrowForwardRounded";
import { useAuthStore } from "../../auth/state/auth.store";
import { carService } from "../../cars/services/car.service";
import { Car } from "../../cars/types/car.types";
import { VehicleCard } from "../../cars/components/VehicleCard";
import { useDocumentTitle } from "../../shared";
import { brandColors } from "../../../styles/theme";

const AddCarDialog = lazy(() =>
  import("../../cars/components/AddCarDialog").then((m) => ({
    default: m.AddCarDialog,
  })),
);

export function DashboardPage() {
  useDocumentTitle("Dashboard");
  const navigate = useNavigate();
  const user = useAuthStore((state) => state.user);
  const [cars, setCars] = useState<Car[]>([]);
  const [loading, setLoading] = useState(true);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [openAddDialog, setOpenAddDialog] = useState(false);

  const fetchCars = async () => {
    try {
      setLoading(true);
      setErrorMsg(null);
      const data = await carService.list();
      setCars(data);
    } catch {
      setErrorMsg(
        "Não foi possível carregar os dados de veículos. O servidor pode estar em inicialização.",
      );
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    let isMounted = true;
    carService
      .list()
      .then((data) => {
        if (isMounted) setCars(data);
      })
      .catch(() => {
        if (isMounted) {
          setErrorMsg(
            "Não foi possível carregar os dados de veículos. O servidor pode estar em inicialização.",
          );
        }
      })
      .finally(() => {
        if (isMounted) setLoading(false);
      });

    return () => {
      isMounted = false;
    };
  }, []);

  const handleCarCreated = (newCar: Car) => {
    setCars((prev) => [newCar, ...prev]);
  };

  const handleDeleteCar = async (carId: string) => {
    try {
      await carService.delete(carId);
      setCars((prev) => prev.filter((c) => c.id !== carId));
    } catch {
      // ignore or show toast
    }
  };

  const totalMileage = cars.reduce(
    (acc, car) => acc + (car.lastMileage || 0),
    0,
  );

  return (
    <Box>
      {/* Welcome Banner */}
      <Box sx={{ mb: 4 }}>
        <Typography
          component="h1"
          variant="h4"
          sx={{ fontWeight: 800, mb: 0.5, color: "text.primary" }}
        >
          Olá, {user?.name?.split(" ")[0] || "Motorista"} 👋
        </Typography>
        <Typography variant="body1" sx={{ color: "text.secondary" }}>
          Aqui está o resumo da sua frota e atividades recentes.
        </Typography>
      </Box>

      {/* KPI Cards Grid */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        {/* KPI 1: Veículos */}
        <Grid size={{ xs: 12, sm: 6, md: 4 }}>
          <Card
            elevation={0}
            sx={{
              p: 1,
              bgcolor: "background.paper",
              border: "1px solid #E2E8F0",
              transition: "transform 0.2s",
              "&:hover": { transform: "translateY(-2px)" },
            }}
          >
            <CardContent>
              <Stack
                direction="row"
                sx={{
                  alignItems: "center",
                  justifyContent: "space-between",
                }}
              >
                <Box>
                  <Typography
                    component="h3"
                    variant="caption"
                    sx={{
                      color: "text.secondary",
                      fontWeight: 600,
                      textTransform: "uppercase",
                    }}
                  >
                    Veículos Cadastrados
                  </Typography>
                  <Typography
                    variant="h4"
                    sx={{ fontWeight: 800, my: 0.5, color: "text.primary" }}
                  >
                    {loading ? (
                      <Skeleton
                        width={48}
                        height={42}
                        sx={{ display: "inline-block" }}
                      />
                    ) : (
                      cars.length
                    )}
                  </Typography>
                  <Typography
                    variant="body2"
                    sx={{ color: "text.secondary", fontSize: "0.85rem" }}
                  >
                    {cars.length === 1
                      ? "1 veículo monitorado"
                      : `${cars.length} veículos monitorados`}
                  </Typography>
                </Box>
                <Box
                  sx={{
                    width: 52,
                    height: 52,
                    borderRadius: 2.5,
                    bgcolor: "rgba(2, 132, 199, 0.1)",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    color: "primary.main",
                  }}
                >
                  <DirectionsCarFilledRoundedIcon sx={{ fontSize: 28 }} />
                </Box>
              </Stack>
            </CardContent>
          </Card>
        </Grid>

        {/* KPI 2: Manutenções */}
        <Grid size={{ xs: 12, sm: 6, md: 4 }}>
          <Card
            elevation={0}
            sx={{
              p: 1,
              bgcolor: "background.paper",
              border: "1px solid #E2E8F0",
              transition: "transform 0.2s",
              "&:hover": { transform: "translateY(-2px)" },
            }}
          >
            <CardContent>
              <Stack
                direction="row"
                sx={{
                  alignItems: "center",
                  justifyContent: "space-between",
                }}
              >
                <Box>
                  <Typography
                    component="h3"
                    variant="caption"
                    sx={{
                      color: "text.secondary",
                      fontWeight: 600,
                      textTransform: "uppercase",
                    }}
                  >
                    Status da Frota
                  </Typography>
                  <Typography
                    variant="h4"
                    sx={{ fontWeight: 800, my: 0.5, color: "#16A34A" }}
                  >
                    {loading ? (
                      <Skeleton
                        width={72}
                        height={42}
                        sx={{ display: "inline-block" }}
                      />
                    ) : cars.length > 0 ? (
                      "Em dia"
                    ) : (
                      "--"
                    )}
                  </Typography>
                  <Typography
                    variant="body2"
                    sx={{ color: "text.secondary", fontSize: "0.85rem" }}
                  >
                    Sem alertas pendentes
                  </Typography>
                </Box>
                <Box
                  sx={{
                    width: 52,
                    height: 52,
                    borderRadius: 2.5,
                    bgcolor: "rgba(22, 163, 74, 0.12)",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    color: "#16A34A",
                  }}
                >
                  <BuildCircleRoundedIcon sx={{ fontSize: 28 }} />
                </Box>
              </Stack>
            </CardContent>
          </Card>
        </Grid>

        {/* KPI 3: Km Total */}
        <Grid size={{ xs: 12, sm: 6, md: 4 }}>
          <Card
            elevation={0}
            sx={{
              p: 1,
              bgcolor: "background.paper",
              border: "1px solid #E2E8F0",
              transition: "transform 0.2s",
              "&:hover": { transform: "translateY(-2px)" },
            }}
          >
            <CardContent>
              <Stack
                direction="row"
                sx={{
                  alignItems: "center",
                  justifyContent: "space-between",
                }}
              >
                <Box>
                  <Typography
                    component="h3"
                    variant="caption"
                    sx={{
                      color: "text.secondary",
                      fontWeight: 600,
                      textTransform: "uppercase",
                    }}
                  >
                    Quilometragem Total
                  </Typography>
                  <Typography
                    variant="h4"
                    sx={{ fontWeight: 800, my: 0.5, color: "text.primary" }}
                  >
                    {loading ? (
                      <Skeleton
                        width={84}
                        height={42}
                        sx={{ display: "inline-block" }}
                      />
                    ) : (
                      totalMileage.toLocaleString("pt-BR")
                    )}
                  </Typography>
                  <Typography
                    variant="body2"
                    sx={{ color: "text.secondary", fontSize: "0.85rem" }}
                  >
                    Km acumulados registrados
                  </Typography>
                </Box>
                <Box
                  sx={{
                    width: 52,
                    height: 52,
                    borderRadius: 2.5,
                    bgcolor: "rgba(2, 132, 199, 0.1)",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    color: "primary.main",
                  }}
                >
                  <SpeedRoundedIcon sx={{ fontSize: 28 }} />
                </Box>
              </Stack>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Quick Action CTA Banner */}
      <Card
        elevation={0}
        sx={{
          mb: 4,
          background: brandColors.gradient,
          color: "#FFFFFF",
          p: { xs: 2.5, sm: 3.5 },
          position: "relative",
          overflow: "hidden",
        }}
      >
        <Stack
          direction={{ xs: "column", sm: "row" }}
          spacing={2}
          sx={{
            alignItems: { xs: "flex-start", sm: "center" },
            justifyContent: "space-between",
            position: "relative",
            zIndex: 1,
          }}
        >
          <Box>
            <Typography
              component="h2"
              variant="h6"
              sx={{ fontWeight: 700, color: "#FFFFFF", mb: 0.5 }}
            >
              Mantenha o histórico dos seus carros sempre atualizado
            </Typography>
            <Typography
              variant="body2"
              sx={{ color: "rgba(255, 255, 255, 0.95)", maxWidth: 600 }}
            >
              Cadastre novos veículos e registre manutenções preventivas para
              evitar surpresas e valorizar o seu patrimônio.
            </Typography>
          </Box>

          <Button
            variant="contained"
            onClick={() => setOpenAddDialog(true)}
            startIcon={<AddRoundedIcon />}
            sx={{
              bgcolor: "#FFFFFF",
              color: "#0369A1",
              fontWeight: 700,
              px: 3,
              py: 1.2,
              whiteSpace: "nowrap",
              "&:hover": {
                bgcolor: "#F1F5F9",
                color: "#0C4A6E",
              },
            }}
          >
            Novo Veículo
          </Button>
        </Stack>
      </Card>

      {/* Section Header: Meus Veículos */}
      <Stack
        direction="row"
        sx={{
          alignItems: "center",
          justifyContent: "space-between",
          mb: 2.5,
        }}
      >
        <Typography
          component="h2"
          variant="h6"
          sx={{ fontWeight: 700, color: "text.primary" }}
        >
          Meus Veículos
        </Typography>

        {cars.length > 0 && (
          <Button
            component={RouterLink}
            to="/vehicles"
            endIcon={<ArrowForwardRoundedIcon />}
            sx={{ color: "primary.dark", fontWeight: 700 }}
          >
            Ver todos
          </Button>
        )}
      </Stack>

      {/* Error state */}
      {errorMsg && (
        <Alert
          severity="warning"
          action={
            <Button color="inherit" size="small" onClick={fetchCars}>
              Recarregar
            </Button>
          }
          sx={{ mb: 3 }}
        >
          {errorMsg}
        </Alert>
      )}

      {/* Loading state */}
      {loading ? (
        <Grid container spacing={3}>
          {[1, 2, 3].map((i) => (
            <Grid size={{ xs: 12, sm: 6, md: 4 }} key={i}>
              <Card
                elevation={0}
                sx={{
                  p: 2.5,
                  border: "1px solid #E2E8F0",
                  bgcolor: "background.paper",
                }}
              >
                <Skeleton variant="text" width="60%" height={32} />
                <Skeleton
                  variant="text"
                  width="40%"
                  height={24}
                  sx={{ mb: 2 }}
                />
                <Skeleton
                  variant="rectangular"
                  height={80}
                  sx={{ borderRadius: 1.5 }}
                />
              </Card>
            </Grid>
          ))}
        </Grid>
      ) : cars.length === 0 ? (
        /* Modern Empty State */
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
              width: 64,
              height: 64,
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
            <DirectionsCarFilledRoundedIcon sx={{ fontSize: 36 }} />
          </Box>
          <Typography
            component="h3"
            variant="h6"
            sx={{ fontWeight: 700, mb: 1 }}
          >
            Nenhum veículo cadastrado ainda
          </Typography>
          <Typography
            variant="body2"
            sx={{ color: "text.secondary", maxWidth: 420, mx: "auto", mb: 3 }}
          >
            Adicione o seu primeiro veículo para começar a registrar histórico
            de revisões, troca de óleo e manutenções.
          </Typography>
          <Button
            variant="contained"
            startIcon={<AddRoundedIcon />}
            onClick={() => setOpenAddDialog(true)}
            sx={{
              px: 3,
              py: 1.2,
              bgcolor: "primary.main",
              color: "#FFFFFF",
              fontWeight: 700,
              "&:hover": {
                bgcolor: "primary.dark",
              },
            }}
          >
            Cadastrar primeiro veículo
          </Button>
        </Card>
      ) : (
        /* Vehicles Grid */
        <Grid container spacing={3}>
          {cars.map((car) => (
            <Grid size={{ xs: 12, sm: 6, md: 4 }} key={car.id}>
              <VehicleCard
                car={car}
                onDelete={handleDeleteCar}
                onClick={(c) => navigate(`/vehicles/${c.id}`)}
              />
            </Grid>
          ))}
        </Grid>
      )}

      {/* Add Car Modal (Lazy Loaded) */}
      {openAddDialog && (
        <Suspense fallback={null}>
          <AddCarDialog
            open={openAddDialog}
            onClose={() => setOpenAddDialog(false)}
            onCarCreated={handleCarCreated}
          />
        </Suspense>
      )}
    </Box>
  );
}
