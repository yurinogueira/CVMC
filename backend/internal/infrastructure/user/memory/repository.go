package memory

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	domainuser "cvmc/internal/domain/user"
)

var ErrNotFound = errors.New("user not found")

type Repository struct {
	mu     sync.RWMutex
	byID   map[string]domainuser.User
	byMail map[string]string
	nextID int64
}

func NewRepository() *Repository {
	return &Repository{byID: make(map[string]domainuser.User), byMail: make(map[string]string)}
}

func (r *Repository) Create(ctx context.Context, user domainuser.User) (domainuser.User, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextID++
	user.ID = r.makeID()
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now().UTC()
	}
	if user.MaxVehicles <= 0 {
		user.MaxVehicles = domainuser.DefaultMaxVehicles
	}
	if user.MaxFiles <= 0 {
		user.MaxFiles = domainuser.DefaultMaxFiles
	}
	if user.MaxStorageBytes <= 0 {
		user.MaxStorageBytes = domainuser.DefaultMaxStorageBytes
	}
	r.byID[user.ID] = user
	r.byMail[user.Email] = user.ID
	return user, nil
}

func (r *Repository) Update(ctx context.Context, user domainuser.User) (domainuser.User, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	old, ok := r.byID[user.ID]
	if !ok {
		return domainuser.User{}, ErrNotFound
	}

	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	if user.MaxVehicles <= 0 {
		user.MaxVehicles = domainuser.DefaultMaxVehicles
	}
	if user.MaxFiles <= 0 {
		user.MaxFiles = domainuser.DefaultMaxFiles
	}
	if user.MaxStorageBytes <= 0 {
		user.MaxStorageBytes = domainuser.DefaultMaxStorageBytes
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = time.Now().UTC()
	}

	if old.Email != user.Email {
		delete(r.byMail, old.Email)
		r.byMail[user.Email] = user.ID
	}

	r.byID[user.ID] = user
	return user, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.byID[id]
	if !ok {
		return ErrNotFound
	}
	delete(r.byMail, user.Email)
	delete(r.byID, id)
	return nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (domainuser.User, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.byMail[strings.ToLower(strings.TrimSpace(email))]
	if !ok {
		return domainuser.User{}, ErrNotFound
	}
	user, ok := r.byID[id]
	if !ok {
		return domainuser.User{}, ErrNotFound
	}
	if user.MaxVehicles <= 0 {
		user.MaxVehicles = 3
	}
	return user, nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (domainuser.User, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.byID[id]
	if !ok {
		return domainuser.User{}, ErrNotFound
	}
	if user.MaxVehicles <= 0 {
		user.MaxVehicles = 3
	}
	return user, nil
}

func (r *Repository) FindByGoogleID(ctx context.Context, googleID string) (domainuser.User, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()

	cleanID := strings.TrimSpace(googleID)
	if cleanID == "" {
		return domainuser.User{}, ErrNotFound
	}

	for _, u := range r.byID {
		if u.GoogleID == cleanID {
			if u.MaxVehicles <= 0 {
				u.MaxVehicles = 3
			}
			return u, nil
		}
	}
	return domainuser.User{}, ErrNotFound
}

func (r *Repository) FindByEmailVerificationTokenHash(ctx context.Context, hash string) (domainuser.User, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()

	cleanHash := strings.TrimSpace(hash)
	if cleanHash == "" {
		return domainuser.User{}, ErrNotFound
	}

	for _, u := range r.byID {
		if u.EmailVerificationTokenHash == cleanHash {
			if u.MaxVehicles <= 0 {
				u.MaxVehicles = 3
			}
			return u, nil
		}
	}
	return domainuser.User{}, ErrNotFound
}

func (r *Repository) FindByPasswordResetTokenHash(ctx context.Context, hash string) (domainuser.User, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()

	cleanHash := strings.TrimSpace(hash)
	if cleanHash == "" {
		return domainuser.User{}, ErrNotFound
	}

	for _, u := range r.byID {
		if u.PasswordResetTokenHash == cleanHash {
			if u.MaxVehicles <= 0 {
				u.MaxVehicles = 3
			}
			return u, nil
		}
	}
	return domainuser.User{}, ErrNotFound
}

func (r *Repository) makeID() string {
	return fmt.Sprintf("user-%d", r.nextID)
}
