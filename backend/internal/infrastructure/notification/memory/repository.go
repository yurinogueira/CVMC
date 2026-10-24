package memory

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	notificationport "cvmc/internal/application/ports/notification"
	domainnotification "cvmc/internal/domain/notification"
)

type Repository struct {
	mu      sync.RWMutex
	byToken map[string]domainnotification.DeviceToken
	counter int
}

func NewRepository() *Repository {
	return &Repository{
		byToken: make(map[string]domainnotification.DeviceToken),
	}
}

func (r *Repository) Upsert(_ context.Context, dt domainnotification.DeviceToken) (domainnotification.DeviceToken, error) {
	cleanUserID := strings.TrimSpace(dt.UserID)
	cleanToken := strings.TrimSpace(dt.Token)
	if cleanUserID == "" || cleanToken == "" || len(cleanToken) > domainnotification.MaxTokenLength {
		return domainnotification.DeviceToken{}, errors.New("invalid device token")
	}

	platform, ok := domainnotification.NormalizePlatform(dt.Platform)
	if !ok {
		return domainnotification.DeviceToken{}, errors.New("invalid platform")
	}

	now := dt.UpdatedAt
	if now.IsZero() {
		now = time.Now().UTC()
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.byToken[cleanToken]
	id := existing.ID
	if !exists || id == "" {
		r.counter++
		id = fmt.Sprintf("dt-%d", r.counter)
	}

	saved := domainnotification.DeviceToken{
		ID:        id,
		UserID:    cleanUserID,
		Token:     cleanToken,
		Platform:  platform,
		UpdatedAt: now,
	}
	r.byToken[cleanToken] = saved
	return saved, nil
}

func (r *Repository) ListByUser(_ context.Context, userID string) ([]domainnotification.DeviceToken, error) {
	cleanUserID := strings.TrimSpace(userID)
	r.mu.RLock()
	defer r.mu.RUnlock()

	var items []domainnotification.DeviceToken
	for _, dt := range r.byToken {
		if dt.UserID == cleanUserID {
			items = append(items, dt)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].UpdatedAt.After(items[j].UpdatedAt)
	})
	return items, nil
}

func (r *Repository) ListDistinctUserIDs(_ context.Context) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	seen := make(map[string]struct{})
	var userIDs []string
	for _, dt := range r.byToken {
		if dt.UserID == "" {
			continue
		}
		if _, exists := seen[dt.UserID]; !exists {
			seen[dt.UserID] = struct{}{}
			userIDs = append(userIDs, dt.UserID)
		}
	}
	sort.Strings(userIDs)
	return userIDs, nil
}

func (r *Repository) DeleteByUserAndToken(_ context.Context, userID, token string) error {
	cleanUserID := strings.TrimSpace(userID)
	cleanToken := strings.TrimSpace(token)

	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.byToken[cleanToken]
	if !exists || existing.UserID != cleanUserID {
		return notificationport.ErrNotFound
	}
	delete(r.byToken, cleanToken)
	return nil
}

func (r *Repository) DeleteByToken(_ context.Context, token string) error {
	cleanToken := strings.TrimSpace(token)

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byToken[cleanToken]; !exists {
		return notificationport.ErrNotFound
	}
	delete(r.byToken, cleanToken)
	return nil
}
