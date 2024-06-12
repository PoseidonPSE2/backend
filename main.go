package main

import (
	"fmt"
	"log"

	"gorm.io/gorm"

	"github.com/PoseidonPSE2/code_backend/api"
	"github.com/PoseidonPSE2/code_backend/database"

	_ "github.com/PoseidonPSE2/code_backend/docs" // swagger docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Database connection
var db *gorm.DB

// Flog for database
var shouldRecreateDatabase = false
var shouldImportTestData = false

// Database configuration variables
var (
	dbHost     = "poseidon-database.fly.dev"
	dbPort     = "5432"
	dbUser     = "postgres"
	dbPassword = "rnJpE83UKr1MyF8"
	dbName     = "poseidon_db"
)

func init() {
	var err error
	log.Print("Starting application")

	// Construct the DSN for the administrative connection
	adminDsn := fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=disable dbname=postgres", dbHost, dbUser, dbPassword, dbPort)

	// Construct the DSN for the target database connection
	targetDsn := fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=disable dbname=%s", dbHost, dbUser, dbPassword, dbPort, dbName)

	if shouldRecreateDatabase {
		// Recreate the target database
		if err = database.RecreateDatabase(adminDsn, dbName, db); err != nil {
			log.Fatalf("Failed to recreate database: %v", err)
		}
	}

	// Connect to the new database
	db, err = database.ConnectDatabase(targetDsn, db)
	if err != nil {
		log.Fatalf("failed to connect to new database: %v", err)
	}

	log.Print("Schema migration starting")

	// Migrate the schema
	db.AutoMigrate(&database.User{}, &database.Bottle{}, &database.ConsumerTest{}, &database.ConsumerTestQuestion{},
		&database.ConsumerTestAnswer{}, &database.RefillStation{}, &database.RefillStationReview{},
		&database.RefillStationProblem{}, &database.WaterTransaction{}, &database.Like{})

	log.Print("Schema migration done")

	if shouldImportTestData {
		db = database.CreateTestData(db)
	}
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

// @host poseidon-backend.fly.dev
// @BasePath /
func main() {
	api.SetDB(db)
	r := gin.Default()

	r.GET("/users", api.GetUsers)
	r.POST("/users", api.CreateUser)
	r.PUT("/users", api.UpdateUser)
	r.DELETE("/users", api.DeleteUser)

	r.GET("/bottles", api.GetBottles)
	r.GET("/bottles/:id", api.GetBottleById)
	r.GET("/bottles/users/:userId", api.GetBottlesByUserID)
	r.GET("/bottles/preferences/:nfcId", api.GetBottlePreferencesByNFCId)
	r.POST("/bottles", api.CreateBottle)
	r.PUT("/bottles", api.UpdateBottle)
	r.DELETE("/bottles/:id", api.DeleteBottle)

	r.GET("/consumer_tests", api.GetConsumerTests)
	r.POST("/consumer_tests", api.CreateConsumerTest)
	r.PUT("/consumer_tests", api.UpdateConsumerTest)
	r.DELETE("/consumer_tests", api.DeleteConsumerTest)

	r.GET("/consumer_test_questions", api.GetConsumerTestQuestions)
	r.POST("/consumer_test_questions", api.CreateConsumerTestQuestion)
	r.PUT("/consumer_test_questions", api.UpdateConsumerTestQuestion)
	r.DELETE("/consumer_test_questions", api.DeleteConsumerTestQuestion)

	r.GET("/consumer_test_answers", api.GetConsumerTestAnswers)
	r.POST("/consumer_test_answers", api.CreateConsumerTestAnswer)
	r.PUT("/consumer_test_answers", api.UpdateConsumerTestAnswer)
	r.DELETE("/consumer_test_answers", api.DeleteConsumerTestAnswer)

	r.GET("/refill_stations", api.GetRefillStations)
	r.GET("/refill_stations/markers", api.GetAllRefillstationMarker)
	r.GET("/refill_stations/:id", api.GetRefillStationById)
	r.GET("/refill_stations/:id/reviews", api.GetRefillStationReviewsAverageByID)
	r.POST("/refill_stations", api.CreateRefillStation)
	r.PUT("/refill_stations", api.UpdateRefillStation)
	r.DELETE("/refill_stations", api.DeleteRefillStation)

	r.GET("/refill_station_reviews", api.GetRefillStationReviews)
	r.POST("/refill_station_reviews", api.CreateRefillStationReview)
	r.PUT("/refill_station_reviews", api.UpdateRefillStationReview)
	r.DELETE("/refill_station_reviews", api.DeleteRefillStationReview)

	r.GET("/refill_station_problems", api.GetRefillStationProblems)
	r.POST("/refill_station_problems", api.CreateRefillStationProblem)
	r.PUT("/refill_station_problems", api.UpdateRefillStationProblem)
	r.DELETE("/refill_station_problems", api.DeleteRefillStationProblem)

	r.GET("/water_transactions", api.GetWaterTransactions)
	r.POST("/water_transactions", api.CreateWaterTransaction)
	r.PUT("/water_transactions", api.UpdateWaterTransaction)
	r.DELETE("/water_transactions", api.DeleteWaterTransaction)

	r.GET("/likes", api.GetLikes)
	r.GET("/likes/:refillstationId/:usedId", api.GetLikeByUserIdAndStationID)
	r.POST("/likes", api.CreateLike)
	r.PUT("/likes", api.UpdateLike)
	r.DELETE("/likes", api.DeleteLike)

	r.GET("/contribution/user/:id", api.GetContributionByUser)
	r.GET("/contribution/community", api.GetContributionCommunity)
	r.GET("/contribution/kl", api.GetContributionKL)

	// Swagger UI endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Server running and serving at Port 8080...")
	err := r.Run(":8080")
	if err != nil {
		log.Fatalf("Server could not start: %v", err)
	}
}
