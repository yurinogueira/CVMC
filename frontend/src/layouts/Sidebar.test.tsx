import { describe, expect, it, vi, beforeEach } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { Sidebar } from "./Sidebar";
import { useAuthStore } from "../features/auth/state/auth.store";

const mockNavigate = vi.fn();
vi.mock("react-router-dom", async () => {
  const actual =
    await vi.importActual<typeof import("react-router-dom")>(
      "react-router-dom",
    );
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

describe("Sidebar", () => {
  const mockOnMobileClose = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    useAuthStore.getState().setUser({
      id: "user-1",
      name: "João Silva",
      email: "joao@cvmc.com",
      emailVerified: true,
      maxVehicles: 3,
      createdAt: "2026-08-25T00:00:00Z",
    });
  });

  it("renders user information and navigates to /profile on click", () => {
    render(
      <MemoryRouter initialEntries={["/dashboard"]}>
        <Sidebar mobileOpen={false} onMobileClose={mockOnMobileClose} />
      </MemoryRouter>,
    );

    expect(screen.getAllByText("João Silva")[0]).toBeInTheDocument();
    expect(screen.getAllByText("joao@cvmc.com")[0]).toBeInTheDocument();

    const userShortcut = screen.getAllByTestId("sidebar-user-shortcut")[0];
    fireEvent.click(userShortcut);

    expect(mockNavigate).toHaveBeenCalledWith("/profile");
    expect(mockOnMobileClose).toHaveBeenCalled();
  });

  it("supports keyboard navigation on user shortcut with Enter and Space keys", () => {
    render(
      <MemoryRouter initialEntries={["/dashboard"]}>
        <Sidebar mobileOpen={true} onMobileClose={mockOnMobileClose} />
      </MemoryRouter>,
    );

    const userShortcut = screen.getAllByTestId("sidebar-user-shortcut")[0];
    fireEvent.keyDown(userShortcut, { key: "Enter" });

    expect(mockNavigate).toHaveBeenCalledWith("/profile");
    expect(mockOnMobileClose).toHaveBeenCalled();

    mockNavigate.mockClear();
    mockOnMobileClose.mockClear();

    fireEvent.keyDown(userShortcut, { key: " " });
    expect(mockNavigate).toHaveBeenCalledWith("/profile");
    expect(mockOnMobileClose).toHaveBeenCalled();
  });

  it("renders brand logo and nav items when mobile drawer is open", () => {
    render(
      <MemoryRouter initialEntries={["/dashboard"]}>
        <Sidebar mobileOpen={true} onMobileClose={mockOnMobileClose} />
      </MemoryRouter>,
    );

    expect(screen.getAllByText("Como Vai Meu Carro")[0]).toBeInTheDocument();
    expect(screen.getAllByText("Dashboard")[0]).toBeInTheDocument();
  });

  it("renders user avatar with avatarUrl and referrerPolicy='no-referrer'", () => {
    useAuthStore.getState().setUser({
      id: "user-1",
      name: "João Silva",
      email: "joao@cvmc.com",
      avatarUrl: "https://lh3.googleusercontent.com/joao.jpg",
      emailVerified: true,
      maxVehicles: 3,
      createdAt: "2026-08-25T00:00:00Z",
    });

    render(
      <MemoryRouter initialEntries={["/dashboard"]}>
        <Sidebar mobileOpen={false} onMobileClose={mockOnMobileClose} />
      </MemoryRouter>,
    );

    const avatarImgs = screen.getAllByRole("img", { name: "João Silva" });
    expect(avatarImgs.length).toBeGreaterThan(0);
    expect(avatarImgs[0]).toHaveAttribute(
      "src",
      "https://lh3.googleusercontent.com/joao.jpg",
    );
    expect(avatarImgs[0]).toHaveAttribute("referrerpolicy", "no-referrer");
  });
});
