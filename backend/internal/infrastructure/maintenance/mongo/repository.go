package mongo

import (
	"context"
	"errors"
	"strings"
	"time"

	maintenanceport "cvmc/internal/application/ports/maintenance"
	domainmaintenance "cvmc/internal/domain/maintenance"
	mongoinfra "cvmc/internal/infrastructure/database/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type maintenanceDoc struct {
	ID          string                         `bson:"_id"`
	CarID       string                         `bson:"carId"`
	Title       string                         `bson:"title"`
	Description string                         `bson:"description"`
	Date        time.Time                      `bson:"date"`
	Mileage     int                            `bson:"mileage"`
	Types       []string                       `bson:"types,omitempty"`
	Cost        *float64                       `bson:"cost,omitempty"`
	Attachments []domainmaintenance.Attachment `bson:"attachments,omitempty"`
	CreatedAt   time.Time                      `bson:"createdAt"`
	UpdatedAt   time.Time                      `bson:"updatedAt"`
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
	cleanID, err := mongoinfra.SanitizeID(m.ID)
	if err != nil {
		return domainmaintenance.Maintenance{}, err
	}
	m.ID = cleanID

	cleanCarID, err := mongoinfra.SanitizeID(m.CarID)
	if err != nil {
		return domainmaintenance.Maintenance{}, err
	}
	m.CarID = cleanCarID

	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now().UTC()
	}
	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = m.CreatedAt
	}

	doc := toMaintenanceDoc(m)
	_, err = r.coll.InsertOne(ctx, doc)
	if err != nil {
		return domainmaintenance.Maintenance{}, err
	}
	return m, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (domainmaintenance.Maintenance, error) {
	cleanID, err := mongoinfra.SanitizeID(id)
	if err != nil {
		return domainmaintenance.Maintenance{}, maintenanceport.ErrNotFound
	}

	filter := bson.D{{Key: "_id", Value: cleanID}}
	var doc maintenanceDoc
	err = r.coll.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domainmaintenance.Maintenance{}, maintenanceport.ErrNotFound
		}
		return domainmaintenance.Maintenance{}, err
	}
	return toDomainMaintenance(doc), nil
}

func (r *Repository) ListByCar(ctx context.Context, carID string) ([]domainmaintenance.Maintenance, error) {
	cleanCarID, err := mongoinfra.SanitizeID(carID)
	if err != nil {
		return []domainmaintenance.Maintenance{}, nil
	}

	filter := bson.D{{Key: "carId", Value: cleanCarID}}
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
	cleanID, err := mongoinfra.SanitizeID(m.ID)
	if err != nil {
		return domainmaintenance.Maintenance{}, maintenanceport.ErrNotFound
	}
	m.ID = cleanID

	cleanCarID, err := mongoinfra.SanitizeID(m.CarID)
	if err != nil {
		return domainmaintenance.Maintenance{}, maintenanceport.ErrNotFound
	}
	m.CarID = cleanCarID

	doc := toMaintenanceDoc(m)
	filter := bson.D{{Key: "_id", Value: cleanID}}
	res, err := r.coll.ReplaceOne(ctx, filter, doc)
	if err != nil {
		return domainmaintenance.Maintenance{}, err
	}
	if res.MatchedCount == 0 {
		return domainmaintenance.Maintenance{}, maintenanceport.ErrNotFound
	}
	return m, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	cleanID, err := mongoinfra.SanitizeID(id)
	if err != nil {
		return maintenanceport.ErrNotFound
	}

	filter := bson.D{{Key: "_id", Value: cleanID}}
	res, err := r.coll.DeleteOne(ctx, filter)
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
		Title:       strings.TrimSpace(m.Title),
		Description: strings.TrimSpace(m.Description),
		Date:        m.Date,
		Mileage:     m.Mileage,
		Types:       m.Types,
		Cost:        m.Cost,
		Attachments: m.Attachments,
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
		Types:       d.Types,
		Cost:        d.Cost,
		Attachments: d.Attachments,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}
