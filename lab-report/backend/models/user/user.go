package user

import "gorm.io/gorm"

//laboratory technician

type User struct {
	gorm.Model
	Name       string `json:"name"`
	Lastname   string `json:"lastname"`
	HospitalID string `json:"hospital_id" gorm:"type:varchar(7);uniqueIndex;not null"`
}
