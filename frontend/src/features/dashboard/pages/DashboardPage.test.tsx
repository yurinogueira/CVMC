import { describe, expect, it, vi } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import { BrowserRouter } from "react-router-dom";
import { DashboardPage } from "./DashboardPage";

vi.mock("../../cars/services/car.service", () => ({
  carService: {
    list: vi.fn().mockResolvedValue([
      {
        id: "1",
        name: "Carro do Trabalho",
        manufacturer: "Toyota",
        model: "Corolla",
        yearManufacture: 2022,
        yearModel: 2023,
        lastMileage: 25000,
        ownerId: "user-1",
      },
    ]),
    delete: vi.fn().mockResolvedValue(undefined),
  },
}));

describe("DashboardPage", () => {
  it("renders KPIs and vehicle list when loaded", async () => {
    render(
      <BrowserRouter>
        <DashboardPage />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(screen.getByText("Veículos Cadastrados")).toBeInTheDocument();
      expect(screen.getByText("Carro do Trabalho")).toBeInTheDocument();
      expect(screen.getByText("Toyota • Corolla")).toBeInTheDocument();
    });
  });
});
