package mongo

import (
	"context"
	"errors"
	"strings"
	"time"

	carport "cvmc/internal/application/ports/car"
	domaincar "cvmc/internal/domain/car"
	mongoinfra "cvmc/internal/infrastructure/database/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type carDoc struct {
	ID              string    `bson:"_id"`
	OwnerID         string    `bson:"ownerId"`
	Name            string    `bson:"name"`
	Manufacturer    string    `bson:"manufacturer"`
	Model           string    `bson:"model"`
	YearManufacture int       `bson:"yearManufacture"`
	YearModel       int       `bson:"yearModel"`
	LastMileage     int       `bson:"lastMileage"`
	VehicleType     string    `bson:"vehicleType,omitempty"`
	ImageUrl        string    `bson:"imageUrl,omitempty"`
	FIPECode        string    `bson:"fipeCode,omitempty"`
	FIPEPrice       string    `bson:"fipePrice,omitempty"`
	Fuel            string    `bson:"fuel,omitempty"`
	SharedWith      []string  `bson:"sharedWith"`
	CreatedAt       time.Time `bson:"createdAt"`
	UpdatedAt       time.Time `bson:"updatedAt"`
}

type Repository struct {
	coll *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		coll: db.Collection("cars"),
	}
}

func (r *Repository) EnsureIndexes(ctx context.Context) error {
	models := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "ownerId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "sharedWith", Value: 1}},
		},
	}
	_, err := r.coll.Indexes().CreateMany(ctx, models)
	return err
}

func (r *Repository) Create(ctx context.Context, car domaincar.Car) (domaincar.Car, error) {
	if car.ID == "" {
		car.ID = bson.NewObjectID().Hex()
	}
	cleanID, err := mongoinfra.SanitizeID(car.ID)
	if err != nil {
		return domaincar.Car{}, err
	}
	car.ID = cleanID

	cleanOwnerID, err := mongoinfra.SanitizeID(car.OwnerID)
	if err != nil {
		return domaincar.Car{}, err
	}
	car.OwnerID = cleanOwnerID

	sanitizedShared := make([]string, 0, len(car.SharedWith))
	for _, s := range car.SharedWith {
		if sID, err := mongoinfra.SanitizeID(s); err == nil {
			sanitizedShared = append(sanitizedShared, sID)
		}
	}
	car.SharedWith = sanitizedShared

	if car.CreatedAt.IsZero() {
		car.CreatedAt = time.Now().UTC()
	}
	if car.UpdatedAt.IsZero() {
		car.UpdatedAt = car.CreatedAt
	}

	doc := toCarDoc(car)
	_, err = r.coll.InsertOne(ctx, doc)
	if err != nil {
		return domaincar.Car{}, err
	}
	return car, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (domaincar.Car, error) {
	cleanID, err := mongoinfra.SanitizeID(id)
	if err != nil {
		return domaincar.Car{}, carport.ErrNotFound
	}

	filter := bson.D{{Key: "_id", Value: cleanID}}
	var doc carDoc
	err = r.coll.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domaincar.Car{}, carport.ErrNotFound
		}
		return domaincar.Car{}, err
	}
	return toDomainCar(doc), nil
}

func (r *Repository) ListByUser(ctx context.Context, userID string) ([]domaincar.Car, error) {
	cleanUserID, err := mongoinfra.SanitizeID(userID)
	if err != nil {
		return []domaincar.Car{}, nil
	}

	filter := bson.D{
		{
			Key: "$or",
			Value: bson.A{
				bson.D{{Key: "ownerId", Value: cleanUserID}},
				bson.D{{Key: "sharedWith", Value: cleanUserID}},
			},
		},
	}
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []carDoc
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	cars := make([]domaincar.Car, 0, len(docs))
	for _, doc := range docs {
		cars = append(cars, toDomainCar(doc))
	}
	return cars, nil
}

func (r *Repository) Update(ctx context.Context, car domaincar.Car) (domaincar.Car, error) {
	cleanID, err := mongoinfra.SanitizeID(car.ID)
	if err != nil {
		return domaincar.Car{}, carport.ErrNotFound
	}
	car.ID = cleanID

	cleanOwnerID, err := mongoinfra.SanitizeID(car.OwnerID)
	if err != nil {
		return domaincar.Car{}, carport.ErrNotFound
	}
	car.OwnerID = cleanOwnerID

	doc := toCarDoc(car)
	filter := bson.D{{Key: "_id", Value: cleanID}}
	res, err := r.coll.ReplaceOne(ctx, filter, doc)
	if err != nil {
		return domaincar.Car{}, err
	}
	if res.MatchedCount == 0 {
		return domaincar.Car{}, carport.ErrNotFound
	}
	return car, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	cleanID, err := mongoinfra.SanitizeID(id)
	if err != nil {
		return carport.ErrNotFound
	}

	filter := bson.D{{Key: "_id", Value: cleanID}}
	res, err := r.coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return carport.ErrNotFound
	}
	return nil
}

func (r *Repository) Share(ctx context.Context, carID string, userID string) (domaincar.Car, error) {
	cleanCarID, err := mongoinfra.SanitizeID(carID)
	if err != nil {
		return domaincar.Car{}, carport.ErrNotFound
	}
	cleanUserID, err := mongoinfra.SanitizeID(userID)
	if err != nil {
		return domaincar.Car{}, carport.ErrNotFound
	}

	filter := bson.D{{Key: "_id", Value: cleanCarID}}
	update := bson.D{
		{
			Key: "$addToSet",
			Value: bson.D{
				{Key: "sharedWith", Value: cleanUserID},
			},
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var doc carDoc
	err = r.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domaincar.Car{}, carport.ErrNotFound
		}
		return domaincar.Car{}, err
	}
	return toDomainCar(doc), nil
}

func (r *Repository) Unshare(ctx context.Context, carID string, userID string) (domaincar.Car, error) {
	cleanCarID, err := mongoinfra.SanitizeID(carID)
	if err != nil {
		return domaincar.Car{}, carport.ErrNotFound
	}
	cleanUserID, err := mongoinfra.SanitizeID(userID)
	if err != nil {
		return domaincar.Car{}, carport.ErrNotFound
	}

	filter := bson.D{{Key: "_id", Value: cleanCarID}}
	update := bson.D{
		{
			Key: "$pull",
			Value: bson.D{
				{Key: "sharedWith", Value: cleanUserID},
			},
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var doc carDoc
	err = r.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domaincar.Car{}, carport.ErrNotFound
		}
		return domaincar.Car{}, err
	}
	return toDomainCar(doc), nil
}

func toCarDoc(c domaincar.Car) carDoc {
	shared := c.SharedWith
	if shared == nil {
		shared = []string{}
	}
	return carDoc{
		ID:              c.ID,
		OwnerID:         c.OwnerID,
		Name:            strings.TrimSpace(c.Name),
		Manufacturer:    strings.TrimSpace(c.Manufacturer),
		Model:           strings.TrimSpace(c.Model),
		YearManufacture: c.YearManufacture,
		YearModel:       c.YearModel,
		LastMileage:     c.LastMileage,
		VehicleType:     strings.TrimSpace(c.VehicleType),
		ImageUrl:        strings.TrimSpace(c.ImageUrl),
		FIPECode:        strings.TrimSpace(c.FIPECode),
		FIPEPrice:       strings.TrimSpace(c.FIPEPrice),
		Fuel:            strings.TrimSpace(c.Fuel),
		SharedWith:      shared,
		CreatedAt:       c.CreatedAt,
		UpdatedAt:       c.UpdatedAt,
	}
}

func toDomainCar(d carDoc) domaincar.Car {
	shared := d.SharedWith
	if shared == nil {
		shared = []string{}
	}
	return domaincar.Car{
		ID:              d.ID,
		OwnerID:         d.OwnerID,
		Name:            d.Name,
		Manufacturer:    d.Manufacturer,
		Model:           d.Model,
		YearManufacture: d.YearManufacture,
		YearModel:       d.YearModel,
		LastMileage:     d.LastMileage,
		VehicleType:     d.VehicleType,
		ImageUrl:        d.ImageUrl,
		FIPECode:        d.FIPECode,
		FIPEPrice:       d.FIPEPrice,
		Fuel:            d.Fuel,
		SharedWith:      shared,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}
}
