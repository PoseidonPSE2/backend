package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/PoseidonPSE2/code_backend/database"
	"github.com/gin-gonic/gin"
)

type PostRequestRefillStationProblem struct {
	StationID                 uint    `gorm:"not null" json:"station_id"`
	Title                     string  `gorm:"size:100;not null" json:"title"`
	Description               string  `gorm:"size:255;not null" json:"description"`
	RefillStationProblemImage *[]byte `gorm:"type:TEXT;default:null" json:"problem_image"`
	Status                    string  `gorm:"size:16;not null" json:"status"`
}

// @Summary Show all refill station problems
// @Description Get all refill station problems
// @Tags Refill Station Problems
// @Accept json
// @Produce json
// @Success 200 {array} database.RefillStationProblem
// @Router /refill_station_problems [get]
func GetRefillStationProblems(c *gin.Context) {
	var problems []database.RefillStationProblem
	result := db.Find(&problems)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, problems)
}

// @Summary Show refill station problem by id
// @Description Get refill station problem
// @Tags Refill Station Problems
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} database.RefillStationProblem
// @Router /refill_station_problems/{id} [get]
func GetRefillStationProblemById(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing ID"})
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var problem database.RefillStationProblem
		result := db.First(&problem, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, problem)
	}
}

// @Summary Create a refill station problem
// @Description Create a new refill station problem
// @Tags Refill Station Problems
// @Accept json
// @Produce json
// @Param problem body PostRequestRefillStationProblem true "Refill Station Problem"
// @Success 201 {object} database.RefillStationProblem
// @Router /refill_station_problems [post]
func CreateRefillStationProblem(c *gin.Context) {
	var requestProblem PostRequestRefillStationProblem
	if err := c.ShouldBindJSON(&requestProblem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var problemToInsert = database.RefillStationProblem{
		StationID:   requestProblem.StationID,
		Status:      requestProblem.Status,
		Title:       requestProblem.Title,
		Description: requestProblem.Description,
		Timestamp:   time.Now(),
	}

	if requestProblem.RefillStationProblemImage != nil {
		base64image := EncodeBytesToBase64(*requestProblem.RefillStationProblemImage)
		problemToInsert.RefillStationProblemImage = &base64image
	}

	result := db.Create(&problemToInsert)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusCreated, requestProblem)
}

// @Summary Update a refill station problem
// @Description Update an existing refill station problem
// @Tags Refill Station Problems
// @Accept  json
// @Produce  json
// @Param problem body database.RefillStationProblem true "Refill Station Problem"
// @Success 200 {object} database.RefillStationProblem
// @Router /refill_station_problems [put]
func UpdateRefillStationProblem(c *gin.Context) {
	var requestProblem database.RefillStationProblem
	if err := c.ShouldBindJSON(&requestProblem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var problem database.RefillStationProblem
	if err := db.Where("id = ?", requestProblem.ID).First(&problem).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	requestProblem.Timestamp = time.Now()
	db.Model(&problem).Updates(requestProblem)

	c.JSON(http.StatusOK, requestProblem)
}

// @Summary Delete a refill station problem
// @Description Delete an existing refill station problem
// @Tags Refill Station Problems
// @Accept  json
// @Produce  json
// @Param id path int true "Refill Station Problem ID"
// @Success 204
// @Router /refill_station_problems/{id} [delete]
func DeleteRefillStationProblem(c *gin.Context) {
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
	// Check for record not found error
	var tempProblem database.RefillStationProblem
	result := db.First(&tempProblem, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Problem with ID not found"})
	}

	result = db.Delete(&database.RefillStationProblem{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
