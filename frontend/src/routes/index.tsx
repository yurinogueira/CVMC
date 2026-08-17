import { Route, Routes, Navigate } from 'react-router-dom';
import { AppLayout } from '../layouts/AppLayout';

function DashboardPage() {
  return <div>CVMC dashboard</div>;
}

export function AppRoutes() {
  return (
    <Routes>
      <Route element={<AppLayout />}>
        <Route path="/" element={<Navigate to="/dashboard" replace />} />
        <Route path="/dashboard" element={<DashboardPage />} />
      </Route>
    </Routes>
  );
}