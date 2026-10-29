import { apiClient } from "../../../services/api/client";
import { ApiEnvelope } from "../../auth/types/auth.types";
import { ReportJob, ListReportsResponse } from "../types/report.types";

export const reportService = {
  async startMaintenanceReport(carId: string): Promise<ReportJob> {
    const response = await apiClient.post<ApiEnvelope<ReportJob>>(
      "/reports/maintenance-pdf",
      { carId },
    );
    return response.data.data;
  },

  async listJobs(params?: {
    carId?: string;
    page?: number;
    limit?: number;
  }): Promise<ListReportsResponse> {
    const response = await apiClient.get<ApiEnvelope<ListReportsResponse>>(
      "/reports/jobs",
      { params },
    );
    return response.data.data || { jobs: [], total: 0, page: 1, limit: 10 };
  },

  async getJob(jobId: string): Promise<ReportJob> {
    const response = await apiClient.get<ApiEnvelope<ReportJob>>(
      `/reports/jobs/${jobId}`,
    );
    return response.data.data;
  },

  async downloadReport(filePath: string, defaultName?: string): Promise<void> {
    const response = await apiClient.get("/reports/download", {
      params: { file: filePath },
      responseType: "blob",
    });

    const blob = new Blob([response.data], { type: "application/pdf" });
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.setAttribute("download", defaultName || "relatorio_manutencoes.pdf");
    document.body.appendChild(link);
    link.click();
    link.remove();
    window.URL.revokeObjectURL(url);
  },
};
