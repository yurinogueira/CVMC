import React, { useState } from "react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Typography,
  Alert,
  CircularProgress,
  Stack,
  Box,
} from "@mui/material";
import PictureAsPdfRoundedIcon from "@mui/icons-material/PictureAsPdfRounded";
import FilterDramaRoundedIcon from "@mui/icons-material/FilterDramaRounded";
import { Car } from "../../cars/types/car.types";
import { reportService } from "../services/report.service";

interface ExportMaintenanceDialogProps {
  open: boolean;
  onClose: () => void;
  cars: Car[];
  initialCarId?: string;
  onSuccess: () => void;
}

export const ExportMaintenanceDialog: React.FC<
  ExportMaintenanceDialogProps
> = ({ open, onClose, cars, initialCarId = "", onSuccess }) => {
  const [selectedCarId, setSelectedCarId] = useState<string>("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const effectiveCarId =
    selectedCarId || initialCarId || (cars.length > 0 ? cars[0].id : "");

  const handleSubmit = async () => {
    if (!effectiveCarId) {
      setError("Por favor, selecione um veículo.");
      return;
    }

    setLoading(true);
    setError(null);
    try {
      await reportService.startMaintenanceReport(effectiveCarId);
      onSuccess();
      onClose();
    } catch (err: unknown) {
      const msg =
        (err as { response?: { data?: { message?: string } } })?.response?.data
          ?.message || "Falha ao iniciar exportação do relatório.";
      setError(msg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog
      open={open}
      onClose={loading ? undefined : onClose}
      maxWidth="sm"
      fullWidth
    >
      <DialogTitle sx={{ display: "flex", alignItems: "center", gap: 1 }}>
        <PictureAsPdfRoundedIcon color="primary" />
        <Typography variant="h6" sx={{ fontWeight: 600 }}>
          Exportar Dossiê de Manutenções (PDF)
        </Typography>
      </DialogTitle>

      <DialogContent dividers>
        <Stack spacing={2.5}>
          <Typography variant="body2" color="text.secondary">
            Gere um documento formal em PDF contendo o histórico detalhado de
            todas as manutenções, peças e custos do seu veículo, com cada
            comprovante e nota fiscal anexada diagramada em sua própria página
            ao final do relatório.
          </Typography>

          <Alert
            severity="info"
            icon={<FilterDramaRoundedIcon fontSize="inherit" />}
          >
            <strong>Geração Otimizada em Fila:</strong> Os relatórios são
            gerados em segundo plano de forma sequencial para máxima
            estabilidade e eficiência. O documento ficará disponível para
            download por <strong>30 dias</strong>.
          </Alert>

          {error && <Alert severity="error">{error}</Alert>}

          <Box sx={{ pt: 1 }}>
            <FormControl fullWidth size="small">
              <InputLabel id="select-car-label">Selecione o Veículo</InputLabel>
              <Select
                labelId="select-car-label"
                value={effectiveCarId}
                label="Selecione o Veículo"
                onChange={(e) => setSelectedCarId(e.target.value)}
                disabled={loading}
              >
                {cars.map((car) => (
                  <MenuItem key={car.id} value={car.id}>
                    {car.name || `${car.manufacturer} ${car.model}`} •{" "}
                    {car.lastMileage.toLocaleString("pt-BR")} km
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Box>
        </Stack>
      </DialogContent>

      <DialogActions sx={{ px: 3, py: 2 }}>
        <Button onClick={onClose} disabled={loading} color="inherit">
          Cancelar
        </Button>
        <Button
          variant="contained"
          onClick={() => void handleSubmit()}
          disabled={loading || !effectiveCarId}
          startIcon={
            loading ? <CircularProgress size={16} color="inherit" /> : null
          }
        >
          {loading ? "Iniciando..." : "Solicitar Exportação"}
        </Button>
      </DialogActions>
    </Dialog>
  );
};
