import { render, screen, waitFor, fireEvent } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { MaintenancePage } from "./MaintenancePage";
import { carService } from "../../cars/services/car.service";
import { maintenanceService } from "../services/maintenance.service";

vi.mock("../../cars/services/car.service", () => ({
  carService: {
    list: vi.fn(),
  },
}));

vi.mock("../services/maintenance.service", () => ({
  maintenanceService: {
    getUpcomingSummary: vi.fn(),
    listByCar: vi.fn(),
    create: vi.fn(),
    delete: vi.fn(),
  },
}));

describe("MaintenancePage", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("displays upcoming preventive maintenance cards and dynamic fleet alert status", async () => {
    vi.mocked(carService.list).mockResolvedValue([
      {
        id: "car-1",
        ownerId: "user-1",
        name: "Honda Civic",
        manufacturer: "Honda",
        model: "Civic Touring",
        yearManufacture: 2022,
        yearModel: 2023,
        lastMileage: 49500,
        vehicleType: "cars",
      },
    ]);

    vi.mocked(maintenanceService.getUpcomingSummary).mockResolvedValue({
      totalVehicles: 1,
      status: "upcoming",
      okCount: 4,
      upcomingCount: 1,
      overdueCount: 0,
      vehicles: [
        {
          carId: "car-1",
          carName: "Honda Civic",
          manufacturer: "Honda",
          model: "Civic Touring",
          vehicleType: "cars",
          currentKm: 49500,
          status: "upcoming",
          okCount: 4,
          upcomingCount: 1,
          overdueCount: 0,
          items: [
            {
              type: "Óleo de Motor",
              label: "Troca de óleo",
              status: "upcoming",
              intervalKm: 10000,
              intervalMonths: 12,
              hasHistory: true,
              lastPerformedAt: "2026-03-10T00:00:00Z",
              lastMileage: 40000,
              nextDueMileage: 50000,
              nextDueDate: "2027-03-10T00:00:00Z",
              remainingKm: 500,
              remainingDays: 140,
              progressPercent: 95,
              alertEligible: true,
            },
          ],
        },
      ],
    });

    vi.mocked(maintenanceService.listByCar).mockResolvedValue([
      {
        id: "maint-1",
        carId: "car-1",
        title: "Revisão 40.000 km",
        description: "Óleo sintético",
        date: "2026-03-10T00:00:00Z",
        mileage: 40000,
        types: ["Óleo de Motor"],
        cost: 380,
      },
    ]);

    render(
      <MemoryRouter initialEntries={["/maintenance?carId=car-1"]}>
        <MaintenancePage />
      </MemoryRouter>,
    );

    await waitFor(() => {
      expect(
        screen.getByText("Revisões Preventivas Próximas"),
      ).toBeInTheDocument();
    });

    expect(
      screen.getByText("Próximas Manutenções Preventivas"),
    ).toBeInTheDocument();
    expect(screen.getByText("Troca de óleo")).toBeInTheDocument();
    expect(screen.getByText("Revisão 40.000 km")).toBeInTheDocument();

    // Clicking quick register button opens AddMaintenanceDialog with preselected type
    const quickBtn = screen.getByRole("button", {
      name: /Registrar troca de óleo/i,
    });
    fireEvent.click(quickBtn);

    await waitFor(() => {
      expect(
        screen.getAllByDisplayValue("Óleo de Motor").length,
      ).toBeGreaterThanOrEqual(1);
    });
  });
});
