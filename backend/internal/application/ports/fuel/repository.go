package fuel

import (
	"context"
	"errors"

	domainfuel "cvmc/internal/domain/fuel"
)

var ErrNotFound = errors.New("fueling not found")

type Repository interface {
	Create(ctx context.Context, fueling domainfuel.Fueling) (domainfuel.Fueling, error)
	GetByID(ctx context.Context, id string) (domainfuel.Fueling, error)
	ListByCar(ctx context.Context, carID string) ([]domainfuel.Fueling, error)
	Update(ctx context.Context, fueling domainfuel.Fueling) (domainfuel.Fueling, error)
	Delete(ctx context.Context, id string) error
}
