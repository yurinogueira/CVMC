package notification

import (
	"context"
	"errors"

	domainnotification "cvmc/internal/domain/notification"
)

var (
	ErrNotFound          = errors.New("device token not found")
	ErrUnregisteredToken = errors.New("device token is unregistered or invalid in FCM")
)

type DeviceTokenRepository interface {
	Upsert(ctx context.Context, dt domainnotification.DeviceToken) (domainnotification.DeviceToken, error)
	ListByUser(ctx context.Context, userID string) ([]domainnotification.DeviceToken, error)
	ListDistinctUserIDs(ctx context.Context) ([]string, error)
	DeleteByUserAndToken(ctx context.Context, userID, token string) error
	DeleteByToken(ctx context.Context, token string) error
}

type PushSender interface {
	Send(ctx context.Context, token string, msg domainnotification.PushMessage) error
}
