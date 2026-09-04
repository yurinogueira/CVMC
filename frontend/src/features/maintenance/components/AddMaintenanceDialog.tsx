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
  Typography,
} from "@mui/material";
import BuildCircleRoundedIcon from "@mui/icons-material/BuildCircleRounded";
import { maintenanceService } from "../services/maintenance.service";
import { Maintenance } from "../types/maintenance.types";

interface AddMaintenanceDialogProps {
  open: boolean;
  carId: string;
  carName?: string;
  lastMileage?: number;
  onClose: () => void;
  onMaintenanceCreated: (maintenance: Maintenance) => void;
}

export function AddMaintenanceDialog({
  open,
  carId,
  carName,
  lastMileage = 0,
  onClose,
  onMaintenanceCreated,
}: AddMaintenanceDialogProps) {
  const getTodayString = () => {
    const today = new Date();
    const yyyy = today.getFullYear();
    const mm = String(today.getMonth() + 1).padStart(2, "0");
    const dd = String(today.getDate()).padStart(2, "0");
    return `${yyyy}-${mm}-${dd}`;
  };

  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [date, setDate] = useState(getTodayString());
  const [mileage, setMileage] = useState<number | string>(lastMileage);
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const resetForm = () => {
    setTitle("");
    setDescription("");
    setDate(getTodayString());
    setMileage(lastMileage || 0);
    setErrorMsg(null);
  };

  const handleClose = () => {
    if (loading) return;
    resetForm();
    onClose();
  };

  const handleSubmit = async (e?: React.FormEvent) => {
    if (e) e.preventDefault();
    setErrorMsg(null);

    const trimmedTitle = title.trim();
    if (!trimmedTitle) {
      setErrorMsg("O título do serviço é obrigatório.");
      return;
    }

    if (!date) {
      setErrorMsg("Informe a data da realização do serviço.");
      return;
    }

    const numMileage = Number(mileage);
    if (isNaN(numMileage) || numMileage < 0) {
      setErrorMsg("A quilometragem deve ser um número maior ou igual a zero.");
      return;
    }

    try {
      setLoading(true);

      // Construct ISO date string (UTC)
      const isoDate = new Date(`${date}T12:00:00Z`).toISOString();

      const created = await maintenanceService.create(carId, {
        title: trimmedTitle,
        description: description.trim(),
        date: isoDate,
        mileage: numMileage,
      });

      onMaintenanceCreated(created);
      resetForm();
      onClose();
    } catch (err: unknown) {
      const errorObj = err as { response?: { data?: { message?: string } } };
      setErrorMsg(
        errorObj.response?.data?.message ||
          "Não foi possível registrar a manutenção. Tente novamente.",
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
      <DialogTitle
        sx={{
          fontWeight: 700,
          pb: 1,
          display: "flex",
          alignItems: "center",
          gap: 1,
        }}
      >
        <BuildCircleRoundedIcon color="primary" />
        Registrar Manutenção
      </DialogTitle>

      <DialogContent dividers sx={{ borderColor: "#F1F5F9" }}>
        {carName && (
          <Typography variant="body2" sx={{ color: "text.secondary", mb: 2 }}>
            Veículo: <strong>{carName}</strong>
          </Typography>
        )}

        {errorMsg && (
          <Alert
            severity="error"
            onClose={() => setErrorMsg(null)}
            sx={{ mb: 2.5, borderRadius: 1.5 }}
          >
            {errorMsg}
          </Alert>
        )}

        <Stack spacing={2.5} sx={{ mt: 1 }}>
          <TextField
            fullWidth
            label="Título do Serviço"
            placeholder="Ex: Troca de Óleo e Filtro"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
            disabled={loading}
            autoFocus
          />

          <TextField
            fullWidth
            label="Data da Realização"
            type="date"
            value={date}
            onChange={(e) => setDate(e.target.value)}
            required
            disabled={loading}
            slotProps={{
              inputLabel: { shrink: true },
            }}
          />

          <TextField
            fullWidth
            label="Quilometragem no Momento do Serviço (km)"
            type="number"
            value={mileage}
            onChange={(e) => setMileage(e.target.value)}
            required
            disabled={loading}
            slotProps={{
              htmlInput: { min: 0 },
            }}
          />

          <TextField
            fullWidth
            label="Descrição / Detalhes (opcional)"
            placeholder="Ex: Óleo sintético 0W20, troca dos filtros de ar e combustível..."
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            multiline
            rows={3}
            disabled={loading}
          />
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
          startIcon={
            loading ? <CircularProgress size={18} color="inherit" /> : null
          }
        >
          Salvar Manutenção
        </Button>
      </DialogActions>
    </Dialog>
  );
}
