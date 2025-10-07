// Package dto define estruturas de transferência de dados (Data Transfer Objects)
// para o fluxo de entrada e saída de informações relacionadas a usuários administradores.
// Este pacote atua como camada intermediária entre os controladores (handlers)
// e o domínio da aplicação, garantindo validação, saneamento e conversão de dados.

package dto

import (
	"organizational-climate-survey/backend/internal/domain/entity"
	"strings"
)

// UsuarioAdministradorCreateRequest representa os dados necessários
// para criação de um novo usuário administrador.
// Inclui validações de integridade e formato em nível de API.
type UsuarioAdministradorCreateRequest struct {
	IDEmpresa int    `json:"id_empresa" binding:"required,gt=0"`                       // Identificador da empresa associada
	NomeAdmin string `json:"nome_admin" binding:"required,min=2,max=255"`              // Nome completo do administrador
	Email     string `json:"email" binding:"required,email,max=255"`                   // E-mail válido do administrador
	Senha     string `json:"senha" binding:"required,min=8,max=128"`                   // Senha em texto plano (antes do hash)
	Status    string `json:"status" binding:"required,oneof=Ativo Inativo Pendente"`   // Estado atual do usuário no sistema
}

// UsuarioAdministradorUpdateRequest representa os campos opcionais
// para atualização de dados de um usuário administrador já existente.
// Os campos são ponteiros para permitir a distinção entre “campo vazio” e “não enviado”.
type UsuarioAdministradorUpdateRequest struct {
	NomeAdmin *string `json:"nome_admin,omitempty" binding:"omitempty,min=2,max=255"`                // Novo nome do administrador (opcional)
	Email     *string `json:"email,omitempty" binding:"omitempty,email,max=255"`                     // Novo e-mail (opcional)
	Status    *string `json:"status,omitempty" binding:"omitempty,oneof=Ativo Inativo Pendente"`     // Novo status (opcional)
}

// UsuarioAdministradorLoginRequest representa as credenciais fornecidas
// para autenticação de um usuário administrador.
type UsuarioAdministradorLoginRequest struct {
	Email string `json:"email" binding:"required,email"`  // E-mail de login
	Senha string `json:"senha" binding:"required"`        // Senha em texto plano
}

// ToEntity converte o DTO de criação em uma entidade de domínio UsuarioAdministrador,
// aplicando normalizações e inserindo o hash da senha fornecido externamente.
func (r *UsuarioAdministradorCreateRequest) ToEntity(senhaHash string) *entity.UsuarioAdministrador {
	return &entity.UsuarioAdministrador{
		IDEmpresa: r.IDEmpresa,
		NomeAdmin: strings.TrimSpace(r.NomeAdmin),
		Email:     strings.ToLower(strings.TrimSpace(r.Email)),
		SenhaHash: senhaHash,
		Status:    r.Status,
	}
}

// ApplyToEntity aplica as modificações do DTO de atualização sobre a entidade existente.
// Somente campos não nulos são atualizados, preservando os valores atuais dos demais atributos.
func (r *UsuarioAdministradorUpdateRequest) ApplyToEntity(user *entity.UsuarioAdministrador) {
	if r.NomeAdmin != nil {
		user.NomeAdmin = strings.TrimSpace(*r.NomeAdmin)
	}
	if r.Email != nil {
		user.Email = strings.ToLower(strings.TrimSpace(*r.Email))
	}
	if r.Status != nil {
		user.Status = *r.Status
	}
}
