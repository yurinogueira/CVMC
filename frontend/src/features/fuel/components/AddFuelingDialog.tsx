import React, { useState } from "react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Grid,
  MenuItem,
  FormControlLabel,
  Switch,
  Alert,
  CircularProgress,
  Stack,
  InputAdornment,
} from "@mui/material";
import LocalGasStationRoundedIcon from "@mui/icons-material/LocalGasStationRounded";
import CalendarMonthRoundedIcon from "@mui/icons-material/CalendarMonthRounded";
import AttachMoneyRoundedIcon from "@mui/icons-material/AttachMoneyRounded";
import PlaceRoundedIcon from "@mui/icons-material/PlaceRounded";

import { Fueling, FUEL_TYPES } from "../types/fuel.types";
import { fuelService } from "../services/fuel.service";

interface AddFuelingDialogProps {
  open: boolean;
  onClose: () => void;
  onSuccess: (fueling: Fueling, isEdit: boolean) => void;
  initialData?: Fueling | null;
  carId: string;
}

function getTodayDate() {
  const today = new Date();
  const year = today.getFullYear();
  const month = String(today.getMonth() + 1).padStart(2, "0");
  const day = String(today.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
}

function AddFuelingDialogModal({
  open,
  onClose,
  onSuccess,
  initialData,
  carId,
}: AddFuelingDialogProps) {
  const isEdit = Boolean(initialData);

  const [date, setDate] = useState<string>(() =>
    initialData?.date ? initialData.date.split("T")[0] : getTodayDate(),
  );
  const [fuelType, setFuelType] = useState<string>(
    () => initialData?.fuelType || "Gasolina Comum",
  );
  const [liters, setLiters] = useState<string>(() =>
    initialData?.liters ? String(initialData.liters) : "",
  );
  const [pricePerLiter, setPricePerLiter] = useState<string>(() =>
    initialData?.pricePerLiter ? String(initialData.pricePerLiter) : "",
  );
  const [totalCost, setTotalCost] = useState<string>(() =>
    initialData?.totalCost ? String(initialData.totalCost) : "",
  );
  const [isFullTank, setIsFullTank] = useState<boolean>(
    () => initialData?.isFullTank ?? true,
  );
  const [gasStation, setGasStation] = useState<string>(
    () => initialData?.gasStation || "",
  );
  const [notes, setNotes] = useState<string>(() => initialData?.notes || "");

  const [loading, setLoading] = useState<boolean>(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  // Recálculo reativo bidirecional inteligente
  const handleLitersChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value;
    setLiters(val);
    const numLiters = parseFloat(val.replace(",", "."));
    const numPrice = parseFloat(pricePerLiter.replace(",", "."));
    const numTotal = parseFloat(totalCost.replace(",", "."));

    if (!isNaN(numLiters) && numLiters > 0) {
      if (!isNaN(numPrice) && numPrice > 0) {
        setTotalCost((numLiters * numPrice).toFixed(2));
      } else if (!isNaN(numTotal) && numTotal > 0) {
        setPricePerLiter((numTotal / numLiters).toFixed(3));
      }
    }
  };

  const handleTotalCostChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value;
    setTotalCost(val);
    const numTotal = parseFloat(val.replace(",", "."));
    const numLiters = parseFloat(liters.replace(",", "."));
    const numPrice = parseFloat(pricePerLiter.replace(",", "."));

    if (!isNaN(numTotal) && numTotal > 0) {
      if (!isNaN(numLiters) && numLiters > 0) {
        setPricePerLiter((numTotal / numLiters).toFixed(3));
      } else if (!isNaN(numPrice) && numPrice > 0) {
        setLiters((numTotal / numPrice).toFixed(2));
      }
    }
  };

  const handlePriceChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value;
    setPricePerLiter(val);
    const numPrice = parseFloat(val.replace(",", "."));
    const numLiters = parseFloat(liters.replace(",", "."));
    const numTotal = parseFloat(totalCost.replace(",", "."));

    if (!isNaN(numPrice) && numPrice > 0) {
      if (!isNaN(numLiters) && numLiters > 0) {
        setTotalCost((numLiters * numPrice).toFixed(2));
      } else if (!isNaN(numTotal) && numTotal > 0) {
        setLiters((numTotal / numPrice).toFixed(2));
      }
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMsg(null);

    const parsedLiters = parseFloat(liters.replace(",", "."));
    let parsedPrice = parseFloat(pricePerLiter.replace(",", "."));
    let parsedTotal = parseFloat(totalCost.replace(",", "."));

    if (isNaN(parsedLiters) || parsedLiters <= 0) {
      setErrorMsg("Informe a quantidade de litros abastecidos.");
      return;
    }

    if (isNaN(parsedTotal) || parsedTotal <= 0) {
      if (!isNaN(parsedPrice) && parsedPrice > 0) {
        parsedTotal = Math.round(parsedLiters * parsedPrice * 100) / 100;
      } else {
        setErrorMsg("Informe o valor total pago ou o preço por litro.");
        return;
      }
    }

    if (isNaN(parsedPrice) || parsedPrice <= 0) {
      parsedPrice = Math.round((parsedTotal / parsedLiters) * 1000) / 1000;
    }

    try {
      setLoading(true);
      const payload = {
        date: date ? new Date(date).toISOString() : new Date().toISOString(),
        fuelType,
        liters: parsedLiters,
        pricePerLiter: parsedPrice,
        totalCost: parsedTotal,
        isFullTank,
        gasStation: gasStation.trim() || undefined,
        notes: notes.trim() || undefined,
      };

      if (isEdit && initialData) {
        const updated = await fuelService.update(initialData.id, payload);
        onSuccess(updated, true);
      } else {
        const created = await fuelService.create(carId, payload);
        onSuccess(created, false);
      }
      onClose();
    } catch (err: unknown) {
      const errorObj = err as { response?: { data?: { message?: string } } };
      setErrorMsg(
        errorObj.response?.data?.message ||
          "Não foi possível salvar o abastecimento. Verifique os dados e tente novamente.",
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog
      open={open}
      onClose={() => !loading && onClose()}
      maxWidth="sm"
      fullWidth
      slotProps={{
        paper: {
          sx: { borderRadius: 3, p: 1 },
        },
      }}
    >
      <form onSubmit={handleSubmit}>
        <DialogTitle sx={{ fontWeight: 800, pb: 1 }}>
          <Stack direction="row" spacing={1.5} sx={{ alignItems: "center" }}>
            <LocalGasStationRoundedIcon color="primary" />
            <span>
              {isEdit ? "Editar Abastecimento" : "Registrar Abastecimento"}
            </span>
          </Stack>
        </DialogTitle>

        <DialogContent dividers sx={{ pt: 2 }}>
          <Stack spacing={2.5}>
            {errorMsg && (
              <Alert severity="error" sx={{ borderRadius: 2 }}>
                {errorMsg}
              </Alert>
            )}

            <Grid container spacing={2}>
              {/* Data do Abastecimento */}
              <Grid size={{ xs: 12, sm: 6 }}>
                <TextField
                  fullWidth
                  label="Data do Abastecimento"
                  type="date"
                  value={date}
                  onChange={(e) => setDate(e.target.value)}
                  required
                  slotProps={{
                    input: {
                      startAdornment: (
                        <InputAdornment position="start">
                          <CalendarMonthRoundedIcon fontSize="small" />
                        </InputAdornment>
                      ),
                    },
                  }}
                />
              </Grid>

              {/* Tipo de Combustível */}
              <Grid size={{ xs: 12, sm: 6 }}>
                <TextField
                  select
                  fullWidth
                  label="Tipo de Combustível"
                  value={fuelType}
                  onChange={(e) => setFuelType(e.target.value)}
                  required
                >
                  {FUEL_TYPES.map((type) => (
                    <MenuItem key={type} value={type}>
                      {type}
                    </MenuItem>
                  ))}
                </TextField>
              </Grid>

              {/* Litros Abastecidos */}
              <Grid size={{ xs: 12, sm: 4 }}>
                <TextField
                  fullWidth
                  label="Litros Abastecidos"
                  placeholder="Ex: 45.5"
                  value={liters}
                  onChange={handleLitersChange}
                  required
                  type="text"
                  slotProps={{
                    input: {
                      endAdornment: (
                        <InputAdornment position="end">L</InputAdornment>
                      ),
                    },
                  }}
                  helperText="Quantidade colocada"
                />
              </Grid>

              {/* Valor Total Pago */}
              <Grid size={{ xs: 12, sm: 4 }}>
                <TextField
                  fullWidth
                  label="Valor Total Pago"
                  placeholder="Ex: 250.00"
                  value={totalCost}
                  onChange={handleTotalCostChange}
                  required
                  type="text"
                  slotProps={{
                    input: {
                      startAdornment: (
                        <InputAdornment position="start">R$</InputAdornment>
                      ),
                    },
                  }}
                  helperText="Total pago no posto"
                />
              </Grid>

              {/* Preço por Litro (Calculado automaticamente ou editável) */}
              <Grid size={{ xs: 12, sm: 4 }}>
                <TextField
                  fullWidth
                  label="Preço por Litro"
                  placeholder="Ex: 5.89"
                  value={pricePerLiter}
                  onChange={handlePriceChange}
                  type="text"
                  slotProps={{
                    input: {
                      startAdornment: (
                        <InputAdornment position="start">R$</InputAdornment>
                      ),
                    },
                  }}
                  helperText="Calculado automático"
                />
              </Grid>

              {/* Posto / Estabelecimento */}
              <Grid size={{ xs: 12, sm: 8 }}>
                <TextField
                  fullWidth
                  label="Posto / Estabelecimento (Opcional)"
                  placeholder="Ex: Posto Ipiranga Av. Central"
                  value={gasStation}
                  onChange={(e) => setGasStation(e.target.value)}
                  slotProps={{
                    input: {
                      startAdornment: (
                        <InputAdornment position="start">
                          <PlaceRoundedIcon fontSize="small" />
                        </InputAdornment>
                      ),
                    },
                  }}
                />
              </Grid>

              {/* Tanque Cheio? */}
              <Grid
                size={{ xs: 12, sm: 4 }}
                sx={{ display: "flex", alignItems: "center" }}
              >
                <FormControlLabel
                  control={
                    <Switch
                      checked={isFullTank}
                      onChange={(e) => setIsFullTank(e.target.checked)}
                      color="primary"
                    />
                  }
                  label="Tanque Cheio?"
                />
              </Grid>

              {/* Observações */}
              <Grid size={{ xs: 12 }}>
                <TextField
                  fullWidth
                  multiline
                  rows={2}
                  label="Observações (Opcional)"
                  placeholder="Anotações sobre a viagem, bandeira do combustível, etc."
                  value={notes}
                  onChange={(e) => setNotes(e.target.value)}
                />
              </Grid>
            </Grid>
          </Stack>
        </DialogContent>

        <DialogActions sx={{ px: 3, py: 2 }}>
          <Button
            variant="outlined"
            color="inherit"
            onClick={onClose}
            disabled={loading}
            sx={{ borderRadius: 2 }}
          >
            Cancelar
          </Button>
          <Button
            type="submit"
            variant="contained"
            disabled={loading}
            startIcon={
              loading ? (
                <CircularProgress size={18} color="inherit" />
              ) : (
                <AttachMoneyRoundedIcon />
              )
            }
            sx={{ borderRadius: 2 }}
          >
            {isEdit ? "Salvar Alterações" : "Registrar Abastecimento"}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
}

export function AddFuelingDialog(props: AddFuelingDialogProps) {
  if (!props.open) return null;
  return (
    <AddFuelingDialogModal
      key={props.initialData?.id || (props.open ? "open" : "closed")}
      {...props}
    />
  );
}
