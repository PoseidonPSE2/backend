package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/PoseidonPSE2/code_backend/database"
	"github.com/gin-gonic/gin"
)

// @Summary Show all consumer test answers
// @Description Get all consumer test answers
// @Tags consumer_test_answers
// @Accept  json
// @Produce  json
// @Success 200 {array} database.ConsumerTestAnswer
// @Router /consumer_test_answers [get]
func GetConsumerTestAnswers(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		var answers []database.ConsumerTestAnswer
		result := db.Find(&answers)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, answers)
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var answer database.ConsumerTestAnswer
		result := db.First(&answer, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, answer)
	}
}

// @Summary Create a consumer test answer
// @Description Create a new consumer test answer
// @Tags consumer_test_answers
// @Accept  json
// @Produce  json
// @Param answer body database.ConsumerTestAnswer true "Consumer Test Answer"
// @Success 201 {object} database.ConsumerTestAnswer
// @Router /consumer_test_answers [post]
func CreateConsumerTestAnswer(c *gin.Context) {
	var answer database.ConsumerTestAnswer
	if err := c.ShouldBindJSON(&answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	answer.Timestamp = time.Now()
	result := db.Create(&answer)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusCreated, answer)
}

// @Summary Update a consumer test answer
// @Description Update an existing consumer test answer
// @Tags consumer_test_answers
// @Accept  json
// @Produce  json
// @Param answer body database.ConsumerTestAnswer true "Consumer Test Answer"
// @Success 200 {object} database.ConsumerTestAnswer
// @Router /consumer_test_answers [put]
func UpdateConsumerTestAnswer(c *gin.Context) {
	var answer database.ConsumerTestAnswer
	if err := c.ShouldBindJSON(&answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	answer.Timestamp = time.Now()
	result := db.Save(&answer)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, answer)
}

// @Summary Delete a consumer test answer
// @Description Delete an existing consumer test answer
// @Tags consumer_test_answers
// @Accept  json
// @Produce  json
// @Param id query int true "Consumer Test Answer ID"
// @Success 204
// @Router /consumer_test_answers [delete]
func DeleteConsumerTestAnswer(c *gin.Context) {
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
	result := db.Delete(&database.ConsumerTestAnswer{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
