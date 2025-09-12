package entity

import "time"

type Resposta struct {
    ID             int       `json:"id_resposta"`
    IDPergunta     int       `json:"id_pergunta"`
    IDPesquisa     int       `json:"id_pesquisa"`
    ValorResposta  string    `json:"valor_resposta"`
    DataResposta   time.Time `json:"data_resposta"`
    DataSubmissao  time.Time `json:"data_submissao"`
    
    // Relacionamentos - carregados sob demanda
    Pergunta *Pergunta `json:"pergunta,omitempty"`
    Pesquisa *Pesquisa `json:"pesquisa,omitempty"`
}