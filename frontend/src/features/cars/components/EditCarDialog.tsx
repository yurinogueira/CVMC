import React, { useState } from "react";

import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  Stack,
  Alert,
  CircularProgress,
  Grid,
  Typography,
  Box,
  Paper,
  Chip,
} from "@mui/material";
import DirectionsCarFilledRoundedIcon from "@mui/icons-material/DirectionsCarFilledRounded";
import MonetizationOnRoundedIcon from "@mui/icons-material/MonetizationOnRounded";
import LocalGasStationRoundedIcon from "@mui/icons-material/LocalGasStationRounded";
import TagRoundedIcon from "@mui/icons-material/TagRounded";
import { ImageUploadField } from "./ImageUploadField";

import { carService } from "../services/car.service";
import { Car } from "../types/car.types";

interface EditCarDialogProps {
  open: boolean;
  car: Car | null;
  onClose: () => void;
  onCarUpdated: (car: Car) => void;
}

interface EditCarFormProps {
  car: Car;
  onClose: () => void;
  onCarUpdated: (car: Car) => void;
}

function EditCarForm({ car, onClose, onCarUpdated }: EditCarFormProps) {
  const currentYear = new Date().getFullYear();

  const [name, setName] = useState(car.name || "");
  const [lastMileage, setLastMileage] = useState<number>(car.lastMileage || 0);
  const [yearManufacture, setYearManufacture] = useState<number>(
    car.yearManufacture || currentYear,
  );
  const [yearModel, setYearModel] = useState<number>(
    car.yearModel || currentYear,
  );
  const [imageUrl, setImageUrl] = useState(car.imageUrl || "");
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const handleSubmit = async (e?: React.FormEvent) => {
    if (e) e.preventDefault();
    setErrorMsg(null);

    if (!name.trim()) {
      setErrorMsg("O apelido/identificador do veículo é obrigatório.");
      return;
    }

    if (Number(lastMileage) < 0) {
      setErrorMsg("A quilometragem não pode ser negativa.");
      return;
    }

    if (
      Number(yearManufacture) < 1900 ||
      Number(yearManufacture) > currentYear + 1 ||
      Number(yearModel) < 1900 ||
      Number(yearModel) > currentYear + 2
    ) {
      setErrorMsg("Informe anos de fabricação e modelo válidos.");
      return;
    }

    try {
      setLoading(true);
      const updatedCar = await carService.update(car.id, {
        name: name.trim(),
        manufacturer: car.manufacturer,
        model: car.model,
        yearManufacture: Number(yearManufacture),
        yearModel: Number(yearModel),
        lastMileage: Number(lastMileage) || 0,
        vehicleType: car.vehicleType,
        imageUrl: imageUrl.trim() || undefined,
        fipeCode: car.fipeCode,
        fipePrice: car.fipePrice,
        fuel: car.fuel,
      });

      onCarUpdated(updatedCar);
      onClose();
    } catch (err: unknown) {
      const errorObj = err as { response?: { data?: { message?: string } } };
      setErrorMsg(
        errorObj.response?.data?.message ||
          "Não foi possível atualizar o veículo. Tente novamente.",
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <DialogTitle sx={{ fontWeight: 700, pb: 1 }}>Editar Veículo</DialogTitle>

      <DialogContent dividers sx={{ borderColor: "#F1F5F9" }}>
        {errorMsg && (
          <Alert
            severity="error"
            onClose={() => setErrorMsg(null)}
            sx={{ mb: 2.5, borderRadius: 1.5 }}
          >
            {errorMsg}
          </Alert>
        )}

        <Stack spacing={3} sx={{ mt: 1 }}>
          {/* Resumo do Veículo FIPE */}
          <Paper
            elevation={0}
            sx={{
              p: 2,
              borderRadius: 2,
              bgcolor: "action.hover",
              border: "1px solid",
              borderColor: "divider",
            }}
          >
            <Stack
              direction="row"
              spacing={1}
              sx={{ mb: 1.5, alignItems: "center" }}
            >
              <DirectionsCarFilledRoundedIcon
                color="primary"
                fontSize="small"
              />
              <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>
                {car.manufacturer} • {car.model}
              </Typography>
              {car.fuel && (
                <Chip
                  label={car.fuel}
                  size="small"
                  variant="outlined"
                  sx={{ fontSize: "0.75rem", ml: "auto" }}
                />
              )}
            </Stack>

            <Grid container spacing={2}>
              {car.fipeCode && (
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
                        {car.fipeCode}
                      </Typography>
                    </Box>
                  </Stack>
                </Grid>
              )}

              {car.fipePrice && (
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
                        Preço FIPE Referência
                      </Typography>
                      <Typography
                        variant="body2"
                        sx={{ fontWeight: 700, color: "success.main" }}
                      >
                        {car.fipePrice}
                      </Typography>
                    </Box>
                  </Stack>
                </Grid>
              )}

              {car.fuel && (
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
                        {car.fuel}
                      </Typography>
                    </Box>
                  </Stack>
                </Grid>
              )}
            </Grid>
          </Paper>

          {/* Campos Editáveis */}
          <Box>
            <Typography
              variant="caption"
              sx={{
                fontWeight: 600,
                color: "text.secondary",
                textTransform: "uppercase",
                letterSpacing: 0.5,
                display: "block",
                mb: 1.5,
              }}
            >
              Informações Editáveis
            </Typography>

            <Grid container spacing={2}>
              <Grid size={{ xs: 12, sm: 6 }}>
                <TextField
                  fullWidth
                  label="Apelido / Identificador"
                  placeholder="Ex: Meu Carro, Carro da Família"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  required
                  disabled={loading}
                />
              </Grid>
              <Grid size={{ xs: 12, sm: 6 }}>
                <TextField
                  fullWidth
                  label="Km Atual"
                  type="number"
                  placeholder="0"
                  value={lastMileage}
                  onChange={(e) => setLastMileage(Number(e.target.value))}
                  disabled={loading}
                  slotProps={{
                    htmlInput: {
                      min: 0,
                    },
                  }}
                />
              </Grid>
              <Grid size={{ xs: 12, sm: 6 }}>
                <TextField
                  fullWidth
                  label="Ano de Fabricação"
                  type="number"
                  placeholder="Ex: 2022"
                  value={yearManufacture || ""}
                  onChange={(e) => setYearManufacture(Number(e.target.value))}
                  disabled={loading}
                  slotProps={{
                    htmlInput: {
                      min: 1900,
                      max: currentYear + 1,
                    },
                  }}
                />
              </Grid>
              <Grid size={{ xs: 12, sm: 6 }}>
                <TextField
                  fullWidth
                  label="Ano do Modelo"
                  type="number"
                  placeholder="Ex: 2023"
                  value={yearModel || ""}
                  onChange={(e) => setYearModel(Number(e.target.value))}
                  disabled={loading}
                  slotProps={{
                    htmlInput: {
                      min: 1900,
                      max: currentYear + 2,
                    },
                  }}
                />
              </Grid>
            </Grid>
          </Box>

          {/* Upload de Foto do Veículo */}
          <ImageUploadField
            value={imageUrl}
            onChange={setImageUrl}
            vehicleType={car.vehicleType}
            disabled={loading}
          />
        </Stack>
      </DialogContent>

      <DialogActions sx={{ p: 2.5, gap: 1 }}>
        <Button
          onClick={onClose}
          variant="outlined"
          color="inherit"
          disabled={loading}
        >
          Cancelar
        </Button>
        <Button
          onClick={() => handleSubmit()}
          variant="contained"
          disabled={
            loading ||
            !name.trim() ||
            Number(lastMileage) < 0 ||
            Number(yearManufacture) < 1900 ||
            Number(yearModel) < 1900
          }
        >
          {loading ? (
            <CircularProgress size={22} sx={{ color: "#FFFFFF" }} />
          ) : (
            "Salvar Alterações"
          )}
        </Button>
      </DialogActions>
    </>
  );
}

export function EditCarDialog({
  open,
  car,
  onClose,
  onCarUpdated,
}: EditCarDialogProps) {
  if (!car) return null;

  return (
    <Dialog
      open={open}
      onClose={onClose}
      fullWidth
      maxWidth="md"
      slotProps={{
        paper: {
          elevation: 0,
          sx: {
            borderRadius: 3,
            border: "1px solid #E2E8F0",
            p: 1,
          },
        },
      }}
    >
      <EditCarForm
        key={car.id}
        car={car}
        onClose={onClose}
        onCarUpdated={onCarUpdated}
      />
    </Dialog>
  );
}
