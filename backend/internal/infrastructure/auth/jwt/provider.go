package jwt

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	portauth "cvmc/internal/application/ports/auth"
	domainuser "cvmc/internal/domain/user"
)

type Provider struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

type jwtClaims struct {
	Subject   string `json:"sub"`
	Email     string `json:"email,omitempty"`
	Type      string `json:"type"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
	JTI       string `json:"jti"`
}

func NewProvider(accessSecret, refreshSecret string) *Provider {
	return NewProviderWithTTL(accessSecret, refreshSecret, 24*time.Hour, 7*24*time.Hour)
}

func NewProviderWithTTL(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) *Provider {
	return &Provider{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

func (p *Provider) GeneratePair(user domainuser.User) (portauth.TokenPair, error) {
	accessToken, err := p.GenerateAccessToken(user)
	if err != nil {
		return portauth.TokenPair{}, err
	}
	refreshToken, err := p.GenerateRefreshToken(user.ID)
	if err != nil {
		return portauth.TokenPair{}, err
	}
	return portauth.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (p *Provider) GenerateAccessToken(user domainuser.User) (string, error) {
	return p.signJWT(p.accessSecret, jwtClaims{
		Subject:   user.ID,
		Email:     user.Email,
		Type:      "access",
		IssuedAt:  time.Now().UTC().Unix(),
		ExpiresAt: time.Now().UTC().Add(p.accessTTL).Unix(),
		JTI:       randomID(),
	})
}

func (p *Provider) GenerateRefreshToken(userID string) (string, error) {
	return p.signJWT(p.refreshSecret, jwtClaims{
		Subject:   userID,
		Type:      "refresh",
		IssuedAt:  time.Now().UTC().Unix(),
		ExpiresAt: time.Now().UTC().Add(p.refreshTTL).Unix(),
		JTI:       randomID(),
	})
}

func (p *Provider) ParseAccessToken(token string) (portauth.TokenClaims, error) {
	claims, err := p.verifyJWT(p.accessSecret, token, "access")
	if err != nil {
		return portauth.TokenClaims{}, err
	}
	return portauth.TokenClaims{UserID: claims.Subject, Email: claims.Email}, nil
}

func (p *Provider) ParseRefreshToken(token string) (portauth.TokenClaims, error) {
	claims, err := p.verifyJWT(p.refreshSecret, token, "refresh")
	if err != nil {
		return portauth.TokenClaims{}, err
	}
	return portauth.TokenClaims{UserID: claims.Subject}, nil
}

func (p *Provider) signJWT(secret []byte, claims jwtClaims) (string, error) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	unsigned := header + "." + payload
	signature := signHMAC(secret, unsigned)
	return unsigned + "." + signature, nil
}

func (p *Provider) verifyJWT(secret []byte, token string, expectedType string) (jwtClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return jwtClaims{}, errors.New("invalid token")
	}
	unsigned := parts[0] + "." + parts[1]
	if !hmac.Equal([]byte(signHMAC(secret, unsigned)), []byte(parts[2])) {
		return jwtClaims{}, errors.New("invalid token")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return jwtClaims{}, err
	}
	var claims jwtClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return jwtClaims{}, err
	}
	if claims.Type != expectedType {
		return jwtClaims{}, errors.New("invalid token type")
	}
	if time.Now().UTC().Unix() > claims.ExpiresAt {
		return jwtClaims{}, errors.New("token expired")
	}
	return claims, nil
}

func signHMAC(secret []byte, unsigned string) string {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(unsigned))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func randomID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("fallback-%d", time.Now().UnixNano())
	}
	return base64.RawURLEncoding.EncodeToString(buf)
}
