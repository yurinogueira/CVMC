import { useEffect, useState, lazy, Suspense } from "react";
import {
  Box,
  Typography,
  Grid,
  Button,
  Stack,
  TextField,
  InputAdornment,
  Skeleton,
  Alert,
  Card,
} from "@mui/material";
import AddRoundedIcon from "@mui/icons-material/AddRounded";
import SearchRoundedIcon from "@mui/icons-material/SearchRounded";
import DirectionsCarFilledRoundedIcon from "@mui/icons-material/DirectionsCarFilledRounded";
import { carService } from "../services/car.service";
import { Car } from "../types/car.types";
import { VehicleCard } from "../components/VehicleCard";
import { useDocumentTitle } from "../../shared";

const AddCarDialog = lazy(() =>
  import("../components/AddCarDialog").then((m) => ({
    default: m.AddCarDialog,
  })),
);

export function VehiclesPage() {
  useDocumentTitle("Meus Veículos");
  const [cars, setCars] = useState<Car[]>([]);
  const [search, setSearch] = useState("");
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
      setErrorMsg("Não foi possível carregar a lista de veículos.");
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
          setErrorMsg("Não foi possível carregar a lista de veículos.");
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

  const filteredCars = cars.filter(
    (car) =>
      car.name.toLowerCase().includes(search.toLowerCase()) ||
      car.manufacturer.toLowerCase().includes(search.toLowerCase()) ||
      car.model.toLowerCase().includes(search.toLowerCase()),
  );

  return (
    <Box>
      {/* Header Actions */}
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
            Meus Veículos
          </Typography>
          <Typography variant="body2" sx={{ color: "text.secondary" }}>
            Gerencie todos os veículos cadastrados na sua conta.
          </Typography>
        </Box>

        <Button
          variant="contained"
          startIcon={<AddRoundedIcon />}
          onClick={() => setOpenAddDialog(true)}
          sx={{ px: 2.5, py: 1.1 }}
        >
          Adicionar Veículo
        </Button>
      </Stack>

      {/* Search Input */}
      <Box sx={{ mb: 3, maxWidth: 400 }}>
        <TextField
          fullWidth
          size="small"
          placeholder="Buscar por apelido, fabricante ou modelo..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          slotProps={{
            input: {
              startAdornment: (
                <InputAdornment position="start">
                  <SearchRoundedIcon
                    sx={{ color: "text.secondary", fontSize: 20 }}
                  />
                </InputAdornment>
              ),
            },
          }}
        />
      </Box>

      {/* Error Alert */}
      {errorMsg && (
        <Alert
          severity="error"
          action={
            <Button color="inherit" size="small" onClick={fetchCars}>
              Tentar novamente
            </Button>
          }
          sx={{ mb: 3 }}
        >
          {errorMsg}
        </Alert>
      )}

      {/* Content */}
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
      ) : filteredCars.length === 0 ? (
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
              bgcolor: "rgba(76, 146, 252, 0.1)",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              mx: "auto",
              mb: 2,
              color: "primary.main",
            }}
          >
            <DirectionsCarFilledRoundedIcon sx={{ fontSize: 32 }} />
          </Box>
          <Typography
            component="h2"
            variant="h6"
            sx={{ fontWeight: 700, mb: 0.5 }}
          >
            {search
              ? "Nenhum veículo encontrado para a busca"
              : "Você ainda não possui veículos cadastrados"}
          </Typography>
          <Typography variant="body2" sx={{ color: "text.secondary", mb: 2.5 }}>
            {search
              ? "Tente pesquisar com outros termos ou limpe o campo de busca."
              : "Cadastre seu primeiro carro para começar."}
          </Typography>
          {!search && (
            <Button
              variant="contained"
              startIcon={<AddRoundedIcon />}
              onClick={() => setOpenAddDialog(true)}
            >
              Novo Veículo
            </Button>
          )}
        </Card>
      ) : (
        <Grid container spacing={3}>
          {filteredCars.map((car) => (
            <Grid size={{ xs: 12, sm: 6, md: 4 }} key={car.id}>
              <VehicleCard car={car} onDelete={handleDeleteCar} />
            </Grid>
          ))}
        </Grid>
      )}

      {/* Add Dialog (Lazy Loaded) */}
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
