// Package usecase implementa os casos de uso para Setores.
// Fornece funcionalidades de CRUD e validações específicas para setores empresariais.
package usecase

import (
	"context"
	"fmt"
	"organizational-climate-survey/backend/internal/domain/entity"
	"organizational-climate-survey/backend/internal/domain/repository"
	"strings"
	"time"
)

// SetorUseCase implementa casos de uso para gerenciamento de setores
type SetorUseCase struct {
	repo             repository.SetorRepository       // Repositório de setores
	empresaRepo      repository.EmpresaRepository     // Repositório de empresas
	logAuditoriaRepo repository.LogAuditoriaRepository // Repositório de logs
}

// NewSetorUseCase cria uma nova instância do caso de uso de setores
func NewSetorUseCase(
	repo repository.SetorRepository,
	empresaRepo repository.EmpresaRepository,
	logRepo repository.LogAuditoriaRepository,
) *SetorUseCase {
	return &SetorUseCase{
		repo:             repo,
		empresaRepo:      empresaRepo,
		logAuditoriaRepo: logRepo,
	}
}

// Create cria um novo setor com validações
func (uc *SetorUseCase) Create(ctx context.Context, setor *entity.Setor, userAdminID int, enderecoIP string) error {
	// Validações básicas
	if setor.IDEmpresa <= 0 {
		return fmt.Errorf("ID da empresa é obrigatório")
	}
	
	if strings.TrimSpace(setor.NomeSetor) == "" {
		return fmt.Errorf("nome do setor é obrigatório")
	}
	
	// Verifica se empresa existe
	_, err := uc.empresaRepo.GetByID(ctx, setor.IDEmpresa)
	if err != nil {
		return fmt.Errorf("empresa não encontrada: %v", err)
	}
	
	// Verifica se já existe setor com mesmo nome na empresa
	existingSetor, err := uc.repo.GetByNome(ctx, setor.IDEmpresa, setor.NomeSetor)
	if err == nil && existingSetor != nil {
		return fmt.Errorf("setor '%s' já existe nesta empresa", setor.NomeSetor)
	}
	
	// Log de auditoria
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:    userAdminID,
			TimeStamp:      time.Now(),
			AcaoRealizada:  "Setor Criado",
			Detalhes:       fmt.Sprintf("Setor criado: %s (ID: %d)", setor.NomeSetor, setor.ID),
			EnderecoIP:     enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}
	
	return nil
}

// GetByID busca um setor pelo seu ID
func (uc *SetorUseCase) GetByID(ctx context.Context, id int) (*entity.Setor, error) {
	if id <= 0 {
		return nil, fmt.Errorf("ID do setor deve ser maior que zero")
	}
	
	return uc.repo.GetByID(ctx, id)
}

// GetByNome busca um setor pelo nome dentro de uma empresa
func (uc *SetorUseCase) GetByNome(ctx context.Context, empresaID int, nome string) (*entity.Setor, error) {
	if empresaID <= 0 {
		return nil, fmt.Errorf("ID da empresa deve ser maior que zero")
	}
	
	if strings.TrimSpace(nome) == "" {
		return nil, fmt.Errorf("nome do setor é obrigatório")
	}
	
	return uc.repo.GetByNome(ctx, empresaID, nome)
}

// ListByEmpresa lista todos os setores de uma empresa
func (uc *SetorUseCase) ListByEmpresa(ctx context.Context, empresaID int) ([]*entity.Setor, error) {
	if empresaID <= 0 {
		return nil, fmt.Errorf("ID da empresa deve ser maior que zero")
	}
	
	return uc.repo.ListByEmpresa(ctx, empresaID)
}

// Update atualiza um setor existente
func (uc *SetorUseCase) Update(ctx context.Context, setor *entity.Setor, userAdminID int, enderecoIP string) error {
	// Validações
	if setor.ID <= 0 {
		return fmt.Errorf("ID do setor inválido")
	}

	if strings.TrimSpace(setor.NomeSetor) == "" {
		return fmt.Errorf("nome do setor é obrigatório")
	}
	
	// Verifica se setor existe
	existing, err := uc.repo.GetByID(ctx, setor.ID)
	if err != nil {
		return fmt.Errorf("setor não encontrado: %v", err)
	}
	
	// Verifica se nome não está sendo usado por outro setor da mesma empresa
	setorComNome, err := uc.repo.GetByNome(ctx, setor.IDEmpresa, setor.NomeSetor)
	if err == nil && setorComNome != nil && setorComNome.ID != setor.ID {
		return fmt.Errorf("nome '%s' já está sendo usado por outro setor", setor.NomeSetor)
	}
	
	if err := uc.repo.Update(ctx, setor); err != nil {
		return fmt.Errorf("erro ao atualizar setor: %v", err)
	}
	
	// Log de auditoria
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:    userAdminID,
			TimeStamp:      time.Now(),
			AcaoRealizada:  "Setor Atualizado",
			Detalhes:       fmt.Sprintf("Setor atualizado: %s -> %s (ID: %d)", existing.NomeSetor, setor.NomeSetor, setor.ID),
			EnderecoIP:     enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}
	
	return nil
}

// Delete remove um setor do sistema
func (uc *SetorUseCase) Delete(ctx context.Context, id int, userAdminID int, enderecoIP string) error {
	if id <= 0 {
		return fmt.Errorf("ID do setor inválido")
	}
	
	// Busca setor para log
	setor, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("setor não encontrado: %v", err)
	}
	
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("erro ao deletar setor: %v", err)
	}
	
	// Log de auditoria
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:    userAdminID,
			TimeStamp:      time.Now(),
			AcaoRealizada:  "Setor Deletado",
			Detalhes:       fmt.Sprintf("Setor deletado: %s (ID: %d)", setor.NomeSetor, setor.ID),
			EnderecoIP:     enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}
	
	return nil
}
