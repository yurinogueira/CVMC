import { describe, expect, it, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import { BrowserRouter } from "react-router-dom";
import { ProfilePage } from "./ProfilePage";
import { profileService } from "../services/profile.service";

vi.mock("../services/profile.service", () => ({
  profileService: {
    getProfile: vi.fn(),
    updateProfile: vi.fn(),
    updatePassword: vi.fn(),
    getAuthProviders: vi.fn().mockResolvedValue({
      hasPassword: true,
      googleLinked: false,
    }),
    linkGoogle: vi.fn(),
    unlinkGoogle: vi.fn(),
    setPassword: vi.fn(),
    deleteAccount: vi.fn(),
  },
}));

describe("ProfilePage", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders profile details, email status, and vehicle quota", async () => {
    vi.mocked(profileService.getProfile).mockResolvedValueOnce({
      user: {
        id: "user-1",
        name: "Yuri Nogueira",
        email: "yuri@cvmc.com",
        emailVerified: true,
        maxVehicles: 3,
        createdAt: "2026-08-25T00:00:00Z",
      },
      vehiclesCount: 1,
      maxVehicles: 3,
    });

    render(
      <BrowserRouter>
        <ProfilePage />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(
        screen.getByText("Meu Perfil & Configurações"),
      ).toBeInTheDocument();
    });

    expect(screen.getByText("E-mail Verificado")).toBeInTheDocument();
    expect(screen.getByText("1 de 3 veículos")).toBeInTheDocument();
    expect(screen.getByDisplayValue("Yuri Nogueira")).toBeInTheDocument();
    expect(screen.getByDisplayValue("yuri@cvmc.com")).toBeInTheDocument();
    expect(
      screen.getByRole("heading", { name: /Alterar Senha/i }),
    ).toBeInTheDocument();
  });

  it("renders unverified email badge and resend button when email is not verified", async () => {
    vi.mocked(profileService.getProfile).mockResolvedValueOnce({
      user: {
        id: "user-2",
        name: "Novo Usuário",
        email: "novo@cvmc.com",
        emailVerified: false,
        maxVehicles: 3,
        createdAt: "2026-08-25T00:00:00Z",
      },
      vehiclesCount: 0,
      maxVehicles: 3,
    });

    render(
      <BrowserRouter>
        <ProfilePage />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(screen.getByText("E-mail Não Verificado")).toBeInTheDocument();
    });

    expect(
      screen.getByRole("button", {
        name: /Reenviar e-mail de confirmação/i,
      }),
    ).toBeInTheDocument();
  });

  it("renders Google linked status and permits unlinking when password is set", async () => {
    vi.mocked(profileService.getProfile).mockResolvedValueOnce({
      user: {
        id: "user-3",
        name: "Google Linked User",
        email: "googleuser@cvmc.com",
        googleId: "google-123",
        emailVerified: true,
        maxVehicles: 3,
        createdAt: "2026-08-25T00:00:00Z",
      },
      vehiclesCount: 0,
      maxVehicles: 3,
    });
    vi.mocked(profileService.getAuthProviders).mockResolvedValueOnce({
      hasPassword: true,
      googleLinked: true,
      googleEmail: "googleuser@cvmc.com",
    });

    render(
      <BrowserRouter>
        <ProfilePage />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(screen.getByText("Contas Conectadas")).toBeInTheDocument();
    });

    expect(screen.getByText("Vinculado")).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: /Desvincular/i }),
    ).toBeInTheDocument();
  });

  it("renders 'Definir Senha Inicial' when user has only social login", async () => {
    vi.mocked(profileService.getProfile).mockResolvedValueOnce({
      user: {
        id: "user-4",
        name: "Social Only User",
        email: "social@cvmc.com",
        googleId: "google-456",
        emailVerified: true,
        maxVehicles: 3,
        createdAt: "2026-08-25T00:00:00Z",
      },
      vehiclesCount: 0,
      maxVehicles: 3,
    });
    vi.mocked(profileService.getAuthProviders).mockResolvedValueOnce({
      hasPassword: false,
      googleLinked: true,
      googleEmail: "social@cvmc.com",
    });

    render(
      <BrowserRouter>
        <ProfilePage />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(
        screen.getByRole("heading", { name: /Definir Senha Inicial/i }),
      ).toBeInTheDocument();
    });

    expect(
      screen.getByRole("button", { name: /Cadastrar Senha/i }),
    ).toBeInTheDocument();
    // Senha Atual should NOT be in the document
    expect(screen.queryByLabelText(/Senha Atual/i)).not.toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: /Excluir Minha Conta/i }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole("link", { name: /Política de Privacidade/i }),
    ).toBeInTheDocument();
  });
});
