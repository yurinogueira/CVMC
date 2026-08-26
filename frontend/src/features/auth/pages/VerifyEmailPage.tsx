import { useState, useEffect } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import {
  Box,
  Card,
  CardContent,
  Typography,
  Button,
  CircularProgress,
  Alert,
} from "@mui/material";
import CheckCircleOutlineRoundedIcon from "@mui/icons-material/CheckCircleOutlineRounded";
import ErrorOutlineRoundedIcon from "@mui/icons-material/ErrorOutlineRounded";
import DirectionsCarIcon from "@mui/icons-material/DirectionsCar";
import { AuthHeroBanner } from "../components/AuthHeroBanner";
import { authService } from "../services/auth.service";
import { useAuthStore } from "../state/auth.store";
import { useDocumentTitle } from "../../shared";
import { brandColors } from "../../../styles/theme";

export function VerifyEmailPage() {
  useDocumentTitle("Confirmação de E-mail");
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const token = searchParams.get("token") || "";
  const hasToken = Boolean(token.trim());

  const { user, setUser, isAuthenticated } = useAuthStore();

  const [loading, setLoading] = useState(hasToken);
  const [success, setSuccess] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(
    hasToken
      ? null
      : "Token de validação ausente. Utilize o link enviado para seu e-mail.",
  );

  useEffect(() => {
    if (!hasToken) return;

    let isMounted = true;
    authService
      .verifyEmail({ token: token.trim() })
      .then(() => {
        if (isMounted) {
          setSuccess(true);
          if (user) {
            setUser({ ...user, emailVerified: true });
          }
        }
      })
      .catch((err: unknown) => {
        if (isMounted) {
          const errorObj = err as {
            response?: { data?: { message?: string } };
          };
          setErrorMsg(
            errorObj.response?.data?.message ||
              "Não foi possível validar seu e-mail. O link pode ter expirado ou já ter sido utilizado.",
          );
        }
      })
      .finally(() => {
        if (isMounted) {
          setLoading(false);
        }
      });

    return () => {
      isMounted = false;
    };
  }, [hasToken, token, user, setUser]);

  return (
    <Box
      sx={{
        display: "flex",
        minHeight: "100vh",
        bgcolor: "background.default",
      }}
    >
      {/* Visual Left Banner */}
      <Box
        sx={{
          display: { xs: "none", md: "flex" },
          flex: { md: "0 0 50%", lg: "0 0 52%" },
        }}
      >
        <AuthHeroBanner />
      </Box>

      {/* Content Right Panel */}
      <Box
        sx={{
          flex: 1,
          display: "flex",
          flexDirection: "column",
          justifyContent: "center",
          alignItems: "center",
          p: { xs: 3, sm: 6, md: 8 },
        }}
      >
        <Card
          elevation={0}
          sx={{
            width: "100%",
            maxWidth: 440,
            border: { xs: "none", sm: "1px solid #E2E8F0" },
            bgcolor: { xs: "transparent", sm: "background.paper" },
            p: { xs: 0, sm: 2 },
          }}
        >
          <CardContent sx={{ p: { xs: 1, sm: 3 }, textAlign: "center" }}>
            {/* Mobile Header / Brand */}
            <Box
              sx={{
                display: { xs: "flex", md: "none" },
                alignItems: "center",
                justifyContent: "center",
                mb: 3,
                gap: 1,
              }}
            >
              <Box
                sx={{
                  bgcolor: brandColors.blue,
                  borderRadius: 1.5,
                  p: 0.75,
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                }}
              >
                <DirectionsCarIcon sx={{ fontSize: 24, color: "#FFFFFF" }} />
              </Box>
              <Typography
                variant="h6"
                sx={{ fontWeight: 800, color: "primary.main" }}
              >
                CVMC
              </Typography>
            </Box>

            {loading && (
              <Box sx={{ py: 6 }}>
                <CircularProgress size={48} sx={{ mb: 3 }} />
                <Typography variant="h6" sx={{ fontWeight: 600, mb: 1 }}>
                  Validando seu e-mail...
                </Typography>
                <Typography variant="body2" sx={{ color: "text.secondary" }}>
                  Aguarde um instante enquanto confirmamos seu endereço de
                  e-mail.
                </Typography>
              </Box>
            )}

            {!loading && success && (
              <Box sx={{ py: 3 }}>
                <CheckCircleOutlineRoundedIcon
                  sx={{ fontSize: 64, color: "success.main", mb: 2 }}
                />
                <Typography variant="h5" sx={{ fontWeight: 700, mb: 1 }}>
                  E-mail Validado com Sucesso!
                </Typography>
                <Typography
                  variant="body2"
                  sx={{ color: "text.secondary", mb: 4 }}
                >
                  Sua conta foi confirmada. Agora você tem acesso liberado para
                  cadastrar e gerenciar todos os seus veículos no CVMC.
                </Typography>
                <Button
                  onClick={() =>
                    navigate(isAuthenticated ? "/dashboard" : "/login", {
                      replace: true,
                    })
                  }
                  variant="contained"
                  fullWidth
                  size="large"
                  sx={{ py: 1.4, fontWeight: 600 }}
                >
                  {isAuthenticated ? "Ir para o Painel" : "Fazer Login"}
                </Button>
              </Box>
            )}

            {!loading && !success && (
              <Box sx={{ py: 3 }}>
                <ErrorOutlineRoundedIcon
                  sx={{ fontSize: 64, color: "error.main", mb: 2 }}
                />
                <Typography variant="h5" sx={{ fontWeight: 700, mb: 1 }}>
                  Falha na Validação
                </Typography>
                {errorMsg && (
                  <Alert severity="error" sx={{ mb: 3, textAlign: "left" }}>
                    {errorMsg}
                  </Alert>
                )}
                <Button
                  onClick={() =>
                    navigate(isAuthenticated ? "/profile" : "/login", {
                      replace: true,
                    })
                  }
                  variant="outlined"
                  fullWidth
                  size="large"
                  sx={{ py: 1.4, fontWeight: 600 }}
                >
                  {isAuthenticated ? "Ir para o Meu Perfil" : "Ir para o Login"}
                </Button>
              </Box>
            )}
          </CardContent>
        </Card>
      </Box>
    </Box>
  );
}
