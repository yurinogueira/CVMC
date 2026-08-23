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
} from "@mui/material";
import { carService } from "../services/car.service";
import { Car } from "../types/car.types";

interface AddCarDialogProps {
  open: boolean;
  onClose: () => void;
  onCarCreated: (car: Car) => void;
}

export function AddCarDialog({
  open,
  onClose,
  onCarCreated,
}: AddCarDialogProps) {
  const currentYear = new Date().getFullYear();

  const [name, setName] = useState("");
  const [manufacturer, setManufacturer] = useState("");
  const [model, setModel] = useState("");
  const [yearManufacture, setYearManufacture] = useState<number>(currentYear);
  const [yearModel, setYearModel] = useState<number>(currentYear);
  const [lastMileage, setLastMileage] = useState<number>(0);
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const resetForm = () => {
    setName("");
    setManufacturer("");
    setModel("");
    setYearManufacture(currentYear);
    setYearModel(currentYear);
    setLastMileage(0);
    setErrorMsg(null);
  };

  const handleClose = () => {
    resetForm();
    onClose();
  };

  const handleSubmit = async (e?: React.FormEvent) => {
    if (e) e.preventDefault();
    setErrorMsg(null);

    if (!name.trim() || !manufacturer.trim() || !model.trim()) {
      setErrorMsg("Preencha o apelido, fabricante e modelo do veículo.");
      return;
    }

    try {
      setLoading(true);
      const newCar = await carService.create({
        name: name.trim(),
        manufacturer: manufacturer.trim(),
        model: model.trim(),
        yearManufacture: Number(yearManufacture),
        yearModel: Number(yearModel),
        lastMileage: Number(lastMileage) || 0,
      });

      onCarCreated(newCar);
      handleClose();
    } catch (err: unknown) {
      const errorObj = err as { response?: { data?: { message?: string } } };
      setErrorMsg(
        errorObj.response?.data?.message ||
          "Não foi possível cadastrar o veículo. Tente novamente.",
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog
      open={open}
      onClose={handleClose}
      fullWidth
      maxWidth="sm"
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
      <DialogTitle sx={{ fontWeight: 700, pb: 1 }}>
        Cadastrar Novo Veículo
      </DialogTitle>

      <DialogContent dividers sx={{ borderColor: "#F1F5F9" }}>
        {errorMsg && (
          <Alert severity="error" sx={{ mb: 2.5, borderRadius: 1.5 }}>
            {errorMsg}
          </Alert>
        )}

        <Stack spacing={2.5} sx={{ mt: 1 }}>
          <TextField
            fullWidth
            label="Apelido / Identificador"
            placeholder="Ex: Meu Civic, Carro do Trabalho"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
            disabled={loading}
          />

          <Grid container spacing={2}>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Fabricante"
                placeholder="Ex: Honda, Toyota, VW"
                value={manufacturer}
                onChange={(e) => setManufacturer(e.target.value)}
                required
                disabled={loading}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Modelo"
                placeholder="Ex: Civic Sedan EXL"
                value={model}
                onChange={(e) => setModel(e.target.value)}
                required
                disabled={loading}
              />
            </Grid>
          </Grid>

          <Grid container spacing={2}>
            <Grid item xs={12} sm={4}>
              <TextField
                fullWidth
                label="Ano Fabricação"
                type="number"
                value={yearManufacture}
                onChange={(e) => setYearManufacture(Number(e.target.value))}
                disabled={loading}
              />
            </Grid>
            <Grid item xs={12} sm={4}>
              <TextField
                fullWidth
                label="Ano Modelo"
                type="number"
                value={yearModel}
                onChange={(e) => setYearModel(Number(e.target.value))}
                disabled={loading}
              />
            </Grid>
            <Grid item xs={12} sm={4}>
              <TextField
                fullWidth
                label="Km Atual"
                type="number"
                value={lastMileage}
                onChange={(e) => setLastMileage(Number(e.target.value))}
                disabled={loading}
              />
            </Grid>
          </Grid>
        </Stack>
      </DialogContent>

      <DialogActions sx={{ p: 2.5, gap: 1 }}>
        <Button
          onClick={handleClose}
          variant="outlined"
          color="inherit"
          disabled={loading}
        >
          Cancelar
        </Button>
        <Button
          onClick={() => handleSubmit()}
          variant="contained"
          disabled={loading}
        >
          {loading ? (
            <CircularProgress size={22} sx={{ color: "#FFFFFF" }} />
          ) : (
            "Cadastrar Veículo"
          )}
        </Button>
      </DialogActions>
    </Dialog>
  );
}
