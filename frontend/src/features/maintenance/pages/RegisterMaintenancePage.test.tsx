import { describe, expect, it, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { BrowserRouter } from "react-router-dom";
import { RegisterMaintenancePage } from "./RegisterMaintenancePage";
import { carService } from "../../cars/services/car.service";
import { maintenanceService } from "../services/maintenance.service";

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

vi.mock("../../cars/services/car.service", () => ({
  carService: {
    get: vi.fn(),
  },
}));

vi.mock("../services/maintenance.service", () => ({
  maintenanceService: {
    create: vi.fn(),
  },
}));

const mockCar = {
  id: "car-123",
  name: "Celtinha",
  manufacturer: "Chevrolet",
  model: "Celta Life",
  yearManufacture: 2012,
  yearModel: 2013,
  lastMileage: 115000,
  vehicleType: "cars",
  ownerId: "user-1",
};

describe("RegisterMaintenancePage", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders page title, vehicle info and all category groups with chips", async () => {
    vi.mocked(carService.get).mockResolvedValue(mockCar);

    render(
      <BrowserRouter>
        <RegisterMaintenancePage />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(screen.getByText("Registrar Manutenção")).toBeInTheDocument();
      expect(screen.getByText(/Celtinha/i)).toBeInTheDocument();
    });

    // Categories
    expect(screen.getByText("Motor & Filtros")).toBeInTheDocument();
    expect(screen.getByText("Fluidos & Arrefecimento")).toBeInTheDocument();
    expect(screen.getByText("Freios, Suspensão & Pneus")).toBeInTheDocument();
    expect(screen.getByText("Elétrica, Conforto & Outros")).toBeInTheDocument();

    // Chips
    expect(screen.getByText("Óleo de Motor")).toBeInTheDocument();
    expect(screen.getByText("Fluido de Freio")).toBeInTheDocument();
    expect(screen.getByText("Pastilhas de Freio")).toBeInTheDocument();
    expect(screen.getByText("Bateria")).toBeInTheDocument();
  });

  it("toggles chips and auto-suggests service title", async () => {
    vi.mocked(carService.get).mockResolvedValue(mockCar);

    render(
      <BrowserRouter>
        <RegisterMaintenancePage />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(screen.getByText(/Celtinha/i)).toBeInTheDocument();
    });

    const oleoChip = screen.getByText("Óleo de Motor");
    fireEvent.click(oleoChip);

    const titleInput = screen.getByLabelText(
      /Título do Serviço/i,
    ) as HTMLInputElement;
    expect(titleInput.value).toBe("Óleo de Motor");

    // Click another chip
    const filtroOleoChip = screen.getByText("Filtro do Óleo de Motor");
    fireEvent.click(filtroOleoChip);

    expect(titleInput.value).toBe("Óleo de Motor, Filtro do Óleo de Motor");
  });

  it("shows conditional custom type text field when 'Outro tipo de manutenção' is clicked", async () => {
    vi.mocked(carService.get).mockResolvedValue(mockCar);

    render(
      <BrowserRouter>
        <RegisterMaintenancePage />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(screen.getByText(/Celtinha/i)).toBeInTheDocument();
    });

    expect(
      screen.queryByLabelText(/Especifique o outro tipo de manutenção/i),
    ).not.toBeInTheDocument();

    const outroChip = screen.getByText("Outro tipo de manutenção");
    fireEvent.click(outroChip);

    expect(
      await screen.findByLabelText(/Especifique o outro tipo de manutenção/i),
    ).toBeInTheDocument();
  });

  it("submits maintenance form and navigates back to vehicle details", async () => {
    vi.mocked(carService.get).mockResolvedValue(mockCar);
    vi.mocked(maintenanceService.create).mockResolvedValue({
      id: "maint-new",
      carId: "car-123",
      title: "Revisão Geral",
      description: "Tudo verificado",
      date: "2026-09-04T12:00:00Z",
      mileage: 120000,
      cost: 450,
      types: ["Revisão Geral / Preventiva"],
    });

    render(
      <BrowserRouter>
        <RegisterMaintenancePage />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(screen.getByText(/Celtinha/i)).toBeInTheDocument();
    });

    // Select Revision chip
    const revisaoChip = screen.getByText("Revisão Geral / Preventiva");
    fireEvent.click(revisaoChip);

    // Set cost
    const costInput = screen.getByLabelText(/Valor Total/i);
    fireEvent.change(costInput, { target: { value: "450.00" } });

    // Submit form
    const submitBtn = screen.getAllByRole("button", {
      name: /^Salvar Manutenção$/i,
    })[0];
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(maintenanceService.create).toHaveBeenCalledWith(
        "car-123",
        expect.objectContaining({
          title: "Revisão Geral / Preventiva",
          cost: 450,
          types: ["Revisão Geral / Preventiva"],
        }),
      );
      expect(mockNavigate).toHaveBeenCalledWith("/vehicles/car-123");
    });
  });

  it("validates PDF file size limit (2MB) on upload", async () => {
    vi.mocked(carService.get).mockResolvedValue(mockCar);

    render(
      <BrowserRouter>
        <RegisterMaintenancePage />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(screen.getByText(/Celtinha/i)).toBeInTheDocument();
    });

    const fileInput = document.querySelector(
      'input[type="file"]',
    ) as HTMLInputElement;

    const oversizedPdf = new File(["data".repeat(600000)], "nota_fiscal.pdf", {
      type: "application/pdf",
    });
    Object.defineProperty(oversizedPdf, "size", { value: 3 * 1024 * 1024 });

    fireEvent.change(fileInput, { target: { files: [oversizedPdf] } });

    await waitFor(() => {
      expect(
        screen.getByText(/excede o limite máximo de 2MB/i),
      ).toBeInTheDocument();
    });
  });
});
