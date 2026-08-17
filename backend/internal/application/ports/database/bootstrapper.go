package database

import "context"

type Bootstrapper interface {
	Ensure(ctx context.Context) error
}
