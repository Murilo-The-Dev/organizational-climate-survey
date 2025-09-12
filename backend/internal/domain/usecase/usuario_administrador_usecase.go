package usecase

import (
	"context"
	"fmt"
	"organizational-climate-survey/backend/internal/domain/entity"
	"organizational-climate-survey/backend/internal/domain/repository"
	"regexp"
	"strings"
	"time"
	"golang.org/x/crypto/bcrypt"
)

type UsuarioAdministradorUseCase struct {
	repo             repository.UsuarioAdministradorRepository
	empresaRepo      repository.EmpresaRepository
	logAuditoriaRepo repository.LogAuditoriaRepository
}

func NewUsuarioAdministradorUseCase(
	repo repository.UsuarioAdministradorRepository,
	empresaRepo repository.EmpresaRepository,
	logRepo repository.LogAuditoriaRepository,
) *UsuarioAdministradorUseCase {
	return &UsuarioAdministradorUseCase{
		repo:             repo,
		empresaRepo:      empresaRepo,
		logAuditoriaRepo: logRepo,
	}
}

// Authenticate verifica as credenciais do usuário e retorna os dados se válidas
func (uc *UsuarioAdministradorUseCase) Authenticate(ctx context.Context, email, senha, clientIP string) (*entity.UsuarioAdministrador, error) {
	// Buscar usuário por email
	usuario, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		// Log de tentativa de login falhada
		if uc.logAuditoriaRepo != nil {
			log := &entity.LogAuditoria{
				IDUserAdmin:   0, // Usuário não autenticado
				TimeStamp:     time.Now(),
				AcaoRealizada: "Tentativa de Login Falhada",
				Detalhes:      fmt.Sprintf("Tentativa de login com email inexistente: %s", email),
				EnderecoIP:    clientIP,
			}
			uc.logAuditoriaRepo.Create(ctx, log)
		}
		return nil, fmt.Errorf("credenciais inválidas")
	}

	// Verificar se usuário está ativo
	if usuario.Status != "Ativo" {
		// Log de tentativa com usuário inativo
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

	// Verificar senha
	if !uc.ValidatePassword(senha, usuario.SenhaHash) {
		// Log de senha incorreta
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

	// Log de login bem-sucedido
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

// ValidatePassword compara a senha em texto plano com o hash
func (uc *UsuarioAdministradorUseCase) ValidatePassword(plainPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

// UpdatePassword atualiza a senha do usuário com nova implementação
func (uc *UsuarioAdministradorUseCase) UpdatePassword(ctx context.Context, userID int, newPassword string, adminID int, clientIP string) error {
	// Validações
	if userID <= 0 {
		return fmt.Errorf("ID do usuário inválido")
	}
	
	if strings.TrimSpace(newPassword) == "" {
		return fmt.Errorf("nova senha é obrigatória")
	}
	
	if len(newPassword) < 8 {
		return fmt.Errorf("nova senha deve ter pelo menos 8 caracteres")
	}

	// Verificar se usuário existe
	usuario, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("usuário não encontrado: %v", err)
	}

	// Gerar hash da nova senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("erro ao gerar hash da senha: %v", err)
	}

	// Atualizar senha usando método do repository
	if err := uc.repo.UpdatePassword(ctx, userID, string(hashedPassword)); err != nil {
		return fmt.Errorf("erro ao atualizar senha: %v", err)
	}

	// Log de auditoria
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

// RequestPasswordReset inicia o processo de recuperação de senha
func (uc *UsuarioAdministradorUseCase) RequestPasswordReset(ctx context.Context, email, clientIP string) error {
	// Buscar usuário por email
	usuario, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		// Por segurança, não revelar se o email existe ou não
		// Apenas log interno
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
		return nil // Não retornar erro por segurança
	}

	// Verificar se usuário está ativo
	if usuario.Status != "Ativo" {
		// Log e não processar
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
		return nil // Não retornar erro por segurança
	}

	// Aqui seria implementado o envio de email com token de recuperação
	// Por enquanto, apenas log
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

	// TODO: Implementar envio de email com token
	// - Gerar token seguro
	// - Salvar token no banco com expiração
	// - Enviar email com link de reset

	return nil
}

// ValidateEmail valida formato do email
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
	// Validações básicas
	if usuario.IDEmpresa <= 0 {
		return fmt.Errorf("ID da empresa é obrigatório")
	}
	
	if strings.TrimSpace(usuario.NomeAdmin) == "" {
		return fmt.Errorf("nome é obrigatório")
	}
	
	if err := uc.ValidateEmail(usuario.Email); err != nil {
		return err
	}
	
	if strings.TrimSpace(usuario.SenhaHash) == "" {
		return fmt.Errorf("senha é obrigatória")
	}
	
	// Verifica se empresa existe
	_, err := uc.empresaRepo.GetByID(ctx, usuario.IDEmpresa)
	if err != nil {
		return fmt.Errorf("empresa não encontrada: %v", err)
	}
	
	// Verifica se email já existe
	existingUser, err := uc.repo.GetByEmail(ctx, usuario.Email)
	if err == nil && existingUser != nil {
		return fmt.Errorf("email '%s' já está sendo usado", usuario.Email)
	}
	
	// Define valores padrão
	usuario.DataCadastro = time.Now()
	if usuario.Status == "" {
		usuario.Status = "Ativo"
	}
	
	// Valida status
	validStatuses := map[string]bool{
		"Ativo":     true,
		"Inativo":   true,
		"Suspenso":  true,
	}
	
	if !validStatuses[usuario.Status] {
		return fmt.Errorf("status inválido: %s", usuario.Status)
	}
	
	if err := uc.repo.Create(ctx, usuario); err != nil {
		return fmt.Errorf("erro ao criar usuário: %v", err)
	}
	
	// Log de auditoria
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:    userAdminID,
			TimeStamp:      time.Now(),
			AcaoRealizada:  "Usuário Administrador Criado",
			Detalhes:       fmt.Sprintf("Usuário criado: %s (%s) (ID: %d)", usuario.NomeAdmin, usuario.Email, usuario.ID),
			EnderecoIP:     enderecoIP,
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
		"Ativo":     true,
		"Inativo":   true,
		"Suspenso":  true,
	}
	
	if !validStatuses[status] {
		return nil, fmt.Errorf("status inválido: %s", status)
	}
	
	return uc.repo.ListByStatus(ctx, empresaID, status)
}

func (uc *UsuarioAdministradorUseCase) Update(ctx context.Context, usuario *entity.UsuarioAdministrador, userAdminID int, enderecoIP string) error {
	// Validações
	if usuario.ID <= 0 {
		return fmt.Errorf("ID do usuário inválido")
	}
	
	if strings.TrimSpace(usuario.NomeAdmin) == "" {
		return fmt.Errorf("nome é obrigatório")
	}
	
	if err := uc.ValidateEmail(usuario.Email); err != nil {
		return err
	}
	
	// Verifica se usuário existe
	existing, err := uc.repo.GetByID(ctx, usuario.ID)
	if err != nil {
		return fmt.Errorf("usuário não encontrado: %v", err)
	}
	
	// Verifica se email não está sendo usado por outro usuário
	userComEmail, err := uc.repo.GetByEmail(ctx, usuario.Email)
	if err == nil && userComEmail != nil && userComEmail.ID != usuario.ID {
		return fmt.Errorf("email '%s' já está sendo usado por outro usuário", usuario.Email)
	}
	
	// Valida status se informado
	if usuario.Status != "" {
		validStatuses := map[string]bool{
			"Ativo":     true,
			"Inativo":   true,
			"Suspenso":  true,
		}
		
		if !validStatuses[usuario.Status] {
			return fmt.Errorf("status inválido: %s", usuario.Status)
		}
	}
	
	if err := uc.repo.Update(ctx, usuario); err != nil {
		return fmt.Errorf("erro ao atualizar usuário: %v", err)
	}
	
	// Log de auditoria
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:    userAdminID,
			TimeStamp:      time.Now(),
			AcaoRealizada:  "Usuário Administrador Atualizado",
			Detalhes:       fmt.Sprintf("Usuário atualizado: %s -> %s (ID: %d)", existing.NomeAdmin, usuario.NomeAdmin, usuario.ID),
			EnderecoIP:     enderecoIP,
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
		"Ativo":     true,
		"Inativo":   true,
		"Suspenso":  true,
	}
	
	if !validStatuses[status] {
		return fmt.Errorf("status inválido: %s", status)
	}
	
	// Verifica se usuário existe
	usuario, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("usuário não encontrado: %v", err)
	}
	
	if err := uc.repo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("erro ao atualizar status: %v", err)
	}
	
	// Log de auditoria
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:    userAdminID,
			TimeStamp:      time.Now(),
			AcaoRealizada:  "Status Usuário Alterado",
			Detalhes:       fmt.Sprintf("Status alterado de '%s' para '%s' - Usuário: %s (ID: %d)", usuario.Status, status, usuario.NomeAdmin, usuario.ID),
			EnderecoIP:     enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}
	
	return nil
}

func (uc *UsuarioAdministradorUseCase) Delete(ctx context.Context, id int, userAdminID int, enderecoIP string) error {
	if id <= 0 {
		return fmt.Errorf("ID do usuário inválido")
	}
	
	// Busca usuário para log
	usuario, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("usuário não encontrado: %v", err)
	}
	
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("erro ao deletar usuário: %v", err)
	}
	
	// Log de auditoria
	if userAdminID > 0 {
		log := &entity.LogAuditoria{
			IDUserAdmin:    userAdminID,
			TimeStamp:      time.Now(),
			AcaoRealizada:  "Usuário Administrador Deletado",
			Detalhes:       fmt.Sprintf("Usuário deletado: %s (%s) (ID: %d)", usuario.NomeAdmin, usuario.Email, usuario.ID),
			EnderecoIP:     enderecoIP,
		}
		uc.logAuditoriaRepo.Create(ctx, log)
	}
	
	return nil
}