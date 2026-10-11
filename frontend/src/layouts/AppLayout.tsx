import { useState } from "react";
import { Box } from "@mui/material";
import { Outlet } from "react-router-dom";
import { Sidebar, DRAWER_WIDTH } from "./Sidebar";
import { Topbar } from "./Topbar";

export function AppLayout() {
  const [mobileOpen, setMobileOpen] = useState(false);

  const handleDrawerToggle = () => {
    setMobileOpen((prev) => !prev);
  };

  const handleMobileClose = () => {
    setMobileOpen(false);
  };

  return (
    <Box
      sx={{
        display: "flex",
        minHeight: "100vh",
        bgcolor: "background.default",
      }}
    >
      {/* Navigation Sidebar */}
      <Sidebar mobileOpen={mobileOpen} onMobileClose={handleMobileClose} />

      {/* Main Content Area */}
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          width: { xs: "100%", md: `calc(100% - ${DRAWER_WIDTH}px)` },
          minWidth: 0,
          maxWidth: "100%",
          minHeight: "100vh",
          display: "flex",
          flexDirection: "column",
          overflowX: "hidden",
        }}
      >
        <Topbar onDrawerToggle={handleDrawerToggle} />

        <Box
          sx={{
            flex: 1,
            p: { xs: 1.5, sm: 2.5, md: 3 },
            minWidth: 0,
            maxWidth: "100%",
          }}
        >
          <Box sx={{ maxWidth: 1600, mx: "auto", width: "100%", minWidth: 0 }}>
            <Outlet />
          </Box>
        </Box>
      </Box>
    </Box>
  );
}
