package maintenance

import (
	"context"

	domainmaintenance "cvmc/internal/domain/maintenance"
)

type Repository interface {
	Create(ctx context.Context, maintenance domainmaintenance.Maintenance) (domainmaintenance.Maintenance, error)
	GetByID(ctx context.Context, id string) (domainmaintenance.Maintenance, error)
	ListByCar(ctx context.Context, carID string) ([]domainmaintenance.Maintenance, error)
	Update(ctx context.Context, maintenance domainmaintenance.Maintenance) (domainmaintenance.Maintenance, error)
	Delete(ctx context.Context, id string) error
}
