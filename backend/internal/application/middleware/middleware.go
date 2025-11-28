// Package middleware fornece componentes intermediários para processamento de requisições.
// Implementa autenticação, autorização e validações de segurança da aplicação.
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"organizational-climate-survey/backend/internal/application/dto/response"

	"github.com/golang-jwt/jwt/v5"
)

// CORSMiddleware configura políticas de compartilhamento de recursos entre origens
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Configurar headers CORS para acesso cross-origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Responder requisições preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware registra informações básicas de cada requisição
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Registrar método, rota, IP e user agent
		fmt.Printf("[%s] %s %s - %s\n", 
			r.Method, 
			r.URL.Path, 
			r.RemoteAddr,
			r.UserAgent(),
		)
		
		next.ServeHTTP(w, r)
	})
}

// JWTAuthMiddleware valida tokens JWT e injeta dados do usuário no contexto
func JWTAuthMiddleware(jwtSecret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extrair token do header Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.WriteError(w, http.StatusUnauthorized, "Token não fornecido", "Header Authorization é obrigatório")
				return
			}

			// Validar formato Bearer
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				response.WriteError(w, http.StatusUnauthorized, "Formato de token inválido", "Use: Bearer <token>")
				return
			}

			tokenString := tokenParts[1]

			// Parsear e validar token JWT
			token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
				}
				return jwtSecret, nil
			})

			if err != nil {
				response.WriteError(w, http.StatusUnauthorized, "Token inválido", err.Error())
				return
			}

			if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
				// Injetar dados do usuário no contexto da requisição
				ctx := context.WithValue(r.Context(), "user_admin_id", claims.UserID)
				ctx = context.WithValue(ctx, "empresa_id", claims.EmpresaID)
				ctx = context.WithValue(ctx, "user_email", claims.Email)
				
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				response.WriteError(w, http.StatusUnauthorized, "Token inválido", "Claims inválidas")
				return
			}
		})
	}
}

// EmpresaAuthMiddleware valida autorização de acesso a recursos da empresa
func EmpresaAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar presença de informações de empresa no contexto
		userEmpresaID := r.Context().Value("empresa_id")
		if userEmpresaID == nil {
			response.WriteError(w, http.StatusUnauthorized, "Contexto inválido", "Informações de empresa não encontradas")
			return
		}

		// Ponto de extensão para validações adicionais de autorização
		
		next.ServeHTTP(w, r)
	})
}

// RateLimitMiddleware implementa controle de taxa de requisições
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implementação placeholder para rate limiting
		// Produção requer solução robusta com Redis ou similar
		
		next.ServeHTTP(w, r)
	})
}

// ContentTypeMiddleware valida Content-Type JSON em requisições de escrita
func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" || r.Method == "PUT" {
			contentType := r.Header.Get("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				response.WriteError(w, http.StatusBadRequest, "Content-Type inválido", "Content-Type deve ser application/json")
				return
			}
		}
		
		next.ServeHTTP(w, r)
	})
}

// RecoveryMiddleware captura panics e retorna erro controlado
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("PANIC: %v\n", err)
				response.WriteError(w, http.StatusInternalServerError, "Erro interno do servidor", "Ocorreu um erro inesperado")
			}
		}()
		
		next.ServeHTTP(w, r)
	})
}

// PublicRouteMiddleware adiciona headers de segurança para rotas públicas
func PublicRouteMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Configurar headers de segurança básicos
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		
		next.ServeHTTP(w, r)
	})
}

// ActiveSurveyMiddleware valida se pesquisa aceita novas respostas
func ActiveSurveyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ponto de extensão para validação de status da pesquisa
		// Validação real implementada no use case
		
		next.ServeHTTP(w, r)
	})
}

// ChainMiddleware compõe múltiplos middlewares em ordem de execução
func ChainMiddleware(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}

// PublicMiddlewares retorna cadeia de middlewares para rotas públicas
func PublicMiddlewares() func(http.Handler) http.Handler {
	return ChainMiddleware(
		RecoveryMiddleware,
		CORSMiddleware,
		LoggingMiddleware,
		RateLimitMiddleware,
		PublicRouteMiddleware,
	)
}

// AuthenticatedMiddlewares retorna cadeia de middlewares para rotas autenticadas
func AuthenticatedMiddlewares(jwtSecret []byte) func(http.Handler) http.Handler {
	return ChainMiddleware(
		RecoveryMiddleware,
		CORSMiddleware,
		LoggingMiddleware,
		RateLimitMiddleware,
		ContentTypeMiddleware,
		JWTAuthMiddleware(jwtSecret),
		EmpresaAuthMiddleware,
	)
}

// AdminMiddlewares retorna cadeia de middlewares para rotas administrativas
func AdminMiddlewares(jwtSecret []byte) func(http.Handler) http.Handler {
	return ChainMiddleware(
		RecoveryMiddleware,
		CORSMiddleware,
		LoggingMiddleware,
		RateLimitMiddleware,
		ContentTypeMiddleware,
		JWTAuthMiddleware(jwtSecret),
		EmpresaAuthMiddleware,
	)
}

// SurveySubmissionMiddlewares retorna cadeia de middlewares para submissão de respostas
func SurveySubmissionMiddlewares() func(http.Handler) http.Handler {
	return ChainMiddleware(
		RecoveryMiddleware,
		CORSMiddleware,
		LoggingMiddleware,
		RateLimitMiddleware,
		ContentTypeMiddleware,
		ActiveSurveyMiddleware,
	)
}