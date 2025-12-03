package main

import (
	"log"
	"os"

	"LogiTrackPro/backend/internal/config"
	"LogiTrackPro/backend/internal/database"
	"LogiTrackPro/backend/internal/handlers"
	"LogiTrackPro/backend/internal/optimizer"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize optimizer client
	optimizerClient := optimizer.NewClient(cfg.OptimizerURL)

	// Initialize handlers
	h := handlers.New(db, optimizerClient, cfg)

	// Setup router
	router := setupRouter(h, cfg)

	// Start server
	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRouter(h *handlers.Handler, cfg *config.Config) *gin.Engine {
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// CORS middleware
	router.Use(corsMiddleware())

	// Health check
	router.GET("/health", h.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
			auth.POST("/refresh", h.RefreshToken)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(h.AuthMiddleware())
		{
			// User routes
			protected.GET("/me", h.GetCurrentUser)

			// Warehouse routes
			warehouses := protected.Group("/warehouses")
			{
				warehouses.GET("", h.ListWarehouses)
				warehouses.POST("", h.CreateWarehouse)
				warehouses.GET("/:id", h.GetWarehouse)
				warehouses.PUT("/:id", h.UpdateWarehouse)
				warehouses.DELETE("/:id", h.DeleteWarehouse)
			}

			// Customer routes
			customers := protected.Group("/customers")
			{
				customers.GET("", h.ListCustomers)
				customers.POST("", h.CreateCustomer)
				customers.GET("/:id", h.GetCustomer)
				customers.PUT("/:id", h.UpdateCustomer)
				customers.DELETE("/:id", h.DeleteCustomer)
			}

			// Vehicle routes
			vehicles := protected.Group("/vehicles")
			{
				vehicles.GET("", h.ListVehicles)
				vehicles.POST("", h.CreateVehicle)
				vehicles.GET("/:id", h.GetVehicle)
				vehicles.PUT("/:id", h.UpdateVehicle)
				vehicles.DELETE("/:id", h.DeleteVehicle)
			}

			// Plan routes
			plans := protected.Group("/plans")
			{
				plans.GET("", h.ListPlans)
				plans.POST("", h.CreatePlan)
				plans.GET("/:id", h.GetPlan)
				plans.DELETE("/:id", h.DeletePlan)
				plans.POST("/:id/optimize", h.OptimizePlan)
				plans.GET("/:id/routes", h.GetPlanRoutes)
			}

			// Analytics routes
			analytics := protected.Group("/analytics")
			{
				analytics.GET("/dashboard", h.GetDashboard)
				analytics.GET("/summary", h.GetSummary)
			}
		}
	}

	return router
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
