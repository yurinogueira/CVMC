import React, { useState } from "react";
import {
  Link as RouterLink,
  useSearchParams,
  useNavigate,
} from "react-router-dom";
import {
  Box,
  Card,
  CardContent,
  Typography,
  TextField,
  Button,
  Stack,
  Alert,
  CircularProgress,
  InputAdornment,
  IconButton,
  Link,
} from "@mui/material";
import Visibility from "@mui/icons-material/Visibility";
import VisibilityOff from "@mui/icons-material/VisibilityOff";
import LockOutlinedIcon from "@mui/icons-material/LockOutlined";
import ArrowBackRoundedIcon from "@mui/icons-material/ArrowBackRounded";
import CheckCircleOutlineRoundedIcon from "@mui/icons-material/CheckCircleOutlineRounded";
import DirectionsCarIcon from "@mui/icons-material/DirectionsCar";
import { AuthHeroBanner } from "../components/AuthHeroBanner";
import { authService } from "../services/auth.service";
import { useDocumentTitle } from "../../shared";
import { brandColors } from "../../../styles/theme";

export function ResetPasswordPage() {
  useDocumentTitle("Redefinir Senha");
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const token = searchParams.get("token") || "";

  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMsg(null);

    if (!token.trim()) {
      setErrorMsg(
        "Token de recuperação inválido ou ausente. Solicite um novo link.",
      );
      return;
    }

    if (password.length < 8 || password.length > 72) {
      setErrorMsg("A nova senha deve ter entre 8 e 72 caracteres.");
      return;
    }

    if (password !== confirmPassword) {
      setErrorMsg("As senhas informadas não coincidem.");
      return;
    }

    try {
      setLoading(true);
      await authService.resetPassword({
        token: token.trim(),
        newPassword: password,
      });
      setSuccess(true);
    } catch (err: unknown) {
      const errorObj = err as { response?: { data?: { message?: string } } };
      setErrorMsg(
        errorObj.response?.data?.message ||
          "Não foi possível redefinir a senha. O token pode ter expirado.",
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box
      sx={{
        display: "flex",
        minHeight: "100vh",
        bgcolor: "background.default",
      }}
    >
      {/* Visual Left Banner (Desktop) */}
      <Box
        sx={{
          display: { xs: "none", md: "flex" },
          flex: { md: "0 0 50%", lg: "0 0 52%" },
        }}
      >
        <AuthHeroBanner />
      </Box>

      {/* Form Right Panel */}
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
          <CardContent sx={{ p: { xs: 1, sm: 3 } }}>
            {/* Mobile Header / Brand */}
            <Box
              sx={{
                display: { xs: "flex", md: "none" },
                alignItems: "center",
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

            <Box sx={{ mb: 3.5 }}>
              <Typography
                component="h1"
                variant="h4"
                sx={{ fontWeight: 700, mb: 1, color: "text.primary" }}
              >
                Nova Senha
              </Typography>
              <Typography variant="body2" sx={{ color: "text.secondary" }}>
                Defina uma nova senha de acesso para sua conta.
              </Typography>
            </Box>

            {!token && (
              <Alert severity="warning" sx={{ mb: 3, borderRadius: 1.5 }}>
                Token de recuperação não identificado na URL. Solicite um novo
                link através da página de recuperação.
              </Alert>
            )}

            {errorMsg && (
              <Alert severity="error" sx={{ mb: 3, borderRadius: 1.5 }}>
                {errorMsg}
              </Alert>
            )}

            {success ? (
              <Box sx={{ textAlign: "center", py: 2 }}>
                <CheckCircleOutlineRoundedIcon
                  sx={{ fontSize: 48, color: "success.main", mb: 2 }}
                />
                <Typography variant="body1" sx={{ fontWeight: 600, mb: 1 }}>
                  Senha alterada com sucesso!
                </Typography>
                <Typography
                  variant="body2"
                  sx={{ color: "text.secondary", mb: 3 }}
                >
                  Você já pode realizar o login com suas novas credenciais.
                </Typography>
                <Button
                  onClick={() => navigate("/login", { replace: true })}
                  variant="contained"
                  fullWidth
                >
                  Acessar Minha Conta
                </Button>
              </Box>
            ) : (
              <Box component="form" onSubmit={handleSubmit} noValidate>
                <Stack spacing={2.5}>
                  <TextField
                    fullWidth
                    label="Nova Senha"
                    type={showPassword ? "text" : "password"}
                    placeholder="Mínimo 8 caracteres"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    disabled={loading || !token}
                    required
                    helperText="Entre 8 e 72 caracteres"
                    slotProps={{
                      htmlInput: { maxLength: 72 },
                      input: {
                        startAdornment: (
                          <InputAdornment position="start">
                            <LockOutlinedIcon
                              sx={{ color: "text.secondary", fontSize: 20 }}
                            />
                          </InputAdornment>
                        ),
                        endAdornment: (
                          <InputAdornment position="end">
                            <IconButton
                              aria-label="alternar visibilidade da senha"
                              onClick={() => setShowPassword((prev) => !prev)}
                              edge="end"
                              size="small"
                            >
                              {showPassword ? (
                                <VisibilityOff fontSize="small" />
                              ) : (
                                <Visibility fontSize="small" />
                              )}
                            </IconButton>
                          </InputAdornment>
                        ),
                      },
                    }}
                  />

                  <TextField
                    fullWidth
                    label="Confirmar Nova Senha"
                    type={showPassword ? "text" : "password"}
                    placeholder="Repita a nova senha"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    disabled={loading || !token}
                    required
                    slotProps={{
                      htmlInput: { maxLength: 72 },
                      input: {
                        startAdornment: (
                          <InputAdornment position="start">
                            <LockOutlinedIcon
                              sx={{ color: "text.secondary", fontSize: 20 }}
                            />
                          </InputAdornment>
                        ),
                      },
                    }}
                  />

                  <Button
                    type="submit"
                    fullWidth
                    variant="contained"
                    size="large"
                    disabled={loading || !token}
                    sx={{
                      py: 1.4,
                      fontSize: "0.95rem",
                      fontWeight: 600,
                    }}
                  >
                    {loading ? (
                      <CircularProgress size={24} sx={{ color: "#FFFFFF" }} />
                    ) : (
                      "Salvar Nova Senha"
                    )}
                  </Button>
                </Stack>
              </Box>
            )}

            {!success && (
              <Box sx={{ mt: 3.5, textAlign: "center" }}>
                <Link
                  component={RouterLink}
                  to="/login"
                  sx={{
                    display: "inline-flex",
                    alignItems: "center",
                    gap: 0.5,
                    color: "primary.main",
                    fontWeight: 600,
                    fontSize: "0.875rem",
                    textDecoration: "none",
                    "&:hover": { textDecoration: "underline" },
                  }}
                >
                  <ArrowBackRoundedIcon fontSize="small" />
                  Voltar para o Login
                </Link>
              </Box>
            )}
          </CardContent>
        </Card>
      </Box>
    </Box>
  );
}
