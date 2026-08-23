package mongo

import (
	"context"
	"errors"
	"strings"
	"time"

	userport "cvmc/internal/application/ports/user"
	domainuser "cvmc/internal/domain/user"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type userDoc struct {
	ID           string    `bson:"_id"`
	Name         string    `bson:"name"`
	Email        string    `bson:"email"`
	PasswordHash string    `bson:"passwordHash"`
	CreatedAt    time.Time `bson:"createdAt"`
}

type Repository struct {
	coll *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		coll: db.Collection("users"),
	}
}

func (r *Repository) EnsureIndexes(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := r.coll.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (r *Repository) Create(ctx context.Context, user domainuser.User) (domainuser.User, error) {
	if user.ID == "" {
		user.ID = bson.NewObjectID().Hex()
	}
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now().UTC()
	}

	doc := userDoc{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
	}

	_, err := r.coll.InsertOne(ctx, doc)
	if err != nil {
		return domainuser.User{}, err
	}
	return user, nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (domainuser.User, error) {
	normalized := strings.ToLower(strings.TrimSpace(email))
	var doc userDoc
	err := r.coll.FindOne(ctx, bson.M{"email": normalized}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domainuser.User{}, userport.ErrNotFound
		}
		return domainuser.User{}, err
	}
	return domainuser.User{
		ID:           doc.ID,
		Name:         doc.Name,
		Email:        doc.Email,
		PasswordHash: doc.PasswordHash,
		CreatedAt:    doc.CreatedAt,
	}, nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (domainuser.User, error) {
	var doc userDoc
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domainuser.User{}, userport.ErrNotFound
		}
		return domainuser.User{}, err
	}
	return domainuser.User{
		ID:           doc.ID,
		Name:         doc.Name,
		Email:        doc.Email,
		PasswordHash: doc.PasswordHash,
		CreatedAt:    doc.CreatedAt,
	}, nil
}
