// Package middleware fornece componentes intermediários para processamento de requisicoes.
// Implementa autenticacao, autorizacao e validacoes de seguranca da aplicacao.
package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"organizational-climate-survey/backend/internal/application/dto/response"
	"organizational-climate-survey/backend/internal/domain/repository"
	"organizational-climate-survey/backend/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

const (
	defaultRequestBodyLimitBytes int64 = 1 << 20 // 1 MiB
	defaultRateLimitWindow             = time.Minute
)

type rateLimitEntry struct {
	windowStart time.Time
	count       int
}

var (
	allowedOriginsMu sync.RWMutex
	allowedOrigins   = map[string]struct{}{"http://localhost:3000": {}}
	middlewareLog    = logger.New(nil)

	maxRequestBodyBytes int64 = defaultRequestBodyLimitBytes

	rateLimiterMu sync.Mutex
	rateLimiters  = make(map[string]rateLimitEntry)
)

// ConfigureSecurityMiddleware aplica configuracoes globais dos middlewares de seguranca.
func ConfigureSecurityMiddleware(frontendURL string, requestBodyLimitBytes int64) {
	setAllowedOrigins(frontendURL)
	if requestBodyLimitBytes > 0 {
		maxRequestBodyBytes = requestBodyLimitBytes
	}
}

func setAllowedOrigins(frontendURL string) {
	allowedOriginsMu.Lock()
	defer allowedOriginsMu.Unlock()

	allowedOrigins = map[string]struct{}{}

	originsValue := strings.TrimSpace(frontendURL)
	if originsValue == "" {
		originsValue = strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS"))
	}
	if originsValue == "" {
		originsValue = "http://localhost:3000"
	}

	for _, origin := range strings.Split(originsValue, ",") {
		clean := strings.TrimSpace(origin)
		if clean != "" {
			allowedOrigins[clean] = struct{}{}
		}
	}
}

func isOriginAllowed(origin string) bool {
	allowedOriginsMu.RLock()
	defer allowedOriginsMu.RUnlock()
	_, ok := allowedOrigins[origin]
	return ok
}

func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		parts := strings.Split(ip, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return strings.TrimSpace(ip)
	}
	return r.RemoteAddr
}

func getRateLimitForPath(path string) (int, time.Duration) {
	switch {
	case strings.HasPrefix(path, "/api/v1/auth/login"):
		return 10, defaultRateLimitWindow
	case strings.HasPrefix(path, "/api/v1/auth/forgot-password"):
		return 5, defaultRateLimitWindow
	case strings.HasPrefix(path, "/api/v1/pesquisas/") && strings.HasSuffix(path, "/token"):
		return 20, defaultRateLimitWindow
	case strings.HasPrefix(path, "/api/v1/respostas/submit"):
		return 30, defaultRateLimitWindow
	default:
		return 120, defaultRateLimitWindow
	}
}

func writeRateLimitHeaders(w http.ResponseWriter, limit int, remaining int, resetAt time.Time) {
	w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
	w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
	w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetAt.Unix(), 10))
}

func generateRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	return hex.EncodeToString(b)
}

// CORSMiddleware configura politicas de compartilhamento de recursos entre origens.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := strings.TrimSpace(r.Header.Get("Origin"))
		if origin != "" {
			if !isOriginAllowed(origin) {
				response.WriteError(w, http.StatusForbidden, "Origem nao permitida", "Origin nao autorizada para esta API")
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-Submission-Token")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, X-Request-ID, X-Response-Time")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware registra informacoes basicas de cada requisicao.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := ""
		if rid := r.Context().Value("request_id"); rid != nil {
			if id, ok := rid.(string); ok {
				requestID = id
			}
		}

		middlewareLog.WithFields(map[string]interface{}{
			"request_id":  requestID,
			"method":      r.Method,
			"path":        r.URL.Path,
			"remote_addr": r.RemoteAddr,
			"user_agent":  r.UserAgent(),
		}).Info("http request")

		next.ServeHTTP(w, r)
	})
}

// RequestIDMiddleware injeta um identificador unico por requisicao para rastreabilidade.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := strings.TrimSpace(r.Header.Get("X-Request-ID"))
		if requestID == "" {
			requestID = generateRequestID()
		}

		ctx := context.WithValue(r.Context(), "request_id", requestID)
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type responseTimerWriter struct {
	http.ResponseWriter
	start       time.Time
	wroteHeader bool
}

func (w *responseTimerWriter) WriteHeader(statusCode int) {
	if !w.wroteHeader {
		w.Header().Set("X-Response-Time", time.Since(w.start).String())
		w.wroteHeader = true
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseTimerWriter) Write(data []byte) (int, error) {
	if !w.wroteHeader {
		w.Header().Set("X-Response-Time", time.Since(w.start).String())
		w.wroteHeader = true
	}
	return w.ResponseWriter.Write(data)
}

// ResponseTimeMiddleware adiciona cabecalho com o tempo de processamento da requisicao.
func ResponseTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timedWriter := &responseTimerWriter{ResponseWriter: w, start: time.Now()}
		next.ServeHTTP(timedWriter, r)
	})
}

// JWTAuthMiddleware valida tokens JWT e injeta dados do usuario no contexto.
func JWTAuthMiddleware(jwtSecret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.WriteError(w, http.StatusUnauthorized, "Token nao fornecido", "Header Authorization e obrigatorio")
				return
			}

			tokenParts := strings.Fields(authHeader)
			if len(tokenParts) != 2 || !strings.EqualFold(tokenParts[0], "Bearer") {
				response.WriteError(w, http.StatusUnauthorized, "Formato de token invalido", "Use: Bearer <token>")
				return
			}

			tokenString := strings.TrimSpace(tokenParts[1])
			if tokenString == "" {
				response.WriteError(w, http.StatusUnauthorized, "Formato de token invalido", "Token JWT vazio")
				return
			}

			if IsTokenRevoked(tokenString) {
				response.WriteError(w, http.StatusUnauthorized, "Token invalido", "Token revogado")
				return
			}

			token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("metodo de assinatura inesperado: %v", token.Header["alg"])
				}
				return jwtSecret, nil
			})

			if err != nil {
				response.WriteError(w, http.StatusUnauthorized, "Token invalido", err.Error())
				return
			}

			if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
				if IsTokenRevoked(claims.ID, tokenString) {
					response.WriteError(w, http.StatusUnauthorized, "Token invalido", "Token revogado")
					return
				}

				ctx := context.WithValue(r.Context(), "user_admin_id", claims.UserID)
				ctx = context.WithValue(ctx, "empresa_id", claims.EmpresaID)
				ctx = context.WithValue(ctx, "user_email", claims.Email)
				ctx = context.WithValue(ctx, "jwt_token", tokenString)
				if claims.ID != "" {
					ctx = context.WithValue(ctx, "jwt_jti", claims.ID)
				}
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			response.WriteError(w, http.StatusUnauthorized, "Token invalido", "Claims invalidas")
		})
	}
}

// EmpresaAuthMiddleware valida autorizacao de acesso a recursos da empresa.
func EmpresaAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userEmpresaID := r.Context().Value("empresa_id")
		if userEmpresaID == nil {
			response.WriteError(w, http.StatusUnauthorized, "Contexto invalido", "Informacoes de empresa nao encontradas")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RateLimitMiddleware implementa controle de taxa de requisicoes em memoria por IP e rota.
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		limit, window := getRateLimitForPath(r.URL.Path)
		key := getClientIP(r) + "|" + r.URL.Path
		now := time.Now()

		rateLimiterMu.Lock()
		entry := rateLimiters[key]
		if entry.windowStart.IsZero() || now.Sub(entry.windowStart) >= window {
			entry = rateLimitEntry{windowStart: now, count: 0}
		}

		entry.count++
		rateLimiters[key] = entry

		if len(rateLimiters) > 5000 {
			for k, v := range rateLimiters {
				if now.Sub(v.windowStart) > (2 * window) {
					delete(rateLimiters, k)
				}
			}
		}

		remaining := limit - entry.count
		if remaining < 0 {
			remaining = 0
		}
		resetAt := entry.windowStart.Add(window)
		writeRateLimitHeaders(w, limit, remaining, resetAt)

		exceeded := entry.count > limit
		rateLimiterMu.Unlock()

		if exceeded {
			response.WriteError(w, http.StatusTooManyRequests, "Muitas requisicoes", "Limite temporario de requisicoes excedido")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// BodySizeLimitMiddleware limita o tamanho do corpo para reduzir risco de DOS.
func BodySizeLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch:
			r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
		}
		next.ServeHTTP(w, r)
	})
}

// ContentTypeMiddleware valida Content-Type JSON em requisicoes de escrita.
func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			contentType := r.Header.Get("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				response.WriteError(w, http.StatusBadRequest, "Content-Type invalido", "Content-Type deve ser application/json")
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// RecoveryMiddleware captura panics e retorna erro controlado.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				middlewareLog.WithFields(map[string]interface{}{
					"panic": err,
					"stack": string(debug.Stack()),
					"path":  r.URL.Path,
				}).Error("panic recovered in middleware")
				response.WriteError(w, http.StatusInternalServerError, "Erro interno do servidor", "Ocorreu um erro inesperado")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// PublicRouteMiddleware adiciona headers de seguranca para rotas publicas.
func PublicRouteMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		next.ServeHTTP(w, r)
	})
}

// ActiveSurveyMiddleware valida se pesquisa esta ativa e no periodo correto.
func ActiveSurveyMiddleware(pesquisaRepo repository.PesquisaRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			pesquisaIDStr := vars["pesquisa_id"]

			if pesquisaIDStr == "" {
				next.ServeHTTP(w, r)
				return
			}

			pesquisaID, err := strconv.Atoi(pesquisaIDStr)
			if err != nil {
				response.WriteError(w, http.StatusBadRequest, "ID invalido", "ID da pesquisa deve ser numerico")
				return
			}

			pesquisa, err := pesquisaRepo.GetByID(r.Context(), pesquisaID)
			if err != nil {
				response.WriteError(w, http.StatusNotFound, "Pesquisa nao encontrada", "A pesquisa nao existe")
				return
			}

			if pesquisa.Status != "Ativa" {
				response.WriteError(w, http.StatusBadRequest, "Pesquisa indisponivel", "Esta pesquisa nao esta aceitando respostas")
				return
			}

			now := time.Now()
			if pesquisa.DataAbertura != nil && now.Before(*pesquisa.DataAbertura) {
				response.WriteError(w, http.StatusBadRequest, "Fora do periodo", fmt.Sprintf("Pesquisa abre em: %s", pesquisa.DataAbertura.Format("02/01/2006 15:04")))
				return
			}
			if pesquisa.DataFechamento != nil && now.After(*pesquisa.DataFechamento) {
				response.WriteError(w, http.StatusBadRequest, "Periodo encerrado", fmt.Sprintf("Pesquisa encerrou em: %s", pesquisa.DataFechamento.Format("02/01/2006 15:04")))
				return
			}

			ctx := context.WithValue(r.Context(), "pesquisa", pesquisa)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ChainMiddleware compoe multiplos middlewares em ordem de execucao.
func ChainMiddleware(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}

// PublicMiddlewares retorna cadeia de middlewares para rotas publicas.
func PublicMiddlewares() func(http.Handler) http.Handler {
	return ChainMiddleware(
		RecoveryMiddleware,
		CORSMiddleware,
		RequestIDMiddleware,
		ResponseTimeMiddleware,
		LoggingMiddleware,
		RateLimitMiddleware,
		BodySizeLimitMiddleware,
		PublicRouteMiddleware,
	)
}

// AuthenticatedMiddlewares retorna cadeia de middlewares para rotas autenticadas.
func AuthenticatedMiddlewares(jwtSecret []byte) func(http.Handler) http.Handler {
	return ChainMiddleware(
		RecoveryMiddleware,
		CORSMiddleware,
		RequestIDMiddleware,
		ResponseTimeMiddleware,
		LoggingMiddleware,
		RateLimitMiddleware,
		BodySizeLimitMiddleware,
		ContentTypeMiddleware,
		JWTAuthMiddleware(jwtSecret),
		EmpresaAuthMiddleware,
	)
}

// AdminMiddlewares retorna cadeia de middlewares para rotas administrativas.
func AdminMiddlewares(jwtSecret []byte) func(http.Handler) http.Handler {
	return ChainMiddleware(
		RecoveryMiddleware,
		CORSMiddleware,
		RequestIDMiddleware,
		ResponseTimeMiddleware,
		LoggingMiddleware,
		RateLimitMiddleware,
		BodySizeLimitMiddleware,
		ContentTypeMiddleware,
		JWTAuthMiddleware(jwtSecret),
		EmpresaAuthMiddleware,
	)
}

// SurveySubmissionMiddlewares retorna cadeia de middlewares para submissao de respostas.
func SurveySubmissionMiddlewares(pesquisaRepo repository.PesquisaRepository) func(http.Handler) http.Handler {
	return ChainMiddleware(
		RecoveryMiddleware,
		CORSMiddleware,
		RequestIDMiddleware,
		ResponseTimeMiddleware,
		LoggingMiddleware,
		RateLimitMiddleware,
		BodySizeLimitMiddleware,
		ContentTypeMiddleware,
		ActiveSurveyMiddleware(pesquisaRepo),
	)
}
