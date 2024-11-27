package controllers

import (
	"backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func BookTicket(c *gin.Context) {
	var booking models.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Check ticket availability
	var ticket models.TicketAvailability
	if err := models.DB.Where("event_id = ?", booking.EventID).First(&ticket).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if ticket.Available < booking.Tickets {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough tickets available"})
		return
	}

	// Deduct tickets and save booking
	ticket.Available -= booking.Tickets
	models.DB.Save(&ticket)

	booking.Status = "pending"
	models.DB.Create(&booking)

	c.JSON(http.StatusOK, gin.H{"message": "Booking created", "booking": booking})
}

func CancelBooking(c *gin.Context) {
	bookingID, _ := strconv.Atoi(c.Param("id"))

	var booking models.Booking
	if err := models.DB.First(&booking, bookingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	if booking.Status != "booked" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot cancel this booking"})
		return
	}

	// Restore ticket availability
	var ticket models.TicketAvailability
	models.DB.Where("event_id = ?", booking.EventID).First(&ticket)
	ticket.Available += booking.Tickets
	models.DB.Save(&ticket)

	// Update booking status
	booking.Status = "canceled"
	models.DB.Save(&booking)

	c.JSON(http.StatusOK, gin.H{"message": "Booking canceled"})
}

func CheckTicketAvailability(c *gin.Context) {
	eventID, _ := strconv.Atoi(c.Query("event_id"))

	var ticket models.TicketAvailability
	if err := models.DB.Where("event_id = ?", eventID).First(&ticket).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"available_tickets": ticket.Available})
}
