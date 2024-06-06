package database

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// User Model
// @swagger:model
type User struct {
	ID        uint                  `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName string                `gorm:"size:100;not null" json:"first_name"`
	LastName  string                `gorm:"size:100;not null" json:"last_name"`
	Email     *string               `gorm:"size:100;unique;default:null" json:"email"`
	Bottles   []Bottle              `gorm:"foreignKey:UserID" json:"-"`
	Answers   []ConsumerTestAnswer  `gorm:"foreignKey:UserID" json:"-"`
	Reviews   []RefillStationReview `gorm:"foreignKey:UserID" json:"-"`
	Likes     []Like                `gorm:"foreignKey:UserID" json:"-"`
}

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

// ConsumerTest Model
// @swagger:model
type ConsumerTest struct {
	ID        uint                   `gorm:"primaryKey" json:"id"`
	Questions []ConsumerTestQuestion `gorm:"foreignKey:TestID" json:"-"`
}

// ConsumerTestQuestion Model
// @swagger:model
type ConsumerTestQuestion struct {
	ID       uint                 `gorm:"primaryKey" json:"id"`
	Text     string               `gorm:"size:255;not null" json:"text"`
	MinValue *float32             `gorm:"default:null" json:"min_value,omitempty"`
	MaxValue *float32             `gorm:"default:null" json:"max_value,omitempty"`
	TestID   uint                 `gorm:"not null" json:"test_id"`
	Answers  []ConsumerTestAnswer `gorm:"foreignKey:QuestionID" json:"-"`
}

// ConsumerTestAnswer Model
// @swagger:model
type ConsumerTestAnswer struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null" json:"user_id"`
	QuestionID uint      `gorm:"not null" json:"question_id"`
	Answer     float32   `gorm:"not null" json:"answer"`
	Timestamp  time.Time `gorm:"autoCreateTime" json:"timestamp"`
}

// Bottle Model (previously NFCChip)
// @swagger:model
type Bottle struct {
	ID                uint               `gorm:"primaryKey" json:"id"`
	UserID            uint               `gorm:"not null" json:"user_id"`
	NFCID             string             `gorm:"size:32;unique;not null" json:"nfc_id"`
	FillVolume        int                `gorm:"not null" json:"fill_volume"`
	WaterType         string             `gorm:"size:16;not null" json:"water_type"`
	Title             string             `gorm:"size:16;not null" json:"title"`
	BottleImage       *string            `gorm:"type:TEXT;default:null" json:"bottle_image,omitempty"`
	Active            bool               `gorm:"default:true" json:"active"`
	WaterTransactions []WaterTransaction `gorm:"foreignKey:BottleID" json:"-"`
}

// RefillStation Model
// @swagger:model
type RefillStation struct {
	ID                 uint                   `gorm:"primaryKey" json:"id"`
	Name               string                 `gorm:"size:100;not null" json:"name"`
	Description        string                 `gorm:"size:255;not null" json:"description"`
	Latitude           float64                `gorm:"not null" json:"latitude"`
	Longitude          float64                `gorm:"not null" json:"longitude"`
	Address            string                 `gorm:"size:255;not null" json:"address"`
	WaterSource        string                 `gorm:"size:50;not null" json:"water_source"`
	OpeningTimes       string                 `gorm:"size:100;not null" json:"opening_times"`
	Active             NullBool               `gorm:"default:true" json:"active"`
	Type               string                 `gorm:"size:16;not null" json:"type"`
	OfferedWaterTypes  string                 `gorm:"size:32;not null" json:"offered_water_types"`
	RefillStationImage *string                `gorm:"type:TEXT;default:null" json:"refill_station_image,omitempty"`
	Reviews            []RefillStationReview  `gorm:"foreignKey:StationID" json:"-"`
	Problems           []RefillStationProblem `gorm:"foreignKey:StationID" json:"-"`
	WaterTransactions  []WaterTransaction     `gorm:"foreignKey:StationID" json:"-"`
	Likes              []Like                 `gorm:"foreignKey:StationID" json:"-"`
}

// RefillStationReview Model
// @swagger:model
type RefillStationReview struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	StationID     uint      `gorm:"not null" json:"station_id"`
	UserID        uint      `gorm:"not null" json:"user_id"`
	Cleanness     int       `gorm:"not null;check:cleanness >= 1 AND cleanness <= 5" json:"cleanness"`
	Accessibility int       `gorm:"not null;check:accessibility >= 1 AND accessibility <= 5" json:"accessibility"`
	WaterQuality  int       `gorm:"not null;check:water_quality >= 1 AND water_quality <= 5" json:"water_quality"`
	Timestamp     time.Time `gorm:"autoCreateTime" json:"timestamp"`
}

func (review *RefillStationReview) BeforeCreate(tx *gorm.DB) (err error) {
	if !isValidRating(review.Cleanness) ||
		!isValidRating(review.Accessibility) ||
		!isValidRating(review.WaterQuality) {
		return fmt.Errorf("rating should be between 1 and 5")
	}
	return nil
}

func (review *RefillStationReview) BeforeUpdate(tx *gorm.DB) (err error) {
	if !isValidRating(review.Cleanness) ||
		!isValidRating(review.Accessibility) ||
		!isValidRating(review.WaterQuality) {
		return fmt.Errorf("rating should be between 1 und 5")
	}
	return nil
}

// RefillStationProblem Model
// @swagger:model
type RefillStationProblem struct {
	ID                        uint      `gorm:"primaryKey" json:"id"`
	StationID                 uint      `gorm:"not null" json:"station_id"`
	Title                     string    `gorm:"size:100;not null" json:"title"`
	Description               string    `gorm:"size:255;not null" json:"description"`
	Status                    string    `gorm:"size:16;not null" json:"status"`
	RefillStationProblemImage *string   `gorm:"type:TEXT;default:null" json:"refill_station_problem_image,omitempty"`
	Timestamp                 time.Time `gorm:"autoCreateTime" json:"timestamp"`
}

// WaterTransaction Model
// @swagger:model
type WaterTransaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StationID uint      `gorm:"not null" json:"station_id"`
	BottleID  *uint     `gorm:"default:null" json:"bottle_id,omitempty"`
	UserID    *uint     `gorm:"default:null" json:"user_id,omitempty"`
	Volume    int       `gorm:"not null" json:"volume"`
	WaterType string    `gorm:"size:16;not null" json:"water_type"`
	Timestamp time.Time `gorm:"autoCreateTime" json:"timestamp"`
	Guest     bool      `gorm:"default:false" json:"guest"`
}

// Like Model
// @swagger:model
type Like struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	StationID uint `gorm:"not null" json:"station_id"`
	UserID    uint `gorm:"not null" json:"user_id"`
}

// Enumeration Constraints:
func (Bottle) TableName() string {
	return "bottle"
}

func (RefillStationProblem) TableName() string {
	return "refill_station_problem"
}

func (bottle *Bottle) BeforeCreate(tx *gorm.DB) (err error) {
	allowedWaterTypes := []string{"tap", "mineral"}
	waterType := strings.ToLower(bottle.WaterType)
	if !contains(allowedWaterTypes, waterType) {
		return fmt.Errorf("invalid water type: %s", bottle.WaterType)
	}
	bottle.WaterType = waterType
	return nil
}

func (station *RefillStation) BeforeCreate(tx *gorm.DB) (err error) {
	allowedTypes := []string{"MANUAL", "SMART"}
	allowedWaterTypes := []string{"MINERAL", "TAP", "MINERALTAP"}

	if !contains(allowedTypes, station.Type) {
		return fmt.Errorf("invalid station type: %s", station.Type)
	}
	if !contains(allowedWaterTypes, station.OfferedWaterTypes) {
		return fmt.Errorf("invalid water types: %s", station.OfferedWaterTypes)
	}
	return nil
}

func (problem *RefillStationProblem) BeforeCreate(tx *gorm.DB) (err error) {
	allowedStatuses := []string{"Inactive", "Active", "In Process"}
	if !contains(allowedStatuses, problem.Status) {
		return fmt.Errorf("invalid problem status: %s", problem.Status)
	}
	return nil
}

func (transaction *WaterTransaction) BeforeCreate(tx *gorm.DB) (err error) {
	allowedWaterTypes := []string{"tap", "mineral"}
	waterType := strings.ToLower(transaction.WaterType)
	if !contains(allowedWaterTypes, waterType) {
		return fmt.Errorf("invalid water type: %s", transaction.WaterType)
	}
	transaction.WaterType = waterType
	return nil
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
