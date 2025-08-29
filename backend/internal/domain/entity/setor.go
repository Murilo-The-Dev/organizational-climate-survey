package entity

type Setor struct {
    ID        int    `json:"id_setor"`
    IDEmpresa int    `json:"id_empresa"`
    NomeSetor string `json:"nome_setor"`
    Descricao string `json:"descricao"`
    
    // Relacionamento inverso - opcional
    Empresa *Empresa `json:"empresa,omitempty"`
}