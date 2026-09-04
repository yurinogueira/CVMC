import React, { useState, useRef } from "react";
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
  Box,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Checkbox,
  ListItemText,
  Chip,
  OutlinedInput,
  InputAdornment,
  IconButton,
  Paper,
  Tooltip,
} from "@mui/material";
import BuildCircleRoundedIcon from "@mui/icons-material/BuildCircleRounded";
import CloudUploadRoundedIcon from "@mui/icons-material/CloudUploadRounded";
import DeleteOutlineRoundedIcon from "@mui/icons-material/DeleteOutlineRounded";
import PictureAsPdfRoundedIcon from "@mui/icons-material/PictureAsPdfRounded";
import ImageRoundedIcon from "@mui/icons-material/ImageRounded";
import ReceiptLongRoundedIcon from "@mui/icons-material/ReceiptLongRounded";
import AttachMoneyRoundedIcon from "@mui/icons-material/AttachMoneyRounded";
import { maintenanceService } from "../services/maintenance.service";
import { Maintenance, MaintenanceAttachment } from "../types/maintenance.types";
import {
  MAINTENANCE_TYPES,
  OTHER_MAINTENANCE_TYPE,
} from "../constants/maintenanceTypes";

interface AddMaintenanceDialogProps {
  open: boolean;
  carId: string;
  carName?: string;
  lastMileage?: number;
  onClose: () => void;
  onMaintenanceCreated: (maintenance: Maintenance) => void;
}

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

  const [types, setTypes] = useState<string[]>([]);
  const [customType, setCustomType] = useState("");
  const [title, setTitle] = useState("");
  const [titleTouched, setTitleTouched] = useState(false);
  const [description, setDescription] = useState("");
  const [date, setDate] = useState(getTodayString());
  const [mileage, setMileage] = useState<number | string>(lastMileage);
  const [cost, setCost] = useState<string>("");
  const [attachments, setAttachments] = useState<MaintenanceAttachment[]>([]);
  const [loading, setLoading] = useState(false);
  const [processingFiles, setProcessingFiles] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const fileInputRef = useRef<HTMLInputElement | null>(null);

  const resetForm = () => {
    setTypes([]);
    setCustomType("");
    setTitle("");
    setTitleTouched(false);
    setDescription("");
    setDate(getTodayString());
    setMileage(lastMileage || 0);
    setCost("");
    setAttachments([]);
    setErrorMsg(null);
  };

  const handleClose = () => {
    if (loading || processingFiles) return;
    resetForm();
    onClose();
  };

  const handleTypesChange = (selected: string[]) => {
    setTypes(selected);

    // Se o usuário não editou o título manualmente ou se o título estiver em branco,
    // sugerimos o título concatenando os tipos selecionados
    if (!titleTouched || !title.trim()) {
      const displayItems = selected
        .map((t) => (t === OTHER_MAINTENANCE_TYPE ? customType.trim() : t))
        .filter(Boolean);
      if (displayItems.length > 0) {
        setTitle(displayItems.join(", "));
      } else {
        setTitle("");
      }
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

  const handleRemoveAttachment = (id: string) => {
    setAttachments((prev) => prev.filter((a) => a.id !== id));
  };

  const handleSubmit = async (e?: React.FormEvent) => {
    if (e) e.preventDefault();
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

    // Mapear tipos finais incluindo a especificação customizada se houver
    const finalTypes = types
      .map((t) => (t === OTHER_MAINTENANCE_TYPE ? customType.trim() : t))
      .filter(Boolean);

    try {
      setLoading(true);

      const isoDate = new Date(`${date}T12:00:00Z`).toISOString();

      const created = await maintenanceService.create(carId, {
        title: trimmedTitle,
        description: description.trim(),
        date: isoDate,
        mileage: numMileage,
        types: finalTypes.length > 0 ? finalTypes : undefined,
        cost: parsedCost,
        attachments: attachments.length > 0 ? attachments : undefined,
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
          {/* Seletor Tipo de Manutenção (Multi-select) */}
          <FormControl fullWidth size="medium">
            <InputLabel id="maintenance-types-label">
              Tipo de Manutenção
            </InputLabel>
            <Select
              labelId="maintenance-types-label"
              multiple
              value={types}
              onChange={(e) => {
                const value = e.target.value;
                handleTypesChange(
                  typeof value === "string" ? value.split(",") : value,
                );
              }}
              input={
                <OutlinedInput
                  id="select-multiple-types"
                  label="Tipo de Manutenção"
                />
              }
              renderValue={(selected) => (
                <Box sx={{ display: "flex", flexWrap: "wrap", gap: 0.5 }}>
                  {selected.map((val) => (
                    <Chip
                      key={val}
                      label={
                        val === OTHER_MAINTENANCE_TYPE && customType.trim()
                          ? `Outro: ${customType.trim()}`
                          : val
                      }
                      size="small"
                      color="primary"
                      variant="outlined"
                      sx={{ borderRadius: 1 }}
                    />
                  ))}
                </Box>
              )}
              disabled={loading}
            >
              {MAINTENANCE_TYPES.map((name) => (
                <MenuItem key={name} value={name}>
                  <Checkbox checked={types.indexOf(name) > -1} size="small" />
                  <ListItemText primary={name} />
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          {/* Campo condicional para Outro tipo de manutenção */}
          {types.includes(OTHER_MAINTENANCE_TYPE) && (
            <TextField
              fullWidth
              label="Especifique o outro tipo de manutenção"
              placeholder="Ex: Troca de correia dentada, junta do cabeçote..."
              value={customType}
              onChange={(e) => handleCustomTypeChange(e.target.value)}
              required
              disabled={loading}
              autoFocus
              helperText="Informe o nome específico do serviço realizado."
            />
          )}

          {/* Título do Serviço */}
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
            disabled={loading}
            helperText="Nome do registro que aparecerá no histórico."
          />

          <Stack direction={{ xs: "column", sm: "row" }} spacing={2}>
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
          </Stack>

          {/* Custo Total (R$) - Opcional */}
          <TextField
            fullWidth
            label="Custo Total (opcional)"
            placeholder="0,00"
            type="number"
            value={cost}
            onChange={(e) => setCost(e.target.value)}
            disabled={loading}
            slotProps={{
              input: {
                startAdornment: (
                  <InputAdornment position="start">
                    <AttachMoneyRoundedIcon
                      sx={{ fontSize: 18, color: "text.secondary" }}
                    />
                    <Typography
                      variant="body2"
                      sx={{ fontWeight: 600, mr: 0.5 }}
                    >
                      R$
                    </Typography>
                  </InputAdornment>
                ),
              },
              htmlInput: { min: 0, step: "0.01" },
            }}
            helperText="Valor total das peças e mão de obra do serviço."
          />

          {/* Descrição / Detalhes */}
          <TextField
            fullWidth
            label="Descrição / Detalhes (opcional)"
            placeholder="Ex: Utilizado óleo sintético 0W20 marca X, substituição do bujão e arruela..."
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            multiline
            rows={2.5}
            disabled={loading}
          />

          {/* Área de Anexo de Comprovantes */}
          <Box>
            <Stack
              direction="row"
              spacing={1}
              sx={{ alignItems: "center", mb: 1 }}
            >
              <ReceiptLongRoundedIcon
                sx={{ fontSize: 18, color: "primary.main" }}
              />
              <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>
                Comprovantes e Recibos (opcional)
              </Typography>
            </Stack>

            <Typography
              variant="caption"
              sx={{ color: "text.secondary", display: "block", mb: 1.5 }}
            >
              Anexe notas fiscais, ordens de serviço ou fotos das peças. Aceita
              PDF (até 2MB) ou imagens JPG, PNG, WebP (até 10MB com otimização
              automática).
            </Typography>

            <input
              ref={fileInputRef}
              type="file"
              multiple
              accept="image/jpeg,image/png,image/webp,application/pdf"
              style={{ display: "none" }}
              onChange={handleFilesSelected}
              disabled={loading || processingFiles}
            />

            <Button
              variant="outlined"
              startIcon={
                processingFiles ? (
                  <CircularProgress size={16} color="inherit" />
                ) : (
                  <CloudUploadRoundedIcon />
                )
              }
              onClick={() => fileInputRef.current?.click()}
              disabled={loading || processingFiles}
              sx={{ borderRadius: 2, textTransform: "none", py: 0.8 }}
            >
              {processingFiles
                ? "Otimizando comprovantes..."
                : "Anexar Comprovantes"}
            </Button>

            {/* Listagem de comprovantes anexados */}
            {attachments.length > 0 && (
              <Stack spacing={1} sx={{ mt: 1.5 }}>
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
                              fontSize: 24,
                              flexShrink: 0,
                            }}
                          />
                        ) : (
                          <ImageRoundedIcon
                            sx={{
                              color: "primary.main",
                              fontSize: 24,
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
                              maxWidth: { xs: 180, sm: 280 },
                            }}
                          >
                            {att.name}
                          </Typography>
                          <Typography
                            variant="caption"
                            sx={{ color: "text.secondary" }}
                          >
                            {formatFileSize(att.size)} •{" "}
                            {isPdf ? "Documento PDF" : "Imagem Otimizada (90%)"}
                          </Typography>
                        </Box>
                      </Stack>

                      <Tooltip title="Remover comprovante">
                        <IconButton
                          size="small"
                          color="error"
                          onClick={() => handleRemoveAttachment(att.id)}
                          disabled={loading || processingFiles}
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
          </Box>
        </Stack>
      </DialogContent>

      <DialogActions sx={{ p: 2.5, gap: 1 }}>
        <Button
          onClick={handleClose}
          variant="outlined"
          color="inherit"
          disabled={loading || processingFiles}
        >
          Cancelar
        </Button>
        <Button
          onClick={() => handleSubmit()}
          variant="contained"
          disabled={loading || processingFiles}
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
