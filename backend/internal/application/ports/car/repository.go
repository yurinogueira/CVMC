package car

import (
	"context"
	"errors"

	domaincar "cvmc/internal/domain/car"
)

var ErrNotFound = errors.New("car not found")

type Repository interface {
	Create(ctx context.Context, car domaincar.Car) (domaincar.Car, error)
	GetByID(ctx context.Context, id string) (domaincar.Car, error)
	ListByUser(ctx context.Context, userID string) ([]domaincar.Car, error)
	Update(ctx context.Context, car domaincar.Car) (domaincar.Car, error)
	Delete(ctx context.Context, id string) error
	Share(ctx context.Context, carID string, userID string) (domaincar.Car, error)
	Unshare(ctx context.Context, carID string, userID string) (domaincar.Car, error)
}

