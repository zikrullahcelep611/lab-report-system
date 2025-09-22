package patientrepository

import (
    "context"

    "gitbub.com/zikrullahcelep611/lab-report/backend/models/patient" // gitbub -> github
    "github.com/rs/zerolog/log"
    "gorm.io/gorm"
)

type Repository struct {
    DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
    return &Repository{DB: db}
}

func (p *Repository) GetPatientByID(ctx context.Context, id uint) (patient.Patient, error) {
    var patientData patient.Patient
    result := p.DB.WithContext(ctx).Where("patient_id = ?", id).First(&patientData)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            log.Warn().Str("operation", "GetPatientByID").Uint("patient_id", id).Msg("Patient not found")
            return patient.Patient{}, result.Error
        }
        log.Error().Str("operation", "GetPatientByID").Err(result.Error).Msg("Failed to get patient")
        return patient.Patient{}, result.Error
    }

    log.Info().Str("operation", "GetPatientByID").Uint("patient_id", id).Msg("Patient retrieved successfully")
    return patientData, nil
}

func (p *Repository) CreatePatient(ctx context.Context, newPatient patient.Patient) (patient.Patient, error) {
    result := p.DB.WithContext(ctx).Create(&newPatient)
    if result.Error != nil {
        log.Error().Str("operation", "CreatePatient").Err(result.Error).Msg("Failed to create patient")
        return patient.Patient{}, result.Error
    }

    log.Info().Str("operation", "CreatePatient").Uint("patient_id", newPatient.ID).Msg("Patient created successfully")
    return newPatient, nil
}

func (p *Repository) UpdatePatient(ctx context.Context, updatePatient patient.Patient) (patient.Patient, error) {
    // Bu yaklaşımda eğer id değeri farklı gelirse yeni record oluşturur
    result := p.DB.WithContext(ctx).Save(&updatePatient)
    if result.Error != nil {
        log.Error().Str("operation", "UpdatePatient").Uint("patient_id", updatePatient.ID).Err(result.Error).Msg("Failed to update patient")
        return patient.Patient{}, result.Error
    }

    log.Info().Str("operation", "UpdatePatient").Uint("patient_id", updatePatient.ID).Msg("Patient updated successfully")
    return updatePatient, nil
}

func (p *Repository) DeletePatient(ctx context.Context, id uint) error {
    var existingPatient patient.Patient
    result := p.DB.WithContext(ctx).First(&existingPatient, id)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            log.Warn().Str("operation", "DeletePatient").Uint("patient_id", id).Msg("Patient not found for deletion")
            return result.Error
        }
        log.Error().Str("operation", "DeletePatient").Err(result.Error).Msg("Failed to find patient for deletion")
        return result.Error
    }

    result = p.DB.WithContext(ctx).Delete(&existingPatient)
    if result.Error != nil {
        log.Error().Str("operation", "DeletePatient").Err(result.Error).Uint("patient_id", id).Msg("Failed to delete patient")
        return result.Error
    }

    log.Info().Str("operation", "DeletePatient").Uint("patient_id", id).Msg("Patient deleted successfully")
    return nil
}