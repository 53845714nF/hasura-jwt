package middleware

import (
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestKeyRateLimiter(t *testing.T) {
	// Limiter: 1 req / second, burst 2
	limiter := NewKeyRateLimiter(rate.Every(1*time.Second), 2)

	email1 := "user1@example.com"
	email2 := "user2@example.com"

	// First 2 attempts for email1 should succeed
	if !limiter.Allow(email1) {
		t.Errorf("Expected attempt 1 for email1 to be allowed")
	}
	if !limiter.Allow(email1) {
		t.Errorf("Expected attempt 2 for email1 to be allowed")
	}

	// 3rd attempt for email1 should be blocked
	if limiter.Allow(email1) {
		t.Errorf("Expected attempt 3 for email1 to be blocked")
	}

	// email2 should still be allowed (independent rate limit)
	if !limiter.Allow(email2) {
		t.Errorf("Expected attempt 1 for email2 to be allowed")
	}
}
