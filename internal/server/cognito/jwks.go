// Package cognito provides JWT validation against an AWS Cognito User Pool.
package cognito

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"
)

// JWKSCache fetches and caches the Cognito User Pool's JWKS keys.
// Keys are keyed by kid (key ID) and cached for a configurable TTL.
// On unknown kid, the cache is re-fetched once to handle key rotation.
type JWKSCache struct {
	mu         sync.RWMutex
	keys       map[string]*rsa.PublicKey
	fetchedAt  time.Time
	ttl        time.Duration
	jwksURL    string
	httpClient *http.Client
}

// NewJWKSCache creates a cache that fetches keys from jwksURL with the given TTL.
func NewJWKSCache(jwksURL string, ttl time.Duration) *JWKSCache {
	return &JWKSCache{
		jwksURL:    jwksURL,
		ttl:        ttl,
		keys:       make(map[string]*rsa.PublicKey),
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Get returns the RSA public key for kid.
// Re-fetches if the cache is stale or the kid is not found.
func (c *JWKSCache) Get(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	// Fast path: read lock.
	c.mu.RLock()
	key, ok := c.keys[kid]
	fresh := time.Since(c.fetchedAt) <= c.ttl
	c.mu.RUnlock()

	if ok && fresh {
		return key, nil
	}

	// Slow path: write lock, re-fetch if still needed.
	c.mu.Lock()
	defer c.mu.Unlock()

	// Re-check after acquiring the write lock (another goroutine may have fetched).
	key, ok = c.keys[kid]
	fresh = time.Since(c.fetchedAt) <= c.ttl
	if ok && fresh {
		return key, nil
	}

	if err := c.fetchLocked(ctx); err != nil {
		return nil, fmt.Errorf("fetch JWKS: %w", err)
	}

	key, ok = c.keys[kid]
	if !ok {
		return nil, fmt.Errorf("key %q not found in JWKS after re-fetch", kid)
	}
	return key, nil
}

// fetchLocked downloads the JWKS and rebuilds the key map.
// Must be called with c.mu held for writing.
func (c *JWKSCache) fetchLocked(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.jwksURL, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http get: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("JWKS endpoint returned %d", resp.StatusCode)
	}

	var payload struct {
		Keys []struct {
			Kid string `json:"kid"`
			Kty string `json:"kty"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return fmt.Errorf("decode JWKS: %w", err)
	}

	newKeys := make(map[string]*rsa.PublicKey, len(payload.Keys))
	for _, k := range payload.Keys {
		if k.Kty != "RSA" {
			continue
		}
		pub, err := parseRSAPublicKey(k.N, k.E)
		if err != nil {
			return fmt.Errorf("parse key %q: %w", k.Kid, err)
		}
		newKeys[k.Kid] = pub
	}

	c.keys = newKeys
	c.fetchedAt = time.Now()
	return nil
}

// parseRSAPublicKey constructs an *rsa.PublicKey from base64url-encoded n and e.
func parseRSAPublicKey(n, e string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(n)
	if err != nil {
		return nil, fmt.Errorf("decode n: %w", err)
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(e)
	if err != nil {
		return nil, fmt.Errorf("decode e: %w", err)
	}
	nInt := new(big.Int).SetBytes(nBytes)
	eInt := int(new(big.Int).SetBytes(eBytes).Int64())
	return &rsa.PublicKey{N: nInt, E: eInt}, nil
}
