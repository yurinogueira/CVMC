package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	notificationusecase "cvmc/internal/application/usecase/notification"
	domainnotification "cvmc/internal/domain/notification"
	domainuser "cvmc/internal/domain/user"
	jwtauth "cvmc/internal/infrastructure/auth/jwt"
	notifMemory "cvmc/internal/infrastructure/notification/memory"
)

type noopSender struct{}

func (noopSender) Send(_ context.Context, _ string, _ domainnotification.PushMessage) error {
	return nil
}

func TestNotificationHandler_Lifecycle(t *testing.T) {
	repo := notifMemory.NewRepository()
	svc := notificationusecase.NewService(repo, noopSender{}, nil, nil)
	tokens := jwtauth.NewProvider("test-access-secret-key-at-least-32-chars", "test-refresh-secret-key-at-least-32-chars")
	accessToken, err := tokens.GenerateAccessToken(domainuser.User{
		ID:    "user-123",
		Email: "user@cvmc.com",
	})
	if err != nil {
		t.Fatalf("failed to generate access token: %v", err)
	}

	handler := NewNotificationHandler(svc, tokens)

	// 1. Register token
	regBody, _ := json.Marshal(RegisterDeviceTokenRequest{
		Token:    "fcm-token-xyz",
		Platform: "android",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/devices/token", bytes.NewReader(regBody))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	rec := httptest.NewRecorder()
	handler.RegisterDeviceToken(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on RegisterDeviceToken, got %d: %s", rec.Code, rec.Body.String())
	}

	// 2. Send test notification
	testReq := httptest.NewRequest(http.MethodPost, "/api/v1/notifications/test", nil)
	testReq.Header.Set("Authorization", "Bearer "+accessToken)
	testRec := httptest.NewRecorder()
	handler.SendTestNotification(testRec, testReq)

	if testRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on SendTestNotification, got %d: %s", testRec.Code, testRec.Body.String())
	}

	// 3. Unregister token
	unregBody, _ := json.Marshal(UnregisterDeviceTokenRequest{
		Token: "fcm-token-xyz",
	})
	delReq := httptest.NewRequest(http.MethodDelete, "/api/v1/devices/token", bytes.NewReader(unregBody))
	delReq.Header.Set("Authorization", "Bearer "+accessToken)
	delRec := httptest.NewRecorder()
	handler.UnregisterDeviceToken(delRec, delReq)

	if delRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on UnregisterDeviceToken, got %d: %s", delRec.Code, delRec.Body.String())
	}

	// 4. Test notification when no devices are registered -> 404
	emptyReq := httptest.NewRequest(http.MethodPost, "/api/v1/notifications/test", nil)
	emptyReq.Header.Set("Authorization", "Bearer "+accessToken)
	emptyRec := httptest.NewRecorder()
	handler.SendTestNotification(emptyRec, emptyReq)

	if emptyRec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 NotFound when no devices registered, got %d", emptyRec.Code)
	}
}
