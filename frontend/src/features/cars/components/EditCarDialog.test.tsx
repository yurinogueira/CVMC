import { describe, expect, it, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { EditCarDialog } from "./EditCarDialog";
import { carService } from "../services/car.service";
import { Car } from "../types/car.types";

vi.mock("../services/car.service", () => ({
  carService: {
    update: vi.fn(),
  },
}));

const mockCar: Car = {
  id: "car-123",
  ownerId: "user-1",
  name: "Civic Antigo",
  manufacturer: "Honda",
  model: "Civic Touring",
  yearManufacture: 2020,
  yearModel: 2021,
  lastMileage: 50000,
  fipeCode: "005487-9",
  fipePrice: "R$ 130.000,00",
  fuel: "Gasolina",
};

describe("EditCarDialog", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders with prefilled vehicle information", () => {
    render(
      <EditCarDialog
        open={true}
        car={mockCar}
        onClose={vi.fn()}
        onCarUpdated={vi.fn()}
      />,
    );

    expect(screen.getByText("Editar Veículo")).toBeInTheDocument();
    expect(screen.getByText("Honda • Civic Touring")).toBeInTheDocument();
    expect(screen.getByText("005487-9")).toBeInTheDocument();
    expect(screen.getByText("R$ 130.000,00")).toBeInTheDocument();

    const nameInput = screen.getByLabelText(
      /Apelido \/ Identificador/i,
    ) as HTMLInputElement;
    expect(nameInput.value).toBe("Civic Antigo");

    const mileageInput = screen.getByLabelText(/Km Atual/i) as HTMLInputElement;
    expect(mileageInput.value).toBe("50000");

    const yearManufactureInput = screen.getByLabelText(
      /Ano de Fabricação/i,
    ) as HTMLInputElement;
    expect(yearManufactureInput.value).toBe("2020");

    const yearModelInput = screen.getByLabelText(
      /Ano do Modelo/i,
    ) as HTMLInputElement;
    expect(yearModelInput.value).toBe("2021");
  });

  it("updates vehicle and calls onCarUpdated when submitted", async () => {
    const handleClose = vi.fn();
    const handleCarUpdated = vi.fn();

    const updatedCarMock: Car = {
      ...mockCar,
      name: "Civic Atualizado",
      lastMileage: 55000,
    };

    vi.mocked(carService.update).mockResolvedValueOnce(updatedCarMock);

    render(
      <EditCarDialog
        open={true}
        car={mockCar}
        onClose={handleClose}
        onCarUpdated={handleCarUpdated}
      />,
    );

    const nameInput = screen.getByLabelText(/Apelido \/ Identificador/i);
    fireEvent.change(nameInput, { target: { value: "Civic Atualizado" } });

    const mileageInput = screen.getByLabelText(/Km Atual/i);
    fireEvent.change(mileageInput, { target: { value: "55000" } });

    const saveButton = screen.getByRole("button", {
      name: /Salvar Alterações/i,
    });
    fireEvent.click(saveButton);

    await waitFor(() => {
      expect(carService.update).toHaveBeenCalledWith("car-123", {
        name: "Civic Atualizado",
        manufacturer: "Honda",
        model: "Civic Touring",
        yearManufacture: 2020,
        yearModel: 2021,
        lastMileage: 55000,
        vehicleType: undefined,
        fipeCode: "005487-9",
        fipePrice: "R$ 130.000,00",
        fuel: "Gasolina",
      });
      expect(handleCarUpdated).toHaveBeenCalledWith(updatedCarMock);
      expect(handleClose).toHaveBeenCalled();
    });
  });
});
