import { describe, expect, it, vi, beforeEach } from "vitest";
import { render, screen, waitFor, fireEvent } from "@testing-library/react";
import {
  ExportHistoryTable,
  formatDuration,
  formatDate,
  formatExpiration,
} from "./ExportHistoryTable";
import { reportService } from "../services/report.service";
import { ReportJob } from "../types/report.types";

vi.mock("../services/report.service", () => ({
  reportService: {
    listJobs: vi.fn(),
    downloadReport: vi.fn(),
  },
}));

describe("ExportHistoryTable helper functions", () => {
  it("formats duration correctly", () => {
    expect(formatDuration(undefined)).toBe("-");
    expect(formatDuration(-1)).toBe("-");
    expect(formatDuration(500)).toBe("500ms");
    expect(formatDuration(1500)).toBe("1.5s");
    expect(formatDuration(65000)).toBe("1m 5s");
    expect(formatDuration(120000)).toBe("2m");
  });

  it("formats dates gracefully", () => {
    expect(formatDate(undefined)).toBe("-");
    expect(formatDate("invalid-date")).toBe("invalid-date");
    const d = formatDate("2026-05-10T12:00:00Z");
    expect(d).toContain("2026");
  });

  it("formats expiration status", () => {
    expect(formatExpiration(undefined, "expired")).toBe("Expirado");
    expect(formatExpiration(undefined, "completed")).toBe("-");

    const futureDate = new Date();
    futureDate.setDate(futureDate.getDate() + 15);
    expect(formatExpiration(futureDate.toISOString(), "completed")).toContain(
      "Expira em",
    );
  });
});

describe("ExportHistoryTable component", () => {
  const sampleJobs: ReportJob[] = [
    {
      id: "job-1",
      userId: "user-1",
      carId: "car-1",
      carName: "Honda Civic",
      type: "maintenance_pdf",
      status: "completed",
      filePath: "reports/user_user-1/job-1.pdf",
      durationMs: 1200,
      createdAt: "2026-06-01T10:00:00Z",
      expiresAt: "2026-07-01T10:00:00Z",
    },
    {
      id: "job-2",
      userId: "user-1",
      carId: "car-1",
      carName: "Honda Civic",
      type: "maintenance_pdf",
      status: "failed",
      error: "Timeout na geração",
      createdAt: "2026-06-02T10:00:00Z",
    },
  ];

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders table with jobs and allows download", async () => {
    vi.mocked(reportService.listJobs).mockResolvedValue({
      jobs: sampleJobs,
      total: 2,
      page: 1,
      limit: 10,
    });
    vi.mocked(reportService.downloadReport).mockResolvedValue();

    render(<ExportHistoryTable />);

    await waitFor(() => {
      expect(screen.getAllByText("Honda Civic").length).toBeGreaterThan(0);
      expect(screen.getByText("Concluído")).toBeInTheDocument();
      expect(screen.getByText("Falhou")).toBeInTheDocument();
    });

    const downloadButton = screen.getByText("Baixar PDF");
    fireEvent.click(downloadButton);

    await waitFor(() => {
      expect(reportService.downloadReport).toHaveBeenCalledWith(
        "reports/user_user-1/job-1.pdf",
        "dossie_manutencoes_Honda_Civic.pdf",
      );
    });
  });

  it("renders empty state when no jobs are returned", async () => {
    vi.mocked(reportService.listJobs).mockResolvedValue({
      jobs: [],
      total: 0,
      page: 1,
      limit: 10,
    });

    render(<ExportHistoryTable />);

    await waitFor(() => {
      expect(
        screen.getByText("Nenhum relatório exportado até o momento."),
      ).toBeInTheDocument();
    });
  });
});
