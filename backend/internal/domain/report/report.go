package report

import "time"

type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
	StatusExpired    Status = "expired"
)

type ReportType string

const (
	TypeMaintenancePDF ReportType = "maintenance_pdf"
)

type ReportJob struct {
	ID          string     `json:"id" bson:"_id,omitempty"`
	UserID      string     `json:"userId" bson:"userId"`
	CarID       string     `json:"carId" bson:"carId"`
	CarName     string     `json:"carName" bson:"carName"`
	Type        ReportType `json:"type" bson:"type"`
	Status      Status     `json:"status" bson:"status"`
	FilePath    string     `json:"filePath,omitempty" bson:"filePath,omitempty"`
	Error       string     `json:"error,omitempty" bson:"error,omitempty"`
	DurationMS  int64      `json:"durationMs,omitempty" bson:"durationMs,omitempty"`
	CreatedAt   time.Time  `json:"createdAt" bson:"createdAt"`
	CompletedAt *time.Time `json:"completedAt,omitempty" bson:"completedAt,omitempty"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty" bson:"expiresAt,omitempty"`
	ExpiredAt   *time.Time `json:"expiredAt,omitempty" bson:"expiredAt,omitempty"`
}
