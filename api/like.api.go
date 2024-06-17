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
// @Tags Likes
// @Accept json
// @Produce json
// @Success 200 {array} database.Like
// @Router /likes [get]
func GetLikes(c *gin.Context) {
	var likes []database.Like
	result := db.Find(&likes)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, likes)
}

// @Summary Check if a user likes a refill station
// @Description Check if a specific user likes a specific refill station
// @Tags Likes
// @Accept json
// @Produce json
// @Param refillstationId path int true "Refill Station ID"
// @Param userId path int true "User ID"
// @Success 200 {object} map[string]bool
// @Router /likes/{refillstationId}/{usedId} [get]
func GetLikeByUserIdAndStationID(c *gin.Context) {
	refillstationIdStr := c.Param("refillstationId")
	userIdStr := c.Param("userId")

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

	// Check for if record exists
	var tempStation database.RefillStation
	if result := db.First(&tempStation, refillstationId); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Refill Station with ID not found"})
		return
	}

	// Check for if record exists
	var tempUser database.User
	if result := db.First(&tempUser, userId); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User with ID not found"})
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
// @Tags Likes
// @Accept json
// @Produce json
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
// @Tags Likes
// @Accept json
// @Produce json
// @Param like body database.Like true "Like"
// @Success 200 {object} database.Like
// @Router /likes [put]
func UpdateLike(c *gin.Context) {
	var requestLike database.Like
	if err := c.ShouldBindJSON(&requestLike); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var helpLike database.Like
	if result := db.First(&helpLike, requestLike.ID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Like with ID not found"})
		return
	}
	db.Model(&helpLike).Updates(requestLike)

	c.JSON(http.StatusOK, requestLike)
}

// @Summary Delete a like
// @Description Delete an existing like
// @Tags Likes
// @Accept json
// @Produce json
// @Param id path int true "Like ID"
// @Success 204
// @Router /likes/{id} [delete]
func DeleteLike(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var like database.Like
	if result := db.First(&like, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Like with ID not found"})
		return
	}
	result := db.Delete(&database.Like{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
