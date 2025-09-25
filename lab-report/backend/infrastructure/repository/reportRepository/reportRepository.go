package reportrepository

import (
	"context"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/report"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CreateReport(ctx context.Context, newReport report.Report) (report.Report, error) {
	usr := r.DB.WithContext(ctx).Create(&newReport)
	if usr.Error != nil {
		log.Error().Str("operation", "CreateReport").Err(usr.Error).Msg("Failed to create report")
		return report.Report{}, usr.Error
	}

	log.Info().Str("operation", "CreateReport").Msg("Report created successfully")
	return newReport, nil
}

//bu fonksiyon hasta adi ile arama yapilacak sekilde duzeltilmeli

func (r *Repository) GetReportsWithPatientName(ctx context.Context, firstName string, lastName string) ([]report.Report, error) {
	var rprt []report.Report

	result := r.DB.WithContext(ctx).Preload("Patient").Joins("JOIN patients ON reports.patient_id = patients.id").
		Where("patients.name = ? AND patients.lastname = ?", firstName, lastName).Find(&rprt)

	if result.Error != nil {
		log.Error().Str("operation", "GetReportWithPatientName").Str("firstname", firstName).Str("lastname", lastName).
			Err(result.Error).Msg("Failed to get reports with patient name")
		return []report.Report{}, result.Error
	}

	log.Info().Str("operation", "GetReportWithPatientName").Msg("Get report with patient name successfully")
	return rprt, nil
}

func (r *Repository) GetReportsWithPatientNationalID(ctx context.Context, nationalID string) ([]report.Report, error) {
	var rprt []report.Report
	result := r.DB.WithContext(ctx).Preload("Patient").Joins("JOIN patients ON reports.patient_id = patients.id").
		Where("patients.national_id = ?", nationalID).Find(&rprt)

	if result.Error != nil {
		log.Error().Str("operation", "GetReportWithPatientNationalID").Str("national_id", nationalID).Err(result.Error).
			Msg("Failed to get reports with patient national id")
		return []report.Report{}, result.Error
	}

	log.Info().Str("operation", "GetReportsWithPatientNationalID").Str("national_id", nationalID).Msg("Get reports with patient national id successfully")
	return rprt, nil
}

func (r *Repository) GetAllReportsOrdered(ctx context.Context) ([]report.Report, error) {
	var reports []report.Report
	result := r.DB.WithContext(ctx).Preload("User").Preload("Patient").
		Order("created_at DESC").Find(&reports)

	if result.Error != nil {
		log.Error().Str("operation", "GetAllReportsOrdered").
			Err(result.Error).
			Msg("Failed to get all reports ordered")
		return []report.Report{}, result.Error
	}

	log.Info().Str("operation", "GetAllReportsOrdered").
		Int("count", len(reports)).
		Msg("All reports retrieved successfully (ordered)")

	return reports, nil
}

func (r *Repository) UpdateReport(ctx context.Context, updateReport report.Report) (report.Report, error) {
	result := r.DB.WithContext(ctx).Save(&updateReport)
	if result.Error != nil {
		log.Error().Str("operation", "UpdateReport").Uint("report_id", updateReport.ID).Err(result.Error).
			Msg("Failed to update user")
		return report.Report{}, result.Error
	}

	log.Info().Str("operation", "UpdateReport").Uint("report_id", updateReport.ID).Msg("Report updated successfully")
	return updateReport, nil
}

func (r *Repository) DeleteReport(ctx context.Context, id uint) (bool, error) {
	result := r.DB.WithContext(ctx).Delete(&report.Report{}, id)
	if result.Error != nil {
		log.Error().Str("operation", "DeleteReport").Uint("report_id", id).Err(result.Error).Msg("Failed to delete report")
		return false, result.Error
	}

	log.Info().Str("operation", "DeleteReport").Uint("report_id", id).Msg("Report deleted successfully")
	return true, nil
}

func (r *Repository) GetReportByID(ctx context.Context, id uint) (report.Report, error) {
	var rprt report.Report

	result := r.DB.WithContext(ctx).
		Preload("Patient").
		Preload("User").
		Where("id = ?", id).
		First(&rprt)

	if result.Error != nil {
		log.Error().Str("operation", "GetReportByID").
			Uint("report_id", id).
			Err(result.Error).
			Msg("Failed to get report by ID")
		return report.Report{}, result.Error
	}

	log.Info().Str("operation", "GetReportByID").
		Uint("report_id", id).
		Msg("Report retrieved successfully")

	return rprt, nil
}
