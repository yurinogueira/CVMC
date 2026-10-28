import { describe, it, expect, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter, Routes, Route } from "react-router-dom";
import { ProtectedRoute } from "./ProtectedRoute";
import { useAuthStore } from "../features/auth/state/auth.store";

describe("ProtectedRoute", () => {
  beforeEach(() => {
    useAuthStore.getState().clear();
  });

  it("renders loading fallback while isInitializing is true", () => {
    useAuthStore.setState({ isInitializing: true, isAuthenticated: false });

    render(
      <MemoryRouter initialEntries={["/protected"]}>
        <Routes>
          <Route element={<ProtectedRoute />}>
            <Route path="/protected" element={<div>Conteúdo Protegido</div>} />
          </Route>
        </Routes>
      </MemoryRouter>,
    );

    expect(screen.getByRole("progressbar")).toBeInTheDocument();
    expect(screen.queryByText("Conteúdo Protegido")).not.toBeInTheDocument();
  });

  it("redirects to /login when isInitializing is false and user is not authenticated", () => {
    useAuthStore.setState({ isInitializing: false, isAuthenticated: false });

    render(
      <MemoryRouter initialEntries={["/protected"]}>
        <Routes>
          <Route path="/login" element={<div>Tela de Login</div>} />
          <Route element={<ProtectedRoute />}>
            <Route path="/protected" element={<div>Conteúdo Protegido</div>} />
          </Route>
        </Routes>
      </MemoryRouter>,
    );

    expect(screen.getByText("Tela de Login")).toBeInTheDocument();
    expect(screen.queryByText("Conteúdo Protegido")).not.toBeInTheDocument();
  });

  it("renders child outlet when isInitializing is false and user is authenticated", () => {
    useAuthStore.setState({
      isInitializing: false,
      isAuthenticated: true,
      user: { id: "1", name: "Tester", email: "test@cvmc.com" },
    });

    render(
      <MemoryRouter initialEntries={["/protected"]}>
        <Routes>
          <Route path="/login" element={<div>Tela de Login</div>} />
          <Route element={<ProtectedRoute />}>
            <Route path="/protected" element={<div>Conteúdo Protegido</div>} />
          </Route>
        </Routes>
      </MemoryRouter>,
    );

    expect(screen.getByText("Conteúdo Protegido")).toBeInTheDocument();
    expect(screen.queryByText("Tela de Login")).not.toBeInTheDocument();
  });
});
