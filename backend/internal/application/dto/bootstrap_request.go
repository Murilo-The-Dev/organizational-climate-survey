// Package dto define objetos de transferência de dados para bootstrap.
package dto

import (
    "fmt"
    "strings"

    "organizational-climate-survey/backend/internal/domain/entity"
)

// BootstrapRequest representa requisição para inicializar o sistema
type BootstrapRequest struct {
    // Dados da empresa
    NomeFantasia string `json:"nome_fantasia"`
    RazaoSocial  string `json:"razao_social"`
    CNPJ         string `json:"cnpj"`
    
    // Dados do administrador
    NomeAdmin string `json:"nome_admin"`
    Email     string `json:"email"`
    Senha     string `json:"senha"`
}

// Validate valida campos obrigatórios
func (r *BootstrapRequest) Validate() error {
    if strings.TrimSpace(r.NomeFantasia) == "" {
        return fmt.Errorf("nome_fantasia é obrigatório")
    }
    
    if strings.TrimSpace(r.RazaoSocial) == "" {
        return fmt.Errorf("razao_social é obrigatória")
    }
    
    if strings.TrimSpace(r.CNPJ) == "" {
        return fmt.Errorf("cnpj é obrigatório")
    }
    
    if strings.TrimSpace(r.NomeAdmin) == "" {
        return fmt.Errorf("nome_admin é obrigatório")
    }
    
    if strings.TrimSpace(r.Email) == "" {
        return fmt.Errorf("email é obrigatório")
    }
    
    if strings.TrimSpace(r.Senha) == "" {
        return fmt.Errorf("senha é obrigatória")
    }
    
    if len(r.Senha) < 8 {
        return fmt.Errorf("senha deve ter pelo menos 8 caracteres")
    }
    
    return nil
}

// ToEntities converte DTO para entidades de domínio
func (r *BootstrapRequest) ToEntities() (*entity.Empresa, *entity.UsuarioAdministrador) {
    empresa := &entity.Empresa{
        NomeFantasia: r.NomeFantasia,
        RazaoSocial:  r.RazaoSocial,
        CNPJ:         r.CNPJ,
    }

    usuario := &entity.UsuarioAdministrador{
        NomeAdmin: r.NomeAdmin,
        Email:     r.Email,
        SenhaHash: r.Senha, // Será hasheada no UseCase
        Status:    "Ativo",
    }

    return empresa, usuario
}