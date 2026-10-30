import React from "react";
import {
  Box,
  Stack,
  Typography,
  Button,
  Skeleton,
  SxProps,
  Theme,
} from "@mui/material";
import ArrowBackRoundedIcon from "@mui/icons-material/ArrowBackRounded";

export interface PageHeaderBackButtonConfig {
  onClick: () => void;
  label?: string;
  ariaLabel?: string;
}

export interface PageHeaderProps {
  /**
   * Título principal da página. Renderizado canonicamente como tag <h1> semântica
   * com variant="h5", peso 800 e tamanho responsivo (24px a 28px).
   */
  title: React.ReactNode;
  /**
   * Subtítulo explicativo opcional. Renderizado como variant="body2" com cor text.secondary.
   */
  subtitle?: React.ReactNode;
  /**
   * Botão de voltar opcional. Pode ser um objeto de configuração `{ onClick, label }`
   * ou um nó React customizado.
   */
  backButton?: PageHeaderBackButtonConfig | React.ReactNode;
  /**
   * Bloco de ações à direita (botões, filtros, grupos de botões).
   * Em dispositivos móveis, empilha responsivamente abaixo do título.
   */
  action?: React.ReactNode;
  /**
   * Trilha de navegação (breadcrumbs) opcional renderizada acima do cabeçalho.
   */
  breadcrumbs?: React.ReactNode;
  /**
   * Indica estado de carregamento do cabeçalho (renderiza Skeleton de altura correspondente).
   */
  loading?: boolean;
  /**
   * Margem inferior customizável (padrão: 3).
   */
  mb?: number | { xs?: number; sm?: number; md?: number };
  /**
   * Estilização adicional no container do cabeçalho.
   */
  sx?: SxProps<Theme>;
  /**
   * Identificador para testes de integração.
   */
  dataTestId?: string;
}

function isBackButtonConfig(
  backButton: PageHeaderBackButtonConfig | React.ReactNode,
): backButton is PageHeaderBackButtonConfig {
  return (
    typeof backButton === "object" &&
    backButton !== null &&
    "onClick" in backButton &&
    typeof (backButton as PageHeaderBackButtonConfig).onClick === "function" &&
    !React.isValidElement(backButton)
  );
}

export function PageHeader({
  title,
  subtitle,
  backButton,
  action,
  breadcrumbs,
  loading = false,
  mb = 3,
  sx,
  dataTestId = "page-header",
}: PageHeaderProps) {
  const renderedBackButton = React.useMemo(() => {
    if (!backButton) return null;

    if (isBackButtonConfig(backButton)) {
      return (
        <Button
          variant="outlined"
          size="small"
          startIcon={<ArrowBackRoundedIcon />}
          onClick={backButton.onClick}
          aria-label={backButton.ariaLabel || backButton.label || "Voltar"}
          sx={{
            borderRadius: 2,
            textTransform: "none",
            borderColor: "#E2E8F0",
            color: "text.secondary",
            fontWeight: 600,
            "&:hover": {
              borderColor: "#CBD5E1",
              bgcolor: "#F8FAFC",
            },
          }}
        >
          {backButton.label || "Voltar"}
        </Button>
      );
    }

    return backButton;
  }, [backButton]);

  const renderedTitle = React.useMemo(() => {
    if (loading) {
      return (
        <Skeleton
          width={220}
          height={36}
          data-testid="page-header-skeleton"
          sx={{ my: 0.5 }}
        />
      );
    }

    return (
      <Typography
        component="h1"
        variant="h5"
        sx={{
          fontWeight: 800,
          color: "text.primary",
          fontSize: { xs: "1.25rem", sm: "1.5rem" },
          letterSpacing: "-0.01em",
          wordBreak: "break-word",
          lineHeight: 1.25,
        }}
      >
        {title}
      </Typography>
    );
  }, [loading, title]);

  const renderedSubtitle = React.useMemo(() => {
    if (!subtitle) return null;

    if (typeof subtitle === "string") {
      return (
        <Typography
          variant="body2"
          sx={{ color: "text.secondary", mt: 0.5, lineHeight: 1.5 }}
        >
          {subtitle}
        </Typography>
      );
    }

    return <Box sx={{ mt: 0.5 }}>{subtitle}</Box>;
  }, [subtitle]);

  return (
    <Box
      component="header"
      data-testid={dataTestId}
      sx={{
        width: "100%",
        mb,
        ...sx,
      }}
    >
      {breadcrumbs && <Box sx={{ mb: 1.5 }}>{breadcrumbs}</Box>}

      <Stack
        direction={{ xs: "column", sm: "row" }}
        spacing={2}
        sx={{
          alignItems: { xs: "flex-start", sm: "center" },
          justifyContent: "space-between",
          gap: 1.5,
        }}
      >
        <Stack
          direction="row"
          spacing={1.5}
          sx={{
            alignItems: "center",
            minWidth: 0,
            flex: 1,
            width: { xs: "100%", sm: "auto" },
          }}
        >
          {renderedBackButton}
          <Box sx={{ minWidth: 0, flex: 1 }}>
            {renderedTitle}
            {renderedSubtitle}
          </Box>
        </Stack>

        {action && (
          <Box
            sx={{
              width: { xs: "100%", sm: "auto" },
              flexShrink: 0,
              display: "flex",
              alignItems: "center",
            }}
          >
            {action}
          </Box>
        )}
      </Stack>
    </Box>
  );
}
