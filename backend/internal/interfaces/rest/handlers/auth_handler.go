package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	portauth "cvmc/internal/application/ports/auth"
	userport "cvmc/internal/application/ports/user"
	authusecase "cvmc/internal/application/usecase/auth"
	"cvmc/internal/shared/httpx"
)

type AuthHandler struct {
	service *authusecase.Service
}

func NewAuthHandler(users userport.Repository, hasher portauth.PasswordHasher, tokens portauth.TokenService) *AuthHandler {
	return &AuthHandler{service: authusecase.NewService(users, hasher, tokens)}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleAuthError(w, err)
		return
	}
	output, err := h.service.Register(r.Context(), authusecase.RegisterInput(input))
	if err != nil {
		handleAuthError(w, err)
		return
	}
	httpx.Created(w, map[string]any{
		"user":   map[string]any{"id": output.User.ID, "name": output.User.Name, "email": output.User.Email, "createdAt": output.User.CreatedAt},
		"tokens": map[string]string{"accessToken": output.AccessToken, "refreshToken": output.RefreshToken},
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleAuthError(w, err)
		return
	}
	output, err := h.service.Login(r.Context(), authusecase.LoginInput(input))
	if err != nil {
		handleAuthError(w, err)
		return
	}
	httpx.Success(w, map[string]any{
		"user":   map[string]any{"id": output.User.ID, "name": output.User.Name, "email": output.User.Email, "createdAt": output.User.CreatedAt},
		"tokens": map[string]string{"accessToken": output.AccessToken, "refreshToken": output.RefreshToken},
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleAuthError(w, err)
		return
	}
	output, err := h.service.Refresh(r.Context(), input.RefreshToken)
	if err != nil {
		handleAuthError(w, err)
		return
	}
	httpx.Success(w, map[string]any{
		"user":   map[string]any{"id": output.User.ID, "name": output.User.Name, "email": output.User.Email, "createdAt": output.User.CreatedAt},
		"tokens": map[string]string{"accessToken": output.AccessToken, "refreshToken": output.RefreshToken},
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	accessToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	user, err := h.service.Me(r.Context(), accessToken)
	if err != nil {
		handleAuthError(w, err)
		return
	}
	httpx.Success(w, map[string]any{"id": user.ID, "name": user.Name, "email": user.Email, "createdAt": user.CreatedAt})
}

func handleAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, authusecase.ErrEmailInUse):
		httpx.Error(w, http.StatusConflict, "Email already in use", nil)
	case errors.Is(err, authusecase.ErrInvalidCredentials):
		httpx.Error(w, http.StatusUnauthorized, "Invalid credentials", nil)
	case errors.Is(err, authusecase.ErrInvalidToken):
		httpx.Error(w, http.StatusUnauthorized, "Invalid token", nil)
	case errors.Is(err, authusecase.ErrUserNotFound):
		httpx.Error(w, http.StatusNotFound, "User not found", nil)
	default:
		httpx.Error(w, http.StatusInternalServerError, "Unexpected error", nil)
	}
}
