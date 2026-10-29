package memory

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	reportport "cvmc/internal/application/ports/report"
	domain "cvmc/internal/domain/report"
)

var ErrNotFound = errors.New("report job not found")

type Repository struct {
	mu   sync.RWMutex
	jobs map[string]*domain.ReportJob
	seq  int64
}

func NewRepository() *Repository {
	return &Repository{
		jobs: make(map[string]*domain.ReportJob),
	}
}

func (r *Repository) makeID() string {
	return fmt.Sprintf("report-%d", r.seq)
}

func (r *Repository) Create(ctx context.Context, job *domain.ReportJob) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	r.seq++
	if job.ID == "" {
		job.ID = r.makeID()
	}
	if job.CreatedAt.IsZero() {
		job.CreatedAt = time.Now().UTC()
	}
	copyJob := *job
	r.jobs[job.ID] = &copyJob
	return nil
}

func (r *Repository) Update(ctx context.Context, job *domain.ReportJob) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.jobs[job.ID]
	if !ok || existing.UserID != job.UserID {
		return ErrNotFound
	}
	copyJob := *job
	r.jobs[job.ID] = &copyJob
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id, userID string) (*domain.ReportJob, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()

	job, ok := r.jobs[id]
	if !ok || (userID != "" && job.UserID != userID) {
		return nil, ErrNotFound
	}
	copyJob := *job
	return &copyJob, nil
}

func (r *Repository) ListByUser(ctx context.Context, filter reportport.ListFilter) (*reportport.ListResult, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()

	var matching []*domain.ReportJob
	for _, job := range r.jobs {
		if job.UserID != filter.UserID {
			continue
		}
		if filter.CarID != "" && job.CarID != filter.CarID {
			continue
		}
		copyJob := *job
		matching = append(matching, &copyJob)
	}

	sort.Slice(matching, func(i, j int) bool {
		return matching[i].CreatedAt.After(matching[j].CreatedAt)
	})

	total := int64(len(matching))
	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 10
	}

	start := (page - 1) * limit
	if start >= len(matching) {
		return &reportport.ListResult{
			Jobs:  []*domain.ReportJob{},
			Total: total,
			Page:  page,
			Limit: limit,
		}, nil
	}

	end := start + limit
	if end > len(matching) {
		end = len(matching)
	}

	return &reportport.ListResult{
		Jobs:  matching[start:end],
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (r *Repository) FindPending(ctx context.Context, limit int) ([]*domain.ReportJob, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()

	var pending []*domain.ReportJob
	for _, job := range r.jobs {
		if job.Status == domain.StatusPending || job.Status == domain.StatusProcessing {
			copyJob := *job
			pending = append(pending, &copyJob)
		}
	}

	sort.Slice(pending, func(i, j int) bool {
		return pending[i].CreatedAt.Before(pending[j].CreatedAt)
	})

	if limit > 0 && len(pending) > limit {
		pending = pending[:limit]
	}

	return pending, nil
}

func (r *Repository) FindExpiredCompleted(ctx context.Context, cutoff time.Time, batchSize int) ([]*domain.ReportJob, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()

	var expired []*domain.ReportJob
	for _, job := range r.jobs {
		if job.Status == domain.StatusCompleted && job.ExpiresAt != nil && job.ExpiresAt.Before(cutoff) {
			copyJob := *job
			expired = append(expired, &copyJob)
			if batchSize > 0 && len(expired) >= batchSize {
				break
			}
		}
	}

	return expired, nil
}

func (r *Repository) MarkAsExpired(ctx context.Context, id string, expiredAt time.Time) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	job, ok := r.jobs[id]
	if !ok {
		return ErrNotFound
	}
	job.Status = domain.StatusExpired
	job.ExpiredAt = &expiredAt
	return nil
}
