package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cvmc/internal/shared/middleware"
)

func TestExtractClientIP_UntrustedProxy_IgnoresHeaders(t *testing.T) {
	limiter := middleware.NewRateLimiter(10, 10, "127.0.0.1", "::1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	req.RemoteAddr = "198.51.100.5:1234"
	req.Header.Set("X-Forwarded-For", "203.0.113.1, 10.0.0.1")
	req.Header.Set("X-Real-IP", "203.0.113.2")

	ip := limiter.ClientIP(req)
	if ip != "198.51.100.5" {
		t.Fatalf("expected remote IP '198.51.100.5', got %q", ip)
	}
}

func TestExtractClientIP_TrustedProxy_UsesXRealIP(t *testing.T) {
	limiter := middleware.NewRateLimiter(10, 10, "127.0.0.1", "::1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	req.RemoteAddr = "127.0.0.1:4567"
	req.Header.Set("X-Real-IP", "203.0.113.50")
	req.Header.Set("X-Forwarded-For", "198.51.100.1")

	ip := limiter.ClientIP(req)
	if ip != "203.0.113.50" {
		t.Fatalf("expected real IP '203.0.113.50', got %q", ip)
	}
}

func TestExtractClientIP_TrustedProxy_UsesXForwardedFor(t *testing.T) {
	limiter := middleware.NewRateLimiter(10, 10, "127.0.0.1", "::1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	req.RemoteAddr = "127.0.0.1:4567"
	req.Header.Set("X-Forwarded-For", " 203.0.113.77 , 198.51.100.1 ")

	ip := limiter.ClientIP(req)
	if ip != "203.0.113.77" {
		t.Fatalf("expected first valid IP '203.0.113.77', got %q", ip)
	}
}

func TestExtractClientIP_TrustedProxy_MalformedHeaders(t *testing.T) {
	limiter := middleware.NewRateLimiter(10, 10, "127.0.0.1", "::1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	req.RemoteAddr = "127.0.0.1:4567"
	req.Header.Set("X-Real-IP", "invalid-ip")
	req.Header.Set("X-Forwarded-For", "also-invalid, unknown")

	ip := limiter.ClientIP(req)
	if ip != "127.0.0.1" {
		t.Fatalf("expected fallback to proxy remote IP '127.0.0.1', got %q", ip)
	}
}

func TestExtractClientIP_CIDR_TrustedProxy(t *testing.T) {
	limiter := middleware.NewRateLimiter(10, 10, "10.0.0.0/8", "172.16.0.0/12")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	req.RemoteAddr = "10.0.5.2:54321"
	req.Header.Set("X-Real-IP", "203.0.113.99")

	ip := limiter.ClientIP(req)
	if ip != "203.0.113.99" {
		t.Fatalf("expected extracted IP '203.0.113.99', got %q", ip)
	}
}

func TestExtractClientIP_IPv6(t *testing.T) {
	limiter := middleware.NewRateLimiter(10, 10, "::1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	req.RemoteAddr = "[::1]:54321"
	req.Header.Set("X-Real-IP", "2001:db8::1")

	ip := limiter.ClientIP(req)
	if ip != "2001:db8::1" {
		t.Fatalf("expected extracted IPv6 '2001:db8::1', got %q", ip)
	}
}

func TestRateLimiter_SpoofingBlocked(t *testing.T) {
	// 2 tokens capacity, almost 0 refill rate
	limiter := middleware.NewRateLimiter(0.001, 2, "127.0.0.1")

	handler := limiter.Limit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	spoofedIPs := []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"}
	responses := make([]int, len(spoofedIPs))

	for i, spoofed := range spoofedIPs {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
		req.RemoteAddr = "198.51.100.5:1234" // Direct connection from untrusted IP
		req.Header.Set("X-Forwarded-For", spoofed)
		req.Header.Set("X-Real-IP", spoofed)

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		responses[i] = rec.Code
	}

	if responses[0] != http.StatusOK {
		t.Errorf("1st request expected 200, got %d", responses[0])
	}
	if responses[1] != http.StatusOK {
		t.Errorf("2nd request expected 200, got %d", responses[1])
	}
	if responses[2] != http.StatusTooManyRequests {
		t.Errorf("3rd request expected 429 Too Many Requests (spoofing prevented), got %d", responses[2])
	}
}

func TestRateLimiter_TrustedProxy_DistinctClients(t *testing.T) {
	limiter := middleware.NewRateLimiter(0.001, 2, "127.0.0.1")

	handler := limiter.Limit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Client A makes 2 requests
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/resource", nil)
		req.RemoteAddr = "127.0.0.1:4567" // From trusted reverse proxy
		req.Header.Set("X-Real-IP", "203.0.113.10")

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("client A request %d expected 200, got %d", i+1, rec.Code)
		}
	}

	// Client B makes 2 requests (should succeed because it's a different client IP)
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/resource", nil)
		req.RemoteAddr = "127.0.0.1:4567" // From trusted reverse proxy
		req.Header.Set("X-Real-IP", "203.0.113.20")

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("client B request %d expected 200, got %d", i+1, rec.Code)
		}
	}

	// Client A 3rd request should be blocked
	req := httptest.NewRequest(http.MethodGet, "/api/v1/resource", nil)
	req.RemoteAddr = "127.0.0.1:4567"
	req.Header.Set("X-Real-IP", "203.0.113.10")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("client A 3rd request expected 429, got %d", rec.Code)
	}
}
