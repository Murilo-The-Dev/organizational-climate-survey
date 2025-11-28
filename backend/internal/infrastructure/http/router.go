// Package http implementa o roteamento HTTP e configuração do servidor web.
// Fornece setup de rotas, middleware e handlers para a API REST.
package http

import (
	"net/http"
	"time"

	"organizational-climate-survey/backend/internal/application/handler"
	"organizational-climate-survey/backend/internal/application/middleware"
	"organizational-climate-survey/backend/internal/domain/usecase"
	"organizational-climate-survey/backend/internal/infrastructure/auth"
	"organizational-climate-survey/backend/pkg/logger"
	"organizational-climate-survey/backend/pkg/validator"

	"github.com/gorilla/mux"
)

// RouterConfig contém todas as dependências necessárias para configurar as rotas
type RouterConfig struct {
	EmpresaUseCase              *usecase.EmpresaUseCase              // Use case de empresa
	UsuarioAdministradorUseCase *usecase.UsuarioAdministradorUseCase // Use case de usuário admin
	SetorUseCase                *usecase.SetorUseCase                // Use case de setor
	PesquisaUseCase             *usecase.PesquisaUseCase             // Use case de pesquisa
	PerguntaUseCase             *usecase.PerguntaUseCase             // Use case de pergunta
	RespostaUseCase             *usecase.RespostaUseCase             // Use case de resposta
	DashboardUseCase            *usecase.DashboardUseCase            // Use case de dashboard
	LogAuditoriaUseCase         *usecase.LogAuditoriaUseCase         // Use case de log
	JWTSecret                   string                               // Chave secreta para JWT
}

// SetupRouter configura todas as rotas da API com seus respectivos handlers
func SetupRouter(config *RouterConfig) *mux.Router {
	router := mux.NewRouter()

	log := logger.New(nil)
	val := validator.New()

	authHandler := auth.NewAuthHandler(
		config.UsuarioAdministradorUseCase,
		config.LogAuditoriaUseCase,
		config.JWTSecret,
	)

	var empresaHandler *handler.EmpresaHandler
	if config.EmpresaUseCase != nil {
		empresaHandler = handler.NewEmpresaHandler(config.EmpresaUseCase, log, val)
	}

	var usuarioHandler *handler.UsuarioAdministradorHandler
	if config.UsuarioAdministradorUseCase != nil {
		usuarioHandler = handler.NewUsuarioAdministradorHandler(config.UsuarioAdministradorUseCase, log, val)
	}

	var setorHandler *handler.SetorHandler
	if config.SetorUseCase != nil {
		setorHandler = handler.NewSetorHandler(config.SetorUseCase, log)
	}

	var pesquisaHandler *handler.PesquisaHandler
	if config.PesquisaUseCase != nil {
		pesquisaHandler = handler.NewPesquisaHandler(config.PesquisaUseCase, log)
	}

	var perguntaHandler *handler.PerguntaHandler
	if config.PerguntaUseCase != nil {
		perguntaHandler = handler.NewPerguntaHandler(config.PerguntaUseCase, log)
	}

	var respostaHandler *handler.RespostaHandler
	if config.RespostaUseCase != nil {
		respostaHandler = handler.NewRespostaHandler(config.RespostaUseCase, log)
	}

	var dashboardHandler *handler.DashboardHandler
	if config.DashboardUseCase != nil {
		dashboardHandler = handler.NewDashboardHandler(config.DashboardUseCase, log)
	}

	var logHandler *handler.LogAuditoriaHandler
	if config.LogAuditoriaUseCase != nil {
		logHandler = handler.NewLogAuditoriaHandler(config.LogAuditoriaUseCase, log)
	}

	api := router.PathPrefix("/api/v1").Subrouter()

	publicRoutes := api.PathPrefix("").Subrouter()
	publicRoutes.Use(middleware.PublicMiddlewares())

	authHandler.RegisterRoutes(publicRoutes)

	if respostaHandler != nil {
		surveyRoutes := api.PathPrefix("").Subrouter()
		surveyRoutes.Use(middleware.SurveySubmissionMiddlewares())
		surveyRoutes.HandleFunc("/respostas/submit", respostaHandler.SubmitRespostas).Methods("POST")
	}

	authRoutes := api.PathPrefix("").Subrouter()
	authRoutes.Use(middleware.AuthenticatedMiddlewares([]byte(config.JWTSecret)))

	if empresaHandler != nil {
		empresaHandler.RegisterRoutes(authRoutes)
	}
	if usuarioHandler != nil {
		usuarioHandler.RegisterRoutes(authRoutes)
	}
	if setorHandler != nil {
		setorHandler.RegisterRoutes(authRoutes)
	}
	if pesquisaHandler != nil {
		pesquisaHandler.RegisterRoutes(authRoutes)
	}
	if perguntaHandler != nil {
		perguntaHandler.RegisterRoutes(authRoutes)
	}
	if dashboardHandler != nil {
		dashboardHandler.RegisterRoutes(authRoutes)
	}

	adminRoutes := api.PathPrefix("").Subrouter()
	adminRoutes.Use(middleware.AdminMiddlewares([]byte(config.JWTSecret)))

	if logHandler != nil {
		logHandler.RegisterRoutes(adminRoutes)
	}

	if respostaHandler != nil {
		respostaAdminRoutes := api.PathPrefix("").Subrouter()
		respostaAdminRoutes.Use(middleware.AuthenticatedMiddlewares([]byte(config.JWTSecret)))

		respostaAdminRoutes.HandleFunc("/pesquisas/{pesquisa_id:[0-9]+}/respostas/stats", respostaHandler.GetRespostaStats).Methods("GET")
		respostaAdminRoutes.HandleFunc("/pesquisas/{pesquisa_id:[0-9]+}/respostas/aggregated", respostaHandler.GetRespostasByPesquisa).Methods("GET")
		respostaAdminRoutes.HandleFunc("/pesquisas/{pesquisa_id:[0-9]+}/respostas/by-date", respostaHandler.GetRespostasByDateRange).Methods("GET")
		respostaAdminRoutes.HandleFunc("/pesquisas/{pesquisa_id:[0-9]+}/respostas/count", respostaHandler.CountRespostasByPesquisa).Methods("GET")
		respostaAdminRoutes.HandleFunc("/pesquisas/{pesquisa_id:[0-9]+}/respostas", respostaHandler.DeleteRespostasByPesquisa).Methods("DELETE")
		respostaAdminRoutes.HandleFunc("/perguntas/{pergunta_id:[0-9]+}/respostas/aggregated", respostaHandler.GetRespostasByPergunta).Methods("GET")
		respostaAdminRoutes.HandleFunc("/perguntas/{pergunta_id:[0-9]+}/respostas/count", respostaHandler.CountRespostasByPergunta).Methods("GET")
		respostaAdminRoutes.HandleFunc("/perguntas/{pergunta_id:[0-9]+}/respostas/stats", respostaHandler.GetStatsByPergunta).Methods("GET")
	}

	router.HandleFunc("/health", HealthCheckHandler).Methods("GET")
	router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs/"))))

	return router
}

// HealthCheckHandler responde às requisições de verificação de saúde da API
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	timestamp := time.Now().UTC().Format(time.RFC3339)

	response := `{
		"status": "ok",
		"message": "API está funcionando",
		"version": "1.0.0",
		"timestamp": "` + timestamp + `"
	}`

	w.Write([]byte(response))
}

// SetupCORSRouter configura o router com middleware CORS habilitado
func SetupCORSRouter(config *RouterConfig) *mux.Router {
	router := SetupRouter(config)
	router.Use(middleware.CORSMiddleware)
	return router
}

// SetupMinimalRouter configura um router mínimo com apenas rotas essenciais
// Útil para ambientes de teste ou desenvolvimento
func SetupMinimalRouter(config *RouterConfig) *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.RecoveryMiddleware)
	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.LoggingMiddleware)

	router.HandleFunc("/health", HealthCheckHandler).Methods("GET")

	if config.UsuarioAdministradorUseCase != nil && config.LogAuditoriaUseCase != nil {
		authHandler := auth.NewAuthHandler(
			config.UsuarioAdministradorUseCase,
			config.LogAuditoriaUseCase,
			config.JWTSecret,
		)

		api := router.PathPrefix("/api/v1").Subrouter()
		authHandler.RegisterRoutes(api)
	}

	return router
}
