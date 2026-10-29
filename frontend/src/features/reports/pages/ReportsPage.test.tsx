import { describe, expect, it, vi, beforeEach } from "vitest";
import { render, screen, waitFor, fireEvent } from "@testing-library/react";
import { ReportsPage } from "./ReportsPage";
import { carService } from "../../cars/services/car.service";
import { reportService } from "../services/report.service";

vi.mock("../../cars/services/car.service", () => ({
  carService: {
    list: vi.fn(),
  },
}));

vi.mock("../services/report.service", () => ({
  reportService: {
    listJobs: vi.fn(),
    startMaintenanceReport: vi.fn(),
    downloadReport: vi.fn(),
  },
}));

describe("ReportsPage", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(carService.list).mockResolvedValue([
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
    ]);
    vi.mocked(reportService.listJobs).mockResolvedValue({
      jobs: [],
      total: 0,
      page: 1,
      limit: 10,
    });
  });

  it("renders page header and opens dialog when clicking export button", async () => {
    render(<ReportsPage />);

    expect(screen.getByText("Relatórios & Exportações")).toBeInTheDocument();

    await waitFor(() => {
      expect(
        screen.getByRole("button", { name: /Exportar Manutenções \(PDF\)/i }),
      ).toBeInTheDocument();
    });

    const exportBtn = screen.getByRole("button", {
      name: /Exportar Manutenções \(PDF\)/i,
    });
    fireEvent.click(exportBtn);

    await waitFor(() => {
      expect(
        screen.getByText("Exportar Dossiê de Manutenções (PDF)"),
      ).toBeInTheDocument();
    });
  });
});
