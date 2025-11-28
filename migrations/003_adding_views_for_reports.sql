--Views para dashs
CREATE OR REPLACE VIEW vw_pesquisa_resumo AS
SELECT 
    p.id_pesquisa,
    p.titulo,
    p.status,
    COUNT(r.id_resposta) AS total_respostas,
    MIN(r.data_submissao) AS primeira_resposta,
    MAX(r.data_submissao) AS ultima_resposta
FROM pesquisa p
LEFT JOIN pergunta pe ON pe.id_pesquisa = p.id_pesquisa
LEFT JOIN resposta r ON r.id_pergunta = pe.id_pergunta
GROUP BY p.id_pesquisa, p.titulo, p.status;

--Relátorio de respostas por pergunta
CREATE OR REPLACE VIEW vw_respostas_por_pergunta AS
SELECT 
    p.id_pesquisa,
    pe.id_pergunta,
    pe.texto_pergunta,
    pe.tipo_pergunta,
    COUNT(r.id_resposta) AS total_respostas
FROM pesquisa p
JOIN pergunta pe ON pe.id_pesquisa = p.id_pesquisa
LEFT JOIN resposta r ON r.id_pergunta = pe.id_pergunta
GROUP BY p.id_pesquisa, pe.id_pergunta, pe.texto_pergunta, pe.tipo_pergunta;

--Relátorio de satisfacao
CREATE OR REPLACE VIEW vw_satisfacao_media AS
SELECT 
    pe.id_pesquisa,
    pe.id_pergunta,
    pe.texto_pergunta,
    ROUND(AVG(CAST(r.valor_resposta AS NUMERIC)), 2) AS media_resposta
FROM pergunta pe
JOIN resposta r ON r.id_pergunta = pe.id_pergunta
WHERE pe.tipo_pergunta = 'EscalaNumerica'
GROUP BY pe.id_pesquisa, pe.id_pergunta, pe.texto_pergunta;

