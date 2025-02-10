package main

import (
	"log"
	"net/http"

	"github.com/pabloantipan/hobe-api-boilerplate/config"
	"github.com/pabloantipan/hobe-api-boilerplate/internal/cloud"
	"github.com/pabloantipan/hobe-api-boilerplate/internal/middleware"
	"github.com/pabloantipan/hobe-api-boilerplate/internal/repositories/datastore"
	"github.com/pabloantipan/hobe-api-boilerplate/internal/services"

	"github.com/pabloantipan/hobe-api-boilerplate/internal/handlers"

	"github.com/gin-gonic/gin"

	_ "github.com/pabloantipan/hobe-api-boilerplate/docs/swagger"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Player API
// @version         1.0
// @description     API Server for Player Management
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	logger, err := cloud.NewCloudLogger(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Cloud Logger: %v", err)
	}

	// Initialize Datastore client
	datastoreClient := datastore.NewDatastoreClient(cfg)

	// Initialize repositories
	playerRepo := datastore.NewDatastorePlayerRepository(datastoreClient)

	// Initialize services
	playerService := services.NewPlayerService(playerRepo)

	// Initialize middlewares
	requestLoggerMiddleware := middleware.NewRequestLoggerMiddleware(logger)
	responseLoggerMiddleware := middleware.NewResponseLoggerMiddleware(logger)

	// Initialize handlers
	playerHandler := handlers.NewPlayerHandler(playerService)

	// Add(rateLimiter.Handle)
	healthHandler := handlers.NewHealthHandler(cfg)

	// Setup router
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.DefaultModelsExpandDepth(-1)),
	)

	router.GET("/health", func(c *gin.Context) {
		healthHandler.ServeHTTP(c.Writer, c.Request)
	})

	api := router.Group("/api/v1")
	{
		players := api.Group("/players")
		{
			players.Use(requestLoggerMiddleware.HandleFunc())
			players.Use(responseLoggerMiddleware.HandleFunc())
			players.POST("", playerHandler.Create)
			players.GET("", playerHandler.GetAll)
			players.GET("/:id", playerHandler.GetByID)
			players.PUT("/:id", playerHandler.Update)
			players.DELETE("/:id", playerHandler.Delete)
		}
	}

	err = router.Run(":" + cfg.Port)
	if err != nil {
		log.Fatalf("Starting API Gateway on port %s: %v", cfg.Port, err)
	}
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
