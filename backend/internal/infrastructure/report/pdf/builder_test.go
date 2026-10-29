package pdf

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"testing"
	"time"

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
