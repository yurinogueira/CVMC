package user

import (
	"context"
	"errors"

	domainuser "cvmc/internal/domain/user"
)

var ErrNotFound = errors.New("user not found")

type Repository interface {
	Create(ctx context.Context, user domainuser.User) (domainuser.User, error)
	FindByEmail(ctx context.Context, email string) (domainuser.User, error)
	FindByID(ctx context.Context, id string) (domainuser.User, error)
}

