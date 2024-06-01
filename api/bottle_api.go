package api

import (
	"net/http"
	"strconv"

	"github.com/PoseidonPSE2/code_backend/database"
	"github.com/gin-gonic/gin"
)

// BottleResponse represents a bottle in the response
type BottleResponse struct {
	ID         uint    `json:"id"`
	UserID     uint    `json:"user_id"`
	NFCID      string  `json:"nfc_id"`
	FillVolume int     `json:"fill_volume"`
	WaterType  string  `json:"water_type"`
	Title      string  `json:"title"`
	PathImage  *string `json:"path_image"`
	Active     bool    `json:"active"`
}

type BottlePreferences struct {
	ID         uint   `json:"id"`
	UserID     uint   `json:"user_id"`
	FillVolume int    `json:"fill_volume"`
	WaterType  string `json:"water_type"`
}

// CreateBottleRequest represents a request to create a bottle
type CreateBottleRequest struct {
	UserID     uint    `json:"user_id"`
	NFCID      string  `json:"nfc_id"`
	FillVolume int     `json:"fill_volume"`
	WaterType  string  `json:"water_type"`
	Title      string  `json:"title"`
	PathImage  *string `json:"path_image"`
	Active     bool    `json:"active"`
}

// UpdateBottleRequest represents a request to update a bottle
type UpdateBottleRequest struct {
	ID         uint    `json:"id"`
	UserID     uint    `json:"user_id"`
	NFCID      string  `json:"nfc_id"`
	FillVolume int     `json:"fill_volume"`
	WaterType  string  `json:"water_type"`
	Title      string  `json:"title"`
	PathImage  *string `json:"path_image"`
	Active     bool    `json:"active"`
}

// @Summary Show all bottles
// @Description Get all bottles
// @Tags bottles
// @Accept  json
// @Produce  json
// @Success 200 {array} BottleResponse
// @Router /bottles [get]
func GetBottles(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		var bottles []database.Bottle
		result := db.Find(&bottles)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		respondWithJSON(c, http.StatusOK, bottles)
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var bottle database.Bottle
		result := db.First(&bottle, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
			return
		}
		respondWithJSON(c, http.StatusOK, bottle)
	}
}

// @Summary Get bottle preferences by the NFC ID
// @Description Get bottle preferences by the NFC ID
// @Tags bottles
// @Accept json
// @Produce json
// @Param nfc_id path string true "NFC ID"
// @Success 200 {object} BottlePreferences
// @Router /bottles/preferences/{nfc-id} [get]
func GetBottlePreferencesByNFCId(c *gin.Context) {
	idStr := c.Param("id")
	var result BottlePreferences

	if err := db.Model(&database.Bottle{}).Where("nfc_id = ?", idStr).First(&result).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Row not found for NFC ID"})
		return
	}

	respondWithJSON(c, http.StatusOK, result)
}

// @Summary Create a bottle
// @Description Create a new bottle
// @Tags bottles
// @Accept  json
// @Produce  json
// @Param bottle body CreateBottleRequest true "Bottle"
// @Success 201 {object} BottleResponse
// @Router /bottles [post]
func CreateBottle(c *gin.Context) {
	var bottle database.Bottle
	if err := c.ShouldBindJSON(&bottle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := db.Create(&bottle)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	respondWithJSON(c, http.StatusCreated, bottle)
}

// @Summary Update a bottle
// @Description Update an existing bottle
// @Tags bottles
// @Accept  json
// @Produce  json
// @Param bottle body UpdateBottleRequest true "Bottle"
// @Success 200 {object} BottleResponse
// @Router /bottles [put]
func UpdateBottle(c *gin.Context) {
	var bottle database.Bottle
	if err := c.ShouldBindJSON(&bottle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := db.Save(&bottle)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	respondWithJSON(c, http.StatusOK, bottle)
}

// @Summary Delete a bottle
// @Description Delete an existing bottle
// @Tags bottles
// @Accept  json
// @Produce  json
// @Param id query int true "Bottle ID"
// @Success 204
// @Router /bottles [delete]
func DeleteBottle(c *gin.Context) {
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
	result := db.Delete(&database.Bottle{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Get all bottles by user ID
// @Description Get all bottles associated with a specific user
// @Tags bottles
// @Accept  json
// @Produce  json
// @Param userId path int true "User ID"
// @Success 200 {array} database.Bottle
// @Router /users/{userId}/bottles [get]
func GetBottlesByUserID(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	var bottles []database.Bottle
	result := db.Where("user_id = ?", userID).Find(&bottles)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	respondWithJSON(c, http.StatusOK, bottles)
}
