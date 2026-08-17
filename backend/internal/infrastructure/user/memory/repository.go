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
	r.byID[user.ID] = user
	r.byMail[user.Email] = user.ID
	return user, nil
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
	return user, nil
}

func (r *Repository) makeID() string {
	return fmt.Sprintf("user-%d", r.nextID)
}
