package usecase

import (
	"context"
	"fmt"
	"organizational-climate-survey/backend/internal/domain/entity"
	"organizational-climate-survey/backend/internal/domain/repository"
	"organizational-climate-survey/backend/pkg/crypto"
	"regexp"
	"strings"
	"time"
)

type UsuarioAdministradorUseCase struct {
	repo             repository.UsuarioAdministradorRepository
	empresaRepo      repository.EmpresaRepository
	logAuditoriaRepo repository.LogAuditoriaRepository
	crypto           crypto.CryptoService
}

func NewUsuarioAdministradorUseCase(
	repo repository.UsuarioAdministradorRepository,
	empresaRepo repository.EmpresaRepository,
	logRepo repository.LogAuditoriaRepository,
	cryptoSvc crypto.CryptoService,
) *UsuarioAdministradorUseCase {
	return &UsuarioAdministradorUseCase{
		repo:             repo,
		empresaRepo:      empresaRepo,
		logAuditoriaRepo: logRepo,
		crypto:           cryptoSvc,
	}
}

func (uc *UsuarioAdministradorUseCase) Authenticate(ctx context.Context, email, senha, clientIP string) (*entity.UsuarioAdministrador, error) {
	usuario, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		if uc.logAuditoriaRepo != nil {
			log := &entity.LogAuditoria{
				IDUserAdmin:   0,
				TimeStamp:     time.Now(),
				AcaoRealizada: "Tentativa de Login Falhada",
				Detalhes:      fmt.Sprintf("Tentativa de login com email inexistente: %s", email),
				EnderecoIP:    clientIP,
			}
			uc.logAuditoriaRepo.Create(ctx, log)
		}
		return nil, fmt.Errorf("credenciais inválidas")
	}

	if usuario.Status != "Ativo" {
		if uc.logAuditoriaRepo != nil {
			log := &entity.LogAuditoria{
				IDUserAdmin:   0,
				TimeStamp:     time.Now(),
				AcaoRealizada: "Tentativa de Login - Usuário Inativo",
				Detalhes:      fmt.Sprintf("Tentativa de login com usuário inativo: %s (ID: %d)", email, usuario.ID),
				EnderecoIP:    clientIP,
			}
			uc.logAuditoriaRepo.Create(ctx, log)
		}
		return nil, fmt.Errorf("usuário inativo")
	}

	if !uc.crypto.CheckPasswordHash(senha, usuario.SenhaHash) {
		if uc.logAuditoriaRepo != nil {
			log := &entity.LogAuditoria{
				IDUserAdmin:   0,
				TimeStamp:     time.Now(),
				AcaoRealizada: "Tentativa de Login - Senha Incorreta",
				Detalhes:      fmt.Sprintf("Senha incorreta para usuário: %s (ID: %d)", email, usuario.ID),
				EnderecoIP:    clientIP,
			}
			uc.logAuditoriaRepo.Create(ctx, log)
		}
		return nil, fmt.Errorf("credenciais inválidas")
	}

	if uc.logAuditoriaRepo != nil {
		log := &entity.LogAuditoria{
			IDUserAdmin:   usuario.ID,
			TimeStamp:     time.Now(),
			AcaoRealizada: "Login Realizado",
			Detalhes:      fmt.Sprintf("Login bem-sucedido: %s (ID: %d)", email, usuario.ID),
			EnderecoIP:    clientIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}

	return usuario, nil
}

func (uc *UsuarioAdministradorUseCase) UpdatePassword(ctx context.Context, userID int, newPassword string, adminID int, clientIP string) error {
	if userID <= 0 {
		return fmt.Errorf("ID do usuário inválido")
	}
	
	if strings.TrimSpace(newPassword) == "" {
		return fmt.Errorf("nova senha é obrigatória")
	}
	
	if len(newPassword) < 8 {
		return fmt.Errorf("nova senha deve ter pelo menos 8 caracteres")
	}

	usuario, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("usuário não encontrado: %v", err)
	}

	hashedPassword, err := uc.crypto.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("erro ao gerar hash da senha: %v", err)
	}

	if err := uc.repo.UpdatePassword(ctx, userID, hashedPassword); err != nil {
		return fmt.Errorf("erro ao atualizar senha: %v", err)
	}

	if uc.logAuditoriaRepo != nil && adminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:   adminID,
			TimeStamp:     time.Now(),
			AcaoRealizada: "Senha Atualizada",
			Detalhes:      fmt.Sprintf("Senha atualizada para usuário: %s (ID: %d)", usuario.Email, usuario.ID),
			EnderecoIP:    clientIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}

	return nil
}

func (uc *UsuarioAdministradorUseCase) RequestPasswordReset(ctx context.Context, email, clientIP string) error {
	usuario, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		if uc.logAuditoriaRepo != nil {
			log := &entity.LogAuditoria{
				IDUserAdmin:   0,
				TimeStamp:     time.Now(),
				AcaoRealizada: "Solicitação Reset Senha - Email Inexistente",
				Detalhes:      fmt.Sprintf("Solicitação de reset para email inexistente: %s", email),
				EnderecoIP:    clientIP,
			}
			uc.logAuditoriaRepo.Create(ctx, log)
		}
		return nil
	}

	if usuario.Status != "Ativo" {
		if uc.logAuditoriaRepo != nil {
			log := &entity.LogAuditoria{
				IDUserAdmin:   0,
				TimeStamp:     time.Now(),
				AcaoRealizada: "Solicitação Reset Senha - Usuário Inativo",
				Detalhes:      fmt.Sprintf("Solicitação de reset para usuário inativo: %s (ID: %d)", email, usuario.ID),
				EnderecoIP:    clientIP,
			}
			uc.logAuditoriaRepo.Create(ctx, log)
		}
		return nil
	}

	if uc.logAuditoriaRepo != nil {
		log := &entity.LogAuditoria{
			IDUserAdmin:   usuario.ID,
			TimeStamp:     time.Now(),
			AcaoRealizada: "Solicitação Reset Senha",
			Detalhes:      fmt.Sprintf("Solicitação de reset de senha para: %s (ID: %d)", email, usuario.ID),
			EnderecoIP:    clientIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}

	return nil
}

func (uc *UsuarioAdministradorUseCase) ValidateEmail(email string) error {
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

func (uc *UsuarioAdministradorUseCase) Create(ctx context.Context, usuario *entity.UsuarioAdministrador, userAdminID int, enderecoIP string) error {
	if usuario.IDEmpresa <= 0 {
		return fmt.Errorf("ID da empresa é obrigatório")
	}
	
	if strings.TrimSpace(usuario.NomeAdmin) == "" {
		return fmt.Errorf("nome é obrigatório")
	}
	
	if err := uc.ValidateEmail(usuario.Email); err != nil {
		return err
	}
	
	plainPassword := usuario.SenhaHash
	if strings.TrimSpace(plainPassword) == "" {
		return fmt.Errorf("senha é obrigatória")
	}
	
	if len(plainPassword) < 8 {
		return fmt.Errorf("senha deve ter pelo menos 8 caracteres")
	}
	
	_, err := uc.empresaRepo.GetByID(ctx, usuario.IDEmpresa)
	if err != nil {
		return fmt.Errorf("empresa não encontrada: %v", err)
	}
	
	existingUser, err := uc.repo.GetByEmail(ctx, usuario.Email)
	if err == nil && existingUser != nil {
		return fmt.Errorf("email '%s' já está sendo usado", usuario.Email)
	}
	
	hashedPassword, err := uc.crypto.HashPassword(plainPassword)
	if err != nil {
		return fmt.Errorf("erro ao processar senha: %v", err)
	}
	usuario.SenhaHash = hashedPassword
	
	usuario.DataCadastro = time.Now()
	if usuario.Status == "" {
		usuario.Status = "Ativo"
	}
	
	validStatuses := map[string]bool{
		"Ativo":    true,
		"Inativo":  true,
		"Suspenso": true,
	}
	
	if !validStatuses[usuario.Status] {
		return fmt.Errorf("status inválido: %s", usuario.Status)
	}
	
	if err := uc.repo.Create(ctx, usuario); err != nil {
		return fmt.Errorf("erro ao criar usuário: %v", err)
	}
	
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:   userAdminID,
			TimeStamp:     time.Now(),
			AcaoRealizada: "Usuário Administrador Criado",
			Detalhes:      fmt.Sprintf("Usuário criado: %s (%s) (ID: %d)", usuario.NomeAdmin, usuario.Email, usuario.ID),
			EnderecoIP:    enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}
	
	return nil
}

func (uc *UsuarioAdministradorUseCase) GetByID(ctx context.Context, id int) (*entity.UsuarioAdministrador, error) {
	if id <= 0 {
		return nil, fmt.Errorf("ID do usuário deve ser maior que zero")
	}
	
	return uc.repo.GetByID(ctx, id)
}

func (uc *UsuarioAdministradorUseCase) GetByEmail(ctx context.Context, email string) (*entity.UsuarioAdministrador, error) {
	if err := uc.ValidateEmail(email); err != nil {
		return nil, err
	}
	
	return uc.repo.GetByEmail(ctx, email)
}

func (uc *UsuarioAdministradorUseCase) ListByEmpresa(ctx context.Context, empresaID int) ([]*entity.UsuarioAdministrador, error) {
	if empresaID <= 0 {
		return nil, fmt.Errorf("ID da empresa deve ser maior que zero")
	}
	
	return uc.repo.ListByEmpresa(ctx, empresaID)
}

func (uc *UsuarioAdministradorUseCase) ListByStatus(ctx context.Context, empresaID int, status string) ([]*entity.UsuarioAdministrador, error) {
	if empresaID <= 0 {
		return nil, fmt.Errorf("ID da empresa deve ser maior que zero")
	}
	
	validStatuses := map[string]bool{
		"Ativo":    true,
		"Inativo":  true,
		"Suspenso": true,
	}
	
	if !validStatuses[status] {
		return nil, fmt.Errorf("status inválido: %s", status)
	}
	
	return uc.repo.ListByStatus(ctx, empresaID, status)
}

func (uc *UsuarioAdministradorUseCase) Update(ctx context.Context, usuario *entity.UsuarioAdministrador, userAdminID int, enderecoIP string) error {
	if usuario.ID <= 0 {
		return fmt.Errorf("ID do usuário inválido")
	}
	
	if strings.TrimSpace(usuario.NomeAdmin) == "" {
		return fmt.Errorf("nome é obrigatório")
	}
	
	if err := uc.ValidateEmail(usuario.Email); err != nil {
		return err
	}
	
	existing, err := uc.repo.GetByID(ctx, usuario.ID)
	if err != nil {
		return fmt.Errorf("usuário não encontrado: %v", err)
	}
	
	userComEmail, err := uc.repo.GetByEmail(ctx, usuario.Email)
	if err == nil && userComEmail != nil && userComEmail.ID != usuario.ID {
		return fmt.Errorf("email '%s' já está sendo usado por outro usuário", usuario.Email)
	}
	
	if usuario.Status != "" {
		validStatuses := map[string]bool{
			"Ativo":    true,
			"Inativo":  true,
			"Suspenso": true,
		}
		
		if !validStatuses[usuario.Status] {
			return fmt.Errorf("status inválido: %s", usuario.Status)
		}
	}
	
	if err := uc.repo.Update(ctx, usuario); err != nil {
		return fmt.Errorf("erro ao atualizar usuário: %v", err)
	}
	
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:   userAdminID,
			TimeStamp:     time.Now(),
			AcaoRealizada: "Usuário Administrador Atualizado",
			Detalhes:      fmt.Sprintf("Usuário atualizado: %s -> %s (ID: %d)", existing.NomeAdmin, usuario.NomeAdmin, usuario.ID),
			EnderecoIP:    enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}
	
	return nil
}

func (uc *UsuarioAdministradorUseCase) UpdateStatus(ctx context.Context, id int, status string, userAdminID int, enderecoIP string) error {
	if id <= 0 {
		return fmt.Errorf("ID do usuário inválido")
	}
	
	validStatuses := map[string]bool{
		"Ativo":    true,
		"Inativo":  true,
		"Suspenso": true,
	}
	
	if !validStatuses[status] {
		return fmt.Errorf("status inválido: %s", status)
	}
	
	usuario, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("usuário não encontrado: %v", err)
	}
	
	if err := uc.repo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("erro ao atualizar status: %v", err)
	}
	
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:   userAdminID,
			TimeStamp:     time.Now(),
			AcaoRealizada: "Status Usuário Alterado",
			Detalhes:      fmt.Sprintf("Status alterado de '%s' para '%s' - Usuário: %s (ID: %d)", usuario.Status, status, usuario.NomeAdmin, usuario.ID),
			EnderecoIP:    enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}
	
	return nil
}

func (uc *UsuarioAdministradorUseCase) Delete(ctx context.Context, id int, userAdminID int, enderecoIP string) error {
	if id <= 0 {
		return fmt.Errorf("ID do usuário inválido")
	}
	
	usuario, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("usuário não encontrado: %v", err)
	}
	
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("erro ao deletar usuário: %v", err)
	}
	
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:   userAdminID,
			TimeStamp:     time.Now(),
			AcaoRealizada: "Usuário Administrador Deletado",
			Detalhes:      fmt.Sprintf("Usuário deletado: %s (%s) (ID: %d)", usuario.NomeAdmin, usuario.Email, usuario.ID),
			EnderecoIP:    enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}
	
	return nil
}