package user

import (
	"context"

	domainuser "cvmc/internal/domain/user"
)

type Repository interface {
	Create(ctx context.Context, user domainuser.User) (domainuser.User, error)
	FindByEmail(ctx context.Context, email string) (domainuser.User, error)
	FindByID(ctx context.Context, id string) (domainuser.User, error)
}
