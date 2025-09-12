package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"organizational-climate-survey/backend/internal/domain/entity"
	"organizational-climate-survey/backend/internal/application/dto"
	"organizational-climate-survey/backend/internal/domain/usecase"
	"organizational-climate-survey/backend/internal/application/dto/response"

	"github.com/gorilla/mux"
)

type PesquisaHandler struct {
	pesquisaUseCase *usecase.PesquisaUseCase
}

func NewPesquisaHandler(pesquisaUseCase *usecase.PesquisaUseCase) *PesquisaHandler {
	return &PesquisaHandler{
		pesquisaUseCase: pesquisaUseCase,
	}
}

// CreatePesquisa godoc
// @Summary Criar nova pesquisa
// @Description Cria uma nova pesquisa de clima organizacional
// @Tags pesquisas
// @Accept json
// @Produce json
// @Param pesquisa body dto.PesquisaCreateRequest true "Dados da pesquisa"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /pesquisas [post]
func (h *PesquisaHandler) CreatePesquisa(w http.ResponseWriter, r *http.Request) {
	var req dto.PesquisaCreateRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Validação básica
	if err := h.validatePesquisaCreateRequest(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Validação falhou", err.Error())
		return
	}

	pesquisa, err := req.ToEntity()
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Erro na conversão de dados", err.Error())
		return
	}
	
	// Obter informações do usuário autenticado
	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.pesquisaUseCase.Create(r.Context(), pesquisa, userAdminID, clientIP); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Criar resposta usando função helper do response
	pesquisaResponse := h.toPesquisaResponse(pesquisa)
	response.WriteSuccess(w, http.StatusCreated, "Pesquisa criada com sucesso", pesquisaResponse)
}

// GetPesquisa godoc
// @Summary Buscar pesquisa por ID
// @Description Retorna uma pesquisa específica pelo ID
// @Tags pesquisas
// @Produce json
// @Param id path int true "ID da pesquisa"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /pesquisas/{id} [get]
func (h *PesquisaHandler) GetPesquisa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	pesquisa, err := h.pesquisaUseCase.GetByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			response.WriteError(w, http.StatusNotFound, "Pesquisa não encontrada", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	pesquisaResponse := h.toPesquisaResponse(pesquisa)
	response.WriteSuccess(w, http.StatusOK, "Pesquisa encontrada", pesquisaResponse)
}

// GetPesquisaByLink godoc
// @Summary Buscar pesquisa por link de acesso
// @Description Retorna uma pesquisa específica pelo link de acesso público
// @Tags pesquisas
// @Produce json
// @Param link path string true "Link de acesso da pesquisa"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /pesquisas/link/{link} [get]
func (h *PesquisaHandler) GetPesquisaByLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	link := vars["link"]

	if strings.TrimSpace(link) == "" {
		response.WriteError(w, http.StatusBadRequest, "Link inválido", "Link de acesso é obrigatório")
		return
	}

	pesquisa, err := h.pesquisaUseCase.GetByLinkAcesso(r.Context(), link)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			response.WriteError(w, http.StatusNotFound, "Pesquisa não encontrada", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	pesquisaResponse := h.toPesquisaResponse(pesquisa)
	response.WriteSuccess(w, http.StatusOK, "Pesquisa encontrada", pesquisaResponse)
}

// ListPesquisasByEmpresa godoc
// @Summary Listar pesquisas por empresa
// @Description Retorna lista de pesquisas de uma empresa específica
// @Tags pesquisas
// @Produce json
// @Param empresa_id path int true "ID da empresa"
// @Param status query string false "Filtrar por status"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /empresas/{empresa_id}/pesquisas [get]
func (h *PesquisaHandler) ListPesquisasByEmpresa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	empresaID, err := strconv.Atoi(vars["empresa_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da empresa inválido", "ID deve ser um número inteiro")
		return
	}

	status := r.URL.Query().Get("status")
	
	var pesquisas []*entity.Pesquisa
	
	if status != "" {
		pesquisas, err = h.pesquisaUseCase.ListByStatus(r.Context(), empresaID, status)
	} else {
		pesquisas, err = h.pesquisaUseCase.ListByEmpresa(r.Context(), empresaID)
	}

	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Converter para response
	pesquisasResponse := make([]interface{}, len(pesquisas))
	for i, pesquisa := range pesquisas {
		pesquisasResponse[i] = h.toPesquisaResponse(pesquisa)
	}

	response.WriteSuccess(w, http.StatusOK, "Pesquisas listadas com sucesso", pesquisasResponse)
}

// ListPesquisasBySetor godoc
// @Summary Listar pesquisas por setor
// @Description Retorna lista de pesquisas de um setor específico
// @Tags pesquisas
// @Produce json
// @Param setor_id path int true "ID do setor"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /setores/{setor_id}/pesquisas [get]
func (h *PesquisaHandler) ListPesquisasBySetor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	setorID, err := strconv.Atoi(vars["setor_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID do setor inválido", "ID deve ser um número inteiro")
		return
	}

	pesquisas, err := h.pesquisaUseCase.ListBySetor(r.Context(), setorID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Converter para response
	pesquisasResponse := make([]interface{}, len(pesquisas))
	for i, pesquisa := range pesquisas {
		pesquisasResponse[i] = h.toPesquisaResponse(pesquisa)
	}

	response.WriteSuccess(w, http.StatusOK, "Pesquisas do setor listadas com sucesso", pesquisasResponse)
}

// ListPesquisasActive godoc
// @Summary Listar pesquisas ativas
// @Description Retorna lista de pesquisas ativas de uma empresa
// @Tags pesquisas
// @Produce json
// @Param empresa_id path int true "ID da empresa"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /empresas/{empresa_id}/pesquisas/active [get]
func (h *PesquisaHandler) ListPesquisasActive(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	empresaID, err := strconv.Atoi(vars["empresa_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID da empresa inválido", "ID deve ser um número inteiro")
		return
	}

	pesquisas, err := h.pesquisaUseCase.ListActive(r.Context(), empresaID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Converter para response
	pesquisasResponse := make([]interface{}, len(pesquisas))
	for i, pesquisa := range pesquisas {
		pesquisasResponse[i] = h.toPesquisaResponse(pesquisa)
	}

	response.WriteSuccess(w, http.StatusOK, "Pesquisas ativas listadas com sucesso", pesquisasResponse)
}

// UpdatePesquisa godoc
// @Summary Atualizar pesquisa
// @Description Atualiza dados de uma pesquisa existente
// @Tags pesquisas
// @Accept json
// @Produce json
// @Param id path int true "ID da pesquisa"
// @Param pesquisa body dto.PesquisaUpdateRequest true "Dados para atualização"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /pesquisas/{id} [put]
func (h *PesquisaHandler) UpdatePesquisa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	var req dto.PesquisaUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Buscar pesquisa existente
	pesquisa, err := h.pesquisaUseCase.GetByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			response.WriteError(w, http.StatusNotFound, "Pesquisa não encontrada", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// Aplicar atualizações
	if err := req.ApplyToEntity(pesquisa); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Erro na aplicação dos dados", err.Error())
		return
	}

	// Obter informações do usuário autenticado
	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.pesquisaUseCase.Update(r.Context(), pesquisa, userAdminID, clientIP); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	pesquisaResponse := h.toPesquisaResponse(pesquisa)
	response.WriteSuccess(w, http.StatusOK, "Pesquisa atualizada com sucesso", pesquisaResponse)
}

// UpdateStatusPesquisa godoc
// @Summary Atualizar status da pesquisa
// @Description Atualiza apenas o status de uma pesquisa
// @Tags pesquisas
// @Accept json
// @Produce json
// @Param id path int true "ID da pesquisa"
// @Param status body StatusUpdateRequest true "Novo status"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /pesquisas/{id}/status [put]
func (h *PesquisaHandler) UpdateStatusPesquisa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	var req StatusUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	// Validação do status
	if !h.isValidPesquisaStatus(req.Status) {
		response.WriteError(w, http.StatusBadRequest, "Status inválido", "Status deve ser: Rascunho, Ativa, Concluída ou Arquivada")
		return
	}

	// Obter informações do usuário autenticado
	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.pesquisaUseCase.UpdateStatus(r.Context(), id, req.Status, userAdminID, clientIP); err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			response.WriteError(w, http.StatusNotFound, "Pesquisa não encontrada", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Status da pesquisa atualizado com sucesso", nil)
}

// DeletePesquisa godoc
// @Summary Deletar pesquisa
// @Description Remove uma pesquisa do sistema
// @Tags pesquisas
// @Produce json
// @Param id path int true "ID da pesquisa"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /pesquisas/{id} [delete]
func (h *PesquisaHandler) DeletePesquisa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	// Obter informações do usuário autenticado
	userAdminID := h.getUserAdminIDFromContext(r)
	clientIP := h.getClientIP(r)

	if err := h.pesquisaUseCase.Delete(r.Context(), id, userAdminID, clientIP); err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			response.WriteError(w, http.StatusNotFound, "Pesquisa não encontrada", err.Error())
			return
		}
		if strings.Contains(err.Error(), "possui") && strings.Contains(err.Error(), "vinculados") {
			response.WriteError(w, http.StatusConflict, "Pesquisa possui dependências", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Pesquisa deletada com sucesso", nil)
}

// GenerateQRCode godoc
// @Summary Gerar QR Code para pesquisa
// @Description Gera um QR Code para acesso à pesquisa
// @Tags pesquisas
// @Produce json
// @Param id path int true "ID da pesquisa"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /pesquisas/{id}/qrcode [post]
func (h *PesquisaHandler) GenerateQRCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "ID inválido", "ID deve ser um número inteiro")
		return
	}

	// Verificar se pesquisa existe
	pesquisa, err := h.pesquisaUseCase.GetByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrada") {
			response.WriteError(w, http.StatusNotFound, "Pesquisa não encontrada", err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Erro interno", err.Error())
		return
	}

	// TODO: Implementar geração real do QR Code quando método estiver disponível
	// Por enquanto, retorna um path simulado baseado no link de acesso
	qrCodePath := fmt.Sprintf("/qrcodes/%s.png", pesquisa.LinkAcesso)

	qrResponse := map[string]string{
		"qr_code_path": qrCodePath,
		"link_acesso":  pesquisa.LinkAcesso,
	}

	response.WriteSuccess(w, http.StatusOK, "QR Code gerado com sucesso", qrResponse)
}

// Métodos auxiliares

func (h *PesquisaHandler) validatePesquisaCreateRequest(req *dto.PesquisaCreateRequest) error {
	if req.IDEmpresa <= 0 {
		return fmt.Errorf("ID da empresa é obrigatório")
	}
	if req.IDUserAdmin <= 0 {
		return fmt.Errorf("ID do usuário administrador é obrigatório")
	}
	if req.IDSetor <= 0 {
		return fmt.Errorf("ID do setor é obrigatório")
	}
	if strings.TrimSpace(req.Titulo) == "" {
		return fmt.Errorf("título da pesquisa é obrigatório")
	}
	if len(req.Titulo) < 3 {
		return fmt.Errorf("título deve ter pelo menos 3 caracteres")
	}
	if !h.isValidPesquisaStatus(req.Status) {
		return fmt.Errorf("status inválido")
	}
	return nil
}

func (h *PesquisaHandler) isValidPesquisaStatus(status string) bool {
	validStatuses := []string{"Rascunho", "Ativa", "Concluída", "Arquivada"}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

func (h *PesquisaHandler) getUserAdminIDFromContext(r *http.Request) int {
	if userID := r.Context().Value("user_admin_id"); userID != nil {
		if id, ok := userID.(int); ok {
			return id
		}
	}
	return 0
}

func (h *PesquisaHandler) getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

// toPesquisaResponse converte entity.Pesquisa para response.PesquisaResponse
func (h *PesquisaHandler) toPesquisaResponse(pesquisa *entity.Pesquisa) *response.PesquisaResponse {
	return &response.PesquisaResponse{
		ID:             pesquisa.ID,
		IDEmpresa:      pesquisa.IDEmpresa,
		IDSetor:        pesquisa.IDSetor,
		Titulo:         pesquisa.Titulo,
		Descricao:      pesquisa.Descricao,
		DataCriacao:    pesquisa.DataCriacao,
		DataAbertura:   pesquisa.DataAbertura,
		DataFechamento: pesquisa.DataFechamento,
		Status:         pesquisa.Status,
		LinkAcesso:     pesquisa.LinkAcesso,
		QRCodePath:     pesquisa.QRCodePath,
		Anonimato:      pesquisa.Anonimato,
	}
}

// RegisterRoutes registra as rotas do handler
func (h *PesquisaHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/pesquisas", h.CreatePesquisa).Methods("POST")
	router.HandleFunc("/pesquisas/{id:[0-9]+}", h.GetPesquisa).Methods("GET")
	router.HandleFunc("/pesquisas/{id:[0-9]+}", h.UpdatePesquisa).Methods("PUT")
	router.HandleFunc("/pesquisas/{id:[0-9]+}", h.DeletePesquisa).Methods("DELETE")
	router.HandleFunc("/pesquisas/{id:[0-9]+}/status", h.UpdateStatusPesquisa).Methods("PUT")
	router.HandleFunc("/pesquisas/{id:[0-9]+}/qrcode", h.GenerateQRCode).Methods("POST")
	router.HandleFunc("/pesquisas/link/{link}", h.GetPesquisaByLink).Methods("GET")
	router.HandleFunc("/empresas/{empresa_id:[0-9]+}/pesquisas", h.ListPesquisasByEmpresa).Methods("GET")
	router.HandleFunc("/empresas/{empresa_id:[0-9]+}/pesquisas/active", h.ListPesquisasActive).Methods("GET")
	router.HandleFunc("/setores/{setor_id:[0-9]+}/pesquisas", h.ListPesquisasBySetor).Methods("GET")
}