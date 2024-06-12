package api

import (
	"net/http"
	"strconv"

	"github.com/PoseidonPSE2/code_backend/database"
	"github.com/gin-gonic/gin"
)

// @Summary Show all bottles
// @Description Get all bottles
// @Tags Bottles
// @Accept json
// @Produce json
// @Success 200 {array} database.Bottle
// @Router /bottles [get]
func GetBottles(c *gin.Context) {
	var bottles []database.Bottle
	result := db.Find(&bottles)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, bottles)
}

// @Summary Get bottle by bottle ID
// @Description Get one bottle with the given ID
// @Tags Bottles
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} database.Bottle
// @Router /bottles/{id} [get]
func GetBottleById(c *gin.Context) {
	idStr := c.Param("id")
	if idStr != "" {
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
		c.JSON(http.StatusOK, bottle)
	}
}

// @Summary Get all bottles by user ID
// @Description Get all bottles associated with a specific user
// @Tags Bottles
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {array} database.Bottle
// @Router /bottles/users/{userId} [get]
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
	c.JSON(http.StatusOK, bottles)
}

// @Summary Get bottle preferences by the NFC ID
// @Description Get bottle preferences by the NFC ID
// @Tags Bottles
// @Accept json
// @Produce json
// @Param nfc_id path string true "NFC ID"
// @Success 200 {object} database.Bottle
// @Router /bottles/preferences/{nfcId} [get]
func GetBottlePreferencesByNFCId(c *gin.Context) {
	nfcID := c.Param("nfcId")
	var bottle database.Bottle

	if nfcID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nfcID cannot be empty"})
		return
	}

	if err := db.Where("nfc_id = ?", nfcID).First(&bottle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Row not found for NFC ID"})
		return
	}

	c.JSON(http.StatusOK, bottle)
}

// @Summary Create a bottle
// @Description Create a new bottle
// @Tags Bottles
// @Accept  json
// @Produce  json
// @Param bottle body database.Bottle true "Bottle"
// @Success 201 {object} database.Bottle
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
	c.JSON(http.StatusCreated, bottle)
}

// @Summary Update a bottle
// @Description Update an existing bottle
// @Tags Bottles
// @Accept  json
// @Produce  json
// @Param bottle body database.Bottle true "Bottle"
// @Success 200 {object} database.Bottle
// @Router /bottles [put]
func UpdateBottle(c *gin.Context) {
	var newBottle database.Bottle
	if err := c.ShouldBindJSON(&newBottle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var bottle database.Bottle
	if err := db.Where("id = ?", newBottle.ID).First(&bottle).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	db.Model(&bottle).Updates(newBottle)

	// Workaround to ensure saving of an empty nfc-id
	if newBottle.NFCID == "" {
		db.Model(&bottle).Select("NFCID").Updates(map[string]interface{}{"NFCID": ""})
	}

	// Respond with the updated bottle data
	c.JSON(http.StatusOK, bottle)
}

// @Summary Delete a bottle
// @Description Delete an existing bottle
// @Tags Bottles
// @Accept  json
// @Produce  json
// @Param id query int true "Bottle ID"
// @Success 204
// @Router /bottles/{id} [delete]
func DeleteBottle(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err = db.Where("id = ?", id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	result := db.Delete(&database.Bottle{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
