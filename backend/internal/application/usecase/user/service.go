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
	ErrGoogleEmailMismatch    = errors.New("google account email does not match profile email")
	ErrGoogleAlreadyLinked    = errors.New("google account is already linked to another user")
	ErrGoogleNotLinked        = errors.New("google account is not linked")
	ErrLockoutPrevented       = errors.New("cannot unlink google account without setting a password first")
	ErrPasswordAlreadySet     = errors.New("password already set, use update password instead")
	ErrInvalidGoogleToken     = errors.New("invalid google token")
)

type AuthProvidersOutput struct {
	HasPassword  bool   `json:"hasPassword"`
	GoogleLinked bool   `json:"googleLinked"`
	GoogleEmail  string `json:"googleEmail,omitempty"`
}

type Service struct {
	users           userport.Repository
	cars            carport.Repository
	hasher          portauth.PasswordHasher
	googleValidator portauth.GoogleTokenValidator
	now             func() time.Time
}

type ProfileOutput struct {
	User          domainuser.User `json:"user"`
	VehiclesCount int             `json:"vehiclesCount"`
	MaxVehicles   int             `json:"maxVehicles"`
}

func NewService(users userport.Repository, cars carport.Repository, hasher portauth.PasswordHasher, googleValidator ...portauth.GoogleTokenValidator) *Service {
	s := &Service{
		users:  users,
		cars:   cars,
		hasher: hasher,
		now:    time.Now,
	}
	if len(googleValidator) > 0 {
		s.googleValidator = googleValidator[0]
	}
	return s
}

func (s *Service) SetGoogleValidator(v portauth.GoogleTokenValidator) {
	s.googleValidator = v
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

func (s *Service) GetAuthProviders(ctx context.Context, userID string) (AuthProvidersOutput, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return AuthProvidersOutput{}, ErrUserNotFound
	}
	googleLinked := user.GoogleID != ""
	googleEmail := ""
	if googleLinked {
		googleEmail = user.Email
	}
	return AuthProvidersOutput{
		HasPassword:  user.PasswordHash != "",
		GoogleLinked: googleLinked,
		GoogleEmail:  googleEmail,
	}, nil
}

func (s *Service) LinkGoogleAccount(ctx context.Context, userID, idToken string) error {
	if s.googleValidator == nil {
		return errors.New("google authentication not configured")
	}
	cleanToken := strings.TrimSpace(idToken)
	if cleanToken == "" {
		return ErrInvalidInput
	}

	payload, err := s.googleValidator.Validate(ctx, cleanToken)
	if err != nil || !payload.EmailVerified {
		return ErrInvalidGoogleToken
	}

	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	if !strings.EqualFold(payload.Email, user.Email) {
		return ErrGoogleEmailMismatch
	}

	existing, err := s.users.FindByGoogleID(ctx, payload.Subject)
	if err == nil && existing.ID != userID {
		return ErrGoogleAlreadyLinked
	}

	now := s.now().UTC()
	user.GoogleID = payload.Subject
	if !user.EmailVerified {
		user.EmailVerified = true
		user.EmailVerifiedAt = &now
	}
	user.UpdatedAt = now

	_, err = s.users.Update(ctx, user)
	return err
}

func (s *Service) UnlinkGoogleAccount(ctx context.Context, userID string) error {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	if user.GoogleID == "" {
		return ErrGoogleNotLinked
	}

	if user.PasswordHash == "" {
		return ErrLockoutPrevented
	}

	user.GoogleID = ""
	user.UpdatedAt = s.now().UTC()

	_, err = s.users.Update(ctx, user)
	return err
}

func (s *Service) SetPassword(ctx context.Context, userID, newPassword string) error {
	if len(newPassword) < MinPasswordLen || len(newPassword) > MaxPasswordLen {
		return ErrWeakPassword
	}

	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	if user.PasswordHash != "" {
		return ErrPasswordAlreadySet
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
