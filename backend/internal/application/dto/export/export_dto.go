package export

type ExportRequest struct {
	IDPesquisa  int      `json:"id_pesquisa" binding:"required,gt=0"`
	Formato     string   `json:"formato" binding:"required,oneof=excel pdf csv json"`
	Filtros     []string `json:"filtros,omitempty"`
	IncluirGraficos bool `json:"incluir_graficos,omitempty"`
}

type ExportResponse struct {
	FileName    string `json:"file_name"`
	FileURL     string `json:"file_url"`
	FileSize    int64  `json:"file_size"`
	ContentType string `json:"content_type"`
	ExpiresAt   string `json:"expires_at"`
}