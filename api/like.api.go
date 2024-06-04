package api

import (
	"net/http"
	"strconv"

	"github.com/PoseidonPSE2/code_backend/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Show all likes
// @Description Get all likes
// @Tags likes
// @Accept  json
// @Produce  json
// @Success 200 {array} database.Like
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
		c.JSON(http.StatusOK, likes)
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
		c.JSON(http.StatusOK, like)
	}
}

// @Summary Check if a user likes a refill station
// @Description Check if a specific user likes a specific refill station
// @Tags likes
// @Accept  json
// @Produce  json
// @Param refillstationId query int true "Refill Station ID"
// @Param userId query int true "User ID"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /refillstation_like [get]
func GetLikeByUserIdAndStationID(c *gin.Context) {
	refillstationIdStr := c.Query("refillstationId")
	userIdStr := c.Query("userId")

	if refillstationIdStr == "" || userIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refillstationId and userId are required"})
		return
	}

	refillstationId, err := strconv.Atoi(refillstationIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid refillstationId"})
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId"})
		return
	}

	var like database.Like
	result := db.Where("station_id = ? AND user_id = ?", refillstationId, userId).First(&like)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	isLiked := result.RowsAffected > 0
	c.JSON(http.StatusOK, gin.H{"isLiked": isLiked})
}

// @Summary Create a like
// @Description Create a new like
// @Tags likes
// @Accept  json
// @Produce  json
// @Param like body database.Like true "Like"
// @Success 201 {object} database.Like
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
	c.JSON(http.StatusCreated, like)
}

// @Summary Update a like
// @Description Update an existing like
// @Tags likes
// @Accept  json
// @Produce  json
// @Param like body database.Like true "Like"
// @Success 200 {object} database.Like
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
	c.JSON(http.StatusOK, like)
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
