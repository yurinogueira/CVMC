package notification

import (
	"context"
	"errors"
	"strings"
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

func TestService_CheckMaintenanceAlerts_ConsolidatesMultiplePendingItemsAndCooldown(t *testing.T) {
	ctx := context.Background()
	notifRepo := notifMemory.NewRepository()
	carRepo := carMemory.NewRepository()
	maintRepo := maintMemory.NewRepository()
	sender := &mockPushSender{}

	svc := NewService(notifRepo, sender, carRepo, maintRepo)
	fixedNow := time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return fixedNow }

	_, _ = svc.RegisterDeviceToken(ctx, "user-1", "valid-token", "android")

	car, err := carRepo.Create(ctx, domaincar.Car{
		OwnerID:         "user-1",
		Name:            "Honda Civic",
		Manufacturer:    "Honda",
		Model:           "Civic",
		YearManufacture: 2022,
		YearModel:       2023,
		InitialMileage:  50000,
		LastMileage:     50000,
		CreatedAt:       fixedNow.AddDate(0, -2, 0),
	})
	if err != nil {
		t.Fatalf("failed to create car: %v", err)
	}

	// Oil and Air Filter changed at 40,200 km -> next due 50,200 km -> 200 km remaining (<= 500 km threshold!)
	_, err = maintRepo.Create(ctx, domainmaintenance.Maintenance{
		CarID:   car.ID,
		Title:   "Revisão de Óleo e Ar",
		Date:    fixedNow.AddDate(0, -6, 0),
		Mileage: 40200,
		Types: []string{
			domainmaintenance.TypeEngineOil,
			domainmaintenance.TypeAirFilter,
		},
	})
	if err != nil {
		t.Fatalf("failed to create maintenance: %v", err)
	}

	// 1. First execution: should consolidate BOTH items into a SINGLE push notification
	alerts, err := svc.CheckMaintenanceAlerts(ctx, "user-1")
	if err != nil {
		t.Fatalf("unexpected error checking maintenance alerts: %v", err)
	}
	if alerts != 1 {
		t.Fatalf("expected 1 consolidated alert triggered, got %d", alerts)
	}
	if len(sender.sent) != 1 {
		t.Fatalf("expected exactly 1 push sent, got %d", len(sender.sent))
	}

	msg := sender.sent[0]
	if msg.Title != "Revisão Preventiva: Honda Civic" {
		t.Errorf("unexpected title: %q", msg.Title)
	}
	if !strings.Contains(msg.Body, "Troca de óleo, Troca de filtro de ar") {
		t.Errorf("expected consolidated body with both items, got %q", msg.Body)
	}
	if msg.Data["route"] != "/maintenance" || msg.Data["carId"] != car.ID {
		t.Errorf("expected route=/maintenance and carId=%s, got %+v", car.ID, msg.Data)
	}

	// 2. Second execution 1 day later (within 7-day cooldown, no new items): should NOT send spam
	svc.now = func() time.Time { return fixedNow.Add(24 * time.Hour) }
	alertsSecond, err := svc.RunScheduledMaintenanceAlerts(ctx)
	if err != nil {
		t.Fatalf("unexpected error on scheduled run: %v", err)
	}
	if alertsSecond != 0 {
		t.Errorf("expected 0 alerts during 7-day cooldown, got %d", alertsSecond)
	}

	// 3. New item becomes critical (Oil Filter recorded at 40,000 km -> overdue): should bypass cooldown!
	_, _ = maintRepo.Create(ctx, domainmaintenance.Maintenance{
		CarID:   car.ID,
		Title:   "Filtro de Óleo",
		Date:    fixedNow.AddDate(0, -8, 0),
		Mileage: 40000,
		Types:   []string{domainmaintenance.TypeOilFilter},
	})
	alertsThird, err := svc.RunScheduledMaintenanceAlerts(ctx)
	if err != nil {
		t.Fatalf("unexpected error on third run: %v", err)
	}
	if alertsThird != 1 {
		t.Errorf("expected 1 alert when a new item becomes critical, got %d", alertsThird)
	}

	// 4. After 8 days (> 7-day cooldown): should send reminder again
	svc.now = func() time.Time { return fixedNow.Add(9 * 24 * time.Hour) }
	alertsFourth, err := svc.RunScheduledMaintenanceAlerts(ctx)
	if err != nil {
		t.Fatalf("unexpected error after cooldown: %v", err)
	}
	if alertsFourth != 1 {
		t.Errorf("expected 1 alert after 7-day cooldown expires, got %d", alertsFourth)
	}
}
