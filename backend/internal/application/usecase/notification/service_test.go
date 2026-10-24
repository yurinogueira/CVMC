package notification

import (
	"context"
	"errors"
	"testing"
	"time"

	notificationport "cvmc/internal/application/ports/notification"
	domaincar "cvmc/internal/domain/car"
	domainmaintenance "cvmc/internal/domain/maintenance"
	domainnotification "cvmc/internal/domain/notification"
	carMemory "cvmc/internal/infrastructure/car/memory"
	maintMemory "cvmc/internal/infrastructure/maintenance/memory"
	notifMemory "cvmc/internal/infrastructure/notification/memory"
)

type mockPushSender struct {
	sent       []domainnotification.PushMessage
	sentTokens []string
	errByToken map[string]error
}

func (m *mockPushSender) Send(_ context.Context, token string, msg domainnotification.PushMessage) error {
	if m.errByToken != nil {
		if err, exists := m.errByToken[token]; exists {
			return err
		}
	}
	m.sentTokens = append(m.sentTokens, token)
	m.sent = append(m.sent, msg)
	return nil
}

func TestService_RegisterAndUnregisterDeviceToken(t *testing.T) {
	ctx := context.Background()
	repo := notifMemory.NewRepository()
	sender := &mockPushSender{}
	svc := NewService(repo, sender, nil, nil)

	dt, err := svc.RegisterDeviceToken(ctx, "user-1", "fcm-token-123", "ANDROID")
	if err != nil {
		t.Fatalf("unexpected error registering token: %v", err)
	}
	if dt.Platform != domainnotification.PlatformAndroid {
		t.Errorf("expected platform android, got %q", dt.Platform)
	}

	// Reject invalid platform
	if _, err := svc.RegisterDeviceToken(ctx, "user-1", "fcm-token-123", "symbian"); !errors.Is(err, ErrInvalidPayload) {
		t.Errorf("expected ErrInvalidPayload for invalid platform, got %v", err)
	}

	// List tokens
	list, err := repo.ListByUser(ctx, "user-1")
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 device token, got %d (err=%v)", len(list), err)
	}

	// Unregister on logout
	if err := svc.UnregisterDeviceToken(ctx, "user-1", "fcm-token-123"); err != nil {
		t.Fatalf("unexpected error unregistering token: %v", err)
	}

	// Idempotent second unregister
	if err := svc.UnregisterDeviceToken(ctx, "user-1", "fcm-token-123"); err != nil {
		t.Fatalf("expected idempotent nil error on second unregister, got: %v", err)
	}

	list, _ = repo.ListByUser(ctx, "user-1")
	if len(list) != 0 {
		t.Fatalf("expected 0 tokens after unregister, got %d", len(list))
	}
}

func TestService_SendToUser_CleansUpUnregisteredTokens(t *testing.T) {
	ctx := context.Background()
	repo := notifMemory.NewRepository()
	sender := &mockPushSender{
		errByToken: map[string]error{
			"stale-token": notificationport.ErrUnregisteredToken,
		},
	}
	svc := NewService(repo, sender, nil, nil)

	_, _ = svc.RegisterDeviceToken(ctx, "user-1", "stale-token", "android")
	_, _ = svc.RegisterDeviceToken(ctx, "user-1", "valid-token", "android")

	sentCount, err := svc.SendToUser(ctx, "user-1", domainnotification.PushMessage{
		Title: "<script>alert(1)</script>Revisão",
		Body:  "Seu veículo precisa de atenção.",
		Data:  map[string]string{"route": "/vehicles/123"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sentCount != 1 {
		t.Errorf("expected 1 delivered message, got %d", sentCount)
	}

	remaining, _ := repo.ListByUser(ctx, "user-1")
	if len(remaining) != 1 || remaining[0].Token != "valid-token" {
		t.Fatalf("expected stale-token to be cleaned up and only valid-token remaining, got %+v", remaining)
	}
	if len(sender.sent) != 1 || sender.sent[0].Title != "alert(1)Revisão" {
		t.Errorf("expected sanitized title 'alert(1)Revisão', got %+v", sender.sent)
	}
}

func TestService_CheckMaintenanceAlerts(t *testing.T) {
	ctx := context.Background()
	notifRepo := notifMemory.NewRepository()
	carRepo := carMemory.NewRepository()
	maintRepo := maintMemory.NewRepository()
	sender := &mockPushSender{}

	svc := NewService(notifRepo, sender, carRepo, maintRepo)
	fixedNow := time.Date(2026, 10, 17, 12, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return fixedNow }

	_, _ = svc.RegisterDeviceToken(ctx, "user-1", "valid-token", "android")

	car, err := carRepo.Create(ctx, domaincar.Car{
		OwnerID:         "user-1",
		Name:            "Honda Civic",
		Manufacturer:    "Honda",
		Model:           "Civic",
		YearManufacture: 2022,
		YearModel:       2023,
		LastMileage:     50000,
	})
	if err != nil {
		t.Fatalf("failed to create car: %v", err)
	}

	_, err = maintRepo.Create(ctx, domainmaintenance.Maintenance{
		CarID:   car.ID,
		Title:   "Troca de Óleo",
		Date:    fixedNow.AddDate(0, -2, 0),
		Mileage: 40000, // 10,000 km ago (>= 9,000 km threshold)
	})
	if err != nil {
		t.Fatalf("failed to create maintenance: %v", err)
	}

	alerts, err := svc.CheckMaintenanceAlerts(ctx, "user-1")
	if err != nil {
		t.Fatalf("unexpected error checking maintenance alerts: %v", err)
	}
	if alerts != 1 {
		t.Fatalf("expected 1 alert triggered, got %d", alerts)
	}
	if len(sender.sent) != 1 {
		t.Fatalf("expected 1 push sent, got %d", len(sender.sent))
	}
	if sender.sent[0].Data["route"] != "/vehicles/"+car.ID {
		t.Errorf("expected deep-link route /vehicles/%s, got %q", car.ID, sender.sent[0].Data["route"])
	}
}
