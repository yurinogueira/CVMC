package bootstrap

import (
	"context"
	"fmt"
	"log"

	carport "cvmc/internal/application/ports/car"
	fipeport "cvmc/internal/application/ports/fipe"
	maintenanceport "cvmc/internal/application/ports/maintenance"
	userport "cvmc/internal/application/ports/user"
	fipeusecase "cvmc/internal/application/usecase/fipe"
	"cvmc/internal/config"
	bcryptinfra "cvmc/internal/infrastructure/auth/bcrypt"
	jwtauth "cvmc/internal/infrastructure/auth/jwt"
	carMemory "cvmc/internal/infrastructure/car/memory"
	carMongo "cvmc/internal/infrastructure/car/mongo"
	mongoinfra "cvmc/internal/infrastructure/database/mongo"
	fipeHttp "cvmc/internal/infrastructure/fipe/http"
	fipeMemory "cvmc/internal/infrastructure/fipe/memory"
	fipeMongo "cvmc/internal/infrastructure/fipe/mongo"
	maintMemory "cvmc/internal/infrastructure/maintenance/memory"
	maintMongo "cvmc/internal/infrastructure/maintenance/mongo"
	userMemory "cvmc/internal/infrastructure/user/memory"
	userMongo "cvmc/internal/infrastructure/user/mongo"
	"cvmc/internal/interfaces/rest"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type App struct {
	handler     *rest.Router
	mongoClient *mongo.Client
}

func New(ctx context.Context, cfg config.Config) (*App, error) {
	var (
		users        userport.Repository
		cars         carport.Repository
		maintenances maintenanceport.Repository
		fipes        fipeport.Repository
		mongoClient  *mongo.Client
	)

	hasher := bcryptinfra.NewHasher()
	tokens := jwtauth.NewProvider(cfg.JWTSecret, cfg.JWTRefreshSecret)
	fipeClient := fipeHttp.NewClient(cfg.FIPEBaseURL, cfg.FIPEToken)

	if cfg.MongoURI != "" && cfg.MongoURI != "memory" {
		log.Printf("Connecting to MongoDB database %q...", cfg.MongoDatabase)
		client, err := mongoinfra.Connect(ctx, cfg.MongoURI)
		if err != nil {
			return nil, fmt.Errorf("could not connect to mongodb: %w", err)
		}
		mongoClient = client

		db := client.Database(cfg.MongoDatabase)
		uMongo := userMongo.NewRepository(db)
		cMongo := carMongo.NewRepository(db)
		mMongo := maintMongo.NewRepository(db)
		fMongo := fipeMongo.NewRepository(db)

		bootstrapper := mongoinfra.NewBootstrapper(
			db,
			[]string{"users", "cars", "maintenances", "marcas"},
			uMongo,
			cMongo,
			mMongo,
			fMongo,
		)

		if err := bootstrapper.Ensure(ctx); err != nil {
			_ = client.Disconnect(ctx)
			return nil, fmt.Errorf("failed to bootstrap mongodb collections and indexes: %w", err)
		}

		log.Printf("Successfully bootstrapped MongoDB collections (users, cars, maintenances, marcas) and indexes on %q", cfg.MongoDatabase)

		users = uMongo
		cars = cMongo
		maintenances = mMongo
		fipes = fMongo
	} else {
		log.Println("Initializing application with in-memory repositories (MONGO_URI=memory)")
		users = userMemory.NewRepository()
		cars = carMemory.NewRepository()
		maintenances = maintMemory.NewRepository()
		fipes = fipeMemory.NewRepository()
	}

	fipeService := fipeusecase.NewService(fipes, fipeClient)
	handler := rest.NewRouter(cfg, users, hasher, tokens, cars, maintenances, fipeService)
	return &App{
		handler:     handler,
		mongoClient: mongoClient,
	}, nil
}

func (a *App) Handler() *rest.Router {
	return a.handler
}

func (a *App) Close(ctx context.Context) error {
	if a.mongoClient != nil {
		return a.mongoClient.Disconnect(ctx)
	}
	return nil
}
