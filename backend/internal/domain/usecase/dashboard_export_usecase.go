package usecase

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"organizational-climate-survey/backend/internal/domain/entity"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

type dashboardExportRow struct {
	PerguntaID     int
	PerguntaTexto  string
	TipoPergunta   string
	TotalRespostas int
	Distribuicao   string
	MediaNumerica  string
}

func (uc *DashboardUseCase) generateDashboardReport(
	ctx context.Context,
	dashboard *entity.Dashboard,
	pesquisa *entity.Pesquisa,
	format string,
) ([]byte, error) {
	rows, err := uc.buildDashboardExportRows(ctx, dashboard.IDPesquisa)
	if err != nil {
		return nil, err
	}

	switch format {
	case "csv":
		return generateDashboardCSV(rows)
	case "xlsx":
		return generateDashboardXLSX(rows)
	case "pdf":
		return generateDashboardPDF(pesquisa, rows)
	default:
		return nil, fmt.Errorf("formato de exportação não suportado: %s", format)
	}
}

func (uc *DashboardUseCase) buildDashboardExportRows(ctx context.Context, pesquisaID int) ([]dashboardExportRow, error) {
	perguntas, err := uc.perguntaRepo.ListByPesquisa(ctx, pesquisaID)
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar perguntas para exportação: %v", err)
	}

	rows := make([]dashboardExportRow, 0, len(perguntas))
	for _, pergunta := range perguntas {
		agregados, err := uc.respostaRepo.GetAggregatedByPergunta(ctx, pergunta.ID)
		if err != nil {
			return nil, fmt.Errorf("erro ao agregar respostas da pergunta %d: %v", pergunta.ID, err)
		}

		total := 0
		for _, count := range agregados {
			total += count
		}

		distBytes, _ := json.Marshal(agregados)
		media := ""
		if strings.EqualFold(pergunta.TipoPergunta, "EscalaNumerica") {
			media = computeNumericAverage(agregados)
		}

		rows = append(rows, dashboardExportRow{
			PerguntaID:     pergunta.ID,
			PerguntaTexto:  pergunta.TextoPergunta,
			TipoPergunta:   pergunta.TipoPergunta,
			TotalRespostas: total,
			Distribuicao:   string(distBytes),
			MediaNumerica:  media,
		})
	}

	return rows, nil
}

func generateDashboardCSV(rows []dashboardExportRow) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	header := []string{"pergunta_id", "pergunta", "tipo", "total_respostas", "distribuicao", "media_numerica"}
	if err := w.Write(header); err != nil {
		return nil, err
	}

	for _, row := range rows {
		rec := []string{
			strconv.Itoa(row.PerguntaID),
			row.PerguntaTexto,
			row.TipoPergunta,
			strconv.Itoa(row.TotalRespostas),
			row.Distribuicao,
			row.MediaNumerica,
		}
		if err := w.Write(rec); err != nil {
			return nil, err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func generateDashboardXLSX(rows []dashboardExportRow) ([]byte, error) {
	f := excelize.NewFile()
	sheet := f.GetSheetName(f.GetActiveSheetIndex())

	headers := []string{"pergunta_id", "pergunta", "tipo", "total_respostas", "distribuicao", "media_numerica"}
	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		if err := f.SetCellValue(sheet, cell, header); err != nil {
			return nil, err
		}
	}

	for i, row := range rows {
		values := []interface{}{
			row.PerguntaID,
			row.PerguntaTexto,
			row.TipoPergunta,
			row.TotalRespostas,
			row.Distribuicao,
			row.MediaNumerica,
		}
		for col, value := range values {
			cell, _ := excelize.CoordinatesToCellName(col+1, i+2)
			if err := f.SetCellValue(sheet, cell, value); err != nil {
				return nil, err
			}
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func generateDashboardPDF(pesquisa *entity.Pesquisa, rows []dashboardExportRow) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(190, 10, fmt.Sprintf("Relatório da Pesquisa: %s", pesquisa.Titulo), "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(190, 6, fmt.Sprintf("Gerado em: %s", time.Now().Format("2006-01-02 15:04:05")), "", 1, "L", false, 0, "")
	pdf.Ln(2)

	for _, row := range rows {
		pdf.SetFont("Arial", "B", 10)
		pdf.MultiCell(190, 6, fmt.Sprintf("Pergunta #%d [%s]", row.PerguntaID, row.TipoPergunta), "1", "L", false)
		pdf.SetFont("Arial", "", 10)
		pdf.MultiCell(190, 6, row.PerguntaTexto, "1", "L", false)
		pdf.MultiCell(190, 6, fmt.Sprintf("Total de respostas: %d", row.TotalRespostas), "1", "L", false)
		if row.MediaNumerica != "" {
			pdf.MultiCell(190, 6, fmt.Sprintf("Média numérica: %s", row.MediaNumerica), "1", "L", false)
		}
		pdf.MultiCell(190, 6, "Distribuição: "+row.Distribuicao, "1", "L", false)
		pdf.Ln(2)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func computeNumericAverage(distribuicao map[string]int) string {
	if len(distribuicao) == 0 {
		return ""
	}

	totalWeighted := 0
	totalCount := 0
	for value, count := range distribuicao {
		n, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			continue
		}
		totalWeighted += n * count
		totalCount += count
	}

	if totalCount == 0 {
		return ""
	}

	avg := float64(totalWeighted) / float64(totalCount)
	return fmt.Sprintf("%.2f", avg)
}
