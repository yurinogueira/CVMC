export type ReportStatus =
  "pending" | "processing" | "completed" | "failed" | "expired";

export type ReportType = "maintenance_pdf";

export interface ReportJob {
  id: string;
  userId: string;
  carId: string;
  carName: string;
  type: ReportType;
  status: ReportStatus;
  filePath?: string;
  error?: string;
  durationMs?: number;
  createdAt: string;
  completedAt?: string;
  expiresAt?: string;
  expiredAt?: string;
}

export interface ListReportsResponse {
  jobs: ReportJob[];
  total: number;
  page: number;
  limit: number;
}
