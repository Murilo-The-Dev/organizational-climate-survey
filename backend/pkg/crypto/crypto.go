package crypto

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword gera o hash de uma senha usando bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("erro ao gerar hash da senha: %w", err)
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash compara uma senha em texto plano com um hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateRandomBytes gera uma sequência de bytes aleatórios
// func GenerateRandomBytes(n int) ([]byte, error) {
// 	// Implementação para gerar bytes aleatórios, útil para tokens, etc.
// 	// Ex: crypto/rand.Read
// 	return nil, fmt.Errorf("não implementado")
// }