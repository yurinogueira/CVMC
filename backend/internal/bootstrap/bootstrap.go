package bootstrap

import (
	"context"
	"fmt"
	"log"

	carport "cvmc/internal/application/ports/car"
	emailport "cvmc/internal/application/ports/email"
	fipeport "cvmc/internal/application/ports/fipe"
	fuelport "cvmc/internal/application/ports/fuel"
	maintenanceport "cvmc/internal/application/ports/maintenance"
	notificationport "cvmc/internal/application/ports/notification"
	userport "cvmc/internal/application/ports/user"
	fipeusecase "cvmc/internal/application/usecase/fipe"
	notificationusecase "cvmc/internal/application/usecase/notification"
	"cvmc/internal/config"
	bcryptinfra "cvmc/internal/infrastructure/auth/bcrypt"
	googleinfra "cvmc/internal/infrastructure/auth/google"
	jwtauth "cvmc/internal/infrastructure/auth/jwt"
	carMemory "cvmc/internal/infrastructure/car/memory"
	carMongo "cvmc/internal/infrastructure/car/mongo"
	mongoinfra "cvmc/internal/infrastructure/database/mongo"
	emailinfra "cvmc/internal/infrastructure/email"
	fipeHttp "cvmc/internal/infrastructure/fipe/http"
	fipeMemory "cvmc/internal/infrastructure/fipe/memory"
	fipeMongo "cvmc/internal/infrastructure/fipe/mongo"
	fuelMemory "cvmc/internal/infrastructure/fuel/memory"
	fuelMongo "cvmc/internal/infrastructure/fuel/mongo"
	maintMemory "cvmc/internal/infrastructure/maintenance/memory"
	maintMongo "cvmc/internal/infrastructure/maintenance/mongo"
	fcmInfra "cvmc/internal/infrastructure/notification/fcm"
	notifMemory "cvmc/internal/infrastructure/notification/memory"
	notifMongo "cvmc/internal/infrastructure/notification/mongo"
	"cvmc/internal/infrastructure/scheduler"
	userMemory "cvmc/internal/infrastructure/user/memory"
	userMongo "cvmc/internal/infrastructure/user/mongo"
	"cvmc/internal/interfaces/rest"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type App struct {
	handler     *rest.Router
	mongoClient *mongo.Client
	cron        *scheduler.MaintenanceCron
}

func New(ctx context.Context, cfg config.Config) (*App, error) {
	var (
		users        userport.Repository
		cars         carport.Repository
		maintenances maintenanceport.Repository
		fuelings     fuelport.Repository
		fipes        fipeport.Repository
		deviceTokens notificationport.DeviceTokenRepository
		mongoClient  *mongo.Client
		emailSender  emailport.Sender
	)

	hasher := bcryptinfra.NewHasher()
	tokens := jwtauth.NewProvider(cfg.JWTSecret, cfg.JWTRefreshSecret)
	fipeClient := fipeHttp.NewClient(cfg.FIPEBaseURL, cfg.FIPEToken)
	emailSender = emailinfra.NewService(emailinfra.Config{
		SMTPHost:   cfg.SMTPHost,
		SMTPPort:   cfg.SMTPPort,
		SMTPUser:   cfg.SMTPUser,
		SMTPPass:   cfg.SMTPPass,
		EmailFrom:  cfg.EmailFrom,
		AppBaseURL: cfg.AppBaseURL,
	})

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
		fuelMongoRepo := fuelMongo.NewRepository(db)
		fMongo := fipeMongo.NewRepository(db)
		dtMongo := notifMongo.NewRepository(db)

		bootstrapper := mongoinfra.NewBootstrapper(
			db,
			[]string{"users", "cars", "maintenances", "fuelings", "brands", "device_tokens"},
			uMongo,
			cMongo,
			mMongo,
			fuelMongoRepo,
			fMongo,
			dtMongo,
		)

		if err := bootstrapper.Ensure(ctx); err != nil {
			_ = client.Disconnect(ctx)
			return nil, fmt.Errorf("failed to bootstrap mongodb collections and indexes: %w", err)
		}

		log.Printf("Successfully bootstrapped MongoDB collections (users, cars, maintenances, fuelings, brands, device_tokens) and indexes on %q", cfg.MongoDatabase)

		users = uMongo
		cars = cMongo
		maintenances = mMongo
		fuelings = fuelMongoRepo
		fipes = fMongo
		deviceTokens = dtMongo
	} else {
		log.Println("Initializing application with in-memory repositories (MONGO_URI=memory)")
		users = userMemory.NewRepository()
		cars = carMemory.NewRepository()
		maintenances = maintMemory.NewRepository()
		fuelings = fuelMemory.NewRepository()
		fipes = fipeMemory.NewRepository()
		deviceTokens = notifMemory.NewRepository()
	}

	pushSender, err := fcmInfra.NewSender(ctx, fcmInfra.Config{
		CredentialsJSON: cfg.FirebaseCredentialsJSON,
		CredentialsFile: cfg.FirebaseCredentialsFile,
	})
	if err != nil {
		if mongoClient != nil {
			_ = mongoClient.Disconnect(ctx)
		}
		return nil, fmt.Errorf("failed to initialize fcm push sender: %w", err)
	}

	fipeService := fipeusecase.NewService(fipes, fipeClient)
	notificationService := notificationusecase.NewService(deviceTokens, pushSender, cars, maintenances)
	googleValidator := googleinfra.NewValidator(cfg.GoogleClientID)
	handler := rest.NewRouter(cfg, users, hasher, tokens, emailSender, cars, maintenances, fuelings, fipeService, notificationService, googleValidator)

	cronInterval := scheduler.ParseCronInterval(cfg.MaintenanceCronInterval)
	maintCron := scheduler.NewMaintenanceCron(notificationService, cronInterval)
	maintCron.Start(context.Background())

	return &App{
		handler:     handler,
		mongoClient: mongoClient,
		cron:        maintCron,
	}, nil
}

func (a *App) Handler() *rest.Router {
	return a.handler
}

func (a *App) Close(ctx context.Context) error {
	if a.cron != nil {
		a.cron.Stop()
	}
	if a.mongoClient != nil {
		return a.mongoClient.Disconnect(ctx)
	}
	return nil
}
