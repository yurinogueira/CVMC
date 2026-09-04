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
