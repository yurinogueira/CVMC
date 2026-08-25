import { describe, expect, it, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { VehicleCard } from "./VehicleCard";
import { Car } from "../types/car.types";

const mockCar: Car = {
  id: "car-123",
  ownerId: "user-1",
  name: "Meu Civic",
  manufacturer: "Honda",
  model: "Civic Touring",
  yearManufacture: 2021,
  yearModel: 2022,
  lastMileage: 45000,
  fipePrice: "R$ 140.000,00",
};

describe("VehicleCard", () => {
  it("renders car information correctly", () => {
    render(<VehicleCard car={mockCar} />);

    expect(screen.getByText("Meu Civic")).toBeInTheDocument();
    expect(screen.getByText("Honda • Civic Touring")).toBeInTheDocument();
    expect(screen.getByText("2021/2022")).toBeInTheDocument();
    expect(screen.getByText("45.000 km")).toBeInTheDocument();
    expect(screen.getByText("FIPE: R$ 140.000,00")).toBeInTheDocument();
  });

  it("calls onEdit when edit button is clicked", () => {
    const handleEdit = vi.fn();
    render(<VehicleCard car={mockCar} onEdit={handleEdit} />);

    const editBtn = screen.getByRole("button", { name: /Editar veículo/i });
    expect(editBtn).toBeInTheDocument();
    fireEvent.click(editBtn);

    expect(handleEdit).toHaveBeenCalledWith(mockCar);
  });

  it("calls onDelete when delete button is clicked", () => {
    const handleDelete = vi.fn();
    render(<VehicleCard car={mockCar} onDelete={handleDelete} />);

    const deleteBtn = screen.getByRole("button", { name: /Remover veículo/i });
    expect(deleteBtn).toBeInTheDocument();
    fireEvent.click(deleteBtn);

    expect(handleDelete).toHaveBeenCalledWith("car-123");
  });
});
