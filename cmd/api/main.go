package main

import (
	"flight_booking_system/internals/config"
	"flight_booking_system/internals/database"
	"flight_booking_system/internals/handlers"
	"flight_booking_system/internals/handlers/admin_handlers"
	"flight_booking_system/internals/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to connect to configuration")
	}

	pool, err := database.Connect(cfg.DBURL)
	if err != nil {
		log.Fatal("Failed to connect ot database")
	}

	defer pool.Close()

	router := gin.Default()

	router.POST("/create", handlers.CreateUserHandler(pool))
	router.POST("/login", handlers.LoginHandler(pool, cfg))

	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware(cfg))
	admin.Use(middleware.AuthAdmin(pool))
	{
		admin.POST("/airports", admin_handlers.AddAirportHandler(pool))
		admin.POST("/airplanes", admin_handlers.AddAirplaneHandler(pool))
		admin.POST("/flights", admin_handlers.CreateFlightHandler(pool))
		admin.POST("/tickets", admin_handlers.GenerateTicketHandler(pool))
	}

	protected := router.Group("/users")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.POST("/search", handlers.SearchEngineHandler(pool))
		protected.POST("/booking", handlers.BookTicketHandler(pool))
	}
	router.Run(":" + cfg.Port)
}
