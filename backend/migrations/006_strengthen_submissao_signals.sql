ALTER TABLE submissao_pesquisa
    ADD COLUMN IF NOT EXISTS user_agent_hash VARCHAR(64),
    ADD COLUMN IF NOT EXISTS accept_language_hash VARCHAR(64),
    ADD COLUMN IF NOT EXISTS fingerprint_composto VARCHAR(64);

CREATE INDEX IF NOT EXISTS idx_submissao_user_agent_hash
    ON submissao_pesquisa(user_agent_hash);

CREATE INDEX IF NOT EXISTS idx_submissao_accept_language_hash
    ON submissao_pesquisa(accept_language_hash);

CREATE INDEX IF NOT EXISTS idx_submissao_fingerprint_composto
    ON submissao_pesquisa(fingerprint_composto);