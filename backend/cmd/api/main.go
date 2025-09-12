package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"organizational-climate-survey/backend/internal/domain/usecase"
	"organizational-climate-survey/backend/internal/infrastructure/postgres"
	"organizational-climate-survey/backend/internal/application/middleware"

	"github.com/joho/godotenv"
)

func main() {
	// Carregar variáveis de ambiente do arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: Não foi possível encontrar o arquivo .env, usando variáveis de ambiente do sistema.")
	}

	// Configurações do servidor
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET não configurado nas variáveis de ambiente")
	}

	// Configuração do banco de dados
	dbConfig := &postgres.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	// Conectar ao banco de dados
	db, err := postgres.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Criar repositórios
	empresaRepo := postgres.NewEmpresaRepository(db)
	usuarioRepo := postgres.NewUsuarioAdministradorRepository(db)
	setorRepo := postgres.NewSetorRepository(db)
	pesquisaRepo := postgres.NewPesquisaRepository(db)
	perguntaRepo := postgres.NewPerguntaRepository(db)
	respostaRepo := postgres.NewRespostaRepository(db)
	dashboardRepo := postgres.NewDashboardRepository(db)
	logRepo := postgres.NewLogAuditoriaRepository(db)

	// Criar use cases
	empresaUseCase := usecase.NewEmpresaUseCase(empresaRepo, logRepo)
	usuarioUseCase := usecase.NewUsuarioAdministradorUseCase(usuarioRepo, empresaRepo, logRepo)
	setorUseCase := usecase.NewSetorUseCase(setorRepo, empresaRepo, logRepo)
	pesquisaUseCase := usecase.NewPesquisaUseCase(pesquisaRepo, empresaRepo, usuarioRepo, setorRepo, logRepo)
	perguntaUseCase := usecase.NewPerguntaUseCase(perguntaRepo, pesquisaRepo, logRepo)
	respostaUseCase := usecase.NewRespostaUseCase(respostaRepo, perguntaRepo, pesquisaRepo, logRepo)
	dashboardUseCase := usecase.NewDashboardUseCase(dashboardRepo, pesquisaRepo, respostaRepo, logRepo)
	logUseCase := usecase.NewLogAuditoriaUseCase(logRepo, usuarioRepo, empresaRepo)

	// Configurar roteador
	routerConfig := &middleware.RouterConfig{
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

	router := middleware.SetupRouter(routerConfig)

	// Configurar servidor
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Iniciar servidor
	fmt.Printf("🚀 Servidor '%s' iniciado na porta %s em modo '%s'\n", os.Getenv("APP_NAME"), port, os.Getenv("APP_ENV"))
	fmt.Printf("🔗 API Base URL: http://localhost:%s/api/v1\n", port)

	log.Fatal(server.ListenAndServe())
}