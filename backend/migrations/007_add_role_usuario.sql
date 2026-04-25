-- Adicionar coluna role na tabela usuario_administrador
-- DB-02: Suportar diferentes papéis/roles para usuários administradores

ALTER TABLE usuario_administrador 
    ADD COLUMN IF NOT EXISTS role VARCHAR(50) DEFAULT 'admin' CHECK (role IN ('admin', 'editor', 'viewer'));

-- Criar índice para queries por role
CREATE INDEX IF NOT EXISTS idx_usuario_admin_role ON usuario_administrador(role);

-- Comentário sobre a nova coluna
COMMENT ON COLUMN usuario_administrador.role IS 'Papel do usuário administrador: admin (acesso total), editor (editar pesquisas), viewer (apenas leitura)';
