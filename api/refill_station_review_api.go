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
	review.Timestamp = time.Now()
	result := db.Create(&review)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusCreated, review)
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
