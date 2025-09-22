package hospitalservice

import (
	"context"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/hospital"
)

type HospitalRepository interface {
	GetHospitalByID(ctx context.Context, id uint) (hospital.Hospital, error)
	CreateHospital(ctx context.Context, newHospital hospital.Hospital) (hospital.Hospital, error)
	UpdateHospital(ctx context.Context, updateHospital hospital.Hospital) (hospital.Hospital, error)
	DeleteHospital(ctx context.Context, id uint) error
}

type HospitalService struct {
	hospitalRepository HospitalRepository
}

func NewHospitalService(hospitalRepository HospitalRepository) *HospitalService {
	return &HospitalService{hospitalRepository: hospitalRepository}
}

func (h *HospitalService) CreateHospital(ctx context.Context, newHospital hospital.Hospital) (hospital.Hospital, error) {
	return h.hospitalRepository.CreateHospital(ctx, newHospital)
}

func (h *HospitalService) GetHospitalByID(ctx context.Context, id uint) (hospital.Hospital, error) {
	return h.hospitalRepository.GetHospitalByID(ctx, id)
}

func (h *HospitalService) UpdateHospital(ctx context.Context, updateHospital hospital.Hospital) (hospital.Hospital, error) {
	return h.hospitalRepository.UpdateHospital(ctx, updateHospital)
}

func (h *HospitalService) DeleteHospital(ctx context.Context, id uint) error {
	return h.hospitalRepository.DeleteHospital(ctx, id)
}
