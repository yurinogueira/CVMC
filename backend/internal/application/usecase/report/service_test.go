package report

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	domaincar "cvmc/internal/domain/car"
	domainmaint "cvmc/internal/domain/maintenance"
	domainnotification "cvmc/internal/domain/notification"
	domainreport "cvmc/internal/domain/report"
	memorycar "cvmc/internal/infrastructure/car/memory"
	memorymaint "cvmc/internal/infrastructure/maintenance/memory"
	memoryreport "cvmc/internal/infrastructure/report/memory"
	memorylocal "cvmc/internal/infrastructure/storage/local"
	"cvmc/internal/shared/filecipher"
)

type mockNotifSender struct {
	mu       sync.Mutex
	messages []domainnotification.PushMessage
	userIDs  []string
}

func (m *mockNotifSender) SendToUser(ctx context.Context, userID string, msg domainnotification.PushMessage) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.userIDs = append(m.userIDs, userID)
	m.messages = append(m.messages, msg)
	return 1, nil
}

func setupTestService(t *testing.T) (*Service, *memorycar.Repository, *memorymaint.Repository, *mockNotifSender) {
	tmpDir := t.TempDir()
	storage := memorylocal.New(tmpDir)
	reportRepo := memoryreport.NewRepository()
	carRepo := memorycar.NewRepository()
	maintRepo := memorymaint.NewRepository()
	notifSender := &mockNotifSender{}
	cipher := filecipher.New("master-secret-key-12345678901234567890")

	svc := NewService(reportRepo, carRepo, maintRepo, storage, notifSender, cipher)
	return svc, carRepo, maintRepo, notifSender
}

func TestStartMaintenanceReportJob_ValidationsAndCreation(t *testing.T) {
	svc, carRepo, _, _ := setupTestService(t)
	ctx := context.Background()

	// 1. Missing IDs
	_, err := svc.StartMaintenanceReportJob(ctx, "", "car-1")
	if err == nil {
		t.Fatal("expected error for empty userID")
	}

	// 2. Car not found
	_, err = svc.StartMaintenanceReportJob(ctx, "user-1", "non-existent")
	if !errors.Is(err, ErrCarNotFound) {
		t.Fatalf("expected ErrCarNotFound, got %v", err)
	}

	// 3. Forbidden car (not owner or shared)
	createdCar, _ := carRepo.Create(ctx, domaincar.Car{
		OwnerID:      "other-user",
		Name:         "Golf",
		Manufacturer: "VW",
		Model:        "Golf",
	})
	_, err = svc.StartMaintenanceReportJob(ctx, "user-1", createdCar.ID)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}

	// 4. Success for owner
	ownedCar, _ := carRepo.Create(ctx, domaincar.Car{
		OwnerID:      "user-1",
		Name:         "Civic",
		Manufacturer: "Honda",
		Model:        "Civic",
	})

	job, err := svc.StartMaintenanceReportJob(ctx, "user-1", ownedCar.ID)
	if err != nil {
		t.Fatalf("unexpected error starting job: %v", err)
	}
	if job.Status != domainreport.StatusPending {
		t.Errorf("expected StatusPending, got %s", job.Status)
	}
	if job.CarName != "Civic" {
		t.Errorf("expected car name Civic, got %s", job.CarName)
	}
}

func TestProcessJob_EndToEndGeneration(t *testing.T) {
	svc, carRepo, maintRepo, notifSender := setupTestService(t)
	ctx := context.Background()

	car, _ := carRepo.Create(ctx, domaincar.Car{
		OwnerID:      "user-1",
		Name:         "Civic Touring",
		Manufacturer: "Honda",
		Model:        "Civic",
		LastMileage:  50000,
	})

	cost := 450.0
	_, _ = maintRepo.Create(ctx, domainmaint.Maintenance{
		CarID:       car.ID,
		Title:       "Troca de Pastilhas de Freio",
		Description: "Pastilhas dianteiras e traseiras",
		Date:        time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC),
		Mileage:     48000,
		Cost:        &cost,
	})

	job, err := svc.StartMaintenanceReportJob(ctx, "user-1", car.ID)
	if err != nil {
		t.Fatalf("failed to start job: %v", err)
	}

	// Process job synchronously
	svc.ProcessJob(ctx, job.ID)

	updatedJob, err := svc.GetJob(ctx, job.ID, "user-1")
	if err != nil {
		t.Fatalf("failed to get updated job: %v", err)
	}
	if updatedJob.Status != domainreport.StatusCompleted {
		t.Fatalf("expected completed status, got %s (error: %s)", updatedJob.Status, updatedJob.Error)
	}
	if updatedJob.FilePath == "" {
		t.Fatalf("expected file path to be populated")
	}
	if updatedJob.ExpiresAt == nil {
		t.Fatalf("expected expires_at to be populated")
	}

	// Verify push notification was triggered
	notifSender.mu.Lock()
	if len(notifSender.messages) != 1 {
		t.Fatalf("expected 1 push notification, got %d", len(notifSender.messages))
	}
	if notifSender.userIDs[0] != "user-1" {
		t.Errorf("expected push to user-1, got %s", notifSender.userIDs[0])
	}
	notifSender.mu.Unlock()

	// Download file
	fileBytes, err := svc.GetReportFile(ctx, "user-1", updatedJob.FilePath)
	if err != nil {
		t.Fatalf("failed to download report file: %v", err)
	}
	if len(fileBytes) < 500 {
		t.Fatalf("downloaded file too small: %d bytes", len(fileBytes))
	}
}

func TestGetReportFile_SecurityAndExpiration(t *testing.T) {
	svc, carRepo, _, _ := setupTestService(t)
	ctx := context.Background()

	car, _ := carRepo.Create(ctx, domaincar.Car{
		OwnerID: "user-1",
		Name:    "Civic",
	})
	job, _ := svc.StartMaintenanceReportJob(ctx, "user-1", car.ID)
	svc.ProcessJob(ctx, job.ID)
	completedJob, _ := svc.GetJob(ctx, job.ID, "user-1")

	// 1. Path traversal attempt
	_, err := svc.GetReportFile(ctx, "user-1", "../secret.pdf")
	if !errors.Is(err, ErrInvalidReportPath) {
		t.Fatalf("expected ErrInvalidReportPath, got %v", err)
	}

	// 2. Cross-user access attempt
	_, err = svc.GetReportFile(ctx, "other-user", completedJob.FilePath)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}

	// 3. Expired file check
	expiredTime := time.Now().UTC().Add(-2 * time.Hour)
	completedJob.ExpiresAt = &expiredTime
	completedJob.Status = domainreport.StatusExpired
	_ = svc.reportRepo.Update(ctx, completedJob)

	_, err = svc.GetReportFile(ctx, "user-1", completedJob.FilePath)
	if !errors.Is(err, ErrReportExpired) {
		t.Fatalf("expected ErrReportExpired, got %v", err)
	}
}

func TestCleanupExpiredReports(t *testing.T) {
	svc, carRepo, _, _ := setupTestService(t)
	ctx := context.Background()

	car, _ := carRepo.Create(ctx, domaincar.Car{
		OwnerID: "user-1",
		Name:    "Civic",
	})
	job, _ := svc.StartMaintenanceReportJob(ctx, "user-1", car.ID)
	svc.ProcessJob(ctx, job.ID)
	completedJob, _ := svc.GetJob(ctx, job.ID, "user-1")

	// Force job to have expired in the past
	pastExpiration := time.Now().UTC().Add(-48 * time.Hour)
	completedJob.ExpiresAt = &pastExpiration
	_ = svc.reportRepo.Update(ctx, completedJob)

	// Run cleanup
	purged, err := svc.CleanupExpiredReports(ctx, 30*24*time.Hour, 50)
	if err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}
	if purged != 1 {
		t.Fatalf("expected 1 purged report, got %d", purged)
	}

	// Verify status updated to expired
	afterPurge, _ := svc.GetJob(ctx, job.ID, "user-1")
	if afterPurge.Status != domainreport.StatusExpired {
		t.Fatalf("expected status expired, got %s", afterPurge.Status)
	}
}

func TestWorkerQueue_Lifecycle(t *testing.T) {
	svc, carRepo, _, _ := setupTestService(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	car, _ := carRepo.Create(ctx, domaincar.Car{
		OwnerID: "user-queue",
		Name:    "Polo",
	})

	svc.StartWorker(ctx)

	job, err := svc.StartMaintenanceReportJob(ctx, "user-queue", car.ID)
	if err != nil {
		t.Fatalf("failed to start job: %v", err)
	}

	// Wait for queue worker to pick up and process
	var finishedJob *domainreport.ReportJob
	for i := 0; i < 20; i++ {
		time.Sleep(50 * time.Millisecond)
		j, err := svc.GetJob(ctx, job.ID, "user-queue")
		if err == nil && (j.Status == domainreport.StatusCompleted || j.Status == domainreport.StatusFailed) {
			finishedJob = j
			break
		}
	}

	if finishedJob == nil || finishedJob.Status != domainreport.StatusCompleted {
		t.Fatalf("expected job to be processed by worker queue within timeout, got %+v", finishedJob)
	}

	svc.StopWorker()
}
