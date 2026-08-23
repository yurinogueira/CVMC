import { useLocation, useNavigate } from "react-router-dom";
import {
  Box,
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Typography,
  Stack,
  Divider,
  Avatar,
  IconButton,
  Tooltip,
} from "@mui/material";
import DashboardRoundedIcon from "@mui/icons-material/DashboardRounded";
import DirectionsCarRoundedIcon from "@mui/icons-material/DirectionsCarRounded";
import BuildRoundedIcon from "@mui/icons-material/BuildRounded";
import LogoutRoundedIcon from "@mui/icons-material/LogoutRounded";
import { brandColors } from "../styles/theme";
import { useAuthStore } from "../features/auth/state/auth.store";

export const DRAWER_WIDTH = 260;

interface SidebarProps {
  mobileOpen: boolean;
  onMobileClose: () => void;
}

const navItems = [
  {
    label: "Dashboard",
    path: "/dashboard",
    icon: <DashboardRoundedIcon fontSize="small" />,
  },
  {
    label: "Meus Veículos",
    path: "/vehicles",
    icon: <DirectionsCarRoundedIcon fontSize="small" />,
  },
  {
    label: "Manutenções",
    path: "/maintenance",
    icon: <BuildRoundedIcon fontSize="small" />,
  },
];

export function Sidebar({ mobileOpen, onMobileClose }: SidebarProps) {
  const location = useLocation();
  const navigate = useNavigate();
  const { user, clear } = useAuthStore();

  const handleLogout = () => {
    clear();
    navigate("/login", { replace: true });
  };

  const drawerContent = (
    <Box
      sx={{
        height: "100%",
        display: "flex",
        flexDirection: "column",
        bgcolor: "background.paper",
        borderRight: "1px solid #E2E8F0",
      }}
    >
      {/* Brand Header */}
      <Box
        sx={{
          height: 64,
          minHeight: 64,
          px: 2.5,
          display: "flex",
          alignItems: "center",
          gap: 1.5,
          boxSizing: "border-box",
        }}
      >
        <Box
          sx={{
            width: 40,
            height: 40,
            borderRadius: 2,
            background: brandColors.gradient,
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            boxShadow: "0 4px 12px rgba(2, 132, 199, 0.25)",
          }}
        >
          <DirectionsCarRoundedIcon sx={{ color: "#FFFFFF", fontSize: 24 }} />
        </Box>
        <Box>
          <Typography
            variant="h6"
            sx={{ fontWeight: 800, lineHeight: 1.2, color: "text.primary" }}
          >
            CVMC
          </Typography>
          <Typography
            variant="caption"
            sx={{ color: "text.secondary", fontWeight: 500 }}
          >
            Gestão de Veículos
          </Typography>
        </Box>
      </Box>

      <Divider sx={{ borderColor: "#E2E8F0" }} />

      {/* Navigation Links */}
      <Box sx={{ flex: 1, py: 2, px: 1.5 }}>
        <Typography
          variant="caption"
          sx={{
            px: 1.5,
            mb: 1,
            display: "block",
            color: "text.secondary",
            fontWeight: 600,
            textTransform: "uppercase",
            letterSpacing: 0.5,
            fontSize: "0.7rem",
          }}
        >
          Menu Principal
        </Typography>

        <List disablePadding>
          {navItems.map((item) => {
            const isActive =
              location.pathname === item.path ||
              (item.path !== "/dashboard" &&
                location.pathname.startsWith(item.path));

            return (
              <ListItem key={item.path} disablePadding sx={{ mb: 0.5 }}>
                <ListItemButton
                  onClick={() => {
                    navigate(item.path);
                    onMobileClose();
                  }}
                  sx={{
                    borderRadius: 1.5,
                    py: 1.1,
                    px: 1.5,
                    bgcolor: isActive
                      ? "rgba(76, 146, 252, 0.1)"
                      : "transparent",
                    color: isActive ? "primary.main" : "text.secondary",
                    "&:hover": {
                      bgcolor: isActive
                        ? "rgba(76, 146, 252, 0.15)"
                        : "rgba(241, 245, 249, 0.8)",
                      color: isActive ? "primary.main" : "text.primary",
                    },
                  }}
                >
                  <ListItemIcon
                    sx={{
                      minWidth: 36,
                      color: isActive ? "primary.main" : "text.secondary",
                    }}
                  >
                    {item.icon}
                  </ListItemIcon>
                  <ListItemText
                    primary={item.label}
                    slotProps={{
                      primary: {
                        sx: {
                          fontSize: "0.9rem",
                          fontWeight: isActive ? 600 : 500,
                        },
                      },
                    }}
                  />
                  {isActive && (
                    <Box
                      sx={{
                        width: 4,
                        height: 20,
                        borderRadius: 2,
                        bgcolor: "primary.main",
                        ml: 1,
                      }}
                    />
                  )}
                </ListItemButton>
              </ListItem>
            );
          })}
        </List>
      </Box>

      {/* User Section / Bottom Bar */}
      <Divider sx={{ borderColor: "#F1F5F9" }} />
      <Box sx={{ p: 2, bgcolor: "#F8FAFC" }}>
        <Stack
          direction="row"
          sx={{
            alignItems: "center",
            justifyContent: "space-between",
          }}
        >
          <Stack
            direction="row"
            spacing={1.5}
            sx={{ alignItems: "center", minWidth: 0, flex: 1 }}
          >
            <Avatar
              sx={{
                width: 36,
                height: 36,
                bgcolor: "primary.main",
                fontSize: "0.875rem",
                fontWeight: 600,
              }}
            >
              {user?.name ? user.name.charAt(0).toUpperCase() : "U"}
            </Avatar>
            <Box sx={{ minWidth: 0, flex: 1 }}>
              <Typography
                variant="body2"
                noWrap
                sx={{ fontWeight: 600, color: "text.primary" }}
              >
                {user?.name || "Usuário"}
              </Typography>
              <Typography
                variant="caption"
                noWrap
                sx={{ color: "text.secondary", display: "block" }}
              >
                {user?.email || "motorista@cvmc.com"}
              </Typography>
            </Box>
          </Stack>

          <Tooltip title="Encerrar sessão">
            <IconButton
              onClick={handleLogout}
              size="small"
              sx={{ color: "text.secondary", "&:hover": { color: "#EF4444" } }}
            >
              <LogoutRoundedIcon fontSize="small" />
            </IconButton>
          </Tooltip>
        </Stack>
      </Box>
    </Box>
  );

  return (
    <>
      {/* Mobile Drawer */}
      <Drawer
        variant="temporary"
        open={mobileOpen}
        onClose={onMobileClose}
        ModalProps={{ keepMounted: true }}
        sx={{
          display: { xs: "block", md: "none" },
          "& .MuiDrawer-paper": {
            boxSizing: "border-box",
            width: DRAWER_WIDTH,
          },
        }}
      >
        {drawerContent}
      </Drawer>

      {/* Desktop Permanent Drawer */}
      <Drawer
        variant="permanent"
        sx={{
          display: { xs: "none", md: "block" },
          width: DRAWER_WIDTH,
          flexShrink: 0,
          "& .MuiDrawer-paper": {
            boxSizing: "border-box",
            width: DRAWER_WIDTH,
            border: "none",
          },
        }}
        open
      >
        {drawerContent}
      </Drawer>
    </>
  );
}
