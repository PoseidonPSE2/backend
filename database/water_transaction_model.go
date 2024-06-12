package database

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

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

func (transaction *WaterTransaction) BeforeCreate(tx *gorm.DB) (err error) {
	allowedWaterTypes := []string{"tap", "mineral"}
	waterType := strings.ToLower(transaction.WaterType)
	if !contains(allowedWaterTypes, waterType) {
		return fmt.Errorf("invalid water type: %s", transaction.WaterType)
	}
	transaction.WaterType = waterType
	return nil
}
