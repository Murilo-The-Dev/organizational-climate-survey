package entity

type Pergunta struct {
    ID             int    `json:"id_pergunta"`
    IDPesquisa     int    `json:"id_pesquisa"`
    TextoPergunta  string `json:"texto_pergunta"`
    TipoPergunta   string `json:"tipo_pergunta"` // MultiplaEscolha, RespostaAberta, EscalaNumerica, SimNao
    OrdemExibicao  int    `json:"ordem_exibicao"`
    OpcoesResposta *string `json:"opcoes_resposta"` // JSON string para múltipla escolha
    
    // Relacionamentos - carregados sob demanda
    Respostas []Resposta `json:"respostas,omitempty"`
}