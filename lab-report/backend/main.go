package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitbub.com/zikrullahcelep611/lab-report/backend/api/auth"
	"gitbub.com/zikrullahcelep611/lab-report/backend/api/hospital"
	"gitbub.com/zikrullahcelep611/lab-report/backend/api/patient"
	"gitbub.com/zikrullahcelep611/lab-report/backend/api/report"
	"gitbub.com/zikrullahcelep611/lab-report/backend/api/user"
	authservice "gitbub.com/zikrullahcelep611/lab-report/backend/application/authService"
	hospitalservice "gitbub.com/zikrullahcelep611/lab-report/backend/application/hospitalService"
	jwtservice "gitbub.com/zikrullahcelep611/lab-report/backend/application/jwtService"
	patientservice "gitbub.com/zikrullahcelep611/lab-report/backend/application/patientService"
	reportservice "gitbub.com/zikrullahcelep611/lab-report/backend/application/reportService"
	tokenservice "gitbub.com/zikrullahcelep611/lab-report/backend/application/tokenService"
	userservice "gitbub.com/zikrullahcelep611/lab-report/backend/application/userService"
	backgroundjobs "gitbub.com/zikrullahcelep611/lab-report/backend/background_jobs"
	"gitbub.com/zikrullahcelep611/lab-report/backend/infrastructure/config"
	postgresDb2 "gitbub.com/zikrullahcelep611/lab-report/backend/infrastructure/postgresDb"
	hospitalrepository "gitbub.com/zikrullahcelep611/lab-report/backend/infrastructure/repository/hospitalRepository"
	loginrepository "gitbub.com/zikrullahcelep611/lab-report/backend/infrastructure/repository/loginRepository"
	patientrepository "gitbub.com/zikrullahcelep611/lab-report/backend/infrastructure/repository/patientRepository"
	reportrepository "gitbub.com/zikrullahcelep611/lab-report/backend/infrastructure/repository/reportRepository"
	tokenrepository "gitbub.com/zikrullahcelep611/lab-report/backend/infrastructure/repository/tokenRepository"
	userrepository "gitbub.com/zikrullahcelep611/lab-report/backend/infrastructure/repository/userRepository"
	authmiddleware "gitbub.com/zikrullahcelep611/lab-report/backend/middleware/authMiddleware"
	contexttimeoutmiddleware "gitbub.com/zikrullahcelep611/lab-report/backend/middleware/contextTimeoutMiddleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func main() {
	configModel := config.SetConfig("resources")
	zerolog.SetGlobalLevel(configModel.Log.Level)

	db := postgresDb2.ConnectDatabase(configModel.Database)
	//redis := redis2.ConnectRedis(configModel.Redis)
	//mailDialer := email.SetupMailDailer(configModel.Email)

	postgresDb2.MigrateDatabase(db)

	newLoginRepository := loginrepository.NewRepository(db)
	newReportRepository := reportrepository.NewRepository(db)
	newTokenRepository := tokenrepository.NewRepository(db)
	newUserRepository := userrepository.NewRepository(db)
	newHospitalRepository := hospitalrepository.NewRepository(db)
	newPatientRepository := patientrepository.NewRepository(db)

	newJwtService := jwtservice.NewJwtService(configModel.JWT.SecretKey)
	newTokenService := tokenservice.NewTokenService(newTokenRepository)
	newAuthService := authservice.NewAuthService(newLoginRepository, newTokenService, newJwtService)
	newReportService := reportservice.NewReportService(newReportRepository)
	newUserService := userservice.NewUserService(newUserRepository)
	newHospitalService := hospitalservice.NewHospitalService(newHospitalRepository)
	newPatientService := patientservice.NewPatientService(newPatientRepository)

	newAuthHandler := auth.NewAuthController(newAuthService, newJwtService)
	newReportHandler := report.NewReportController(newReportService)
	newUserHandler := user.NewUserController(newUserService, newJwtService)
	newHospitalHandler := hospital.NewHospitalHandler(newHospitalService)
	newPatientHandler := patient.NewPatientHandler(newPatientService)

	app := fiber.New(fiber.Config{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(contexttimeoutmiddleware.TimeoutMiddleware(5))

	newAuthMiddleware := authmiddleware.NewAuthMiddleware(newTokenService, newJwtService)

	auth.RegisterAuthRoutes(app, newAuthHandler)
	user.RegisterUserRoutes(app, newUserHandler)

	protected := app.Group("/api")
	protected.Use(newAuthMiddleware.Authenticate)

	report.RegisterReportRoutes(protected, newReportHandler)
	hospital.RegisterHospitalRoutes(app, newHospitalHandler)
	patient.RegisterPatientRouter(protected, newPatientHandler)

	backgroundjobs.StartCleanExpiredJwtTokens(newTokenService)

	go func() {
		addr := fmt.Sprintf(":%d", configModel.Server.Port)
		log.Info().Msgf("Server started on port %d", configModel.Server.Port)
		if err := app.Listen(addr); err != nil {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Error().Err(err).Msg("Server shutdown error")
	}

	gracefulShutdown(db)
	log.Info().Msg("Server shutdown complete")
}

func gracefulShutdown(db *gorm.DB) {
	log.Info().Msg("Shutting down server...")

	// Database connection'ı kapat
	if db != nil {
		log.Info().Msg("Closing database connection...")
		sqlDB, err := db.DB()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get sql.DB from gorm.DB")
		} else {
			if err := sqlDB.Close(); err != nil {
				log.Error().Err(err).Msg("Failed to close database connection")
			} else {
				log.Info().Msg("Database connection closed.")
			}
		}
	}

	log.Info().Msg("Server shutdown complete.")
}
