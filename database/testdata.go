package database

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
)

type BottleJSON struct {
	UserID     uint
	NFCID      string
	FillVolume int
	WaterType  string
	Title      string
	ImagePath  string
}

type UsersJSON struct {
	FirstName string
	LastName  string
}

type RefillStationJSON struct {
	Name              string
	Description       string
	Latitude          float64
	Longitude         float64
	Address           string
	WaterSource       string
	OpeningTimes      string
	Type              string
	OfferedWaterTypes string
	ImagePath         string
}

type RefillStationProblemJSON struct {
	StationID                 uint
	Title                     string
	Description               string
	Status                    string
	RefillStationProblemImage string
}

type RefillStationReviewJSON struct {
	StationID     uint
	UserID        uint
	Cleanness     int
	Accessibility int
	WaterQuality  int
}

type WaterTransactionJSON struct {
	StationID uint
	BottleID  *uint
	UserID    *uint
	Volume    int
	WaterType string
	Timestamp time.Time
	Guest     bool
}

func CreateTestData(db *gorm.DB) *gorm.DB {
	log.Print("Test data creation started")

	db = CreateUsers(db)
	db = CreateBottles(db)
	db = CreateRefillStations(db)
	db = CreateRefillStationReviews(db)
	db = CreateRefillStationProblems(db)
	db = CreateWaterTransactions(db)
	db = CreateLikes(db)

	log.Print("Test data creation finished")

	return db
}

func CreateUsers(db *gorm.DB) *gorm.DB {
	// Read the JSON file
	file, err := os.Open("./testdata/users.json")
	if err != nil {
		log.Fatalf("failed to open JSON file: %v", err)
	}
	defer file.Close()

	// Read the file content
	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON data into intermediate struct
	var usersJson []UsersJSON
	if err := json.Unmarshal(bytes, &usersJson); err != nil {
		log.Fatalf("failed to unmarshal JSON data: %v", err)
	}

	// Convert intermediate data to actual User structs
	var users []User
	for _, userData := range usersJson {
		user := User{
			FirstName: userData.FirstName,
			LastName:  userData.LastName,
		}
		users = append(users, user)
	}

	// Create users in the database
	if err := db.Create(&users).Error; err != nil {
		log.Fatalf("failed to create users: %v", err)
	}

	log.Print("Created users successfully")

	return db
}

func CreateBottles(db *gorm.DB) *gorm.DB {
	// Read the JSON file
	file, err := os.Open("./testdata/bottles.json")
	if err != nil {
		log.Fatalf("failed to open JSON file: %v", err)
	}
	defer file.Close()

	// Read the file content
	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON data into a slice of BottleJSON
	var bottlesJSON []BottleJSON
	if err := json.Unmarshal(bytes, &bottlesJSON); err != nil {
		log.Fatalf("failed to unmarshal JSON data: %v", err)
	}

	// Convert image paths to base64 strings and create Bottle slice
	var bottles []Bottle
	for _, bottleJSON := range bottlesJSON {
		bottleImage := ImageToBase64(bottleJSON.ImagePath)
		bottles = append(bottles, Bottle{
			UserID:      bottleJSON.UserID,
			NFCID:       bottleJSON.NFCID,
			FillVolume:  bottleJSON.FillVolume,
			WaterType:   bottleJSON.WaterType,
			Title:       bottleJSON.Title,
			BottleImage: &bottleImage,
		})
	}

	// Create bottles in the database
	if err := db.Create(&bottles).Error; err != nil {
		log.Fatalf("failed to create bottles: %v", err)
	}

	log.Print("Created bottles successfully")

	return db
}

func CreateRefillStations(db *gorm.DB) *gorm.DB {
	// Read the JSON file
	file, err := os.Open("./testdata/refill_stations.json")
	if err != nil {
		log.Fatalf("failed to open JSON file: %v", err)
	}
	defer file.Close()

	// Read the file content
	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON data into a slice of RefillStationJSON
	var refillStationsJSON []RefillStationJSON
	if err := json.Unmarshal(bytes, &refillStationsJSON); err != nil {
		log.Fatalf("failed to unmarshal JSON data: %v", err)
	}

	// Convert image paths to base64 strings and create RefillStation slice
	var refillStations []RefillStation
	for _, stationJSON := range refillStationsJSON {
		refillStationImage := ImageToBase64(stationJSON.ImagePath)
		refillStations = append(refillStations, RefillStation{
			Name:               stationJSON.Name,
			Description:        stationJSON.Description,
			Latitude:           stationJSON.Latitude,
			Longitude:          stationJSON.Longitude,
			Address:            stationJSON.Address,
			WaterSource:        stationJSON.WaterSource,
			OpeningTimes:       stationJSON.OpeningTimes,
			Type:               stationJSON.Type,
			OfferedWaterTypes:  stationJSON.OfferedWaterTypes,
			RefillStationImage: &refillStationImage,
		})
	}

	// Create refill stations in the database
	if err := db.Create(&refillStations).Error; err != nil {
		log.Fatalf("failed to create refill stations: %v", err)
	}

	log.Print("Created refill stations successfully")

	return db
}

func CreateRefillStationReviews(db *gorm.DB) *gorm.DB {
	// Read the JSON file
	file, err := os.Open("./testdata/refill_station_reviews.json")
	if err != nil {
		log.Fatalf("failed to open JSON file: %v", err)
	}
	defer file.Close()

	// Read the file content
	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON data into a slice of the temporary struct
	var reviewsJSON []RefillStationReviewJSON
	if err := json.Unmarshal(bytes, &reviewsJSON); err != nil {
		log.Fatalf("failed to unmarshal JSON data: %v", err)
	}

	// Convert the temporary struct to the actual model
	var reviews []RefillStationReview
	for _, problemJSON := range reviewsJSON {
		problem := RefillStationReview{
			StationID:     problemJSON.StationID,
			UserID:        problemJSON.UserID,
			Cleanness:     problemJSON.Cleanness,
			WaterQuality:  problemJSON.WaterQuality,
			Accessibility: problemJSON.Accessibility,
		}
		reviews = append(reviews, problem)
	}

	// Create refill station problems in the database
	if err := db.Create(&reviews).Error; err != nil {
		log.Fatalf("failed to create refill station reviews: %v", err)
	}

	log.Print("Created refill station reviews successfully")

	return db
}

func CreateRefillStationProblems(db *gorm.DB) *gorm.DB {
	// Read the JSON file
	file, err := os.Open("./testdata/refill_station_problems.json")
	if err != nil {
		log.Fatalf("failed to open JSON file: %v", err)
	}
	defer file.Close()

	// Read the file content
	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON data into a slice of the temporary struct
	var problemsJSON []RefillStationProblemJSON
	if err := json.Unmarshal(bytes, &problemsJSON); err != nil {
		log.Fatalf("failed to unmarshal JSON data: %v", err)
	}

	// Convert the temporary struct to the actual model
	var problems []RefillStationProblem
	for _, problemJSON := range problemsJSON {
		imageBase64 := ImageToBase64(problemJSON.RefillStationProblemImage)
		problem := RefillStationProblem{
			StationID:                 problemJSON.StationID,
			Title:                     problemJSON.Title,
			Description:               problemJSON.Description,
			Status:                    problemJSON.Status,
			RefillStationProblemImage: &imageBase64,
		}
		problems = append(problems, problem)
	}

	// Create refill station problems in the database
	if err := db.Create(&problems).Error; err != nil {
		log.Fatalf("failed to create refill station problems: %v", err)
	}

	log.Print("Created refill station problems successfully")

	return db
}

func CreateWaterTransactions(db *gorm.DB) *gorm.DB {
	// Read the JSON file
	file, err := os.Open("./testdata/water_transactions.json")
	if err != nil {
		log.Fatalf("failed to open JSON file: %v", err)
	}
	defer file.Close()

	// Read the file content
	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON data into a slice of WaterTransaction
	var json_transactions []WaterTransactionJSON
	if err := json.Unmarshal(bytes, &json_transactions); err != nil {
		log.Fatalf("failed to unmarshal JSON data: %v", err)
	}

	// Convert the temporary struct to the actual model
	var transactions []WaterTransaction
	for _, transJSON := range json_transactions {
		transaction := WaterTransaction{
			StationID: transJSON.StationID,
			BottleID:  transJSON.BottleID,
			UserID:    transJSON.UserID,
			Volume:    transJSON.Volume,
			WaterType: transJSON.WaterType,
			Timestamp: transJSON.Timestamp,
			Guest:     transJSON.Guest,
		}
		transactions = append(transactions, transaction)
	}

	// Create water transactions in the database
	if err := db.Create(&transactions).Error; err != nil {
		log.Fatalf("failed to create water transactions: %v", err)
	}

	log.Print("Created water transactions successfully")

	return db
}

func CreateLikes(db *gorm.DB) *gorm.DB {
	// Read the JSON file
	file, err := os.Open("./testdata/likes.json")
	if err != nil {
		log.Fatalf("failed to open JSON file: %v", err)
	}
	defer file.Close()

	// Read the file content
	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON data into a slice of Like structs
	var likes []Like
	if err := json.Unmarshal(bytes, &likes); err != nil {
		log.Fatalf("failed to unmarshal JSON data: %v", err)
	}

	// Create likes in the database
	if err := db.Create(&likes).Error; err != nil {
		log.Fatalf("failed to create likes: %v", err)
	}

	log.Print("Created likes successfully")

	return db
}
