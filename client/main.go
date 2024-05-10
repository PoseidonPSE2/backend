package main

import (
	"database/crud"
	"database/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// PostgreSQL-Verbindungszeichenfolge
	dsn := "host=localhost user=postgres password=yourpassword dbname=yourdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	fmt.Println("Successfully connected to the database.")

	// Beispiel-User erstellen
	user := models.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     strPtr("john.doe@example.com"),
	}

	// User erstellen
	if err := crud.CreateUser(db, &user); err != nil {
		log.Fatalf("failed to create user: %v", err)
	}

	// Alle User finden
	users, err := crud.FindAllUsers(db)
	if err != nil {
		log.Fatalf("failed to find all users: %v", err)
	}
	fmt.Println("Users:", users)

	// Einen bestimmten User löschen
	if err := crud.DeleteUserByID(db, user.ID); err != nil {
		log.Fatalf("failed to delete user: %v", err)
	}
}

// Hilfsfunktion, um eine String-Pointer zu erstellen
func strPtr(s string) *string {
	return &s
}
