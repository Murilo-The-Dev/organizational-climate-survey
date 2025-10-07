// Package usecase implementa os casos de uso para Empresas.
// Fornece funcionalidades de CRUD e validações específicas para empresas.
package usecase

import (
	"context"
	"fmt"
	"organizational-climate-survey/backend/internal/domain/entity"
	"organizational-climate-survey/backend/internal/domain/repository"
	"regexp"
	"strings"
	"time"
)

// EmpresaUseCase implementa casos de uso para gerenciamento de empresas
type EmpresaUseCase struct {
	empresaRepo      repository.EmpresaRepository      // Repositório de empresas
	logAuditoriaRepo repository.LogAuditoriaRepository // Repositório de logs
}

// NewEmpresaUseCase cria uma nova instância do caso de uso de empresas
func NewEmpresaUseCase(empresaRepo repository.EmpresaRepository, logRepo repository.LogAuditoriaRepository) *EmpresaUseCase {
	return &EmpresaUseCase{
		empresaRepo:      empresaRepo,
		logAuditoriaRepo: logRepo,
	}
}

// ValidateCNPJ valida formato básico do CNPJ
func (uc *EmpresaUseCase) ValidateCNPJ(cnpj string) error {
	// Remove caracteres não numéricos
	cnpjNumbers := regexp.MustCompile(`\D`).ReplaceAllString(cnpj, "")

	if len(cnpjNumbers) != 14 {
		return fmt.Errorf("CNPJ deve conter 14 dígitos")
	}

	// Verifica se não são todos iguais (ex: 11111111111111)
	if regexp.MustCompile(`^(\d)\1{13}$`).MatchString(cnpjNumbers) {
		return fmt.Errorf("CNPJ inválido")
	}

	return nil
}

// Create cria uma nova empresa com validações de negócio
func (uc *EmpresaUseCase) Create(ctx context.Context, empresa *entity.Empresa, userAdminID int, enderecoIP string) error {
	// Validações de negócio
	if strings.TrimSpace(empresa.NomeFantasia) == "" {
		return fmt.Errorf("nome fantasia é obrigatório")
	}

	if strings.TrimSpace(empresa.RazaoSocial) == "" {
		return fmt.Errorf("razão social é obrigatória")
	}

	if err := uc.ValidateCNPJ(empresa.CNPJ); err != nil {
		return fmt.Errorf("CNPJ inválido: %v", err)
	}

	// Verifica se CNPJ já existe
	existingEmpresa, err := uc.empresaRepo.GetByCNPJ(ctx, empresa.CNPJ)
	if err == nil && existingEmpresa != nil {
		return fmt.Errorf("empresa com CNPJ %s já cadastrada", empresa.CNPJ)
	}

	// Define data de cadastro
	empresa.DataCadastro = time.Now()

	// Cria a empresa
	if err := uc.empresaRepo.Create(ctx, empresa); err != nil {
		return fmt.Errorf("erro ao criar empresa: %v", err)
	}

	// Log de auditoria
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:   userAdminID,
			TimeStamp:     time.Now(),
			AcaoRealizada: "Empresa Criada",
			Detalhes:      fmt.Sprintf("Empresa criada: %s (ID: %d)", empresa.NomeFantasia, empresa.ID),
			EnderecoIP:    enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}

	return nil
}

// GetByID busca uma empresa pelo seu ID
func (uc *EmpresaUseCase) GetByID(ctx context.Context, id int) (*entity.Empresa, error) {
	if id <= 0 {
		return nil, fmt.Errorf("ID da empresa deve ser maior que zero")
	}

	return uc.empresaRepo.GetByID(ctx, id)
}

// GetByCNPJ busca uma empresa pelo CNPJ
func (uc *EmpresaUseCase) GetByCNPJ(ctx context.Context, cnpj string) (*entity.Empresa, error) {
	if err := uc.ValidateCNPJ(cnpj); err != nil {
		return nil, err
	}

	return uc.empresaRepo.GetByCNPJ(ctx, cnpj)
}

// List lista empresas com paginação
func (uc *EmpresaUseCase) List(ctx context.Context, limit, offset int) ([]*entity.Empresa, error) {
	if limit <= 0 || limit > 100 {
		limit = 20 // Limite padrão
	}
	if offset < 0 {
		offset = 0
	}

	return uc.empresaRepo.List(ctx, limit, offset)
}

// Update atualiza dados de uma empresa existente
func (uc *EmpresaUseCase) Update(ctx context.Context, empresa *entity.Empresa, userAdminID int, enderecoIP string) error {
	// Validações
	if empresa.ID <= 0 {
		return fmt.Errorf("ID da empresa inválido")
	}

	if strings.TrimSpace(empresa.NomeFantasia) == "" {
		return fmt.Errorf("nome fantasia é obrigatório")
	}

	if strings.TrimSpace(empresa.RazaoSocial) == "" {
		return fmt.Errorf("razão social é obrigatória")
	}

	if err := uc.ValidateCNPJ(empresa.CNPJ); err != nil {
		return fmt.Errorf("CNPJ inválido: %v", err)
	}

	// Verifica se empresa existe
	existing, err := uc.empresaRepo.GetByID(ctx, empresa.ID)
	if err != nil {
		return fmt.Errorf("empresa não encontrada: %v", err)
	}

	// Verifica se CNPJ não está sendo usado por outra empresa
	empresaComCNPJ, err := uc.empresaRepo.GetByCNPJ(ctx, empresa.CNPJ)
	if err == nil && empresaComCNPJ != nil && empresaComCNPJ.ID != empresa.ID {
		return fmt.Errorf("CNPJ %s já está sendo usado por outra empresa", empresa.CNPJ)
	}

	// Atualiza
	if err := uc.empresaRepo.Update(ctx, empresa); err != nil {
		return fmt.Errorf("erro ao atualizar empresa: %v", err)
	}

	// Log de auditoria
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:   userAdminID,
			TimeStamp:     time.Now(),
			AcaoRealizada: "Empresa Atualizada",
			Detalhes:      fmt.Sprintf("Empresa atualizada: %s -> %s (ID: %d)", existing.NomeFantasia, empresa.NomeFantasia, empresa.ID),
			EnderecoIP:    enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}

	return nil
}

// Delete remove uma empresa do sistema
func (uc *EmpresaUseCase) Delete(ctx context.Context, id int, userAdminID int, enderecoIP string) error {
	if id <= 0 {
		return fmt.Errorf("ID da empresa inválido")
	}

	// Busca empresa para log
	empresa, err := uc.empresaRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("empresa não encontrada: %v", err)
	}

	// Tenta deletar
	if err := uc.empresaRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("erro ao deletar empresa: %v", err)
	}

	// Log de auditoria
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:   userAdminID,
			TimeStamp:     time.Now(),
			AcaoRealizada: "Empresa Deletada",
			Detalhes:      fmt.Sprintf("Empresa deletada: %s (ID: %d)", empresa.NomeFantasia, empresa.ID),
			EnderecoIP:    enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}

	return nil
}
