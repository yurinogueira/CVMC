import React, { useEffect, useState, useCallback, useRef } from "react";
import {
  Box,
  Card,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TablePagination,
  Chip,
  IconButton,
  Tooltip,
  Typography,
  CircularProgress,
  Snackbar,
  Alert,
  Stack,
  Button,
} from "@mui/material";
import DownloadRoundedIcon from "@mui/icons-material/DownloadRounded";
import RefreshRoundedIcon from "@mui/icons-material/RefreshRounded";
import ErrorOutlineRoundedIcon from "@mui/icons-material/ErrorOutlineRounded";
import HourglassEmptyRoundedIcon from "@mui/icons-material/HourglassEmptyRounded";
import CheckCircleRoundedIcon from "@mui/icons-material/CheckCircleRounded";
import CancelRoundedIcon from "@mui/icons-material/CancelRounded";
import HistoryRoundedIcon from "@mui/icons-material/HistoryRounded";
import { ReportJob, ReportStatus } from "../types/report.types";
import { reportService } from "../services/report.service";

interface ExportHistoryTableProps {
  carId?: string;
  refreshTrigger?: number;
  onRefreshRequested?: () => void;
}

export const formatDuration = (ms?: number): string => {
  if (ms === undefined || ms === null || ms < 0) return "-";
  if (ms < 1000) return `${ms}ms`;
  const seconds = ms / 1000;
  if (seconds < 60) return `${seconds.toFixed(1)}s`;
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = Math.round(seconds % 60);
  return remainingSeconds > 0
    ? `${minutes}m ${remainingSeconds}s`
    : `${minutes}m`;
};

export const formatDate = (dateStr?: string): string => {
  if (!dateStr) return "-";
  try {
    const d = new Date(dateStr);
    if (isNaN(d.getTime())) return dateStr;
    return d.toLocaleString("pt-BR", {
      day: "2-digit",
      month: "2-digit",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  } catch {
    return dateStr;
  }
};

export const formatExpiration = (
  expiresAt?: string,
  status?: ReportStatus,
): string => {
  if (status === "expired") return "Expirado";
  if (!expiresAt) return "-";
  try {
    const exp = new Date(expiresAt);
    const now = new Date();
    const diffMs = exp.getTime() - now.getTime();
    if (diffMs <= 0) return "Expirado";
    const days = Math.ceil(diffMs / (1000 * 60 * 60 * 24));
    return days === 1 ? "Expira amanhã" : `Expira em ${days} dias`;
  } catch {
    return "-";
  }
};

export const ExportHistoryTable: React.FC<ExportHistoryTableProps> = ({
  carId,
  refreshTrigger = 0,
  onRefreshRequested,
}) => {
  const [jobs, setJobs] = useState<ReportJob[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(0);
  const [limit, setLimit] = useState(10);
  const [loading, setLoading] = useState(true);
  const [downloadingId, setDownloadingId] = useState<string | null>(null);
  const [snackbar, setSnackbar] = useState<{
    open: boolean;
    message: string;
    severity: "success" | "error";
  }>({
    open: false,
    message: "",
    severity: "success",
  });

  const pollingRef = useRef<ReturnType<typeof setInterval> | null>(null);

  const fetchJobs = useCallback(
    async (isPolling = false) => {
      if (!isPolling) setLoading(true);
      try {
        const res = await reportService.listJobs({
          page: page + 1,
          limit,
          carId: carId || undefined,
        });
        setJobs(res.jobs || []);
        setTotal(res.total || 0);
      } catch (err) {
        if (!isPolling) {
          console.error("Falha ao carregar relatórios:", err);
        }
      } finally {
        if (!isPolling) setLoading(false);
      }
    },
    [carId, page, limit],
  );

  useEffect(() => {
    let isMounted = true;
    reportService
      .listJobs({
        page: page + 1,
        limit,
        carId: carId || undefined,
      })
      .then((res) => {
        if (isMounted) {
          setJobs(res.jobs || []);
          setTotal(res.total || 0);
          setLoading(false);
        }
      })
      .catch((err) => {
        if (isMounted) {
          console.error("Falha ao carregar relatórios:", err);
          setLoading(false);
        }
      });

    return () => {
      isMounted = false;
    };
  }, [carId, page, limit, refreshTrigger]);

  // Polling management: automatically poll while jobs are pending or processing
  useEffect(() => {
    const hasActiveJobs = jobs.some(
      (j) => j.status === "pending" || j.status === "processing",
    );

    if (hasActiveJobs) {
      if (!pollingRef.current) {
        pollingRef.current = setInterval(() => {
          void fetchJobs(true);
        }, 3000);
      }
    } else {
      if (pollingRef.current) {
        clearInterval(pollingRef.current);
        pollingRef.current = null;
      }
    }

    return () => {
      if (pollingRef.current) {
        clearInterval(pollingRef.current);
        pollingRef.current = null;
      }
    };
  }, [jobs, fetchJobs]);

  const handleDownload = async (job: ReportJob) => {
    if (!job.filePath) return;
    setDownloadingId(job.id);
    try {
      const fileName = `dossie_manutencoes_${job.carName.replace(/\s+/g, "_")}.pdf`;
      await reportService.downloadReport(job.filePath, fileName);
      setSnackbar({
        open: true,
        message: "Download iniciado com sucesso!",
        severity: "success",
      });
    } catch {
      setSnackbar({
        open: true,
        message:
          "Erro ao baixar relatório. Verifique se o arquivo não expirou.",
        severity: "error",
      });
    } finally {
      setDownloadingId(null);
    }
  };

  const getStatusChip = (status: ReportStatus) => {
    switch (status) {
      case "completed":
        return (
          <Chip
            size="small"
            icon={<CheckCircleRoundedIcon fontSize="small" />}
            label="Concluído"
            color="success"
            variant="filled"
          />
        );
      case "processing":
        return (
          <Chip
            size="small"
            icon={<CircularProgress size={12} color="inherit" />}
            label="Gerando PDF..."
            color="primary"
            variant="outlined"
          />
        );
      case "pending":
        return (
          <Chip
            size="small"
            icon={<HourglassEmptyRoundedIcon fontSize="small" />}
            label="Na fila..."
            color="warning"
            variant="outlined"
          />
        );
      case "failed":
        return (
          <Chip
            size="small"
            icon={<CancelRoundedIcon fontSize="small" />}
            label="Falhou"
            color="error"
            variant="filled"
          />
        );
      case "expired":
        return (
          <Chip
            size="small"
            icon={<HistoryRoundedIcon fontSize="small" />}
            label="Expirado (30d)"
            variant="outlined"
            sx={{ color: "text.secondary", borderColor: "#CBD5E1" }}
          />
        );
      default:
        return <Chip size="small" label={status} />;
    }
  };

  return (
    <Card variant="outlined" sx={{ borderRadius: 2 }}>
      <Box
        sx={{
          p: 2,
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          borderBottom: "1px solid #E2E8F0",
        }}
      >
        <Stack direction="row" spacing={1} sx={{ alignItems: "center" }}>
          <HistoryRoundedIcon color="primary" />
          <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>
            Histórico de Exportações
          </Typography>
        </Stack>

        <Stack direction="row" spacing={1}>
          <Tooltip title="Atualizar listagem">
            <span>
              <IconButton
                size="small"
                onClick={() => {
                  void fetchJobs();
                  if (onRefreshRequested) onRefreshRequested();
                }}
                disabled={loading}
              >
                <RefreshRoundedIcon
                  fontSize="small"
                  sx={{
                    animation: loading ? "spin 1s linear infinite" : "none",
                    "@keyframes spin": {
                      "0%": { transform: "rotate(0deg)" },
                      "100%": { transform: "rotate(360deg)" },
                    },
                  }}
                />
              </IconButton>
            </span>
          </Tooltip>
        </Stack>
      </Box>

      <TableContainer>
        <Table size="small">
          <TableHead sx={{ bgcolor: "#F8FAFC" }}>
            <TableRow>
              <TableCell sx={{ fontWeight: 600 }}>Veículo</TableCell>
              <TableCell sx={{ fontWeight: 600 }}>Solicitado em</TableCell>
              <TableCell sx={{ fontWeight: 600 }}>Duração</TableCell>
              <TableCell sx={{ fontWeight: 600 }}>Status</TableCell>
              <TableCell sx={{ fontWeight: 600 }}>Validade (30 dias)</TableCell>
              <TableCell align="right" sx={{ fontWeight: 600 }}>
                Ações
              </TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {loading && jobs.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} align="center" sx={{ py: 4 }}>
                  <CircularProgress size={24} />
                  <Typography
                    variant="body2"
                    sx={{ mt: 1, color: "text.secondary" }}
                  >
                    Carregando histórico...
                  </Typography>
                </TableCell>
              </TableRow>
            ) : jobs.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} align="center" sx={{ py: 4 }}>
                  <Typography variant="body2" color="text.secondary">
                    Nenhum relatório exportado até o momento.
                  </Typography>
                </TableCell>
              </TableRow>
            ) : (
              jobs.map((job) => (
                <TableRow key={job.id} hover>
                  <TableCell>
                    <Typography variant="body2" sx={{ fontWeight: 600 }}>
                      {job.carName}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      Dossiê de Manutenções (PDF)
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Typography variant="body2">
                      {formatDate(job.createdAt)}
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Typography variant="body2">
                      {formatDuration(job.durationMs)}
                    </Typography>
                  </TableCell>
                  <TableCell>{getStatusChip(job.status)}</TableCell>
                  <TableCell>
                    <Typography variant="body2" color="text.secondary">
                      {formatExpiration(job.expiresAt, job.status)}
                    </Typography>
                  </TableCell>
                  <TableCell align="right">
                    {job.status === "completed" && (
                      <Button
                        size="small"
                        variant="contained"
                        startIcon={
                          downloadingId === job.id ? (
                            <CircularProgress size={14} color="inherit" />
                          ) : (
                            <DownloadRoundedIcon fontSize="small" />
                          )
                        }
                        disabled={downloadingId === job.id}
                        onClick={() => void handleDownload(job)}
                        sx={{ textTransform: "none", py: 0.5 }}
                      >
                        Baixar PDF
                      </Button>
                    )}

                    {job.status === "failed" && job.error && (
                      <Tooltip title={`Erro: ${job.error}`}>
                        <IconButton size="small" color="error">
                          <ErrorOutlineRoundedIcon fontSize="small" />
                        </IconButton>
                      </Tooltip>
                    )}

                    {job.status === "expired" && (
                      <Typography variant="caption" color="text.secondary">
                        Indisponível
                      </Typography>
                    )}
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </TableContainer>

      <TablePagination
        component="div"
        count={total}
        page={page}
        onPageChange={(_, newPage) => setPage(newPage)}
        rowsPerPage={limit}
        onRowsPerPageChange={(e) => {
          setLimit(parseInt(e.target.value, 10));
          setPage(0);
        }}
        rowsPerPageOptions={[5, 10, 25]}
        labelRowsPerPage="Linhas por página:"
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
    </Card>
  );
};
