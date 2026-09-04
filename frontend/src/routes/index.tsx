import { lazy, Suspense } from "react";
import { Route, Routes, Navigate } from "react-router-dom";
import { AppLayout } from "../layouts/AppLayout";
import { ProtectedRoute } from "./ProtectedRoute";
import { PageLoadingFallback } from "../features/shared";

// Lazy-loaded route components for optimal bundle splitting and initial load performance
const LoginPage = lazy(() =>
  import("../features/auth/pages/LoginPage").then((m) => ({
    default: m.LoginPage,
  })),
);
const RegisterPage = lazy(() =>
  import("../features/auth/pages/RegisterPage").then((m) => ({
    default: m.RegisterPage,
  })),
);
const ForgotPasswordPage = lazy(() =>
  import("../features/auth/pages/ForgotPasswordPage").then((m) => ({
    default: m.ForgotPasswordPage,
  })),
);
const ResetPasswordPage = lazy(() =>
  import("../features/auth/pages/ResetPasswordPage").then((m) => ({
    default: m.ResetPasswordPage,
  })),
);
const VerifyEmailPage = lazy(() =>
  import("../features/auth/pages/VerifyEmailPage").then((m) => ({
    default: m.VerifyEmailPage,
  })),
);
const DashboardPage = lazy(() =>
  import("../features/dashboard/pages/DashboardPage").then((m) => ({
    default: m.DashboardPage,
  })),
);
const VehiclesPage = lazy(() =>
  import("../features/cars/pages/VehiclesPage").then((m) => ({
    default: m.VehiclesPage,
  })),
);
const VehicleDetailsPage = lazy(() =>
  import("../features/cars/pages/VehicleDetailsPage").then((m) => ({
    default: m.VehicleDetailsPage,
  })),
);
const MaintenancePage = lazy(() =>
  import("../features/maintenance/pages/MaintenancePage").then((m) => ({
    default: m.MaintenancePage,
  })),
);
const RegisterMaintenancePage = lazy(() =>
  import("../features/maintenance/pages/RegisterMaintenancePage").then((m) => ({
    default: m.RegisterMaintenancePage,
  })),
);
const ProfilePage = lazy(() =>
  import("../features/profile/pages/ProfilePage").then((m) => ({
    default: m.ProfilePage,
  })),
);

export function AppRoutes() {
  return (
    <Suspense fallback={<PageLoadingFallback />}>
      <Routes>
        {/* Public Authentication Routes */}
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/forgot-password" element={<ForgotPasswordPage />} />
        <Route path="/reset-password" element={<ResetPasswordPage />} />
        <Route path="/verify-email" element={<VerifyEmailPage />} />

        {/* Protected App Routes */}
        <Route element={<ProtectedRoute />}>
          <Route element={<AppLayout />}>
            <Route path="/" element={<Navigate to="/dashboard" replace />} />
            <Route path="/dashboard" element={<DashboardPage />} />
            <Route path="/vehicles" element={<VehiclesPage />} />
            <Route path="/vehicles/:id" element={<VehicleDetailsPage />} />
            <Route
              path="/vehicles/:id/maintenance/new"
              element={<RegisterMaintenancePage />}
            />
            <Route path="/maintenance" element={<MaintenancePage />} />
            <Route path="/profile" element={<ProfilePage />} />
          </Route>
        </Route>

        {/* Fallback */}
        <Route path="*" element={<Navigate to="/dashboard" replace />} />
      </Routes>
    </Suspense>
  );
}
