package patient

import "gorm.io/gorm"

type Patient struct {
	gorm.Model
	Name       string `json:"name" validate:"required"`
	Lastname   string `json:"lastname" validate:"required"`
	NationalID string `json:"national_id" gorm:"type:varchar(11);uniqueIndex;not null" validate:"required,len=11,numeric"`
}
