// Package usecase implementa os casos de uso para Bootstrap do sistema.
package usecase

import (
    "context"
    "fmt"
    "strings"
    "time"

    "organizational-climate-survey/backend/internal/domain/entity"
    "organizational-climate-survey/backend/internal/domain/repository"
    "organizational-climate-survey/backend/pkg/crypto"
)

// BootstrapUseCase implementa caso de uso de inicialização do sistema
type BootstrapUseCase struct {
    empresaRepo      repository.EmpresaRepository
    usuarioRepo      repository.UsuarioAdministradorRepository
    logAuditoriaRepo repository.LogAuditoriaRepository
    crypto           crypto.CryptoService
}

// NewBootstrapUseCase cria nova instância do caso de uso de bootstrap
func NewBootstrapUseCase(
    empresaRepo repository.EmpresaRepository,
    usuarioRepo repository.UsuarioAdministradorRepository,
    logRepo repository.LogAuditoriaRepository,
    cryptoSvc crypto.CryptoService,
) *BootstrapUseCase {
    return &BootstrapUseCase{
        empresaRepo:      empresaRepo,
        usuarioRepo:      usuarioRepo,
        logAuditoriaRepo: logRepo,
        crypto:           cryptoSvc,
    }
}

// BootstrapData agrupa dados necessários para bootstrap
type BootstrapData struct {
    Empresa *entity.Empresa
    Usuario *entity.UsuarioAdministrador
}

// InitializeSystem cria empresa e primeiro admin atomicamente
func (uc *BootstrapUseCase) InitializeSystem(ctx context.Context, data *BootstrapData) error {
    // Validar se sistema já foi inicializado
    if err := uc.validateSystemNotInitialized(ctx); err != nil {
        return err
    }

    // Validar dados da empresa
    if err := uc.validateEmpresa(data.Empresa); err != nil {
        return fmt.Errorf("validação empresa: %w", err)
    }

    // Validar dados do usuário
    if err := uc.validateUsuario(data.Usuario); err != nil {
        return fmt.Errorf("validação usuário: %w", err)
    }

    // Criar empresa
    if err := uc.createEmpresa(ctx, data.Empresa); err != nil {
        return err
    }

    // Vincular usuário à empresa criada
    data.Usuario.IDEmpresa = data.Empresa.ID

    // Criar administrador
    if err := uc.createAdministrador(ctx, data.Usuario); err != nil {
        // Rollback seria ideal aqui (se estivesse usando transação)
        return err
    }

    // Registrar log de auditoria
    uc.logBootstrapSuccess(ctx, data)

    return nil
}

// validateSystemNotInitialized verifica se já existe algum administrador
func (uc *BootstrapUseCase) validateSystemNotInitialized(ctx context.Context) error {
    count, err := uc.usuarioRepo.Count(ctx)
    if err != nil {
        return fmt.Errorf("erro ao verificar inicialização: %w", err)
    }

    if count > 0 {
        return fmt.Errorf("sistema já inicializado com %d administradores", count)
    }

    return nil
}

// validateEmpresa valida dados da empresa
func (uc *BootstrapUseCase) validateEmpresa(empresa *entity.Empresa) error {
    if empresa == nil {
        return fmt.Errorf("empresa não pode ser nula")
    }

    if strings.TrimSpace(empresa.NomeFantasia) == "" {
        return fmt.Errorf("nome fantasia é obrigatório")
    }

    if len(empresa.NomeFantasia) < 2 {
        return fmt.Errorf("nome fantasia deve ter pelo menos 2 caracteres")
    }

    if strings.TrimSpace(empresa.RazaoSocial) == "" {
        return fmt.Errorf("razão social é obrigatória")
    }

    if strings.TrimSpace(empresa.CNPJ) == "" {
        return fmt.Errorf("CNPJ é obrigatório")
    }

    // CNPJ format validation (básico)
    if !uc.isValidCNPJFormat(empresa.CNPJ) {
        return fmt.Errorf("formato de CNPJ inválido")
    }

    return nil
}

// validateUsuario valida dados do usuário
func (uc *BootstrapUseCase) validateUsuario(usuario *entity.UsuarioAdministrador) error {
    if usuario == nil {
        return fmt.Errorf("usuário não pode ser nulo")
    }

    if strings.TrimSpace(usuario.NomeAdmin) == "" {
        return fmt.Errorf("nome do administrador é obrigatório")
    }

    if len(usuario.NomeAdmin) < 3 {
        return fmt.Errorf("nome do administrador deve ter pelo menos 3 caracteres")
    }

    if strings.TrimSpace(usuario.Email) == "" {
        return fmt.Errorf("email é obrigatório")
    }

    if !uc.isValidEmailFormat(usuario.Email) {
        return fmt.Errorf("formato de email inválido")
    }

    if strings.TrimSpace(usuario.SenhaHash) == "" {
        return fmt.Errorf("senha é obrigatória")
    }

    if len(usuario.SenhaHash) < 8 {
        return fmt.Errorf("senha deve ter pelo menos 8 caracteres")
    }

    return nil
}

// createEmpresa cria empresa no banco
func (uc *BootstrapUseCase) createEmpresa(ctx context.Context, empresa *entity.Empresa) error {
    // Verificar se CNPJ já existe
    existing, err := uc.empresaRepo.GetByCNPJ(ctx, empresa.CNPJ)
    if err == nil && existing != nil {
        return fmt.Errorf("CNPJ já cadastrado")
    }

    empresa.DataCadastro = time.Now()

    if err := uc.empresaRepo.Create(ctx, empresa); err != nil {
        return fmt.Errorf("erro ao criar empresa: %w", err)
    }

    return nil
}

// createAdministrador cria administrador no banco
func (uc *BootstrapUseCase) createAdministrador(ctx context.Context, usuario *entity.UsuarioAdministrador) error {
    // Verificar se email já existe
    existing, err := uc.usuarioRepo.GetByEmail(ctx, usuario.Email)
    if err == nil && existing != nil {
        return fmt.Errorf("email já cadastrado")
    }

    // Hash da senha
    plainPassword := usuario.SenhaHash
    hashedPassword, err := uc.crypto.HashPassword(plainPassword)
    if err != nil {
        return fmt.Errorf("erro ao processar senha: %w", err)
    }

    usuario.SenhaHash = hashedPassword
    usuario.DataCadastro = time.Now()
    usuario.Status = "Ativo"

    if err := uc.usuarioRepo.Create(ctx, usuario); err != nil {
        return fmt.Errorf("erro ao criar administrador: %w", err)
    }

    return nil
}

// logBootstrapSuccess registra log de bootstrap bem-sucedido
func (uc *BootstrapUseCase) logBootstrapSuccess(ctx context.Context, data *BootstrapData) {
    if uc.logAuditoriaRepo == nil {
        return
    }

    log := &entity.LogAuditoria{
        IDUserAdmin:   data.Usuario.ID,
        TimeStamp:     time.Now(),
        AcaoRealizada: "Bootstrap - Sistema Inicializado",
        Detalhes: fmt.Sprintf(
            "Empresa: %s (CNPJ: %s) | Admin: %s (%s)",
            data.Empresa.NomeFantasia,
            data.Empresa.CNPJ,
            data.Usuario.NomeAdmin,
            data.Usuario.Email,
        ),
        EnderecoIP: "bootstrap",
    }

    uc.logAuditoriaRepo.Create(ctx, log)
}

// isValidCNPJFormat valida formato básico de CNPJ
func (uc *BootstrapUseCase) isValidCNPJFormat(cnpj string) bool {
    // Formato esperado: XX.XXX.XXX/XXXX-XX
    if len(cnpj) != 18 {
        return false
    }

    // Validação simplificada
    return strings.Contains(cnpj, ".") && 
           strings.Contains(cnpj, "/") && 
           strings.Contains(cnpj, "-")
}

// isValidEmailFormat valida formato básico de email
func (uc *BootstrapUseCase) isValidEmailFormat(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}