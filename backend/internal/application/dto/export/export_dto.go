// Package export contém structs usadas para requisições e respostas de exportação de dados.
package export

// ExportRequest define os parâmetros para exportar dados de uma pesquisa.
type ExportRequest struct {
	IDPesquisa      int      `json:"id_pesquisa" binding:"required,gt=0" example:"10"`                  // ID da pesquisa a ser exportada
	Formato         string   `json:"formato" binding:"required,oneof=excel pdf csv json" example:"pdf"` // Formato do arquivo de exportação
	Filtros         []string `json:"filtros,omitempty" example:"setor=RH,periodo=Q1"`                   // Lista de filtros aplicados à exportação
	IncluirGraficos bool     `json:"incluir_graficos,omitempty" example:"true"`                         // Indica se gráficos devem ser incluídos
}

// ExportResponse representa a resposta após a criação do arquivo de exportação.
type ExportResponse struct {
	FileName    string `json:"file_name" example:"relatorio_pesquisa_10.pdf"`                                // Nome do arquivo gerado
	FileURL     string `json:"file_url" example:"https://api.exemplo.com/exports/relatorio_pesquisa_10.pdf"` // URL para download do arquivo
	FileSize    int64  `json:"file_size" example:"245760"`                                                   // Tamanho do arquivo em bytes
	ContentType string `json:"content_type" example:"application/pdf"`                                       // Tipo MIME do arquivo
	ExpiresAt   string `json:"expires_at" example:"2026-02-01T10:00:00Z"`                                    // Data/hora de expiração do arquivo
}
