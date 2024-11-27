package controllers

import (
	"backend/models"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Dummy Payment Processor
func processPayment(_ float64) string {
	// Simulate payment statuses randomly
	statuses := []string{"success", "failure", "pending"}
	rand.Seed(time.Now().UnixNano())
	return statuses[rand.Intn(len(statuses))]
}

// Handle Payment Status
func HandlePayment(c *gin.Context) {
	var transaction models.Transaction

	// Bind JSON data to transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment data"})
		return
	}

	// Simulate payment processing
	transaction.Status = processPayment(transaction.Amount)
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = time.Now()

	// Save the transaction
	if err := models.DB.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record transaction"})
		return
	}

	// Update booking based on payment status
	var booking models.Booking
	if err := models.DB.First(&booking, "user_id = ? AND event_id = ?", transaction.UserID, transaction.EventID).Error; err == nil {
		if transaction.Status == "success" {
			booking.Status = "booked"
		} else {
			booking.Status = "pending"
		}
		models.DB.Save(&booking)
	}

	// Respond with the payment status
	c.JSON(http.StatusOK, gin.H{
		"transaction": transaction,
		"message":     "Payment processed successfully",
	})
}
