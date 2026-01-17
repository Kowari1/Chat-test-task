package models

import "time"

type Message struct {
	ID        uint      `gorm:"primaryKey"`
	ChatID    uint      `gorm:"not null;index"`
	Chat      Chat      `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE"`
	Text      string    `gorm:"not null;size:5000"`
	CreatedAt time.Time `json:"created_at"`
}
