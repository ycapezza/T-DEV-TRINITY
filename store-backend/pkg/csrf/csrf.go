package csrf

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"
)

type TokenStore struct {
    tokens map[string]time.Time
    mu     sync.RWMutex
}

func NewTokenStore() *TokenStore {
    return &TokenStore{
        tokens: make(map[string]time.Time),
    }
}

func (s *TokenStore) Generate() (string, error) {
    b := make([]byte, 32)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }

    token := base64.URLEncoding.EncodeToString(b)
    
    s.mu.Lock()
    s.tokens[token] = time.Now().Add(24 * time.Hour)
    s.mu.Unlock()

    return token, nil
}

func (s *TokenStore) Verify(token string) bool {
    s.mu.RLock()
    expiry, exists := s.tokens[token]
    s.mu.RUnlock()

    if !exists || time.Now().After(expiry) {
        return false
    }

    return true
}