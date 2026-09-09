import { describe, expect, it, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { FuelCard } from "./FuelCard";
import { Fueling } from "../types/fuel.types";

describe("FuelCard", () => {
  const mockFueling: Fueling = {
    id: "fuel-1",
    carId: "car-123",
    date: "2026-09-08T12:00:00Z",
    fuelType: "Gasolina Comum",
    liters: 45.5,
    pricePerLiter: 5.89,
    totalCost: 267.99,
    isFullTank: true,
    gasStation: "Posto Shell Central",
    notes: "Abastecido antes de viajar",
  };

  it("renders fuel type, date, liters, price per liter and total cost", () => {
    render(<FuelCard fueling={mockFueling} />);

    expect(screen.getByText("Gasolina Comum")).toBeInTheDocument();
    expect(screen.getByText("08/09/2026")).toBeInTheDocument();
    expect(screen.getByText(/45,5 Litros/)).toBeInTheDocument();
    expect(screen.getByText(/5,89/)).toBeInTheDocument();
    expect(screen.getByText(/267,99/)).toBeInTheDocument();
    expect(screen.getByText("Tanque Cheio")).toBeInTheDocument();
    expect(screen.getByText("Posto Shell Central")).toBeInTheDocument();
    expect(screen.getByText("Abastecido antes de viajar")).toBeInTheDocument();
  });

  it("renders action buttons and triggers callbacks when provided", () => {
    const handleEdit = vi.fn();
    const handleDelete = vi.fn();

    render(
      <FuelCard
        fueling={mockFueling}
        onEdit={handleEdit}
        onDelete={handleDelete}
      />,
    );

    const editBtn = screen.getByRole("button", {
      name: "Editar abastecimento",
    });
    const deleteBtn = screen.getByRole("button", {
      name: "Excluir abastecimento",
    });

    expect(editBtn).toBeInTheDocument();
    expect(deleteBtn).toBeInTheDocument();

    fireEvent.click(editBtn);
    expect(handleEdit).toHaveBeenCalledWith(mockFueling);

    fireEvent.click(deleteBtn);
    expect(handleDelete).toHaveBeenCalledWith(mockFueling.id);
  });
});
