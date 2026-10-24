export interface MaintenanceAttachment {
  id: string;
  name: string;
  size: number;
  mimeType: string;
  dataUrl: string;
  createdAt: string;
}

export interface Maintenance {
  id: string;
  carId: string;
  title: string;
  description: string;
  date: string;
  mileage: number;
  types?: string[];
  cost?: number;
  attachments?: MaintenanceAttachment[];
  createdAt?: string;
  updatedAt?: string;
}

export interface CreateMaintenancePayload {
  title: string;
  description?: string;
  date: string;
  mileage: number;
  types?: string[];
  cost?: number;
  attachments?: MaintenanceAttachment[];
}

export interface UpdateMaintenancePayload {
  title: string;
  description?: string;
  date: string;
  mileage: number;
  types?: string[];
  cost?: number;
  attachments?: MaintenanceAttachment[];
}

export type UpcomingMaintenanceStatus = "ok" | "upcoming" | "overdue";

export interface UpcomingMaintenanceItem {
  type: string;
  label: string;
  status: UpcomingMaintenanceStatus;
  intervalKm: number;
  intervalMonths: number;
  hasHistory: boolean;
  lastPerformedAt?: string;
  lastMileage: number;
  nextDueMileage: number;
  nextDueDate: string;
  remainingKm: number;
  remainingDays: number;
  progressPercent: number;
  alertEligible: boolean;
}

export interface CarUpcomingMaintenances {
  carId: string;
  carName: string;
  manufacturer: string;
  model: string;
  vehicleType: string;
  currentKm: number;
  status: UpcomingMaintenanceStatus;
  okCount: number;
  upcomingCount: number;
  overdueCount: number;
  items: UpcomingMaintenanceItem[];
}

export interface FleetUpcomingSummary {
  totalVehicles: number;
  status: UpcomingMaintenanceStatus;
  okCount: number;
  upcomingCount: number;
  overdueCount: number;
  vehicles: CarUpcomingMaintenances[];
}
