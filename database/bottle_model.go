package database

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

var BottleWaterTypes []string = []string{"tap", "mineral"}

// Bottle Model (previously NFCChip)
// @swagger:model
type Bottle struct {
	ID                uint               `gorm:"primaryKey" json:"id"`
	UserID            uint               `gorm:"not null" json:"user_id"`
	NFCID             string             `gorm:"size:20; json:"nfc_id"`
	FillVolume        int                `gorm:"not null" json:"fill_volume"`
	WaterType         string             `gorm:"size:16;not null" json:"water_type"`
	Title             string             `gorm:"size:16;not null" json:"title"`
	BottleImage       *string            `gorm:"type:TEXT;default:null" json:"bottle_image,omitempty"`
	Active            bool               `gorm:"default:true" json:"active"`
	WaterTransactions []WaterTransaction `gorm:"foreignKey:BottleID" json:"-"`
}

// Enumeration Constraints:
func (Bottle) TableName() string {
	return "bottle"
}

func (bottle *Bottle) BeforeCreate(tx *gorm.DB) (err error) {
	waterType := strings.ToLower(bottle.WaterType)
	if !contains(BottleWaterTypes, waterType) {
		return fmt.Errorf("invalid water type: %s", bottle.WaterType)
	}
	bottle.WaterType = waterType

	// If there is already a bottle with this nfc id no second can be added with same id, empty id is always possible
	var count int64
	if result := tx.Model(&Bottle{}).Where("nfc_id = ?", bottle.NFCID).Count(&count); result.Error == nil {
		if count != 0 && bottle.NFCID != "" {
			return fmt.Errorf("NFC ID '%s' already exists", bottle.NFCID)
		}
	}

	return nil
}
