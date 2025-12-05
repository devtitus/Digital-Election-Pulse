package models

import (
	"time"

	"gorm.io/gorm"
)

type Party struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `json:"name"`
	Leader    string         `json:"leader"`
	ColorHex  string         `json:"color_hex"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type SentimentSnapshot struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	PartyID         uint      `json:"party_id"`
	Party           Party     `gorm:"foreignKey:PartyID" json:"-"`
	Score           float64   `json:"score"`
	KeyTopics       string    `gorm:"type:jsonb" json:"key_topics"` // Stores JSON array of strings
	Emotion         string    `json:"emotion"`
	SourceBreakdown string    `gorm:"type:jsonb" json:"source_breakdown"`
	CreatedAt       time.Time `json:"created_at"`
}
