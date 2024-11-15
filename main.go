package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/Peeranut-Kit/health_api_assignment/docs" // Import Swagger docs
	"github.com/Peeranut-Kit/health_api_assignment/internal/patient"
	"github.com/Peeranut-Kit/health_api_assignment/internal/staff"
	"github.com/Peeranut-Kit/health_api_assignment/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// @title Hospital Middleware API
// @version 1.0
// @description This is a hospital middleware API.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.email peeranut.kit.work@gmail.com

// @host localhost:8080
// @BasePath /

// docker compose up -d --scale api-service=3 --build
func main() {
	defer gracefulShutdown()

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// Initialize database
	db, err := initDatabase()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to the database: %v", err))
	}

	fmt.Println("Database connected successfully")

	// Gin Framework
	r := gin.Default()

	// Dependency Injection
	patientRepo := patient.NewGormPatientRepository(db)
	staffRepo := staff.NewGormStaffRepository(db)

	patientService := patient.NewPatientService(patientRepo)
	staffService := staff.NewStaffService(staffRepo)

	patientHandler := patient.NewHttpPatientHandler(patientService)
	staffHandler := staff.NewHttpStaffHandler(staffService)

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// test nginx reverse proxy endpoint
	r.GET("/ping", func(c *gin.Context) {
		hostname, _ := os.Hostname()
		log.Println(hostname)
		c.JSON(http.StatusOK, gin.H{
			"message":  "pong",
			"hostname": hostname,
		})
	})

	// API to create a new hospital staff member
	r.POST("/staff/create", staffHandler.CreateStaff)
	// API for staff login
	r.POST("/staff/login", staffHandler.SignInStaff)

	// API to search for a patient
	r.GET("/patient/search", middleware.AuthRequiredMiddleware, patientHandler.SearchPatient)

	r.Run(":" + os.Getenv("PORT")) // listen and serve on port 8080
}

func initDatabase() (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return nil, err
	}
	return db, nil
}

func gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Println("Shutting down server...")
}
