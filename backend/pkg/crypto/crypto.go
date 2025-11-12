// Package crypto fornece operações criptográficas seguras incluindo hash de senhas,
// geração de tokens e comparações seguras usando algoritmos padrão da indústria.
package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"organizational-climate-survey/backend/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

// Constantes que definem parâmetros e restrições de segurança
const (
	DefaultBcryptCost = 12 // Custo bcrypt recomendado para 2024+ (equilíbrio entre segurança/performance)
	MinTokenBytes     = 32 // Tamanho mínimo do token para segurança criptográfica
)

// Interface Hasher define operações de hash de senha
// Abstrai o mecanismo de hash subjacente para testabilidade e flexibilidade
type Hasher interface {
	HashPassword(password string) (string, error) // Gera hash seguro a partir da senha
	CheckPasswordHash(password, hash string) bool // Verifica senha contra hash
}

// Interface TokenGenerator define geração segura de tokens e bytes aleatórios
// Usado para tokens de sessão, chaves API e nonces criptográficos
type TokenGenerator interface {
	GenerateToken(length int) (string, error)       // Gera token base64 seguro para URL
	GenerateRandomBytes(length int) ([]byte, error) // Gera bytes aleatórios criptograficamente seguros
}

// CryptoService combina todas as operações criptográficas em uma única interface
// Fornece uma API unificada para toda funcionalidade relacionada à criptografia
type CryptoService interface {
	Hasher
	TokenGenerator
}

// cryptoService implementa CryptoService com custo bcrypt configurável
type cryptoService struct {
	bcryptCost int           // Parâmetro de custo para bcrypt (maior = mais seguro mas mais lento)
	log        logger.Logger // Logger para debug e auditoria
}

// NewCryptoService cria um novo CryptoService com custo bcrypt especificado
// Valida parâmetro de custo contra limites min/max do bcrypt por segurança
func NewCryptoService(bcryptCost int) CryptoService {
	if bcryptCost < bcrypt.MinCost || bcryptCost > bcrypt.MaxCost {
		bcryptCost = DefaultBcryptCost
	}
	return &cryptoService{
		bcryptCost: bcryptCost,
		log:        logger.New(nil),
	}
}

// NewDefaultCryptoService cria um CryptoService com configurações padrão recomendadas
// Adequado para maioria das aplicações em produção
func NewDefaultCryptoService() CryptoService {
	return NewCryptoService(DefaultBcryptCost)
}

// HashPassword gera um hash bcrypt a partir de uma senha em texto plano
// Retorna erro se senha estiver vazia ou bcrypt falhar
func (c *cryptoService) HashPassword(password string) (string, error) {
	if password == "" {
		c.log.Warn("HashPassword: tentativa de hash com senha vazia")
		return "", fmt.Errorf("senha não pode estar vazia")
	}

	c.log.Debug("HashPassword: gerando hash com custo %d", c.bcryptCost)

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), c.bcryptCost)
	if err != nil {
		c.log.Error("HashPassword: erro ao gerar hash - %v", err)
		return "", err
	}

	c.log.Debug("HashPassword: hash gerado com sucesso (len=%d)", len(hashedBytes))
	return string(hashedBytes), nil
}

// CheckPasswordHash verifica uma senha em texto plano contra um hash bcrypt
// Retorna falso para entradas vazias para prevenir ataques de tempo
func (c *cryptoService) CheckPasswordHash(password, hash string) bool {
	if password == "" || hash == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateRandomBytes produz bytes aleatórios criptograficamente seguros
// Usa crypto/rand para qualidade CSPRNG de aleatoriedade
func (c *cryptoService) GenerateRandomBytes(length int) ([]byte, error) {
	if length <= 0 {
		c.log.Warn("GenerateRandomBytes: comprimento inválido %d", length)
		return nil, fmt.Errorf("comprimento deve ser positivo")
	}

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		c.log.Error("GenerateRandomBytes: erro ao gerar bytes aleatórios - %v", err)
		return nil, err
	}

	c.log.Debug("GenerateRandomBytes: %d bytes gerados com sucesso", length)
	return bytes, nil
}

// GenerateToken cria um token codificado em base64 seguro para URL
// Força comprimento mínimo para segurança e usa aleatoriedade criptograficamente segura
func (c *cryptoService) GenerateToken(length int) (string, error) {
	if length < MinTokenBytes {
		c.log.Warn("GenerateToken: comprimento %d menor que mínimo %d", length, MinTokenBytes)
		return "", fmt.Errorf("comprimento do token deve ser pelo menos %d bytes", MinTokenBytes)
	}

	bytes, err := c.GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}

	token := base64.URLEncoding.EncodeToString(bytes)
	c.log.Debug("GenerateToken: token gerado com sucesso (len=%d)", len(token))
	return token, nil
}

// SecureCompare realiza comparação de string em tempo constante para prevenir ataques de tempo
// Essencial para comparar tokens, hashes e outros valores sensíveis à segurança
func SecureCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// SecureCompareBytes realiza comparação de slice de bytes em tempo constante
// Versão de baixo nível do SecureCompare para operações com bytes
func SecureCompareBytes(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}