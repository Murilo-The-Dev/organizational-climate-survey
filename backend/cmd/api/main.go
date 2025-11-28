// Package main inicializa a aplicação, conectando ao banco, configurando repositórios, use cases e servidor HTTP.
package main

import (
	"fmt"
	"log"
	"net/http"

	"organizational-climate-survey/backend/config"
	"organizational-climate-survey/backend/internal/domain/usecase"
	httpRouter "organizational-climate-survey/backend/internal/infrastructure/http"
	"organizational-climate-survey/backend/internal/infrastructure/postgres"
	"organizational-climate-survey/backend/pkg/crypto"

	"github.com/joho/godotenv"
)

func main() {
	// Carrega variáveis de ambiente do .env, se existir
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: Não foi possível encontrar o arquivo .env, usando variáveis de ambiente do sistema.")
	}

	// Carregar configurações
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	// Configuração e conexão com o banco de dados
	db, err := postgres.NewDB(
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()
	log.Println("✅ Conexão com banco de dados estabelecida")

	// Inicializa repositórios
	repos := postgres.NewRepositories(db)
	log.Println("✅ Repositórios inicializados")

	// Inicializa serviço de criptografia
	cryptoSvc := crypto.NewDefaultCryptoService()
	log.Println("✅ Crypto service inicializado")

	// Bootstrap Use Case (não depende de outros use cases)
	var bootstrapUseCase *usecase.BootstrapUseCase
	if repos.Empresa != nil && repos.UsuarioAdministrador != nil {
    bootstrapUseCase = usecase.NewBootstrapUseCase(
        repos.Empresa,
        repos.UsuarioAdministrador,
        repos.LogAuditoria,
        cryptoSvc,
    )
	}
			
	// Inicializa Use Cases
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

	// NOVO: SubmissaoPesquisaUseCase
	var submissaoUseCase *usecase.SubmissaoPesquisaUseCase
	if repos.SubmissaoPesquisa != nil && repos.Pesquisa != nil {
		submissaoUseCase = usecase.NewSubmissaoPesquisaUseCase(
			repos.SubmissaoPesquisa,
			repos.Pesquisa,
			cryptoSvc,
			cfg.Crypto.HashSalt,
		)
	}

	// MODIFICADO: RespostaUseCase agora depende de SubmissaoUseCase
	var respostaUseCase *usecase.RespostaUseCase
	if repos.Resposta != nil && repos.Pergunta != nil && repos.Pesquisa != nil && submissaoUseCase != nil {
		respostaUseCase = usecase.NewRespostaUseCase(
			repos.Resposta, 
			repos.Pergunta, 
			repos.Pesquisa,
			submissaoUseCase,
		)
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

	// Configuração do router HTTP
	routerConfig := &httpRouter.RouterConfig{
		EmpresaUseCase:              empresaUseCase,
		UsuarioAdministradorUseCase: usuarioUseCase,
		SetorUseCase:                setorUseCase,
		PesquisaUseCase:             pesquisaUseCase,
		PerguntaUseCase:             perguntaUseCase,
		RespostaUseCase:             respostaUseCase,
		SubmissaoUseCase:            submissaoUseCase, 
		DashboardUseCase:            dashboardUseCase,
		LogAuditoriaUseCase:         logUseCase,
		PesquisaRepo:                repos.Pesquisa,   
		JWTSecret:                   cfg.JWT.Secret,
		BootstrapUseCase: 			 bootstrapUseCase, 
	}
	router := httpRouter.SetupRouter(routerConfig)
	log.Println("✅ Router configurado")

	// Inicializa servidor HTTP
	server := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: router,
	}

	fmt.Printf("🚀 Servidor '%s' iniciado na porta %s em modo '%s'\n", cfg.App.Name, cfg.App.Port, cfg.App.Env)
	fmt.Printf("🔗 API Base URL: http://localhost:%s/api/v1\n", cfg.App.Port)
	fmt.Printf("📊 Health Check: http://localhost:%s/health\n", cfg.App.Port)
	if cfg.App.Env == "development" {
		fmt.Printf("📚 Documentação: http://localhost:%s/docs/\n", cfg.App.Port)
	}

	log.Fatal(server.ListenAndServe())
}