import { describe, expect, it, vi, beforeEach } from "vitest";
import { render, screen, waitFor, fireEvent } from "@testing-library/react";
import { BrowserRouter } from "react-router-dom";
import { VehicleDetailsPage } from "./VehicleDetailsPage";
import { carService } from "../services/car.service";
import { maintenanceService } from "../../maintenance/services/maintenance.service";
import { fuelService } from "../../fuel/services/fuel.service";

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
    delete: vi.fn(),
  },
}));

vi.mock("../../fuel/services/fuel.service", () => ({
  fuelService: {
    listByCar: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    delete: vi.fn(),
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
    cost: 350.0,
  },
  {
    id: "maint-2",
    carId: "car-123",
    title: "Troca de Pastilhas de Freio",
    description: "Pastilhas dianteiras cerâmica",
    date: "2026-08-20T12:00:00Z",
    mileage: 42000,
    cost: 450.0,
  },
];

const mockFuelings = [
  {
    id: "fuel-1",
    carId: "car-123",
    date: "2026-09-01T12:00:00Z",
    fuelType: "Gasolina Comum",
    liters: 40,
    pricePerLiter: 5.5,
    totalCost: 220.0,
    isFullTank: true,
    gasStation: "Posto Shell",
    notes: "Abastecimento pós-estrada",
  },
  {
    id: "fuel-2",
    carId: "car-123",
    date: "2026-09-05T12:00:00Z",
    fuelType: "Gasolina Aditivada",
    liters: 30,
    pricePerLiter: 6.0,
    totalCost: 180.0,
    isFullTank: false,
    gasStation: "Posto Ipiranga",
  },
];

describe("VehicleDetailsPage", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders vehicle details, specs, FIPE info and financial KPIs correctly", async () => {
    vi.mocked(carService.get).mockResolvedValue(mockCar);
    vi.mocked(maintenanceService.listByCar).mockResolvedValue(mockMaintenances);
    vi.mocked(fuelService.listByCar).mockResolvedValue(mockFuelings);

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

      // Financial KPIs:
      // Fuel: 220 + 180 = 400.00 (2 abastecimentos)
      expect(screen.getByText("Gastos com Combustível")).toBeInTheDocument();
      expect(screen.getByText(/400,00/)).toBeInTheDocument();
      expect(screen.getByText("2 abastecimentos")).toBeInTheDocument();

      // Maintenance: 350 + 450 = 800.00 (2 serviços)
      expect(screen.getByText("Gastos com Manutenção")).toBeInTheDocument();
      expect(screen.getByText(/800,00/)).toBeInTheDocument();
      expect(screen.getByText("2 serviços")).toBeInTheDocument();

      // Operational Total: 400 + 800 = 1200.00
      expect(screen.getByText("Custo Operacional Total")).toBeInTheDocument();
      expect(screen.getByText(/1\.200,00/)).toBeInTheDocument();
    });
  });

  it("renders fuelings list by default on Tab 0", async () => {
    vi.mocked(carService.get).mockResolvedValue(mockCar);
    vi.mocked(maintenanceService.listByCar).mockResolvedValue([]);
    vi.mocked(fuelService.listByCar).mockResolvedValue(mockFuelings);

    render(
      <BrowserRouter>
        <VehicleDetailsPage />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(
        screen.getByText("Histórico de Abastecimentos"),
      ).toBeInTheDocument();
      expect(screen.getByText("Gasolina Comum")).toBeInTheDocument();
      expect(screen.getByText("Gasolina Aditivada")).toBeInTheDocument();
      expect(screen.getByText("Posto Shell")).toBeInTheDocument();
      expect(screen.getByText("Posto Ipiranga")).toBeInTheDocument();
      expect(screen.getByText("Tanque Cheio")).toBeInTheDocument();
    });
  });

  it("renders empty state when there are no fuelings", async () => {
    vi.mocked(carService.get).mockResolvedValue(mockCar);
    vi.mocked(maintenanceService.listByCar).mockResolvedValue([]);
    vi.mocked(fuelService.listByCar).mockResolvedValue([]);

    render(
      <BrowserRouter>
        <VehicleDetailsPage />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(
        screen.getByText("Nenhum abastecimento registrado para este veículo"),
      ).toBeInTheDocument();
      expect(
        screen.getByRole("button", {
          name: /Registrar Primeiro Abastecimento/i,
        }),
      ).toBeInTheDocument();
    });
  });

  it("switches to Tab 1 (Manutenções & Revisões) and renders maintenance history", async () => {
    vi.mocked(carService.get).mockResolvedValue(mockCar);
    vi.mocked(maintenanceService.listByCar).mockResolvedValue(mockMaintenances);
    vi.mocked(fuelService.listByCar).mockResolvedValue(mockFuelings);

    render(
      <BrowserRouter>
        <VehicleDetailsPage />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(
        screen.getByText("Abastecimentos & Combustível"),
      ).toBeInTheDocument();
    });

    const maintTab = screen.getByRole("tab", {
      name: /Manutenções & Revisões/i,
    });
    fireEvent.click(maintTab);

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

  it("renders empty state on Tab 1 when there are no maintenances", async () => {
    vi.mocked(carService.get).mockResolvedValue(mockCar);
    vi.mocked(maintenanceService.listByCar).mockResolvedValue([]);
    vi.mocked(fuelService.listByCar).mockResolvedValue(mockFuelings);

    render(
      <BrowserRouter>
        <VehicleDetailsPage />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(screen.getByText("Manutenções & Revisões")).toBeInTheDocument();
    });

    const maintTab = screen.getByRole("tab", {
      name: /Manutenções & Revisões/i,
    });
    fireEvent.click(maintTab);

    await waitFor(() => {
      expect(
        screen.getByText("Nenhuma manutenção registrada para este veículo"),
      ).toBeInTheDocument();
      expect(
        screen.getByRole("button", { name: /Registrar Primeira Manutenção/i }),
      ).toBeInTheDocument();
    });
  });

  it("renders error state when car is not found", async () => {
    vi.mocked(carService.get).mockRejectedValue(new Error("Not found"));
    vi.mocked(maintenanceService.listByCar).mockResolvedValue([]);
    vi.mocked(fuelService.listByCar).mockResolvedValue([]);

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
    vi.mocked(fuelService.listByCar).mockResolvedValue([]);

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

  it("navigates to maintenance registration page when 'Registrar Manutenção' button is clicked on Tab 1", async () => {
    vi.mocked(carService.get).mockResolvedValue(mockCar);
    vi.mocked(maintenanceService.listByCar).mockResolvedValue([]);
    vi.mocked(fuelService.listByCar).mockResolvedValue([]);

    render(
      <BrowserRouter>
        <VehicleDetailsPage />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(screen.getByText("Manutenções & Revisões")).toBeInTheDocument();
    });

    fireEvent.click(
      screen.getByRole("tab", { name: /Manutenções & Revisões/i }),
    );

    await waitFor(() => {
      const buttons = screen.getAllByRole("button", {
        name: /^Registrar Manutenção$/i,
      });
      expect(buttons.length).toBeGreaterThanOrEqual(1);
    });

    const buttons = screen.getAllByRole("button", {
      name: /^Registrar Manutenção$/i,
    });
    fireEvent.click(buttons[0]);

    expect(mockNavigate).toHaveBeenCalledWith(
      "/vehicles/car-123/maintenance/new",
    );
  });

  it("opens delete dialog, confirms deletion and updates fueling list", async () => {
    vi.mocked(carService.get).mockResolvedValue(mockCar);
    vi.mocked(maintenanceService.listByCar).mockResolvedValue([]);
    vi.mocked(fuelService.listByCar).mockResolvedValue(mockFuelings);
    vi.mocked(fuelService.delete).mockResolvedValue();

    render(
      <BrowserRouter>
        <VehicleDetailsPage />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(screen.getByText("Gasolina Aditivada")).toBeInTheDocument();
    });

    const deleteButtons = screen.getAllByRole("button", {
      name: "Excluir abastecimento",
    });
    fireEvent.click(deleteButtons[0]);

    // Dialog appears
    expect(screen.getByText("Excluir Abastecimento")).toBeInTheDocument();
    expect(
      screen.getByText(/Tem certeza que deseja excluir o abastecimento/i),
    ).toBeInTheDocument();

    // Click confirm Excluir inside dialog
    const confirmDeleteBtn = screen.getByRole("button", { name: "Excluir" });
    fireEvent.click(confirmDeleteBtn);

    await waitFor(() => {
      expect(fuelService.delete).toHaveBeenCalledWith("fuel-2");
      expect(screen.queryByText("Gasolina Aditivada")).not.toBeInTheDocument();
      expect(
        screen.getByText("Abastecimento removido com sucesso"),
      ).toBeInTheDocument();
    });
  });

  it("opens delete dialog for maintenance on Tab 1, confirms deletion and updates maintenance list", async () => {
    vi.mocked(carService.get).mockResolvedValue(mockCar);
    vi.mocked(maintenanceService.listByCar).mockResolvedValue(mockMaintenances);
    vi.mocked(fuelService.listByCar).mockResolvedValue([]);
    vi.mocked(maintenanceService.delete).mockResolvedValue();

    render(
      <BrowserRouter>
        <VehicleDetailsPage />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(screen.getByText("Manutenções & Revisões")).toBeInTheDocument();
    });

    fireEvent.click(
      screen.getByRole("tab", { name: /Manutenções & Revisões/i }),
    );

    await waitFor(() => {
      expect(
        screen.getByText("Troca de Pastilhas de Freio"),
      ).toBeInTheDocument();
    });

    const deleteButtons = screen.getAllByRole("button", {
      name: "Excluir manutenção",
    });
    fireEvent.click(deleteButtons[0]);

    expect(screen.getByText("Excluir Manutenção")).toBeInTheDocument();

    const confirmDeleteBtn = screen.getByRole("button", { name: "Excluir" });
    fireEvent.click(confirmDeleteBtn);

    await waitFor(() => {
      expect(maintenanceService.delete).toHaveBeenCalledWith("maint-2");
      expect(
        screen.queryByText("Troca de Pastilhas de Freio"),
      ).not.toBeInTheDocument();
      expect(
        screen.getByText("Manutenção removida com sucesso"),
      ).toBeInTheDocument();
    });
  });
});
