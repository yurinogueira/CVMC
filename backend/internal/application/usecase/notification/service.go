package notification

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	carport "cvmc/internal/application/ports/car"
	maintenanceport "cvmc/internal/application/ports/maintenance"
	notificationport "cvmc/internal/application/ports/notification"
	domainmaintenance "cvmc/internal/domain/maintenance"
	domainnotification "cvmc/internal/domain/notification"
)

const (
	PreventiveMileageThreshold = 9000
	PreventiveDaysThreshold    = 150
	UpcomingWindowDays         = 7
)

var (
	ErrInvalidPayload      = errors.New("invalid notification payload")
	ErrNoDevicesRegistered = errors.New("no active device tokens registered for user")
	ErrDeliveryFailed      = errors.New("failed to deliver push notification")
)

type Service struct {
	tokens       notificationport.DeviceTokenRepository
	sender       notificationport.PushSender
	cars         carport.Repository
	maintenances maintenanceport.Repository
	now          func() time.Time
}

func NewService(
	tokens notificationport.DeviceTokenRepository,
	sender notificationport.PushSender,
	cars carport.Repository,
	maintenances maintenanceport.Repository,
) *Service {
	return &Service{
		tokens:       tokens,
		sender:       sender,
		cars:         cars,
		maintenances: maintenances,
		now:          time.Now,
	}
}

func (s *Service) RegisterDeviceToken(ctx context.Context, userID, token, platform string) (domainnotification.DeviceToken, error) {
	cleanUserID := strings.TrimSpace(userID)
	cleanToken := strings.TrimSpace(token)
	if cleanUserID == "" || cleanToken == "" || len(cleanToken) > domainnotification.MaxTokenLength {
		return domainnotification.DeviceToken{}, ErrInvalidPayload
	}

	normPlatform, ok := domainnotification.NormalizePlatform(platform)
	if !ok {
		return domainnotification.DeviceToken{}, ErrInvalidPayload
	}

	dt := domainnotification.DeviceToken{
		UserID:    cleanUserID,
		Token:     cleanToken,
		Platform:  normPlatform,
		UpdatedAt: s.now().UTC(),
	}
	return s.tokens.Upsert(ctx, dt)
}

func (s *Service) UnregisterDeviceToken(ctx context.Context, userID, token string) error {
	cleanUserID := strings.TrimSpace(userID)
	cleanToken := strings.TrimSpace(token)
	if cleanUserID == "" || cleanToken == "" {
		return ErrInvalidPayload
	}

	err := s.tokens.DeleteByUserAndToken(ctx, cleanUserID, cleanToken)
	if errors.Is(err, notificationport.ErrNotFound) {
		// Idempotent removal on logout
		return nil
	}
	return err
}

func (s *Service) SendToUser(ctx context.Context, userID string, msg domainnotification.PushMessage) (int, error) {
	cleanUserID := strings.TrimSpace(userID)
	if cleanUserID == "" {
		return 0, ErrInvalidPayload
	}

	sanitized := domainnotification.SanitizeMessage(msg)
	if sanitized.Title == "" || sanitized.Body == "" {
		return 0, ErrInvalidPayload
	}

	devices, err := s.tokens.ListByUser(ctx, cleanUserID)
	if err != nil {
		return 0, err
	}
	if len(devices) == 0 {
		return 0, ErrNoDevicesRegistered
	}

	sentCount := 0
	var lastErr error
	for _, dt := range devices {
		if s.sender == nil {
			sentCount++
			continue
		}
		err := s.sender.Send(ctx, dt.Token, sanitized)
		if err != nil {
			if errors.Is(err, notificationport.ErrUnregisteredToken) {
				_ = s.tokens.DeleteByToken(ctx, dt.Token)
				continue
			}
			lastErr = err
			continue
		}
		sentCount++
	}

	if sentCount == 0 && lastErr != nil {
		return 0, fmt.Errorf("%w: %v", ErrDeliveryFailed, lastErr)
	}
	return sentCount, nil
}

func (s *Service) SendTestNotification(ctx context.Context, userID string) (int, error) {
	msg := domainnotification.PushMessage{
		Title: "Notificações CVMC Ativas",
		Body:  "Tudo pronto! Você receberá alertas de revisões e manutenções do seu veículo em tempo real.",
		Data: map[string]string{
			"type":  "test",
			"route": "/dashboard",
		},
	}
	return s.SendToUser(ctx, userID, msg)
}

func (s *Service) CheckMaintenanceAlerts(ctx context.Context, userID string) (int, error) {
	if s.cars == nil || s.maintenances == nil {
		return 0, nil
	}
	cars, err := s.cars.ListByUser(ctx, userID)
	if err != nil {
		return 0, err
	}

	now := s.now().UTC()
	alertsTriggered := 0

	for _, car := range cars {
		items, err := s.maintenances.ListByCar(ctx, car.ID)
		if err != nil || len(items) == 0 {
			continue
		}

		alertMsg, shouldAlert := evaluateCarMaintenanceAlert(car.ID, car.Name, car.LastMileage, items, now)
		if !shouldAlert {
			continue
		}

		sent, err := s.SendToUser(ctx, userID, alertMsg)
		if err == nil && sent > 0 {
			alertsTriggered++
		}
	}

	return alertsTriggered, nil
}

func evaluateCarMaintenanceAlert(carID, carName string, currentMileage int, maintenances []domainmaintenance.Maintenance, now time.Time) (domainnotification.PushMessage, bool) {
	var latestPast *domainmaintenance.Maintenance
	var upcoming *domainmaintenance.Maintenance

	for i := range maintenances {
		m := &maintenances[i]
		if m.Date.After(now) {
			diffDays := int(m.Date.Sub(now).Hours() / 24)
			if diffDays <= UpcomingWindowDays {
				if upcoming == nil || m.Date.Before(upcoming.Date) {
					upcoming = m
				}
			}
			continue
		}
		if latestPast == nil || m.Date.After(latestPast.Date) || (m.Date.Equal(latestPast.Date) && m.Mileage > latestPast.Mileage) {
			latestPast = m
		}
	}

	route := fmt.Sprintf("/vehicles/%s", carID)

	if upcoming != nil {
		return domainnotification.PushMessage{
			Title: fmt.Sprintf("Revisão Próxima: %s", carName),
			Body:  fmt.Sprintf("A manutenção %q está agendada para %s.", upcoming.Title, upcoming.Date.Format("02/01/2006")),
			Data: map[string]string{
				"type":          "maintenance_upcoming",
				"carId":         carID,
				"maintenanceId": upcoming.ID,
				"route":         route,
			},
		}, true
	}

	if latestPast != nil {
		kmSinceLast := currentMileage - latestPast.Mileage
		daysSinceLast := int(now.Sub(latestPast.Date).Hours() / 24)

		if kmSinceLast >= PreventiveMileageThreshold {
			return domainnotification.PushMessage{
				Title: fmt.Sprintf("Alerta de Quilometragem: %s", carName),
				Body:  fmt.Sprintf("Seu veículo rodou %d km desde a última manutenção (%s). Verifique a revisão preventiva!", kmSinceLast, latestPast.Title),
				Data: map[string]string{
					"type":  "maintenance_mileage_due",
					"carId": carID,
					"route": route,
				},
			}, true
		}

		if daysSinceLast >= PreventiveDaysThreshold {
			return domainnotification.PushMessage{
				Title: fmt.Sprintf("Lembrete de Revisão: %s", carName),
				Body:  fmt.Sprintf("Já se passaram %d dias desde a última manutenção (%s). Mantenha seu veículo em dia!", daysSinceLast, latestPast.Title),
				Data: map[string]string{
					"type":  "maintenance_date_due",
					"carId": carID,
					"route": route,
				},
			}, true
		}
	}

	return domainnotification.PushMessage{}, false
}
