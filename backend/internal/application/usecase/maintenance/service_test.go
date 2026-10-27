package maintenance

import (
	"context"
	"testing"
	"time"

	"cvmc/internal/application/usecase/auth"
	"cvmc/internal/application/usecase/car"
	domaincar "cvmc/internal/domain/car"
	domainmaintenance "cvmc/internal/domain/maintenance"
	"cvmc/internal/infrastructure/auth/bcrypt"
	jwtauth "cvmc/internal/infrastructure/auth/jwt"
	carrepo "cvmc/internal/infrastructure/car/memory"
	maintrepo "cvmc/internal/infrastructure/maintenance/memory"
	memoryuser "cvmc/internal/infrastructure/user/memory"
	"cvmc/internal/shared/filecipher"
)

func TestMaintenanceServiceUpdatesMileage(t *testing.T) {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret", "refresh-secret")
	authService := auth.NewService(users, hasher, tokens, nil)
	ctx := context.Background()

	owner, err := authService.Register(ctx, auth.RegisterInput{Name: "Owner", Email: "owner@example.com", Password: "secret123"})
	if err != nil {
		t.Fatalf("owner register failed: %v", err)
	}

	owner.User.EmailVerified = true
	if _, err := users.Update(ctx, owner.User); err != nil {
		t.Fatalf("update user failed: %v", err)
	}

	cars := carrepo.NewRepository()
	carService := car.NewService(cars, users)
	createdCar, err := carService.Create(ctx, owner.User.ID, car.CreateInput{Name: "Carro 1", Manufacturer: "Fiat", Model: "Uno", YearManufacture: 2010, YearModel: 2011, LastMileage: 10000})
	if err != nil {
		t.Fatalf("create car failed: %v", err)
	}

	service := NewService(maintrepo.NewRepository(), cars)
	createdMaintenance, err := service.Create(ctx, owner.User.ID, createdCar.ID, CreateInput{Title: "Troca de óleo", Description: "Filtro e óleo", Date: time.Now().UTC(), Mileage: 12000})
	if err != nil {
		t.Fatalf("create maintenance failed: %v", err)
	}
	if createdMaintenance.Mileage != 12000 {
		t.Fatalf("unexpected mileage: %d", createdMaintenance.Mileage)
	}
	updatedCar, err := cars.GetByID(ctx, createdCar.ID)
	if err != nil {
		t.Fatalf("get car failed: %v", err)
	}
	if updatedCar.LastMileage != 12000 {
		t.Fatalf("expected last mileage updated, got %d", updatedCar.LastMileage)
	}
}

func TestMaintenanceServiceUpdateAndForbidden(t *testing.T) {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret", "refresh-secret")
	authService := auth.NewService(users, hasher, tokens, nil)
	ctx := context.Background()

	owner, err := authService.Register(ctx, auth.RegisterInput{Name: "Owner", Email: "owner2@example.com", Password: "secret123"})
	if err != nil {
		t.Fatalf("owner register failed: %v", err)
	}
	other, err := authService.Register(ctx, auth.RegisterInput{Name: "Other", Email: "other2@example.com", Password: "secret123"})
	if err != nil {
		t.Fatalf("other register failed: %v", err)
	}

	owner.User.EmailVerified = true
	other.User.EmailVerified = true
	_, _ = users.Update(ctx, owner.User)
	_, _ = users.Update(ctx, other.User)

	cars := carrepo.NewRepository()
	carService := car.NewService(cars, users)
	createdCar, err := carService.Create(ctx, owner.User.ID, car.CreateInput{Name: "Carro 2", Manufacturer: "VW", Model: "Gol", YearManufacture: 2015, YearModel: 2016, LastMileage: 50000})
	if err != nil {
		t.Fatalf("create car failed: %v", err)
	}

	service := NewService(maintrepo.NewRepository(), cars)
	created, err := service.Create(ctx, owner.User.ID, createdCar.ID, CreateInput{
		Title:   "Troca de óleo inicial",
		Date:    time.Now().UTC(),
		Mileage: 51000,
	})
	if err != nil {
		t.Fatalf("create maintenance failed: %v", err)
	}

	// 1. Successful update of fields and mileage > car.LastMileage
	costVal := 450.0
	updated, err := service.Update(ctx, owner.User.ID, created.ID, UpdateInput{
		Title:       "Revisão dos 55.000 km",
		Description: "Troca de óleo e filtro de combustível",
		Date:        time.Now().UTC(),
		Mileage:     55000,
		Types:       []string{"Óleo de Motor", "Filtro de Combustível"},
		Cost:        &costVal,
	})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Title != "Revisão dos 55.000 km" || updated.Mileage != 55000 || updated.Cost == nil || *updated.Cost != 450.0 {
		t.Fatalf("unexpected updated fields: %+v", updated)
	}

	// Verify car mileage was bumped to 55000
	c, err := cars.GetByID(ctx, createdCar.ID)
	if err != nil || c.LastMileage != 55000 {
		t.Fatalf("expected car last mileage 55000, got %d", c.LastMileage)
	}

	// 2. Forbidden when other user tries to update
	_, err = service.Update(ctx, other.User.ID, created.ID, UpdateInput{Title: "Hacked"})
	if err != ErrMaintenanceForbidden {
		t.Fatalf("expected ErrMaintenanceForbidden, got %v", err)
	}

	// 3. Invalid when title is empty
	_, err = service.Update(ctx, owner.User.ID, created.ID, UpdateInput{Title: ""})
	if err != ErrMaintenanceInvalid {
		t.Fatalf("expected ErrMaintenanceInvalid, got %v", err)
	}

	// 4. Not found when maintenance does not exist
	_, err = service.Update(ctx, owner.User.ID, "non-existent-id", UpdateInput{Title: "Alguma coisa"})
	if err != ErrMaintenanceNotFound {
		t.Fatalf("expected ErrMaintenanceNotFound, got %v", err)
	}
}

func TestMaintenanceServiceDeleteAndForbidden(t *testing.T) {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret", "refresh-secret")
	authService := auth.NewService(users, hasher, tokens, nil)
	ctx := context.Background()

	owner, _ := authService.Register(ctx, auth.RegisterInput{Name: "Owner", Email: "owner3@example.com", Password: "secret123"})
	other, _ := authService.Register(ctx, auth.RegisterInput{Name: "Other", Email: "other3@example.com", Password: "secret123"})

	owner.User.EmailVerified = true
	other.User.EmailVerified = true
	_, _ = users.Update(ctx, owner.User)
	_, _ = users.Update(ctx, other.User)

	cars := carrepo.NewRepository()
	carService := car.NewService(cars, users)
	createdCar, _ := carService.Create(ctx, owner.User.ID, car.CreateInput{Name: "Carro 3", Manufacturer: "Honda", Model: "Fit", YearManufacture: 2018, YearModel: 2019, LastMileage: 60000})

	service := NewService(maintrepo.NewRepository(), cars)
	created, err := service.Create(ctx, owner.User.ID, createdCar.ID, CreateInput{
		Title:   "Troca de pastilhas",
		Date:    time.Now().UTC(),
		Mileage: 60000,
	})
	if err != nil {
		t.Fatalf("create maintenance failed: %v", err)
	}

	// 1. Forbidden when other user tries to delete
	err = service.Delete(ctx, other.User.ID, created.ID)
	if err != ErrMaintenanceForbidden {
		t.Fatalf("expected ErrMaintenanceForbidden, got %v", err)
	}

	// 2. Successful delete by owner
	err = service.Delete(ctx, owner.User.ID, created.ID)
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	// 3. Not found when trying to delete again
	err = service.Delete(ctx, owner.User.ID, created.ID)
	if err != ErrMaintenanceNotFound {
		t.Fatalf("expected ErrMaintenanceNotFound, got %v", err)
	}
}

func TestMaintenanceServiceGet(t *testing.T) {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret", "refresh-secret")
	authService := auth.NewService(users, hasher, tokens, nil)
	ctx := context.Background()

	owner, _ := authService.Register(ctx, auth.RegisterInput{Name: "Owner", Email: "owner4@example.com", Password: "secret123"})
	other, _ := authService.Register(ctx, auth.RegisterInput{Name: "Other", Email: "other4@example.com", Password: "secret123"})

	owner.User.EmailVerified = true
	other.User.EmailVerified = true
	_, _ = users.Update(ctx, owner.User)
	_, _ = users.Update(ctx, other.User)

	cars := carrepo.NewRepository()
	carService := car.NewService(cars, users)
	createdCar, _ := carService.Create(ctx, owner.User.ID, car.CreateInput{Name: "Carro 4", Manufacturer: "Toyota", Model: "Etios", YearManufacture: 2017, YearModel: 2018, LastMileage: 80000})

	service := NewService(maintrepo.NewRepository(), cars)
	created, err := service.Create(ctx, owner.User.ID, createdCar.ID, CreateInput{
		Title:   "Alinhamento e Balanceamento",
		Date:    time.Now().UTC(),
		Mileage: 80000,
	})
	if err != nil {
		t.Fatalf("create maintenance failed: %v", err)
	}

	// 1. Get success by owner
	item, err := service.Get(ctx, owner.User.ID, created.ID)
	if err != nil || item.ID != created.ID {
		t.Fatalf("expected item, got err: %v", err)
	}

	// 2. Forbidden when other user tries to get
	_, err = service.Get(ctx, other.User.ID, created.ID)
	if err != ErrMaintenanceForbidden {
		t.Fatalf("expected ErrMaintenanceForbidden, got %v", err)
	}

	// 3. Not found when ID does not exist
	_, err = service.Get(ctx, owner.User.ID, "non-existent")
	if err != ErrMaintenanceNotFound {
		t.Fatalf("expected ErrMaintenanceNotFound, got %v", err)
	}
}

func TestMaintenanceServiceAttachmentEncryptionAndQuota(t *testing.T) {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret-32-chars-long-security", "refresh-secret-32-chars-long-security")
	authService := auth.NewService(users, hasher, tokens, nil)
	ctx := context.Background()

	owner, _ := authService.Register(ctx, auth.RegisterInput{Name: "Owner", Email: "owner5@example.com", Password: "secret123"})
	owner.User.EmailVerified = true
	owner.User.MaxFiles = 2
	owner.User.MaxStorageBytes = 500
	_, _ = users.Update(ctx, owner.User)

	cars := carrepo.NewRepository()
	maintRepo := maintrepo.NewRepository()
	fc := filecipher.New("master-file-encryption-key-with-32-chars!!")

	carService := car.NewService(cars, users)
	carService.SetSecurityDependencies(maintRepo, fc)

	sampleDataURL := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="
	createdCar, err := carService.Create(ctx, owner.User.ID, car.CreateInput{
		Name:            "Carro Criptografado",
		Manufacturer:    "Honda",
		Model:           "Civic",
		YearManufacture: 2022,
		YearModel:       2023,
		LastMileage:     10000,
		ImageUrl:        sampleDataURL,
	})
	if err != nil {
		t.Fatalf("create car failed: %v", err)
	}
	if createdCar.ImageUrl != sampleDataURL {
		t.Fatalf("expected decrypted ImageUrl on return, got %q", createdCar.ImageUrl)
	}

	// Verify stored car image in repo is encrypted at rest
	rawCar, _ := cars.GetByID(ctx, createdCar.ID)
	if !filecipher.IsEncrypted(rawCar.ImageUrl) {
		t.Fatalf("expected car ImageUrl to be encrypted at rest, got %q", rawCar.ImageUrl)
	}

	service := NewService(maintRepo, cars)
	service.SetSecurityDependencies(users, fc)

	// 1 file is already used by Car.ImageUrl; adding 1 attachment (total 2 <= MaxFiles 2) succeeds
	createdMaint, err := service.Create(ctx, owner.User.ID, createdCar.ID, CreateInput{
		Title:   "Revisão com Nota Fiscal",
		Date:    time.Now().UTC(),
		Mileage: 15000,
		Attachments: []domainmaintenance.Attachment{
			{
				ID:       "att-1",
				Name:     "nota.png",
				MimeType: "image/png",
				DataUrl:  sampleDataURL,
			},
		},
	})
	if err != nil {
		t.Fatalf("expected maintenance creation to succeed within quota, got %v", err)
	}
	if len(createdMaint.Attachments) != 1 || createdMaint.Attachments[0].DataUrl != sampleDataURL {
		t.Fatalf("expected decrypted attachment DataUrl on return, got %+v", createdMaint.Attachments)
	}

	// Verify stored attachment in repository is encrypted at rest
	rawMaint, _ := maintRepo.GetByID(ctx, createdMaint.ID)
	if !filecipher.IsEncrypted(rawMaint.Attachments[0].DataUrl) {
		t.Fatalf("expected attachment DataUrl to be encrypted at rest, got %q", rawMaint.Attachments[0].DataUrl)
	}

	// Adding another maintenance with 1 attachment exceeds MaxFiles (2 + 1 = 3 > 2)
	_, err = service.Create(ctx, owner.User.ID, createdCar.ID, CreateInput{
		Title:   "Excedendo Cota",
		Date:    time.Now().UTC(),
		Mileage: 16000,
		Attachments: []domainmaintenance.Attachment{
			{
				ID:       "att-2",
				Name:     "nota2.png",
				MimeType: "image/png",
				DataUrl:  sampleDataURL,
			},
		},
	})
	if err != ErrStorageQuotaExceeded {
		t.Fatalf("expected ErrStorageQuotaExceeded when exceeding MaxFiles, got %v", err)
	}

	// Verify persisted counters on the User record
	persistedUser, _ := users.FindByID(ctx, owner.User.ID)
	if !persistedUser.StorageQuotaSynced || persistedUser.FilesCount != 2 || persistedUser.StorageUsedBytes <= 0 {
		t.Fatalf("expected persisted FilesCount=2 and StorageQuotaSynced=true, got %+v", persistedUser)
	}

	// Deleting the maintenance frees 1 file in the persisted quota
	if err := service.Delete(ctx, owner.User.ID, createdMaint.ID); err != nil {
		t.Fatalf("delete maintenance failed: %v", err)
	}
	persistedAfterDelete, _ := users.FindByID(ctx, owner.User.ID)
	if persistedAfterDelete.FilesCount != 1 {
		t.Fatalf("expected persisted FilesCount=1 after deleting maintenance, got %d", persistedAfterDelete.FilesCount)
	}
}

func TestLegacyFilesLazyMigrationAndPersistedQuota(t *testing.T) {
	users := memoryuser.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := jwtauth.NewProvider("access-secret-32-chars-long-security", "refresh-secret-32-chars-long-security")
	authService := auth.NewService(users, hasher, tokens, nil)
	ctx := context.Background()

	owner, _ := authService.Register(ctx, auth.RegisterInput{Name: "Legacy User", Email: "legacy@example.com", Password: "secret123"})
	owner.User.EmailVerified = true
	owner.User.StorageQuotaSynced = false
	_, _ = users.Update(ctx, owner.User)

	cars := carrepo.NewRepository()
	maintRepo := maintrepo.NewRepository()

	sampleDataURL := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="

	// Pre-populate legacy unencrypted car and maintenance directly in repositories
	legacyCar, _ := cars.Create(ctx, domaincar.Car{
		OwnerID:         owner.User.ID,
		Name:            "Carro Antigo",
		Manufacturer:    "VW",
		Model:           "Gol",
		YearManufacture: 2015,
		YearModel:       2016,
		ImageUrl:        sampleDataURL,
	})
	legacyMaint, _ := maintRepo.Create(ctx, domainmaintenance.Maintenance{
		CarID:   legacyCar.ID,
		Title:   "Manutenção Antiga",
		Date:    time.Now().UTC(),
		Mileage: 50000,
		Attachments: []domainmaintenance.Attachment{
			{
				ID:       "legacy-att-1",
				Name:     "recibo_antigo.png",
				MimeType: "image/png",
				DataUrl:  sampleDataURL,
			},
		},
	})

	fc := filecipher.New("master-file-encryption-key-with-32-chars!!")
	carService := car.NewService(cars, users)
	carService.SetSecurityDependencies(maintRepo, fc)

	// Calling List triggers one-time lazy sync and in-place encryption of legacy files
	listedCars, err := carService.List(ctx, owner.User.ID)
	if err != nil || len(listedCars) != 1 {
		t.Fatalf("expected 1 listed car, got %v (err: %v)", len(listedCars), err)
	}
	if listedCars[0].ImageUrl != sampleDataURL {
		t.Fatalf("expected decrypted ImageUrl on List, got %q", listedCars[0].ImageUrl)
	}

	// Verify legacy car and maintenance were encrypted in-place in the repositories
	storedCar, _ := cars.GetByID(ctx, legacyCar.ID)
	if !filecipher.IsEncrypted(storedCar.ImageUrl) {
		t.Fatalf("expected legacy car ImageUrl to be encrypted in-place, got %q", storedCar.ImageUrl)
	}
	storedMaint, _ := maintRepo.GetByID(ctx, legacyMaint.ID)
	if !filecipher.IsEncrypted(storedMaint.Attachments[0].DataUrl) {
		t.Fatalf("expected legacy attachment DataUrl to be encrypted in-place, got %q", storedMaint.Attachments[0].DataUrl)
	}

	// Verify user record now has StorageQuotaSynced=true and persisted FilesCount=2
	syncedUser, _ := users.FindByID(ctx, owner.User.ID)
	if !syncedUser.StorageQuotaSynced || syncedUser.FilesCount != 2 || syncedUser.StorageUsedBytes <= 0 {
		t.Fatalf("expected synced user with FilesCount=2, got %+v", syncedUser)
	}
}
