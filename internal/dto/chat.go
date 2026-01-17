package dto

import "time"

type ChatResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatWithMessagesResponse struct {
	Chat     ChatResponse      `json:"chat"`
	Messages []MessageResponse `json:"messages"`
}
