// Package export contém structs usadas para requisições e respostas de exportação de dados.
package export

// ExportRequest define os parâmetros para exportar dados de uma pesquisa.
type ExportRequest struct {
	IDPesquisa      int      `json:"id_pesquisa" binding:"required,gt=0"`             // ID da pesquisa a ser exportada
	Formato         string   `json:"formato" binding:"required,oneof=excel pdf csv json"` // Formato do arquivo de exportação
	Filtros         []string `json:"filtros,omitempty"`                               // Lista de filtros aplicados à exportação
	IncluirGraficos bool     `json:"incluir_graficos,omitempty"`                     // Indica se gráficos devem ser incluídos
}

// ExportResponse representa a resposta após a criação do arquivo de exportação.
type ExportResponse struct {
	FileName    string `json:"file_name"`    // Nome do arquivo gerado
	FileURL     string `json:"file_url"`     // URL para download do arquivo
	FileSize    int64  `json:"file_size"`    // Tamanho do arquivo em bytes
	ContentType string `json:"content_type"` // Tipo MIME do arquivo
	ExpiresAt   string `json:"expires_at"`   // Data/hora de expiração do arquivo
}
