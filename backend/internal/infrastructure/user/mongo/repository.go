package mongo

import (
	"context"
	"errors"
	"strings"
	"time"

	userport "cvmc/internal/application/ports/user"
	domainuser "cvmc/internal/domain/user"
	mongoinfra "cvmc/internal/infrastructure/database/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type userDoc struct {
	ID                         string     `bson:"_id"`
	Name                       string     `bson:"name"`
	Email                      string     `bson:"email"`
	PasswordHash               string     `bson:"passwordHash"`
	GoogleID                   string     `bson:"googleId,omitempty"`
	EmailVerified              bool       `bson:"emailVerified"`
	EmailVerifiedAt            *time.Time `bson:"emailVerifiedAt,omitempty"`
	EmailVerificationTokenHash string     `bson:"emailVerificationTokenHash,omitempty"`
	EmailVerificationExpiresAt *time.Time `bson:"emailVerificationExpiresAt,omitempty"`
	PasswordResetTokenHash     string     `bson:"passwordResetTokenHash,omitempty"`
	PasswordResetExpiresAt     *time.Time `bson:"passwordResetExpiresAt,omitempty"`
	MaxVehicles                int        `bson:"maxVehicles"`
	MaxFiles                   int        `bson:"maxFiles,omitempty"`
	MaxStorageBytes            int64      `bson:"maxStorageBytes,omitempty"`
	FilesCount                 int        `bson:"filesCount,omitempty"`
	StorageUsedBytes           int64      `bson:"storageUsedBytes,omitempty"`
	StorageQuotaSynced         bool       `bson:"storageQuotaSynced,omitempty"`
	CreatedAt                  time.Time  `bson:"createdAt"`
	UpdatedAt                  time.Time  `bson:"updatedAt,omitempty"`
}

func (d userDoc) toDomain() domainuser.User {
	maxVehicles := d.MaxVehicles
	if maxVehicles <= 0 {
		maxVehicles = domainuser.DefaultMaxVehicles
	}
	maxFiles := d.MaxFiles
	if maxFiles <= 0 {
		maxFiles = domainuser.DefaultMaxFiles
	}
	maxStorage := d.MaxStorageBytes
	if maxStorage <= 0 {
		maxStorage = domainuser.DefaultMaxStorageBytes
	}
	return domainuser.User{
		ID:                         d.ID,
		Name:                       d.Name,
		Email:                      d.Email,
		PasswordHash:               d.PasswordHash,
		GoogleID:                   d.GoogleID,
		EmailVerified:              d.EmailVerified,
		EmailVerifiedAt:            d.EmailVerifiedAt,
		EmailVerificationTokenHash: d.EmailVerificationTokenHash,
		EmailVerificationExpiresAt: d.EmailVerificationExpiresAt,
		PasswordResetTokenHash:     d.PasswordResetTokenHash,
		PasswordResetExpiresAt:     d.PasswordResetExpiresAt,
		MaxVehicles:                maxVehicles,
		MaxFiles:                   maxFiles,
		MaxStorageBytes:            maxStorage,
		FilesCount:                 d.FilesCount,
		StorageUsedBytes:           d.StorageUsedBytes,
		StorageQuotaSynced:         d.StorageQuotaSynced,
		CreatedAt:                  d.CreatedAt,
		UpdatedAt:                  d.UpdatedAt,
	}
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
	models := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "googleId", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true),
		},
		{
			Keys:    bson.D{{Key: "emailVerificationTokenHash", Value: 1}},
			Options: options.Index().SetSparse(true),
		},
		{
			Keys:    bson.D{{Key: "passwordResetTokenHash", Value: 1}},
			Options: options.Index().SetSparse(true),
		},
	}
	_, err := r.coll.Indexes().CreateMany(ctx, models)
	return err
}

func (r *Repository) Create(ctx context.Context, user domainuser.User) (domainuser.User, error) {
	if user.ID == "" {
		user.ID = bson.NewObjectID().Hex()
	}
	cleanID, err := mongoinfra.SanitizeID(user.ID)
	if err != nil {
		return domainuser.User{}, err
	}
	user.ID = cleanID

	cleanEmail, err := mongoinfra.SanitizeEmail(user.Email)
	if err != nil {
		return domainuser.User{}, err
	}
	user.Email = cleanEmail

	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now().UTC()
	}
	if user.MaxVehicles <= 0 {
		user.MaxVehicles = domainuser.DefaultMaxVehicles
	}
	if user.MaxFiles <= 0 {
		user.MaxFiles = domainuser.DefaultMaxFiles
	}
	if user.MaxStorageBytes <= 0 {
		user.MaxStorageBytes = domainuser.DefaultMaxStorageBytes
	}

	doc := userDoc{
		ID:                         user.ID,
		Name:                       strings.TrimSpace(user.Name),
		Email:                      user.Email,
		PasswordHash:               user.PasswordHash,
		GoogleID:                   strings.TrimSpace(user.GoogleID),
		EmailVerified:              user.EmailVerified,
		EmailVerifiedAt:            user.EmailVerifiedAt,
		EmailVerificationTokenHash: user.EmailVerificationTokenHash,
		EmailVerificationExpiresAt: user.EmailVerificationExpiresAt,
		PasswordResetTokenHash:     user.PasswordResetTokenHash,
		PasswordResetExpiresAt:     user.PasswordResetExpiresAt,
		MaxVehicles:                user.MaxVehicles,
		FilesCount:                 user.FilesCount,
		MaxFiles:                   user.MaxFiles,
		StorageUsedBytes:           user.StorageUsedBytes,
		MaxStorageBytes:            user.MaxStorageBytes,
		StorageQuotaSynced:         user.StorageQuotaSynced,
		CreatedAt:                  user.CreatedAt,
		UpdatedAt:                  user.UpdatedAt,
	}

	_, err = r.coll.InsertOne(ctx, doc)
	if err != nil {
		return domainuser.User{}, err
	}
	return user, nil
}

func (r *Repository) Update(ctx context.Context, user domainuser.User) (domainuser.User, error) {
	cleanID, err := mongoinfra.SanitizeID(user.ID)
	if err != nil {
		return domainuser.User{}, userport.ErrNotFound
	}
	user.ID = cleanID

	if user.MaxVehicles <= 0 {
		user.MaxVehicles = domainuser.DefaultMaxVehicles
	}
	if user.MaxFiles <= 0 {
		user.MaxFiles = domainuser.DefaultMaxFiles
	}
	if user.MaxStorageBytes <= 0 {
		user.MaxStorageBytes = domainuser.DefaultMaxStorageBytes
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = time.Now().UTC()
	}

	doc := bson.M{
		"$set": bson.M{
			"name":                       strings.TrimSpace(user.Name),
			"email":                      user.Email,
			"passwordHash":               user.PasswordHash,
			"googleId":                   strings.TrimSpace(user.GoogleID),
			"emailVerified":              user.EmailVerified,
			"emailVerifiedAt":            user.EmailVerifiedAt,
			"emailVerificationTokenHash": user.EmailVerificationTokenHash,
			"emailVerificationExpiresAt": user.EmailVerificationExpiresAt,
			"passwordResetTokenHash":     user.PasswordResetTokenHash,
			"passwordResetExpiresAt":     user.PasswordResetExpiresAt,
			"maxVehicles":                user.MaxVehicles,
			"filesCount":                 user.FilesCount,
			"maxFiles":                   user.MaxFiles,
			"storageUsedBytes":           user.StorageUsedBytes,
			"maxStorageBytes":            user.MaxStorageBytes,
			"storageQuotaSynced":         user.StorageQuotaSynced,
			"updatedAt":                  user.UpdatedAt,
		},
	}

	filter := bson.D{{Key: "_id", Value: cleanID}}
	res, err := r.coll.UpdateOne(ctx, filter, doc)
	if err != nil {
		return domainuser.User{}, err
	}
	if res.MatchedCount == 0 {
		return domainuser.User{}, userport.ErrNotFound
	}

	return user, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	cleanID, err := mongoinfra.SanitizeID(id)
	if err != nil {
		return userport.ErrNotFound
	}

	filter := bson.D{{Key: "_id", Value: cleanID}}
	res, err := r.coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return userport.ErrNotFound
	}
	return nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (domainuser.User, error) {
	cleanEmail, err := mongoinfra.SanitizeEmail(email)
	if err != nil {
		return domainuser.User{}, userport.ErrNotFound
	}

	filter := bson.D{{Key: "email", Value: cleanEmail}}
	var doc userDoc
	err = r.coll.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domainuser.User{}, userport.ErrNotFound
		}
		return domainuser.User{}, err
	}
	return doc.toDomain(), nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (domainuser.User, error) {
	cleanID, err := mongoinfra.SanitizeID(id)
	if err != nil {
		return domainuser.User{}, userport.ErrNotFound
	}

	filter := bson.D{{Key: "_id", Value: cleanID}}
	var doc userDoc
	err = r.coll.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domainuser.User{}, userport.ErrNotFound
		}
		return domainuser.User{}, err
	}
	return doc.toDomain(), nil
}

func (r *Repository) FindByGoogleID(ctx context.Context, googleID string) (domainuser.User, error) {
	cleanID := strings.TrimSpace(googleID)
	if cleanID == "" {
		return domainuser.User{}, userport.ErrNotFound
	}

	filter := bson.D{{Key: "googleId", Value: cleanID}}
	var doc userDoc
	err := r.coll.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domainuser.User{}, userport.ErrNotFound
		}
		return domainuser.User{}, err
	}
	return doc.toDomain(), nil
}

func (r *Repository) FindByEmailVerificationTokenHash(ctx context.Context, hash string) (domainuser.User, error) {
	if strings.TrimSpace(hash) == "" {
		return domainuser.User{}, userport.ErrNotFound
	}

	filter := bson.D{{Key: "emailVerificationTokenHash", Value: strings.TrimSpace(hash)}}
	var doc userDoc
	err := r.coll.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domainuser.User{}, userport.ErrNotFound
		}
		return domainuser.User{}, err
	}
	return doc.toDomain(), nil
}

func (r *Repository) FindByPasswordResetTokenHash(ctx context.Context, hash string) (domainuser.User, error) {
	if strings.TrimSpace(hash) == "" {
		return domainuser.User{}, userport.ErrNotFound
	}

	filter := bson.D{{Key: "passwordResetTokenHash", Value: strings.TrimSpace(hash)}}
	var doc userDoc
	err := r.coll.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domainuser.User{}, userport.ErrNotFound
		}
		return domainuser.User{}, err
	}
	return doc.toDomain(), nil
}
