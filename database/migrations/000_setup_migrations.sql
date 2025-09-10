-- migrations/000_setup_migrations.sql
-- Este arquivo cria uma tabela para controlar quais migrations já foram executadas e registrar todas ao longo do projeto

-- Criar tabela schema_migrations
CREATE TABLE IF NOT EXISTS schema_migrations (
    id SERIAL PRIMARY KEY,
    migration_name VARCHAR(255) NOT NULL UNIQUE,
    executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Verifica se o estado incial já foi criado e adiciona as migrations posteriores
INSERT INTO schema_migrations (migration_name) 
VALUES ('001_estado_inicial.sql')
ON CONFLICT (migration_name) DO NOTHING;

-- Verificação do registro
SELECT * FROM schema_migrations;
