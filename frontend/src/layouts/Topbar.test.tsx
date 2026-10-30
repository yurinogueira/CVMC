import { describe, expect, it, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { Topbar, getTopbarTitle } from "./Topbar";

vi.mock("../features/auth/services/auth.service", () => ({
  authService: {
    logout: vi.fn().mockResolvedValue(undefined),
  },
}));

vi.mock("../services/notifications/push.service", () => ({
  pushNotificationService: {
    unregisterTokenFromBackend: vi.fn().mockResolvedValue(undefined),
    cleanup: vi.fn().mockResolvedValue(undefined),
  },
}));

let mockUser: { name?: string; email?: string; avatarUrl?: string } | null = {
  name: "Yuri Nogueira",
  email: "yuri@cvmc.dev.br",
};

vi.mock("../features/auth/state/auth.store", () => ({
  useAuthStore: () => ({
    user: mockUser,
    clear: vi.fn(),
  }),
}));

describe("getTopbarTitle", () => {
  it("resolves static routes correctly", () => {
    expect(getTopbarTitle("/dashboard")).toBe("Dashboard");
    expect(getTopbarTitle("/vehicles")).toBe("Meus Veículos");
    expect(getTopbarTitle("/maintenance")).toBe("Manutenções & Revisões");
    expect(getTopbarTitle("/reports")).toBe("Relatórios & Exportações");
    expect(getTopbarTitle("/profile")).toBe("Meu Perfil");
    expect(getTopbarTitle("/privacy")).toBe("Política de Privacidade");
  });

  it("normalizes trailing slashes properly", () => {
    expect(getTopbarTitle("/dashboard/")).toBe("Dashboard");
    expect(getTopbarTitle("/vehicles/")).toBe("Meus Veículos");
    expect(getTopbarTitle("/reports/")).toBe("Relatórios & Exportações");
    expect(getTopbarTitle("/profile///")).toBe("Meu Perfil");
  });

  it("resolves dynamic routes properly", () => {
    expect(getTopbarTitle("/vehicles/car-123")).toBe("Detalhes do Veículo");
    expect(getTopbarTitle("/vehicles/car-123/maintenance/new")).toBe(
      "Registrar Manutenção",
    );
    expect(getTopbarTitle("/vehicles/car-123/maintenance/m-456/edit")).toBe(
      "Editar Manutenção",
    );
  });

  it("returns fallback for unknown paths", () => {
    expect(getTopbarTitle("/unknown-route")).toBe("Painel CVMC");
    expect(getTopbarTitle("")).toBe("Painel CVMC");
  });
});

describe("Topbar Component", () => {
  it("renders normalized title and user avatar", () => {
    render(
      <MemoryRouter initialEntries={["/reports"]}>
        <Topbar onDrawerToggle={vi.fn()} />
      </MemoryRouter>,
    );

    expect(screen.getByText("Relatórios & Exportações")).toBeInTheDocument();
    expect(screen.getByText("Y")).toBeInTheDocument();
  });

  it("renders dynamic vehicle details title", () => {
    render(
      <MemoryRouter initialEntries={["/vehicles/car-uuid-99"]}>
        <Topbar onDrawerToggle={vi.fn()} />
      </MemoryRouter>,
    );

    expect(screen.getByText("Detalhes do Veículo")).toBeInTheDocument();
  });

  it("renders avatar image with referrerPolicy='no-referrer' when avatarUrl is present", () => {
    mockUser = {
      name: "Yuri Nogueira",
      email: "yuri@cvmc.dev.br",
      avatarUrl: "https://lh3.googleusercontent.com/avatar.jpg",
    };

    render(
      <MemoryRouter initialEntries={["/dashboard"]}>
        <Topbar onDrawerToggle={vi.fn()} />
      </MemoryRouter>,
    );

    const avatarImg = screen.getByRole("img", { name: "Yuri Nogueira" });
    expect(avatarImg).toHaveAttribute(
      "src",
      "https://lh3.googleusercontent.com/avatar.jpg",
    );
    expect(avatarImg).toHaveAttribute("referrerpolicy", "no-referrer");

    // Reset mockUser
    mockUser = {
      name: "Yuri Nogueira",
      email: "yuri@cvmc.dev.br",
    };
  });
});
