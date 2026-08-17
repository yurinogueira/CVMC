package mongo

import "context"

type Collection interface {
	Name() string
	EnsureIndexes(ctx context.Context) error
	Validate(ctx context.Context) error
	Seed(ctx context.Context) error
}

type Bootstrapper struct {
	collections []Collection
}

func New(collections ...Collection) *Bootstrapper {
	return &Bootstrapper{collections: collections}
}

func (b *Bootstrapper) Ensure(ctx context.Context) error {
	for _, collection := range b.collections {
		if err := collection.Validate(ctx); err != nil {
			return err
		}
		if err := collection.EnsureIndexes(ctx); err != nil {
			return err
		}
		if err := collection.Seed(ctx); err != nil {
			return err
		}
	}
	return nil
}
