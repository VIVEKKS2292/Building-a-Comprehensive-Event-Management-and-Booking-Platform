package controllers

import (
	"backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Add an event to the wishlist
func AddToWishlist(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		return
	}

	// Parse eventID from the URL parameter
	eventID, err := strconv.Atoi(c.Param("event_id"))
	if err != nil || eventID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	// Check if the event exists
	var event models.Event
	if err := models.DB.First(&event, eventID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	// Add the event to the user's wishlist
	wishlist := models.Wishlist{
		UserID:  userID.(uint), // Ensure proper type casting
		EventID: uint(eventID),
	}
	if err := models.DB.Create(&wishlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to wishlist"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Event added to wishlist"})
}

func GetWishlist(c *gin.Context) {
	// Retrieve userID from the context (from AuthMiddleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		return
	}

	// Fetch the list of wishlist entries for the given userID
	var wishlists []models.Wishlist
	if err := models.DB.Where("user_id = ?", userID).Find(&wishlists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No wishlist found for user"})
		return
	}

	// Prepare the response structure
	var response []map[string]interface{}

	for _, wishlist := range wishlists {
		// Fetch the event details for each event ID in the wishlist
		var event models.Event
		if err := models.DB.Where("id = ?", wishlist.EventID).First(&event).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch event details"})
			return
		}

		// Append the wishlist entry with the associated event details
		response = append(response, map[string]interface{}{
			"wishlist_id": wishlist.ID,
			"user_id":     wishlist.UserID,
			"event":       event,
		})
	}

	// Return the constructed response
	c.JSON(http.StatusOK, gin.H{"wishlist": response})
}

// Remove an event from the wishlist
func RemoveFromWishlist(c *gin.Context) {
	// Retrieve userID from the context (from AuthMiddleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		return
	}

	// Parse eventID from the URL parameter
	eventID, err := strconv.Atoi(c.Param("event_id"))
	if err != nil || eventID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	// Delete the event from the user's wishlist
	if err := models.DB.Where("user_id = ? AND event_id = ?", userID, eventID).Delete(&models.Wishlist{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove event from wishlist"})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{"message": "Event removed from wishlist"})
}
