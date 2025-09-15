package reportservice

import (
    "context"

    "gitbub.com/zikrullahcelep611/lab-report/backend/models/report" // gitbub -> github
)

type ReportRepository interface{
    CreateReport(ctx context.Context, newReport report.Report) (report.Report, error)
    GetReportsWithPatientName(ctx context.Context, firstName string, lastName string) ([]report.Report, error)
    GetReportsWithPatientNationalID(ctx context.Context, nationalID string) ([]report.Report, error)
    GetAllReportsOrdered(ctx context.Context) ([]report.Report, error)
    UpdateReport(ctx context.Context, updateReport report.Report) (report.Report, error)
    DeleteReport(ctx context.Context, id uint) (bool, error)
    GetReportByID(ctx context.Context, id uint) (report.Report, error)
}

type ReportService struct{
    reportRepository ReportRepository
}

func NewReportService(reportRepository ReportRepository) *ReportService{
    return &ReportService{reportRepository: reportRepository}
}

func (r *ReportService) CreateReport(ctx context.Context, newReport report.Report) (report.Report, error){
    return r.reportRepository.CreateReport(ctx, newReport)
}

func (r *ReportService) GetReportWithPatientName(ctx context.Context, firstName string, lastName string) ([]report.Report, error){
    return r.reportRepository.GetReportsWithPatientName(ctx, firstName, lastName)
}

func (r *ReportService) GetReportsWithPatientNationalID(ctx context.Context, nationalID string) ([]report.Report, error){
    return r.reportRepository.GetReportsWithPatientNationalID(ctx, nationalID)
}

func (r *ReportService) GetAllReportsOrdered(ctx context.Context) ([]report.Report, error){
    return r.reportRepository.GetAllReportsOrdered(ctx)
}

func (r *ReportService) UpdateReport(ctx context.Context, updateReport report.Report) (report.Report, error){
    return r.reportRepository.UpdateReport(ctx, updateReport)
}

func (r *ReportService) DeleteReport(ctx context.Context, id uint) (bool, error){
    return r.reportRepository.DeleteReport(ctx, id)
}

func (r *ReportService) GetReportByID(ctx context.Context, id uint) (report.Report, error){
    return r.reportRepository.GetReportByID(ctx, id)
}