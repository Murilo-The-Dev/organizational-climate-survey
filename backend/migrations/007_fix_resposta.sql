-- Adicionar colunas id_pesquisa e data_resposta na tabela resposta
-- DB-01: Relacionar respostas à pesquisa e rastrear data da resposta

ALTER TABLE resposta 
    ADD COLUMN IF NOT EXISTS id_pesquisa INTEGER NOT NULL REFERENCES pesquisa(id_pesquisa) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS data_resposta TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Criar índices para melhor performance
CREATE INDEX IF NOT EXISTS idx_resposta_pesquisa ON resposta(id_pesquisa);
CREATE INDEX IF NOT EXISTS idx_resposta_data_resposta ON resposta(data_resposta);

-- Comentário sobre as novas colunas
COMMENT ON COLUMN resposta.id_pesquisa IS 'Identificador da pesquisa à qual a resposta pertence (redundante com pergunta.id_pesquisa mas útil para queries)';
COMMENT ON COLUMN resposta.data_resposta IS 'Data e hora em que a resposta foi registrada';
