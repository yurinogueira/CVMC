import React, { useState, useEffect } from "react";
import {
  Box,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Snackbar,
  Alert,
} from "@mui/material";
import AddRoundedIcon from "@mui/icons-material/AddRounded";
import { useDocumentTitle, PageHeader } from "../../shared";
import { carService } from "../../cars/services/car.service";
import { Car } from "../../cars/types/car.types";
import { ExportHistoryTable } from "../components/ExportHistoryTable";
import { ExportMaintenanceDialog } from "../components/ExportMaintenanceDialog";

export function ReportsPage() {
  useDocumentTitle("Relatórios & Exportações");

  const [cars, setCars] = useState<Car[]>([]);
  const [selectedFilterCarId, setSelectedFilterCarId] = useState<string>("");
  const [dialogOpen, setDialogOpen] = useState(false);
  const [refreshTrigger, setRefreshTrigger] = useState(0);
  const [snackbar, setSnackbar] = useState<{
    open: boolean;
    message: string;
    severity: "success" | "error";
  }>({
    open: false,
    message: "",
    severity: "success",
  });

  useEffect(() => {
    void (async () => {
      try {
        const userCars = await carService.list();
        setCars(userCars);
      } catch (err) {
        console.error("Falha ao carregar lista de veículos:", err);
      }
    })();
  }, []);

  const handleExportStarted = () => {
    setRefreshTrigger((prev) => prev + 1);
    setSnackbar({
      open: true,
      message:
        "Exportação solicitada com sucesso! O relatório está na fila de processamento.",
      severity: "success",
    });
  };

  return (
    <Box sx={{ width: "100%", pb: 4 }}>
      {/* Page Header */}
      <PageHeader
        title="Relatórios & Exportações"
        subtitle="Gere e baixe dossiês completos de manutenções com comprovantes anexados"
        action={
          <Button
            variant="contained"
            startIcon={<AddRoundedIcon />}
            onClick={() => setDialogOpen(true)}
            disabled={cars.length === 0}
            sx={{ textTransform: "none", fontWeight: 600, px: 2.5 }}
          >
            Exportar Manutenções (PDF)
          </Button>
        }
      />

      {/* Filter Toolbar */}
      {cars.length > 0 && (
        <Box sx={{ mb: 2.5, maxWidth: 300 }}>
          <FormControl fullWidth size="small">
            <InputLabel id="filter-car-label">Filtrar por Veículo</InputLabel>
            <Select
              labelId="filter-car-label"
              value={selectedFilterCarId}
              label="Filtrar por Veículo"
              onChange={(e) => setSelectedFilterCarId(e.target.value)}
            >
              <MenuItem value="">Todos os Veículos</MenuItem>
              {cars.map((c) => (
                <MenuItem key={c.id} value={c.id}>
                  {c.name || `${c.manufacturer} ${c.model}`}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Box>
      )}

      {/* Reports History Table */}
      <ExportHistoryTable
        carId={selectedFilterCarId}
        refreshTrigger={refreshTrigger}
      />

      {/* Dialog for starting new export */}
      <ExportMaintenanceDialog
        open={dialogOpen}
        onClose={() => setDialogOpen(false)}
        cars={cars}
        initialCarId={selectedFilterCarId}
        onSuccess={handleExportStarted}
      />

      <Snackbar
        open={snackbar.open}
        autoHideDuration={4000}
        onClose={() => setSnackbar((prev) => ({ ...prev, open: false }))}
      >
        <Alert severity={snackbar.severity} sx={{ width: "100%" }}>
          {snackbar.message}
        </Alert>
      </Snackbar>
    </Box>
  );
}
