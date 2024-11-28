package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/models"
)

func CreateEvent(c *gin.Context) {
	var event models.Event

	// Fetch the UserID from the context (set by the authentication middleware)
	userID, exists := c.Get("userID") // Assuming "userID" is set in the middleware
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: user ID not found"})
		return
	}

	// Use userID (it will be of type uint)
	fmt.Println("UserID:", userID.(uint))

	// Bind the incoming JSON to the Event struct
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the UserID in the Event struct
	event.UserID = userID.(uint)

	// Save the event to the database
	if result := models.DB.Create(&event); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Event created successfully"})
}

func GetEvents(c *gin.Context) {
	var events []models.Event
	if result := models.DB.Find(&events); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

func UpdateEvent(c *gin.Context) {
	var event models.Event
	id := c.Param("id")

	// Extract user ID from middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		return
	}

	// Fetch the event by ID
	if err := models.DB.First(&event, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	// Check if the event belongs to the logged-in user
	if event.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this event"})
		return
	}

	// Bind JSON payload to the event model
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the updated event
	if err := models.DB.Save(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event updated successfully"})
}

func DeleteEvent(c *gin.Context) {
	id := c.Param("id")

	// Extract user ID from middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		return
	}

	// Fetch the event to check ownership
	var event models.Event
	if err := models.DB.First(&event, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	// Check if the event belongs to the logged-in user
	if event.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this event"})
		return
	}

	// Delete the event
	if result := models.DB.Delete(&event); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}

func FilterEvents(c *gin.Context) {
	var events []models.Event

	// Retrieve query parameters for filters
	date := c.Query("date")
	location := c.Query("location")
	category := c.Query("category")

	// Build the query dynamically based on filters
	query := models.DB
	if date != "" {
		query = query.Where("date = ?", date)
	}
	if location != "" {
		query = query.Where("location = ?", location)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// Fetch filtered events
	if result := query.Find(&events); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}

	c.JSON(http.StatusOK, events)
}
