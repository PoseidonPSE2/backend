package database

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

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
