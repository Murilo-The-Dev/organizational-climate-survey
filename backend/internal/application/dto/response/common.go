package response

import (
    "encoding/json"
    "net/http"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success    bool           `json:"success"`
	Message    string         `json:"message,omitempty"`
	Data       interface{}    `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
	Error      string         `json:"error,omitempty"`
}

type PaginationInfo struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

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

func NewErrorResponse(error string) APIResponse {
	return APIResponse{
		Success: false,
		Error:   error,
	}
}

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

// Adicionar estas funções:

func WriteSuccess(w http.ResponseWriter, status int, message string, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    
    response := APIResponse{
        Success: true,
        Message: message,
        Data:    data,
    }
    
    json.NewEncoder(w).Encode(response)
}

func WriteError(w http.ResponseWriter, status int, message, error string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    
    response := APIResponse{
        Success: false,
        Message: message,
        Error:   error,
    }
    
    json.NewEncoder(w).Encode(response)
}

func WritePaginated(w http.ResponseWriter, status int, message string, data interface{}, pagination PaginationInfo) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    
    response := PaginatedResponse{
        Success:    true,
        Message:    message,
        Data:       data,
        Pagination: pagination,
    }
    
    json.NewEncoder(w).Encode(response)
}