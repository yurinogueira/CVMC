import React, { useState, useEffect, useRef } from "react";
import { useParams, useNavigate } from "react-router-dom";
import {
  Box,
  Typography,
  Card,
  Button,
  Stack,
  TextField,
  Chip,
  InputAdornment,
  Alert,
  CircularProgress,
  IconButton,
  Paper,
  Tooltip,
} from "@mui/material";
import Grid from "@mui/material/Grid";
import ArrowBackRoundedIcon from "@mui/icons-material/ArrowBackRounded";
import BuildCircleRoundedIcon from "@mui/icons-material/BuildCircleRounded";
import OpacityRoundedIcon from "@mui/icons-material/OpacityRounded";
import SpeedRoundedIcon from "@mui/icons-material/SpeedRounded";
import ElectricBoltRoundedIcon from "@mui/icons-material/ElectricBoltRounded";
import CheckRoundedIcon from "@mui/icons-material/CheckRounded";
import CloudUploadRoundedIcon from "@mui/icons-material/CloudUploadRounded";
import DeleteOutlineRoundedIcon from "@mui/icons-material/DeleteOutlineRounded";
import PictureAsPdfRoundedIcon from "@mui/icons-material/PictureAsPdfRounded";
import ImageRoundedIcon from "@mui/icons-material/ImageRounded";
import AttachMoneyRoundedIcon from "@mui/icons-material/AttachMoneyRounded";
import ReceiptLongRoundedIcon from "@mui/icons-material/ReceiptLongRounded";
import DirectionsCarRoundedIcon from "@mui/icons-material/DirectionsCarRounded";

import { useDocumentTitle } from "../../shared";
import { carService } from "../../cars/services/car.service";
import { Car } from "../../cars/types/car.types";
import { maintenanceService } from "../services/maintenance.service";
import { MaintenanceAttachment } from "../types/maintenance.types";
import {
  MAINTENANCE_CATEGORIES,
  MaintenanceCategoryGroup,
} from "../constants/maintenanceCategories";
import { OTHER_MAINTENANCE_TYPE } from "../constants/maintenanceTypes";

function formatFileSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

async function processFile(file: File): Promise<MaintenanceAttachment> {
  const isPdf = file.type === "application/pdf";
  const isImage = ["image/jpeg", "image/png", "image/webp"].includes(file.type);

  if (!isPdf && !isImage) {
    throw new Error(
      `Arquivo "${file.name}" não suportado. Utilize PDF ou imagens JPG, PNG, WebP.`,
    );
  }

  if (isPdf && file.size > 2 * 1024 * 1024) {
    throw new Error(
      `O arquivo PDF "${file.name}" (${formatFileSize(file.size)}) excede o limite máximo de 2MB.`,
    );
  }

  if (isImage && file.size > 10 * 1024 * 1024) {
    throw new Error(
      `A imagem "${file.name}" (${formatFileSize(file.size)}) excede o limite máximo de 10MB.`,
    );
  }

  if (isPdf) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = () => {
        resolve({
          id: `att-${Date.now()}-${Math.random().toString(36).substring(2, 9)}`,
          name: file.name,
          size: file.size,
          mimeType: file.type,
          dataUrl: reader.result as string,
          createdAt: new Date().toISOString(),
        });
      };
      reader.onerror = () =>
        reject(new Error(`Falha ao ler o arquivo PDF "${file.name}".`));
      reader.readAsDataURL(file);
    });
  }

  // Otimização de imagem para 90% de qualidade via Canvas
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      const img = new Image();
      img.onload = () => {
        const canvas = document.createElement("canvas");
        const maxDimension = 1920;
        let width = img.width;
        let height = img.height;

        if (width > maxDimension || height > maxDimension) {
          if (width > height) {
            height = Math.round((height * maxDimension) / width);
            width = maxDimension;
          } else {
            width = Math.round((width * maxDimension) / height);
            height = maxDimension;
          }
        }

        canvas.width = width;
        canvas.height = height;
        const ctx = canvas.getContext("2d");
        if (ctx) {
          ctx.drawImage(img, 0, 0, width, height);
          const compressedDataUrl = canvas.toDataURL("image/jpeg", 0.9);
          const base64Content =
            compressedDataUrl.indexOf(",") !== -1
              ? compressedDataUrl.split(",")[1]
              : compressedDataUrl;
          const estimatedBytes = Math.round((base64Content.length * 3) / 4);
          resolve({
            id: `att-${Date.now()}-${Math.random().toString(36).substring(2, 9)}`,
            name: file.name,
            size: estimatedBytes,
            mimeType: "image/jpeg",
            dataUrl: compressedDataUrl,
            createdAt: new Date().toISOString(),
          });
        } else {
          resolve({
            id: `att-${Date.now()}-${Math.random().toString(36).substring(2, 9)}`,
            name: file.name,
            size: file.size,
            mimeType: file.type,
            dataUrl: e.target?.result as string,
            createdAt: new Date().toISOString(),
          });
        }
      };
      img.onerror = () =>
        reject(new Error(`Falha ao decodificar a imagem "${file.name}".`));
      img.src = e.target?.result as string;
    };
    reader.onerror = () =>
      reject(new Error(`Falha ao ler a imagem "${file.name}".`));
    reader.readAsDataURL(file);
  });
}

function getCategoryIcon(categoryId: string) {
  switch (categoryId) {
    case "engine_filters":
      return (
        <BuildCircleRoundedIcon sx={{ fontSize: 20, color: "primary.main" }} />
      );
    case "fluids":
      return <OpacityRoundedIcon sx={{ fontSize: 20, color: "#0284C7" }} />;
    case "brakes_suspension_tires":
      return <SpeedRoundedIcon sx={{ fontSize: 20, color: "#EA580C" }} />;
    case "electrical_comfort_others":
      return (
        <ElectricBoltRoundedIcon sx={{ fontSize: 20, color: "#CA8A04" }} />
      );
    default:
      return (
        <BuildCircleRoundedIcon sx={{ fontSize: 20, color: "primary.main" }} />
      );
  }
}

export function RegisterMaintenancePage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const getTodayString = () => {
    const today = new Date();
    const yyyy = today.getFullYear();
    const mm = String(today.getMonth() + 1).padStart(2, "0");
    const dd = String(today.getDate()).padStart(2, "0");
    return `${yyyy}-${mm}-${dd}`;
  };

  const [car, setCar] = useState<Car | null>(null);
  const [loadingCar, setLoadingCar] = useState(true);

  const [types, setTypes] = useState<string[]>([]);
  const [customType, setCustomType] = useState("");
  const [title, setTitle] = useState("");
  const [titleTouched, setTitleTouched] = useState(false);
  const [date, setDate] = useState(getTodayString());
  const [mileage, setMileage] = useState<number | string>(0);
  const [cost, setCost] = useState("");
  const [description, setDescription] = useState("");
  const [attachments, setAttachments] = useState<MaintenanceAttachment[]>([]);
  const [processingFiles, setProcessingFiles] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const fileInputRef = useRef<HTMLInputElement | null>(null);

  useDocumentTitle(
    car ? `Registrar Manutenção - ${car.name}` : "Registrar Manutenção",
  );

  useEffect(() => {
    if (!id) return;
    const fetchCar = async () => {
      try {
        setLoadingCar(true);
        const data = await carService.get(id);
        setCar(data);
        if (data.lastMileage) {
          setMileage(data.lastMileage);
        }
      } catch {
        setErrorMsg("Não foi possível carregar as informações do veículo.");
      } finally {
        setLoadingCar(false);
      }
    };
    fetchCar();
  }, [id]);

  const toggleType = (item: string) => {
    const isSelected = types.includes(item);
    let next: string[];
    if (isSelected) {
      next = types.filter((t) => t !== item);
    } else {
      next = [...types, item];
    }
    setTypes(next);

    // Atualiza sugestão de título se o usuário não tiver alterado manualmente
    if (!titleTouched || !title.trim()) {
      const displayItems = next
        .map((t) => (t === OTHER_MAINTENANCE_TYPE ? customType.trim() : t))
        .filter(Boolean);
      setTitle(displayItems.join(", "));
    }
  };

  const handleCustomTypeChange = (val: string) => {
    setCustomType(val);
    if (!titleTouched || !title.trim()) {
      const displayItems = types
        .map((t) => (t === OTHER_MAINTENANCE_TYPE ? val.trim() : t))
        .filter(Boolean);
      setTitle(displayItems.join(", "));
    }
  };

  const handleFilesSelected = async (
    e: React.ChangeEvent<HTMLInputElement>,
  ) => {
    const files = e.target.files;
    if (!files || files.length === 0) return;

    setErrorMsg(null);
    setProcessingFiles(true);

    try {
      const newAttachments: MaintenanceAttachment[] = [];
      for (let i = 0; i < files.length; i++) {
        const file = files[i];
        const processed = await processFile(file);
        newAttachments.push(processed);
      }
      setAttachments((prev) => [...prev, ...newAttachments]);
    } catch (err: unknown) {
      const message =
        err instanceof Error ? err.message : "Erro ao processar comprovantes.";
      setErrorMsg(message);
    } finally {
      setProcessingFiles(false);
      if (fileInputRef.current) {
        fileInputRef.current.value = "";
      }
    }
  };

  const handleRemoveAttachment = (attId: string) => {
    setAttachments((prev) => prev.filter((a) => a.id !== attId));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!id) return;
    setErrorMsg(null);

    const trimmedTitle = title.trim();
    if (!trimmedTitle) {
      setErrorMsg("O título do serviço é obrigatório.");
      return;
    }

    if (types.includes(OTHER_MAINTENANCE_TYPE) && !customType.trim()) {
      setErrorMsg(
        "Por favor, especifique o nome da manutenção no campo 'Outro tipo de manutenção'.",
      );
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

    let parsedCost: number | undefined = undefined;
    if (cost.trim()) {
      parsedCost = parseFloat(cost.replace(",", "."));
      if (isNaN(parsedCost) || parsedCost < 0) {
        setErrorMsg("O custo informado deve ser um valor numérico válido.");
        return;
      }
    }

    const finalTypes = types
      .map((t) => (t === OTHER_MAINTENANCE_TYPE ? customType.trim() : t))
      .filter(Boolean);

    try {
      setSubmitting(true);
      const isoDate = new Date(`${date}T12:00:00Z`).toISOString();

      await maintenanceService.create(id, {
        title: trimmedTitle,
        description: description.trim(),
        date: isoDate,
        mileage: numMileage,
        types: finalTypes.length > 0 ? finalTypes : undefined,
        cost: parsedCost,
        attachments: attachments.length > 0 ? attachments : undefined,
      });

      navigate(`/vehicles/${id}`);
    } catch (err: unknown) {
      const errorObj = err as { response?: { data?: { message?: string } } };
      setErrorMsg(
        errorObj.response?.data?.message ||
          "Não foi possível registrar a manutenção. Tente novamente.",
      );
    } finally {
      setSubmitting(false);
    }
  };

  if (loadingCar) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", py: 12 }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box component="form" onSubmit={handleSubmit} sx={{ pb: 6 }}>
      {/* Header & Back Navigation */}
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
            onClick={() => navigate(`/vehicles/${id}`)}
            sx={{ borderRadius: 2, textTransform: "none" }}
          >
            Voltar
          </Button>
          <Box>
            <Typography
              component="h1"
              variant="h5"
              sx={{ fontWeight: 800, color: "text.primary" }}
            >
              Registrar Manutenção
            </Typography>
            {car && (
              <Stack
                direction="row"
                spacing={1}
                sx={{ alignItems: "center", mt: 0.2 }}
              >
                <DirectionsCarRoundedIcon
                  sx={{ fontSize: 16, color: "text.secondary" }}
                />
                <Typography variant="body2" sx={{ color: "text.secondary" }}>
                  Veículo: <strong>{car.name}</strong> ({car.manufacturer}{" "}
                  {car.model})
                </Typography>
              </Stack>
            )}
          </Box>
        </Stack>

        <Stack direction="row" spacing={1.5}>
          <Button
            variant="outlined"
            color="inherit"
            onClick={() => navigate(`/vehicles/${id}`)}
            disabled={submitting || processingFiles}
            sx={{ borderRadius: 2 }}
          >
            Cancelar
          </Button>
          <Button
            type="submit"
            variant="contained"
            disabled={submitting || processingFiles}
            startIcon={
              submitting ? (
                <CircularProgress size={18} color="inherit" />
              ) : (
                <BuildCircleRoundedIcon />
              )
            }
            sx={{ px: 3, borderRadius: 2 }}
          >
            Salvar Manutenção
          </Button>
        </Stack>
      </Stack>

      {errorMsg && (
        <Alert
          severity="error"
          onClose={() => setErrorMsg(null)}
          sx={{ mb: 3, borderRadius: 2 }}
        >
          {errorMsg}
        </Alert>
      )}

      {/* Main Two-Column Layout */}
      <Grid container spacing={3}>
        {/* Left Column (8 cols): Service Types & Details */}
        <Grid size={{ xs: 12, md: 8 }}>
          <Stack spacing={3}>
            {/* Card: Tipos de Manutenção em Categorias */}
            <Card
              elevation={0}
              sx={{
                borderRadius: 3,
                border: "1px solid #E2E8F0",
                bgcolor: "background.paper",
                p: { xs: 2.5, sm: 3 },
              }}
            >
              <Typography
                variant="subtitle1"
                component="h2"
                sx={{ fontWeight: 800, mb: 0.5 }}
              >
                Tipo de Manutenção
              </Typography>
              <Typography
                variant="body2"
                sx={{ color: "text.secondary", mb: 2.5 }}
              >
                Clique nos serviços realizados para selecionar. Você pode marcar
                múltiplos itens para compor a revisão.
              </Typography>

              <Stack spacing={2.5}>
                {MAINTENANCE_CATEGORIES.map(
                  (category: MaintenanceCategoryGroup) => (
                    <Box key={category.id}>
                      <Stack
                        direction="row"
                        spacing={1}
                        sx={{ alignItems: "center", mb: 1 }}
                      >
                        {getCategoryIcon(category.id)}
                        <Typography
                          variant="caption"
                          sx={{
                            fontWeight: 700,
                            color: "text.secondary",
                            textTransform: "uppercase",
                            letterSpacing: 0.5,
                          }}
                        >
                          {category.name}
                        </Typography>
                      </Stack>

                      <Box
                        sx={{
                          display: "flex",
                          flexWrap: "wrap",
                          gap: 1,
                        }}
                      >
                        {category.items.map((item) => {
                          const isSelected = types.includes(item);
                          const isOther = item === OTHER_MAINTENANCE_TYPE;
                          return (
                            <Chip
                              key={item}
                              label={
                                isOther && customType.trim()
                                  ? `Outro: ${customType.trim()}`
                                  : item
                              }
                              onClick={() => toggleType(item)}
                              color={isSelected ? "primary" : "default"}
                              variant={isSelected ? "filled" : "outlined"}
                              icon={
                                isSelected ? (
                                  <CheckRoundedIcon sx={{ fontSize: 16 }} />
                                ) : undefined
                              }
                              sx={{
                                borderRadius: 2,
                                py: 2.2,
                                px: 1,
                                fontSize: "0.85rem",
                                fontWeight: isSelected ? 700 : 500,
                                cursor: "pointer",
                                borderColor: isSelected
                                  ? "primary.main"
                                  : "#CBD5E1",
                                bgcolor: isSelected
                                  ? "primary.main"
                                  : "background.paper",
                                color: isSelected ? "#FFFFFF" : "text.primary",
                                transition: "all 0.15s ease",
                                "&:hover": {
                                  bgcolor: isSelected
                                    ? "primary.dark"
                                    : "rgba(2, 132, 199, 0.08)",
                                  borderColor: "primary.main",
                                },
                              }}
                            />
                          );
                        })}
                      </Box>
                    </Box>
                  ),
                )}

                {/* Campo condicional para outro tipo */}
                {types.includes(OTHER_MAINTENANCE_TYPE) && (
                  <Box
                    sx={{
                      p: 2,
                      borderRadius: 2,
                      bgcolor: "#F8FAFC",
                      border: "1px dashed #CBD5E1",
                      mt: 1,
                    }}
                  >
                    <TextField
                      fullWidth
                      label="Especifique o outro tipo de manutenção"
                      placeholder="Ex: Troca de junta do cabeçote, reparo no escapamento..."
                      value={customType}
                      onChange={(e) => handleCustomTypeChange(e.target.value)}
                      required
                      autoFocus
                      disabled={submitting}
                      helperText="Informe o nome específico do serviço realizado para constar no histórico."
                    />
                  </Box>
                )}
              </Stack>
            </Card>

            {/* Card: Dados Gerais do Serviço */}
            <Card
              elevation={0}
              sx={{
                borderRadius: 3,
                border: "1px solid #E2E8F0",
                bgcolor: "background.paper",
                p: { xs: 2.5, sm: 3 },
              }}
            >
              <Typography
                variant="subtitle1"
                component="h2"
                sx={{ fontWeight: 800, mb: 2 }}
              >
                Detalhes do Serviço
              </Typography>

              <Stack spacing={2.5}>
                <TextField
                  fullWidth
                  label="Título do Serviço"
                  placeholder="Ex: Troca de Óleo e Filtro"
                  value={title}
                  onChange={(e) => {
                    setTitle(e.target.value);
                    setTitleTouched(true);
                  }}
                  required
                  disabled={submitting}
                  helperText="Nome do registro que aparecerá com destaque no histórico."
                />

                <Stack direction={{ xs: "column", sm: "row" }} spacing={2}>
                  <TextField
                    fullWidth
                    label="Data da Realização"
                    type="date"
                    value={date}
                    onChange={(e) => setDate(e.target.value)}
                    required
                    disabled={submitting}
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
                    disabled={submitting}
                    slotProps={{
                      htmlInput: { min: 0 },
                    }}
                    helperText="Quilometragem do odômetro na data da realização."
                  />
                </Stack>

                <TextField
                  fullWidth
                  label="Descrição / Observações (opcional)"
                  placeholder="Ex: Utilizado óleo 0W20 sintético com troca do anel do bujão. Garantia de 10.000 km..."
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  multiline
                  rows={3.5}
                  disabled={submitting}
                />
              </Stack>
            </Card>
          </Stack>
        </Grid>

        {/* Right Column (4 cols): Finance, Receipts & Actions */}
        <Grid size={{ xs: 12, md: 4 }}>
          <Stack spacing={3}>
            {/* Card: Custo Financeiro */}
            <Card
              elevation={0}
              sx={{
                borderRadius: 3,
                border: "1px solid #E2E8F0",
                bgcolor: "background.paper",
                p: 2.5,
              }}
            >
              <Stack
                direction="row"
                spacing={1}
                sx={{ alignItems: "center", mb: 1.5 }}
              >
                <AttachMoneyRoundedIcon color="primary" />
                <Typography variant="subtitle1" sx={{ fontWeight: 800 }}>
                  Custo da Manutenção
                </Typography>
              </Stack>

              <TextField
                fullWidth
                label="Valor Total (opcional)"
                placeholder="0,00"
                type="number"
                value={cost}
                onChange={(e) => setCost(e.target.value)}
                disabled={submitting}
                slotProps={{
                  input: {
                    startAdornment: (
                      <InputAdornment position="start">
                        <Typography
                          variant="body2"
                          sx={{ fontWeight: 700, color: "text.secondary" }}
                        >
                          R$
                        </Typography>
                      </InputAdornment>
                    ),
                  },
                  htmlInput: { min: 0, step: "0.01" },
                }}
                helperText="Total gasto somando peças e mão de obra."
              />
            </Card>

            {/* Card: Comprovantes & Recibos */}
            <Card
              elevation={0}
              sx={{
                borderRadius: 3,
                border: "1px solid #E2E8F0",
                bgcolor: "background.paper",
                p: 2.5,
              }}
            >
              <Stack
                direction="row"
                spacing={1}
                sx={{ alignItems: "center", mb: 1 }}
              >
                <ReceiptLongRoundedIcon color="primary" />
                <Typography variant="subtitle1" sx={{ fontWeight: 800 }}>
                  Comprovantes & Recibos
                </Typography>
              </Stack>

              <Typography
                variant="body2"
                sx={{ color: "text.secondary", mb: 2 }}
              >
                Anexe ordens de serviço, notas fiscais ou fotos. Aceita PDF (até
                2MB) e imagens JPG, PNG, WebP (até 10MB com otimização
                automática a 90%).
              </Typography>

              <input
                ref={fileInputRef}
                type="file"
                multiple
                accept="image/jpeg,image/png,image/webp,application/pdf"
                style={{ display: "none" }}
                onChange={handleFilesSelected}
                disabled={submitting || processingFiles}
              />

              <Button
                fullWidth
                variant="outlined"
                startIcon={
                  processingFiles ? (
                    <CircularProgress size={16} color="inherit" />
                  ) : (
                    <CloudUploadRoundedIcon />
                  )
                }
                onClick={() => fileInputRef.current?.click()}
                disabled={submitting || processingFiles}
                sx={{
                  borderRadius: 2,
                  py: 1.2,
                  textTransform: "none",
                  fontWeight: 600,
                }}
              >
                {processingFiles
                  ? "Otimizando comprovantes..."
                  : "Anexar Arquivos"}
              </Button>

              {/* Lista de comprovantes */}
              {attachments.length > 0 && (
                <Stack spacing={1} sx={{ mt: 2 }}>
                  {attachments.map((att) => {
                    const isPdf = att.mimeType === "application/pdf";
                    return (
                      <Paper
                        key={att.id}
                        variant="outlined"
                        sx={{
                          p: 1.2,
                          display: "flex",
                          alignItems: "center",
                          justifyContent: "space-between",
                          borderRadius: 2,
                          bgcolor: "#F8FAFC",
                          borderColor: "#E2E8F0",
                        }}
                      >
                        <Stack
                          direction="row"
                          spacing={1.5}
                          sx={{ alignItems: "center", minWidth: 0 }}
                        >
                          {isPdf ? (
                            <PictureAsPdfRoundedIcon
                              sx={{
                                color: "#EF4444",
                                fontSize: 22,
                                flexShrink: 0,
                              }}
                            />
                          ) : (
                            <ImageRoundedIcon
                              sx={{
                                color: "primary.main",
                                fontSize: 22,
                                flexShrink: 0,
                              }}
                            />
                          )}
                          <Box sx={{ minWidth: 0 }}>
                            <Typography
                              variant="body2"
                              sx={{
                                fontWeight: 600,
                                overflow: "hidden",
                                textOverflow: "ellipsis",
                                whiteSpace: "nowrap",
                                maxWidth: 160,
                              }}
                            >
                              {att.name}
                            </Typography>
                            <Typography
                              variant="caption"
                              sx={{ color: "text.secondary" }}
                            >
                              {formatFileSize(att.size)} •{" "}
                              {isPdf ? "PDF" : "Imagem 90%"}
                            </Typography>
                          </Box>
                        </Stack>

                        <Tooltip title="Remover comprovante">
                          <IconButton
                            size="small"
                            color="error"
                            onClick={() => handleRemoveAttachment(att.id)}
                            disabled={submitting || processingFiles}
                            aria-label={`Remover comprovante ${att.name}`}
                          >
                            <DeleteOutlineRoundedIcon sx={{ fontSize: 18 }} />
                          </IconButton>
                        </Tooltip>
                      </Paper>
                    );
                  })}
                </Stack>
              )}
            </Card>

            {/* Card: Ações de Salvar e Cancelar */}
            <Card
              elevation={0}
              sx={{
                borderRadius: 3,
                border: "1px solid #E2E8F0",
                bgcolor: "background.paper",
                p: 2.5,
              }}
            >
              <Stack spacing={1.5}>
                <Button
                  type="submit"
                  variant="contained"
                  fullWidth
                  size="large"
                  disabled={submitting || processingFiles}
                  startIcon={
                    submitting ? (
                      <CircularProgress size={18} color="inherit" />
                    ) : (
                      <BuildCircleRoundedIcon />
                    )
                  }
                  sx={{ py: 1.3, borderRadius: 2, fontWeight: 700 }}
                >
                  Salvar Manutenção
                </Button>
                <Button
                  variant="outlined"
                  color="inherit"
                  fullWidth
                  onClick={() => navigate(`/vehicles/${id}`)}
                  disabled={submitting || processingFiles}
                  sx={{ py: 1.1, borderRadius: 2 }}
                >
                  Cancelar
                </Button>
              </Stack>
            </Card>
          </Stack>
        </Grid>
      </Grid>
    </Box>
  );
}
