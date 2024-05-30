package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/PoseidonPSE2/code_backend/database"

	_ "github.com/PoseidonPSE2/code_backend/docs" // swagger docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Initialize the database connection
var db *gorm.DB

func init() {
	var err error
	log.Print("Starting application")

	connectionString := "postgres://postgres:rnJpE83UKr1MyF8@poseidon-database.flycast:5432"

	db, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	fmt.Println("Successfully connected to the database!")

	log.Print("Schema migration starting")

	// Migrate the schema
	db.AutoMigrate(&database.User{}, &database.Bottle{}, &database.ConsumerTest{}, &database.ConsumerTestQuestion{},
		&database.ConsumerTestAnswer{}, &database.RefillStation{}, &database.RefillStationReview{},
		&database.RefillStationProblem{}, &database.WaterTransaction{}, &database.Like{})
	log.Print("Schema migration done")
}

func respondWithJSON(c *gin.Context, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(status, "application/json", response)
}

// UserResponse represents a user in the response
type UserResponse struct {
	ID        uint    `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     *string `json:"email"`
}

// @Summary Show all users
// @Description Get all users
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} UserResponse
// @Router /users [get]
func getUsers(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		var users []database.User
		result := db.Find(&users)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		respondWithJSON(c, http.StatusOK, users)
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var user database.User
		result := db.First(&user, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
			return
		}
		respondWithJSON(c, http.StatusOK, user)
	}
}

// CreateUserRequest represents a request to create a user
type CreateUserRequest struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     *string `json:"email"`
}

// @Summary Create a user
// @Description Create a new user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body CreateUserRequest true "User"
// @Success 201 {object} UserResponse
// @Router /users [post]
func createUser(c *gin.Context) {
	var user database.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := db.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	respondWithJSON(c, http.StatusCreated, user)
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	ID        uint    `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     *string `json:"email"`
}

// @Summary Update a user
// @Description Update an existing user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body UpdateUserRequest true "User"
// @Success 200 {object} UserResponse
// @Router /users [put]
func updateUser(c *gin.Context) {
	var user database.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := db.Save(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	respondWithJSON(c, http.StatusOK, user)
}

// @Summary Delete a user
// @Description Delete an existing user
// @Tags users
// @Accept  json
// @Produce  json
// @Param id query int true "User ID"
// @Success 204
// @Router /users [delete]
func deleteUser(c *gin.Context) {
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
	result := db.Delete(&database.User{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

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

// @Summary Show all bottles
// @Description Get all bottles
// @Tags bottles
// @Accept  json
// @Produce  json
// @Success 200 {array} BottleResponse
// @Router /bottles [get]
func getBottles(c *gin.Context) {
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

// @Summary Create a bottle
// @Description Create a new bottle
// @Tags bottles
// @Accept  json
// @Produce  json
// @Param bottle body CreateBottleRequest true "Bottle"
// @Success 201 {object} BottleResponse
// @Router /bottles [post]
func createBottle(c *gin.Context) {
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

// @Summary Update a bottle
// @Description Update an existing bottle
// @Tags bottles
// @Accept  json
// @Produce  json
// @Param bottle body UpdateBottleRequest true "Bottle"
// @Success 200 {object} BottleResponse
// @Router /bottles [put]
func updateBottle(c *gin.Context) {
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
func deleteBottle(c *gin.Context) {
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

// ConsumerTestResponse represents a consumer test in the response
type ConsumerTestResponse struct {
	ID uint `json:"id"`
}

// @Summary Show all consumer tests
// @Description Get all consumer tests
// @Tags consumer_tests
// @Accept  json
// @Produce  json
// @Success 200 {array} ConsumerTestResponse
// @Router /consumer_tests [get]
func getConsumerTests(c *gin.Context) {
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

// CreateConsumerTestRequest represents a request to create a consumer test
type CreateConsumerTestRequest struct {
	Questions []database.ConsumerTestQuestion `json:"questions"`
}

// @Summary Create a consumer test
// @Description Create a new consumer test
// @Tags consumer_tests
// @Accept  json
// @Produce  json
// @Param test body CreateConsumerTestRequest true "Consumer Test"
// @Success 201 {object} ConsumerTestResponse
// @Router /consumer_tests [post]
func createConsumerTest(c *gin.Context) {
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

// UpdateConsumerTestRequest represents a request to update a consumer test
type UpdateConsumerTestRequest struct {
	ID        uint                            `json:"id"`
	Questions []database.ConsumerTestQuestion `json:"questions"`
}

// @Summary Update a consumer test
// @Description Update an existing consumer test
// @Tags consumer_tests
// @Accept  json
// @Produce  json
// @Param test body UpdateConsumerTestRequest true "Consumer Test"
// @Success 200 {object} ConsumerTestResponse
// @Router /consumer_tests [put]
func updateConsumerTest(c *gin.Context) {
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
func deleteConsumerTest(c *gin.Context) {
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

// ConsumerTestQuestionResponse represents a consumer test question in the response
type ConsumerTestQuestionResponse struct {
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
func getConsumerTestQuestions(c *gin.Context) {
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

// CreateConsumerTestQuestionRequest represents a request to create a consumer test question
type CreateConsumerTestQuestionRequest struct {
	Text     string   `json:"text"`
	MinValue *float32 `json:"min_value"`
	MaxValue *float32 `json:"max_value"`
	TestID   uint     `json:"test_id"`
}

// @Summary Create a consumer test question
// @Description Create a new consumer test question
// @Tags consumer_test_questions
// @Accept  json
// @Produce  json
// @Param question body CreateConsumerTestQuestionRequest true "Consumer Test Question"
// @Success 201 {object} ConsumerTestQuestionResponse
// @Router /consumer_test_questions [post]
func createConsumerTestQuestion(c *gin.Context) {
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

// UpdateConsumerTestQuestionRequest represents a request to update a consumer test question
type UpdateConsumerTestQuestionRequest struct {
	ID       uint     `json:"id"`
	Text     string   `json:"text"`
	MinValue *float32 `json:"min_value"`
	MaxValue *float32 `json:"max_value"`
	TestID   uint     `json:"test_id"`
}

// @Summary Update a consumer test question
// @Description Update an existing consumer test question
// @Tags consumer_test_questions
// @Accept  json
// @Produce  json
// @Param question body UpdateConsumerTestQuestionRequest true "Consumer Test Question"
// @Success 200 {object} ConsumerTestQuestionResponse
// @Router /consumer_test_questions [put]
func updateConsumerTestQuestion(c *gin.Context) {
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
func deleteConsumerTestQuestion(c *gin.Context) {
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

// ConsumerTestAnswerResponse represents a consumer test answer in the response
type ConsumerTestAnswerResponse struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	QuestionID uint      `json:"question_id"`
	Answer     float32   `json:"answer"`
	Timestamp  time.Time `json:"timestamp"`
}

// @Summary Show all consumer test answers
// @Description Get all consumer test answers
// @Tags consumer_test_answers
// @Accept  json
// @Produce  json
// @Success 200 {array} ConsumerTestAnswerResponse
// @Router /consumer_test_answers [get]
func getConsumerTestAnswers(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		var answers []database.ConsumerTestAnswer
		result := db.Find(&answers)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		respondWithJSON(c, http.StatusOK, answers)
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
		respondWithJSON(c, http.StatusOK, answer)
	}
}

// CreateConsumerTestAnswerRequest represents a request to create a consumer test answer
type CreateConsumerTestAnswerRequest struct {
	UserID     uint    `json:"user_id"`
	QuestionID uint    `json:"question_id"`
	Answer     float32 `json:"answer"`
}

// @Summary Create a consumer test answer
// @Description Create a new consumer test answer
// @Tags consumer_test_answers
// @Accept  json
// @Produce  json
// @Param answer body CreateConsumerTestAnswerRequest true "Consumer Test Answer"
// @Success 201 {object} ConsumerTestAnswerResponse
// @Router /consumer_test_answers [post]
func createConsumerTestAnswer(c *gin.Context) {
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
	respondWithJSON(c, http.StatusCreated, answer)
}

// UpdateConsumerTestAnswerRequest represents a request to update a consumer test answer
type UpdateConsumerTestAnswerRequest struct {
	ID         uint    `json:"id"`
	UserID     uint    `json:"user_id"`
	QuestionID uint    `json:"question_id"`
	Answer     float32 `json:"answer"`
}

// @Summary Update a consumer test answer
// @Description Update an existing consumer test answer
// @Tags consumer_test_answers
// @Accept  json
// @Produce  json
// @Param answer body UpdateConsumerTestAnswerRequest true "Consumer Test Answer"
// @Success 200 {object} ConsumerTestAnswerResponse
// @Router /consumer_test_answers [put]
func updateConsumerTestAnswer(c *gin.Context) {
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
	respondWithJSON(c, http.StatusOK, answer)
}

// @Summary Delete a consumer test answer
// @Description Delete an existing consumer test answer
// @Tags consumer_test_answers
// @Accept  json
// @Produce  json
// @Param id query int true "Consumer Test Answer ID"
// @Success 204
// @Router /consumer_test_answers [delete]
func deleteConsumerTestAnswer(c *gin.Context) {
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

// RefillStationResponse represents a refill station in the response
type RefillStationResponse struct {
	ID                uint    `json:"id"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	Address           string  `json:"address"`
	LikeCounter       int     `json:"like_counter"`
	WaterSource       string  `json:"water_source"`
	OpeningTimes      string  `json:"opening_times"`
	Active            bool    `json:"active"`
	Type              string  `json:"type"`
	OfferedWaterTypes string  `json:"offered_water_types"`
	ImagePath         *string `json:"image_path"`
}

// @Summary Show all refill stations
// @Description Get all refill stations
// @Tags refill_stations
// @Accept  json
// @Produce  json
// @Success 200 {array} RefillStationResponse
// @Router /refill_stations [get]
func getRefillStations(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		var stations []database.RefillStation
		result := db.Find(&stations)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		respondWithJSON(c, http.StatusOK, stations)
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var station database.RefillStation
		result := db.First(&station, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
			return
		}
		respondWithJSON(c, http.StatusOK, station)
	}
}

// CreateRefillStationRequest represents a request to create a refill station
type CreateRefillStationRequest struct {
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	Address           string  `json:"address"`
	LikeCounter       int     `json:"like_counter"`
	WaterSource       string  `json:"water_source"`
	OpeningTimes      string  `json:"opening_times"`
	Active            bool    `json:"active"`
	Type              string  `json:"type"`
	OfferedWaterTypes string  `json:"offered_water_types"`
	ImagePath         *string `json:"image_path"`
}

// @Summary Create a refill station
// @Description Create a new refill station
// @Tags refill_stations
// @Accept  json
// @Produce  json
// @Param station body CreateRefillStationRequest true "Refill Station"
// @Success 201 {object} RefillStationResponse
// @Router /refill_stations [post]
func createRefillStation(c *gin.Context) {
	var station database.RefillStation
	if err := c.ShouldBindJSON(&station); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := db.Create(&station)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	respondWithJSON(c, http.StatusCreated, station)
}

// UpdateRefillStationRequest represents a request to update a refill station
type UpdateRefillStationRequest struct {
	ID                uint    `json:"id"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	Address           string  `json:"address"`
	LikeCounter       int     `json:"like_counter"`
	WaterSource       string  `json:"water_source"`
	OpeningTimes      string  `json:"opening_times"`
	Active            bool    `json:"active"`
	Type              string  `json:"type"`
	OfferedWaterTypes string  `json:"offered_water_types"`
	ImagePath         *string `json:"image_path"`
}

// @Summary Update a refill station
// @Description Update an existing refill station
// @Tags refill_stations
// @Accept  json
// @Produce  json
// @Param station body UpdateRefillStationRequest true "Refill Station"
// @Success 200 {object} RefillStationResponse
// @Router /refill_stations [put]
func updateRefillStation(c *gin.Context) {
	var station database.RefillStation
	if err := c.ShouldBindJSON(&station); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := db.Save(&station)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	respondWithJSON(c, http.StatusOK, station)
}

// @Summary Delete a refill station
// @Description Delete an existing refill station
// @Tags refill_stations
// @Accept  json
// @Produce  json
// @Param id query int true "Refill Station ID"
// @Success 204
// @Router /refill_stations [delete]
func deleteRefillStation(c *gin.Context) {
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
	result := db.Delete(&database.RefillStation{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// RefillStationReviewResponse represents a refill station review in the response
type RefillStationReviewResponse struct {
	ID            uint      `json:"id"`
	StationID     uint      `json:"station_id"`
	UserID        uint      `json:"user_id"`
	Cleanness     int       `json:"cleanness"`
	Accessibility int       `json:"accessibility"`
	WaterQuality  int       `json:"water_quality"`
	Timestamp     time.Time `json:"timestamp"`
}

// @Summary Show all refill station reviews
// @Description Get all refill station reviews
// @Tags refill_station_reviews
// @Accept  json
// @Produce  json
// @Success 200 {array} RefillStationReviewResponse
// @Router /refill_station_reviews [get]
func getRefillStationReviews(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		var reviews []database.RefillStationReview
		result := db.Find(&reviews)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		respondWithJSON(c, http.StatusOK, reviews)
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var review database.RefillStationReview
		result := db.First(&review, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
			return
		}
		respondWithJSON(c, http.StatusOK, review)
	}
}

// CreateRefillStationReviewRequest represents a request to create a refill station review
type CreateRefillStationReviewRequest struct {
	StationID     uint `json:"station_id"`
	UserID        uint `json:"user_id"`
	Cleanness     int  `json:"cleanness"`
	Accessibility int  `json:"accessibility"`
	WaterQuality  int  `json:"water_quality"`
}

// @Summary Create a refill station review
// @Description Create a new refill station review
// @Tags refill_station_reviews
// @Accept  json
// @Produce  json
// @Param review body CreateRefillStationReviewRequest true "Refill Station Review"
// @Success 201 {object} RefillStationReviewResponse
// @Router /refill_station_reviews [post]
func createRefillStationReview(c *gin.Context) {
	var review database.RefillStationReview
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	review.Timestamp = time.Now()
	result := db.Create(&review)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	respondWithJSON(c, http.StatusCreated, review)
}

// UpdateRefillStationReviewRequest represents a request to update a refill station review
type UpdateRefillStationReviewRequest struct {
	ID            uint      `json:"id"`
	StationID     uint      `json:"station_id"`
	UserID        uint      `json:"user_id"`
	Cleanness     int       `json:"cleanness"`
	Accessibility int       `json:"accessibility"`
	WaterQuality  int       `json:"water_quality"`
	Timestamp     time.Time `json:"timestamp"`
}

// @Summary Update a refill station review
// @Description Update an existing refill station review
// @Tags refill_station_reviews
// @Accept  json
// @Produce  json
// @Param review body UpdateRefillStationReviewRequest true "Refill Station Review"
// @Success 200 {object} RefillStationReviewResponse
// @Router /refill_station_reviews [put]
func updateRefillStationReview(c *gin.Context) {
	var review database.RefillStationReview
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	review.Timestamp = time.Now()
	result := db.Save(&review)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	respondWithJSON(c, http.StatusOK, review)
}

// @Summary Delete a refill station review
// @Description Delete an existing refill station review
// @Tags refill_station_reviews
// @Accept  json
// @Produce  json
// @Param id query int true "Refill Station Review ID"
// @Success 204
// @Router /refill_station_reviews [delete]
func deleteRefillStationReview(c *gin.Context) {
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
	result := db.Delete(&database.RefillStationReview{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// RefillStationProblemResponse represents a refill station problem in the response
type RefillStationProblemResponse struct {
	ID          uint      `json:"id"`
	StationID   uint      `json:"station_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	LinkToMedia *string   `json:"link_to_media"`
	Timestamp   time.Time `json:"timestamp"`
}

// @Summary Show all refill station problems
// @Description Get all refill station problems
// @Tags refill_station_problems
// @Accept  json
// @Produce  json
// @Success 200 {array} RefillStationProblemResponse
// @Router /refill_station_problems [get]
func getRefillStationProblems(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		var problems []database.RefillStationProblem
		result := db.Find(&problems)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		respondWithJSON(c, http.StatusOK, problems)
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
		respondWithJSON(c, http.StatusOK, problem)
	}
}

// CreateRefillStationProblemRequest represents a request to create a refill station problem
type CreateRefillStationProblemRequest struct {
	StationID   uint    `json:"station_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	LinkToMedia *string `json:"link_to_media"`
}

// @Summary Create a refill station problem
// @Description Create a new refill station problem
// @Tags refill_station_problems
// @Accept  json
// @Produce  json
// @Param problem body CreateRefillStationProblemRequest true "Refill Station Problem"
// @Success 201 {object} RefillStationProblemResponse
// @Router /refill_station_problems [post]
func createRefillStationProblem(c *gin.Context) {
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
	respondWithJSON(c, http.StatusCreated, problem)
}

// UpdateRefillStationProblemRequest represents a request to update a refill station problem
type UpdateRefillStationProblemRequest struct {
	ID          uint      `json:"id"`
	StationID   uint      `json:"station_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	LinkToMedia *string   `json:"link_to_media"`
	Timestamp   time.Time `json:"timestamp"`
}

// @Summary Update a refill station problem
// @Description Update an existing refill station problem
// @Tags refill_station_problems
// @Accept  json
// @Produce  json
// @Param problem body UpdateRefillStationProblemRequest true "Refill Station Problem"
// @Success 200 {object} RefillStationProblemResponse
// @Router /refill_station_problems [put]
func updateRefillStationProblem(c *gin.Context) {
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
	respondWithJSON(c, http.StatusOK, problem)
}

// @Summary Delete a refill station problem
// @Description Delete an existing refill station problem
// @Tags refill_station_problems
// @Accept  json
// @Produce  json
// @Param id query int true "Refill Station Problem ID"
// @Success 204
// @Router /refill_station_problems [delete]
func deleteRefillStationProblem(c *gin.Context) {
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

// WaterTransactionResponse represents a water transaction in the response
type WaterTransactionResponse struct {
	ID        uint      `json:"id"`
	StationID uint      `json:"station_id"`
	BottleID  *uint     `json:"bottle_id"`
	UserID    *uint     `json:"user_id"`
	Volume    int       `json:"volume"`
	WaterType string    `json:"water_type"`
	Timestamp time.Time `json:"timestamp"`
	Guest     bool      `json:"guest"`
}

// @Summary Show all water transactions
// @Description Get all water transactions
// @Tags water_transactions
// @Accept  json
// @Produce  json
// @Success 200 {array} WaterTransactionResponse
// @Router /water_transactions [get]
func getWaterTransactions(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		var transactions []database.WaterTransaction
		result := db.Find(&transactions)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		respondWithJSON(c, http.StatusOK, transactions)
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
		respondWithJSON(c, http.StatusOK, transaction)
	}
}

// CreateWaterTransactionRequest represents a request to create a water transaction
type CreateWaterTransactionRequest struct {
	StationID uint   `json:"station_id"`
	BottleID  *uint  `json:"bottle_id"`
	UserID    *uint  `json:"user_id"`
	Volume    int    `json:"volume"`
	WaterType string `json:"water_type"`
	Guest     bool   `json:"guest"`
}

// @Summary Create a water transaction
// @Description Create a new water transaction
// @Tags water_transactions
// @Accept  json
// @Produce  json
// @Param transaction body CreateWaterTransactionRequest true "Water Transaction"
// @Success 201 {object} WaterTransactionResponse
// @Router /water_transactions [post]
func createWaterTransaction(c *gin.Context) {
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
	respondWithJSON(c, http.StatusCreated, transaction)
}

// UpdateWaterTransactionRequest represents a request to update a water transaction
type UpdateWaterTransactionRequest struct {
	ID        uint      `json:"id"`
	StationID uint      `json:"station_id"`
	BottleID  *uint     `json:"bottle_id"`
	UserID    *uint     `json:"user_id"`
	Volume    int       `json:"volume"`
	WaterType string    `json:"water_type"`
	Timestamp time.Time `json:"timestamp"`
	Guest     bool      `json:"guest"`
}

// @Summary Update a water transaction
// @Description Update an existing water transaction
// @Tags water_transactions
// @Accept  json
// @Produce  json
// @Param transaction body UpdateWaterTransactionRequest true "Water Transaction"
// @Success 200 {object} WaterTransactionResponse
// @Router /water_transactions [put]
func updateWaterTransaction(c *gin.Context) {
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
	respondWithJSON(c, http.StatusOK, transaction)
}

// @Summary Delete a water transaction
// @Description Delete an existing water transaction
// @Tags water_transactions
// @Accept  json
// @Produce  json
// @Param id query int true "Water Transaction ID"
// @Success 204
// @Router /water_transactions [delete]
func deleteWaterTransaction(c *gin.Context) {
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

// LikeResponse represents a like in the response
type LikeResponse struct {
	ID        uint `json:"id"`
	StationID uint `json:"station_id"`
	UserID    uint `json:"user_id"`
}

// @Summary Show all likes
// @Description Get all likes
// @Tags likes
// @Accept  json
// @Produce  json
// @Success 200 {array} LikeResponse
// @Router /likes [get]
func getLikes(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		var likes []database.Like
		result := db.Find(&likes)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		respondWithJSON(c, http.StatusOK, likes)
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var like database.Like
		result := db.First(&like, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
			return
		}
		respondWithJSON(c, http.StatusOK, like)
	}
}

// CreateLikeRequest represents a request to create a like
type CreateLikeRequest struct {
	StationID uint `json:"station_id"`
	UserID    uint `json:"user_id"`
}

// @Summary Create a like
// @Description Create a new like
// @Tags likes
// @Accept  json
// @Produce  json
// @Param like body CreateLikeRequest true "Like"
// @Success 201 {object} LikeResponse
// @Router /likes [post]
func createLike(c *gin.Context) {
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
	respondWithJSON(c, http.StatusCreated, like)
}

// UpdateLikeRequest represents a request to update a like
type UpdateLikeRequest struct {
	ID        uint `json:"id"`
	StationID uint `json:"station_id"`
	UserID    uint `json:"user_id"`
}

// @Summary Update a like
// @Description Update an existing like
// @Tags likes
// @Accept  json
// @Produce  json
// @Param like body UpdateLikeRequest true "Like"
// @Success 200 {object} LikeResponse
// @Router /likes [put]
func updateLike(c *gin.Context) {
	var like database.Like
	if err := c.ShouldBindJSON(&like); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := db.Save(&like)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	respondWithJSON(c, http.StatusOK, like)
}

// @Summary Delete a like
// @Description Delete an existing like
// @Tags likes
// @Accept  json
// @Produce  json
// @Param id query int true "Like ID"
// @Success 204
// @Router /likes [delete]
func deleteLike(c *gin.Context) {
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
	result := db.Delete(&database.Like{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Get a refill station by ID
// @Description Get a refill station by its ID
// @Tags refill_stations
// @Accept  json
// @Produce  json
// @Param id path int true "Refill Station ID"
// @Success 200 {object} RefillStationResponse
// @Router /refill_stations/{id} [get]
func getRefillStationById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var station database.RefillStation
	result := db.First(&station, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		return
	}
	respondWithJSON(c, http.StatusOK, station)
}

// @Summary Get the average review score for a refill station
// @Description Get the average review score for a refill station by its ID
// @Tags refill_station_reviews
// @Accept  json
// @Produce  json
// @Param id query int true "Refill Station ID"
// @Success 200 {number} float64
// @Router /refill_station_review/average [get]
func getRefillStationReview(c *gin.Context) {
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

	var reviews []database.RefillStationReview
	result := db.Where("station_id = ?", id).Find(&reviews)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if len(reviews) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No reviews found"})
		return
	}

	var totalCleanness, totalAccessibility, totalWaterQuality float64
	for _, review := range reviews {
		totalCleanness += float64(review.Cleanness)
		totalAccessibility += float64(review.Accessibility)
		totalWaterQuality += float64(review.WaterQuality)
	}

	average := (totalCleanness + totalAccessibility + totalWaterQuality) / (float64(len(reviews)) * 3)
	average = math.Round(average*10) / 10 // Round to 1 decimal place

	c.JSON(http.StatusOK, gin.H{"average": average})
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// IsLikedResponse represents a response indicating if a user likes a refill station
type IsLikedResponse struct {
	IsLiked bool `json:"isLiked"`
}

// @Summary Check if a user likes a refill station
// @Description Check if a specific user likes a specific refill station
// @Tags likes
// @Accept  json
// @Produce  json
// @Param refillstationId query int true "Refill Station ID"
// @Param userId query int true "User ID"
// @Success 200 {object} IsLikedResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /refillstation_like [get]
func getRefillstationLike(c *gin.Context) {
	refillstationIdStr := c.Query("refillstationId")
	userIdStr := c.Query("userId")

	if refillstationIdStr == "" || userIdStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{"refillstationId and userId are required"})
		return
	}

	refillstationId, err := strconv.Atoi(refillstationIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{"Invalid refillstationId"})
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{"Invalid userId"})
		return
	}

	var like database.Like
	result := db.Where("station_id = ? AND user_id = ?", refillstationId, userId).First(&like)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, ErrorResponse{result.Error.Error()})
		return
	}

	isLiked := result.RowsAffected > 0
	respondWithJSON(c, http.StatusOK, IsLikedResponse{IsLiked: isLiked})
}

// @Summary Create a like for a refill station
// @Description Create a like for a specific refill station by a specific user if it doesn't already exist
// @Tags likes
// @Accept  json
// @Produce  json
// @Param refillstationId query int true "Refill Station ID"
// @Param userId query int true "User ID"
// @Success 201 {object} LikeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /refillstation_like [post]
func postRefillstationLike(c *gin.Context) {
	refillstationIdStr := c.Query("refillstationId")
	userIdStr := c.Query("userId")

	if refillstationIdStr == "" || userIdStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{"refillstationId and userId are required"})
		return
	}

	refillstationId, err := strconv.Atoi(refillstationIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{"Invalid refillstationId"})
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{"Invalid userId"})
		return
	}

	var like database.Like
	result := db.Where("station_id = ? AND user_id = ?", refillstationId, userId).First(&like)
	if result.Error == nil {
		c.JSON(http.StatusConflict, ErrorResponse{"Like already exists"})
		return
	}

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, ErrorResponse{result.Error.Error()})
		return
	}

	newLike := database.Like{
		StationID: uint(refillstationId),
		UserID:    uint(userId),
	}

	result = db.Create(&newLike)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{result.Error.Error()})
		return
	}

	respondWithJSON(c, http.StatusCreated, newLike)
}

// ContributionUserResponse represents the user contribution response
type ContributionUserResponse struct {
	AmountFillings int64   `json:"amountFillings"`
	AmountWater    int64   `json:"amountWater"`
	SavedMoney     float64 `json:"savedMoney"`
	SavedTrash     float64 `json:"savedTrash"`
}

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
// @Success 200 {object} ContributionUserResponse
// @Router /contribution/user [get]
func getContributionByUser(c *gin.Context) {
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

	response := ContributionUserResponse{
		AmountFillings: totalFillings,
		AmountWater:    totalVolume,
		SavedMoney:     savedMoney,
		SavedTrash:     savedTrash,
	}

	respondWithJSON(c, http.StatusOK, response)
}

// ContributionCommunityResponse represents the community contribution response
type ContributionCommunityResponse struct {
	AmountFillings int64   `json:"amountFillings"`
	AmountWater    int64   `json:"amountWater"`
	SavedMoney     float64 `json:"savedMoney"`
	SavedTrash     float64 `json:"savedTrash"`
	AmountUser     int64   `json:"amountUser"`
}

// @Summary Get community contribution
// @Description Get the total water amount and savings for the community
// @Tags contribution
// @Accept  json
// @Produce  json
// @Success 200 {object} ContributionCommunityResponse
// @Router /contribution/community [get]
func getContributionCommunity(c *gin.Context) {
	var totalVolume int64
	var totalFillings int64
	var totalUsers int64

	db.Model(&database.WaterTransaction{}).Count(&totalFillings)
	db.Model(&database.WaterTransaction{}).Select("sum(volume)").Row().Scan(&totalVolume)
	db.Model(&database.User{}).Count(&totalUsers)

	savedMoney, savedTrash := calculateSavings(int(totalVolume))

	response := ContributionCommunityResponse{
		AmountFillings: totalFillings,
		AmountWater:    totalVolume,
		SavedMoney:     savedMoney,
		SavedTrash:     savedTrash,
		AmountUser:     totalUsers,
	}

	respondWithJSON(c, http.StatusOK, response)
}

// ContributionKLResponse represents the contribution by station type response
type ContributionKLResponse struct {
	AmountRefillStationSmart  int64 `json:"amountRefillStationSmart"`
	AmountRefillStationManual int64 `json:"amountRefillStationManual"`
}

// @Summary Get contribution by station type
// @Description Get the number of smart and manual refill stations
// @Tags contribution
// @Accept  json
// @Produce  json
// @Success 200 {object} ContributionKLResponse
// @Router /contribution/kl [get]
func getContributionKL(c *gin.Context) {
	var smartStations int64
	var manualStations int64

	db.Model(&database.RefillStation{}).Where("type = ?", "Smart").Count(&smartStations)
	db.Model(&database.RefillStation{}).Where("type = ?", "Manual").Count(&manualStations)

	response := ContributionKLResponse{
		AmountRefillStationSmart:  smartStations,
		AmountRefillStationManual: manualStations,
	}

	respondWithJSON(c, http.StatusOK, response)
}

// RefillStationMarkerResponse represents a refill station marker in the response
type RefillStationMarkerResponse struct {
	ID        uint    `json:"id"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Status    bool    `json:"status"`
}

// @Summary Get all refill station markers
// @Description Get all refill station markers with specific attributes
// @Tags refill_stations
// @Accept  json
// @Produce  json
// @Success 200 {array} RefillStationMarkerResponse
// @Router /refill_stations/markers [get]
func getAllRefillstationMarker(c *gin.Context) {
	var stations []database.RefillStation
	result := db.Select("id, longitude, latitude, active").Find(&stations)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	var markers []RefillStationMarkerResponse
	for _, station := range stations {
		markers = append(markers, RefillStationMarkerResponse{
			ID:        station.ID,
			Longitude: station.Longitude,
			Latitude:  station.Latitude,
			Status:    station.Active,
		})
	}

	respondWithJSON(c, http.StatusOK, markers)
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server for a water station.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	r := gin.Default()

	r.GET("/users", getUsers)
	r.POST("/users", createUser)
	r.PUT("/users", updateUser)
	r.DELETE("/users", deleteUser)

	r.GET("/bottles", getBottles)
	r.POST("/bottles", createBottle)
	r.PUT("/bottles", updateBottle)
	r.DELETE("/bottles", deleteBottle)

	r.GET("/consumer_tests", getConsumerTests)
	r.POST("/consumer_tests", createConsumerTest)
	r.PUT("/consumer_tests", updateConsumerTest)
	r.DELETE("/consumer_tests", deleteConsumerTest)

	r.GET("/consumer_test_questions", getConsumerTestQuestions)
	r.POST("/consumer_test_questions", createConsumerTestQuestion)
	r.PUT("/consumer_test_questions", updateConsumerTestQuestion)
	r.DELETE("/consumer_test_questions", deleteConsumerTestQuestion)

	r.GET("/consumer_test_answers", getConsumerTestAnswers)
	r.POST("/consumer_test_answers", createConsumerTestAnswer)
	r.PUT("/consumer_test_answers", updateConsumerTestAnswer)
	r.DELETE("/consumer_test_answers", deleteConsumerTestAnswer)

	r.GET("/refill_stations", getRefillStations)
	r.POST("/refill_stations", createRefillStation)
	r.PUT("/refill_stations", updateRefillStation)
	r.DELETE("/refill_stations", deleteRefillStation)

	r.GET("/refill_station_reviews", getRefillStationReviews)
	r.POST("/refill_station_reviews", createRefillStationReview)
	r.PUT("/refill_station_reviews", updateRefillStationReview)
	r.DELETE("/refill_station_reviews", deleteRefillStationReview)

	r.GET("/refill_station_problems", getRefillStationProblems)
	r.POST("/refill_station_problems", createRefillStationProblem)
	r.PUT("/refill_station_problems", updateRefillStationProblem)
	r.DELETE("/refill_station_problems", deleteRefillStationProblem)

	r.GET("/water_transactions", getWaterTransactions)
	r.POST("/water_transactions", createWaterTransaction)
	r.PUT("/water_transactions", updateWaterTransaction)
	r.DELETE("/water_transactions", deleteWaterTransaction)

	r.GET("/likes", getLikes)
	r.POST("/likes", createLike)
	r.PUT("/likes", updateLike)
	r.DELETE("/likes", deleteLike)

	r.GET("/refill_stations/:id", getRefillStationById)
	r.GET("/refill_station_review/average", getRefillStationReview)
	r.GET("/refillstation_like", getRefillstationLike)
	r.POST("/refillstation_like", postRefillstationLike)
	r.GET("/contribution/user", getContributionByUser)
	r.GET("/contribution/community", getContributionCommunity)
	r.GET("/contribution/kl", getContributionKL)
	r.GET("/refill_stations/markers", getAllRefillstationMarker)

	// Swagger UI endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Server läuft und hört auf Port 8080...")
	err := r.Run(":8080")
	if err != nil {
		log.Fatalf("Server konnte nicht gestartet werden: %v", err)
	}
}
