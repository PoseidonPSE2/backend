package api

import (
	"net/http"
	"strconv"

	"github.com/PoseidonPSE2/code_backend/database"
	"github.com/gin-gonic/gin"
)

// ConsumerTestResponse represents a consumer test in the response
type ConsumerTestResponse struct {
	ID uint `json:"id"`
}

// CreateConsumerTestRequest represents a request to create a consumer test
type CreateConsumerTestRequest struct {
	Questions []database.ConsumerTestQuestion `json:"questions"`
}

// UpdateConsumerTestRequest represents a request to update a consumer test
type UpdateConsumerTestRequest struct {
	ID        uint                            `json:"id"`
	Questions []database.ConsumerTestQuestion `json:"questions"`
}

// @Summary Show all consumer tests
// @Description Get all consumer tests
// @Tags consumer_tests
// @Accept  json
// @Produce  json
// @Success 200 {array} ConsumerTestResponse
// @Router /consumer_tests [get]
func GetConsumerTests(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		var tests []database.ConsumerTest
		result := db.Find(&tests)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		respondWithJSON(c, http.StatusOK, tests)
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var test database.ConsumerTest
		result := db.First(&test, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
			return
		}
		respondWithJSON(c, http.StatusOK, test)
	}
}

// @Summary Create a consumer test
// @Description Create a new consumer test
// @Tags consumer_tests
// @Accept  json
// @Produce  json
// @Param test body CreateConsumerTestRequest true "Consumer Test"
// @Success 201 {object} ConsumerTestResponse
// @Router /consumer_tests [post]
func CreateConsumerTest(c *gin.Context) {
	var test database.ConsumerTest
	if err := c.ShouldBindJSON(&test); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := db.Create(&test)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	respondWithJSON(c, http.StatusCreated, test)
}

// @Summary Update a consumer test
// @Description Update an existing consumer test
// @Tags consumer_tests
// @Accept  json
// @Produce  json
// @Param test body UpdateConsumerTestRequest true "Consumer Test"
// @Success 200 {object} ConsumerTestResponse
// @Router /consumer_tests [put]
func UpdateConsumerTest(c *gin.Context) {
	var test database.ConsumerTest
	if err := c.ShouldBindJSON(&test); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := db.Save(&test)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	respondWithJSON(c, http.StatusOK, test)
}

// @Summary Delete a consumer test
// @Description Delete an existing consumer test
// @Tags consumer_tests
// @Accept  json
// @Produce  json
// @Param id query int true "Consumer Test ID"
// @Success 204
// @Router /consumer_tests [delete]
func DeleteConsumerTest(c *gin.Context) {
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
	result := db.Delete(&database.ConsumerTest{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
