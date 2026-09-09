import { describe, expect, it, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { AddFuelingDialog } from "./AddFuelingDialog";
import { fuelService } from "../services/fuel.service";

vi.mock("../services/fuel.service", () => ({
  fuelService: {
    create: vi.fn(),
    update: vi.fn(),
  },
}));

describe("AddFuelingDialog", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("calculates TotalCost automatically when typing Liters and PricePerLiter", () => {
    render(
      <AddFuelingDialog
        open={true}
        onClose={vi.fn()}
        onSuccess={vi.fn()}
        carId="car-123"
      />,
    );

    const litersInput = screen.getByLabelText(/Litros Abastecidos/i);
    const priceInput = screen.getByLabelText(/Preço por Litro/i);
    const totalInput = screen.getByLabelText(/Valor Total Pago/i);

    fireEvent.change(litersInput, { target: { value: "40" } });
    fireEvent.change(priceInput, { target: { value: "5.50" } });

    expect((totalInput as HTMLInputElement).value).toBe("220.00");
  });

  it("calculates PricePerLiter automatically when typing Liters and TotalCost", () => {
    render(
      <AddFuelingDialog
        open={true}
        onClose={vi.fn()}
        onSuccess={vi.fn()}
        carId="car-123"
      />,
    );

    const litersInput = screen.getByLabelText(/Litros Abastecidos/i);
    const totalInput = screen.getByLabelText(/Valor Total Pago/i);
    const priceInput = screen.getByLabelText(/Preço por Litro/i);

    fireEvent.change(litersInput, { target: { value: "50" } });
    fireEvent.change(totalInput, { target: { value: "250.00" } });

    expect((priceInput as HTMLInputElement).value).toBe("5.000");
  });

  it("submits fueling data on create", async () => {
    const handleSuccess = vi.fn();
    const handleClose = vi.fn();

    vi.mocked(fuelService.create).mockResolvedValue({
      id: "fuel-99",
      carId: "car-123",
      date: "2026-09-08T00:00:00.000Z",
      fuelType: "Gasolina Comum",
      liters: 40,
      pricePerLiter: 5.5,
      totalCost: 220.0,
      isFullTank: true,
      gasStation: "Shell",
    });

    render(
      <AddFuelingDialog
        open={true}
        onClose={handleClose}
        onSuccess={handleSuccess}
        carId="car-123"
      />,
    );

    fireEvent.change(screen.getByLabelText(/Litros Abastecidos/i), {
      target: { value: "40" },
    });
    fireEvent.change(screen.getByLabelText(/Valor Total Pago/i), {
      target: { value: "220" },
    });

    fireEvent.click(
      screen.getByRole("button", { name: /Registrar Abastecimento/i }),
    );

    await waitFor(() => {
      expect(fuelService.create).toHaveBeenCalledWith(
        "car-123",
        expect.objectContaining({
          fuelType: "Gasolina Comum",
          liters: 40,
          totalCost: 220,
        }),
      );
      expect(handleSuccess).toHaveBeenCalled();
      expect(handleClose).toHaveBeenCalled();
    });
  });
});
