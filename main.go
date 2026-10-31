package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RaihanStark/clockworth/internal/auth"
	"github.com/RaihanStark/clockworth/internal/config"
	"github.com/RaihanStark/clockworth/internal/database"
	"github.com/RaihanStark/clockworth/internal/handlers"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	// Configure connection pool optimized for production with network latency
	db.SetMaxOpenConns(50)                          // Increased for production load
	db.SetMaxIdleConns(10)                          // Keep more connections ready
	db.SetConnMaxLifetime(10 * time.Minute)         // Longer lifetime for remote DB
	db.SetConnMaxIdleTime(5 * time.Minute)          // Close idle connections after 5 min

	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping database:", err)
	}

	log.Println("Database connection pool configured: max_open=50, max_idle=10")

	// Initialize database queries
	queries := database.New(db)

	// Initialize OAuth
	googleOAuth := auth.NewGoogleOAuth(
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.GoogleRedirectURL,
	)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(queries, googleOAuth, cfg)
	settingsHandler := handlers.NewSettingsHandler(queries)
	activitiesHandler := handlers.NewActivitiesHandler(queries)

	// Setup Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{cfg.FrontendURL},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	// Public routes
	e.GET("/health", healthCheck)
	e.POST("/auth/register", authHandler.Register)
	e.POST("/auth/login", authHandler.EmailLogin)
	e.GET("/auth/google/login", authHandler.Login)
	e.GET("/auth/google/callback", authHandler.Callback)

	// Protected routes
	protected := e.Group("/api")
	protected.Use(auth.AuthMiddleware(queries))
	protected.GET("/me", authHandler.Me)
	protected.POST("/logout", authHandler.Logout)

	// Settings routes
	protected.GET("/settings", settingsHandler.GetSettings)
	protected.PUT("/settings/hourly-worth", settingsHandler.UpdateHourlyWorth)

	// Activities routes
	protected.GET("/activities", activitiesHandler.ListActivities)
	protected.GET("/activities/heatmap", activitiesHandler.GetLast30DaysHeatmap)
	protected.GET("/activities/:id", activitiesHandler.GetActivity)
	protected.POST("/activities", activitiesHandler.CreateActivity)
	protected.PUT("/activities/:id", activitiesHandler.UpdateActivity)
	protected.DELETE("/activities/:id", activitiesHandler.DeleteActivity)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.Port)))
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
