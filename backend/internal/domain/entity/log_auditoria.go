package entity

import "time"

type LogAuditoria struct {
    ID            int       `json:"id_log"`
    IDUserAdmin   int       `json:"id_user_admin"`
    TimeStamp     time.Time `json:"timestamp"`
    AcaoRealizada string    `json:"acao_realizada"`
    Detalhes      string    `json:"detalhes"`
    EnderecoIP    string    `json:"endereco_ip"`
    
    // Relacionamento inverso - opcional
    UsuarioAdministrador *UsuarioAdministrador `json:"usuario_administrador,omitempty"`
}