package user

import (
	"gitbub.com/zikrullahcelep611/lab-report/backend/models/hospital"
	"gitbub.com/zikrullahcelep611/lab-report/backend/models/role"
	"gorm.io/gorm"
)

//laboratory technician

type User struct {
	gorm.Model
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`

	HospitalID uint              `json:"hospital_id"`
	Hospital   hospital.Hospital `gorm:"foreignKey:HospitalID" json:"hospital"`

	Role role.Role `json:"role" gorm:"type:varchar(20);default:'user'"`
}
