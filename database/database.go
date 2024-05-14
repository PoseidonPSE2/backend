package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// groß geschrieben ist public - deswegen Init! InitDatabase
func InitDatabase() {
	dsn := "host=localhost user=developer password=pw dbname=poseidon port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	fmt.Println("Successfully connected to the database.")

	// Automatische Migration der Modelle
	err = db.AutoMigrate(
		&User{},
		&ConsumerTest{},
		&ConsumerTestQuestion{},
		&ConsumerTestAnswer{},
		&NFCChip{},
		&RefillStation{},
		&RefillStationReview{},
		&RefillStationProblem{},
		&WaterTransaction{},
	)
	if err != nil {
		log.Fatalf("failed to migrate the database schema: %v", err)
	}
	fmt.Println("Database schema migrated successfully.")
}
