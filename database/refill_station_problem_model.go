package database

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

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

func (RefillStationProblem) TableName() string {
	return "refill_station_problem"
}

func (problem *RefillStationProblem) BeforeCreate(tx *gorm.DB) (err error) {
	allowedStatuses := []string{"Inactive", "Active", "In Process"}
	if !contains(allowedStatuses, problem.Status) {
		return fmt.Errorf("invalid problem status: %s", problem.Status)
	}
	return nil
}
