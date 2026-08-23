package mongo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Indexable interface {
	EnsureIndexes(ctx context.Context) error
}

type Bootstrapper struct {
	db          *mongo.Database
	collections []string
	indexers    []Indexable
}

func Connect(ctx context.Context, uri string) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		_ = client.Disconnect(ctx)
		return nil, fmt.Errorf("failed to ping mongodb cluster: %w", err)
	}
	return client, nil
}

func NewBootstrapper(db *mongo.Database, collections []string, indexers ...Indexable) *Bootstrapper {
	return &Bootstrapper{
		db:          db,
		collections: collections,
		indexers:    indexers,
	}
}

func (b *Bootstrapper) Ensure(ctx context.Context) error {
	for _, name := range b.collections {
		err := b.db.CreateCollection(ctx, name)
		if err != nil && !isCollectionExistsError(err) {
			return fmt.Errorf("failed to create collection %q: %w", name, err)
		}
	}

	for _, indexer := range b.indexers {
		if err := indexer.EnsureIndexes(ctx); err != nil {
			return fmt.Errorf("failed to ensure indexes: %w", err)
		}
	}
	return nil
}

func isCollectionExistsError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "already exists") || strings.Contains(msg, "NamespaceExists") || strings.Contains(msg, "code 48")
}
