package report

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/report"
	"github.com/gorilla/mux"
)

type ReportService interface {
	CreateReport(ctx context.Context, newReport report.Report) (report.Report, error)
	GetReportsWithPatientName(ctx context.Context, firstName string, lastName string) ([]report.Report, error)
	GetReportsWithPatientNationalID(ctx context.Context, nationalID string) ([]report.Report, error)
	GetAllReportsOrdered(ctx context.Context) ([]report.Report, error)
	UpdateReport(ctx context.Context, updateReport report.Report) (report.Report, error)
	DeleteReport(ctx context.Context, id uint) (bool, error)
	GetReportByID(ctx context.Context, id uint) (report.Report, error)
}

type ReportHandler struct {
	reportService ReportService
}

func NewReportController(reportService ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

func (h *ReportHandler) CreateReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var newReport report.Report
	if err := json.NewDecoder(r.Body).Decode(&newReport); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	createdReport, err := h.reportService.CreateReport(ctx, newReport)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdReport)
}

func (h *ReportHandler) GetReportsWithPatientName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	firstName := r.URL.Query().Get("firstName")
	lastName := r.URL.Query().Get("lastName")

	if firstName == "" || lastName == "" {
        http.Error(w, "firstName and lastName parameters are required", http.StatusBadRequest)
        return
    }

	reports, err := h.reportService.GetReportsWithPatientName(ctx, firstName, lastName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

func (h *ReportHandler) GetReportsWithPatientNationalID(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	nationalId := r.URL.Query().Get("nationalID")
	if nationalId == "" {
        http.Error(w, "nationalID parameter is required", http.StatusBadRequest)
        return
    }

	reports, err := h.reportService.GetReportsWithPatientNationalID(ctx, nationalId)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

func (h *ReportHandler) GetAllReportsOrdered(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	reports, err := h.reportService.GetAllReportsOrdered(ctx)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

func (h *ReportHandler) GetReportByID(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil{
		http.Error(w, "Invalid report ID", http.StatusBadRequest)
		return 
	} 

	reportData, err := h.reportService.GetReportByID(ctx, uint(id))
	if err != nil{
		http.Error(w, err.Error(), http.StatusNotFound)
		return 
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reportData)
}

func (h *ReportHandler) UpdateReport(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil{
		http.Error(w, "Invalid report ID", http.StatusBadRequest)
		return 
	}

	var updateReport report.Report
	if err := json.NewDecoder(r.Body).Decode(&updateReport); err != nil{
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	updateReport.ID = uint(id)
	updatedReport, err := h.reportService.UpdateReport(ctx, updateReport)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedReport)
}

func (h *ReportHandler) DeleteReport(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil{
		http.Error(w, "Invalid report ID", http.StatusBadRequest)
		return 
	}

	deletedReport, err := h.reportService.DeleteReport(ctx, uint(id))
	if err != nil{
		http.Error(w,err.Error(), http.StatusNotFound)
		return 
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"deleted": deletedReport,
		"message": "Report deleted successfully",
	})
}

