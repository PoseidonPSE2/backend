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
	db.AutoMigrate(&database.User{}, &database.NFCChip{}, &database.ConsumerTest{}, &database.ConsumerTestQuestion{},
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

// @Summary Show all users
// @Description Get all users
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} database.User
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

// @Summary Create a user
// @Description Create a new user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body database.User true "User"
// @Success 201 {object} database.User
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

// @Summary Update a user
// @Description Update an existing user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body database.User true "User"
// @Success 200 {object} database.User
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

// @Summary Show all NFC chips
// @Description Get all NFC chips
// @Tags nfchips
// @Accept  json
// @Produce  json
// @Success 200 {array} database.NFCChip
// @Router /nfchips [get]
func getNFCChips(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		var chips []database.NFCChip
		result := db.Find(&chips)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		respondWithJSON(c, http.StatusOK, chips)
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var chip database.NFCChip
		result := db.First(&chip, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
			return
		}
		respondWithJSON(c, http.StatusOK, chip)
	}
}

// @Summary Create an NFC chip
// @Description Create a new NFC chip
// @Tags nfchips
// @Accept  json
// @Produce  json
// @Param chip body database.NFCChip true "NFC Chip"
// @Success 201 {object} database.NFCChip
// @Router /nfchips [post]
func createNFCChip(c *gin.Context) {
	var chip database.NFCChip
	if err := c.ShouldBindJSON(&chip); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := db.Create(&chip)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	respondWithJSON(c, http.StatusCreated, chip)
}

// @Summary Update an NFC chip
// @Description Update an existing NFC chip
// @Tags nfchips
// @Accept  json
// @Produce  json
// @Param chip body database.NFCChip true "NFC Chip"
// @Success 200 {object} database.NFCChip
// @Router /nfchips [put]
func updateNFCChip(c *gin.Context) {
	var chip database.NFCChip
	if err := c.ShouldBindJSON(&chip); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := db.Save(&chip)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	respondWithJSON(c, http.StatusOK, chip)
}

// @Summary Delete an NFC chip
// @Description Delete an existing NFC chip
// @Tags nfchips
// @Accept  json
// @Produce  json
// @Param id query int true "NFC Chip ID"
// @Success 204
// @Router /nfchips [delete]
func deleteNFCChip(c *gin.Context) {
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
	result := db.Delete(&database.NFCChip{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Show all consumer tests
// @Description Get all consumer tests
// @Tags consumer_tests
// @Accept  json
// @Produce  json
// @Success 200 {array} database.ConsumerTest
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

// @Summary Create a consumer test
// @Description Create a new consumer test
// @Tags consumer_tests
// @Accept  json
// @Produce  json
// @Param test body database.ConsumerTest true "Consumer Test"
// @Success 201 {object} database.ConsumerTest
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

// @Summary Update a consumer test
// @Description Update an existing consumer test
// @Tags consumer_tests
// @Accept  json
// @Produce  json
// @Param test body database.ConsumerTest true "Consumer Test"
// @Success 200 {object} database.ConsumerTest
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

// @Summary Show all consumer test questions
// @Description Get all consumer test questions
// @Tags consumer_test_questions
// @Accept  json
// @Produce  json
// @Success 200 {array} database.ConsumerTestQuestion
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

// @Summary Create a consumer test question
// @Description Create a new consumer test question
// @Tags consumer_test_questions
// @Accept  json
// @Produce  json
// @Param question body database.ConsumerTestQuestion true "Consumer Test Question"
// @Success 201 {object} database.ConsumerTestQuestion
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

// @Summary Update a consumer test question
// @Description Update an existing consumer test question
// @Tags consumer_test_questions
// @Accept  json
// @Produce  json
// @Param question body database.ConsumerTestQuestion true "Consumer Test Question"
// @Success 200 {object} database.ConsumerTestQuestion
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

// @Summary Show all consumer test answers
// @Description Get all consumer test answers
// @Tags consumer_test_answers
// @Accept  json
// @Produce  json
// @Success 200 {array} database.ConsumerTestAnswer
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

// @Summary Create a consumer test answer
// @Description Create a new consumer test answer
// @Tags consumer_test_answers
// @Accept  json
// @Produce  json
// @Param answer body database.ConsumerTestAnswer true "Consumer Test Answer"
// @Success 201 {object} database.ConsumerTestAnswer
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

// @Summary Update a consumer test answer
// @Description Update an existing consumer test answer
// @Tags consumer_test_answers
// @Accept  json
// @Produce  json
// @Param answer body database.ConsumerTestAnswer true "Consumer Test Answer"
// @Success 200 {object} database.ConsumerTestAnswer
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

// @Summary Show all refill stations
// @Description Get all refill stations
// @Tags refill_stations
// @Accept  json
// @Produce  json
// @Success 200 {array} database.RefillStation
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

// @Summary Create a refill station
// @Description Create a new refill station
// @Tags refill_stations
// @Accept  json
// @Produce  json
// @Param station body database.RefillStation true "Refill Station"
// @Success 201 {object} database.RefillStation
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

// @Summary Update a refill station
// @Description Update an existing refill station
// @Tags refill_stations
// @Accept  json
// @Produce  json
// @Param station body database.RefillStation true "Refill Station"
// @Success 200 {object} database.RefillStation
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

// @Summary Show all refill station reviews
// @Description Get all refill station reviews
// @Tags refill_station_reviews
// @Accept  json
// @Produce  json
// @Success 200 {array} database.RefillStationReview
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

// @Summary Create a refill station review
// @Description Create a new refill station review
// @Tags refill_station_reviews
// @Accept  json
// @Produce  json
// @Param review body database.RefillStationReview true "Refill Station Review"
// @Success 201 {object} database.RefillStationReview
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

// @Summary Update a refill station review
// @Description Update an existing refill station review
// @Tags refill_station_reviews
// @Accept  json
// @Produce  json
// @Param review body database.RefillStationReview true "Refill Station Review"
// @Success 200 {object} database.RefillStationReview
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

// @Summary Show all refill station problems
// @Description Get all refill station problems
// @Tags refill_station_problems
// @Accept  json
// @Produce  json
// @Success 200 {array} database.RefillStationProblem
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

// @Summary Create a refill station problem
// @Description Create a new refill station problem
// @Tags refill_station_problems
// @Accept  json
// @Produce  json
// @Param problem body database.RefillStationProblem true "Refill Station Problem"
// @Success 201 {object} database.RefillStationProblem
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

// @Summary Update a refill station problem
// @Description Update an existing refill station problem
// @Tags refill_station_problems
// @Accept  json
// @Produce  json
// @Param problem body database.RefillStationProblem true "Refill Station Problem"
// @Success 200 {object} database.RefillStationProblem
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

// @Summary Show all water transactions
// @Description Get all water transactions
// @Tags water_transactions
// @Accept  json
// @Produce  json
// @Success 200 {array} database.WaterTransaction
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

// @Summary Create a water transaction
// @Description Create a new water transaction
// @Tags water_transactions
// @Accept  json
// @Produce  json
// @Param transaction body database.WaterTransaction true "Water Transaction"
// @Success 201 {object} database.WaterTransaction
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

// @Summary Update a water transaction
// @Description Update an existing water transaction
// @Tags water_transactions
// @Accept  json
// @Produce  json
// @Param transaction body database.WaterTransaction true "Water Transaction"
// @Success 200 {object} database.WaterTransaction
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

// @Summary Show all likes
// @Description Get all likes
// @Tags likes
// @Accept  json
// @Produce  json
// @Success 200 {array} database.Like
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

// @Summary Create a like
// @Description Create a new like
// @Tags likes
// @Accept  json
// @Produce  json
// @Param like body database.Like true "Like"
// @Success 201 {object} database.Like
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

// @Summary Update a like
// @Description Update an existing like
// @Tags likes
// @Accept  json
// @Produce  json
// @Param like body database.Like true "Like"
// @Success 200 {object} database.Like
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

//extrawünsche

// @Summary Get a refill station by ID
// @Description Get a refill station by its ID
// @Tags refill_stations
// @Accept  json
// @Produce  json
// @Param id path int true "Refill Station ID"
// @Success 200 {object} database.RefillStation
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

	r.GET("/nfchips", getNFCChips)
	r.POST("/nfchips", createNFCChip)
	r.PUT("/nfchips", updateNFCChip)
	r.DELETE("/nfchips", deleteNFCChip)

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

	// Swagger UI endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Server läuft und hört auf Port 8080...")
	err := r.Run(":8080")
	if err != nil {
		log.Fatalf("Server konnte nicht gestartet werden: %v", err)
	}
}
