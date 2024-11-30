package controllers

import (
	"net/http"
	"strconv"

	"backend/models"

	"github.com/gin-gonic/gin"
)

// CreateOrganization - Create a new organization
func CreateOrganization(c *gin.Context) {
	// Extract user ID from middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		return
	}

	// Check if the user already has an organization
	var existingOrg models.Organization
	if err := models.DB.Where("user_id = ?", userID).First(&existingOrg).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already has an organization"})
		return
	}

	// Bind JSON payload to the Organization model
	var organization models.Organization
	if err := c.ShouldBindJSON(&organization); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Associate the organization with the user
	organization.UserID = userID.(uint)

	// Save the organization in the database
	if err := models.DB.Create(&organization).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organization created successfully", "organization": organization})
}

// UpdateOrganization - Update an existing organization
func UpdateOrganization(c *gin.Context) {
	// Extract user ID from middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		return
	}

	// Parse and validate organization ID from the request
	organizationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	// Fetch the organization to validate ownership
	var organization models.Organization
	if err := models.DB.First(&organization, organizationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	// Check if the user owns the organization
	if organization.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this organization"})
		return
	}

	// Bind the request JSON to the Organization model
	if err := c.ShouldBindJSON(&organization); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the organization in the database
	if err := models.DB.Model(&organization).Updates(organization).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update organization"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organization updated successfully", "organization": organization})
}
