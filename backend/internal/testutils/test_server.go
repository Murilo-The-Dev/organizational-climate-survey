package testutils

import (
	"net/http"
	"time"

	appMiddleware "organizational-climate-survey/backend/internal/application/middleware"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

func SetupTestRouter(jwtSecret string) *mux.Router {
	appMiddleware.ConfigureSecurityMiddleware("http://localhost:3000", 1<<20)

	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()

	publicRoutes := api.PathPrefix("").Subrouter()
	publicRoutes.Use(appMiddleware.PublicMiddlewares())
	publicRoutes.HandleFunc("/public/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pong"))
	}).Methods(http.MethodGet)

	protectedRoutes := api.PathPrefix("").Subrouter()
	protectedRoutes.Use(appMiddleware.AuthenticatedMiddlewares([]byte(jwtSecret)))
	protectedRoutes.HandleFunc("/protected/ping", func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value("user_admin_id") == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}).Methods(http.MethodGet)
	protectedRoutes.HandleFunc("/protected/admin-only", func(w http.ResponseWriter, r *http.Request) {
		empresaID, _ := r.Context().Value("empresa_id").(int)
		if empresaID != 999 {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte("forbidden"))
			return
		}
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	return router
}

func GenerateTestToken(secret string, userID, empresaID int, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := appMiddleware.JWTClaims{
		UserID:    userID,
		EmpresaID: empresaID,
		Email:     "admin@test.com",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "organizational-climate-survey-tests",
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
