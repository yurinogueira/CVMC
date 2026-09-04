import { describe, expect, it, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { AddMaintenanceDialog } from "./AddMaintenanceDialog";
import { maintenanceService } from "../services/maintenance.service";

vi.mock("../services/maintenance.service", () => ({
  maintenanceService: {
    create: vi.fn(),
  },
}));

describe("AddMaintenanceDialog", () => {
  const mockOnClose = vi.fn();
  const mockOnMaintenanceCreated = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders all form elements correctly", () => {
    render(
      <AddMaintenanceDialog
        open={true}
        carId="car-1"
        carName="Honda Civic"
        lastMileage={50000}
        onClose={mockOnClose}
        onMaintenanceCreated={mockOnMaintenanceCreated}
      />,
    );

    expect(screen.getByText("Registrar Manutenção")).toBeInTheDocument();
    expect(screen.getByText(/Honda Civic/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Tipo de Manutenção/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Título do Serviço/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Data da Realização/i)).toBeInTheDocument();
    expect(
      screen.getByLabelText(/Quilometragem no Momento do Serviço/i),
    ).toBeInTheDocument();
    expect(screen.getByLabelText(/Custo Total/i)).toBeInTheDocument();
    expect(screen.getByText("Anexar Comprovantes")).toBeInTheDocument();
  });

  it("opens conditional custom type input when 'Outro tipo de manutenção' is selected", async () => {
    render(
      <AddMaintenanceDialog
        open={true}
        carId="car-1"
        onClose={mockOnClose}
        onMaintenanceCreated={mockOnMaintenanceCreated}
      />,
    );

    // Initial state: custom type input should not be visible
    expect(
      screen.queryByLabelText(/Especifique o outro tipo de manutenção/i),
    ).not.toBeInTheDocument();

    // Click on Select Tipo de Manutenção
    const select = screen.getByLabelText(/Tipo de Manutenção/i);
    fireEvent.mouseDown(select);

    // Click on Outro tipo de manutenção in menu
    const outroOption = await screen.findByRole("option", {
      name: /Outro tipo de manutenção/i,
    });
    fireEvent.click(outroOption);

    // Custom type input should now appear
    expect(
      await screen.findByLabelText(/Especifique o outro tipo de manutenção/i),
    ).toBeInTheDocument();
  });

  it("validates required title and creates maintenance successfully", async () => {
    vi.mocked(maintenanceService.create).mockResolvedValue({
      id: "maint-10",
      carId: "car-1",
      title: "Troca de Óleo de Motor",
      description: "Óleo 0W20",
      date: "2026-09-04T12:00:00.000Z",
      mileage: 52000,
      cost: 350.5,
      types: ["Óleo de Motor"],
    });

    render(
      <AddMaintenanceDialog
        open={true}
        carId="car-1"
        lastMileage={50000}
        onClose={mockOnClose}
        onMaintenanceCreated={mockOnMaintenanceCreated}
      />,
    );

    // Select Óleo de Motor
    const select = screen.getByLabelText(/Tipo de Manutenção/i);
    fireEvent.mouseDown(select);
    const oleoOption = await screen.findByRole("option", {
      name: /^Óleo de Motor$/i,
    });
    fireEvent.click(oleoOption);
    fireEvent.keyDown(oleoOption, { key: "Escape" });

    // Title should be auto-suggested
    const titleInput = screen.getByLabelText(
      /Título do Serviço/i,
    ) as HTMLInputElement;
    expect(titleInput.value).toBe("Óleo de Motor");

    // Set cost
    const costInput = screen.getByLabelText(/Custo Total/i);
    fireEvent.change(costInput, { target: { value: "350.50" } });

    // Set mileage
    const mileageInput = screen.getByLabelText(
      /Quilometragem no Momento do Serviço/i,
    );
    fireEvent.change(mileageInput, { target: { value: "52000" } });

    // Submit
    const submitBtn = screen.getByRole("button", {
      name: /^Salvar Manutenção$/i,
    });
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(maintenanceService.create).toHaveBeenCalledWith(
        "car-1",
        expect.objectContaining({
          title: "Óleo de Motor",
          mileage: 52000,
          cost: 350.5,
          types: ["Óleo de Motor"],
        }),
      );
      expect(mockOnMaintenanceCreated).toHaveBeenCalled();
      expect(mockOnClose).toHaveBeenCalled();
    });
  });

  it("shows error if PDF file exceeds 2MB limit", async () => {
    render(
      <AddMaintenanceDialog
        open={true}
        carId="car-1"
        onClose={mockOnClose}
        onMaintenanceCreated={mockOnMaintenanceCreated}
      />,
    );

    const fileInput = document.querySelector(
      'input[type="file"]',
    ) as HTMLInputElement;
    expect(fileInput).toBeInTheDocument();

    // 2.5 MB PDF file
    const largePdf = new File(["dummy".repeat(600000)], "recibo_grande.pdf", {
      type: "application/pdf",
    });
    Object.defineProperty(largePdf, "size", { value: 2.5 * 1024 * 1024 });

    fireEvent.change(fileInput, { target: { files: [largePdf] } });

    await waitFor(() => {
      expect(
        screen.getByText(/excede o limite máximo de 2MB/i),
      ).toBeInTheDocument();
    });
  });

  it("shows error if unsupported file type is attached", async () => {
    render(
      <AddMaintenanceDialog
        open={true}
        carId="car-1"
        onClose={mockOnClose}
        onMaintenanceCreated={mockOnMaintenanceCreated}
      />,
    );

    const fileInput = document.querySelector(
      'input[type="file"]',
    ) as HTMLInputElement;

    const invalidFile = new File(["content"], "arquivo.exe", {
      type: "application/x-msdownload",
    });

    fireEvent.change(fileInput, { target: { files: [invalidFile] } });

    await waitFor(() => {
      expect(
        screen.getByText(/não suportado. Utilize PDF ou imagens/i),
      ).toBeInTheDocument();
    });
  });
});
