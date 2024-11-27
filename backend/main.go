package main

import (
	"github.com/gin-gonic/gin"

	"backend/controllers"
	"backend/middleware"
	"backend/models"
)

func main() {
	// Initialize the database
	models.InitDB()

	// Create Gin router
	r := gin.Default()

	// Public routes
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// Protected routes
	eventRoutes := r.Group("/events")
	{
		eventRoutes.Use(middleware.AuthMiddleware("Admin", "Organizer"))
		eventRoutes.POST("/", controllers.CreateEvent)
		eventRoutes.PUT("/:id", controllers.UpdateEvent)
		eventRoutes.DELETE("/:id", controllers.DeleteEvent)
	}

	r.GET("/events", middleware.AuthMiddleware("Admin", "User", "Organizer"), controllers.GetEvents)
	// Public route for filtering events
	r.GET("/events/filter", middleware.AuthMiddleware("Admin", "User", "Organizer"), controllers.FilterEvents)

	// Start the server
	r.Run(":8080")
}
