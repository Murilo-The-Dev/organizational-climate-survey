// Package entity define as entidades principais do domínio da aplicação.
// Fornece as estruturas de dados para perguntas do sistema de pesquisa.
package entity

// Pergunta representa uma questão individual dentro de uma pesquisa de clima
type Pergunta struct {
    ID             int    `json:"id_pergunta"`       // Identificador único da pergunta
    IDPesquisa     int    `json:"id_pesquisa"`       // ID da pesquisa associada
    TextoPergunta  string `json:"texto_pergunta"`    // Texto exibido ao respondente
    TipoPergunta   string `json:"tipo_pergunta"`     // Tipo de resposta esperada
    OrdemExibicao  int    `json:"ordem_exibicao"`    // Sequência de apresentação
    OpcoesResposta *string `json:"opcoes_resposta"`  // JSON com opções para múltipla escolha
    
    // Relacionamento com respostas (carregamento opcional)
    Respostas []Resposta `json:"respostas,omitempty"` // Respostas coletadas
}