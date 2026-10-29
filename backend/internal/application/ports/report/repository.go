package report

import (
	"context"
	"time"

	domain "cvmc/internal/domain/report"
)

type ListFilter struct {
	UserID string
	CarID  string
	Page   int
	Limit  int
}

type ListResult struct {
	Jobs  []*domain.ReportJob
	Total int64
	Page  int
	Limit int
}

type Repository interface {
	Create(ctx context.Context, job *domain.ReportJob) error
	Update(ctx context.Context, job *domain.ReportJob) error
	GetByID(ctx context.Context, id, userID string) (*domain.ReportJob, error)
	ListByUser(ctx context.Context, filter ListFilter) (*ListResult, error)
	FindPending(ctx context.Context, limit int) ([]*domain.ReportJob, error)
	FindExpiredCompleted(ctx context.Context, cutoff time.Time, batchSize int) ([]*domain.ReportJob, error)
	MarkAsExpired(ctx context.Context, id string, expiredAt time.Time) error
}
