package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"organizational-climate-survey/backend/internal/application/dto/response"

	"github.com/golang-jwt/jwt/v5"
)

// Middleware de CORS
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Permitir todas as origens (em produção, especificar origens específicas)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 horas

		// Responder a requisições OPTIONS (preflight)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Middleware de logging
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log básico da requisição
		fmt.Printf("[%s] %s %s - %s\n", 
			r.Method, 
			r.URL.Path, 
			r.RemoteAddr,
			r.UserAgent(),
		)
		
		next.ServeHTTP(w, r)
	})
}

// Middleware de autenticação JWT
func JWTAuthMiddleware(jwtSecret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extrair token do header Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.WriteError(w, http.StatusUnauthorized, "Token não fornecido", "Header Authorization é obrigatório")
				return
			}

			// Verificar formato "Bearer <token>"
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				response.WriteError(w, http.StatusUnauthorized, "Formato de token inválido", "Use: Bearer <token>")
				return
			}

			tokenString := tokenParts[1]

			// Validar token JWT
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
				// Adicionar informações do usuário ao contexto
				ctx := context.WithValue(r.Context(), "user_admin_id", claims.UserID)
				ctx = context.WithValue(ctx, "empresa_id", claims.EmpresaID)
				ctx = context.WithValue(ctx, "user_email", claims.Email)
				
				// Continuar com o contexto atualizado
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				response.WriteError(w, http.StatusUnauthorized, "Token inválido", "Claims inválidas")
				return
			}
		})
	}
}

// Middleware de autorização por empresa
func EmpresaAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Este middleware verifica se o usuário tem acesso aos recursos da empresa
		// Pode ser usado em rotas que incluem empresa_id nos parâmetros
		
		userEmpresaID := r.Context().Value("empresa_id")
		if userEmpresaID == nil {
			response.WriteError(w, http.StatusUnauthorized, "Contexto inválido", "Informações de empresa não encontradas")
			return
		}

		// Aqui poderia ser implementada lógica adicional de autorização
		// Por exemplo, verificar se o usuário tem permissão para acessar recursos de outras empresas
		
		next.ServeHTTP(w, r)
	})
}

// Middleware de rate limiting básico
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implementação básica de rate limiting
		// Em produção, usar uma solução mais robusta como Redis
		
		// Por enquanto, apenas passa adiante
		next.ServeHTTP(w, r)
	})
}

// Middleware de validação de Content-Type para requisições POST/PUT
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

// Middleware de recuperação de panic
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

// Middleware para rotas públicas (sem autenticação)
func PublicRouteMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Adicionar headers de segurança básicos
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		
		next.ServeHTTP(w, r)
	})
}

// Middleware para validar se a pesquisa está ativa (para submissão de respostas)
func ActiveSurveyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Este middleware seria usado na rota de submissão de respostas
		// para verificar se a pesquisa está ativa e aceita respostas
		
		// Por enquanto, apenas passa adiante
		// A validação real seria feita no use case
		next.ServeHTTP(w, r)
	})
}

// Função auxiliar para aplicar múltiplos middlewares
func ChainMiddleware(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}

// Configuração de middlewares para diferentes tipos de rotas

// Middlewares para rotas públicas
func PublicMiddlewares() func(http.Handler) http.Handler {
	return ChainMiddleware(
		RecoveryMiddleware,
		CORSMiddleware,
		LoggingMiddleware,
		RateLimitMiddleware,
		PublicRouteMiddleware,
	)
}

// Middlewares para rotas autenticadas
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

// Middlewares para rotas administrativas (com validações extras)
func AdminMiddlewares(jwtSecret []byte) func(http.Handler) http.Handler {
	return ChainMiddleware(
		RecoveryMiddleware,
		CORSMiddleware,
		LoggingMiddleware,
		RateLimitMiddleware,
		ContentTypeMiddleware,
		JWTAuthMiddleware(jwtSecret),
		EmpresaAuthMiddleware,
		// Aqui poderiam ser adicionados middlewares específicos para admins
	)
}

// Middlewares para submissão de respostas (público, mas com validações específicas)
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