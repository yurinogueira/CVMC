import { describe, expect, it, vi, beforeEach } from "vitest";
import { render, screen, waitFor, fireEvent } from "@testing-library/react";
import { BrowserRouter } from "react-router-dom";
import { VehicleDetailsPage } from "./VehicleDetailsPage";
import { carService } from "../services/car.service";
import { maintenanceService } from "../../maintenance/services/maintenance.service";

const mockNavigate = vi.fn();

vi.mock("react-router-dom", async () => {
  const actual =
    await vi.importActual<typeof import("react-router-dom")>(
      "react-router-dom",
    );
  return {
    ...actual,
    useParams: () => ({ id: "car-123" }),
    useNavigate: () => mockNavigate,
  };
});

vi.mock("../services/car.service", () => ({
  carService: {
    get: vi.fn(),
  },
}));

vi.mock("../../maintenance/services/maintenance.service", () => ({
  maintenanceService: {
    listByCar: vi.fn(),
    create: vi.fn(),
  },
}));

const mockCar = {
  id: "car-123",
  name: "Civic de Passeio",
  manufacturer: "Honda",
  model: "Civic Touring",
  yearManufacture: 2021,
  yearModel: 2022,
  lastMileage: 45000,
  vehicleType: "cars",
  fipeCode: "014095-3",
  fipePrice: "R$ 138.500,00",
  fuel: "Gasolina",
  ownerId: "user-1",
};

const mockMaintenances = [
  {
    id: "maint-1",
    carId: "car-123",
    title: "Troca de Óleo e Filtro",
    description: "Óleo sintético 0W20 e filtro novo",
    date: "2026-08-15T12:00:00Z",
    mileage: 40000,
  },
  {
    id: "maint-2",
    carId: "car-123",
    title: "Troca de Pastilhas de Freio",
    description: "Pastilhas dianteiras cerâmica",
    date: "2026-08-20T12:00:00Z",
    mileage: 42000,
  },
];

describe("VehicleDetailsPage", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders vehicle details, specs and FIPE info correctly", async () => {
    vi.mocked(carService.get).mockResolvedValue(mockCar);
    vi.mocked(maintenanceService.listByCar).mockResolvedValue([]);

    render(
      <BrowserRouter>
        <VehicleDetailsPage />
      </BrowserRouter>,
    );

    await waitFor(() => {
      // Vehicle name appears in title and card
      const names = screen.getAllByText("Civic de Passeio");
      expect(names.length).toBeGreaterThanOrEqual(1);
      expect(screen.getByText("Honda • Civic Touring")).toBeInTheDocument();
      expect(screen.getByText("Ano: 2021/2022")).toBeInTheDocument();
      expect(screen.getByText("45.000 km")).toBeInTheDocument();
      expect(screen.getByText("014095-3")).toBeInTheDocument();
      expect(screen.getByText("R$ 138.500,00")).toBeInTheDocument();
      expect(screen.getByText("Gasolina")).toBeInTheDocument();
    });
  });

  it("renders empty state when there are no maintenances", async () => {
    vi.mocked(carService.get).mockResolvedValue(mockCar);
    vi.mocked(maintenanceService.listByCar).mockResolvedValue([]);

    render(
      <BrowserRouter>
        <VehicleDetailsPage />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(
        screen.getByText("Nenhuma manutenção registrada para este veículo"),
      ).toBeInTheDocument();
      expect(
        screen.getByRole("button", { name: /Registrar Primeira Manutenção/i }),
      ).toBeInTheDocument();
    });
  });

  it("renders maintenance history list when maintenances exist", async () => {
    vi.mocked(carService.get).mockResolvedValue(mockCar);
    vi.mocked(maintenanceService.listByCar).mockResolvedValue(mockMaintenances);

    render(
      <BrowserRouter>
        <VehicleDetailsPage />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(screen.getByText("Troca de Óleo e Filtro")).toBeInTheDocument();
      expect(
        screen.getByText("Troca de Pastilhas de Freio"),
      ).toBeInTheDocument();
      expect(
        screen.getByText("Óleo sintético 0W20 e filtro novo"),
      ).toBeInTheDocument();
      expect(screen.getByText("40.000 km")).toBeInTheDocument();
      expect(screen.getByText("42.000 km")).toBeInTheDocument();
    });
  });

  it("renders error state when car is not found", async () => {
    vi.mocked(carService.get).mockRejectedValue(new Error("Not found"));
    vi.mocked(maintenanceService.listByCar).mockResolvedValue([]);

    render(
      <BrowserRouter>
        <VehicleDetailsPage />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(
        screen.getByText(
          "Veículo não encontrado ou você não possui permissão para visualizá-lo.",
        ),
      ).toBeInTheDocument();
    });
  });

  it("navigates to /vehicles when 'Voltar' button is clicked", async () => {
    vi.mocked(carService.get).mockResolvedValue(mockCar);
    vi.mocked(maintenanceService.listByCar).mockResolvedValue([]);

    render(
      <BrowserRouter>
        <VehicleDetailsPage />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(
        screen.getByRole("button", { name: /Voltar/i }),
      ).toBeInTheDocument();
    });

    fireEvent.click(screen.getByRole("button", { name: /Voltar/i }));
    expect(mockNavigate).toHaveBeenCalledWith("/vehicles");
  });

  it("navigates to maintenance registration page when 'Registrar Manutenção' button is clicked", async () => {
    vi.mocked(carService.get).mockResolvedValue(mockCar);
    vi.mocked(maintenanceService.listByCar).mockResolvedValue([]);

    render(
      <BrowserRouter>
        <VehicleDetailsPage />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(
        screen.getByRole("button", { name: /^Registrar Manutenção$/i }),
      ).toBeInTheDocument();
    });

    fireEvent.click(
      screen.getByRole("button", { name: /^Registrar Manutenção$/i }),
    );

    expect(mockNavigate).toHaveBeenCalledWith(
      "/vehicles/car-123/maintenance/new",
    );
  });
});
