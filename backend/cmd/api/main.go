package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"organizational-climate-survey/backend/internal/domain/usecase"
	"organizational-climate-survey/backend/internal/infrastructure/postgres"
	httpRouter "organizational-climate-survey/backend/internal/infrastructure/http"
	"organizational-climate-survey/backend/pkg/crypto"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: Não foi possível encontrar o arquivo .env, usando variáveis de ambiente do sistema.")
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET não configurado nas variáveis de ambiente")
	}

	dbConfig := &postgres.Config{
		Host:     getEnvWithDefault("DB_HOST", "localhost"),
		Port:     getEnvWithDefault("DB_PORT", "5432"),
		User:     getEnvWithDefault("DB_USER", "postgres"),
		Password: os.Getenv("DB_PASS"),
		DBName:   getEnvWithDefault("DB_NAME", "Atmos"),
		SSLMode:  getEnvWithDefault("DB_SSLMODE", "disable"),
	}

	if dbConfig.Password == "" {
		log.Fatal("DB_PASS não configurado nas variáveis de ambiente")
	}

	db, err := postgres.NewDB(dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	log.Println("✅ Conexão com banco de dados estabelecida")

	repos := postgres.NewRepositories(db)
	log.Println("✅ Repositórios inicializados")

	// Inicializar crypto service
	cryptoSvc := crypto.NewDefaultCryptoService()
	log.Println("✅ Crypto service inicializado")

	var empresaUseCase *usecase.EmpresaUseCase
	if repos.Empresa != nil && repos.LogAuditoria != nil {
		empresaUseCase = usecase.NewEmpresaUseCase(repos.Empresa, repos.LogAuditoria)
	}

	var usuarioUseCase *usecase.UsuarioAdministradorUseCase  
	if repos.UsuarioAdministrador != nil && repos.Empresa != nil && repos.LogAuditoria != nil {
		usuarioUseCase = usecase.NewUsuarioAdministradorUseCase(
			repos.UsuarioAdministrador, 
			repos.Empresa, 
			repos.LogAuditoria,
			cryptoSvc,
		)
	}

	var setorUseCase *usecase.SetorUseCase
	if repos.Setor != nil && repos.Empresa != nil && repos.LogAuditoria != nil {
		setorUseCase = usecase.NewSetorUseCase(repos.Setor, repos.Empresa, repos.LogAuditoria)
	}

	var pesquisaUseCase *usecase.PesquisaUseCase
	if repos.Pesquisa != nil && repos.Empresa != nil && repos.Setor != nil && repos.Dashboard != nil && repos.LogAuditoria != nil {
		pesquisaUseCase = usecase.NewPesquisaUseCase(repos.Pesquisa, repos.Empresa, repos.Setor, repos.Dashboard, repos.LogAuditoria)
	}

	var perguntaUseCase *usecase.PerguntaUseCase
	if repos.Pergunta != nil && repos.Resposta != nil && repos.Pesquisa != nil && repos.LogAuditoria != nil {
		perguntaUseCase = usecase.NewPerguntaUseCase(repos.Pergunta, repos.Resposta, repos.Pesquisa, repos.LogAuditoria)
	}

	var respostaUseCase *usecase.RespostaUseCase
	if repos.Resposta != nil && repos.Pergunta != nil && repos.Pesquisa != nil {
		respostaUseCase = usecase.NewRespostaUseCase(repos.Resposta, repos.Pergunta, repos.Pesquisa)
	}
	
	var logUseCase *usecase.LogAuditoriaUseCase
	if repos.LogAuditoria != nil && repos.UsuarioAdministrador != nil && repos.Empresa != nil {
		logUseCase = usecase.NewLogAuditoriaUseCase(repos.LogAuditoria, repos.UsuarioAdministrador, repos.Empresa)
	}

	var dashboardUseCase *usecase.DashboardUseCase
	if repos.Dashboard != nil && repos.Pesquisa != nil && repos.Empresa != nil && repos.LogAuditoria != nil {
		dashboardUseCase = usecase.NewDashboardUseCase(repos.Dashboard, repos.Pesquisa, repos.Empresa, repos.LogAuditoria)
	}

	log.Println("✅ Use cases inicializados")

	routerConfig := &httpRouter.RouterConfig{
		EmpresaUseCase:              empresaUseCase,
		UsuarioAdministradorUseCase: usuarioUseCase,
		SetorUseCase:                setorUseCase,
		PesquisaUseCase:             pesquisaUseCase,
		PerguntaUseCase:             perguntaUseCase,
		RespostaUseCase:             respostaUseCase,
		DashboardUseCase:            dashboardUseCase,
		LogAuditoriaUseCase:         logUseCase,
		JWTSecret:                   jwtSecret,
	}

	router := httpRouter.SetupRouter(routerConfig)
	log.Println("✅ Router configurado")

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	appName := getEnvWithDefault("APP_NAME", "organizational-climate-survey")
	appEnv := getEnvWithDefault("APP_ENV", "development")
	
	fmt.Printf("🚀 Servidor '%s' iniciado na porta %s em modo '%s'\n", appName, port, appEnv)
	fmt.Printf("🔗 API Base URL: http://localhost:%s/api/v1\n", port)
	fmt.Printf("📊 Health Check: http://localhost:%s/health\n", port)
	
	if appEnv == "development" {
		fmt.Printf("📚 Documentação: http://localhost:%s/docs/\n", port)
	}

	log.Fatal(server.ListenAndServe())
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}