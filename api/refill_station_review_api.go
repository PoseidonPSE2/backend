package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/PoseidonPSE2/code_backend/database"
	"github.com/gin-gonic/gin"
)

// @Summary Show all refill station reviews
// @Description Get all refill station reviews
// @Tags refill_station_reviews
// @Accept  json
// @Produce  json
// @Success 200 {array} database.RefillStationReview
// @Router /refill_station_reviews [get]
func GetRefillStationReviews(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		var reviews []database.RefillStationReview
		result := db.Find(&reviews)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, reviews)
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var review database.RefillStationReview
		result := db.First(&review, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, review)
	}
}

// @Summary Show all refill station reviews by user ID and station ID
// @Description Get all refill station reviews by user ID and station ID
// @Tags refill_station_reviews
// @Accept  json
// @Produce  json
// @Param userId path int true "User ID"
// @Param stationId path int true "Station ID"
// @Success 200 {array} database.RefillStationReview
// @Router /refill_station_reviews/{userId}/{stationId} [get]
func GetRefillStationReviewsByUserId(c *gin.Context) {
	userIdStr := c.Param("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	stationIdStr := c.Param("stationId")
	stationId, err := strconv.Atoi(stationIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Station ID"})
		return
	}

	var reviews []database.RefillStationReview
	result := db.Where("user_id = ? AND station_id = ?", userId, stationId).Find(&reviews)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if len(reviews) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No reviews found for this user and station"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// @Summary Create a refill station review
// @Description Create a new refill station review
// @Tags refill_station_reviews
// @Accept  json
// @Produce  json
// @Param review body database.RefillStationReview true "Refill Station Review"
// @Success 201 {object} database.RefillStationReview
// @Router /refill_station_reviews [post]
func CreateRefillStationReview(c *gin.Context) {
	var review database.RefillStationReview
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user has already reviewed this station
	var existingReview database.RefillStationReview
	result := db.Where("user_id = ? AND station_id = ?", review.UserID, review.StationID).First(&existingReview)
	if result.Error == nil {
		// User has already reviewed this station, update the existing review
		existingReview.Cleanness = review.Cleanness
		existingReview.Accessibility = review.Accessibility
		existingReview.WaterQuality = review.WaterQuality
		existingReview.Timestamp = time.Now()

		if err := db.Save(&existingReview).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, existingReview)
	} else {
		// No existing review found, create a new one
		review.Timestamp = time.Now()
		if err := db.Create(&review).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, review)
	}
}

// @Summary Update a refill station review
// @Description Update an existing refill station review
// @Tags refill_station_reviews
// @Accept  json
// @Produce  json
// @Param review body database.RefillStationReview true "Refill Station Review"
// @Success 200 {object} database.RefillStationReview
// @Router /refill_station_reviews [put]
func UpdateRefillStationReview(c *gin.Context) {
	var review database.RefillStationReview
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	review.Timestamp = time.Now()
	result := db.Save(&review)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, review)
}

// @Summary Delete a refill station review
// @Description Delete an existing refill station review
// @Tags refill_station_reviews
// @Accept  json
// @Produce  json
// @Param id query int true "Refill Station Review ID"
// @Success 204
// @Router /refill_station_reviews [delete]
func DeleteRefillStationReview(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	result := db.Delete(&database.RefillStationReview{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
