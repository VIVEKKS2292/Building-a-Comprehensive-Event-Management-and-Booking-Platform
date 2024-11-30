package main

import (
	"github.com/gin-contrib/cors" // Import the CORS middleware
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

	// Configure CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Allow frontend to connect
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

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

	// Event and user routes with authorization checks
	r.GET("/events", middleware.AuthMiddleware("Admin", "User", "Organizer"), controllers.GetEvents)
	r.GET("/events/filter", middleware.AuthMiddleware("Admin", "User", "Organizer"), controllers.FilterEvents)

	// Ticket Booking routes
	bookingRoutes := r.Group("/bookings")
	{
		bookingRoutes.Use(middleware.AuthMiddleware("User"))
		bookingRoutes.POST("/", controllers.BookTicket)
		bookingRoutes.DELETE("/:id", controllers.CancelBooking)
		bookingRoutes.GET("/availability", controllers.CheckTicketAvailability)
		bookingRoutes.POST("/payments", controllers.HandlePayment)
	}

	// Ticket Availability routes (only Admin and Organizer)
	ticketAvailabilityRoutes := r.Group("/ticket-availability")
	{
		ticketAvailabilityRoutes.Use(middleware.AuthMiddleware("Admin", "Organizer"))
		ticketAvailabilityRoutes.POST("/", controllers.CreateTicketAvailability)         // Create new ticket availability
		ticketAvailabilityRoutes.GET("/:event_id", controllers.GetTicketAvailability)    // Get ticket availability for specific event
		ticketAvailabilityRoutes.PUT("/:event_id", controllers.UpdateTicketAvailability) // Update ticket availability for specific event
	}

	r.GET("/realtime/tickets", controllers.TicketUpdates)

	// Wishlist routes
	wishlistRoutes := r.Group("/wishlist")
	{
		wishlistRoutes.Use(middleware.AuthMiddleware("User"))
		wishlistRoutes.POST("/:event_id", controllers.AddToWishlist)        // Add event to wishlist
		wishlistRoutes.GET("/", controllers.GetWishlist)                    // Get user's wishlist
		wishlistRoutes.DELETE("/:event_id", controllers.RemoveFromWishlist) // Remove event from wishlist
	}

	// Event Analytics routes
	analyticsRoutes := r.Group("/analytics")
	{
		analyticsRoutes.Use(middleware.AuthMiddleware("Organizer"))
		analyticsRoutes.GET("/ticket-sales", controllers.GetTicketSalesAnalytics)          // Get ticket sales analytics
		analyticsRoutes.GET("/revenue", controllers.GetRevenueAnalytics)                   // Get revenue analytics
		analyticsRoutes.GET("/attendee-demographics", controllers.GetAttendeeDemographics) // Get attendee demographics
	}

	// Organization routes
	organizationRoutes := r.Group("/organizations")
	{
		organizationRoutes.Use(middleware.AuthMiddleware("Organizer"))
		organizationRoutes.POST("/", controllers.CreateOrganization)
		organizationRoutes.PUT("/:id", controllers.UpdateOrganization)
	}

	// Start the server
	r.Run(":8080")
}
