package memory

import (
	"context"
	"testing"
	"time"

	reportport "cvmc/internal/application/ports/report"
	domain "cvmc/internal/domain/report"
)

func TestMemoryReportRepository(t *testing.T) {
	repo := NewRepository()
	ctx := context.Background()

	now := time.Now().UTC()
	job := &domain.ReportJob{
		UserID:    "user-1",
		CarID:     "car-1",
		CarName:   "Civic",
		Type:      domain.TypeMaintenancePDF,
		Status:    domain.StatusPending,
		CreatedAt: now,
	}

	if err := repo.Create(ctx, job); err != nil {
		t.Fatalf("failed to create job: %v", err)
	}
	if job.ID == "" {
		t.Fatalf("expected job ID to be generated")
	}

	retrieved, err := repo.GetByID(ctx, job.ID, "user-1")
	if err != nil || retrieved == nil {
		t.Fatalf("failed to get job: %v", err)
	}
	if retrieved.CarName != "Civic" {
		t.Errorf("expected Civic, got %s", retrieved.CarName)
	}

	// Update
	completed := now.Add(time.Second)
	retrieved.Status = domain.StatusCompleted
	retrieved.CompletedAt = &completed
	retrieved.FilePath = "reports/user_user-1/job.pdf"
	if err := repo.Update(ctx, retrieved); err != nil {
		t.Fatalf("failed to update job: %v", err)
	}

	// List
	res, err := repo.ListByUser(ctx, reportport.ListFilter{UserID: "user-1", Page: 1, Limit: 10})
	if err != nil || res.Total != 1 || len(res.Jobs) != 1 {
		t.Fatalf("unexpected list result: %+v, err: %v", res, err)
	}

	// FindPending (none pending now)
	pending, err := repo.FindPending(ctx, 10)
	if err != nil || len(pending) != 0 {
		t.Fatalf("expected 0 pending, got %d", len(pending))
	}

	// Expired
	expiresAt := now.Add(-time.Hour)
	retrieved.ExpiresAt = &expiresAt
	_ = repo.Update(ctx, retrieved)

	expired, err := repo.FindExpiredCompleted(ctx, now, 10)
	if err != nil || len(expired) != 1 {
		t.Fatalf("expected 1 expired, got %d", len(expired))
	}

	if err := repo.MarkAsExpired(ctx, job.ID, now); err != nil {
		t.Fatalf("failed to mark expired: %v", err)
	}
	afterExpire, _ := repo.GetByID(ctx, job.ID, "user-1")
	if afterExpire.Status != domain.StatusExpired {
		t.Fatalf("expected expired status, got %s", afterExpire.Status)
	}
}
