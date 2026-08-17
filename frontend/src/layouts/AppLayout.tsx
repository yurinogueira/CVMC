import { Container, Box, Typography } from '@mui/material';
import { Outlet } from 'react-router-dom';

export function AppLayout() {
  return (
    <Box sx={{ minHeight: '100vh', bgcolor: 'background.default', py: 4 }}>
      <Container maxWidth="lg">
        <Typography variant="h4" component="h1" sx={{ mb: 3, fontWeight: 700 }}>
          CVMC
        </Typography>
        <Outlet />
      </Container>
    </Box>
  );
}