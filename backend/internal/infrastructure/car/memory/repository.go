package memory

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	domaincar "cvmc/internal/domain/car"
)

var ErrNotFound = errors.New("car not found")

type Repository struct {
	mu   sync.RWMutex
	byID map[string]domaincar.Car
	seq  int64
}

func NewRepository() *Repository {
	return &Repository{byID: make(map[string]domaincar.Car)}
}

func (r *Repository) Create(ctx context.Context, car domaincar.Car) (domaincar.Car, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	car.ID = r.makeID()
	car.SharedWith = append([]string{}, car.SharedWith...)
	r.byID[car.ID] = car
	return car, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (domaincar.Car, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	car, ok := r.byID[id]
	if !ok {
		return domaincar.Car{}, ErrNotFound
	}
	return car, nil
}

func (r *Repository) ListByUser(ctx context.Context, userID string) ([]domaincar.Car, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]domaincar.Car, 0)
	for _, car := range r.byID {
		if car.OwnerID == userID || contains(car.SharedWith, userID) {
			items = append(items, car)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })
	return items, nil
}

func (r *Repository) Update(ctx context.Context, car domaincar.Car) (domaincar.Car, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byID[car.ID]; !ok {
		return domaincar.Car{}, ErrNotFound
	}
	r.byID[car.ID] = car
	return car, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byID[id]; !ok {
		return ErrNotFound
	}
	delete(r.byID, id)
	return nil
}

func (r *Repository) Share(ctx context.Context, carID string, userID string) (domaincar.Car, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	car, ok := r.byID[carID]
	if !ok {
		return domaincar.Car{}, ErrNotFound
	}
	if !contains(car.SharedWith, userID) {
		car.SharedWith = append(car.SharedWith, userID)
	}
	r.byID[carID] = car
	return car, nil
}

func (r *Repository) Unshare(ctx context.Context, carID string, userID string) (domaincar.Car, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	car, ok := r.byID[carID]
	if !ok {
		return domaincar.Car{}, ErrNotFound
	}
	shared := make([]string, 0, len(car.SharedWith))
	for _, current := range car.SharedWith {
		if current != userID {
			shared = append(shared, current)
		}
	}
	car.SharedWith = shared
	r.byID[carID] = car
	return car, nil
}

func (r *Repository) makeID() string {
	return fmt.Sprintf("car-%d", r.seq)
}

func contains(items []string, value string) bool {
	for _, item := range items {
		if strings.EqualFold(item, value) {
			return true
		}
	}
	return false
}
