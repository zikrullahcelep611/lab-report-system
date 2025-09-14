package user

import "gorm.io/gorm"

//laboratory technician

type User struct {
	gorm.Model
	Firstname       string `json:"firstname"`
	Lastname   string `json:"lastname"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	HospitalID string `json:"hospital_id" gorm:"type:varchar(7);uniqueIndex;not null"`
}
