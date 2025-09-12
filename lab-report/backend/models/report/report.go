package report

import (
	"gitbub.com/zikrullahcelep611/lab-report/backend/models/patient"
	"gitbub.com/zikrullahcelep611/lab-report/backend/models/user"
	"gorm.io/gorm"
)

type Report struct {
	gorm.Model
	DiagnosisTitle   string          `json:"diagnosis_title" validate:"required"`
	DiagnosisDetails string          `json:"diagnosis_details" validate:"required"`
	ImagePath        string          `json:"image_path" gorm:"type:varchar(255)" `
	UserID           uint            `json:"user_id"`
	User             user.User       `gorm:"foreignKey:UserID" json:"user"`
	PatientID        uint            `json:"patient_id"`
	Patient          patient.Patient `gorm:"foreignKey:PatientID" json:"patient"`
}
