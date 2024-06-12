package database

import (
	"fmt"

	"gorm.io/gorm"
)

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
