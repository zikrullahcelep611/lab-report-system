package patientservice

import (
    "context"

    "gitbub.com/zikrullahcelep611/lab-report/backend/models/patient" // gitbub -> github
)

type PatientRepository interface {
    GetPatientByID(ctx context.Context, id uint) (patient.Patient, error)
    CreatePatient(ctx context.Context, newPatient patient.Patient) (patient.Patient, error)
    UpdatePatient(ctx context.Context, updatePatient patient.Patient) (patient.Patient, error)
    DeletePatient(ctx context.Context, id uint) error
}

type PatientService struct {
    patientRepository PatientRepository
}

func NewPatientService(patientRepository PatientRepository) *PatientService {
    return &PatientService{patientRepository: patientRepository}
}

func (p *PatientService) CreatePatient(ctx context.Context, newPatient patient.Patient) (patient.Patient, error) {
    return p.patientRepository.CreatePatient(ctx, newPatient)
}

func (p *PatientService) GetPatientByID(ctx context.Context, id uint) (patient.Patient, error) {
    return p.patientRepository.GetPatientByID(ctx, id)
}

func (p *PatientService) UpdatePatient(ctx context.Context, updatePatient patient.Patient) (patient.Patient, error) {
    return p.patientRepository.UpdatePatient(ctx, updatePatient)
}

func (p *PatientService) DeletePatient(ctx context.Context, id uint) error {
    return p.patientRepository.DeletePatient(ctx, id)
}