package memory

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	domainmaintenance "cvmc/internal/domain/maintenance"
)

var ErrNotFound = errors.New("maintenance not found")

type Repository struct {
	mu    sync.RWMutex
	byID  map[string]domainmaintenance.Maintenance
	byCar map[string][]string
	seq   int64
}

func NewRepository() *Repository {
	return &Repository{byID: make(map[string]domainmaintenance.Maintenance), byCar: make(map[string][]string)}
}

func (r *Repository) Create(ctx context.Context, maintenance domainmaintenance.Maintenance) (domainmaintenance.Maintenance, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	maintenance.ID = r.makeID()
	if maintenance.CreatedAt.IsZero() {
		maintenance.CreatedAt = time.Now().UTC()
	}
	if maintenance.UpdatedAt.IsZero() {
		maintenance.UpdatedAt = maintenance.CreatedAt
	}
	r.byID[maintenance.ID] = maintenance
	r.byCar[maintenance.CarID] = append(r.byCar[maintenance.CarID], maintenance.ID)
	return maintenance, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (domainmaintenance.Maintenance, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	maintenance, ok := r.byID[id]
	if !ok {
		return domainmaintenance.Maintenance{}, ErrNotFound
	}
	return maintenance, nil
}

func (r *Repository) ListByCar(ctx context.Context, carID string) ([]domainmaintenance.Maintenance, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := r.byCar[carID]
	items := make([]domainmaintenance.Maintenance, 0, len(ids))
	for _, id := range ids {
		if maintenance, ok := r.byID[id]; ok {
			items = append(items, maintenance)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Date.After(items[j].Date) })
	return items, nil
}

func (r *Repository) Update(ctx context.Context, maintenance domainmaintenance.Maintenance) (domainmaintenance.Maintenance, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byID[maintenance.ID]; !ok {
		return domainmaintenance.Maintenance{}, ErrNotFound
	}
	r.byID[maintenance.ID] = maintenance
	return maintenance, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	maintenance, ok := r.byID[id]
	if !ok {
		return ErrNotFound
	}
	delete(r.byID, id)
	ids := r.byCar[maintenance.CarID]
	filtered := ids[:0]
	for _, current := range ids {
		if current != id {
			filtered = append(filtered, current)
		}
	}
	r.byCar[maintenance.CarID] = filtered
	return nil
}

func (r *Repository) makeID() string {
	return fmt.Sprintf("maintenance-%d", r.seq)
}
