package mongo

import (
	"context"
	"errors"
	"strings"
	"time"

	notificationport "cvmc/internal/application/ports/notification"
	domainnotification "cvmc/internal/domain/notification"
	mongoinfra "cvmc/internal/infrastructure/database/mongo"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type deviceTokenDoc struct {
	ID        string    `bson:"_id"`
	UserID    string    `bson:"userId"`
	Token     string    `bson:"token"`
	Platform  string    `bson:"platform"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

type Repository struct {
	coll *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		coll: db.Collection("device_tokens"),
	}
}

func (r *Repository) EnsureIndexes(ctx context.Context) error {
	models := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "token", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "userId", Value: 1},
				{Key: "updatedAt", Value: -1},
			},
		},
	}
	_, err := r.coll.Indexes().CreateMany(ctx, models)
	return err
}

func (r *Repository) Upsert(ctx context.Context, dt domainnotification.DeviceToken) (domainnotification.DeviceToken, error) {
	cleanUserID, err := mongoinfra.SanitizeID(dt.UserID)
	if err != nil {
		return domainnotification.DeviceToken{}, err
	}
	cleanToken := strings.TrimSpace(dt.Token)
	if cleanToken == "" || len(cleanToken) > domainnotification.MaxTokenLength {
		return domainnotification.DeviceToken{}, errors.New("invalid device token")
	}

	platform, ok := domainnotification.NormalizePlatform(dt.Platform)
	if !ok {
		return domainnotification.DeviceToken{}, errors.New("invalid platform")
	}

	now := dt.UpdatedAt
	if now.IsZero() {
		now = time.Now().UTC()
	}

	newID := bson.NewObjectID().Hex()
	filter := bson.D{{Key: "token", Value: cleanToken}}
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "userId", Value: cleanUserID},
				{Key: "token", Value: cleanToken},
				{Key: "platform", Value: platform},
				{Key: "updatedAt", Value: now},
			},
		},
		{
			Key: "$setOnInsert",
			Value: bson.D{
				{Key: "_id", Value: newID},
			},
		},
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	var doc deviceTokenDoc
	if err := r.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&doc); err != nil {
		return domainnotification.DeviceToken{}, err
	}
	return toDomain(doc), nil
}

func (r *Repository) ListByUser(ctx context.Context, userID string) ([]domainnotification.DeviceToken, error) {
	cleanUserID, err := mongoinfra.SanitizeID(userID)
	if err != nil {
		return []domainnotification.DeviceToken{}, nil
	}

	filter := bson.D{{Key: "userId", Value: cleanUserID}}
	opts := options.Find().SetSort(bson.D{{Key: "updatedAt", Value: -1}})
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []deviceTokenDoc
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	items := make([]domainnotification.DeviceToken, 0, len(docs))
	for _, d := range docs {
		items = append(items, toDomain(d))
	}
	return items, nil
}

func (r *Repository) ListDistinctUserIDs(ctx context.Context) ([]string, error) {
	var rawIDs []string
	if err := r.coll.Distinct(ctx, "userId", bson.D{}).Decode(&rawIDs); err != nil {
		return nil, err
	}
	userIDs := make([]string, 0, len(rawIDs))
	for _, id := range rawIDs {
		if cleanID, err := mongoinfra.SanitizeID(id); err == nil && cleanID != "" {
			userIDs = append(userIDs, cleanID)
		}
	}
	return userIDs, nil
}

func (r *Repository) DeleteByUserAndToken(ctx context.Context, userID, token string) error {
	cleanUserID, err := mongoinfra.SanitizeID(userID)
	if err != nil {
		return notificationport.ErrNotFound
	}
	cleanToken := strings.TrimSpace(token)
	if cleanToken == "" {
		return notificationport.ErrNotFound
	}

	filter := bson.D{
		{Key: "userId", Value: cleanUserID},
		{Key: "token", Value: cleanToken},
	}
	res, err := r.coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return notificationport.ErrNotFound
	}
	return nil
}

func (r *Repository) DeleteByToken(ctx context.Context, token string) error {
	cleanToken := strings.TrimSpace(token)
	if cleanToken == "" {
		return notificationport.ErrNotFound
	}

	filter := bson.D{{Key: "token", Value: cleanToken}}
	res, err := r.coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return notificationport.ErrNotFound
	}
	return nil
}

func toDomain(d deviceTokenDoc) domainnotification.DeviceToken {
	return domainnotification.DeviceToken{
		ID:        d.ID,
		UserID:    d.UserID,
		Token:     d.Token,
		Platform:  d.Platform,
		UpdatedAt: d.UpdatedAt,
	}
}
