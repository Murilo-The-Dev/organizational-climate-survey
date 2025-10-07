// Package response fornece structs e funções para respostas padrão da API.
package response

import (
    "encoding/json"
    "net/http"
)

// APIResponse representa uma resposta padrão da API.
type APIResponse struct {
	Success bool        `json:"success"`           // Indica se a operação foi bem-sucedida
	Message string      `json:"message,omitempty"` // Mensagem opcional para o usuário
	Data    interface{} `json:"data,omitempty"`    // Dados retornados, genéricos
	Error   string      `json:"error,omitempty"`   // Mensagem de erro, se houver
}

// PaginatedResponse representa resposta paginada.
type PaginatedResponse struct {
	Success    bool           `json:"success"`
	Message    string         `json:"message,omitempty"`
	Data       interface{}    `json:"data"`
	Pagination PaginationInfo `json:"pagination"` // Informações de paginação
	Error      string         `json:"error,omitempty"`
}

// PaginationInfo mantém detalhes da paginação.
type PaginationInfo struct {
	Page       int `json:"page"`        // Página atual
	Limit      int `json:"limit"`       // Itens por página
	Total      int `json:"total"`       // Total de itens
	TotalPages int `json:"total_pages"` // Total de páginas
}

// NewSuccessResponse cria resposta de sucesso simples.
func NewSuccessResponse(data interface{}, message ...string) APIResponse {
	msg := "Operação realizada com sucesso"
	if len(message) > 0 {
		msg = message[0]
	}
	return APIResponse{
		Success: true,
		Message: msg,
		Data:    data,
	}
}

// NewErrorResponse cria resposta de erro simples.
func NewErrorResponse(err string) APIResponse {
	return APIResponse{
		Success: false,
		Error:   err,
	}
}

// NewPaginatedResponse cria resposta paginada de sucesso.
func NewPaginatedResponse(data interface{}, pagination PaginationInfo, message ...string) PaginatedResponse {
	msg := "Consulta realizada com sucesso"
	if len(message) > 0 {
		msg = message[0]
	}
	return PaginatedResponse{
		Success:    true,
		Message:    msg,
		Data:       data,
		Pagination: pagination,
	}
}

// WriteSuccess escreve resposta JSON de sucesso no ResponseWriter HTTP.
func WriteSuccess(w http.ResponseWriter, status int, message string, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(APIResponse{
        Success: true,
        Message: message,
        Data:    data,
    })
}

// WriteError escreve resposta JSON de erro no ResponseWriter HTTP.
func WriteError(w http.ResponseWriter, status int, message, err string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(APIResponse{
        Success: false,
        Message: message,
        Error:   err,
    })
}

// WritePaginated escreve resposta JSON paginada no ResponseWriter HTTP.
func WritePaginated(w http.ResponseWriter, status int, message string, data interface{}, pagination PaginationInfo) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(PaginatedResponse{
        Success:    true,
        Message:    message,
        Data:       data,
        Pagination: pagination,
    })
}
