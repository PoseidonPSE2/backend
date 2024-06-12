package database

// Like Model
// @swagger:model
type Like struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	StationID uint `gorm:"not null" json:"station_id"`
	UserID    uint `gorm:"not null" json:"user_id"`
}
