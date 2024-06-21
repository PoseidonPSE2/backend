package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/PoseidonPSE2/code_backend/database"
	"github.com/gin-gonic/gin"
)

// @Summary Show all water transactions
// @Description Get all water transactions
// @Tags Water Transactions
// @Accept json
// @Produce json
// @Success 200 {array} database.WaterTransaction
// @Router /water_transactions [get]
func GetWaterTransactions(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		var transactions []database.WaterTransaction
		result := db.Find(&transactions)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, transactions)
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var transaction database.WaterTransaction
		result := db.First(&transaction, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, transaction)
	}
}

// @Summary Create a water transaction
// @Description Create a new water transaction
// @Tags Water Transactions
// @Accept json
// @Produce json
// @Param transaction body database.WaterTransaction true "Water Transaction"
// @Success 201 {object} database.WaterTransaction
// @Router /water_transactions [post]
func CreateWaterTransaction(c *gin.Context) {
	var transaction database.WaterTransaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	transaction.Timestamp = time.Now()
	result := db.Create(&transaction)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusCreated, transaction)
}

// @Summary Update a water transaction
// @Description Update an existing water transaction
// @Tags Water Transactions
// @Accept json
// @Produce json
// @Param transaction body database.WaterTransaction true "Water Transaction"
// @Success 200 {object} database.WaterTransaction
// @Router /water_transactions [put]
func UpdateWaterTransaction(c *gin.Context) {
	var transaction database.WaterTransaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	transaction.Timestamp = time.Now()
	result := db.Save(&transaction)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, transaction)
}

// @Summary Delete a water transaction
// @Description Delete an existing water transaction
// @Tags Water Transactions
// @Accept json
// @Produce json
// @Param id query int true "Water Transaction ID"
// @Success 204
// @Router /water_transactions [delete]
func DeleteWaterTransaction(c *gin.Context) {
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
	result := db.Delete(&database.WaterTransaction{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
