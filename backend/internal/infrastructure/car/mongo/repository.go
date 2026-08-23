package mongo

import (
	"context"
	"errors"
	"time"

	carport "cvmc/internal/application/ports/car"
	domaincar "cvmc/internal/domain/car"
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
	if car.SharedWith == nil {
		car.SharedWith = []string{}
	}
	if car.CreatedAt.IsZero() {
		car.CreatedAt = time.Now().UTC()
	}
	if car.UpdatedAt.IsZero() {
		car.UpdatedAt = car.CreatedAt
	}

	doc := toCarDoc(car)
	_, err := r.coll.InsertOne(ctx, doc)
	if err != nil {
		return domaincar.Car{}, err
	}
	return car, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (domaincar.Car, error) {
	var doc carDoc
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domaincar.Car{}, carport.ErrNotFound
		}
		return domaincar.Car{}, err
	}
	return toDomainCar(doc), nil
}

func (r *Repository) ListByUser(ctx context.Context, userID string) ([]domaincar.Car, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"ownerId": userID},
			{"sharedWith": userID},
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
	doc := toCarDoc(car)
	res, err := r.coll.ReplaceOne(ctx, bson.M{"_id": car.ID}, doc)
	if err != nil {
		return domaincar.Car{}, err
	}
	if res.MatchedCount == 0 {
		return domaincar.Car{}, carport.ErrNotFound
	}
	return car, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	res, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return carport.ErrNotFound
	}
	return nil
}

func (r *Repository) Share(ctx context.Context, carID string, userID string) (domaincar.Car, error) {
	update := bson.M{
		"$addToSet": bson.M{"sharedWith": userID},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var doc carDoc
	err := r.coll.FindOneAndUpdate(ctx, bson.M{"_id": carID}, update, opts).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domaincar.Car{}, carport.ErrNotFound
		}
		return domaincar.Car{}, err
	}
	return toDomainCar(doc), nil
}

func (r *Repository) Unshare(ctx context.Context, carID string, userID string) (domaincar.Car, error) {
	update := bson.M{
		"$pull": bson.M{"sharedWith": userID},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var doc carDoc
	err := r.coll.FindOneAndUpdate(ctx, bson.M{"_id": carID}, update, opts).Decode(&doc)
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
		Name:            c.Name,
		Manufacturer:    c.Manufacturer,
		Model:           c.Model,
		YearManufacture: c.YearManufacture,
		YearModel:       c.YearModel,
		LastMileage:     c.LastMileage,
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
		SharedWith:      shared,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}
}
