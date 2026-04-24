package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"organizational-climate-survey/backend/internal/testutils"
)

func TestAuthMiddlewareIntegration(t *testing.T) {
	secret := "test-secret"
	router := testutils.SetupTestRouter(secret)

	t.Run("no token returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/protected/ping", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("invalid token returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/protected/ping", nil)
		req.Header.Set("Authorization", "Bearer invalid")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("expired token returns 401", func(t *testing.T) {
		token, err := testutils.GenerateTestToken(secret, 1, 10, -time.Minute)
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}
		req := httptest.NewRequest(http.MethodGet, "/api/v1/protected/ping", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("valid token passes", func(t *testing.T) {
		token, err := testutils.GenerateTestToken(secret, 1, 10, time.Hour)
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}
		req := httptest.NewRequest(http.MethodGet, "/api/v1/protected/ping", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("valid token wrong permission returns 403", func(t *testing.T) {
		token, err := testutils.GenerateTestToken(secret, 1, 10, time.Hour)
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}
		req := httptest.NewRequest(http.MethodGet, "/api/v1/protected/admin-only", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", w.Code)
		}
	})
}
