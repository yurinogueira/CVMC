import { Link as RouterLink } from "react-router-dom";
import {
  Box,
  Container,
  Card,
  CardContent,
  Typography,
  Stack,
  Divider,
  Button,
  Chip,
  Alert,
} from "@mui/material";
import DirectionsCarIcon from "@mui/icons-material/DirectionsCar";
import ShieldOutlinedIcon from "@mui/icons-material/ShieldOutlined";
import DeleteForeverRoundedIcon from "@mui/icons-material/DeleteForeverRounded";
import ArrowBackRoundedIcon from "@mui/icons-material/ArrowBackRounded";
import LockOutlinedIcon from "@mui/icons-material/LockOutlined";
import { useDocumentTitle } from "../../shared";
import { brandColors } from "../../../styles/theme";
import { useAuthStore } from "../../auth/state/auth.store";

export function PrivacyPolicyPage() {
  useDocumentTitle("Política de Privacidade e Exclusão de Dados");
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);

  return (
    <Box
      sx={{
        minHeight: "100vh",
        bgcolor: "background.default",
        py: { xs: 3, sm: 5, md: 6 },
      }}
    >
      <Container maxWidth="md">
        {/* Top Navigation Bar */}
        <Stack
          direction={{ xs: "column", sm: "row" }}
          spacing={2}
          sx={{
            justifyContent: "space-between",
            alignItems: { xs: "flex-start", sm: "center" },
            mb: 3,
          }}
        >
          <Stack direction="row" spacing={1.5} sx={{ alignItems: "center" }}>
            <Box
              sx={{
                bgcolor: brandColors.blue,
                borderRadius: 1.5,
                p: 0.9,
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
              }}
            >
              <DirectionsCarIcon sx={{ fontSize: 26, color: "#FFFFFF" }} />
            </Box>
            <Box>
              <Typography
                variant="h6"
                sx={{ fontWeight: 800, color: "primary.main", lineHeight: 1.1 }}
              >
                CVMC — Como Vai Meu Carro
              </Typography>
              <Typography variant="caption" sx={{ color: "text.secondary" }}>
                Central de Privacidade, Segurança e Proteção de Dados (LGPD)
              </Typography>
            </Box>
          </Stack>

          <Button
            component={RouterLink}
            to={isAuthenticated ? "/profile" : "/login"}
            variant="outlined"
            size="small"
            startIcon={<ArrowBackRoundedIcon />}
            sx={{ fontWeight: 600 }}
          >
            {isAuthenticated ? "Voltar ao Meu Perfil" : "Ir para Login"}
          </Button>
        </Stack>

        <Card
          elevation={0}
          sx={{
            borderRadius: 3,
            border: "1px solid #E2E8F0",
            bgcolor: "background.paper",
          }}
        >
          <CardContent sx={{ p: { xs: 2.5, sm: 4, md: 5 } }}>
            {/* Header */}
            <Stack spacing={1.5} sx={{ mb: 3 }}>
              <Stack
                direction="row"
                spacing={1}
                sx={{ alignItems: "center", flexWrap: "wrap", gap: 1 }}
              >
                <Chip
                  icon={<ShieldOutlinedIcon />}
                  label="Conformidade LGPD & Google Play"
                  color="primary"
                  size="small"
                  variant="outlined"
                  sx={{ fontWeight: 700 }}
                />
                <Typography variant="caption" sx={{ color: "text.secondary" }}>
                  Última atualização: 27 de outubro de 2026
                </Typography>
              </Stack>

              <Typography
                component="h1"
                variant="h4"
                sx={{ fontWeight: 800, color: "text.primary" }}
              >
                Política de Privacidade
              </Typography>

              <Typography variant="body1" sx={{ color: "text.secondary" }}>
                Esta Política de Privacidade descreve de forma transparente como
                o aplicativo móvel Android (pacote <code>com.cvmc.app</code>) e
                a plataforma web <strong>CVMC (Como Vai Meu Carro)</strong>{" "}
                coletam, utilizam, armazenam, protegem e excluem os dados dos
                usuários, em estrita conformidade com a Lei Geral de Proteção de
                Dados Pessoais (Lei nº 13.709/2018 — LGPD) e com a Política de
                Dados do Usuário da Google Play Store.
              </Typography>
            </Stack>

            <Divider sx={{ my: 3 }} />

            <Stack spacing={3.5}>
              {/* 1. Dados Coletados */}
              <Box component="section">
                <Typography
                  component="h2"
                  variant="h6"
                  sx={{ fontWeight: 700, mb: 1, color: "text.primary" }}
                >
                  1. Dados Coletados e Finalidade
                </Typography>
                <Typography
                  variant="body2"
                  sx={{ color: "text.secondary", mb: 1.5 }}
                >
                  Coletamos apenas as informações estritamente necessárias para
                  o funcionamento das funcionalidades de gestão veicular:
                </Typography>
                <Stack component="ul" spacing={1} sx={{ pl: 2.5, m: 0 }}>
                  <Typography
                    component="li"
                    variant="body2"
                    sx={{ color: "text.secondary" }}
                  >
                    <strong>Dados de Identificação e Autenticação:</strong> Nome
                    completo, endereço de e-mail, hash criptográfico de senha
                    (via <code>bcrypt</code>) ou identificador único de login
                    social (Google OAuth 2.0 / Credential Manager) para criação
                    e proteção da conta do condutor.
                  </Typography>
                  <Typography
                    component="li"
                    variant="body2"
                    sx={{ color: "text.secondary" }}
                  >
                    <strong>Dados de Veículos, Manutenções e Custos:</strong>{" "}
                    Marca, modelo, ano, quilometragem, registros de revisões
                    preventivas/corretivas, histórico de abastecimentos e
                    comprovantes/imagens enviados voluntariamente pelo usuário.
                  </Typography>
                  <Typography
                    component="li"
                    variant="body2"
                    sx={{ color: "text.secondary" }}
                  >
                    <strong>
                      Identificadores de Dispositivo para Notificações Push:
                    </strong>{" "}
                    No aplicativo Android, caso o usuário conceda a permissão{" "}
                    <code>POST_NOTIFICATIONS</code>, coletamos o token de
                    registro do Firebase Cloud Messaging (FCM) exclusivamente
                    para enviar alertas de revisões preventivas do veículo. O
                    token é desvinculado automaticamente no logout ou na
                    exclusão da conta.
                  </Typography>
                </Stack>
              </Box>

              {/* 2. Uso e Compartilhamento */}
              <Box component="section">
                <Typography
                  component="h2"
                  variant="h6"
                  sx={{ fontWeight: 700, mb: 1, color: "text.primary" }}
                >
                  2. Uso e Compartilhamento de Dados
                </Typography>
                <Typography variant="body2" sx={{ color: "text.secondary" }}>
                  O CVMC <strong>não vende, aluga ou comercializa</strong> dados
                  pessoais com terceiros ou redes de publicidade. Os dados são
                  processados unicamente para prover os serviços da plataforma,
                  utilizando provedores de infraestrutura seguros (Google
                  Identity Services / Firebase Cloud Messaging para autenticação
                  e notificações push, MongoDB com isolamento lógico por usuário
                  e envio transacional de e-mails de segurança).
                </Typography>
              </Box>

              {/* 3. Segurança e Armazenamento */}
              <Box component="section">
                <Typography
                  component="h2"
                  variant="h6"
                  sx={{ fontWeight: 700, mb: 1, color: "text.primary" }}
                >
                  3. Segurança e Criptografia em Trânsito e Repouso
                </Typography>
                <Stack
                  direction="row"
                  spacing={1.5}
                  sx={{ alignItems: "flex-start" }}
                >
                  <LockOutlinedIcon
                    sx={{ color: "primary.main", mt: 0.3, fontSize: 20 }}
                  />
                  <Typography variant="body2" sx={{ color: "text.secondary" }}>
                    Todas as comunicações entre o aplicativo Android/Web e a API
                    do CVMC ocorrem exclusivamente sobre conexões criptografadas
                    HTTPS/TLS (<code>cleartextTraffic: false</code>). Senhas são
                    armazenadas exclusivamente sob hash <code>bcrypt</code> e
                    tokens de verificação ou redefinição de senha são protegidos
                    por hash SHA-256.
                  </Typography>
                </Stack>
              </Box>

              {/* 4. Exclusão de Conta e Dados */}
              <Box
                component="section"
                id="exclusao-de-conta"
                sx={{
                  p: { xs: 2, sm: 3 },
                  borderRadius: 2.5,
                  border: "1px solid",
                  borderColor: "error.200",
                  bgcolor: "error.50",
                }}
              >
                <Stack
                  direction="row"
                  spacing={1.2}
                  sx={{ alignItems: "center", mb: 1.5 }}
                >
                  <DeleteForeverRoundedIcon color="error" />
                  <Typography
                    component="h2"
                    variant="h6"
                    sx={{ fontWeight: 800, color: "error.main" }}
                  >
                    4. Exclusão de Conta e Dados do Usuário
                  </Typography>
                </Stack>

                <Typography
                  variant="body2"
                  sx={{ color: "text.primary", mb: 1.5 }}
                >
                  Você pode excluir permanentemente sua conta e todos os dados
                  associados a qualquer momento, diretamente pelo aplicativo
                  Android ou pela plataforma Web:
                </Typography>

                <Stack component="ol" spacing={1} sx={{ pl: 2.5, mb: 2 }}>
                  <Typography component="li" variant="body2">
                    Acesse sua conta no aplicativo <strong>CVMC</strong> ou em{" "}
                    <RouterLink to="/login">
                      cvmc.yurinogueira.dev.br/login
                    </RouterLink>
                    .
                  </Typography>
                  <Typography component="li" variant="body2">
                    Navegue até o menu <strong>Meu Perfil</strong> (
                    <code>/profile</code>).
                  </Typography>
                  <Typography component="li" variant="body2">
                    Na seção <strong>Zona de Perigo — Excluir Conta</strong>,
                    clique em <strong>Excluir Minha Conta</strong> e confirme a
                    operação.
                  </Typography>
                </Stack>

                <Alert severity="warning" sx={{ mb: 2 }}>
                  <strong>O que é excluído imediatamente:</strong> Ao confirmar
                  a exclusão, são removidos de forma imediata e irreversível o
                  seu cadastro de usuário (nome, e-mail, credenciais), todos os
                  veículos de sua propriedade, todo o histórico de manutenções,
                  abastecimentos, vínculos de compartilhamento e tokens de
                  notificações push (nenhum dado pessoal é retido após a
                  exclusão).
                </Alert>

                <Stack
                  direction={{ xs: "column", sm: "row" }}
                  spacing={1.5}
                  sx={{ alignItems: { xs: "stretch", sm: "center" } }}
                >
                  <Button
                    component={RouterLink}
                    to={isAuthenticated ? "/profile" : "/login"}
                    variant="contained"
                    color="error"
                    size="small"
                    sx={{ fontWeight: 700 }}
                  >
                    {isAuthenticated
                      ? "Gerenciar / Excluir Conta no Perfil"
                      : "Entrar para Excluir Minha Conta"}
                  </Button>
                  <Typography
                    variant="caption"
                    sx={{ color: "text.secondary" }}
                  >
                    Caso não consiga acessar sua conta, você também pode
                    solicitar a exclusão enviando um e-mail para{" "}
                    <strong>suporte@yurinogueira.dev.br</strong>.
                  </Typography>
                </Stack>
              </Box>

              {/* 5. Direitos do Titular e Contato */}
              <Box component="section">
                <Typography
                  component="h2"
                  variant="h6"
                  sx={{ fontWeight: 700, mb: 1, color: "text.primary" }}
                >
                  5. Direitos do Titular (LGPD) e Contato do Controlador
                </Typography>
                <Typography variant="body2" sx={{ color: "text.secondary" }}>
                  Em conformidade com o Art. 18 da LGPD, você possui direito de
                  confirmar a existência de tratamento, acessar, corrigir dados
                  incompletos ou desatualizados, revogar consentimentos e
                  eliminar seus dados pessoais a qualquer momento através do
                  painel <strong>Meu Perfil</strong>. Para dúvidas ou
                  solicitações adicionais sobre privacidade, entre em contato
                  com o responsável pelo tratamento de dados:
                </Typography>
                <Typography
                  variant="body2"
                  sx={{ color: "text.primary", fontWeight: 600, mt: 1 }}
                >
                  Responsável: Yuri Nogueira — E-mail:
                  suporte@yurinogueira.dev.br
                </Typography>
              </Box>
            </Stack>
          </CardContent>
        </Card>
      </Container>
    </Box>
  );
}
