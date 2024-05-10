package main

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// User Model
type User struct {
	ID        uint    `gorm:"primaryKey;autoIncrement"`
	FirstName string  `gorm:"size:100;not null"`
	LastName  string  `gorm:"size:100;not null"`
	Email     *string `gorm:"size:100;unique;default:null"`
	NFCChips  []NFCChip
	Answers   []ConsumerTestAnswer
	Reviews   []RefillStationReview
}

// ConsumerTest Model
type ConsumerTest struct {
	ID        uint `gorm:"primaryKey"`
	Questions []ConsumerTestQuestion
}

// ConsumerTestQuestion Model
type ConsumerTestQuestion struct {
	ID       uint     `gorm:"primaryKey"`
	Text     string   `gorm:"size:255;not null"`
	MinValue *float32 `gorm:"default:null"`
	MaxValue *float32 `gorm:"default:null"`
	TestID   uint     `gorm:"not null"`
	Answers  []ConsumerTestAnswer
}

// ConsumerTestAnswer Model
type ConsumerTestAnswer struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"not null"`
	QuestionID uint      `gorm:"not null"`
	Answer     float32   `gorm:"not null"`
	Timestamp  time.Time `gorm:"autoCreateTime"`
}

// NFCChip Model
type NFCChip struct {
	ID                uint    `gorm:"primaryKey"`
	UserID            uint    `gorm:"not null"`
	HardwareID        string  `gorm:"size:32;unique;not null"`
	FillVolume        int     `gorm:"not null"`
	WaterType         string  `gorm:"size:16;not null"`
	PathImage         *string `gorm:"size:255;default:null"`
	Active            bool    `gorm:"default:true"`
	WaterTransactions []WaterTransaction
}

// RefillStation Model
type RefillStation struct {
	ID                uint    `gorm:"primaryKey"`
	Name              string  `gorm:"size:100;not null"`
	Description       string  `gorm:"size:255;not null"`
	Latitude          float64 `gorm:"not null"`
	Longitude         float64 `gorm:"not null"`
	Address           string  `gorm:"size:255;not null"`
	LikeCounter       int     `gorm:"default:0"`
	WaterSource       string  `gorm:"size:50;not null"`
	OpeningTimes      string  `gorm:"size:100;not null"`
	Active            bool    `gorm:"default:true"`
	Type              string  `gorm:"size:16;not null"`
	OfferedWaterTypes string  `gorm:"size:32;not null"`
	Reviews           []RefillStationReview
	Problems          []RefillStationProblem
	WaterTransactions []WaterTransaction
}

// RefillStationReview Model
type RefillStationReview struct {
	ID            uint      `gorm:"primaryKey"`
	StationID     uint      `gorm:"not null"`
	UserID        uint      `gorm:"not null"`
	Cleanness     int       `gorm:"not null;check:cleanness >= 1 AND cleanness <= 5"`
	Accessibility int       `gorm:"not null;check:accessibility >= 1 AND accessibility <= 5"`
	WaterQuality  int       `gorm:"not null;check:water_quality >= 1 AND water_quality <= 5"`
	Timestamp     time.Time `gorm:"autoCreateTime"`
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
type RefillStationProblem struct {
	ID          uint      `gorm:"primaryKey"`
	StationID   uint      `gorm:"not null"`
	Title       string    `gorm:"size:100;not null"`
	Description string    `gorm:"size:255;not null"`
	Status      string    `gorm:"size:16;not null"`
	LinkToMedia *string   `gorm:"size:255;default:null"`
	Timestamp   time.Time `gorm:"autoCreateTime"`
}

// WaterTransaction Model
type WaterTransaction struct {
	ID        uint      `gorm:"primaryKey"`
	StationID uint      `gorm:"not null"`
	ChipID    *uint     `gorm:"default:null"`
	Volume    int       `gorm:"not null"`
	WaterType string    `gorm:"size:16;not null"`
	Timestamp time.Time `gorm:"autoCreateTime"`
	Guest     bool      `gorm:"default:false"`
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
