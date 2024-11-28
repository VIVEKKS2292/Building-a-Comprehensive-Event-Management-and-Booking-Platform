package controllers

import (
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get ticket sales analytics
func GetTicketSalesAnalytics(c *gin.Context) {
	eventID := c.Query("event_id")

	var result []struct {
		TimePeriod string `json:"time_period"`
		Sales      int    `json:"sales"`
	}
	query := `
		SELECT DATE(created_at) AS time_period, COUNT(*) AS sales
		FROM bookings
		WHERE event_id = ?
		GROUP BY DATE(created_at)
	`
	if err := models.DB.Raw(query, eventID).Scan(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ticket sales analytics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ticket_sales": result})
}

// Get revenue analytics
func GetRevenueAnalytics(c *gin.Context) {
	eventID := c.Query("event_id")

	var result struct {
		TotalRevenue float64 `json:"total_revenue"`
	}
	query := `
		SELECT SUM(amount) AS total_revenue
		FROM transactions
		WHERE event_id = ? AND status = "success"
	`
	if err := models.DB.Raw(query, eventID).Scan(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch revenue analytics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"revenue": result.TotalRevenue})
}

// Get attendee demographics
func GetAttendeeDemographics(c *gin.Context) {
	eventID := c.Query("event_id")

	var result []struct {
		AgeGroup string `json:"age_group"`
		Count    int    `json:"count"`
	}
	query := `
		SELECT 
			CASE 
				WHEN age BETWEEN 18 AND 24 THEN '18-24'
				WHEN age BETWEEN 25 AND 34 THEN '25-34'
				WHEN age BETWEEN 35 AND 44 THEN '35-44'
				ELSE '45+'
			END AS age_group,
			COUNT(*) AS count
		FROM users
		JOIN bookings ON users.id = bookings.user_id
		WHERE bookings.event_id = ?
		GROUP BY age_group
	`
	if err := models.DB.Raw(query, eventID).Scan(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attendee demographics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"attendee_demographics": result})
}
