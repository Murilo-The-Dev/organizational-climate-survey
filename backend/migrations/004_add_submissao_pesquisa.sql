CREATE TABLE submissao_pesquisa (
    id_submissao SERIAL PRIMARY KEY,
    id_pesquisa INTEGER NOT NULL REFERENCES pesquisa(id_pesquisa) ON DELETE CASCADE,
    token_acesso VARCHAR(255) UNIQUE NOT NULL,
    ip_hash VARCHAR(64),
    fingerprint_hash VARCHAR(64),
    status VARCHAR(20) DEFAULT 'pendente' CHECK (status IN ('pendente', 'completa', 'expirada')),
    data_criacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    data_expiracao TIMESTAMP NOT NULL,
    data_conclusao TIMESTAMP,
    UNIQUE(id_pesquisa, token_acesso)
);

CREATE INDEX idx_submissao_token ON submissao_pesquisa(token_acesso);
CREATE INDEX idx_submissao_pesquisa ON submissao_pesquisa(id_pesquisa);
CREATE INDEX idx_submissao_status ON submissao_pesquisa(status, data_expiracao);
