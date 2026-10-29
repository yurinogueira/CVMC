package mongo

import (
	"context"
	"errors"
	"time"

	reportport "cvmc/internal/application/ports/report"
	domain "cvmc/internal/domain/report"
	mongoinfra "cvmc/internal/infrastructure/database/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var ErrNotFound = errors.New("report job not found")

type reportJobDoc struct {
	ID          string            `bson:"_id"`
	UserID      string            `bson:"userId"`
	CarID       string            `bson:"carId"`
	CarName     string            `bson:"carName"`
	Type        domain.ReportType `bson:"type"`
	Status      domain.Status     `bson:"status"`
	FilePath    string            `bson:"filePath,omitempty"`
	Error       string            `bson:"error,omitempty"`
	DurationMS  int64             `bson:"durationMs,omitempty"`
	CreatedAt   time.Time         `bson:"createdAt"`
	CompletedAt *time.Time        `bson:"completedAt,omitempty"`
	ExpiresAt   *time.Time        `bson:"expiresAt,omitempty"`
	ExpiredAt   *time.Time        `bson:"expiredAt,omitempty"`
}

type Repository struct {
	coll *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		coll: db.Collection("report_jobs"),
	}
}

func (r *Repository) EnsureIndexes(ctx context.Context) error {
	models := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "userId", Value: 1},
				{Key: "createdAt", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "createdAt", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "expiresAt", Value: 1},
			},
		},
	}
	_, err := r.coll.Indexes().CreateMany(ctx, models)
	return err
}

func (r *Repository) Create(ctx context.Context, job *domain.ReportJob) error {
	cleanUserID, err := mongoinfra.SanitizeID(job.UserID)
	if err != nil {
		return err
	}
	job.UserID = cleanUserID

	if job.ID == "" {
		job.ID = bson.NewObjectID().Hex()
	}
	cleanID, err := mongoinfra.SanitizeID(job.ID)
	if err != nil {
		return err
	}
	job.ID = cleanID

	if job.CreatedAt.IsZero() {
		job.CreatedAt = time.Now().UTC()
	}

	doc := reportJobDoc{
		ID:          job.ID,
		UserID:      job.UserID,
		CarID:       job.CarID,
		CarName:     job.CarName,
		Type:        job.Type,
		Status:      job.Status,
		FilePath:    job.FilePath,
		Error:       job.Error,
		DurationMS:  job.DurationMS,
		CreatedAt:   job.CreatedAt,
		CompletedAt: job.CompletedAt,
		ExpiresAt:   job.ExpiresAt,
		ExpiredAt:   job.ExpiredAt,
	}

	_, err = r.coll.InsertOne(ctx, doc)
	return err
}

func (r *Repository) Update(ctx context.Context, job *domain.ReportJob) error {
	cleanID, err := mongoinfra.SanitizeID(job.ID)
	if err != nil {
		return err
	}
	cleanUserID, err := mongoinfra.SanitizeID(job.UserID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: cleanID},
		{Key: "userId", Value: cleanUserID},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "status", Value: job.Status},
			{Key: "filePath", Value: job.FilePath},
			{Key: "error", Value: job.Error},
			{Key: "durationMs", Value: job.DurationMS},
			{Key: "completedAt", Value: job.CompletedAt},
			{Key: "expiresAt", Value: job.ExpiresAt},
			{Key: "expiredAt", Value: job.ExpiredAt},
		}},
	}

	res, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id, userID string) (*domain.ReportJob, error) {
	cleanID, err := mongoinfra.SanitizeID(id)
	if err != nil {
		return nil, err
	}
	cleanUserID, err := mongoinfra.SanitizeID(userID)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		{Key: "_id", Value: cleanID},
		{Key: "userId", Value: cleanUserID},
	}

	var doc reportJobDoc
	if err := r.coll.FindOne(ctx, filter).Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return toDomain(doc), nil
}

func (r *Repository) ListByUser(ctx context.Context, filter reportport.ListFilter) (*reportport.ListResult, error) {
	cleanUserID, err := mongoinfra.SanitizeID(filter.UserID)
	if err != nil {
		return nil, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	skip := int64((page - 1) * limit)

	filterConditions := bson.D{{Key: "userId", Value: cleanUserID}}
	if filter.CarID != "" {
		filterConditions = append(filterConditions, bson.E{Key: "carId", Value: filter.CarID})
	}

	total, err := r.coll.CountDocuments(ctx, filterConditions)
	if err != nil {
		return nil, err
	}

	findOpts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(skip).
		SetLimit(int64(limit))

	cursor, err := r.coll.Find(ctx, filterConditions, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []reportJobDoc
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	jobs := make([]*domain.ReportJob, len(docs))
	for i, d := range docs {
		jobs[i] = toDomain(d)
	}

	return &reportport.ListResult{
		Jobs:  jobs,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (r *Repository) FindPending(ctx context.Context, limit int) ([]*domain.ReportJob, error) {
	filter := bson.D{
		{Key: "status", Value: bson.D{
			{Key: "$in", Value: bson.A{domain.StatusPending, domain.StatusProcessing}},
		}},
	}

	findOpts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: 1}})

	if limit > 0 {
		findOpts.SetLimit(int64(limit))
	}

	cursor, err := r.coll.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []reportJobDoc
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	jobs := make([]*domain.ReportJob, len(docs))
	for i, d := range docs {
		jobs[i] = toDomain(d)
	}
	return jobs, nil
}

func (r *Repository) FindExpiredCompleted(ctx context.Context, cutoff time.Time, batchSize int) ([]*domain.ReportJob, error) {
	filter := bson.D{
		{Key: "status", Value: domain.StatusCompleted},
		{Key: "expiresAt", Value: bson.D{
			{Key: "$lte", Value: cutoff},
		}},
	}

	findOpts := options.Find().
		SetSort(bson.D{{Key: "expiresAt", Value: 1}})

	if batchSize > 0 {
		findOpts.SetLimit(int64(batchSize))
	}

	cursor, err := r.coll.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []reportJobDoc
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	jobs := make([]*domain.ReportJob, len(docs))
	for i, d := range docs {
		jobs[i] = toDomain(d)
	}
	return jobs, nil
}

func (r *Repository) MarkAsExpired(ctx context.Context, id string, expiredAt time.Time) error {
	cleanID, err := mongoinfra.SanitizeID(id)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: cleanID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "status", Value: domain.StatusExpired},
			{Key: "expiredAt", Value: expiredAt},
		}},
	}

	res, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func toDomain(d reportJobDoc) *domain.ReportJob {
	return &domain.ReportJob{
		ID:          d.ID,
		UserID:      d.UserID,
		CarID:       d.CarID,
		CarName:     d.CarName,
		Type:        d.Type,
		Status:      d.Status,
		FilePath:    d.FilePath,
		Error:       d.Error,
		DurationMS:  d.DurationMS,
		CreatedAt:   d.CreatedAt,
		CompletedAt: d.CompletedAt,
		ExpiresAt:   d.ExpiresAt,
		ExpiredAt:   d.ExpiredAt,
	}
}
