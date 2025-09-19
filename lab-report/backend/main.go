package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitbub.com/zikrullahcelep611/lab-report/backend/api/auth"
	"gitbub.com/zikrullahcelep611/lab-report/backend/api/report"
	"gitbub.com/zikrullahcelep611/lab-report/backend/api/user"
	authservice "gitbub.com/zikrullahcelep611/lab-report/backend/application/authService"
	jwtservice "gitbub.com/zikrullahcelep611/lab-report/backend/application/jwtService"
	reportservice "gitbub.com/zikrullahcelep611/lab-report/backend/application/reportService"
	tokenservice "gitbub.com/zikrullahcelep611/lab-report/backend/application/tokenService"
	userservice "gitbub.com/zikrullahcelep611/lab-report/backend/application/userService"
	backgroundjobs "gitbub.com/zikrullahcelep611/lab-report/backend/background_jobs"
	"gitbub.com/zikrullahcelep611/lab-report/backend/infrastructure/config"
	postgresDb2 "gitbub.com/zikrullahcelep611/lab-report/backend/infrastructure/postgresDb"
	loginrepository "gitbub.com/zikrullahcelep611/lab-report/backend/infrastructure/repository/loginRepository"
	reportrepository "gitbub.com/zikrullahcelep611/lab-report/backend/infrastructure/repository/reportRepository"
	tokenrepository "gitbub.com/zikrullahcelep611/lab-report/backend/infrastructure/repository/tokenRepository"
	userrepository "gitbub.com/zikrullahcelep611/lab-report/backend/infrastructure/repository/userRepository"
	authmiddleware "gitbub.com/zikrullahcelep611/lab-report/backend/middleware/authMiddleware"
	contexttimeoutmiddleware "gitbub.com/zikrullahcelep611/lab-report/backend/middleware/contextTimeoutMiddleware"
	"github.com/gorilla/mux"
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

	newJwtService := jwtservice.NewJwtService(configModel.JWT.SecretKey)
	newTokenService := tokenservice.NewTokenService(newTokenRepository)
	newAuthService := authservice.NewAuthService(newLoginRepository, newTokenService, newJwtService)
	newReportService := reportservice.NewReportService(newReportRepository)
	newUserService := userservice.NewUserService(newUserRepository)

	newAuthHandler := auth.NewAuthController(newAuthService, newJwtService)
	newReportHandler := report.NewReportController(newReportService)
	newUserHandler := user.NewUserController(newUserService, newJwtService)

	router := mux.NewRouter()

	// Create a subrouter
	securedRouter := router.PathPrefix("/api").Subrouter()

	newAuthMiddleware := authmiddleware.NewAuthMiddleware(newTokenService, newJwtService)

	router.Use(contexttimeoutmiddleware.TimeoutMiddleware(5))
	securedRouter.Use(newAuthMiddleware.Authenticate)

	auth.RegisterAuthRoutes(router, newAuthHandler)

	report.RegisterReportRoutes(securedRouter, newReportHandler)
	user.RegisterUserRoutes(router, newUserHandler)

	backgroundjobs.StartCleanExpiredJwtTokens(newTokenService)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", configModel.Server.Port),
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info().Msg(fmt.Sprintf("Server started on port %d", configModel.Server.Port))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	<-quit
	log.Info().Msg("Closing signal received...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Info().Msgf("The server could not be shut down: %v", err)
	}
	gracefulShutdown(ctx, server, db)
	log.Info().Msg("Successful shutdown of the server.")

}

func gracefulShutdown(ctx context.Context, server *http.Server, db *gorm.DB) {
	log.Info().Msg("Shutting down server...")

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to gracefully shutdown server")
	} else {
		log.Info().Msg("Server stopped gracefully.")
	}

	// Close database connection
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
