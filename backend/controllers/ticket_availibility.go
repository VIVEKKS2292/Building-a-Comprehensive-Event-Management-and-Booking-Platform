package controllers

import (
	"backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateTicketAvailability - Create a new ticket availability entry
func CreateTicketAvailability(c *gin.Context) {
	var ticketAvailability models.TicketAvailability
	if err := c.ShouldBindJSON(&ticketAvailability); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	eventID, err := strconv.Atoi(c.Param("event_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	var ticketAvailability models.TicketAvailability
	if err := models.DB.Where("event_id = ?", eventID).First(&ticketAvailability).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket availability not found for event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ticketAvailability})
}

// UpdateTicketAvailability - Update ticket availability for a specific event
func UpdateTicketAvailability(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("event_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

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
