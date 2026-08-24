package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"cvmc/internal/shared/httpx"
)

type bucket struct {
	tokens    float64
	lastCheck time.Time
}

// RateLimiter implements a per-IP token bucket rate limiter.
type RateLimiter struct {
	mu          sync.Mutex
	buckets     map[string]*bucket
	rate        float64 // tokens per second
	capacity    float64
	trustedIPs  []net.IP
	trustedNets []*net.IPNet
}

// NewRateLimiter creates a rate limiter that allows 'rate' requests per second
// with a burst capacity of 'capacity'. Optional trustedProxies (IPs or CIDRs) determine
// whether X-Forwarded-For and X-Real-IP headers are trusted. Defaults to loopback (127.0.0.1, ::1).
func NewRateLimiter(rate float64, capacity float64, trustedProxies ...string) *RateLimiter {
	if len(trustedProxies) == 0 {
		trustedProxies = []string{"127.0.0.1", "::1"}
	}

	trustedIPs, trustedNets := parseTrustedProxies(trustedProxies)

	rl := &RateLimiter{
		buckets:     make(map[string]*bucket),
		rate:        rate,
		capacity:    capacity,
		trustedIPs:  trustedIPs,
		trustedNets: trustedNets,
	}
	go rl.cleanup()
	return rl
}

func parseTrustedProxies(proxies []string) ([]net.IP, []*net.IPNet) {
	var ips []net.IP
	var nets []*net.IPNet

	for _, p := range proxies {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if strings.Contains(p, "/") {
			_, ipNet, err := net.ParseCIDR(p)
			if err == nil && ipNet != nil {
				nets = append(nets, ipNet)
			}
			continue
		}
		if ip := net.ParseIP(p); ip != nil {
			ips = append(ips, ip)
		}
	}

	return ips, nets
}

func isTrusted(ip net.IP, trustedIPs []net.IP, trustedNets []*net.IPNet) bool {
	if ip == nil {
		return false
	}
	for _, tip := range trustedIPs {
		if ip.Equal(tip) {
			return true
		}
	}
	for _, tnet := range trustedNets {
		if tnet.Contains(ip) {
			return true
		}
	}
	return false
}

// ExtractClientIP securely extracts the client IP address from the HTTP request.
// If the remote address is not in trustedProxies, forwarded headers (X-Forwarded-For, X-Real-IP)
// are strictly ignored to prevent spoofing and rate limit bypass.
func ExtractClientIP(r *http.Request, trustedIPs []net.IP, trustedNets []*net.IPNet) string {
	remoteHost, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteHost = r.RemoteAddr
	}

	remoteIP := net.ParseIP(remoteHost)
	if remoteIP == nil {
		return remoteHost
	}

	// If the request does not come from a trusted proxy, always use remote IP.
	if !isTrusted(remoteIP, trustedIPs, trustedNets) {
		return remoteIP.String()
	}

	// Trusted proxy: check X-Real-IP first
	if xri := strings.TrimSpace(r.Header.Get("X-Real-IP")); xri != "" {
		if parsed := net.ParseIP(xri); parsed != nil {
			return parsed.String()
		}
	}

	// Trusted proxy: check X-Forwarded-For (extract first valid IP from the client side)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		for _, part := range strings.Split(xff, ",") {
			part = strings.TrimSpace(part)
			if parsed := net.ParseIP(part); parsed != nil {
				return parsed.String()
			}
		}
	}

	return remoteIP.String()
}

// ClientIP returns the client IP for the request based on the RateLimiter's trusted proxies.
func (rl *RateLimiter) ClientIP(r *http.Request) string {
	return ExtractClientIP(r, rl.trustedIPs, rl.trustedNets)
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, exists := rl.buckets[ip]
	if !exists {
		rl.buckets[ip] = &bucket{tokens: rl.capacity - 1, lastCheck: now}
		return true
	}

	elapsed := now.Sub(b.lastCheck).Seconds()
	b.tokens += elapsed * rl.rate
	if b.tokens > rl.capacity {
		b.tokens = rl.capacity
	}
	b.lastCheck = now

	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(5 * time.Minute)
		rl.mu.Lock()
		now := time.Now()
		for ip, b := range rl.buckets {
			if now.Sub(b.lastCheck) > 10*time.Minute {
				delete(rl.buckets, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Limit returns a middleware that rate-limits requests by client IP.
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := rl.ClientIP(r)
		if !rl.allow(ip) {
			httpx.Error(w, http.StatusTooManyRequests, "Too many requests", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}
