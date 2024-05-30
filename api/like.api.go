package api

import (
	"net/http"
	"strconv"

	"github.com/PoseidonPSE2/code_backend/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LikeResponse represents a like in the response
type LikeResponse struct {
	ID        uint `json:"id"`
	StationID uint `json:"station_id"`
	UserID    uint `json:"user_id"`
}

// CreateLikeRequest represents a request to create a like
type CreateLikeRequest struct {
	StationID uint `json:"station_id"`
	UserID    uint `json:"user_id"`
}

// UpdateLikeRequest represents a request to update a like
type UpdateLikeRequest struct {
	ID        uint `json:"id"`
	StationID uint `json:"station_id"`
	UserID    uint `json:"user_id"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// IsLikedResponse represents a response indicating if a user likes a refill station
type IsLikedResponse struct {
	IsLiked bool `json:"isLiked"`
}

// @Summary Show all likes
// @Description Get all likes
// @Tags likes
// @Accept  json
// @Produce  json
// @Success 200 {array} LikeResponse
// @Router /likes [get]
func GetLikes(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		var likes []database.Like
		result := db.Find(&likes)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		respondWithJSON(c, http.StatusOK, likes)
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var like database.Like
		result := db.First(&like, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
			return
		}
		respondWithJSON(c, http.StatusOK, like)
	}
}

// @Summary Check if a user likes a refill station
// @Description Check if a specific user likes a specific refill station
// @Tags likes
// @Accept  json
// @Produce  json
// @Param refillstationId query int true "Refill Station ID"
// @Param userId query int true "User ID"
// @Success 200 {object} IsLikedResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /refillstation_like [get]
func GetLikeByUserIdAndStationID(c *gin.Context) {
	refillstationIdStr := c.Query("refillstationId")
	userIdStr := c.Query("userId")

	if refillstationIdStr == "" || userIdStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{"refillstationId and userId are required"})
		return
	}

	refillstationId, err := strconv.Atoi(refillstationIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{"Invalid refillstationId"})
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{"Invalid userId"})
		return
	}

	var like database.Like
	result := db.Where("station_id = ? AND user_id = ?", refillstationId, userId).First(&like)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, ErrorResponse{result.Error.Error()})
		return
	}

	isLiked := result.RowsAffected > 0
	respondWithJSON(c, http.StatusOK, IsLikedResponse{IsLiked: isLiked})
}

// @Summary Create a like
// @Description Create a new like
// @Tags likes
// @Accept  json
// @Produce  json
// @Param like body CreateLikeRequest true "Like"
// @Success 201 {object} LikeResponse
// @Router /likes [post]
func CreateLike(c *gin.Context) {
	var like database.Like
	if err := c.ShouldBindJSON(&like); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := db.Create(&like)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	respondWithJSON(c, http.StatusCreated, like)
}

// @Summary Update a like
// @Description Update an existing like
// @Tags likes
// @Accept  json
// @Produce  json
// @Param like body UpdateLikeRequest true "Like"
// @Success 200 {object} LikeResponse
// @Router /likes [put]
func UpdateLike(c *gin.Context) {
	var like database.Like
	if err := c.ShouldBindJSON(&like); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := db.Save(&like)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	respondWithJSON(c, http.StatusOK, like)
}

// @Summary Delete a like
// @Description Delete an existing like
// @Tags likes
// @Accept  json
// @Produce  json
// @Param id query int true "Like ID"
// @Success 204
// @Router /likes [delete]
func DeleteLike(c *gin.Context) {
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
	result := db.Delete(&database.Like{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
