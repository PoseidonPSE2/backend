package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/postgres/pgxv5"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/PoseidonPSE2/code_backend/database"
)

// Initialize the database connection
var db *gorm.DB

func init() {
	log.Print("Starting application")

	cleanup, err := pgxv5.RegisterDriver(
		"cloudsql-postgres",
		cloudsqlconn.WithLazyRefresh(),
		cloudsqlconn.WithIAMAuthN(),
	)
	if err != nil {
		panic(err)
	}

	log.Print("Connecting to database")

	user := "developer"
	password := "pw"
	dbHost := "unique-machine-422214-b0:europe-west3:poseidon-database"
	//dbHost := "35.246.250.79"
	databaseName := "poseidon"

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, user, databaseName, password)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "cloudsql-postgres",
		DSN:        dsn,
	}))
	if err != nil {
		panic(err)
	} else {
		// get the underlying *sql.DB type to verify the connection
		sdb, err := db.DB()
		if err != nil {
			panic(err)
		}

		var t time.Time
		if err := sdb.QueryRow("select now()").Scan(&t); err != nil {
			panic(err)
		}

		fmt.Println(t)

		log.Print("succesfuly connected to database")

		log.Print("Schema migration starting")
		// Migrate the schema
		db.AutoMigrate(&database.User{}, &database.NFCChip{}, &database.ConsumerTest{}, &database.ConsumerTestQuestion{},
			&database.ConsumerTestAnswer{}, &database.RefillStation{}, &database.RefillStationReview{},
			&database.RefillStationProblem{}, &database.WaterTransaction{}, &database.Like{})
		log.Print("Schema migration done")
	}

	// cleanup will stop the driver from retrieving ephemeral certificates
	// Don't call cleanup until you're done with your database connections
	defer cleanup()
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// Handlers for User model
func handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsers(w, r)
	case http.MethodPost:
		createUser(w, r)
	case http.MethodPut:
		updateUser(w, r)
	case http.MethodDelete:
		deleteUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		var users []database.User
		result := db.Find(&users)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		respondWithJSON(w, http.StatusOK, users)
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		var user database.User
		result := db.First(&user, id)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusNotFound)
			return
		}
		respondWithJSON(w, http.StatusOK, user)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user database.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := db.Create(&user)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	var user database.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := db.Save(&user)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	result := db.Delete(&database.User{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Handlers for NFCChip model
func handleNFCChips(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getNFCChips(w, r)
	case http.MethodPost:
		createNFCChip(w, r)
	case http.MethodPut:
		updateNFCChip(w, r)
	case http.MethodDelete:
		deleteNFCChip(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getNFCChips(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		var chips []database.NFCChip
		result := db.Find(&chips)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		respondWithJSON(w, http.StatusOK, chips)
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		var chip database.NFCChip
		result := db.First(&chip, id)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusNotFound)
			return
		}
		respondWithJSON(w, http.StatusOK, chip)
	}
}

func createNFCChip(w http.ResponseWriter, r *http.Request) {
	var chip database.NFCChip
	if err := json.NewDecoder(r.Body).Decode(&chip); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := db.Create(&chip)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusCreated, chip)
}

func updateNFCChip(w http.ResponseWriter, r *http.Request) {
	var chip database.NFCChip
	if err := json.NewDecoder(r.Body).Decode(&chip); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := db.Save(&chip)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, chip)
}

func deleteNFCChip(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	result := db.Delete(&database.NFCChip{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Implement handlers for ConsumerTest model
func handleConsumerTests(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getConsumerTests(w, r)
	case http.MethodPost:
		createConsumerTest(w, r)
	case http.MethodPut:
		updateConsumerTest(w, r)
	case http.MethodDelete:
		deleteConsumerTest(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getConsumerTests(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		var tests []database.ConsumerTest
		result := db.Find(&tests)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		respondWithJSON(w, http.StatusOK, tests)
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		var test database.ConsumerTest
		result := db.First(&test, id)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusNotFound)
			return
		}
		respondWithJSON(w, http.StatusOK, test)
	}
}

func createConsumerTest(w http.ResponseWriter, r *http.Request) {
	var test database.ConsumerTest
	if err := json.NewDecoder(r.Body).Decode(&test); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := db.Create(&test)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusCreated, test)
}

func updateConsumerTest(w http.ResponseWriter, r *http.Request) {
	var test database.ConsumerTest
	if err := json.NewDecoder(r.Body).Decode(&test); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := db.Save(&test)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, test)
}

func deleteConsumerTest(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	result := db.Delete(&database.ConsumerTest{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Implement handlers for ConsumerTestQuestion model
func handleConsumerTestQuestions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getConsumerTestQuestions(w, r)
	case http.MethodPost:
		createConsumerTestQuestion(w, r)
	case http.MethodPut:
		updateConsumerTestQuestion(w, r)
	case http.MethodDelete:
		deleteConsumerTestQuestion(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getConsumerTestQuestions(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		var questions []database.ConsumerTestQuestion
		result := db.Find(&questions)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		respondWithJSON(w, http.StatusOK, questions)
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		var question database.ConsumerTestQuestion
		result := db.First(&question, id)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusNotFound)
			return
		}
		respondWithJSON(w, http.StatusOK, question)
	}
}

func createConsumerTestQuestion(w http.ResponseWriter, r *http.Request) {
	var question database.ConsumerTestQuestion
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := db.Create(&question)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusCreated, question)
}

func updateConsumerTestQuestion(w http.ResponseWriter, r *http.Request) {
	var question database.ConsumerTestQuestion
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := db.Save(&question)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, question)
}

func deleteConsumerTestQuestion(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	result := db.Delete(&database.ConsumerTestQuestion{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Implement handlers for ConsumerTestAnswer model
func handleConsumerTestAnswers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getConsumerTestAnswers(w, r)
	case http.MethodPost:
		createConsumerTestAnswer(w, r)
	case http.MethodPut:
		updateConsumerTestAnswer(w, r)
	case http.MethodDelete:
		deleteConsumerTestAnswer(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getConsumerTestAnswers(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		var answers []database.ConsumerTestAnswer
		result := db.Find(&answers)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		respondWithJSON(w, http.StatusOK, answers)
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		var answer database.ConsumerTestAnswer
		result := db.First(&answer, id)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusNotFound)
			return
		}
		respondWithJSON(w, http.StatusOK, answer)
	}
}

func createConsumerTestAnswer(w http.ResponseWriter, r *http.Request) {
	var answer database.ConsumerTestAnswer
	if err := json.NewDecoder(r.Body).Decode(&answer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	answer.Timestamp = time.Now()
	result := db.Create(&answer)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusCreated, answer)
}

func updateConsumerTestAnswer(w http.ResponseWriter, r *http.Request) {
	var answer database.ConsumerTestAnswer
	if err := json.NewDecoder(r.Body).Decode(&answer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	answer.Timestamp = time.Now()
	result := db.Save(&answer)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, answer)
}

func deleteConsumerTestAnswer(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	result := db.Delete(&database.ConsumerTestAnswer{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Implement handlers for RefillStation model
func handleRefillStations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getRefillStations(w, r)
	case http.MethodPost:
		createRefillStation(w, r)
	case http.MethodPut:
		updateRefillStation(w, r)
	case http.MethodDelete:
		deleteRefillStation(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getRefillStations(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		var stations []database.RefillStation
		result := db.Find(&stations)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		respondWithJSON(w, http.StatusOK, stations)
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		var station database.RefillStation
		result := db.First(&station, id)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusNotFound)
			return
		}
		respondWithJSON(w, http.StatusOK, station)
	}
}

func createRefillStation(w http.ResponseWriter, r *http.Request) {
	var station database.RefillStation
	if err := json.NewDecoder(r.Body).Decode(&station); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := db.Create(&station)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusCreated, station)
}

func updateRefillStation(w http.ResponseWriter, r *http.Request) {
	var station database.RefillStation
	if err := json.NewDecoder(r.Body).Decode(&station); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := db.Save(&station)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, station)
}

func deleteRefillStation(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	result := db.Delete(&database.RefillStation{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Implement handlers for RefillStationReview model
func handleRefillStationReviews(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getRefillStationReviews(w, r)
	case http.MethodPost:
		createRefillStationReview(w, r)
	case http.MethodPut:
		updateRefillStationReview(w, r)
	case http.MethodDelete:
		deleteRefillStationReview(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getRefillStationReviews(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		var reviews []database.RefillStationReview
		result := db.Find(&reviews)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		respondWithJSON(w, http.StatusOK, reviews)
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		var review database.RefillStationReview
		result := db.First(&review, id)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusNotFound)
			return
		}
		respondWithJSON(w, http.StatusOK, review)
	}
}

func createRefillStationReview(w http.ResponseWriter, r *http.Request) {
	var review database.RefillStationReview
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	review.Timestamp = time.Now()
	result := db.Create(&review)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusCreated, review)
}

func updateRefillStationReview(w http.ResponseWriter, r *http.Request) {
	var review database.RefillStationReview
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	review.Timestamp = time.Now()
	result := db.Save(&review)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, review)
}

func deleteRefillStationReview(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	result := db.Delete(&database.RefillStationReview{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Implement handlers for RefillStationProblem model
func handleRefillStationProblems(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getRefillStationProblems(w, r)
	case http.MethodPost:
		createRefillStationProblem(w, r)
	case http.MethodPut:
		updateRefillStationProblem(w, r)
	case http.MethodDelete:
		deleteRefillStationProblem(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getRefillStationProblems(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		var problems []database.RefillStationProblem
		result := db.Find(&problems)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		respondWithJSON(w, http.StatusOK, problems)
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		var problem database.RefillStationProblem
		result := db.First(&problem, id)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusNotFound)
			return
		}
		respondWithJSON(w, http.StatusOK, problem)
	}
}

func createRefillStationProblem(w http.ResponseWriter, r *http.Request) {
	var problem database.RefillStationProblem
	if err := json.NewDecoder(r.Body).Decode(&problem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	problem.Timestamp = time.Now()
	result := db.Create(&problem)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusCreated, problem)
}

func updateRefillStationProblem(w http.ResponseWriter, r *http.Request) {
	var problem database.RefillStationProblem
	if err := json.NewDecoder(r.Body).Decode(&problem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	problem.Timestamp = time.Now()
	result := db.Save(&problem)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, problem)
}

func deleteRefillStationProblem(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	result := db.Delete(&database.RefillStationProblem{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Implement handlers for WaterTransaction model
func handleWaterTransactions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getWaterTransactions(w, r)
	case http.MethodPost:
		createWaterTransaction(w, r)
	case http.MethodPut:
		updateWaterTransaction(w, r)
	case http.MethodDelete:
		deleteWaterTransaction(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getWaterTransactions(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		var transactions []database.WaterTransaction
		result := db.Find(&transactions)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		respondWithJSON(w, http.StatusOK, transactions)
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		var transaction database.WaterTransaction
		result := db.First(&transaction, id)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusNotFound)
			return
		}
		respondWithJSON(w, http.StatusOK, transaction)
	}
}

func createWaterTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction database.WaterTransaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	transaction.Timestamp = time.Now()
	result := db.Create(&transaction)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusCreated, transaction)
}

func updateWaterTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction database.WaterTransaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	transaction.Timestamp = time.Now()
	result := db.Save(&transaction)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, transaction)
}

func deleteWaterTransaction(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	result := db.Delete(&database.WaterTransaction{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Implement handlers for Like model
func handleLikes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getLikes(w, r)
	case http.MethodPost:
		createLike(w, r)
	case http.MethodPut:
		updateLike(w, r)
	case http.MethodDelete:
		deleteLike(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getLikes(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		var likes []database.Like
		result := db.Find(&likes)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		respondWithJSON(w, http.StatusOK, likes)
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		var like database.Like
		result := db.First(&like, id)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusNotFound)
			return
		}
		respondWithJSON(w, http.StatusOK, like)
	}
}

func createLike(w http.ResponseWriter, r *http.Request) {
	var like database.Like
	if err := json.NewDecoder(r.Body).Decode(&like); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := db.Create(&like)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusCreated, like)
}

func updateLike(w http.ResponseWriter, r *http.Request) {
	var like database.Like
	if err := json.NewDecoder(r.Body).Decode(&like); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := db.Save(&like)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, like)
}

func deleteLike(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	result := db.Delete(&database.Like{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	http.HandleFunc("/users", handleUsers)
	http.HandleFunc("/nfchips", handleNFCChips)
	http.HandleFunc("/consumer_tests", handleConsumerTests)
	http.HandleFunc("/consumer_test_questions", handleConsumerTestQuestions)
	http.HandleFunc("/consumer_test_answers", handleConsumerTestAnswers)
	http.HandleFunc("/refill_stations", handleRefillStations)
	http.HandleFunc("/refill_station_reviews", handleRefillStationReviews)
	http.HandleFunc("/refill_station_problems", handleRefillStationProblems)
	http.HandleFunc("/water_transactions", handleWaterTransactions)
	http.HandleFunc("/likes", handleLikes)

	log.Println("Server läuft und hört auf Port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server konnte nicht gestartet werden: %v", err)
	}
}
