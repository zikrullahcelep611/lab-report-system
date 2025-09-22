package hospitalrepository

import (
	"context"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/hospital"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Repository struct{
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository{
	return &Repository{DB: db}
}

func (h *Repository) GetHospitalByID(ctx context.Context, id uint) (hospital.Hospital, error){
	var hospitalData hospital.Hospital
    result := h.DB.WithContext(ctx).Where("hospital_id = ?", id).First(&hospitalData)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            log.Warn().Str("operation", "GetHospitalByID").Uint("hospital_id", id).Msg("Hospital not found")
            return hospital.Hospital{}, result.Error
        }
        log.Error().Str("operation", "GetHospitalByID").Err(result.Error).Msg("Failed to get hospital")
        return hospital.Hospital{}, result.Error
    }

    log.Info().Str("operation", "GetHospitalByID").Uint("hospital_id", id).Msg("Hospital retrieved successfully")
    return hospitalData, nil
}

func (h *Repository) CreateHospital(ctx context.Context, newHospital hospital.Hospital) (hospital.Hospital, error){
	result := h.DB.WithContext(ctx).Create(&newHospital)
	if result.Error != nil{
		log.Error().Str("operation", "CreateHospital").Err(result.Error).Msg("Failed to create hospital")
		return hospital.Hospital{}, result.Error
	}

	log.Info().Str("operation", "CreateHospital").Uint("hospital_id", newHospital.ID).Msg("Hospital created successfully")
	return newHospital, nil
}

func (h *Repository) UpdateHospital(ctx context.Context, updateHospital hospital.Hospital) (hospital.Hospital, error){
	
	//Bu yaklaşımda eğer id değeri farklı gelirse yeni record oluşturur
	result := h.DB.WithContext(ctx).Save(&updateHospital)
	if result.Error != nil{
		log.Error().Str("operation", "UpdateHospital").Uint("hospital_id", updateHospital.ID).Err(result.Error).Msg("Failed to update hospital")
		return hospital.Hospital{}, result.Error
	}

	log.Info().Str("operation", "UpdateHospital").Uint("hospital_id", updateHospital.ID).Msg("Hospital updated successfully")
	return updateHospital, nil
}

func (h *Repository) DeleteHospital(ctx context.Context, id uint) error {
	var existingHospital hospital.Hospital
    result := h.DB.WithContext(ctx).First(&existingHospital, id)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            log.Warn().Str("operation", "DeleteHospital").Uint("hospital_id", id).Msg("Hospital not found for deletion")
            return result.Error
        }
        log.Error().Str("operation", "DeleteHospital").Err(result.Error).Msg("Failed to find hospital for deletion")
        return result.Error
    }

    result = h.DB.WithContext(ctx).Delete(&existingHospital)
    if result.Error != nil {
        log.Error().Str("operation", "DeleteHospital").Err(result.Error).Uint("hospital_id", id).Msg("Failed to delete hospital")
        return result.Error
    }

    log.Info().Str("operation", "DeleteHospital").Uint("hospital_id", id).Msg("Hospital deleted successfully")
    return nil
}