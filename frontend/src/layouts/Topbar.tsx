import React, { useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import {
  AppBar,
  Toolbar,
  Typography,
  IconButton,
  Box,
  Avatar,
  Menu,
  MenuItem,
  ListItemIcon,
  ListItemText,
  Divider,
} from "@mui/material";
import MenuRoundedIcon from "@mui/icons-material/MenuRounded";
import LogoutRoundedIcon from "@mui/icons-material/LogoutRounded";
import PersonRoundedIcon from "@mui/icons-material/PersonRounded";
import DirectionsCarRoundedIcon from "@mui/icons-material/DirectionsCarRounded";
import { useAuthStore } from "../features/auth/state/auth.store";

interface TopbarProps {
  onDrawerToggle: () => void;
}

const routeTitles: Record<string, string> = {
  "/dashboard": "Visão Geral",
  "/vehicles": "Meus Veículos",
  "/maintenance": "Histórico de Manutenções",
};

export function Topbar({ onDrawerToggle }: TopbarProps) {
  const location = useLocation();
  const navigate = useNavigate();
  const { user, clear } = useAuthStore();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const currentTitle = routeTitles[location.pathname] || "Painel CVMC";

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = () => {
    handleMenuClose();
    clear();
    navigate("/login", { replace: true });
  };

  return (
    <AppBar
      position="sticky"
      elevation={0}
      sx={{
        width: "100%",
        bgcolor: "background.paper",
        color: "text.primary",
        borderBottom: "1px solid #E2E8F0",
      }}
    >
      <Toolbar
        sx={{
          justifyContent: "space-between",
          height: 64,
          minHeight: "64px !important",
          px: { xs: 2, sm: 3 },
        }}
      >
        <Box sx={{ display: "flex", alignItems: "center", gap: 1.5 }}>
          <IconButton
            color="inherit"
            aria-label="abrir menu lateral"
            edge="start"
            onClick={onDrawerToggle}
            sx={{ display: { md: "none" } }}
          >
            <MenuRoundedIcon />
          </IconButton>
          <Typography
            variant="h6"
            sx={{ fontWeight: 700, fontSize: { xs: "1.1rem", sm: "1.25rem" } }}
          >
            {currentTitle}
          </Typography>
        </Box>

        <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
          <IconButton onClick={handleMenuOpen} size="small" sx={{ p: 0.5 }}>
            <Avatar
              sx={{
                width: 38,
                height: 38,
                bgcolor: "primary.main",
                fontSize: "0.9rem",
                fontWeight: 600,
                border: "2px solid #E2E8F0",
              }}
            >
              {user?.name ? user.name.charAt(0).toUpperCase() : "U"}
            </Avatar>
          </IconButton>

          <Menu
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleMenuClose}
            transformOrigin={{ horizontal: "right", vertical: "top" }}
            anchorOrigin={{ horizontal: "right", vertical: "bottom" }}
            slotProps={{
              paper: {
                elevation: 0,
                sx: {
                  overflow: "visible",
                  filter: "drop-shadow(0px 8px 16px rgba(15, 23, 42, 0.08))",
                  mt: 1.5,
                  minWidth: 200,
                  borderRadius: 2,
                  border: "1px solid #E2E8F0",
                },
              },
            }}
          >
            <Box sx={{ px: 2, py: 1.5 }}>
              <Typography variant="subtitle2" sx={{ fontWeight: 600 }}>
                {user?.name || "Usuário"}
              </Typography>
              <Typography
                variant="caption"
                sx={{ color: "text.secondary", display: "block" }}
              >
                {user?.email || "usuario@cvmc.com"}
              </Typography>
            </Box>

            <Divider />

            <MenuItem
              onClick={() => {
                handleMenuClose();
                navigate("/vehicles");
              }}
              sx={{ py: 1 }}
            >
              <ListItemIcon>
                <DirectionsCarRoundedIcon fontSize="small" />
              </ListItemIcon>
              <ListItemText primary="Meus Veículos" />
            </MenuItem>

            <MenuItem
              onClick={() => {
                handleMenuClose();
                navigate("/dashboard");
              }}
              sx={{ py: 1 }}
            >
              <ListItemIcon>
                <PersonRoundedIcon fontSize="small" />
              </ListItemIcon>
              <ListItemText primary="Meu Painel" />
            </MenuItem>

            <Divider />

            <MenuItem
              onClick={handleLogout}
              sx={{ py: 1, color: "error.main" }}
            >
              <ListItemIcon sx={{ color: "error.main" }}>
                <LogoutRoundedIcon fontSize="small" />
              </ListItemIcon>
              <ListItemText primary="Sair do CVMC" />
            </MenuItem>
          </Menu>
        </Box>
      </Toolbar>
    </AppBar>
  );
}
