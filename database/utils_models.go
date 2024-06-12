package database

import (
	"database/sql"
	"database/sql/driver"
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
