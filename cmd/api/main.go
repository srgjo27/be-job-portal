package main

import (
	"log"

	"be-job-portal/internal/config"
	"be-job-portal/internal/delivery/http"
	"be-job-portal/internal/domain"
	"be-job-portal/internal/repository"
	"be-job-portal/internal/usecase"
	"be-job-portal/pkg/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	// Connect Database
	db := database.ConnectDB(cfg)

	// Auto Migrate
	db.AutoMigrate(&domain.User{}, &domain.Job{}, &domain.Application{}, &domain.SeekerProfile{}, &domain.CompanyProfile{}, &domain.Experience{}, &domain.Education{})

	// Init Router
	r := gin.Default()

	// CORS Configuration
	r.Use(cors.Default())

	// Repositories
	userRepo := repository.NewUserRepository(db)
	jobRepo := repository.NewJobRepository(db)
	appRepo := repository.NewApplicationRepository(db)
	profileRepo := repository.NewProfileRepository(db)

	// Usecases
	authUsecase := usecase.NewAuthUsecase(userRepo, cfg)
	jobUsecase := usecase.NewJobUsecase(jobRepo)
	appUsecase := usecase.NewApplicationUsecase(appRepo, jobRepo)
	profileUsecase := usecase.NewProfileUsecase(profileRepo)

	// Handlers
	// Handlers
	authHandler := http.NewAuthHandler(authUsecase)
	jobHandler := http.NewJobHandler(jobUsecase)
	appHandler := http.NewApplicationHandler(appUsecase)
	profileHandler := http.NewProfileHandler(profileUsecase, userRepo)
	dashboardHandler := http.NewDashboardHandler(appUsecase)

	// Register Routes
	http.RegisterRoutes(r, authHandler, jobHandler, appHandler, profileHandler, dashboardHandler)

	// Run Server
	r.Run(":" + cfg.ServerPort)
}
