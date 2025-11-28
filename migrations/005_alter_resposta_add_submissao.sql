-- Adicionar coluna id_submissao
ALTER TABLE resposta 
    ADD COLUMN id_submissao INTEGER REFERENCES submissao_pesquisa(id_submissao) ON DELETE CASCADE;

-- Remover coluna redundante data_resposta
ALTER TABLE resposta DROP COLUMN IF EXISTS data_resposta;

-- Remover coluna redundante id_pesquisa (pergunta já tem)
-- CUIDADO: Se já tem dados, migrar primeiro
-- ALTER TABLE resposta DROP COLUMN id_pesquisa;

CREATE INDEX idx_resposta_submissao ON resposta(id_submissao);
