package pdf

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	_ "golang.org/x/image/webp"

	"github.com/go-pdf/fpdf"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"

	domaincar "cvmc/internal/domain/car"
	domainmaint "cvmc/internal/domain/maintenance"
	"cvmc/internal/shared/filecipher"
)

type DecryptFunc func(ownerID, cipherText string) (string, error)

type Builder struct {
	cipher *filecipher.FileCipher
}

func NewBuilder(cipher *filecipher.FileCipher) *Builder {
	return &Builder{cipher: cipher}
}

func (b *Builder) BuildMaintenanceReport(ctx context.Context, car domaincar.Car, maintenances []domainmaint.Maintenance) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(14, 14, 14)
	pdf.SetAutoPageBreak(true, 18)
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	// Color palette (Slate / CVMC Theme)
	primaryR, primaryG, primaryB := 15, 23, 42          // Slate 900
	accentR, accentG, accentB := 37, 99, 235            // Blue 600
	textDarkR, textDarkG, textDarkB := 30, 41, 59       // Slate 800
	textMutedR, textMutedG, textMutedB := 100, 116, 139 // Slate 500
	borderR, borderG, borderB := 226, 232, 240          // Slate 200
	bgLightR, bgLightG, bgLightB := 248, 250, 252       // Slate 50

	// Configure Footer with page numbers
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)
		pdf.SetTextColor(textMutedR, textMutedG, textMutedB)
		footerText := fmt.Sprintf("Como Vai Meu Carro (CVMC) • Relatório gerado em %s • Página %d", time.Now().Format("02/01/2006 às 15:04"), pdf.PageNo())
		pdf.CellFormat(182, 10, tr(footerText), "", 0, "C", false, 0, "")
	})

	pdf.AddPage()

	// Header Banner
	pdf.SetFillColor(primaryR, primaryG, primaryB)
	pdf.Rect(14, 14, 182, 24, "F")

	pdf.SetXY(18, 18)
	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(255, 255, 255)
	pdf.Cell(174, 8, tr("COMO VAI MEU CARRO • DOSSIÊ DE MANUTENÇÕES"))

	pdf.SetXY(18, 26)
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(148, 163, 184) // Slate 400
	pdf.Cell(174, 6, tr("Relatório consolidado de serviços, revisões periódicas e comprovantes de despesas"))

	pdf.SetY(43)

	// Vehicle Overview Card
	pdf.SetFillColor(bgLightR, bgLightG, bgLightB)
	pdf.SetDrawColor(borderR, borderG, borderB)
	pdf.SetLineWidth(0.3)
	pdf.RoundedRect(14, 42, 182, 38, 3, "1234", "FD")

	pdf.SetXY(18, 45)
	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(accentR, accentG, accentB)
	pdf.Cell(174, 6, tr("INFORMAÇÕES DO VEÍCULO"))

	// Compute total investment
	var totalCost float64
	for _, m := range maintenances {
		if m.Cost != nil && *m.Cost > 0 {
			totalCost += *m.Cost
		}
	}

	col0W := 40.0
	col1W := 64.0
	col2W := 36.0
	col3W := 34.0

	pdf.SetFont("Arial", "B", 9)
	pdf.SetTextColor(textMutedR, textMutedG, textMutedB)
	pdf.SetXY(18, 52)
	pdf.Cell(col0W, 5, tr("Veículo:"))
	pdf.Cell(col1W, 5, tr("Fabricante / Modelo:"))
	pdf.Cell(col2W, 5, tr("Ano Fab / Modelo:"))
	pdf.Cell(col3W, 5, tr("Km Atual:"))

	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(textDarkR, textDarkG, textDarkB)
	pdf.SetXY(18, 57)
	carTitle := car.Name
	if carTitle == "" {
		carTitle = car.Model
	}
	makerModel := fmt.Sprintf("%s %s", car.Manufacturer, car.Model)
	yearFabMod := fmt.Sprintf("%d / %d", car.YearManufacture, car.YearModel)
	kmCurrent := fmt.Sprintf("%d km", car.LastMileage)

	pdf.Cell(col0W, 5, fitText(pdf, tr, carTitle, col0W-2))
	pdf.Cell(col1W, 5, fitText(pdf, tr, makerModel, col1W-2))
	pdf.Cell(col2W, 5, fitText(pdf, tr, yearFabMod, col2W-2))
	pdf.Cell(col3W, 5, fitText(pdf, tr, kmCurrent, col3W-2))

	pdf.SetFont("Arial", "B", 9)
	pdf.SetTextColor(textMutedR, textMutedG, textMutedB)
	pdf.SetXY(18, 64)
	pdf.Cell(col0W, 5, tr("Total de Manutenções:"))
	pdf.Cell(col1W, 5, tr("Investimento Total:"))
	pdf.Cell(col2W, 5, tr("Combustível:"))
	pdf.Cell(col3W, 5, tr("Cód. FIPE:"))

	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(textDarkR, textDarkG, textDarkB)
	pdf.SetXY(18, 69)
	totalMaintStr := fmt.Sprintf("%d registro(s)", len(maintenances))
	totalCostStr := fmt.Sprintf("R$ %.2f", totalCost)
	fuelLabel := car.Fuel
	if fuelLabel == "" {
		fuelLabel = "Flex / Não inf."
	}
	fipeLabel := car.FIPECode
	if fipeLabel == "" {
		fipeLabel = "-"
	}

	pdf.Cell(col0W, 5, fitText(pdf, tr, totalMaintStr, col0W-2))
	pdf.Cell(col1W, 5, fitText(pdf, tr, totalCostStr, col1W-2))
	pdf.Cell(col2W, 5, fitText(pdf, tr, fuelLabel, col2W-2))
	pdf.Cell(col3W, 5, fitText(pdf, tr, fipeLabel, col3W-2))

	pdf.SetY(85)

	// Section Title
	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(primaryR, primaryG, primaryB)
	pdf.Cell(182, 7, tr("HISTÓRICO CRONOLÓGICO DE MANUTENÇÕES"))
	pdf.Ln(8)

	// Table Header
	colW := []float64{22, 22, 90, 24, 24} // Total = 182mm

	renderTableHeader := func() {
		pdf.SetFillColor(primaryR, primaryG, primaryB)
		pdf.SetTextColor(255, 255, 255)
		pdf.SetFont("Arial", "B", 8)

		pdf.CellFormat(colW[0], 7, tr("Data"), "1", 0, "C", true, 0, "")
		pdf.CellFormat(colW[1], 7, tr("Km"), "1", 0, "C", true, 0, "")
		pdf.CellFormat(colW[2], 7, tr("Descrição / Serviços"), "1", 0, "L", true, 0, "")
		pdf.CellFormat(colW[3], 7, tr("Comprovantes"), "1", 0, "C", true, 0, "")
		pdf.CellFormat(colW[4], 7, tr("Custo (R$)"), "1", 1, "R", true, 0, "")

		pdf.SetFont("Arial", "", 8)
		pdf.SetTextColor(textDarkR, textDarkG, textDarkB)
	}

	renderTableHeader()

	if len(maintenances) == 0 {
		pdf.SetFont("Arial", "I", 9)
		pdf.SetTextColor(textMutedR, textMutedG, textMutedB)
		pdf.CellFormat(182, 10, tr("Nenhuma manutenção registrada para este veículo até o momento."), "1", 1, "C", false, 0, "")
	} else {
		setRowFill := func(idx int) {
			if idx%2 == 0 {
				pdf.SetFillColor(255, 255, 255)
			} else {
				pdf.SetFillColor(bgLightR, bgLightG, bgLightB)
			}
		}

		for i, m := range maintenances {
			if err := ctx.Err(); err != nil {
				return nil, err
			}

			dateStr := m.Date.Format("02/01/2006")
			kmStr := fmt.Sprintf("%d km", m.Mileage)

			costStr := "-"
			if m.Cost != nil && *m.Cost > 0 {
				costStr = fmt.Sprintf("R$ %.2f", *m.Cost)
			}

			attStr := "-"
			if len(m.Attachments) > 0 {
				attStr = fmt.Sprintf("%d anexo(s)", len(m.Attachments))
			}

			desc := m.Title
			if len(m.Types) > 0 {
				desc = fmt.Sprintf("%s (%s)", desc, strings.Join(m.Types, ", "))
			}
			if m.Description != "" && m.Description != m.Title {
				desc = fmt.Sprintf("%s - %s", desc, m.Description)
			}

			// Pre-calculate line height to avoid mid-row page breaks
			pdf.SetFont("Arial", "", 8)
			lines := pdf.SplitLines([]byte(tr(desc)), colW[2])
			estHeight := float64(len(lines)) * 5.0
			if estHeight < 7.0 {
				estHeight = 7.0
			}

			// Check if new page needed
			if pdf.GetY()+estHeight > 265 {
				pdf.AddPage()
				renderTableHeader()
			}

			// Restore alternating row fill color (fixes dark background on new page)
			setRowFill(i)
			pdf.SetFont("Arial", "", 8)
			pdf.SetTextColor(textDarkR, textDarkG, textDarkB)

			startX := pdf.GetX()
			startY := pdf.GetY()

			// MultiCell for description
			pdf.SetXY(startX+colW[0]+colW[1], startY)
			pdf.MultiCell(colW[2], 5, tr(desc), "1", "L", true)
			rowHeight := pdf.GetY() - startY
			if rowHeight < 7 {
				rowHeight = 7
			}

			// Render other columns with matching row height
			pdf.SetXY(startX, startY)
			pdf.CellFormat(colW[0], rowHeight, tr(dateStr), "1", 0, "C", true, 0, "")
			pdf.CellFormat(colW[1], rowHeight, tr(kmStr), "1", 0, "C", true, 0, "")

			pdf.SetXY(startX+colW[0]+colW[1]+colW[2], startY)
			pdf.CellFormat(colW[3], rowHeight, tr(attStr), "1", 0, "C", true, 0, "")
			pdf.CellFormat(colW[4], rowHeight, tr(costStr), "1", 0, "R", true, 0, "")

			pdf.SetXY(startX, startY+rowHeight)
		}
	}

	// ----------------------------------------------------
	// Processing Attachments & Receipts
	// ----------------------------------------------------
	var pdfReceiptStreams [][]byte
	imgReceiptCount := 0

	for _, m := range maintenances {
		for attIdx, att := range m.Attachments {
			if err := ctx.Err(); err != nil {
				return nil, err
			}

			decryptedDataURL := att.DataUrl
			if b.cipher != nil && filecipher.IsEncrypted(decryptedDataURL) {
				if plain, err := b.cipher.DecryptForOwner(car.OwnerID, decryptedDataURL); err == nil {
					decryptedDataURL = plain
				}
			}

			mimeType, rawBytes, err := parseDataURL(decryptedDataURL)
			if err != nil || len(rawBytes) == 0 {
				continue
			}

			if isPDFMime(mimeType) {
				pdfReceiptStreams = append(pdfReceiptStreams, rawBytes)
				// Transition / bookmark page in main PDF
				pdf.AddPage()
				headerBottomY := b.renderReceiptHeader(pdf, tr, m, att, attIdx+1, true)
				pdf.SetY(headerBottomY + 16.0)
				pdf.SetFont("Arial", "I", 10)
				pdf.SetTextColor(textMutedR, textMutedG, textMutedB)
				pdf.CellFormat(182, 10, tr("O documento original em formato PDF (Nota Fiscal / Ordem de Serviço) foi anexado nas páginas a seguir."), "", 1, "C", false, 0, "")
				continue
			}

			// Handle Image Receipt (JPG, PNG, WEBP, etc.)
			imgCfg, format, err := image.DecodeConfig(bytes.NewReader(rawBytes))
			if err != nil {
				// Try full decode
				img, _, decErr := image.Decode(bytes.NewReader(rawBytes))
				if decErr != nil {
					continue
				}
				imgCfg = image.Config{Width: img.Bounds().Dx(), Height: img.Bounds().Dy()}
				format = "png"
				var buf bytes.Buffer
				if err := png.Encode(&buf, img); err == nil {
					rawBytes = buf.Bytes()
				}
			}

			// If WebP, convert to PNG for fpdf compatibility
			if format == "webp" {
				img, _, decErr := image.Decode(bytes.NewReader(rawBytes))
				if decErr == nil {
					var buf bytes.Buffer
					if err := png.Encode(&buf, img); err == nil {
						rawBytes = buf.Bytes()
						format = "png"
						imgCfg.Width = img.Bounds().Dx()
						imgCfg.Height = img.Bounds().Dy()
					}
				}
			}

			pdf.AddPage()
			headerBottomY := b.renderReceiptHeader(pdf, tr, m, att, attIdx+1, false)

			// Calculate dimensions to fit inside A4 page (max width: 182mm, page bottom limit: 275mm)
			margin := 6.0
			topY := headerBottomY + margin
			availW := 182.0
			availH := 275.0 - topY

			scaleX := availW / float64(imgCfg.Width)
			scaleY := availH / float64(imgCfg.Height)
			scale := scaleX
			if scaleY < scale {
				scale = scaleY
			}

			finalW := float64(imgCfg.Width) * scale
			finalH := float64(imgCfg.Height) * scale
			posX := 14.0 + (availW-finalW)/2.0
			posY := topY + (availH-finalH)/2.0

			imageName := fmt.Sprintf("receipt_%s_%d", m.ID, imgReceiptCount)
			imgReceiptCount++

			imageType := "JPG"
			if strings.EqualFold(format, "png") {
				imageType = "PNG"
			}

			opts := fpdf.ImageOptions{
				ImageType:             imageType,
				ReadDpi:               false,
				AllowNegativePosition: false,
			}
			pdf.RegisterImageOptionsReader(imageName, opts, bytes.NewReader(rawBytes))
			pdf.ImageOptions(imageName, posX, posY, finalW, finalH, false, opts, 0, "")
		}
	}

	var mainPdfBuf bytes.Buffer
	if err := pdf.Output(&mainPdfBuf); err != nil {
		return nil, fmt.Errorf("failed to generate main pdf: %w", err)
	}

	// If no PDF attachments to merge, return main PDF directly
	if len(pdfReceiptStreams) == 0 {
		return mainPdfBuf.Bytes(), nil
	}

	// Merge PDF receipts using pdfcpu
	readers := make([]io.ReadSeeker, 0, len(pdfReceiptStreams)+1)
	readers = append(readers, bytes.NewReader(mainPdfBuf.Bytes()))
	for _, stream := range pdfReceiptStreams {
		readers = append(readers, bytes.NewReader(stream))
	}

	var mergedBuf bytes.Buffer
	conf := model.NewDefaultConfiguration()
	conf.ValidationMode = model.ValidationRelaxed

	if err := api.MergeRaw(ctx, readers, &mergedBuf, false, conf); err != nil {
		// Fallback to main PDF if merge encounters corrupt PDF receipt
		return mainPdfBuf.Bytes(), nil
	}

	return mergedBuf.Bytes(), nil
}

var uuidRegex = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)

func sanitizeAttachmentName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "comprovante"
	}

	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	extClean := strings.TrimPrefix(ext, ".")

	// If the base starts with or is a UUID, shorten it (e.g. 7dd77e08...jpeg)
	if uuidRegex.MatchString(base) {
		if extClean != "" {
			return base[:8] + "..." + extClean
		}
		return base[:8] + "..."
	}

	// If name exceeds 25 characters, truncate gracefully
	const maxLen = 25
	if len(name) > maxLen {
		if len(extClean) > 0 && len(extClean) < 10 {
			maxBaseLen := maxLen - len(extClean) - 3 // 3 for "..."
			if maxBaseLen < 6 {
				maxBaseLen = 6
			}
			runes := []rune(base)
			if len(runes) > maxBaseLen {
				base = string(runes[:maxBaseLen])
			}
			return base + "..." + extClean
		}
		runes := []rune(name)
		return string(runes[:maxLen-3]) + "..."
	}

	return name
}

func fitText(pdf *fpdf.Fpdf, tr func(string) string, text string, maxW float64) string {
	str := tr(text)
	if pdf.GetStringWidth(str) <= maxW {
		return str
	}
	runes := []rune(text)
	for len(runes) > 0 {
		candidate := tr(string(runes) + "...")
		if pdf.GetStringWidth(candidate) <= maxW {
			return candidate
		}
		runes = runes[:len(runes)-1]
	}
	return tr("...")
}

func (b *Builder) renderReceiptHeader(pdf *fpdf.Fpdf, tr func(string) string, m domainmaint.Maintenance, att domainmaint.Attachment, num int, isPDF bool) float64 {
	primaryR, primaryG, primaryB := 15, 23, 42
	accentR, accentG, accentB := 37, 99, 235
	bgLightR, bgLightG, bgLightB := 248, 250, 252
	borderR, borderG, borderB := 226, 232, 240
	textMutedR, textMutedG, textMutedB := 100, 116, 139

	cardX := 14.0
	cardY := 14.0
	cardW := 182.0
	contentX := 18.0
	contentW := 174.0

	// Check title line wrapping
	pdf.SetFont("Arial", "B", 9)
	titleText := fmt.Sprintf("Manutenção: %s", m.Title)
	titleLines := pdf.SplitLines([]byte(tr(titleText)), contentW)

	cardH := 28.0
	hasMultipleLines := len(titleLines) > 1
	if hasMultipleLines {
		cardH = 34.0
	}

	pdf.SetFillColor(bgLightR, bgLightG, bgLightB)
	pdf.SetDrawColor(borderR, borderG, borderB)
	pdf.SetLineWidth(0.3)
	pdf.RoundedRect(cardX, cardY, cardW, cardH, 2, "1234", "FD")

	// 1. Type header
	pdf.SetXY(contentX, 17)
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(accentR, accentG, accentB)
	typeLabel := "COMPROVANTE DE MANUTENÇÃO (IMAGEM)"
	if isPDF {
		typeLabel = "COMPROVANTE DE MANUTENÇÃO (DOCUMENTO PDF)"
	}
	pdf.Cell(contentW, 5, tr(fmt.Sprintf("%s #%d", typeLabel, num)))

	// 2. Maintenance Title
	pdf.SetFont("Arial", "B", 9)
	pdf.SetTextColor(primaryR, primaryG, primaryB)
	if !hasMultipleLines {
		pdf.SetXY(contentX, 23)
		pdf.Cell(contentW, 5, fitText(pdf, tr, titleText, contentW-2))
	} else {
		// First line
		pdf.SetXY(contentX, 23)
		pdf.Cell(contentW, 5, string(titleLines[0]))
		// Second line (if more than 2 lines, truncate with fitText)
		pdf.SetXY(contentX, 28)
		if len(titleLines) == 2 {
			pdf.Cell(contentW, 5, string(titleLines[1]))
		} else {
			line2 := fitText(pdf, func(s string) string { return s }, string(titleLines[1]), contentW-8)
			if !strings.HasSuffix(line2, "...") {
				line2 += "..."
			}
			pdf.Cell(contentW, 5, line2)
		}
	}

	// 3. Metadata line
	metaY := 29.0
	if hasMultipleLines {
		metaY = 34.0
	}
	pdf.SetXY(contentX, metaY)
	pdf.SetFont("Arial", "", 8)
	pdf.SetTextColor(textMutedR, textMutedG, textMutedB)
	costStr := "-"
	if m.Cost != nil && *m.Cost > 0 {
		costStr = fmt.Sprintf("R$ %.2f", *m.Cost)
	}
	attName := sanitizeAttachmentName(att.Name)
	metaStr := fmt.Sprintf("Data: %s  •  Odômetro: %d km  •  Valor: %s  •  Arquivo: %s", m.Date.Format("02/01/2006"), m.Mileage, costStr, attName)
	pdf.Cell(contentW, 5, fitText(pdf, tr, metaStr, contentW-2))

	return cardY + cardH
}

func isPDFMime(mime string) bool {
	m := strings.ToLower(strings.TrimSpace(mime))
	return m == "application/pdf" || strings.HasSuffix(m, "/pdf")
}

func parseDataURL(raw string) (string, []byte, error) {
	trimmed := strings.TrimSpace(raw)
	if strings.HasPrefix(trimmed, "data:") {
		parts := strings.SplitN(trimmed, ",", 2)
		if len(parts) == 2 {
			meta := parts[0]
			b64Part := parts[1]
			mimeType := "image/jpeg"
			if semiIdx := strings.Index(meta, ";"); semiIdx != -1 {
				mimeType = strings.TrimPrefix(meta[:semiIdx], "data:")
			}
			data, err := base64.StdEncoding.DecodeString(b64Part)
			if err != nil {
				return "", nil, err
			}
			return mimeType, data, nil
		}
	}

	// Try direct base64 decode
	data, err := base64.StdEncoding.DecodeString(trimmed)
	if err == nil && len(data) > 0 {
		// Detect mime from magic bytes
		if len(data) >= 4 && string(data[:4]) == "%PDF" {
			return "application/pdf", data, nil
		}
		if len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
			return "image/jpeg", data, nil
		}
		if len(data) >= 4 && data[0] == 0x89 && data[1] == 'P' && data[2] == 'N' && data[3] == 'G' {
			return "image/png", data, nil
		}
		return "application/octet-stream", data, nil
	}

	return "", nil, fmt.Errorf("invalid data url payload")
}
