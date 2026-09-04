import {
  Card,
  CardContent,
  Typography,
  Box,
  Stack,
  Chip,
  Button,
} from "@mui/material";
import BuildCircleRoundedIcon from "@mui/icons-material/BuildCircleRounded";
import CalendarMonthRoundedIcon from "@mui/icons-material/CalendarMonthRounded";
import SpeedRoundedIcon from "@mui/icons-material/SpeedRounded";
import AttachMoneyRoundedIcon from "@mui/icons-material/AttachMoneyRounded";
import PictureAsPdfRoundedIcon from "@mui/icons-material/PictureAsPdfRounded";
import ImageRoundedIcon from "@mui/icons-material/ImageRounded";
import ReceiptLongRoundedIcon from "@mui/icons-material/ReceiptLongRounded";
import { Maintenance, MaintenanceAttachment } from "../types/maintenance.types";
import { brandColors } from "../../../styles/theme";

interface MaintenanceCardProps {
  maintenance: Maintenance;
}

export function MaintenanceCard({ maintenance }: MaintenanceCardProps) {
  const formatDate = (dateStr: string): string => {
    try {
      const datePart = dateStr.split("T")[0];
      const parts = datePart.split("-");
      if (parts.length === 3) {
        return `${parts[2]}/${parts[1]}/${parts[0]}`;
      }
      return new Date(dateStr).toLocaleDateString("pt-BR");
    } catch {
      return dateStr;
    }
  };

  const formatCost = (val: number): string => {
    return new Intl.NumberFormat("pt-BR", {
      style: "currency",
      currency: "BRL",
    }).format(val);
  };

  const handleOpenAttachment = (att: MaintenanceAttachment) => {
    if (!att.dataUrl) return;

    if (att.dataUrl.startsWith("data:application/pdf")) {
      // Abre PDF em uma nova aba convertendo dataURL em blob para compatibilidade ampla
      const byteString = atob(att.dataUrl.split(",")[1]);
      const ab = new ArrayBuffer(byteString.length);
      const ia = new Uint8Array(ab);
      for (let i = 0; i < byteString.length; i++) {
        ia[i] = byteString.charCodeAt(i);
      }
      const blob = new Blob([ab], { type: "application/pdf" });
      const blobUrl = URL.createObjectURL(blob);
      window.open(blobUrl, "_blank");
    } else {
      // Imagens ou links diretos
      const win = window.open();
      if (win) {
        win.document.write(
          `<img src="${att.dataUrl}" style="max-width:100%; height:auto;" alt="${att.name}" />`,
        );
        win.document.title = att.name;
      }
    }
  };

  return (
    <Card
      elevation={0}
      sx={{
        borderRadius: 2.5,
        border: "1px solid #E2E8F0",
        bgcolor: "background.paper",
        transition: "transform 0.2s ease, box-shadow 0.2s ease",
        "&:hover": {
          transform: "translateY(-2px)",
          boxShadow: "0 8px 24px -4px rgba(2, 132, 199, 0.08)",
        },
      }}
    >
      <CardContent sx={{ p: 2.5 }}>
        <Stack
          direction={{ xs: "column", sm: "row" }}
          spacing={2}
          sx={{
            alignItems: { xs: "flex-start", sm: "center" },
            justifyContent: "space-between",
            mb: 1.5,
          }}
        >
          <Stack direction="row" spacing={1.5} sx={{ alignItems: "center" }}>
            <Box
              sx={{
                width: 42,
                height: 42,
                borderRadius: 2,
                bgcolor: "rgba(2, 132, 199, 0.1)",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                color: "primary.main",
                flexShrink: 0,
              }}
            >
              <BuildCircleRoundedIcon sx={{ fontSize: 24 }} />
            </Box>
            <Box>
              <Typography
                variant="subtitle1"
                component="h3"
                sx={{ fontWeight: 800, color: "text.primary", lineHeight: 1.3 }}
              >
                {maintenance.title}
              </Typography>
            </Box>
          </Stack>

          <Stack direction="row" spacing={1} sx={{ flexWrap: "wrap", gap: 1 }}>
            {maintenance.cost !== undefined && maintenance.cost !== null && (
              <Chip
                icon={
                  <AttachMoneyRoundedIcon
                    sx={{ "&&": { fontSize: 16, color: "#1E3A8A" } }}
                  />
                }
                label={formatCost(maintenance.cost)}
                size="small"
                sx={{
                  bgcolor: "#EFF6FF",
                  color: "#1E3A8A",
                  fontWeight: 700,
                  borderRadius: 1,
                  border: "1px solid #BFDBFE",
                }}
              />
            )}
            <Chip
              icon={
                <CalendarMonthRoundedIcon sx={{ "&&": { fontSize: 16 } }} />
              }
              label={formatDate(maintenance.date)}
              size="small"
              sx={{
                bgcolor: "#F1F5F9",
                color: "text.primary",
                fontWeight: 500,
                borderRadius: 1,
              }}
            />
            <Chip
              icon={
                <SpeedRoundedIcon
                  sx={{ "&&": { fontSize: 16, color: "#064E3B" } }}
                />
              }
              label={`${(maintenance.mileage || 0).toLocaleString("pt-BR")} km`}
              size="small"
              sx={{
                bgcolor: brandColors.mint,
                color: "#064E3B",
                fontWeight: 600,
                borderRadius: 1,
              }}
            />
          </Stack>
        </Stack>

        {/* Tipos de Manutenção */}
        {maintenance.types && maintenance.types.length > 0 && (
          <Stack
            direction="row"
            spacing={0.8}
            sx={{ flexWrap: "wrap", gap: 0.8, my: 1.5 }}
          >
            {maintenance.types.map((type, idx) => (
              <Chip
                key={`${type}-${idx}`}
                label={type}
                size="small"
                variant="outlined"
                sx={{
                  borderRadius: 1,
                  fontSize: "0.75rem",
                  borderColor: "#CBD5E1",
                  bgcolor: "#F8FAFC",
                  color: "text.primary",
                }}
              />
            ))}
          </Stack>
        )}

        {/* Descrição */}
        {maintenance.description && (
          <Typography
            variant="body2"
            sx={{
              color: "text.secondary",
              bgcolor: "#F8FAFC",
              p: 1.5,
              borderRadius: 1.5,
              border: "1px solid #F1F5F9",
              mt: 1,
              whiteSpace: "pre-line",
            }}
          >
            {maintenance.description}
          </Typography>
        )}

        {/* Comprovantes Anexados */}
        {maintenance.attachments && maintenance.attachments.length > 0 && (
          <Box sx={{ mt: 2, pt: 1.5, borderTop: "1px dashed #E2E8F0" }}>
            <Stack
              direction="row"
              spacing={1}
              sx={{ alignItems: "center", mb: 1 }}
            >
              <ReceiptLongRoundedIcon
                sx={{ fontSize: 16, color: "text.secondary" }}
              />
              <Typography
                variant="caption"
                sx={{
                  fontWeight: 700,
                  color: "text.secondary",
                  textTransform: "uppercase",
                }}
              >
                Comprovantes ({maintenance.attachments.length})
              </Typography>
            </Stack>

            <Stack
              direction="row"
              spacing={1}
              sx={{ flexWrap: "wrap", gap: 1 }}
            >
              {maintenance.attachments.map((att) => {
                const isPdf = att.mimeType === "application/pdf";
                return (
                  <Button
                    key={att.id}
                    variant="outlined"
                    size="small"
                    startIcon={
                      isPdf ? (
                        <PictureAsPdfRoundedIcon sx={{ color: "#EF4444" }} />
                      ) : (
                        <ImageRoundedIcon color="primary" />
                      )
                    }
                    onClick={() => handleOpenAttachment(att)}
                    sx={{
                      textTransform: "none",
                      fontSize: "0.78rem",
                      borderRadius: 1.5,
                      borderColor: "#CBD5E1",
                      bgcolor: "background.paper",
                      color: "text.primary",
                      py: 0.4,
                      px: 1,
                      "&:hover": {
                        borderColor: "primary.main",
                        bgcolor: "rgba(2, 132, 199, 0.04)",
                      },
                    }}
                  >
                    {att.name}
                  </Button>
                );
              })}
            </Stack>
          </Box>
        )}
      </CardContent>
    </Card>
  );
}
