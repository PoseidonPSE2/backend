package database

import (
	"fmt"

	"gorm.io/gorm"
)

// Like Model
// @swagger:model
type Like struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	StationID uint `gorm:"not null" json:"station_id"`
	UserID    uint `gorm:"not null" json:"user_id"`
}

func (like *Like) BeforeCreate(tx *gorm.DB) (err error) {
	var count int64
	// Check if like is already present in database (Same Refill Station and User)
	if result := tx.Model(&Like{}).Where(&Like{StationID: like.StationID, UserID: like.UserID}).Count(&count); result.Error == nil {
		if count != 0 {
			return fmt.Errorf("Like already exists")
		}
	}
	return nil
}
