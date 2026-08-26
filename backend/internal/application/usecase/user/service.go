package user

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode"

	portauth "cvmc/internal/application/ports/auth"
	carport "cvmc/internal/application/ports/car"
	userport "cvmc/internal/application/ports/user"
	domainuser "cvmc/internal/domain/user"
)

const (
	MinPasswordLen = 8
	MaxPasswordLen = 72
	MinNameLen     = 2
	MaxNameLen     = 100
)

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrInvalidInput           = errors.New("invalid input")
	ErrInvalidCurrentPassword = errors.New("invalid current password")
	ErrWeakPassword           = errors.New("weak password")
)

type Service struct {
	users  userport.Repository
	cars   carport.Repository
	hasher portauth.PasswordHasher
	now    func() time.Time
}

type ProfileOutput struct {
	User          domainuser.User `json:"user"`
	VehiclesCount int             `json:"vehiclesCount"`
	MaxVehicles   int             `json:"maxVehicles"`
}

func NewService(users userport.Repository, cars carport.Repository, hasher portauth.PasswordHasher) *Service {
	return &Service{
		users:  users,
		cars:   cars,
		hasher: hasher,
		now:    time.Now,
	}
}

func isValidName(name string) bool {
	if len(name) < MinNameLen || len(name) > MaxNameLen {
		return false
	}
	for _, r := range name {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}

func (s *Service) GetProfile(ctx context.Context, userID string) (ProfileOutput, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return ProfileOutput{}, ErrUserNotFound
	}

	userCars, err := s.cars.ListByUser(ctx, userID)
	if err != nil {
		return ProfileOutput{}, err
	}

	ownedCount := 0
	for _, c := range userCars {
		if c.OwnerID == userID {
			ownedCount++
		}
	}

	maxVehicles := user.MaxVehicles
	if maxVehicles <= 0 {
		maxVehicles = 3
	}

	return ProfileOutput{
		User:          user,
		VehiclesCount: ownedCount,
		MaxVehicles:   maxVehicles,
	}, nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID, name string) (domainuser.User, error) {
	name = strings.TrimSpace(name)
	if !isValidName(name) {
		return domainuser.User{}, ErrInvalidInput
	}

	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return domainuser.User{}, ErrUserNotFound
	}

	user.Name = name
	user.UpdatedAt = s.now().UTC()

	return s.users.Update(ctx, user)
}

func (s *Service) UpdatePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	if len(newPassword) < MinPasswordLen || len(newPassword) > MaxPasswordLen {
		return ErrWeakPassword
	}

	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	if err := s.hasher.Compare(user.PasswordHash, currentPassword); err != nil {
		return ErrInvalidCurrentPassword
	}

	newHash, err := s.hasher.Hash(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = newHash
	user.UpdatedAt = s.now().UTC()

	_, err = s.users.Update(ctx, user)
	return err
}
