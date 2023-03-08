package main

import (
	"time"

	"github.com/clementb49/welsh_academy/config"
	"github.com/clementb49/welsh_academy/middlewares"
	"github.com/clementb49/welsh_academy/models"
	"github.com/clementb49/welsh_academy/routes"

	"github.com/fvbock/endless"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"moul.io/zapgorm2"
)

// Build the Gin HTTP server with routes and middleware, and initialize it with an endless server.
func main() {
	// Load the application configuration from the configuration source.
	waCfg := config.GetWaConfig()
	// Build the logger with the specified configuration.
	logger := buildLogger(waCfg.ConfigureLogger())
	defer logger.Sync()

	// Configure the Gin framework with the application configuration.
	waCfg.ConfigureGin()
	logger.Info("Configuration loaded")
	// Build the GORM database object with the specified configuration.
	dsn := waCfg.DbCfg.BuildDsn()
	db := initGorm(logger, dsn)

	// Migrate the database schema.
	migrateDb(db, logger)
	logger.Info("Initializing HTTP server...")
	// Create a new Gin HTTP server.
	eng := gin.New()
	// Register the Gin middleware for logging HTTP requests and responses.
	registerLogMiddleware(eng, logger)
	logger.Info("server initialized")
	// Register the API routes.
	registerApiRoutes(db, eng, logger)
	// Initialize the endless HTTP server with the Gin HTTP server.
	srv := endless.NewServer(":8000", eng)
	srv.ErrorLog = zap.NewStdLog(logger)
	// Start serving HTTP requests.
	srv.ListenAndServe()
}

// Register the Gin middleware for logging HTTP requests and responses with zap logger.
func registerLogMiddleware(eng *gin.Engine, logger *zap.Logger) {
	// Log HTTP requests and responses with the Zap logger.
	eng.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	// Recover from panics in the HTTP request handler and log the error with the Zap logger.
	eng.Use(ginzap.RecoveryWithZap(logger, true))
}

// Build a Zap logger with the specified configuration.
func buildLogger(cfg zap.Config) *zap.Logger {
	// Build the logger with the specified configuration.
	logger := zap.Must(cfg.Build())
	// Replace the zap global logger with the built logger.
	zap.ReplaceGlobals(logger)
	// Redirect standard Go log to the logger.
	zap.RedirectStdLog(logger)
	return logger
}

// Initialize a GORM database object with the specified configuration.
func initGorm(logger *zap.Logger, dsn string) *gorm.DB {
	logger.Sugar().Debugf("Intializing DB with the following database connection string: %s", dsn)
	// Create a Zap logger for GORM to log SQL statements and errors.
	gLogger := zapgorm2.New(logger)
	gLogger.SetAsDefault()
	// Open the database connection with GORM.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "wac_",
		},
	})
	if err != nil {
		logger.Sugar().Fatalf("Unable to connect to the database: %w", err)
	}
	return db
}

// Migrate the database schema.
func migrateDb(db *gorm.DB, logger *zap.Logger) {
	logger.Info("Begin database migration ...")
	// Auto-migrate the database schema for the specified models.
	err := db.AutoMigrate(&models.User{}, &models.Ingredient{}, &models.Recipe{})
	if err != nil {
		logger.Sugar().Fatalf("The database migration encounter the folowing error: %w", err)
	}
	logger.Info("Database migration terminated successfully")
}

// register the API route in the gin framework
func registerApiRoutes(db *gorm.DB, eng *gin.Engine, logger *zap.Logger) {
	logger.Info("Registering API routes ...")
	// Create two Gin router groups for authenticated and unauthenticated routes with the same base path
	unauthApiRouter := eng.Group("/api/v1")
	authApiRouter := eng.Group("/api/v1")
	// Apply an authentication middleware to the authenticated API router
	authApiRouter.Use(middlewares.Auth())
	// Register the API routes for ingredients, recipes, and users
	routes.InitIngredientRoute(db, unauthApiRouter, authApiRouter)
	routes.InitRecipeRoute(db, unauthApiRouter, authApiRouter)
	routes.InitUserRoutes(db, unauthApiRouter, authApiRouter)
	logger.Info("API routes registered")
}
