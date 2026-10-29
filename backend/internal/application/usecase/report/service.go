package report

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	carport "cvmc/internal/application/ports/car"
	maintenanceport "cvmc/internal/application/ports/maintenance"
	reportport "cvmc/internal/application/ports/report"
	storageport "cvmc/internal/application/ports/storage"
	domaincar "cvmc/internal/domain/car"
	domainnotification "cvmc/internal/domain/notification"
	domainreport "cvmc/internal/domain/report"
	pdfinfra "cvmc/internal/infrastructure/report/pdf"
	"cvmc/internal/shared/filecipher"
)

var (
	ErrCarNotFound        = errors.New("car not found")
	ErrForbidden          = errors.New("forbidden access to car or report")
	ErrInvalidReportPath  = errors.New("invalid report path")
	ErrReportExpired      = errors.New("report has expired")
	ErrStorageUnavailable = errors.New("storage provider unavailable")
)

type NotificationSender interface {
	SendToUser(ctx context.Context, userID string, msg domainnotification.PushMessage) (int, error)
}

type Service struct {
	reportRepo      reportport.Repository
	carRepo         carport.Repository
	maintRepo       maintenanceport.Repository
	storageProvider storageport.Provider
	notifSender     NotificationSender
	pdfBuilder      *pdfinfra.Builder
	queue           *WorkerQueue
	now             func() time.Time
}

type WorkerQueue struct {
	jobsChan chan string
	service  *Service
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	mu       sync.Mutex
	running  bool
}

func NewWorkerQueue(bufferSize int) *WorkerQueue {
	if bufferSize <= 0 {
		bufferSize = 100
	}
	return &WorkerQueue{
		jobsChan: make(chan string, bufferSize),
	}
}

func (q *WorkerQueue) Start(parentCtx context.Context) {
	q.mu.Lock()
	if q.running {
		q.mu.Unlock()
		return
	}
	ctx, cancel := context.WithCancel(parentCtx)
	q.cancel = cancel
	q.running = true
	q.wg.Add(1)
	q.mu.Unlock()

	go func() {
		defer q.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case jobID, ok := <-q.jobsChan:
				if !ok {
					return
				}
				if q.service != nil {
					q.service.ProcessJob(context.Background(), jobID)
					runtime.GC() // Clean up heavy buffers after PDF processing
				}
			}
		}
	}()
}

func (q *WorkerQueue) Stop() {
	q.mu.Lock()
	if !q.running {
		q.mu.Unlock()
		return
	}
	cancel := q.cancel
	q.running = false
	q.cancel = nil
	q.mu.Unlock()

	if cancel != nil {
		cancel()
	}
	q.wg.Wait()
}

func (q *WorkerQueue) Enqueue(jobID string) {
	select {
	case q.jobsChan <- jobID:
	default:
		// Fallback: spawn one-time processor if channel is saturated
		go func(id string) {
			if q.service != nil {
				q.service.ProcessJob(context.Background(), id)
				runtime.GC()
			}
		}(jobID)
	}
}

func NewService(
	reportRepo reportport.Repository,
	carRepo carport.Repository,
	maintRepo maintenanceport.Repository,
	storageProvider storageport.Provider,
	notifSender NotificationSender,
	cipher *filecipher.FileCipher,
) *Service {
	queue := NewWorkerQueue(100)
	svc := &Service{
		reportRepo:      reportRepo,
		carRepo:         carRepo,
		maintRepo:       maintRepo,
		storageProvider: storageProvider,
		notifSender:     notifSender,
		pdfBuilder:      pdfinfra.NewBuilder(cipher),
		queue:           queue,
		now:             time.Now,
	}
	queue.service = svc
	return svc
}

func (s *Service) StartWorker(ctx context.Context) {
	if s.queue != nil {
		s.queue.Start(ctx)
		// Reconcile and recover pending jobs on startup
		if s.reportRepo != nil {
			pending, err := s.reportRepo.FindPending(ctx, 100)
			if err == nil {
				for _, j := range pending {
					s.queue.Enqueue(j.ID)
				}
			}
		}
	}
}

func (s *Service) StopWorker() {
	if s.queue != nil {
		s.queue.Stop()
	}
}

func (s *Service) StartMaintenanceReportJob(ctx context.Context, userID, carID string) (*domainreport.ReportJob, error) {
	cleanUserID := strings.TrimSpace(userID)
	cleanCarID := strings.TrimSpace(carID)
	if cleanUserID == "" || cleanCarID == "" {
		return nil, errors.New("userID and carID are required")
	}

	car, err := s.carRepo.GetByID(ctx, cleanCarID)
	if err != nil {
		return nil, ErrCarNotFound
	}

	// Verify car access
	if !canAccessCar(cleanUserID, car) {
		return nil, ErrForbidden
	}

	carName := car.Name
	if carName == "" {
		carName = fmt.Sprintf("%s %s", car.Manufacturer, car.Model)
	}

	job := &domainreport.ReportJob{
		UserID:    cleanUserID,
		CarID:     cleanCarID,
		CarName:   carName,
		Type:      domainreport.TypeMaintenancePDF,
		Status:    domainreport.StatusPending,
		CreatedAt: s.now().UTC(),
	}

	if err := s.reportRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create report job: %w", err)
	}

	// Enqueue in sequential worker
	if s.queue != nil {
		s.queue.Enqueue(job.ID)
	}

	return job, nil
}

func (s *Service) ProcessJob(ctx context.Context, jobID string) {
	jobCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	// Locate job across any user for processing
	pendingList, err := s.reportRepo.FindPending(jobCtx, 100)
	if err != nil {
		return
	}
	var targetJob *domainreport.ReportJob
	for _, j := range pendingList {
		if j.ID == jobID {
			targetJob = j
			break
		}
	}
	if targetJob == nil {
		return
	}

	start := s.now().UTC()
	targetJob.Status = domainreport.StatusProcessing
	_ = s.reportRepo.Update(jobCtx, targetJob)

	car, err := s.carRepo.GetByID(jobCtx, targetJob.CarID)
	if err != nil {
		s.failJob(jobCtx, targetJob, start, fmt.Sprintf("car not found: %v", err))
		return
	}

	maintenances, err := s.maintRepo.ListByCar(jobCtx, targetJob.CarID)
	if err != nil {
		s.failJob(jobCtx, targetJob, start, fmt.Sprintf("failed to list maintenances: %v", err))
		return
	}

	// Sort chronologically (oldest to newest)
	sort.Slice(maintenances, func(i, j int) bool {
		return maintenances[i].Date.Before(maintenances[j].Date)
	})

	pdfBytes, err := s.pdfBuilder.BuildMaintenanceReport(jobCtx, car, maintenances)
	if err != nil {
		s.failJob(jobCtx, targetJob, start, fmt.Sprintf("pdf generation failed: %v", err))
		return
	}

	if s.storageProvider == nil {
		s.failJob(jobCtx, targetJob, start, "storage provider not available")
		return
	}

	fileName := fmt.Sprintf("manutencoes_%s_%d.pdf", targetJob.CarID, start.Unix())
	relativePath := fmt.Sprintf("reports/user_%s/%s", targetJob.UserID, fileName)

	_, err = s.storageProvider.Save(jobCtx, relativePath, storageport.File{
		Name:        fileName,
		Data:        pdfBytes,
		ContentType: "application/pdf",
	})
	if err != nil {
		s.failJob(jobCtx, targetJob, start, fmt.Sprintf("failed to store report file: %v", err))
		return
	}

	completed := s.now().UTC()
	expires := completed.Add(30 * 24 * time.Hour)

	targetJob.Status = domainreport.StatusCompleted
	targetJob.FilePath = relativePath
	targetJob.CompletedAt = &completed
	targetJob.ExpiresAt = &expires
	targetJob.DurationMS = completed.Sub(start).Milliseconds()
	targetJob.Error = ""

	_ = s.reportRepo.Update(jobCtx, targetJob)

	// Send Push Notification
	if s.notifSender != nil {
		msg := domainnotification.PushMessage{
			Title: "Relatório de Manutenções Pronto 📄",
			Body:  fmt.Sprintf("O relatório do seu %s está pronto para download!", targetJob.CarName),
			Data: map[string]string{
				"type":  "report_ready",
				"jobId": targetJob.ID,
				"carId": targetJob.CarID,
				"route": "/reports",
			},
		}
		_, _ = s.notifSender.SendToUser(jobCtx, targetJob.UserID, msg)
	}
}

func (s *Service) failJob(ctx context.Context, job *domainreport.ReportJob, start time.Time, reason string) {
	completed := s.now().UTC()
	job.Status = domainreport.StatusFailed
	job.Error = reason
	job.CompletedAt = &completed
	job.DurationMS = completed.Sub(start).Milliseconds()
	_ = s.reportRepo.Update(ctx, job)
	log.Printf("[REPORT-JOB-FAILED] Job %s failed: %s", job.ID, reason)
}

func (s *Service) GetJob(ctx context.Context, id, userID string) (*domainreport.ReportJob, error) {
	return s.reportRepo.GetByID(ctx, id, userID)
}

func (s *Service) ListJobs(ctx context.Context, filter reportport.ListFilter) (*reportport.ListResult, error) {
	return s.reportRepo.ListByUser(ctx, filter)
}

func (s *Service) GetReportFile(ctx context.Context, userID, relativePath string) ([]byte, error) {
	cleanPath := filepath.ToSlash(filepath.Clean(strings.TrimPrefix(relativePath, "/")))
	if strings.Contains(cleanPath, "..") {
		return nil, ErrInvalidReportPath
	}

	expectedPrefix := fmt.Sprintf("reports/user_%s/", userID)
	if !strings.HasPrefix(cleanPath, expectedPrefix) {
		return nil, ErrForbidden
	}

	// Verify expiration from repository
	listRes, err := s.reportRepo.ListByUser(ctx, reportport.ListFilter{UserID: userID, Limit: 100})
	if err == nil && listRes != nil {
		for _, j := range listRes.Jobs {
			if j.FilePath == cleanPath {
				if j.Status == domainreport.StatusExpired {
					return nil, ErrReportExpired
				}
				if j.ExpiresAt != nil && j.ExpiresAt.Before(s.now().UTC()) {
					return nil, ErrReportExpired
				}
				break
			}
		}
	}

	if s.storageProvider == nil {
		return nil, ErrStorageUnavailable
	}

	return s.storageProvider.Get(ctx, cleanPath)
}

func (s *Service) CleanupExpiredReports(ctx context.Context, retentionDuration time.Duration, batchSize int) (int, error) {
	if s.reportRepo == nil {
		return 0, nil
	}
	if retentionDuration <= 0 {
		retentionDuration = 30 * 24 * time.Hour
	}
	if batchSize <= 0 {
		batchSize = 50
	}

	cutoff := s.now().UTC()
	totalPurged := 0

	for {
		select {
		case <-ctx.Done():
			return totalPurged, ctx.Err()
		default:
		}

		jobs, err := s.reportRepo.FindExpiredCompleted(ctx, cutoff, batchSize)
		if err != nil {
			return totalPurged, fmt.Errorf("failed to find expired report jobs: %w", err)
		}
		if len(jobs) == 0 {
			break
		}

		now := s.now().UTC()
		batchProcessed := 0

		for _, job := range jobs {
			select {
			case <-ctx.Done():
				return totalPurged, ctx.Err()
			default:
			}

			if job.FilePath != "" && s.storageProvider != nil {
				cleanPath := filepath.ToSlash(filepath.Clean(strings.TrimPrefix(job.FilePath, "/")))
				expectedPrefix := fmt.Sprintf("reports/user_%s/", job.UserID)
				if strings.HasPrefix(cleanPath, expectedPrefix) && !strings.Contains(cleanPath, "..") {
					if delErr := s.storageProvider.Delete(ctx, cleanPath); delErr != nil {
						isAbsent := errors.Is(delErr, os.ErrNotExist) ||
							strings.Contains(strings.ToLower(delErr.Error()), "not found") ||
							strings.Contains(strings.ToLower(delErr.Error()), "does not exist")
						if !isAbsent {
							log.Printf("[REPORT-CLEANUP] Failed to delete file %s: %v", cleanPath, delErr)
							continue
						}
					}
				}
			}

			if err := s.reportRepo.MarkAsExpired(ctx, job.ID, now); err != nil {
				log.Printf("[REPORT-CLEANUP] Failed to mark job %s as expired: %v", job.ID, err)
				continue
			}

			totalPurged++
			batchProcessed++
		}

		if len(jobs) < batchSize || batchProcessed == 0 {
			break
		}
	}

	return totalPurged, nil
}

func canAccessCar(userID string, car domaincar.Car) bool {
	if car.OwnerID == userID {
		return true
	}
	for _, id := range car.SharedWith {
		if id == userID {
			return true
		}
	}
	return false
}
