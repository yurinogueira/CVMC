import { Card, CardContent, Typography, Box, Stack, Chip } from "@mui/material";
import BuildCircleRoundedIcon from "@mui/icons-material/BuildCircleRounded";
import CalendarMonthRoundedIcon from "@mui/icons-material/CalendarMonthRounded";
import SpeedRoundedIcon from "@mui/icons-material/SpeedRounded";
import { Maintenance } from "../types/maintenance.types";
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
      </CardContent>
    </Card>
  );
}
