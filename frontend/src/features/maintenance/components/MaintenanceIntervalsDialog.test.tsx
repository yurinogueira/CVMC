import { describe, expect, it, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { MaintenanceIntervalsDialog } from "./MaintenanceIntervalsDialog";
import { carService } from "../../cars/services/car.service";
import { Car } from "../../cars/types/car.types";

vi.mock("../../cars/services/car.service", () => ({
  carService: {
    update: vi.fn(),
  },
}));

const mockCar: Car = {
  id: "car-123",
  ownerId: "user-1",
  name: "Corollinha",
  manufacturer: "Toyota",
  model: "Corolla XEi",
  yearManufacture: 2020,
  yearModel: 2021,
  lastMileage: 50000,
  vehicleType: "cars",
  maintenanceIntervals: {
    "Óleo de Motor": {
      intervalKm: 5000,
      intervalMonths: 6,
    },
  },
};

const mockMoto: Car = {
  id: "moto-456",
  ownerId: "user-1",
  name: "Yamaha MT-03",
  manufacturer: "Yamaha",
  model: "MT-03",
  yearManufacture: 2022,
  yearModel: 2023,
  lastMileage: 12000,
  vehicleType: "motorcycles",
};

describe("MaintenanceIntervalsDialog", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders all 5 periodic items for a car", () => {
    render(
      <MaintenanceIntervalsDialog
        open={true}
        car={mockCar}
        onClose={vi.fn()}
        onSaved={vi.fn()}
      />,
    );

    expect(screen.getByText("Intervalos de Revisão")).toBeInTheDocument();
    expect(screen.getByText("Troca de óleo")).toBeInTheDocument();
    expect(screen.getByText("Troca de filtro de óleo")).toBeInTheDocument();
    expect(screen.getByText("Troca de filtro de ar")).toBeInTheDocument();
    expect(
      screen.getByText("Troca do filtro de ar-condicionado"),
    ).toBeInTheDocument();
    expect(
      screen.getByText("Troca do filtro de combustível"),
    ).toBeInTheDocument();
    expect(screen.getByText("Personalizado")).toBeInTheDocument(); // Since engine oil is 5000 / 6
  });

  it("renders only 3 periodic items for a motorcycle", () => {
    render(
      <MaintenanceIntervalsDialog
        open={true}
        car={mockMoto}
        onClose={vi.fn()}
        onSaved={vi.fn()}
      />,
    );

    expect(screen.getByText("Troca de óleo")).toBeInTheDocument();
    expect(screen.getByText("Troca de filtro de óleo")).toBeInTheDocument();
    expect(screen.getByText("Troca de filtro de ar")).toBeInTheDocument();
    expect(
      screen.queryByText("Troca do filtro de ar-condicionado"),
    ).not.toBeInTheDocument();
    expect(
      screen.queryByText("Troca do filtro de combustível"),
    ).not.toBeInTheDocument();
  });

  it("calls carService.update and onSaved when saving intervals", async () => {
    const onSavedMock = vi.fn();
    const onCloseMock = vi.fn();
    const updatedMockCar = { ...mockCar };
    vi.mocked(carService.update).mockResolvedValueOnce(updatedMockCar);

    render(
      <MaintenanceIntervalsDialog
        open={true}
        car={mockCar}
        onClose={onCloseMock}
        onSaved={onSavedMock}
      />,
    );

    const saveButton = screen.getByRole("button", {
      name: /Salvar Intervalos/i,
    });
    fireEvent.click(saveButton);

    await waitFor(() => {
      expect(carService.update).toHaveBeenCalledWith(
        "car-123",
        expect.objectContaining({
          maintenanceIntervals: expect.objectContaining({
            "Óleo de Motor": {
              intervalKm: 5000,
              intervalMonths: 6,
            },
          }),
        }),
      );
      expect(onSavedMock).toHaveBeenCalledWith(updatedMockCar);
      expect(onCloseMock).toHaveBeenCalled();
    });
  });

  it("resets values to recommended defaults when clicking Restaurar Padrões Recomendados", async () => {
    const onCloseMock = vi.fn();
    const onSavedMock = vi.fn();
    vi.mocked(carService.update).mockResolvedValueOnce(mockCar);

    render(
      <MaintenanceIntervalsDialog
        open={true}
        car={mockCar}
        onClose={onCloseMock}
        onSaved={onSavedMock}
      />,
    );

    const resetButton = screen.getByRole("button", {
      name: /Restaurar Padrões Recomendados/i,
    });
    fireEvent.click(resetButton);

    const saveButton = screen.getByRole("button", {
      name: /Salvar Intervalos/i,
    });
    fireEvent.click(saveButton);

    await waitFor(() => {
      expect(carService.update).toHaveBeenCalledWith(
        "car-123",
        expect.objectContaining({
          maintenanceIntervals: expect.objectContaining({
            "Óleo de Motor": {
              intervalKm: 10000,
              intervalMonths: 12,
            },
          }),
        }),
      );
    });
  });

  it("displays error message when update fails", async () => {
    vi.mocked(carService.update).mockRejectedValueOnce({
      response: { data: { message: "Erro ao atualizar intervalos" } },
    });

    render(
      <MaintenanceIntervalsDialog
        open={true}
        car={mockCar}
        onClose={vi.fn()}
        onSaved={vi.fn()}
      />,
    );

    const saveButton = screen.getByRole("button", {
      name: /Salvar Intervalos/i,
    });
    fireEvent.click(saveButton);

    await waitFor(() => {
      expect(
        screen.getByText("Erro ao atualizar intervalos"),
      ).toBeInTheDocument();
    });
  });
});
