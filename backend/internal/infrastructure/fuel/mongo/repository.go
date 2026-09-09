package mongo

import (
	"context"
	"errors"
	"strings"
	"time"

	fuelport "cvmc/internal/application/ports/fuel"
	domainfuel "cvmc/internal/domain/fuel"
	mongoinfra "cvmc/internal/infrastructure/database/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type fuelingDoc struct {
	ID            string    `bson:"_id"`
	CarID         string    `bson:"carId"`
	Date          time.Time `bson:"date"`
	FuelType      string    `bson:"fuelType"`
	Liters        float64   `bson:"liters"`
	PricePerLiter float64   `bson:"pricePerLiter"`
	TotalCost     float64   `bson:"totalCost"`
	IsFullTank    bool      `bson:"isFullTank"`
	GasStation    string    `bson:"gasStation,omitempty"`
	Notes         string    `bson:"notes,omitempty"`
	CreatedAt     time.Time `bson:"createdAt"`
	UpdatedAt     time.Time `bson:"updatedAt"`
}

type Repository struct {
	coll *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		coll: db.Collection("fuelings"),
	}
}

func (r *Repository) EnsureIndexes(ctx context.Context) error {
	models := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "carId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "date", Value: -1}},
		},
	}
	_, err := r.coll.Indexes().CreateMany(ctx, models)
	return err
}

func (r *Repository) Create(ctx context.Context, f domainfuel.Fueling) (domainfuel.Fueling, error) {
	if f.ID == "" {
		f.ID = bson.NewObjectID().Hex()
	}
	cleanID, err := mongoinfra.SanitizeID(f.ID)
	if err != nil {
		return domainfuel.Fueling{}, err
	}
	f.ID = cleanID

	cleanCarID, err := mongoinfra.SanitizeID(f.CarID)
	if err != nil {
		return domainfuel.Fueling{}, err
	}
	f.CarID = cleanCarID

	doc := toDoc(f)
	if _, err := r.coll.InsertOne(ctx, doc); err != nil {
		return domainfuel.Fueling{}, err
	}
	return f, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (domainfuel.Fueling, error) {
	cleanID, err := mongoinfra.SanitizeID(id)
	if err != nil {
		return domainfuel.Fueling{}, fuelport.ErrNotFound
	}
	var doc fuelingDoc
	if err := r.coll.FindOne(ctx, bson.M{"_id": cleanID}).Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domainfuel.Fueling{}, fuelport.ErrNotFound
		}
		return domainfuel.Fueling{}, err
	}
	return fromDoc(doc), nil
}

func (r *Repository) ListByCar(ctx context.Context, carID string) ([]domainfuel.Fueling, error) {
	cleanCarID, err := mongoinfra.SanitizeID(carID)
	if err != nil {
		return nil, fuelport.ErrNotFound
	}
	cursor, err := r.coll.Find(ctx, bson.M{"carId": cleanCarID}, options.Find().SetSort(bson.D{{Key: "date", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []fuelingDoc
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}
	items := make([]domainfuel.Fueling, 0, len(docs))
	for _, doc := range docs {
		items = append(items, fromDoc(doc))
	}
	return items, nil
}

func (r *Repository) Update(ctx context.Context, f domainfuel.Fueling) (domainfuel.Fueling, error) {
	cleanID, err := mongoinfra.SanitizeID(f.ID)
	if err != nil {
		return domainfuel.Fueling{}, fuelport.ErrNotFound
	}
	f.ID = cleanID

	doc := toDoc(f)
	res, err := r.coll.ReplaceOne(ctx, bson.M{"_id": cleanID}, doc)
	if err != nil {
		return domainfuel.Fueling{}, err
	}
	if res.MatchedCount == 0 {
		return domainfuel.Fueling{}, fuelport.ErrNotFound
	}
	return f, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	cleanID, err := mongoinfra.SanitizeID(id)
	if err != nil {
		return fuelport.ErrNotFound
	}
	res, err := r.coll.DeleteOne(ctx, bson.M{"_id": cleanID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return fuelport.ErrNotFound
	}
	return nil
}

func toDoc(f domainfuel.Fueling) fuelingDoc {
	return fuelingDoc{
		ID:            f.ID,
		CarID:         f.CarID,
		Date:          f.Date,
		FuelType:      f.FuelType,
		Liters:        f.Liters,
		PricePerLiter: f.PricePerLiter,
		TotalCost:     f.TotalCost,
		IsFullTank:    f.IsFullTank,
		GasStation:    strings.TrimSpace(f.GasStation),
		Notes:         strings.TrimSpace(f.Notes),
		CreatedAt:     f.CreatedAt,
		UpdatedAt:     f.UpdatedAt,
	}
}

func fromDoc(d fuelingDoc) domainfuel.Fueling {
	return domainfuel.Fueling{
		ID:            d.ID,
		CarID:         d.CarID,
		Date:          d.Date,
		FuelType:      d.FuelType,
		Liters:        d.Liters,
		PricePerLiter: d.PricePerLiter,
		TotalCost:     d.TotalCost,
		IsFullTank:    d.IsFullTank,
		GasStation:    d.GasStation,
		Notes:         d.Notes,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
}
