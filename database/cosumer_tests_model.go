package database

import "time"

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
