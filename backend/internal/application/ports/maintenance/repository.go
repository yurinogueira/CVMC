package maintenance

import (
	"context"
	"errors"

	domainmaintenance "cvmc/internal/domain/maintenance"
)

var ErrNotFound = errors.New("maintenance not found")

type Repository interface {
	Create(ctx context.Context, maintenance domainmaintenance.Maintenance) (domainmaintenance.Maintenance, error)
	GetByID(ctx context.Context, id string) (domainmaintenance.Maintenance, error)
	ListByCar(ctx context.Context, carID string) ([]domainmaintenance.Maintenance, error)
	Update(ctx context.Context, maintenance domainmaintenance.Maintenance) (domainmaintenance.Maintenance, error)
	Delete(ctx context.Context, id string) error
}

