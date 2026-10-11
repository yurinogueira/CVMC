import React, { useState, useEffect } from "react";
import {
  Box,
  Card,
  CardContent,
  Typography,
  Grid,
  TextField,
  Button,
  Stack,
  Alert,
  Chip,
  LinearProgress,
  Divider,
  InputAdornment,
  IconButton,
  CircularProgress,
} from "@mui/material";
import VerifiedUserRoundedIcon from "@mui/icons-material/VerifiedUserRounded";
import GppMaybeRoundedIcon from "@mui/icons-material/GppMaybeRounded";
import MarkEmailReadRoundedIcon from "@mui/icons-material/MarkEmailReadRounded";
import DirectionsCarRoundedIcon from "@mui/icons-material/DirectionsCarRounded";
import PersonRoundedIcon from "@mui/icons-material/PersonRounded";
import LockResetRoundedIcon from "@mui/icons-material/LockResetRounded";
import Visibility from "@mui/icons-material/Visibility";
import VisibilityOff from "@mui/icons-material/VisibilityOff";
import SaveRoundedIcon from "@mui/icons-material/SaveRounded";
import SendRoundedIcon from "@mui/icons-material/SendRounded";
import LinkRoundedIcon from "@mui/icons-material/LinkRounded";
import LinkOffRoundedIcon from "@mui/icons-material/LinkOffRounded";

import { profileService } from "../services/profile.service";
import { authService } from "../../auth/services/auth.service";
import { useAuthStore } from "../../auth/state/auth.store";
import { AuthProvidersData } from "../../auth/types/auth.types";
import {
  GoogleIcon,
  GoogleSignInButton,
} from "../../auth/components/GoogleSignInButton";
import { ProfileData } from "../types/profile.types";
import { useDocumentTitle, PageLoadingFallback } from "../../shared";

export function ProfilePage() {
  useDocumentTitle("Meu Perfil");
  const user = useAuthStore((state) => state.user);
  const setUser = useAuthStore((state) => state.setUser);

  const [profileData, setProfileData] = useState<ProfileData | null>(null);
  const [loading, setLoading] = useState(true);
  const [pageError, setPageError] = useState<string | null>(null);

  // Name form state
  const [name, setName] = useState("");
  const [updatingName, setUpdatingName] = useState(false);
  const [nameSuccess, setNameSuccess] = useState<string | null>(null);
  const [nameError, setNameError] = useState<string | null>(null);

  // Password form state
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showCurrentPassword, setShowCurrentPassword] = useState(false);
  const [showNewPassword, setShowNewPassword] = useState(false);
  const [updatingPassword, setUpdatingPassword] = useState(false);
  const [passwordSuccess, setPasswordSuccess] = useState<string | null>(null);
  const [passwordError, setPasswordError] = useState<string | null>(null);

  // Auth providers state
  const [authProviders, setAuthProviders] = useState<AuthProvidersData | null>(
    null,
  );
  const [unlinkingGoogle, setUnlinkingGoogle] = useState(false);
  const [providerSuccess, setProviderSuccess] = useState<string | null>(null);
  const [providerError, setProviderError] = useState<string | null>(null);

  // Resend email verification state
  const [resendingEmail, setResendingEmail] = useState(false);
  const [emailSuccess, setEmailSuccess] = useState<string | null>(null);
  const [emailError, setEmailError] = useState<string | null>(null);
  const [cooldown, setCooldown] = useState(0);

  useEffect(() => {
    let isMounted = true;
    Promise.all([
      profileService.getProfile(),
      profileService.getAuthProviders().catch(() => null),
    ])
      .then(([data, providers]) => {
        if (isMounted) {
          setProfileData(data);
          setName(data.user.name);
          if (providers) {
            setAuthProviders(providers);
          }
          const currentUser = useAuthStore.getState().user;
          if (currentUser) {
            setUser({ ...currentUser, ...data.user });
          } else {
            setUser(data.user);
          }
        }
      })
      .catch(() => {
        if (isMounted) {
          setPageError("Não foi possível carregar as informações do perfil.");
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
  }, [setUser]);

  // Cooldown countdown timer
  useEffect(() => {
    if (cooldown <= 0) return;
    const timer = setInterval(() => {
      setCooldown((prev) => prev - 1);
    }, 1000);
    return () => clearInterval(timer);
  }, [cooldown]);

  const handleUpdateName = async (e: React.FormEvent) => {
    e.preventDefault();
    setNameSuccess(null);
    setNameError(null);

    if (name.trim().length < 2 || name.trim().length > 100) {
      setNameError("O nome deve ter entre 2 e 100 caracteres.");
      return;
    }

    try {
      setUpdatingName(true);
      const updatedUser = await profileService.updateProfile({
        name: name.trim(),
      });
      setNameSuccess("Nome atualizado com sucesso!");
      if (user) {
        setUser({ ...user, name: updatedUser.name });
      }
      if (profileData) {
        setProfileData({
          ...profileData,
          user: { ...profileData.user, name: updatedUser.name },
        });
      }
    } catch (err: unknown) {
      const errorObj = err as { response?: { data?: { message?: string } } };
      setNameError(
        errorObj.response?.data?.message || "Erro ao atualizar o nome.",
      );
    } finally {
      setUpdatingName(false);
    }
  };

  const handleLinkGoogle = async (idToken: string) => {
    setProviderSuccess(null);
    setProviderError(null);
    try {
      const res = await profileService.linkGoogle(idToken);
      setProviderSuccess(res.message || "Conta Google vinculada com sucesso!");
      const updatedProviders = await profileService.getAuthProviders();
      setAuthProviders(updatedProviders);
      const updatedProfile = await profileService.getProfile();
      setProfileData(updatedProfile);
    } catch (err: unknown) {
      const errorObj = err as { response?: { data?: { message?: string } } };
      setProviderError(
        errorObj.response?.data?.message || "Erro ao vincular conta Google.",
      );
    }
  };

  const handleUnlinkGoogle = async () => {
    setProviderSuccess(null);
    setProviderError(null);
    try {
      setUnlinkingGoogle(true);
      const res = await profileService.unlinkGoogle();
      setProviderSuccess(
        res.message || "Conta Google desvinculada com sucesso!",
      );
      const updatedProviders = await profileService.getAuthProviders();
      setAuthProviders(updatedProviders);
      const updatedProfile = await profileService.getProfile();
      setProfileData(updatedProfile);
    } catch (err: unknown) {
      const errorObj = err as { response?: { data?: { message?: string } } };
      setProviderError(
        errorObj.response?.data?.message || "Erro ao desvincular conta Google.",
      );
    } finally {
      setUnlinkingGoogle(false);
    }
  };

  const handleUpdatePassword = async (e: React.FormEvent) => {
    e.preventDefault();
    setPasswordSuccess(null);
    setPasswordError(null);

    const hasPassword = authProviders ? authProviders.hasPassword : true;

    if (hasPassword && !currentPassword) {
      setPasswordError("Informe sua senha atual.");
      return;
    }

    if (newPassword.length < 8 || newPassword.length > 72) {
      setPasswordError("A nova senha deve ter entre 8 e 72 caracteres.");
      return;
    }

    if (newPassword !== confirmPassword) {
      setPasswordError("A confirmação de senha não coincide com a nova senha.");
      return;
    }

    try {
      setUpdatingPassword(true);
      if (!hasPassword) {
        const res = await profileService.setPassword(newPassword);
        setPasswordSuccess(res.message || "Senha cadastrada com sucesso!");
        const updatedProviders = await profileService.getAuthProviders();
        setAuthProviders(updatedProviders);
      } else {
        const res = await profileService.updatePassword({
          currentPassword,
          newPassword,
        });
        setPasswordSuccess(res.message || "Senha alterada com sucesso!");
      }
      setCurrentPassword("");
      setNewPassword("");
      setConfirmPassword("");
    } catch (err: unknown) {
      const errorObj = err as { response?: { data?: { message?: string } } };
      setPasswordError(
        errorObj.response?.data?.message || "Erro ao atualizar a senha.",
      );
    } finally {
      setUpdatingPassword(false);
    }
  };

  const handleResendVerification = async () => {
    setEmailSuccess(null);
    setEmailError(null);

    try {
      setResendingEmail(true);
      const res = await authService.resendVerification();
      setEmailSuccess(
        res.message || "E-mail de confirmação reenviado com sucesso!",
      );
      setCooldown(60); // 60s cooldown
    } catch (err: unknown) {
      const errorObj = err as { response?: { data?: { message?: string } } };
      setEmailError(
        errorObj.response?.data?.message ||
          "Não foi possível reenviar o e-mail de confirmação.",
      );
    } finally {
      setResendingEmail(false);
    }
  };

  if (loading) {
    return <PageLoadingFallback />;
  }

  const isEmailVerified = Boolean(profileData?.user?.emailVerified);
  const vehiclesCount = profileData?.vehiclesCount || 0;
  const maxVehicles = profileData?.maxVehicles || 3;
  const quotaPercent = Math.min(
    Math.round((vehiclesCount / maxVehicles) * 100),
    100,
  );

  return (
    <Box>
      {/* Page Title */}
      <Box sx={{ mb: 4 }}>
        <Typography
          variant="h4"
          sx={{ fontWeight: 800, color: "text.primary" }}
        >
          Meu Perfil & Configurações
        </Typography>
        <Typography variant="body2" sx={{ color: "text.secondary", mt: 0.5 }}>
          Gerencie seus dados de acesso, status de e-mail e limites de frota.
        </Typography>
      </Box>

      {pageError && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {pageError}
        </Alert>
      )}

      <Grid container spacing={3}>
        {/* Card 1: Status de E-mail */}
        <Grid size={{ xs: 12, md: 6 }}>
          <Card
            elevation={0}
            sx={{
              height: "100%",
              borderRadius: 3,
              border: "1px solid #E2E8F0",
              p: 1,
            }}
          >
            <CardContent>
              <Stack
                direction="row"
                spacing={1.5}
                sx={{ alignItems: "center", mb: 2 }}
              >
                <MarkEmailReadRoundedIcon
                  sx={{ color: "primary.main", fontSize: 28 }}
                />
                <Typography variant="h6" sx={{ fontWeight: 700 }}>
                  Status do E-mail
                </Typography>
              </Stack>

              <Box
                sx={{
                  p: 2,
                  borderRadius: 2,
                  bgcolor: isEmailVerified ? "success.50" : "warning.50",
                  border: "1px solid",
                  borderColor: isEmailVerified ? "success.200" : "warning.200",
                  mb: 2.5,
                }}
              >
                <Stack
                  direction="row"
                  spacing={1}
                  sx={{ alignItems: "center", mb: 1 }}
                >
                  {isEmailVerified ? (
                    <VerifiedUserRoundedIcon color="success" />
                  ) : (
                    <GppMaybeRoundedIcon color="warning" />
                  )}
                  <Chip
                    label={
                      isEmailVerified
                        ? "E-mail Verificado"
                        : "E-mail Não Verificado"
                    }
                    color={isEmailVerified ? "success" : "warning"}
                    size="small"
                    sx={{ fontWeight: 700 }}
                  />
                </Stack>
                <Typography
                  variant="body2"
                  sx={{
                    color: isEmailVerified ? "success.900" : "warning.900",
                    fontWeight: 500,
                  }}
                >
                  {isEmailVerified
                    ? "Seu endereço de e-mail está verificado e seguro. Você tem permissão total para cadastrar veículos."
                    : "Você precisa confirmar seu e-mail para cadastrar novos veículos na plataforma."}
                </Typography>
              </Box>

              {emailSuccess && (
                <Alert
                  severity="success"
                  onClose={() => setEmailSuccess(null)}
                  sx={{ mb: 2 }}
                >
                  {emailSuccess}
                </Alert>
              )}

              {emailError && (
                <Alert
                  severity="error"
                  onClose={() => setEmailError(null)}
                  sx={{ mb: 2 }}
                >
                  {emailError}
                </Alert>
              )}

              {!isEmailVerified && (
                <Button
                  variant="contained"
                  color="warning"
                  fullWidth
                  onClick={handleResendVerification}
                  disabled={resendingEmail || cooldown > 0}
                  startIcon={
                    resendingEmail ? (
                      <CircularProgress size={18} color="inherit" />
                    ) : (
                      <SendRoundedIcon />
                    )
                  }
                  sx={{ fontWeight: 600 }}
                >
                  {cooldown > 0
                    ? `Reenviar em ${cooldown}s`
                    : resendingEmail
                      ? "Enviando..."
                      : "Reenviar e-mail de confirmação"}
                </Button>
              )}
            </CardContent>
          </Card>
        </Grid>

        {/* Card 2: Cota de Veículos */}
        <Grid size={{ xs: 12, md: 6 }}>
          <Card
            elevation={0}
            sx={{
              height: "100%",
              borderRadius: 3,
              border: "1px solid #E2E8F0",
              p: 1,
            }}
          >
            <CardContent>
              <Stack
                direction="row"
                spacing={1.5}
                sx={{ alignItems: "center", mb: 2 }}
              >
                <DirectionsCarRoundedIcon
                  sx={{ color: "primary.main", fontSize: 28 }}
                />
                <Typography variant="h6" sx={{ fontWeight: 700 }}>
                  Cota de Veículos
                </Typography>
              </Stack>

              <Box sx={{ mb: 3 }}>
                <Stack
                  direction="row"
                  sx={{
                    justifyContent: "space-between",
                    alignItems: "center",
                    mb: 1,
                  }}
                >
                  <Typography variant="body2" sx={{ fontWeight: 600 }}>
                    Consumo Atual
                  </Typography>
                  <Typography
                    variant="body2"
                    sx={{ fontWeight: 700, color: "primary.main" }}
                  >
                    {vehiclesCount} de {maxVehicles} veículos
                  </Typography>
                </Stack>
                <LinearProgress
                  variant="determinate"
                  value={quotaPercent}
                  sx={{
                    height: 10,
                    borderRadius: 5,
                    bgcolor: "#E2E8F0",
                    "& .MuiLinearProgress-bar": {
                      bgcolor:
                        quotaPercent >= 100
                          ? "error.main"
                          : quotaPercent >= 66
                            ? "warning.main"
                            : "primary.main",
                      borderRadius: 5,
                    },
                  }}
                />
              </Box>

              <Typography variant="body2" sx={{ color: "text.secondary" }}>
                {vehiclesCount >= maxVehicles
                  ? "Você atingiu a cota máxima de veículos cadastrados na sua conta. Entre em contato com o suporte caso necessite expandir o limite."
                  : `Você ainda pode cadastrar mais ${maxVehicles - vehiclesCount} veículo(s) com sua cota atual.`}
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        {/* Card 3: Dados Pessoais (Nome) */}
        <Grid size={{ xs: 12, md: 6 }}>
          <Card
            elevation={0}
            sx={{
              borderRadius: 3,
              border: "1px solid #E2E8F0",
              p: 1,
            }}
          >
            <CardContent>
              <Stack
                direction="row"
                spacing={1.5}
                sx={{ alignItems: "center", mb: 2 }}
              >
                <PersonRoundedIcon
                  sx={{ color: "primary.main", fontSize: 28 }}
                />
                <Typography variant="h6" sx={{ fontWeight: 700 }}>
                  Dados Pessoais
                </Typography>
              </Stack>

              {nameSuccess && (
                <Alert
                  severity="success"
                  onClose={() => setNameSuccess(null)}
                  sx={{ mb: 2 }}
                >
                  {nameSuccess}
                </Alert>
              )}

              {nameError && (
                <Alert
                  severity="error"
                  onClose={() => setNameError(null)}
                  sx={{ mb: 2 }}
                >
                  {nameError}
                </Alert>
              )}

              <Box component="form" onSubmit={handleUpdateName} noValidate>
                <Stack spacing={2.5}>
                  <TextField
                    fullWidth
                    label="E-mail (Não editável)"
                    value={profileData?.user?.email || ""}
                    disabled
                  />

                  <TextField
                    fullWidth
                    label="Nome Completo"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    disabled={updatingName}
                    required
                  />

                  <Button
                    type="submit"
                    variant="contained"
                    startIcon={
                      updatingName ? (
                        <CircularProgress size={18} color="inherit" />
                      ) : (
                        <SaveRoundedIcon />
                      )
                    }
                    disabled={updatingName || !name.trim()}
                    sx={{ alignSelf: "flex-start", fontWeight: 600 }}
                  >
                    {updatingName ? "Salvando..." : "Salvar Alterações"}
                  </Button>
                </Stack>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        {/* Card 4: Contas Conectadas */}
        <Grid size={{ xs: 12, md: 6 }}>
          <Card
            elevation={0}
            sx={{
              height: "100%",
              borderRadius: 3,
              border: "1px solid #E2E8F0",
              p: 1,
            }}
          >
            <CardContent>
              <Stack
                direction="row"
                spacing={1.5}
                sx={{ alignItems: "center", mb: 2 }}
              >
                <LinkRoundedIcon sx={{ color: "primary.main", fontSize: 28 }} />
                <Typography variant="h6" sx={{ fontWeight: 700 }}>
                  Contas Conectadas
                </Typography>
              </Stack>

              <Typography
                variant="body2"
                sx={{ color: "text.secondary", mb: 2.5 }}
              >
                Vincule provedores sociais para acessar o CVMC com agilidade e
                segurança.
              </Typography>

              {providerSuccess && (
                <Alert
                  severity="success"
                  onClose={() => setProviderSuccess(null)}
                  sx={{ mb: 2 }}
                >
                  {providerSuccess}
                </Alert>
              )}

              {providerError && (
                <Alert
                  severity="error"
                  onClose={() => setProviderError(null)}
                  sx={{ mb: 2 }}
                >
                  {providerError}
                </Alert>
              )}

              {/* Google Provider Row */}
              <Box
                sx={{
                  p: 2,
                  borderRadius: 2,
                  border: "1px solid #E2E8F0",
                  bgcolor: "background.default",
                }}
              >
                <Stack
                  direction="row"
                  spacing={2}
                  sx={{
                    alignItems: "center",
                    justifyContent: "space-between",
                    flexWrap: "wrap",
                    gap: 1.5,
                  }}
                >
                  <Stack
                    direction="row"
                    spacing={1.5}
                    sx={{ alignItems: "center" }}
                  >
                    <GoogleIcon />
                    <Box>
                      <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>
                        Google
                      </Typography>
                      <Typography
                        variant="caption"
                        sx={{ color: "text.secondary" }}
                      >
                        {authProviders?.googleLinked
                          ? authProviders.googleEmail ||
                            user?.email ||
                            "Conectado"
                          : "Não conectado"}
                      </Typography>
                    </Box>
                  </Stack>

                  <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                    {authProviders?.googleLinked ? (
                      <>
                        <Chip
                          label="Vinculado"
                          size="small"
                          color="success"
                          variant="outlined"
                          sx={{ fontWeight: 600 }}
                        />
                        <Button
                          size="small"
                          variant="outlined"
                          color="error"
                          startIcon={
                            unlinkingGoogle ? (
                              <CircularProgress size={14} color="inherit" />
                            ) : (
                              <LinkOffRoundedIcon fontSize="small" />
                            )
                          }
                          disabled={unlinkingGoogle}
                          onClick={handleUnlinkGoogle}
                          data-testid="unlink-google-btn"
                          sx={{ textTransform: "none", fontWeight: 600 }}
                        >
                          {unlinkingGoogle ? "Desvinculando..." : "Desvincular"}
                        </Button>
                      </>
                    ) : (
                      <Box sx={{ minWidth: 160 }}>
                        <GoogleSignInButton
                          onSuccess={handleLinkGoogle}
                          text="signin_with"
                          label="Vincular Conta"
                        />
                      </Box>
                    )}
                  </Box>
                </Stack>

                {authProviders?.googleLinked && !authProviders?.hasPassword && (
                  <Alert severity="info" sx={{ mt: 2, fontSize: "0.8rem" }}>
                    Você se autenticou via Google e ainda não possui senha.
                    Cadastre uma senha de acesso ao lado antes de desvincular
                    para evitar bloqueio.
                  </Alert>
                )}
              </Box>
            </CardContent>
          </Card>
        </Grid>

        {/* Card 5: Segurança & Alteração / Definição de Senha */}
        <Grid size={{ xs: 12, md: 6 }}>
          <Card
            elevation={0}
            sx={{
              borderRadius: 3,
              border: "1px solid #E2E8F0",
              p: 1,
            }}
          >
            <CardContent>
              <Stack
                direction="row"
                spacing={1.5}
                sx={{ alignItems: "center", mb: 2 }}
              >
                <LockResetRoundedIcon
                  sx={{ color: "primary.main", fontSize: 28 }}
                />
                <Typography variant="h6" sx={{ fontWeight: 700 }}>
                  {authProviders && !authProviders.hasPassword
                    ? "Definir Senha Inicial"
                    : "Alterar Senha"}
                </Typography>
              </Stack>

              {authProviders && !authProviders.hasPassword && (
                <Typography
                  variant="body2"
                  sx={{ color: "text.secondary", mb: 2 }}
                >
                  Sua conta foi criada via Google. Defina uma senha caso queira
                  fazer login também utilizando e-mail e senha.
                </Typography>
              )}

              {passwordSuccess && (
                <Alert
                  severity="success"
                  onClose={() => setPasswordSuccess(null)}
                  sx={{ mb: 2 }}
                >
                  {passwordSuccess}
                </Alert>
              )}

              {passwordError && (
                <Alert
                  severity="error"
                  onClose={() => setPasswordError(null)}
                  sx={{ mb: 2 }}
                >
                  {passwordError}
                </Alert>
              )}

              <Box component="form" onSubmit={handleUpdatePassword} noValidate>
                <Stack spacing={2}>
                  {authProviders?.hasPassword !== false && (
                    <>
                      <TextField
                        fullWidth
                        label="Senha Atual"
                        type={showCurrentPassword ? "text" : "password"}
                        value={currentPassword}
                        onChange={(e) => setCurrentPassword(e.target.value)}
                        disabled={updatingPassword}
                        required
                        slotProps={{
                          input: {
                            endAdornment: (
                              <InputAdornment position="end">
                                <IconButton
                                  onClick={() =>
                                    setShowCurrentPassword((prev) => !prev)
                                  }
                                  edge="end"
                                  size="small"
                                >
                                  {showCurrentPassword ? (
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
                      <Divider />
                    </>
                  )}

                  <TextField
                    fullWidth
                    label="Nova Senha"
                    type={showNewPassword ? "text" : "password"}
                    placeholder="Mínimo 8 caracteres"
                    value={newPassword}
                    onChange={(e) => setNewPassword(e.target.value)}
                    disabled={updatingPassword}
                    required
                    helperText="Entre 8 e 72 caracteres"
                    slotProps={{
                      input: {
                        endAdornment: (
                          <InputAdornment position="end">
                            <IconButton
                              onClick={() =>
                                setShowNewPassword((prev) => !prev)
                              }
                              edge="end"
                              size="small"
                            >
                              {showNewPassword ? (
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
                    type={showNewPassword ? "text" : "password"}
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    disabled={updatingPassword}
                    required
                  />

                  <Button
                    type="submit"
                    variant="contained"
                    color="primary"
                    startIcon={
                      updatingPassword ? (
                        <CircularProgress size={18} color="inherit" />
                      ) : (
                        <SaveRoundedIcon />
                      )
                    }
                    disabled={
                      updatingPassword ||
                      (authProviders?.hasPassword !== false &&
                        !currentPassword) ||
                      !newPassword ||
                      !confirmPassword
                    }
                    sx={{ alignSelf: "flex-start", fontWeight: 600 }}
                  >
                    {updatingPassword
                      ? "Salvando..."
                      : authProviders && !authProviders.hasPassword
                        ? "Cadastrar Senha"
                        : "Alterar Senha"}
                  </Button>
                </Stack>
              </Box>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
}
