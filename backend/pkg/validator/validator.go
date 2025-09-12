package validator

import (
	"fmt"
	"regexp"
	"strings"
)

// IsEmail valida se a string é um formato de email válido
func IsEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return fmt.Errorf("email é obrigatório")
	}
	
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("formato de email inválido")
	}
	return nil
}

// IsCNPJ valida o formato básico de um CNPJ
func IsCNPJ(cnpj string) error {
	cnpjNumbers := regexp.MustCompile(`\\D`).ReplaceAllString(cnpj, "")
	
	if len(cnpjNumbers) != 14 {
		return fmt.Errorf("CNPJ deve conter 14 dígitos")
	}
	
	// Verifica se não são todos iguais (ex: 11111111111111)
	if regexp.MustCompile(`^(\\d)\\1{13}$`).MatchString(cnpjNumbers) {
		return fmt.Errorf("CNPJ inválido")
	}
	
	// Implementação mais completa de validação de CNPJ (dígitos verificadores) pode ser adicionada aqui
	return nil
}

// IsPasswordStrong valida a força da senha (ex: min 8 caracteres)
func IsPasswordStrong(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("senha deve ter pelo menos 8 caracteres")
	}
	// Adicionar mais regras de força de senha aqui (ex: maiúscula, minúscula, número, caractere especial)
	return nil
}

// IsValidStatus verifica se o status fornecido é um dos status válidos
func IsValidStatus(status string, validStatuses []string) bool {
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}