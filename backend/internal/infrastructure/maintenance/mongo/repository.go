package mongo

import (
	"context"
	"errors"
	"time"

	maintenanceport "cvmc/internal/application/ports/maintenance"
	domainmaintenance "cvmc/internal/domain/maintenance"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type maintenanceDoc struct {
	ID          string    `bson:"_id"`
	CarID       string    `bson:"carId"`
	Title       string    `bson:"title"`
	Description string    `bson:"description"`
	Date        time.Time `bson:"date"`
	Mileage     int       `bson:"mileage"`
	CreatedAt   time.Time `bson:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt"`
}

type Repository struct {
	coll *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		coll: db.Collection("maintenances"),
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

func (r *Repository) Create(ctx context.Context, m domainmaintenance.Maintenance) (domainmaintenance.Maintenance, error) {
	if m.ID == "" {
		m.ID = bson.NewObjectID().Hex()
	}
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now().UTC()
	}
	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = m.CreatedAt
	}

	doc := toMaintenanceDoc(m)
	_, err := r.coll.InsertOne(ctx, doc)
	if err != nil {
		return domainmaintenance.Maintenance{}, err
	}
	return m, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (domainmaintenance.Maintenance, error) {
	var doc maintenanceDoc
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domainmaintenance.Maintenance{}, maintenanceport.ErrNotFound
		}
		return domainmaintenance.Maintenance{}, err
	}
	return toDomainMaintenance(doc), nil
}

func (r *Repository) ListByCar(ctx context.Context, carID string) ([]domainmaintenance.Maintenance, error) {
	filter := bson.M{"carId": carID}
	opts := options.Find().SetSort(bson.D{{Key: "date", Value: -1}})
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []maintenanceDoc
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	items := make([]domainmaintenance.Maintenance, 0, len(docs))
	for _, doc := range docs {
		items = append(items, toDomainMaintenance(doc))
	}
	return items, nil
}

func (r *Repository) Update(ctx context.Context, m domainmaintenance.Maintenance) (domainmaintenance.Maintenance, error) {
	doc := toMaintenanceDoc(m)
	res, err := r.coll.ReplaceOne(ctx, bson.M{"_id": m.ID}, doc)
	if err != nil {
		return domainmaintenance.Maintenance{}, err
	}
	if res.MatchedCount == 0 {
		return domainmaintenance.Maintenance{}, maintenanceport.ErrNotFound
	}
	return m, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	res, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return maintenanceport.ErrNotFound
	}
	return nil
}

func toMaintenanceDoc(m domainmaintenance.Maintenance) maintenanceDoc {
	return maintenanceDoc{
		ID:          m.ID,
		CarID:       m.CarID,
		Title:       m.Title,
		Description: m.Description,
		Date:        m.Date,
		Mileage:     m.Mileage,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func toDomainMaintenance(d maintenanceDoc) domainmaintenance.Maintenance {
	return domainmaintenance.Maintenance{
		ID:          d.ID,
		CarID:       d.CarID,
		Title:       d.Title,
		Description: d.Description,
		Date:        d.Date,
		Mileage:     d.Mileage,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}
