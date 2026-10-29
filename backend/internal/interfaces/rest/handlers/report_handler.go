package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	portauth "cvmc/internal/application/ports/auth"
	reportport "cvmc/internal/application/ports/report"
	reportusecase "cvmc/internal/application/usecase/report"
	"cvmc/internal/shared/httpx"
)

type ReportHandler struct {
	service *reportusecase.Service
	tokens  portauth.TokenService
}

func NewReportHandler(service *reportusecase.Service, tokens ...portauth.TokenService) *ReportHandler {
	var tokenService portauth.TokenService
	if len(tokens) > 0 {
		tokenService = tokens[0]
	}
	return &ReportHandler{
		service: service,
		tokens:  tokenService,
	}
}

func (h *ReportHandler) extractUserID(r *http.Request) string {
	if cookie, err := r.Cookie("cvmc_access_token"); err == nil && cookie.Value != "" {
		if h.tokens != nil {
			claims, err := h.tokens.ParseAccessToken(cookie.Value)
			if err == nil && claims.UserID != "" {
				return claims.UserID
			}
		}
	}
	raw := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
	if raw == "" {
		return ""
	}
	if h.tokens != nil {
		claims, err := h.tokens.ParseAccessToken(raw)
		if err == nil && claims.UserID != "" {
			return claims.UserID
		}
	}
	return ""
}

type StartReportRequest struct {
	CarID string `json:"carId" example:"60c72b2f9b1d8b001c8e4d2a"`
}

// StartMaintenanceReport godoc
// @Summary      Solicitar exportação de relatório de manutenções em PDF
// @Description  Inicia o processamento assíncrono do relatório de manutenções com comprovantes paginados via fila sequencial
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body StartReportRequest true "ID do veículo"
// @Success      202 {object} httpx.SuccessEnvelope
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      403 {object} httpx.ErrorEnvelope
// @Failure      404 {object} httpx.ErrorEnvelope
// @Failure      500 {object} httpx.ErrorEnvelope
// @Router       /api/v1/reports/maintenance-pdf [post]
func (h *ReportHandler) StartMaintenanceReport(w http.ResponseWriter, r *http.Request) {
	actorID := h.extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var req StartReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.CarID) == "" {
		httpx.Error(w, http.StatusBadRequest, "carId is required", nil)
		return
	}

	job, err := h.service.StartMaintenanceReportJob(r.Context(), actorID, req.CarID)
	if err != nil {
		if errors.Is(err, reportusecase.ErrCarNotFound) {
			httpx.Error(w, http.StatusNotFound, "Veículo não encontrado", nil)
			return
		}
		if errors.Is(err, reportusecase.ErrForbidden) {
			httpx.Error(w, http.StatusForbidden, "Você não tem permissão para exportar dados deste veículo", nil)
			return
		}
		httpx.Error(w, http.StatusInternalServerError, fmt.Sprintf("Falha ao iniciar exportação: %v", err), nil)
		return
	}

	httpx.Accepted(w, job)
}

// ListJobs godoc
// @Summary      Listar histórico de relatórios
// @Description  Retorna o histórico paginado de relatórios e exportações geradas para o usuário autenticado
// @Tags         Reports
// @Produce      json
// @Security     BearerAuth
// @Param        carId query string false "Filtrar por ID do veículo"
// @Param        page  query int    false "Página (padrão: 1)"
// @Param        limit query int    false "Itens por página (padrão: 10, max: 100)"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      500 {object} httpx.ErrorEnvelope
// @Router       /api/v1/reports/jobs [get]
func (h *ReportHandler) ListJobs(w http.ResponseWriter, r *http.Request) {
	actorID := h.extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	carID := r.URL.Query().Get("carId")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	filter := reportport.ListFilter{
		UserID: actorID,
		CarID:  carID,
		Page:   page,
		Limit:  limit,
	}

	res, err := h.service.ListJobs(r.Context(), filter)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Falha ao listar relatórios", nil)
		return
	}

	httpx.Success(w, res)
}

// GetJob godoc
// @Summary      Consultar status de um relatório
// @Description  Obtém o status de processamento e metadados de um job de relatório individual
// @Tags         Reports
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID do Job"
// @Success      200 {object} httpx.SuccessEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      404 {object} httpx.ErrorEnvelope
// @Router       /api/v1/reports/jobs/{id} [get]
func (h *ReportHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	actorID := h.extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	jobID := r.PathValue("id")
	if jobID == "" {
		httpx.Error(w, http.StatusBadRequest, "ID do job é obrigatório", nil)
		return
	}

	job, err := h.service.GetJob(r.Context(), jobID, actorID)
	if err != nil || job == nil {
		httpx.Error(w, http.StatusNotFound, "Relatório não encontrado", nil)
		return
	}

	httpx.Success(w, job)
}

// DownloadReport godoc
// @Summary      Baixar arquivo PDF de relatório gerado
// @Description  Permite o download seguro do arquivo de relatório em PDF gerado com validação de isolamento de usuário e expiração de 30 dias
// @Tags         Reports
// @Produce      application/pdf
// @Security     BearerAuth
// @Param        file query string true "Caminho relativo do arquivo de relatório"
// @Success      200 {string} string "Arquivo PDF"
// @Failure      400 {object} httpx.ErrorEnvelope
// @Failure      401 {object} httpx.ErrorEnvelope
// @Failure      403 {object} httpx.ErrorEnvelope
// @Failure      404 {object} httpx.ErrorEnvelope
// @Failure      410 {object} httpx.ErrorEnvelope "Relatório expirado"
// @Failure      500 {object} httpx.ErrorEnvelope
// @Router       /api/v1/reports/download [get]
func (h *ReportHandler) DownloadReport(w http.ResponseWriter, r *http.Request) {
	actorID := h.extractUserID(r)
	if actorID == "" {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	filePath := r.URL.Query().Get("file")
	if filePath == "" {
		httpx.Error(w, http.StatusBadRequest, "Parâmetro 'file' é obrigatório", nil)
		return
	}

	data, err := h.service.GetReportFile(r.Context(), actorID, filePath)
	if err != nil {
		if errors.Is(err, reportusecase.ErrInvalidReportPath) || errors.Is(err, reportusecase.ErrForbidden) {
			httpx.Error(w, http.StatusForbidden, "Acesso não autorizado a este relatório", nil)
			return
		}
		if errors.Is(err, reportusecase.ErrReportExpired) {
			httpx.Error(w, http.StatusGone, "Este relatório expirou após 30 dias e não está mais disponível para download", nil)
			return
		}
		httpx.Error(w, http.StatusNotFound, "Relatório não encontrado", nil)
		return
	}

	fileName := filepath.Base(filePath)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileName))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
