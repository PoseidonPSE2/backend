package database

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
