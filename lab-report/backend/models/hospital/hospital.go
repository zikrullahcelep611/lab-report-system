package hospital

import "gorm.io/gorm"

type Hospital struct {
	gorm.Model
	Name string `json:"name"`
}
