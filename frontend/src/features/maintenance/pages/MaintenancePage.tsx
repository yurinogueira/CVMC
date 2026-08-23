import { Box, Typography, Card, Button, Stack, Chip } from "@mui/material";
import BuildCircleRoundedIcon from "@mui/icons-material/BuildCircleRounded";
import AddRoundedIcon from "@mui/icons-material/AddRounded";
import CheckCircleRoundedIcon from "@mui/icons-material/CheckCircleRounded";
import { useDocumentTitle } from "../../shared";
import { brandColors } from "../../../styles/theme";

export function MaintenancePage() {
  useDocumentTitle("Histórico de Manutenções");
  return (
    <Box>
      {/* Header */}
      <Stack
        direction={{ xs: "column", sm: "row" }}
        spacing={2}
        sx={{
          alignItems: { xs: "flex-start", sm: "center" },
          justifyContent: "space-between",
          mb: 3,
        }}
      >
        <Box>
          <Typography
            component="h1"
            variant="h5"
            sx={{ fontWeight: 800, color: "text.primary" }}
          >
            Manutenções & Revisões
          </Typography>
          <Typography variant="body2" sx={{ color: "text.secondary" }}>
            Histórico completo de revisões, trocas de óleo, pneus e reparos.
          </Typography>
        </Box>

        <Button
          variant="contained"
          startIcon={<AddRoundedIcon />}
          sx={{ px: 2.5, py: 1.1 }}
        >
          Registrar Manutenção
        </Button>
      </Stack>

      {/* Info Status Banner */}
      <Card
        elevation={0}
        sx={{
          mb: 4,
          p: 3,
          bgcolor: "background.paper",
          border: "1px solid #E2E8F0",
          display: "flex",
          alignItems: "center",
          gap: 2,
        }}
      >
        <Box
          sx={{
            width: 48,
            height: 48,
            borderRadius: 2,
            bgcolor: "rgba(22, 163, 74, 0.12)",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            color: "#16A34A",
          }}
        >
          <CheckCircleRoundedIcon sx={{ fontSize: 28 }} />
        </Box>
        <Box sx={{ flex: 1 }}>
          <Stack
            direction="row"
            spacing={1}
            sx={{ alignItems: "center", mb: 0.5 }}
          >
            <Typography
              component="h2"
              variant="subtitle1"
              sx={{ fontWeight: 700 }}
            >
              Sua frota está em dia
            </Typography>
            <Chip
              label="Tudo Certo"
              size="small"
              sx={{
                bgcolor: brandColors.green,
                color: "#064E3B",
                fontWeight: 700,
                borderRadius: 1,
              }}
            />
          </Stack>
          <Typography variant="body2" sx={{ color: "text.secondary" }}>
            Nenhum alerta de manutenção atrasada ou revisão pendente no momento.
          </Typography>
        </Box>
      </Card>

      {/* Empty State */}
      <Card
        elevation={0}
        sx={{
          py: 8,
          px: 3,
          textAlign: "center",
          border: "1px dashed #CBD5E1",
          bgcolor: "background.paper",
        }}
      >
        <Box
          sx={{
            width: 56,
            height: 56,
            borderRadius: "50%",
            bgcolor: "rgba(2, 132, 199, 0.1)",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            mx: "auto",
            mb: 2,
            color: "primary.main",
          }}
        >
          <BuildCircleRoundedIcon sx={{ fontSize: 32 }} />
        </Box>
        <Typography
          component="h2"
          variant="h6"
          sx={{ fontWeight: 700, mb: 0.5 }}
        >
          Nenhuma manutenção registrada
        </Typography>
        <Typography
          variant="body2"
          sx={{ color: "text.secondary", maxWidth: 460, mx: "auto", mb: 3 }}
        >
          Registre trocas de óleo, pastilhas de freio, alinhamentos ou revisões
          para acompanhar custos e emitir relatórios.
        </Typography>
        <Button variant="contained" startIcon={<AddRoundedIcon />}>
          Nova Manutenção
        </Button>
      </Card>
    </Box>
  );
}
