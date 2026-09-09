import { describe, expect, it, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { MaintenanceCard } from "./MaintenanceCard";
import { Maintenance } from "../types/maintenance.types";

describe("MaintenanceCard", () => {
  const mockMaintenance: Maintenance = {
    id: "maint-1",
    carId: "car-123",
    title: "Revisão dos 40.000 km",
    description: "Troca completa de filtros e fluidos",
    date: "2026-08-20T12:00:00Z",
    mileage: 40000,
    types: ["Óleo de Motor", "Filtro do Óleo de Motor", "Pastilhas de Freio"],
    cost: 780.5,
    attachments: [
      {
        id: "att-1",
        name: "nota_fiscal.pdf",
        size: 1024 * 500,
        mimeType: "application/pdf",
        dataUrl: "data:application/pdf;base64,JVBERi0xLjQK",
        createdAt: "2026-08-20T12:00:00Z",
      },
      {
        id: "att-2",
        name: "foto_pecas.jpg",
        size: 1024 * 300,
        mimeType: "image/jpeg",
        dataUrl: "data:image/jpeg;base64,/9j/4AAQSkZJRg==",
        createdAt: "2026-08-20T12:00:00Z",
      },
    ],
  };

  it("renders title, formatted date and mileage", () => {
    render(<MaintenanceCard maintenance={mockMaintenance} />);

    expect(screen.getByText("Revisão dos 40.000 km")).toBeInTheDocument();
    expect(screen.getByText("20/08/2026")).toBeInTheDocument();
    expect(screen.getByText("40.000 km")).toBeInTheDocument();
    expect(
      screen.getByText("Troca completa de filtros e fluidos"),
    ).toBeInTheDocument();
  });

  it("renders maintenance types as chips", () => {
    render(<MaintenanceCard maintenance={mockMaintenance} />);

    expect(screen.getByText("Óleo de Motor")).toBeInTheDocument();
    expect(screen.getByText("Filtro do Óleo de Motor")).toBeInTheDocument();
    expect(screen.getByText("Pastilhas de Freio")).toBeInTheDocument();
  });

  it("renders formatted total cost when provided", () => {
    render(<MaintenanceCard maintenance={mockMaintenance} />);

    expect(screen.getByText(/780,50/)).toBeInTheDocument();
  });

  it("renders attachments buttons and handles click", () => {
    const originalOpen = window.open;
    window.open = vi.fn();

    render(<MaintenanceCard maintenance={mockMaintenance} />);

    expect(screen.getByText(/Comprovantes \(2\)/i)).toBeInTheDocument();
    const pdfBtn = screen.getByRole("button", { name: /nota_fiscal.pdf/i });
    const imgBtn = screen.getByRole("button", { name: /foto_pecas.jpg/i });

    expect(pdfBtn).toBeInTheDocument();
    expect(imgBtn).toBeInTheDocument();

    fireEvent.click(pdfBtn);
    expect(window.open).toHaveBeenCalled();

    window.open = originalOpen;
  });

  it("renders action buttons when onEdit and onDelete callbacks are provided", () => {
    const handleEdit = vi.fn();
    const handleDelete = vi.fn();

    render(
      <MaintenanceCard
        maintenance={mockMaintenance}
        onEdit={handleEdit}
        onDelete={handleDelete}
      />,
    );

    const editBtn = screen.getByRole("button", { name: "Editar manutenção" });
    const deleteBtn = screen.getByRole("button", {
      name: "Excluir manutenção",
    });

    expect(editBtn).toBeInTheDocument();
    expect(deleteBtn).toBeInTheDocument();

    fireEvent.click(editBtn);
    expect(handleEdit).toHaveBeenCalledWith(mockMaintenance);

    fireEvent.click(deleteBtn);
    expect(handleDelete).toHaveBeenCalledWith(mockMaintenance.id);
  });

  it("does not render action buttons when onEdit and onDelete are omitted", () => {
    render(<MaintenanceCard maintenance={mockMaintenance} />);

    expect(
      screen.queryByRole("button", { name: "Editar manutenção" }),
    ).not.toBeInTheDocument();
    expect(
      screen.queryByRole("button", { name: "Excluir manutenção" }),
    ).not.toBeInTheDocument();
  });
});
