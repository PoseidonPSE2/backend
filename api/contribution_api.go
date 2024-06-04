package api

import (
	"net/http"
	"strconv"

	"github.com/PoseidonPSE2/code_backend/database"
	"github.com/gin-gonic/gin"
)

func calculateSavings(volume int) (float64, float64) {
	const moneyFactor = 0.50
	const trashFactor = 0.10

	savedMoney := float64(volume) * moneyFactor / 1000 // Convert volume to liters
	savedTrash := float64(volume) * trashFactor / 1000

	return savedMoney, savedTrash
}

// @Summary Get user contribution
// @Description Get the total water amount and savings for a user
// @Tags contribution
// @Accept  json
// @Produce  json
// @Param userId query int true "User ID"
// @Success 200 {object} database.WaterTransaction
// @Router /contribution/user [get]
func GetContributionByUser(c *gin.Context) {
	userIdStr := c.Query("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	var totalVolume int64
	var totalFillings int64

	db.Model(&database.WaterTransaction{}).Where("user_id = ?", userId).Count(&totalFillings)
	db.Model(&database.WaterTransaction{}).Where("user_id = ?", userId).Select("sum(volume)").Row().Scan(&totalVolume)

	savedMoney, savedTrash := calculateSavings(int(totalVolume))

	response := gin.H{
		"amountFillings": totalFillings,
		"amountWater":    totalVolume,
		"savedMoney":     savedMoney,
		"savedTrash":     savedTrash,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get community contribution
// @Description Get the total water amount and savings for the community
// @Tags contribution
// @Accept  json
// @Produce  json
// @Success 200 {object} database.WaterTransaction
// @Router /contribution/community [get]
func GetContributionCommunity(c *gin.Context) {
	var totalVolume int64
	var totalFillings int64
	var totalUsers int64

	db.Model(&database.WaterTransaction{}).Count(&totalFillings)
	db.Model(&database.WaterTransaction{}).Select("sum(volume)").Row().Scan(&totalVolume)
	db.Model(&database.User{}).Count(&totalUsers)

	savedMoney, savedTrash := calculateSavings(int(totalVolume))

	response := gin.H{
		"amountFillings": totalFillings,
		"amountWater":    totalVolume,
		"savedMoney":     savedMoney,
		"savedTrash":     savedTrash,
		"amountUser":     totalUsers,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get contribution by station type
// @Description Get the number of smart and manual refill stations
// @Tags contribution
// @Accept  json
// @Produce  json
// @Success 200 {object} database.RefillStation
// @Router /contribution/kl [get]
func GetContributionKL(c *gin.Context) {
	var smartStations int64
	var manualStations int64

	db.Model(&database.RefillStation{}).Where("type = ?", "Smart").Count(&smartStations)
	db.Model(&database.RefillStation{}).Where("type = ?", "Manual").Count(&manualStations)

	response := gin.H{
		"amountRefillStationSmart":  smartStations,
		"amountRefillStationManual": manualStations,
	}

	c.JSON(http.StatusOK, response)
}
