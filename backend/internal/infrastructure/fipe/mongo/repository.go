package mongo

import (
	"context"
	"errors"
	"time"

	domainfipe "cvmc/internal/domain/fipe"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository struct {
	coll *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		coll: db.Collection("marcas"),
	}
}

func (r *Repository) EnsureIndexes(ctx context.Context) error {
	models := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "vehicleType", Value: 1}, {Key: "code", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "models.code", Value: 1}},
		},
	}
	_, err := r.coll.Indexes().CreateMany(ctx, models)
	return err
}

func (r *Repository) GetBrands(ctx context.Context, vehicleType domainfipe.VehicleType) ([]domainfipe.Brand, time.Time, error) {
	filter := bson.D{{Key: "vehicleType", Value: string(vehicleType)}}
	opts := options.Find().
		SetProjection(bson.D{
			{Key: "code", Value: 1},
			{Key: "name", Value: 1},
			{Key: "vehicleType", Value: 1},
			{Key: "updatedAt", Value: 1},
		}).
		SetSort(bson.D{{Key: "name", Value: 1}})

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, time.Time{}, err
	}
	defer cursor.Close(ctx)

	type brandProj struct {
		Code        string    `bson:"code"`
		Name        string    `bson:"name"`
		VehicleType string    `bson:"vehicleType"`
		UpdatedAt   time.Time `bson:"updatedAt"`
	}

	var results []brandProj
	if err := cursor.All(ctx, &results); err != nil {
		return nil, time.Time{}, err
	}

	brands := make([]domainfipe.Brand, 0, len(results))
	var latestSync time.Time
	for _, b := range results {
		brands = append(brands, domainfipe.Brand{
			Code:        b.Code,
			Name:        b.Name,
			VehicleType: b.VehicleType,
		})
		if b.UpdatedAt.After(latestSync) {
			latestSync = b.UpdatedAt
		}
	}

	return brands, latestSync, nil
}

func (r *Repository) UpsertBrands(ctx context.Context, vehicleType domainfipe.VehicleType, brands []domainfipe.Brand, syncTime time.Time) error {
	if len(brands) == 0 {
		return nil
	}

	models := make([]mongo.WriteModel, 0, len(brands))
	for _, b := range brands {
		filter := bson.D{
			{Key: "vehicleType", Value: string(vehicleType)},
			{Key: "code", Value: b.Code},
		}
		update := bson.D{
			{
				Key: "$set",
				Value: bson.D{
					{Key: "name", Value: b.Name},
					{Key: "vehicleType", Value: string(vehicleType)},
					{Key: "updatedAt", Value: syncTime},
				},
			},
			{
				Key: "$setOnInsert",
				Value: bson.D{
					{Key: "createdAt", Value: syncTime},
					{Key: "models", Value: []domainfipe.ModelDocument{}},
				},
			},
		}
		model := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
		models = append(models, model)
	}

	_, err := r.coll.BulkWrite(ctx, models)
	return err
}

func (r *Repository) GetBrandWithModels(ctx context.Context, vehicleType domainfipe.VehicleType, brandCode string) (*domainfipe.BrandDocument, error) {
	filter := bson.D{
		{Key: "vehicleType", Value: string(vehicleType)},
		{Key: "code", Value: brandCode},
	}

	var doc domainfipe.BrandDocument
	err := r.coll.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domainfipe.ErrBrandNotFound
		}
		return nil, err
	}
	return &doc, nil
}

func (r *Repository) UpdateModels(ctx context.Context, vehicleType domainfipe.VehicleType, brandCode string, models []domainfipe.Model, syncTime time.Time) error {
	filter := bson.D{
		{Key: "vehicleType", Value: string(vehicleType)},
		{Key: "code", Value: brandCode},
	}

	// Fetch existing brand doc to merge existing years
	var existing domainfipe.BrandDocument
	err := r.coll.FindOne(ctx, filter).Decode(&existing)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}

	existingModelMap := make(map[string]domainfipe.ModelDocument)
	for _, m := range existing.Models {
		existingModelMap[m.Code] = m
	}

	modelDocs := make([]domainfipe.ModelDocument, 0, len(models))
	for _, m := range models {
		mDoc := domainfipe.ModelDocument{
			Code: m.Code,
			Name: m.Name,
		}
		if prev, ok := existingModelMap[m.Code]; ok {
			mDoc.Years = prev.Years
			mDoc.YearsLastSyncAt = prev.YearsLastSyncAt
		} else {
			mDoc.Years = []domainfipe.YearDocument{}
		}
		modelDocs = append(modelDocs, mDoc)
	}

	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "models", Value: modelDocs},
				{Key: "modelsLastSyncAt", Value: syncTime},
				{Key: "updatedAt", Value: syncTime},
			},
		},
	}

	res, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return domainfipe.ErrBrandNotFound
	}
	return nil
}

func (r *Repository) UpdateModelYears(ctx context.Context, vehicleType domainfipe.VehicleType, brandCode string, modelCode string, years []domainfipe.Year, syncTime time.Time) error {
	filter := bson.D{
		{Key: "vehicleType", Value: string(vehicleType)},
		{Key: "code", Value: brandCode},
		{Key: "models.code", Value: modelCode},
	}

	// Fetch existing to preserve year details
	var existing domainfipe.BrandDocument
	err := r.coll.FindOne(ctx, bson.D{
		{Key: "vehicleType", Value: string(vehicleType)},
		{Key: "code", Value: brandCode},
	}).Decode(&existing)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domainfipe.ErrBrandNotFound
		}
		return err
	}

	existingYearsMap := make(map[string]domainfipe.YearDocument)
	for _, m := range existing.Models {
		if m.Code == modelCode {
			for _, y := range m.Years {
				existingYearsMap[y.Code] = y
			}
			break
		}
	}

	yearDocs := make([]domainfipe.YearDocument, 0, len(years))
	for _, y := range years {
		yDoc := domainfipe.YearDocument{
			Code: y.Code,
			Name: y.Name,
		}
		if prev, ok := existingYearsMap[y.Code]; ok {
			yDoc.Price = prev.Price
			yDoc.PriceValue = prev.PriceValue
			yDoc.FIPECode = prev.FIPECode
			yDoc.Fuel = prev.Fuel
			yDoc.ReferenceMonth = prev.ReferenceMonth
			yDoc.DetailsLastSyncAt = prev.DetailsLastSyncAt
		}
		yearDocs = append(yearDocs, yDoc)
	}

	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "models.$.years", Value: yearDocs},
				{Key: "models.$.yearsLastSyncAt", Value: syncTime},
				{Key: "updatedAt", Value: syncTime},
			},
		},
	}

	res, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return domainfipe.ErrModelNotFound
	}
	return nil
}

func (r *Repository) UpdateYearDetail(ctx context.Context, vehicleType domainfipe.VehicleType, brandCode string, modelCode string, yearCode string, detail domainfipe.VehicleDetail, syncTime time.Time) error {
	filter := bson.D{
		{Key: "vehicleType", Value: string(vehicleType)},
		{Key: "code", Value: brandCode},
	}

	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "models.$[m].years.$[y].price", Value: detail.Price},
				{Key: "models.$[m].years.$[y].priceValue", Value: detail.PriceValue},
				{Key: "models.$[m].years.$[y].fipeCode", Value: detail.CodeFipe},
				{Key: "models.$[m].years.$[y].fuel", Value: detail.Fuel},
				{Key: "models.$[m].years.$[y].referenceMonth", Value: detail.ReferenceMonth},
				{Key: "models.$[m].years.$[y].detailsLastSyncAt", Value: syncTime},
				{Key: "updatedAt", Value: syncTime},
			},
		},
	}

	opts := options.UpdateOne().SetArrayFilters(bson.A{
		bson.D{{Key: "m.code", Value: modelCode}},
		bson.D{{Key: "y.code", Value: yearCode}},
	})

	res, err := r.coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return domainfipe.ErrBrandNotFound
	}
	return nil
}
