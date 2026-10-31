package notification

import (
	"context"
	"testing"
	"time"

	domaincar "cvmc/internal/domain/car"
	carMemory "cvmc/internal/infrastructure/car/memory"
	notifMemory "cvmc/internal/infrastructure/notification/memory"
)

func TestIsMileageStale(t *testing.T) {
	now := time.Date(2026, 10, 30, 12, 0, 0, 0, time.UTC)
	cases := []struct {
		name string
		car  domaincar.Car
		want bool
	}{
		{"recent update", domaincar.Car{MileageUpdatedAt: now.Add(-2 * 24 * time.Hour)}, false},
		{"exactly 3 days", domaincar.Car{MileageUpdatedAt: now.Add(-MileageStaleThreshold)}, false},
		{"stale update", domaincar.Car{MileageUpdatedAt: now.Add(-4 * 24 * time.Hour)}, true},
		{"legacy fallback to updatedAt", domaincar.Car{UpdatedAt: now.Add(-10 * 24 * time.Hour)}, true},
		{"legacy fallback to createdAt", domaincar.Car{CreatedAt: now.Add(-1 * time.Hour)}, false},
		{"no timestamps", domaincar.Car{}, false},
	}
	for _, tc := range cases {
		if got := IsMileageStale(tc.car, now); got != tc.want {
			t.Errorf("%s: IsMileageStale = %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestService_RunMileageReminderAlerts_SendsAndRespectsCooldown(t *testing.T) {
	ctx := context.Background()
	tokens := notifMemory.NewRepository()
	cars := carMemory.NewRepository()
	sender := &mockPushSender{}
	svc := NewService(tokens, sender, cars, nil)

	now := time.Date(2026, 10, 30, 12, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return now }

	staleCar, _ := cars.Create(ctx, domaincar.Car{
		OwnerID:          "user-1",
		Name:             "Civic <b>Touring</b>",
		Model:            "Civic",
		MileageUpdatedAt: now.Add(-5 * 24 * time.Hour),
		CreatedAt:        now.Add(-30 * 24 * time.Hour),
	})
	_, _ = cars.Create(ctx, domaincar.Car{
		OwnerID:          "user-1",
		Name:             "Fresh",
		MileageUpdatedAt: now.Add(-1 * 24 * time.Hour),
		CreatedAt:        now.Add(-30 * 24 * time.Hour),
	})
	// User without devices should be ignored.
	_, _ = cars.Create(ctx, domaincar.Car{
		OwnerID:          "user-2",
		Name:             "Ignored",
		MileageUpdatedAt: now.Add(-20 * 24 * time.Hour),
	})

	if _, err := svc.RegisterDeviceToken(ctx, "user-1", "token-1", "android"); err != nil {
		t.Fatalf("register token: %v", err)
	}

	sent, err := svc.RunMileageReminderAlerts(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sent != 1 || len(sender.sent) != 1 {
		t.Fatalf("expected 1 reminder, got sent=%d pushes=%d", sent, len(sender.sent))
	}

	msg := sender.sent[0]
	if msg.Title != "Lembrete de Quilometragem 🚗" {
		t.Errorf("unexpected title %q", msg.Title)
	}
	wantBody := "Faz mais de 3 dias que a quilometragem do seu Civic Touring não é atualizada. Mantenha os dados em dia para acompanhar suas revisões!"
	if msg.Body != wantBody {
		t.Errorf("unexpected (or unsanitized) body %q", msg.Body)
	}
	if msg.Data["type"] != "mileage_reminder" || msg.Data["carId"] != staleCar.ID || msg.Data["route"] != "/cars/"+staleCar.ID {
		t.Errorf("unexpected data payload %#v", msg.Data)
	}

	// Within cooldown: no duplicate reminder.
	now = now.Add(2 * 24 * time.Hour)
	if sent, _ := svc.RunMileageReminderAlerts(ctx); sent != 0 {
		t.Fatalf("expected cooldown to suppress reminder, got %d", sent)
	}

	// After cooldown: the stale car is reminded again and the previously fresh car is now stale too.
	now = now.Add(MileageReminderCooldown)
	if sent, _ := svc.RunMileageReminderAlerts(ctx); sent != 2 {
		t.Fatalf("expected 2 reminders after cooldown, got %d", sent)
	}
}

func TestService_RunMileageReminderAlerts_NilDependencies(t *testing.T) {
	svc := NewService(nil, nil, nil, nil)
	if sent, err := svc.RunMileageReminderAlerts(context.Background()); sent != 0 || err != nil {
		t.Fatalf("expected no-op, got sent=%d err=%v", sent, err)
	}
}
