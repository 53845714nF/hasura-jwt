package middleware

import (
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type limiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// KeyRateLimiter rate limits operations based on arbitrary keys (e.g. Email address or Client IP)
type KeyRateLimiter struct {
	mu      sync.Mutex
	entries map[string]*limiterEntry
	r       rate.Limit
	b       int
}

// NewKeyRateLimiter creates a key-based rate limiter.
// r: rate limit (e.g. rate.Every(3 * time.Second))
// b: burst size (e.g. 5 attempts)
func NewKeyRateLimiter(rateLimit rate.Limit, burstSize int) *KeyRateLimiter {
	limiter := &KeyRateLimiter{
		entries: make(map[string]*limiterEntry),
		r:       rateLimit,
		b:       burstSize,
	}

	// Clean up stale entries periodically
	go func() {
		for {
			time.Sleep(2 * time.Minute)
			limiter.mu.Lock()
			for key, entry := range limiter.entries {
				if time.Since(entry.lastSeen) > 10*time.Minute {
					delete(limiter.entries, key)
				}
			}
			limiter.mu.Unlock()
		}
	}()

	return limiter
}

// Allow checks if the given key is allowed to proceed under the rate limit.
func (k *KeyRateLimiter) Allow(key string) bool {
	if key == "" {
		return true
	}
	key = strings.ToLower(strings.TrimSpace(key))

	k.mu.Lock()
	defer k.mu.Unlock()

	entry, exists := k.entries[key]
	if !exists {
		limiter := rate.NewLimiter(k.r, k.b)
		k.entries[key] = &limiterEntry{limiter: limiter, lastSeen: time.Now()}
		return limiter.Allow()
	}

	entry.lastSeen = time.Now()
	return entry.limiter.Allow()
}
