package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/PoseidonPSE2/code_backend/database"
	"github.com/gin-gonic/gin"
)

// @Summary Show all refill station problems
// @Description Get all refill station problems
// @Tags refill_station_problems
// @Accept  json
// @Produce  json
// @Success 200 {array} database.RefillStationProblem
// @Router /refill_station_problems [get]
func GetRefillStationProblems(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		var problems []database.RefillStationProblem
		result := db.Find(&problems)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, problems)
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
// @Tags refill_station_problems
// @Accept  json
// @Produce  json
// @Param problem body database.RefillStationProblem true "Refill Station Problem"
// @Success 201 {object} database.RefillStationProblem
// @Router /refill_station_problems [post]
func CreateRefillStationProblem(c *gin.Context) {
	var problem database.RefillStationProblem
	if err := c.ShouldBindJSON(&problem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	problem.Timestamp = time.Now()
	result := db.Create(&problem)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusCreated, problem)
}

// @Summary Update a refill station problem
// @Description Update an existing refill station problem
// @Tags refill_station_problems
// @Accept  json
// @Produce  json
// @Param problem body database.RefillStationProblem true "Refill Station Problem"
// @Success 200 {object} database.RefillStationProblem
// @Router /refill_station_problems [put]
func UpdateRefillStationProblem(c *gin.Context) {
	var problem database.RefillStationProblem
	if err := c.ShouldBindJSON(&problem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	problem.Timestamp = time.Now()
	result := db.Save(&problem)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, problem)
}

// @Summary Delete a refill station problem
// @Description Delete an existing refill station problem
// @Tags refill_station_problems
// @Accept  json
// @Produce  json
// @Param id query int true "Refill Station Problem ID"
// @Success 204
// @Router /refill_station_problems [delete]
func DeleteRefillStationProblem(c *gin.Context) {
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
	result := db.Delete(&database.RefillStationProblem{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
