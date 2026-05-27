package middleware

import (
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type revokedTokenEntry struct {
	expiresAt time.Time
}

var (
	revokedTokens       sync.Map
	revocationCheckHits uint64
)

// RevokeToken registra um token/JTI revogado até sua expiração.
func RevokeToken(tokenOrJTI string, expiresAt time.Time) {
	key := strings.TrimSpace(tokenOrJTI)
	if key == "" {
		return
	}
	revokedTokens.Store(key, revokedTokenEntry{expiresAt: expiresAt})
}

// IsTokenRevoked verifica se qualquer chave fornecida está revogada e ainda válida.
func IsTokenRevoked(keys ...string) bool {
	now := time.Now()
	for _, raw := range keys {
		key := strings.TrimSpace(raw)
		if key == "" {
			continue
		}

		value, ok := revokedTokens.Load(key)
		if !ok {
			continue
		}

		entry, ok := value.(revokedTokenEntry)
		if !ok {
			revokedTokens.Delete(key)
			continue
		}

		if !entry.expiresAt.IsZero() && now.After(entry.expiresAt) {
			revokedTokens.Delete(key)
			continue
		}

		return true
	}

	if atomic.AddUint64(&revocationCheckHits, 1)%512 == 0 {
		cleanupExpiredRevocations(now)
	}

	return false
}

func cleanupExpiredRevocations(now time.Time) {
	revokedTokens.Range(func(key, value interface{}) bool {
		entry, ok := value.(revokedTokenEntry)
		if !ok {
			revokedTokens.Delete(key)
			return true
		}

		if !entry.expiresAt.IsZero() && now.After(entry.expiresAt) {
			revokedTokens.Delete(key)
		}
		return true
	})
}

func resetRevokedTokensForTests() {
	revokedTokens.Range(func(key, value interface{}) bool {
		revokedTokens.Delete(key)
		return true
	})
	atomic.StoreUint64(&revocationCheckHits, 0)
}
