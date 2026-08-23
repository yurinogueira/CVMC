import { Box, CircularProgress } from "@mui/material";

export function PageLoadingFallback() {
  return (
    <Box
      sx={{
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        minHeight: "60vh",
        width: "100%",
      }}
      aria-busy="true"
      aria-label="Carregando conteúdo da página"
    >
      <CircularProgress size={40} thickness={4} color="primary" />
    </Box>
  );
}
