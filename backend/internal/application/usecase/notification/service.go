package notification

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	carport "cvmc/internal/application/ports/car"
	maintenanceport "cvmc/internal/application/ports/maintenance"
	notificationport "cvmc/internal/application/ports/notification"
	domaincar "cvmc/internal/domain/car"
	domainmaintenance "cvmc/internal/domain/maintenance"
	domainnotification "cvmc/internal/domain/notification"
)

const (
	PreventiveMileageThreshold = 9000
	PreventiveDaysThreshold    = 150
	UpcomingWindowDays         = 7
	AlertCooldownDuration      = 7 * 24 * time.Hour
	// MileageStaleThreshold is the minimum time without mileage updates before a reminder is sent.
	MileageStaleThreshold = 3 * 24 * time.Hour
	// MileageReminderCooldown is the minimum time between consecutive mileage reminders for the same vehicle.
	MileageReminderCooldown = 7 * 24 * time.Hour
)

var (
	ErrInvalidPayload      = errors.New("invalid notification payload")
	ErrNoDevicesRegistered = errors.New("no active device tokens registered for user")
	ErrDeliveryFailed      = errors.New("failed to deliver push notification")
)

type carAlertState struct {
	lastSentAt   time.Time
	alertedTypes map[string]struct{}
}

type Service struct {
	tokens       notificationport.DeviceTokenRepository
	sender       notificationport.PushSender
	cars         carport.Repository
	maintenances maintenanceport.Repository
	now          func() time.Time
	mu           sync.RWMutex
	alertHistory map[string]carAlertState
	// mileageReminderHistory tracks the last mileage reminder sent per user:car key.
	mileageReminderHistory map[string]time.Time
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
		alertHistory: make(map[string]carAlertState),

		mileageReminderHistory: make(map[string]time.Time),
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

func (s *Service) DeleteUserDeviceTokens(ctx context.Context, userID string) error {
	if s == nil || s.tokens == nil {
		return nil
	}
	cleanUserID := strings.TrimSpace(userID)
	if cleanUserID == "" {
		return nil
	}
	devices, err := s.tokens.ListByUser(ctx, cleanUserID)
	if err != nil {
		return err
	}
	for _, dt := range devices {
		_ = s.tokens.DeleteByUserAndToken(ctx, cleanUserID, dt.Token)
	}
	return nil
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

// RunScheduledMaintenanceAlerts scans all users with active FCM tokens and dispatches consolidated alerts.
func (s *Service) RunScheduledMaintenanceAlerts(ctx context.Context) (int, error) {
	if s.tokens == nil || s.cars == nil || s.maintenances == nil {
		return 0, nil
	}

	userIDs, err := s.tokens.ListDistinctUserIDs(ctx)
	if err != nil {
		return 0, err
	}

	totalAlerts := 0
	for _, userID := range userIDs {
		if err := ctx.Err(); err != nil {
			return totalAlerts, err
		}
		triggered, err := s.CheckMaintenanceAlerts(ctx, userID)
		if err != nil {
			continue
		}
		totalAlerts += triggered
	}
	return totalAlerts, nil
}

// CheckMaintenanceAlerts evaluates all vehicles of a user and sends at most ONE consolidated notification per vehicle, respecting the 7-day cooldown.
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
		if err != nil {
			continue
		}

		alertMsg, pendingTypes, shouldAlert := evaluateConsolidatedCarAlert(car, items, now)
		if !shouldAlert {
			continue
		}

		if !s.shouldDispatchAlert(userID, car.ID, pendingTypes, now) {
			continue
		}

		sent, err := s.SendToUser(ctx, userID, alertMsg)
		if err == nil && sent > 0 {
			s.recordDispatchedAlert(userID, car.ID, pendingTypes, now)
			alertsTriggered++
		}
	}

	return alertsTriggered, nil
}

func (s *Service) shouldDispatchAlert(userID, carID string, pendingTypes []string, now time.Time) bool {
	key := userID + ":" + carID
	s.mu.RLock()
	prev, exists := s.alertHistory[key]
	s.mu.RUnlock()

	if !exists {
		return true
	}

	// If cooldown period (7 days) has elapsed, allow re-notifying
	if now.Sub(prev.lastSentAt) >= AlertCooldownDuration {
		return true
	}

	// Within cooldown: only dispatch if a NEW item became critical/eligible
	for _, t := range pendingTypes {
		if _, alreadyAlerted := prev.alertedTypes[t]; !alreadyAlerted {
			return true
		}
	}
	return false
}

func (s *Service) recordDispatchedAlert(userID, carID string, pendingTypes []string, now time.Time) {
	key := userID + ":" + carID
	typeSet := make(map[string]struct{}, len(pendingTypes))
	for _, t := range pendingTypes {
		typeSet[t] = struct{}{}
	}

	s.mu.Lock()
	if prev, exists := s.alertHistory[key]; exists && now.Sub(prev.lastSentAt) < AlertCooldownDuration {
		for k := range prev.alertedTypes {
			typeSet[k] = struct{}{}
		}
	}
	s.alertHistory[key] = carAlertState{
		lastSentAt:   now,
		alertedTypes: typeSet,
	}
	s.mu.Unlock()
}

func evaluateConsolidatedCarAlert(car domaincar.Car, maintenances []domainmaintenance.Maintenance, now time.Time) (domainnotification.PushMessage, []string, bool) {
	upcomingReport := domainmaintenance.CalculateCarUpcoming(car, maintenances, now)

	var pendingLabels []string
	var pendingTypes []string

	for _, item := range upcomingReport.Items {
		if item.AlertEligible {
			pendingLabels = append(pendingLabels, item.Label)
			pendingTypes = append(pendingTypes, item.Type)
		}
	}

	if len(pendingLabels) > 0 {
		return domainnotification.PushMessage{
			Title: fmt.Sprintf("Revisão Preventiva: %s", car.Name),
			Body:  fmt.Sprintf("Itens pendentes: %s. Agende a revisão para manter seu veículo em dia.", strings.Join(pendingLabels, ", ")),
			Data: map[string]string{
				"type":  "maintenance_preventive",
				"carId": car.ID,
				"route": "/maintenance",
			},
		}, pendingTypes, true
	}

	// Fallback for future scheduled maintenances (<= 7 days)
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
		}
	}

	if upcoming != nil {
		return domainnotification.PushMessage{
			Title: fmt.Sprintf("Revisão Preventiva: %s", car.Name),
			Body:  fmt.Sprintf("Itens pendentes: %s (%s). Agende a revisão para manter seu veículo em dia.", upcoming.Title, upcoming.Date.Format("02/01/2006")),
			Data: map[string]string{
				"type":          "maintenance_upcoming",
				"carId":         car.ID,
				"maintenanceId": upcoming.ID,
				"route":         "/maintenance",
			},
		}, []string{"scheduled:" + upcoming.ID}, true
	}

	return domainnotification.PushMessage{}, nil, false
}

// RunMileageReminderAlerts scans all users with active FCM tokens and reminds them to update the mileage
// of vehicles that have not been updated for more than MileageStaleThreshold, respecting a per-vehicle cooldown.
func (s *Service) RunMileageReminderAlerts(ctx context.Context) (int, error) {
	if s.tokens == nil || s.cars == nil {
		return 0, nil
	}

	userIDs, err := s.tokens.ListDistinctUserIDs(ctx)
	if err != nil {
		return 0, err
	}

	total := 0
	for _, userID := range userIDs {
		if err := ctx.Err(); err != nil {
			return total, err
		}
		sent, err := s.CheckMileageReminders(ctx, userID)
		if err != nil {
			continue
		}
		total += sent
	}
	return total, nil
}

// CheckMileageReminders evaluates the vehicles accessible by a user and sends one reminder per stale vehicle.
func (s *Service) CheckMileageReminders(ctx context.Context, userID string) (int, error) {
	if s.cars == nil {
		return 0, nil
	}
	cleanUserID := strings.TrimSpace(userID)
	if cleanUserID == "" {
		return 0, ErrInvalidPayload
	}
	cars, err := s.cars.ListByUser(ctx, cleanUserID)
	if err != nil {
		return 0, err
	}

	now := s.now().UTC()
	sentReminders := 0
	for _, car := range cars {
		if err := ctx.Err(); err != nil {
			return sentReminders, err
		}
		if !IsMileageStale(car, now) {
			continue
		}
		key := cleanUserID + ":" + car.ID
		if !s.shouldDispatchMileageReminder(key, now) {
			continue
		}
		sent, err := s.SendToUser(ctx, cleanUserID, buildMileageReminderMessage(car))
		if err == nil && sent > 0 {
			s.mu.Lock()
			s.mileageReminderHistory[key] = now
			s.mu.Unlock()
			sentReminders++
		}
	}
	return sentReminders, nil
}

// IsMileageStale reports whether the vehicle mileage has not been updated for more than MileageStaleThreshold.
func IsMileageStale(car domaincar.Car, now time.Time) bool {
	last := car.EffectiveMileageUpdatedAt()
	if last.IsZero() {
		return false
	}
	return now.Sub(last) > MileageStaleThreshold
}

func (s *Service) shouldDispatchMileageReminder(key string, now time.Time) bool {
	s.mu.RLock()
	lastSent, exists := s.mileageReminderHistory[key]
	s.mu.RUnlock()
	return !exists || now.Sub(lastSent) >= MileageReminderCooldown
}

func buildMileageReminderMessage(car domaincar.Car) domainnotification.PushMessage {
	label := strings.TrimSpace(car.Name)
	if label == "" {
		label = strings.TrimSpace(car.Model)
	}
	if label == "" {
		label = "veículo"
	}
	return domainnotification.SanitizeMessage(domainnotification.PushMessage{
		Title: "Lembrete de Quilometragem 🚗",
		Body:  fmt.Sprintf("Faz mais de 3 dias que a quilometragem do seu %s não é atualizada. Mantenha os dados em dia para acompanhar suas revisões!", label),
		Data: map[string]string{
			"type":  "mileage_reminder",
			"carId": car.ID,
			"route": "/cars/" + car.ID,
		},
	})
}
