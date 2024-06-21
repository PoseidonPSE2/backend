package database

import (
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"log"
	"os"
)

// NullBool is a custom struct for handling sql.NullBool in Swagger
// @swagger:model
type NullBool struct {
	Bool  bool `json:"bool"`
	Valid bool `json:"valid"`
}

// Implementing the Scanner and Valuer interfaces for NullBool
func (nb *NullBool) Scan(value interface{}) error {
	sqlBool := sql.NullBool{}
	err := sqlBool.Scan(value)
	if err != nil {
		return err
	}
	nb.Bool = sqlBool.Bool
	nb.Valid = sqlBool.Valid
	return nil
}

func (nb NullBool) Value() (driver.Value, error) {
	if !nb.Valid {
		return nil, nil
	}
	return nb.Bool, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func isValidRating(rating int) bool {
	return rating >= 1 && rating <= 5
}

// Function to read an image file and return its base64 representation
func ImageToBase64(filePath string) string {
	// Read the entire image file
	imageData, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed to read image: %v", err)
		return ""
	}

	// Encode image data to base64 string
	base64Encoded := base64.StdEncoding.EncodeToString(imageData)

	return base64Encoded
}
