package http

import (
	"net/http"
	"time"

	"organizational-climate-survey/backend/internal/application/handler"
	"organizational-climate-survey/backend/internal/application/middleware"
	"organizational-climate-survey/backend/internal/domain/usecase"
	"organizational-climate-survey/backend/internal/infrastructure/auth"
	"organizational-climate-survey/backend/pkg/logger"

	"github.com/gorilla/mux"
)

// RouterConfig contém as configurações necessárias para o roteador
type RouterConfig struct {
	EmpresaUseCase              *usecase.EmpresaUseCase
	UsuarioAdministradorUseCase *usecase.UsuarioAdministradorUseCase
	SetorUseCase                *usecase.SetorUseCase
	PesquisaUseCase             *usecase.PesquisaUseCase
	PerguntaUseCase             *usecase.PerguntaUseCase
	RespostaUseCase             *usecase.RespostaUseCase
	DashboardUseCase            *usecase.DashboardUseCase
	LogAuditoriaUseCase         *usecase.LogAuditoriaUseCase
	JWTSecret                   string
}

// SetupRouter configura e retorna o roteador principal com todas as rotas
func SetupRouter(config *RouterConfig) *mux.Router {
	// Criar roteador principal
	router := mux.NewRouter()

	// Criar instância de logger (usar DefaultConfig se necessário)
	log := logger.New(nil)

	// Criar handlers
	authHandler := auth.NewAuthHandler(
		config.UsuarioAdministradorUseCase,
		config.LogAuditoriaUseCase,
		config.JWTSecret,
	)

	var empresaHandler *handler.EmpresaHandler
	if config.EmpresaUseCase != nil {
		empresaHandler = handler.NewEmpresaHandler(config.EmpresaUseCase, log)
	}

	var usuarioHandler *handler.UsuarioAdministradorHandler
	if config.UsuarioAdministradorUseCase != nil {
		usuarioHandler = handler.NewUsuarioAdministradorHandler(config.UsuarioAdministradorUseCase, log)
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

	// Configurar rotas de API
	api := router.PathPrefix("/api/v1").Subrouter()

	// Rotas públicas (sem autenticação)
	publicRoutes := api.PathPrefix("").Subrouter()
	publicRoutes.Use(middleware.PublicMiddlewares())

	// Rotas de autenticação
	authHandler.RegisterRoutes(publicRoutes)

	// Rota pública para submissão de respostas
	if respostaHandler != nil {
		surveyRoutes := api.PathPrefix("").Subrouter()
		surveyRoutes.Use(middleware.SurveySubmissionMiddlewares())
		surveyRoutes.HandleFunc("/respostas/submit", respostaHandler.SubmitRespostas).Methods("POST")
	}

	// Rotas autenticadas
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

	// Rotas administrativas (com middlewares extras)
	adminRoutes := api.PathPrefix("").Subrouter()
	adminRoutes.Use(middleware.AdminMiddlewares([]byte(config.JWTSecret)))

	if logHandler != nil {
		logHandler.RegisterRoutes(adminRoutes)
	}

	// Rotas de respostas administrativas (exceto submissão pública)
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

	// Rota de health check
	router.HandleFunc("/health", HealthCheckHandler).Methods("GET")

	// Rota de documentação da API (Swagger)
	router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs/"))))

	return router
}

// HealthCheckHandler verifica se a API está funcionando
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

// SetupCORSRouter configura um roteador com CORS habilitado para desenvolvimento
func SetupCORSRouter(config *RouterConfig) *mux.Router {
	router := SetupRouter(config)
	router.Use(middleware.CORSMiddleware)
	return router
}

// SetupMinimalRouter cria um router mínimo apenas com as funcionalidades básicas
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
