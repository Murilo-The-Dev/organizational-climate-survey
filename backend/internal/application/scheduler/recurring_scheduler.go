package scheduler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"organizational-climate-survey/backend/internal/domain/entity"
	"organizational-climate-survey/backend/internal/domain/repository"
	"organizational-climate-survey/backend/pkg/logger"

	"github.com/robfig/cron/v3"
)

type recurrenceConfig struct {
	Enabled      bool
	IntervalDays int
	DurationDays int
	AutoActivate bool
	NextRunAt    time.Time
}

// RecurringSurveyScheduler processa pesquisas com config_recorrencia em background.
type RecurringSurveyScheduler struct {
	cron          *cron.Cron
	empresaRepo   repository.EmpresaRepository
	pesquisaRepo  repository.PesquisaRepository
	perguntaRepo  repository.PerguntaRepository
	dashboardRepo repository.DashboardRepository
	log           logger.Logger
}

func NewRecurringSurveyScheduler(
	empresaRepo repository.EmpresaRepository,
	pesquisaRepo repository.PesquisaRepository,
	perguntaRepo repository.PerguntaRepository,
	dashboardRepo repository.DashboardRepository,
	log logger.Logger,
) *RecurringSurveyScheduler {
	if log == nil {
		log = logger.New(nil)
	}

	return &RecurringSurveyScheduler{
		cron:          cron.New(cron.WithSeconds()),
		empresaRepo:   empresaRepo,
		pesquisaRepo:  pesquisaRepo,
		perguntaRepo:  perguntaRepo,
		dashboardRepo: dashboardRepo,
		log:           log,
	}
}

func (s *RecurringSurveyScheduler) Start() error {
	_, err := s.cron.AddFunc("@every 1m", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.RunCycle(ctx); err != nil {
			s.log.Error("erro no ciclo do scheduler de recorrencia: %v", err)
		}
	})
	if err != nil {
		return err
	}

	s.cron.Start()
	s.log.Info("scheduler de recorrencia iniciado")
	return nil
}

func (s *RecurringSurveyScheduler) Stop() {
	s.cron.Stop()
	s.log.Info("scheduler de recorrencia finalizado")
}

// RunCycle executa um ciclo de processamento de recorrências; útil para testes e execução manual.
func (s *RecurringSurveyScheduler) RunCycle(ctx context.Context) error {
	empresas, err := s.empresaRepo.List(ctx, 10000, 0)
	if err != nil {
		return fmt.Errorf("erro ao listar empresas para recorrencia: %v", err)
	}

	now := time.Now()
	for _, empresa := range empresas {
		pesquisas, err := s.listEmpresaSurveys(ctx, empresa.ID)
		if err != nil {
			s.log.Error("erro ao listar pesquisas da empresa %d: %v", empresa.ID, err)
			continue
		}

		for _, pesquisa := range pesquisas {
			if pesquisa.ConfigRecorrencia == nil || strings.TrimSpace(*pesquisa.ConfigRecorrencia) == "" {
				continue
			}

			cfg, err := parseRecurrenceConfig(*pesquisa.ConfigRecorrencia, pesquisa)
			if err != nil {
				s.log.Warn("config_recorrencia inválida na pesquisa %d: %v", pesquisa.ID, err)
				continue
			}

			if !cfg.Enabled || cfg.NextRunAt.After(now) {
				continue
			}

			if err := s.materializeRecurringSurvey(ctx, pesquisa, cfg); err != nil {
				s.log.Error("erro ao materializar recorrencia da pesquisa %d: %v", pesquisa.ID, err)
				continue
			}

			cfg.NextRunAt = advanceNextRun(cfg.NextRunAt, cfg.IntervalDays, now)
			updatedCfg := serializeRecurrenceConfig(cfg)
			pesquisa.ConfigRecorrencia = &updatedCfg

			if err := s.pesquisaRepo.Update(ctx, pesquisa); err != nil {
				s.log.Error("erro ao atualizar proxima execução da pesquisa %d: %v", pesquisa.ID, err)
			}
		}
	}

	return nil
}

func (s *RecurringSurveyScheduler) listEmpresaSurveys(ctx context.Context, empresaID int) ([]*entity.Pesquisa, error) {
	statuses := []string{"Rascunho", "Ativa", "Concluída", "Arquivada"}
	seen := make(map[int]*entity.Pesquisa)

	for _, status := range statuses {
		items, err := s.pesquisaRepo.ListByStatus(ctx, empresaID, status)
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			seen[item.ID] = item
		}
	}

	out := make([]*entity.Pesquisa, 0, len(seen))
	for _, item := range seen {
		out = append(out, item)
	}
	return out, nil
}

func (s *RecurringSurveyScheduler) materializeRecurringSurvey(ctx context.Context, template *entity.Pesquisa, cfg recurrenceConfig) error {
	link, err := generateUniqueLink()
	if err != nil {
		return err
	}

	openAt := cfg.NextRunAt
	closeAt := openAt.AddDate(0, 0, cfg.DurationDays)
	status := "Rascunho"
	if cfg.AutoActivate {
		status = "Ativa"
	}

	recurring := &entity.Pesquisa{
		IDEmpresa:      template.IDEmpresa,
		IDUserAdmin:    template.IDUserAdmin,
		IDSetor:        template.IDSetor,
		Titulo:         template.Titulo,
		Descricao:      template.Descricao,
		DataCriacao:    time.Now(),
		DataAbertura:   &openAt,
		DataFechamento: &closeAt,
		Status:         status,
		LinkAcesso:     link,
		Anonimato:      template.Anonimato,
	}

	if err := s.pesquisaRepo.Create(ctx, recurring); err != nil {
		return fmt.Errorf("erro ao criar pesquisa recorrente: %v", err)
	}

	perguntas, err := s.perguntaRepo.ListByPesquisa(ctx, template.ID)
	if err != nil {
		return fmt.Errorf("erro ao listar perguntas da pesquisa modelo: %v", err)
	}

	if len(perguntas) > 0 {
		clones := make([]*entity.Pergunta, 0, len(perguntas))
		for _, pergunta := range perguntas {
			clones = append(clones, &entity.Pergunta{
				IDPesquisa:     recurring.ID,
				TextoPergunta:  pergunta.TextoPergunta,
				TipoPergunta:   pergunta.TipoPergunta,
				OrdemExibicao:  pergunta.OrdemExibicao,
				OpcoesResposta: pergunta.OpcoesResposta,
			})
		}

		if err := s.perguntaRepo.CreateBatch(ctx, clones); err != nil {
			return fmt.Errorf("erro ao clonar perguntas da recorrencia: %v", err)
		}
	}

	dashboardTitle := fmt.Sprintf("Dashboard - %s", recurring.Titulo)
	dashboardFilters := `{"filtros_padrao": true}`
	dashboard := &entity.Dashboard{
		IDPesquisa:    recurring.ID,
		Titulo:        dashboardTitle,
		DataCriacao:   time.Now(),
		ConfigFiltros: &dashboardFilters,
	}

	if err := s.dashboardRepo.Create(ctx, dashboard); err != nil {
		s.log.Warn("erro ao criar dashboard da recorrencia para pesquisa %d: %v", recurring.ID, err)
	}

	s.log.Info("pesquisa recorrente criada a partir do modelo %d com novo ID %d", template.ID, recurring.ID)
	return nil
}

func parseRecurrenceConfig(raw string, pesquisa *entity.Pesquisa) (recurrenceConfig, error) {
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return recurrenceConfig{}, err
	}

	intervalDays := asInt(payload["interval_days"], 0)
	if intervalDays <= 0 {
		intervalDays = asInt(payload["intervalo_dias"], 30)
	}

	durationDays := asInt(payload["duration_days"], 0)
	if durationDays <= 0 {
		durationDays = asInt(payload["duracao_dias"], defaultDurationDays(pesquisa))
	}

	autoActivate := asBool(payload["auto_activate"], true)
	if _, ok := payload["ativar_automaticamente"]; ok {
		autoActivate = asBool(payload["ativar_automaticamente"], autoActivate)
	}

	enabled := asBool(payload["enabled"], true)
	if _, ok := payload["habilitado"]; ok {
		enabled = asBool(payload["habilitado"], enabled)
	}

	nextRunRaw := asString(payload["next_run_at"])
	if nextRunRaw == "" {
		nextRunRaw = asString(payload["proxima_execucao"])
	}

	nextRunAt, err := time.Parse(time.RFC3339, nextRunRaw)
	if err != nil {
		if nextRunRaw == "" {
			nextRunAt = time.Now().AddDate(0, 0, intervalDays)
		} else {
			return recurrenceConfig{}, fmt.Errorf("next_run_at inválido: %v", err)
		}
	}

	return recurrenceConfig{
		Enabled:      enabled,
		IntervalDays: intervalDays,
		DurationDays: durationDays,
		AutoActivate: autoActivate,
		NextRunAt:    nextRunAt,
	}, nil
}

func serializeRecurrenceConfig(cfg recurrenceConfig) string {
	payload := map[string]interface{}{
		"enabled":       cfg.Enabled,
		"interval_days": cfg.IntervalDays,
		"duration_days": cfg.DurationDays,
		"auto_activate": cfg.AutoActivate,
		"next_run_at":   cfg.NextRunAt.UTC().Format(time.RFC3339),
	}

	b, _ := json.Marshal(payload)
	return string(b)
}

func generateUniqueLink() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("erro ao gerar link de recorrencia: %v", err)
	}
	return hex.EncodeToString(b), nil
}

func defaultDurationDays(p *entity.Pesquisa) int {
	if p.DataAbertura != nil && p.DataFechamento != nil {
		d := int(p.DataFechamento.Sub(*p.DataAbertura).Hours() / 24)
		if d > 0 {
			return d
		}
	}
	return 30
}

func advanceNextRun(nextRun time.Time, intervalDays int, now time.Time) time.Time {
	if intervalDays <= 0 {
		intervalDays = 30
	}
	for !nextRun.After(now) {
		nextRun = nextRun.AddDate(0, 0, intervalDays)
	}
	return nextRun
}

func asInt(value interface{}, fallback int) int {
	switch v := value.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		parsed, err := strconv.Atoi(strings.TrimSpace(v))
		if err == nil {
			return parsed
		}
	}
	return fallback
}

func asBool(value interface{}, fallback bool) bool {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		s := strings.ToLower(strings.TrimSpace(v))
		if s == "true" || s == "1" || s == "yes" || s == "sim" {
			return true
		}
		if s == "false" || s == "0" || s == "no" || s == "nao" || s == "não" {
			return false
		}
	}
	return fallback
}

func asString(value interface{}) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprintf("%v", value))
}
