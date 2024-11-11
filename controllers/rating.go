package controllers

import (
	"yak.rabai/config"
	"yak.rabai/models"

	"github.com/gin-gonic/gin"
)

// RateListener saves a rating and feedback from Speaker to Listener
func RateListener(c *gin.Context) {
	var input struct {
		UserID  string `json:"user_id"`
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	rating := models.Rating{
		UserID:  input.UserID,
		Rating:  input.Rating,
		Comment: input.Comment,
	}

	// Save rating to database
	if err := config.DB.Create(&rating).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to save rating"})
		return
	}

	c.JSON(200, gin.H{"message": "Rating saved successfully"})
}
