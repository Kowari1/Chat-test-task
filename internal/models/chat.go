package models

import "time"

type Chat struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"not null;size:200"`
	CreatedAt time.Time `json:"created_at"`
}
