package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	reportusecase "cvmc/internal/application/usecase/report"
	domaincar "cvmc/internal/domain/car"
	domainuser "cvmc/internal/domain/user"
	jwtauth "cvmc/internal/infrastructure/auth/jwt"
	memorycar "cvmc/internal/infrastructure/car/memory"
	memorymaint "cvmc/internal/infrastructure/maintenance/memory"
	memoryreport "cvmc/internal/infrastructure/report/memory"
	memorylocal "cvmc/internal/infrastructure/storage/local"
	"cvmc/internal/shared/filecipher"
	"cvmc/internal/shared/httpx"
)

func setupTestReportHandler(t *testing.T) (*ReportHandler, *reportusecase.Service, *memorycar.Repository, *jwtauth.Provider) {
	tmpDir := t.TempDir()
	storage := memorylocal.New(tmpDir)
	reportRepo := memoryreport.NewRepository()
	carRepo := memorycar.NewRepository()
	maintRepo := memorymaint.NewRepository()
	tokens := jwtauth.NewProvider("test-secret-key-12345678901234567890", "test-refresh-key-12345678901234567890")
	cipher := filecipher.New("cipher-secret")

	svc := reportusecase.NewService(reportRepo, carRepo, maintRepo, storage, nil, cipher)
	handler := NewReportHandler(svc, tokens)
	return handler, svc, carRepo, tokens
}

func TestReportHandler_StartMaintenanceReport(t *testing.T) {
	handler, _, carRepo, tokens := setupTestReportHandler(t)
	ctx := context.Background()

	car, _ := carRepo.Create(ctx, domaincar.Car{
		OwnerID:      "user-1",
		Name:         "Civic",
		Manufacturer: "Honda",
		Model:        "Civic",
	})

	token, _ := tokens.GenerateAccessToken(domainuser.User{ID: "user-1", Email: "user@cvmc.com"})

	body, _ := json.Marshal(StartReportRequest{CarID: car.ID})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reports/maintenance-pdf", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.StartMaintenanceReport(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("expected status 202 Accepted, got %d (body: %s)", w.Code, w.Body.String())
	}

	var env httpx.SuccessEnvelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("failed to parse response envelope: %v", err)
	}
}

func TestReportHandler_ListAndGetJobs(t *testing.T) {
	handler, svc, carRepo, tokens := setupTestReportHandler(t)
	ctx := context.Background()

	car, _ := carRepo.Create(ctx, domaincar.Car{
		OwnerID: "user-1",
		Name:    "Civic",
	})

	job, _ := svc.StartMaintenanceReportJob(ctx, "user-1", car.ID)

	token, _ := tokens.GenerateAccessToken(domainuser.User{ID: "user-1", Email: "user@cvmc.com"})

	// List
	reqList := httptest.NewRequest(http.MethodGet, "/api/v1/reports/jobs", nil)
	reqList.Header.Set("Authorization", "Bearer "+token)
	wList := httptest.NewRecorder()
	handler.ListJobs(wList, reqList)

	if wList.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", wList.Code)
	}

	// Get Individual
	reqGet := httptest.NewRequest(http.MethodGet, "/api/v1/reports/jobs/"+job.ID, nil)
	reqGet.SetPathValue("id", job.ID)
	reqGet.Header.Set("Authorization", "Bearer "+token)
	wGet := httptest.NewRecorder()
	handler.GetJob(wGet, reqGet)

	if wGet.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", wGet.Code)
	}
}

func TestReportHandler_DownloadReport(t *testing.T) {
	handler, svc, carRepo, tokens := setupTestReportHandler(t)
	ctx := context.Background()

	car, _ := carRepo.Create(ctx, domaincar.Car{
		OwnerID: "user-1",
		Name:    "Civic",
	})

	job, _ := svc.StartMaintenanceReportJob(ctx, "user-1", car.ID)
	svc.ProcessJob(ctx, job.ID)
	completed, _ := svc.GetJob(ctx, job.ID, "user-1")

	token, _ := tokens.GenerateAccessToken(domainuser.User{ID: "user-1", Email: "user@cvmc.com"})

	reqDownload := httptest.NewRequest(http.MethodGet, "/api/v1/reports/download?file="+completed.FilePath, nil)
	reqDownload.Header.Set("Authorization", "Bearer "+token)
	wDownload := httptest.NewRecorder()

	handler.DownloadReport(wDownload, reqDownload)

	if wDownload.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body: %s)", wDownload.Code, wDownload.Body.String())
	}
	if wDownload.Header().Get("Content-Type") != "application/pdf" {
		t.Errorf("expected application/pdf, got %s", wDownload.Header().Get("Content-Type"))
	}
}

func TestReportHandler_Unauthorized(t *testing.T) {
	handler, _, _, _ := setupTestReportHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reports/jobs", nil)
	w := httptest.NewRecorder()
	handler.ListJobs(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized, got %d", w.Code)
	}
}
