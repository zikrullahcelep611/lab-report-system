package main

import (
	"gitbub.com/zikrullahcelep611/lab-report/backend/infrastructure/config"
	"gitbub.com/zikrullahcelep611/lab-report/backend/infrastructure/email"
	postgresDb2 "gitbub.com/zikrullahcelep611/lab-report/backend/infrastructure/postgresDb"
	"github.com/rs/zerolog"
)



func main() {
	configModel := config.SetConfig("resources")
	zerolog.SetGlobalLevel(configModel.Log.Level)

	db := postgresDb2.ConnectDatabase(configModel.Database)
	//redis := redis2.ConnectRedis(configModel.Redis)

	postgresDb2.MigrateDatabase(db)
	mailDialer := email.SetupMailDailer(configModel.Email)

	if mailDialer != nil {
		print(mailDialer.From)
	}
}