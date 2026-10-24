import { render, screen, fireEvent } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { UpcomingMaintenanceCard } from "./UpcomingMaintenanceCard";
import { UpcomingMaintenanceItem } from "../types/maintenance.types";

describe("UpcomingMaintenanceCard", () => {
  const sampleItem: UpcomingMaintenanceItem = {
    type: "Óleo de Motor",
    label: "Troca de óleo",
    status: "upcoming",
    intervalKm: 10000,
    intervalMonths: 12,
    hasHistory: true,
    lastPerformedAt: "2026-04-15T12:00:00Z",
    lastMileage: 40000,
    nextDueMileage: 50000,
    nextDueDate: "2027-04-15T12:00:00Z",
    remainingKm: 600,
    remainingDays: 170,
    progressPercent: 94,
    alertEligible: false,
  };

  it("renders status chip, remaining km, next due mileage and triggers quick register", () => {
    const onQuickRegister = vi.fn();
    render(
      <UpcomingMaintenanceCard
        item={sampleItem}
        onQuickRegister={onQuickRegister}
      />,
    );

    expect(screen.getByText("Troca de óleo")).toBeInTheDocument();
    expect(screen.getByText("Próximo")).toBeInTheDocument();
    expect(screen.getByText(/Faltam 600 km/i)).toBeInTheDocument();
    expect(screen.getByText(/50\.000 km/i)).toBeInTheDocument();
    expect(screen.getByText(/94% do ciclo/i)).toBeInTheDocument();

    const btn = screen.getByRole("button", {
      name: /Registrar troca de óleo/i,
    });
    fireEvent.click(btn);
    expect(onQuickRegister).toHaveBeenCalledWith("Óleo de Motor");
  });

  it("renders overdue state when remainingKm is negative", () => {
    render(
      <UpcomingMaintenanceCard
        item={{
          ...sampleItem,
          status: "overdue",
          remainingKm: -350,
          progressPercent: 100,
        }}
      />,
    );

    expect(screen.getByText("Vencido")).toBeInTheDocument();
    expect(screen.getByText(/Excedido em 350 km/i)).toBeInTheDocument();
  });
});
