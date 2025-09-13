package postgresdb

import (
	"gitbub.com/zikrullahcelep611/lab-report/backend/models/patient"
	"gitbub.com/zikrullahcelep611/lab-report/backend/models/report"
	"gitbub.com/zikrullahcelep611/lab-report/backend/models/token"
	"gitbub.com/zikrullahcelep611/lab-report/backend/models/user"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func MigrateDatabase(db *gorm.DB) {
	err := db.AutoMigrate(
		&patient.Patient{},
		&report.Report{},
		&user.User{},
		&token.ExpiredTokens{},
	)

	if err != nil {
		log.Fatal().Err(err).Msg("Error migrating models")
		panic(err)
	}
}
