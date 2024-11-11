// services/rating.go
package services

import (
	"yak.rabai/config"
	"yak.rabai/models"
)

// SaveRating saves a rating and comment for a listener
func SaveRating(userID string, rating int, comment string) error {
	ratingRecord := models.Rating{
		UserID:  userID,
		Rating:  rating,
		Comment: comment,
	}

	// Store the rating in the database
	if err := config.DB.Create(&ratingRecord).Error; err != nil {
		return err
	}

	return nil
}

// GetUserRatings retrieves all ratings for a specific user
func GetUserRatings(userID uint) ([]models.Rating, error) {
	var ratings []models.Rating
	if err := config.DB.Where("user_id = ?", userID).Find(&ratings).Error; err != nil {
		return nil, err
	}
	return ratings, nil
}
