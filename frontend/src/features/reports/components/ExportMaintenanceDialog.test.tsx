import { describe, expect, it, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { ExportMaintenanceDialog } from "./ExportMaintenanceDialog";
import { reportService } from "../services/report.service";
import { Car } from "../../cars/types/car.types";

vi.mock("../services/report.service", () => ({
  reportService: {
    startMaintenanceReport: vi.fn(),
  },
}));

describe("ExportMaintenanceDialog", () => {
  const sampleCars: Car[] = [
    {
      id: "car-1",
      ownerId: "user-1",
      name: "Civic",
      manufacturer: "Honda",
      model: "Civic",
      yearManufacture: 2018,
      yearModel: 2019,
      lastMileage: 45000,
      createdAt: "2026-01-01T00:00:00Z",
      updatedAt: "2026-01-01T00:00:00Z",
    },
    {
      id: "car-2",
      ownerId: "user-1",
      name: "Gol",
      manufacturer: "VW",
      model: "Gol",
      yearManufacture: 2015,
      yearModel: 2016,
      lastMileage: 90000,
      createdAt: "2026-01-01T00:00:00Z",
      updatedAt: "2026-01-01T00:00:00Z",
    },
  ];

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders dialog and triggers startMaintenanceReport on submit", async () => {
    const mockOnSuccess = vi.fn();
    const mockOnClose = vi.fn();
    vi.mocked(reportService.startMaintenanceReport).mockResolvedValue({
      id: "job-new",
      userId: "user-1",
      carId: "car-1",
      carName: "Civic",
      type: "maintenance_pdf",
      status: "pending",
      createdAt: new Date().toISOString(),
    });

    render(
      <ExportMaintenanceDialog
        open={true}
        onClose={mockOnClose}
        cars={sampleCars}
        initialCarId="car-1"
        onSuccess={mockOnSuccess}
      />,
    );

    expect(
      screen.getByText("Exportar Dossiê de Manutenções (PDF)"),
    ).toBeInTheDocument();

    const submitBtn = screen.getByText("Solicitar Exportação");
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(reportService.startMaintenanceReport).toHaveBeenCalledWith(
        "car-1",
      );
      expect(mockOnSuccess).toHaveBeenCalled();
      expect(mockOnClose).toHaveBeenCalled();
    });
  });

  it("handles errors during startMaintenanceReport", async () => {
    vi.mocked(reportService.startMaintenanceReport).mockRejectedValue({
      response: { data: { message: "Veículo bloqueado para exportação" } },
    });

    render(
      <ExportMaintenanceDialog
        open={true}
        onClose={vi.fn()}
        cars={sampleCars}
        initialCarId="car-1"
        onSuccess={vi.fn()}
      />,
    );

    const submitBtn = screen.getByText("Solicitar Exportação");
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(
        screen.getByText("Veículo bloqueado para exportação"),
      ).toBeInTheDocument();
    });
  });
});
