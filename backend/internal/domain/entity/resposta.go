package entity

import "time"

type Resposta struct {
    ID            int       `json:"id_resposta"`
    IDPergunta    int       `json:"id_pergunta"`
    DataSubmissao time.Time `json:"data_submissao"`
    ValorResposta string    `json:"valor_resposta"`
    
    // Sem campos identificadores - garantia de anonimato
    // Não armazenar IP, user-agent, session, etc.
}