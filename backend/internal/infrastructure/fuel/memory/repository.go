package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	fuelport "cvmc/internal/application/ports/fuel"
	domainfuel "cvmc/internal/domain/fuel"
)

type Repository struct {
	mu       sync.RWMutex
	fuelings map[string]domainfuel.Fueling
	seq      int64
}

func NewRepository() *Repository {
	return &Repository{
		fuelings: make(map[string]domainfuel.Fueling),
	}
}

func (r *Repository) Create(ctx context.Context, fueling domainfuel.Fueling) (domainfuel.Fueling, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if fueling.ID == "" {
		r.seq++
		fueling.ID = fmt.Sprintf("fuel-%d", r.seq)
	}
	if fueling.CreatedAt.IsZero() {
		fueling.CreatedAt = time.Now().UTC()
	}
	if fueling.UpdatedAt.IsZero() {
		fueling.UpdatedAt = time.Now().UTC()
	}
	r.fuelings[fueling.ID] = fueling
	return fueling, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (domainfuel.Fueling, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.fuelings[id]
	if !ok {
		return domainfuel.Fueling{}, fuelport.ErrNotFound
	}
	return item, nil
}

func (r *Repository) ListByCar(ctx context.Context, carID string) ([]domainfuel.Fueling, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domainfuel.Fueling
	for _, item := range r.fuelings {
		if item.CarID == carID {
			result = append(result, item)
		}
	}
	return result, nil
}

func (r *Repository) Update(ctx context.Context, fueling domainfuel.Fueling) (domainfuel.Fueling, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.fuelings[fueling.ID]; !ok {
		return domainfuel.Fueling{}, fuelport.ErrNotFound
	}
	fueling.UpdatedAt = time.Now().UTC()
	r.fuelings[fueling.ID] = fueling
	return fueling, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.fuelings[id]; !ok {
		return fuelport.ErrNotFound
	}
	delete(r.fuelings, id)
	return nil
}
