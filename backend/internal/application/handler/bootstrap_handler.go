// Package handler implementa os controladores HTTP da aplicação.
// Processa requisição de bootstrap (inicialização) do sistema.
package handler

import (
    "encoding/json"
    "net/http"
    "strings"

    "organizational-climate-survey/backend/internal/application/dto"
    "organizational-climate-survey/backend/internal/application/dto/response"
    "organizational-climate-survey/backend/internal/domain/usecase"
    "organizational-climate-survey/backend/pkg/logger"
    "organizational-climate-survey/backend/pkg/validator"

    "github.com/gorilla/mux"
)

// BootstrapHandler gerencia inicialização do sistema
type BootstrapHandler struct {
    bootstrapUseCase *usecase.BootstrapUseCase
    log              logger.Logger
    validator        *validator.Validator
}

// NewBootstrapHandler cria nova instância do handler de bootstrap
func NewBootstrapHandler(
    bootstrapUseCase *usecase.BootstrapUseCase,
    log logger.Logger,
    val *validator.Validator,
) *BootstrapHandler {
    return &BootstrapHandler{
        bootstrapUseCase: bootstrapUseCase,
        log:              log,
        validator:        val,
    }
}

// Bootstrap inicializa o sistema criando empresa e primeiro admin
// POST /bootstrap (SEM AUTENTICAÇÃO)
func (h *BootstrapHandler) Bootstrap(w http.ResponseWriter, r *http.Request) {
    var req dto.BootstrapRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.log.WithContext(r.Context()).Warn("Bootstrap decode erro: %v", err)
        response.WriteError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
        return
    }

    // Validação básica do DTO
    if err := req.Validate(); err != nil {
        h.log.WithContext(r.Context()).Info("Bootstrap validação falhou: %v", err)
        response.WriteError(w, http.StatusBadRequest, "Validação falhou", err.Error())
        return
    }

    // Validação de formato de email
    if err := h.validator.IsEmail(req.Email); err != nil {
        h.log.WithContext(r.Context()).Info("Bootstrap email inválido: %v", err)
        response.WriteError(w, http.StatusBadRequest, "Email inválido", err.Error())
        return
    }

    // Validação de força da senha
    if err := h.validator.IsPasswordStrong(req.Senha); err != nil {
        h.log.WithContext(r.Context()).Info("Bootstrap senha fraca: %v", err)
        response.WriteError(w, http.StatusBadRequest, "Senha não atende requisitos de segurança", err.Error())
        return
    }

    // Validação de CNPJ
    if err := h.validator.IsCNPJ(req.CNPJ); err != nil {
        h.log.WithContext(r.Context()).Info("Bootstrap CNPJ inválido: %v", err)
        response.WriteError(w, http.StatusBadRequest, "CNPJ inválido", err.Error())
        return
    }

    // Converter DTO para entidades
    empresa, usuario := req.ToEntities()

    // Preparar dados de bootstrap
    bootstrapData := &usecase.BootstrapData{
        Empresa: empresa,
        Usuario: usuario,
    }

    // Executar bootstrap
    if err := h.bootstrapUseCase.InitializeSystem(r.Context(), bootstrapData); err != nil {
        h.log.WithContext(r.Context()).Error("Bootstrap erro: %v", err)
        
        if strings.Contains(err.Error(), "já inicializado") {
            response.WriteError(w, http.StatusForbidden, "Sistema já inicializado", err.Error())
            return
        }
        
        if strings.Contains(err.Error(), "já cadastrado") {
            response.WriteError(w, http.StatusConflict, "Dados já cadastrados", err.Error())
            return
        }
        
        response.WriteError(w, http.StatusInternalServerError, "Erro ao inicializar sistema", "Erro interno ao processar bootstrap")
        return
    }

    h.log.WithContext(r.Context()).Info(
        "Bootstrap concluído: empresa_id=%d, admin_id=%d, email=%s",
        empresa.ID,
        usuario.ID,
        usuario.Email,
    )
    
    // Retornar confirmação (sem dados sensíveis)
    resp := map[string]interface{}{
        "message":    "Sistema inicializado com sucesso",
        "empresa_id": empresa.ID,
        "admin_id":   usuario.ID,
    }
    
    response.WriteSuccess(w, http.StatusCreated, "Sistema inicializado com sucesso", resp)
}

// RegisterRoutes registra todas as rotas HTTP do handler no roteador
func (h *BootstrapHandler) RegisterRoutes(router *mux.Router) {
    // Rota pública (SEM autenticação) - executável apenas uma vez
    router.HandleFunc("/bootstrap", h.Bootstrap).Methods("POST")
}