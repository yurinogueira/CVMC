package pdf

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"strings"
	"testing"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"

	domaincar "cvmc/internal/domain/car"
	domainmaint "cvmc/internal/domain/maintenance"
	"cvmc/internal/shared/filecipher"
)

func generateTestPNG(w, h int) string {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, color.RGBA{R: 50, G: 150, B: 250, A: 255})
		}
	}
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
}

func TestBuildMaintenanceReport_EmptyAndWithData(t *testing.T) {
	cipher := filecipher.New("test-secret-key-12345678901234567890")
	builder := NewBuilder(cipher)

	car := domaincar.Car{
		ID:              "car-1",
		OwnerID:         "owner-1",
		Name:            "Civic LXR",
		Manufacturer:    "Honda",
		Model:           "Civic",
		YearManufacture: 2014,
		YearModel:       2015,
		LastMileage:     85000,
		Fuel:            "Flex",
		FIPECode:        "004389-3",
	}

	// 1. Empty maintenances
	pdfBytes, err := builder.BuildMaintenanceReport(context.Background(), car, nil)
	if err != nil {
		t.Fatalf("unexpected error building empty report: %v", err)
	}
	if len(pdfBytes) < 500 {
		t.Fatalf("pdf output too small: %d bytes", len(pdfBytes))
	}
	if string(pdfBytes[:4]) != "%PDF" {
		t.Fatalf("expected valid PDF header, got %q", string(pdfBytes[:4]))
	}

	// 2. Maintenances with Image Attachment
	cost1 := 350.50
	cost2 := 120.00
	pngData := generateTestPNG(200, 200)

	// Encrypt attachment using cipher to verify decryption flow
	encryptedPng, _ := cipher.EncryptForOwner("owner-1", pngData)

	maintenances := []domainmaint.Maintenance{
		{
			ID:          "m-1",
			CarID:       "car-1",
			Title:       "Troca de Óleo e Filtros",
			Description: "Óleo 0W20 Sintético e filtro original Honda",
			Date:        time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC),
			Mileage:     80000,
			Types:       []string{"oil", "oil_filter"},
			Cost:        &cost1,
			Attachments: []domainmaint.Attachment{
				{
					ID:       "att-1",
					Name:     "nota_fiscal.png",
					MimeType: "image/png",
					DataUrl:  encryptedPng,
				},
			},
		},
		{
			ID:          "m-2",
			CarID:       "car-1",
			Title:       "Alinhamento e Balanceamento",
			Description: "Alinhamento 3D das 4 rodas",
			Date:        time.Date(2026, 3, 20, 14, 30, 0, 0, time.UTC),
			Mileage:     83000,
			Types:       []string{"tires"},
			Cost:        &cost2,
		},
	}

	fullPdf, err := builder.BuildMaintenanceReport(context.Background(), car, maintenances)
	if err != nil {
		t.Fatalf("unexpected error building full report: %v", err)
	}
	if len(fullPdf) < 2000 {
		t.Fatalf("expected larger pdf with embedded image, got %d bytes", len(fullPdf))
	}
	if string(fullPdf[:4]) != "%PDF" {
		t.Fatalf("expected valid PDF header, got %q", string(fullPdf[:4]))
	}
}

func TestBuildMaintenanceReport_WithPDFAttachment(t *testing.T) {
	cipher := filecipher.New("test-secret-key-12345678901234567890")
	builder := NewBuilder(cipher)

	car := domaincar.Car{
		ID:           "car-1",
		OwnerID:      "owner-1",
		Name:         "Corolla",
		Manufacturer: "Toyota",
		Model:        "Corolla",
	}

	// First create a mini valid PDF to use as attachment
	miniPdf, err := builder.BuildMaintenanceReport(context.Background(), car, nil)
	if err != nil {
		t.Fatalf("failed to create mini pdf: %v", err)
	}

	pdfDataURL := "data:application/pdf;base64," + base64.StdEncoding.EncodeToString(miniPdf)

	cost := 500.0
	maintenances := []domainmaint.Maintenance{
		{
			ID:      "m-pdf",
			CarID:   "car-1",
			Title:   "Revisão dos 50.000 km",
			Date:    time.Now().UTC(),
			Mileage: 50000,
			Cost:    &cost,
			Attachments: []domainmaint.Attachment{
				{
					ID:       "att-pdf",
					Name:     "nf_eletronica.pdf",
					MimeType: "application/pdf",
					DataUrl:  pdfDataURL,
				},
			},
		},
	}

	merged, err := builder.BuildMaintenanceReport(context.Background(), car, maintenances)
	if err != nil {
		t.Fatalf("failed to build report with merged PDF receipt: %v", err)
	}
	if len(merged) < len(miniPdf) {
		t.Fatalf("expected merged PDF to be larger than single PDF, got %d vs %d", len(merged), len(miniPdf))
	}
}

func TestSanitizeAttachmentName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", "comprovante"},
		{"   ", "comprovante"},
		{"7dd77e08-8fcc-4ba1-9550-0957170817b8.jpeg", "7dd77e08...jpeg"},
		{"7DD77E08-8FCC-4BA1-9550-0957170817B8.PNG", "7DD77E08...PNG"},
		{"7dd77e08-8fcc-4ba1-9550-0957170817b8", "7dd77e08..."},
		{"nota_fiscal.png", "nota_fiscal.png"},
		{"recibo_oficina_autorizada_honda_2026.pdf", "recibo_oficina_auto...pdf"},
	}

	for _, tt := range tests {
		got := sanitizeAttachmentName(tt.input)
		if got != tt.expected {
			t.Errorf("sanitizeAttachmentName(%q) = %q, expected %q", tt.input, got, tt.expected)
		}
		if len(got) > 25 {
			t.Errorf("sanitizeAttachmentName(%q) produced %q with length %d > 25", tt.input, got, len(got))
		}
	}
}

func TestFitText(t *testing.T) {
	doc := fpdf.New("P", "mm", "A4", "")
	doc.AddPage()
	doc.SetFont("Arial", "", 9)
	tr := doc.UnicodeTranslatorFromDescriptor("")

	// 1. Fits within max width
	short := "Civic"
	fittedShort := fitText(doc, tr, short, 50.0)
	if fittedShort != tr(short) {
		t.Errorf("expected %q, got %q", tr(short), fittedShort)
	}

	// 2. Longer than max width gets truncated with ellipsis
	long := "Toyota Corolla XEi 1.8/1.8 Flex 16V Automático Especial 2010"
	maxW := 30.0
	fittedLong := fitText(doc, tr, long, maxW)
	if !strings.HasSuffix(fittedLong, "...") {
		t.Errorf("expected ellipsis at end of %q", fittedLong)
	}
	measuredW := doc.GetStringWidth(fittedLong)
	if measuredW > maxW {
		t.Errorf("fitted text width %.2f exceeds maxW %.2f", measuredW, maxW)
	}

	// 3. Multi-byte UTF-8 string handles rune boundaries safely
	utf8Text := "Manutenção com Óleo e Filtros em Ação Especializada"
	fittedUTF8 := fitText(doc, tr, utf8Text, 25.0)
	if !strings.HasSuffix(fittedUTF8, "...") {
		t.Errorf("expected ellipsis for utf8 text: %q", fittedUTF8)
	}
	if doc.GetStringWidth(fittedUTF8) > 25.0 {
		t.Errorf("fitted utf8 text width %.2f exceeds maxW 25.0", doc.GetStringWidth(fittedUTF8))
	}
}

func TestBuildMaintenanceReport_LongFIPENameAndMultiPage(t *testing.T) {
	cipher := filecipher.New("test-secret-key-12345678901234567890")
	builder := NewBuilder(cipher)

	car := domaincar.Car{
		ID:              "car-long",
		OwnerID:         "owner-1",
		Name:            "Meu Carro Corolla",
		Manufacturer:    "Toyota",
		Model:           "Corolla XEi 1.8/1.8 Flex 16V Aut. Versao Especial FIPE",
		YearManufacture: 2009,
		YearModel:       2010,
		LastMileage:     145000,
		Fuel:            "Gasolina e Álcool (Flex)",
		FIPECode:        "004389-3",
	}

	pngData := generateTestPNG(300, 400)
	cost := 250.0

	// Create 15 maintenances to force multi-page table breakdown
	var maintenances []domainmaint.Maintenance
	for i := 1; i <= 15; i++ {
		var atts []domainmaint.Attachment
		if i == 1 {
			atts = []domainmaint.Attachment{
				{
					ID:       "att-uuid",
					Name:     "7dd77e08-8fcc-4ba1-9550-0957170817b8.jpeg",
					MimeType: "image/png",
					DataUrl:  pngData,
				},
			}
		}
		maintenances = append(maintenances, domainmaint.Maintenance{
			ID:          fmt.Sprintf("m-%d", i),
			CarID:       "car-long",
			Title:       "Óleo de Motor, Filtro do Óleo de Motor, Filtro de Ar do Motor, Filtro de Combustível, Fluido da Transmissão",
			Description: fmt.Sprintf("Serviço preventivo detalhado número %d realizado em oficina credenciada com peças originais e garantia estendida.", i),
			Date:        time.Date(2025, time.Month((i%12)+1), 10, 8, 0, 0, 0, time.UTC),
			Mileage:     100000 + i*3000,
			Types:       []string{"oil", "oil_filter", "air_filter", "fuel_filter", "transmission_fluid"},
			Cost:        &cost,
			Attachments: atts,
		})
	}

	pdfBytes, err := builder.BuildMaintenanceReport(context.Background(), car, maintenances)
	if err != nil {
		t.Fatalf("unexpected error building multi-page report: %v", err)
	}

	if string(pdfBytes[:4]) != "%PDF" {
		t.Fatalf("expected valid PDF header, got %q", string(pdfBytes[:4]))
	}

	// Verify page count is >= 3 (page 1 + page 2 for table, page 3 for attachment)
	conf := model.NewDefaultConfiguration()
	pageCount, err := api.PageCount(context.Background(), bytes.NewReader(pdfBytes), conf)
	if err != nil {
		t.Fatalf("failed to count pages in generated pdf: %v", err)
	}
	if pageCount < 3 {
		t.Errorf("expected at least 3 pages (multi-page table + attachment), got %d", pageCount)
	}
}
