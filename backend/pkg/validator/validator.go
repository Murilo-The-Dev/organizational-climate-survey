// Package validator implementa funcionalidades de validação para dados da aplicação.
package validator

// Importação dos pacotes necessários
import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Constantes do pacote
const (
	MinPasswordLength = 8  // Comprimento mínimo para senhas
	CNPJLength        = 14 // Comprimento padrão do CNPJ
)

// Expressões regulares pré-compiladas para validações
var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`) // Validação de email
	cnpjRegex  = regexp.MustCompile(`\D`)                                               // Remove caracteres não numéricos do CNPJ
)

// Validator: estrutura principal para validações
type Validator struct{}

// New: construtor que cria uma nova instância do Validator
func New() *Validator {
	return &Validator{}
}

// ValidationError: estrutura para padronização de erros de validação
type ValidationError struct {
	Field   string // Campo que gerou o erro
	Message string // Mensagem descritiva do erro
}

// Error: implementa a interface error para ValidationError
func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// IsEmail: valida o formato de um endereço de email
// Verifica se o email está vazio e se corresponde ao padrão esperado
func (v *Validator) IsEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return ValidationError{"email", "required"}
	}
	if !emailRegex.MatchString(email) {
		return ValidationError{"email", "invalid format"}
	}
	return nil
}

// IsCNPJ: realiza a validação completa de um CNPJ
// Inclui verificação de formato, dígitos repetidos e algoritmo dos dígitos verificadores
func (v *Validator) IsCNPJ(cnpj string) error {
	fmt.Printf("DEBUG: CNPJ recebido: '%s' (len=%d)\n", cnpj, len(cnpj))
	
	cnpjNumbers := cnpjRegex.ReplaceAllString(cnpj, "")
	
	fmt.Printf("DEBUG: CNPJ apenas números: '%s' (len=%d)\n", cnpjNumbers, len(cnpjNumbers))

	if len(cnpjNumbers) != CNPJLength {
		return ValidationError{"cnpj", "must contain 14 digits"}
	}

	if isAllSameDigits(cnpjNumbers) {
		return ValidationError{"cnpj", "invalid sequence"}
	}

	if !isValidCNPJCheckDigits(cnpjNumbers) {
		fmt.Println("DEBUG: Dígitos verificadores inválidos")
		return ValidationError{"cnpj", "invalid check digits"}
	}

	return nil
}

// isAllSameDigits: verifica se todos os dígitos são iguais
// Útil para evitar CNPJs inválidos como "11111111111111"
func isAllSameDigits(s string) bool {
	if len(s) == 0 {
		return false
	}
	first := s[0]
	for _, char := range s {
		if byte(char) != first {
			return false
		}
	}
	return true
}

// isValidCNPJCheckDigits: implementa o algoritmo oficial de validação dos dígitos verificadores do CNPJ
// Utiliza os pesos e cálculos conforme especificação da Receita Federal
func isValidCNPJCheckDigits(cnpj string) bool {
	// Arrays de pesos para cálculo dos dígitos verificadores
	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	// Converte string em array de inteiros
	digits := make([]int, CNPJLength)
	for i, char := range cnpj {
		digits[i], _ = strconv.Atoi(string(char))
	}

	// Calcula primeiro dígito verificador
	sum := 0
	for i := 0; i < 12; i++ {
		sum += digits[i] * weights1[i]
	}
	remainder := sum % 11
	checkDigit1 := 0
	if remainder >= 2 {
		checkDigit1 = 11 - remainder
	}

	if digits[12] != checkDigit1 {
		return false
	}

	// Calcula segundo dígito verificador
	sum = 0
	for i := 0; i < 13; i++ {
		sum += digits[i] * weights2[i]
	}
	remainder = sum % 11
	checkDigit2 := 0
	if remainder >= 2 {
		checkDigit2 = 11 - remainder
	}

	return digits[13] == checkDigit2
}

// IsPasswordStrong: valida a força da senha segundo boas práticas de segurança
// Verifica comprimento mínimo e presença de caracteres maiúsculos, minúsculos, números e especiais
func (v *Validator) IsPasswordStrong(password string) error {
	if len(password) < MinPasswordLength {
		return ValidationError{"password", fmt.Sprintf("must be at least %d characters", MinPasswordLength)}
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool

	// Verifica cada tipo de caractere
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// Retorna erro específico para cada requisito não atendido
	if !hasUpper {
		return ValidationError{"password", "must contain uppercase letter"}
	}
	if !hasLower {
		return ValidationError{"password", "must contain lowercase letter"}
	}
	if !hasNumber {
		return ValidationError{"password", "must contain number"}
	}
	if !hasSpecial {
		return ValidationError{"password", "must contain special character"}
	}

	return nil
}

// IsValidStatus: verifica se um status está presente na lista de status válidos
// Retorna erro com as opções disponíveis caso o status seja inválido
func (v *Validator) IsValidStatus(status string, validStatuses []string) error {
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return nil
		}
	}
	return ValidationError{"status", fmt.Sprintf("must be one of: %v", validStatuses)}
}
