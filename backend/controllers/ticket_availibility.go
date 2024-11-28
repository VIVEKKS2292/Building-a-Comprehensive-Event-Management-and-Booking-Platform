package controllers

import (
	"backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateTicketAvailability - Create a new ticket availability entry
func CreateTicketAvailability(c *gin.Context) {
	// Extract user ID from middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		return
	}

	var ticketAvailability models.TicketAvailability
	if err := c.ShouldBindJSON(&ticketAvailability); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch the event to validate ownership
	var event models.Event
	if err := models.DB.First(&event, ticketAvailability.EventID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	// Check if the user owns the event
	if event.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to create ticket availability for this event"})
		return
	}

	// Save the ticket availability in the database
	if err := models.DB.Create(&ticketAvailability).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create ticket availability"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Ticket availability created", "data": ticketAvailability})
}

// GetTicketAvailability - Get the ticket availability for a specific event
func GetTicketAvailability(c *gin.Context) {
	// Extract user ID from middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		return
	}

	// Parse and validate event ID from the request
	eventID, err := strconv.Atoi(c.Param("event_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	// Fetch the event to validate ownership
	var event models.Event
	if err := models.DB.First(&event, eventID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	// Check if the user owns the event
	if event.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to view ticket availability for this event"})
		return
	}

	// Fetch ticket availability for the given event ID
	var ticketAvailability models.TicketAvailability
	if err := models.DB.Where("event_id = ?", eventID).First(&ticketAvailability).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket availability not found for event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ticketAvailability})
}

// UpdateTicketAvailability - Update ticket availability for a specific event
func UpdateTicketAvailability(c *gin.Context) {
	// Extract user ID from middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		return
	}

	// Parse and validate event ID from the request
	eventID, err := strconv.Atoi(c.Param("event_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	// Fetch the event to validate ownership
	var event models.Event
	if err := models.DB.First(&event, eventID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	// Check if the user owns the event
	if event.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update ticket availability for this event"})
		return
	}

	// Bind the request JSON to the TicketAvailability struct
	var ticketAvailability models.TicketAvailability
	if err := c.ShouldBindJSON(&ticketAvailability); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the ticket availability in the database
	if err := models.DB.Model(&models.TicketAvailability{}).
		Where("event_id = ?", eventID).
		Update("available", ticketAvailability.Available).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update ticket availability"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket availability updated"})
}
