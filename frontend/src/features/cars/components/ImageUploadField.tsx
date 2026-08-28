import React, { useRef, useState } from "react";
import {
  Box,
  Button,
  Stack,
  Typography,
  IconButton,
  TextField,
  Tabs,
  Tab,
  CircularProgress,
  Tooltip,
  Alert,
} from "@mui/material";
import CloudUploadRoundedIcon from "@mui/icons-material/CloudUploadRounded";
import LinkRoundedIcon from "@mui/icons-material/LinkRounded";
import DeleteOutlineRoundedIcon from "@mui/icons-material/DeleteOutlineRounded";
import { VehicleImage } from "./VehicleImage";

interface ImageUploadFieldProps {
  value?: string;
  onChange: (imageUrl: string) => void;
  vehicleType?: string;
  disabled?: boolean;
}

const MAX_FILE_SIZE_BYTES = 5 * 1024 * 1024; // 5MB

export function ImageUploadField({
  value = "",
  onChange,
  vehicleType = "cars",
  disabled = false,
}: ImageUploadFieldProps) {
  const [tab, setTab] = useState<0 | 1>(value.startsWith("http") ? 1 : 0);
  const [urlInput, setUrlInput] = useState(
    value.startsWith("http") ? value : "",
  );
  const [compressing, setCompressing] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setErrorMsg(null);

    // Validate mime type
    if (!["image/jpeg", "image/png", "image/webp"].includes(file.type)) {
      setErrorMsg("Formato não suportado. Utilize imagens JPG, PNG ou WebP.");
      return;
    }

    // Validate size (5MB)
    if (file.size > MAX_FILE_SIZE_BYTES) {
      setErrorMsg("Arquivo muito grande. O tamanho máximo permitido é 5MB.");
      return;
    }

    setCompressing(true);
    const reader = new FileReader();

    reader.onload = (event) => {
      const img = new Image();
      img.onload = () => {
        // Resize and compress via canvas
        const canvas = document.createElement("canvas");
        const maxDimension = 1280;
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
          const compressedDataUrl = canvas.toDataURL("image/jpeg", 0.85);
          onChange(compressedDataUrl);
        } else {
          // Fallback to original read
          onChange(event.target?.result as string);
        }
        setCompressing(false);
      };

      img.onerror = () => {
        setErrorMsg("Falha ao carregar e processar a imagem.");
        setCompressing(false);
      };

      img.src = event.target?.result as string;
    };

    reader.onerror = () => {
      setErrorMsg("Erro ao ler o arquivo selecionado.");
      setCompressing(false);
    };

    reader.readAsDataURL(file);
  };

  const handleUrlBlur = () => {
    const trimmed = urlInput.trim();
    if (trimmed && !trimmed.match(/^https?:\/\/.+/i)) {
      setErrorMsg("Informe uma URL HTTP ou HTTPS válida para a foto.");
      return;
    }
    setErrorMsg(null);
    onChange(trimmed);
  };

  const handleClear = () => {
    onChange("");
    setUrlInput("");
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
    setErrorMsg(null);
  };

  return (
    <Box sx={{ width: "100%" }}>
      <Typography
        variant="caption"
        sx={{
          fontWeight: 600,
          color: "text.secondary",
          textTransform: "uppercase",
          letterSpacing: 0.5,
          display: "block",
          mb: 1,
        }}
      >
        Foto do Veículo (Opcional)
      </Typography>

      {errorMsg && (
        <Alert
          severity="error"
          onClose={() => setErrorMsg(null)}
          sx={{ mb: 1.5, borderRadius: 1.5 }}
        >
          {errorMsg}
        </Alert>
      )}

      <Stack
        direction={{ xs: "column", sm: "row" }}
        spacing={2}
        sx={{ alignItems: "stretch" }}
      >
        {/* Preview Area */}
        <Box sx={{ width: { xs: "100%", sm: 200 }, flexShrink: 0 }}>
          <VehicleImage
            imageUrl={value}
            vehicleType={vehicleType}
            alt="Foto do veículo"
            aspectRatio="16 / 10"
            borderRadius={2}
          />
          {value && (
            <Button
              size="small"
              color="error"
              onClick={handleClear}
              startIcon={<DeleteOutlineRoundedIcon />}
              disabled={disabled || compressing}
              sx={{ mt: 1, width: "100%", textTransform: "none" }}
            >
              Remover Foto
            </Button>
          )}
        </Box>

        {/* Input Controls */}
        <Box sx={{ flex: 1, minWidth: 0 }}>
          <Tabs
            value={tab}
            onChange={(_, next) => setTab(next)}
            sx={{
              minHeight: 36,
              mb: 1.5,
              "& .MuiTab-root": {
                minHeight: 36,
                py: 0.5,
                fontSize: "0.82rem",
                textTransform: "none",
                fontWeight: 600,
              },
            }}
          >
            <Tab
              icon={<CloudUploadRoundedIcon sx={{ fontSize: 18 }} />}
              iconPosition="start"
              label="Fazer Upload"
            />
            <Tab
              icon={<LinkRoundedIcon sx={{ fontSize: 18 }} />}
              iconPosition="start"
              label="URL da Foto"
            />
          </Tabs>

          {tab === 0 ? (
            <Box
              sx={{
                p: 2,
                border: "2px dashed #CBD5E1",
                borderRadius: 2,
                textAlign: "center",
                bgcolor: "#F8FAFC",
                transition: "all 0.2s ease",
                "&:hover": {
                  borderColor: "primary.main",
                  bgcolor: "rgba(2, 132, 199, 0.04)",
                },
              }}
            >
              <input
                type="file"
                ref={fileInputRef}
                accept="image/jpeg,image/png,image/webp"
                onChange={handleFileSelect}
                style={{ display: "none" }}
                disabled={disabled || compressing}
              />

              <Stack spacing={1} sx={{ alignItems: "center" }}>
                <CloudUploadRoundedIcon
                  sx={{ fontSize: 32, color: "primary.main" }}
                />
                <Typography variant="body2" sx={{ fontWeight: 600 }}>
                  {compressing
                    ? "Otimizando imagem..."
                    : "Selecione uma foto do seu veículo"}
                </Typography>
                <Typography variant="caption" sx={{ color: "text.secondary" }}>
                  Formatos aceitos: JPG, PNG, WebP (Máx. 5MB)
                </Typography>

                <Button
                  variant="outlined"
                  size="small"
                  onClick={() => fileInputRef.current?.click()}
                  disabled={disabled || compressing}
                  startIcon={
                    compressing ? (
                      <CircularProgress size={16} />
                    ) : (
                      <CloudUploadRoundedIcon />
                    )
                  }
                  sx={{ mt: 1, textTransform: "none", borderRadius: 1.5 }}
                >
                  {value ? "Trocar Arquivo" : "Escolher Arquivo"}
                </Button>
              </Stack>
            </Box>
          ) : (
            <Stack spacing={1}>
              <TextField
                fullWidth
                size="small"
                label="URL da Imagem"
                placeholder="https://exemplo.com/foto-do-carro.jpg"
                value={urlInput}
                onChange={(e) => setUrlInput(e.target.value)}
                onBlur={handleUrlBlur}
                disabled={disabled}
                slotProps={{
                  input: {
                    endAdornment: urlInput ? (
                      <Tooltip title="Limpar URL">
                        <IconButton
                          size="small"
                          onClick={() => {
                            setUrlInput("");
                            onChange("");
                          }}
                        >
                          <DeleteOutlineRoundedIcon fontSize="small" />
                        </IconButton>
                      </Tooltip>
                    ) : null,
                  },
                }}
              />
              <Typography variant="caption" sx={{ color: "text.secondary" }}>
                Cole o link direto para a foto do seu veículo na internet.
              </Typography>
            </Stack>
          )}
        </Box>
      </Stack>
    </Box>
  );
}
