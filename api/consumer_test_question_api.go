package api

import (
	"net/http"
	"strconv"

	"github.com/PoseidonPSE2/code_backend/database"
	"github.com/gin-gonic/gin"
)

// ConsumerTestQuestionResponse represents a consumer test question in the response
type ConsumerTestQuestionResponse struct {
	ID       uint     `json:"id"`
	Text     string   `json:"text"`
	MinValue *float32 `json:"min_value"`
	MaxValue *float32 `json:"max_value"`
	TestID   uint     `json:"test_id"`
}

// CreateConsumerTestQuestionRequest represents a request to create a consumer test question
type CreateConsumerTestQuestionRequest struct {
	Text     string   `json:"text"`
	MinValue *float32 `json:"min_value"`
	MaxValue *float32 `json:"max_value"`
	TestID   uint     `json:"test_id"`
}

// UpdateConsumerTestQuestionRequest represents a request to update a consumer test question
type UpdateConsumerTestQuestionRequest struct {
	ID       uint     `json:"id"`
	Text     string   `json:"text"`
	MinValue *float32 `json:"min_value"`
	MaxValue *float32 `json:"max_value"`
	TestID   uint     `json:"test_id"`
}

// @Summary Show all consumer test questions
// @Description Get all consumer test questions
// @Tags consumer_test_questions
// @Accept  json
// @Produce  json
// @Success 200 {array} ConsumerTestQuestionResponse
// @Router /consumer_test_questions [get]
func GetConsumerTestQuestions(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		var questions []database.ConsumerTestQuestion
		result := db.Find(&questions)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		respondWithJSON(c, http.StatusOK, questions)
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var question database.ConsumerTestQuestion
		result := db.First(&question, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
			return
		}
		respondWithJSON(c, http.StatusOK, question)
	}
}

// @Summary Create a consumer test question
// @Description Create a new consumer test question
// @Tags consumer_test_questions
// @Accept  json
// @Produce  json
// @Param question body CreateConsumerTestQuestionRequest true "Consumer Test Question"
// @Success 201 {object} ConsumerTestQuestionResponse
// @Router /consumer_test_questions [post]
func CreateConsumerTestQuestion(c *gin.Context) {
	var question database.ConsumerTestQuestion
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := db.Create(&question)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	respondWithJSON(c, http.StatusCreated, question)
}

// @Summary Update a consumer test question
// @Description Update an existing consumer test question
// @Tags consumer_test_questions
// @Accept  json
// @Produce  json
// @Param question body UpdateConsumerTestQuestionRequest true "Consumer Test Question"
// @Success 200 {object} ConsumerTestQuestionResponse
// @Router /consumer_test_questions [put]
func UpdateConsumerTestQuestion(c *gin.Context) {
	var question database.ConsumerTestQuestion
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := db.Save(&question)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	respondWithJSON(c, http.StatusOK, question)
}

// @Summary Delete a consumer test question
// @Description Delete an existing consumer test question
// @Tags consumer_test_questions
// @Accept  json
// @Produce  json
// @Param id query int true "Consumer Test Question ID"
// @Success 204
// @Router /consumer_test_questions [delete]
func DeleteConsumerTestQuestion(c *gin.Context) {
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
	result := db.Delete(&database.ConsumerTestQuestion{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
