package storage

import "github.com/Kowari1/TestTask/internal/models"

type DB interface {
	NewDB(dsn string) (*DB, error)
}

type ChatStore interface {
	CreateChat(title string) (*models.Chat, error)
	GetChatByID(id uint) (*models.Chat, error)
	DeleteChat(id uint) error
}

type MessageStore interface {
	CreateMessage(chatID uint, text string) (*models.Message, error)
	GetLastMessageByChatID(chatID uint, limit int) ([]models.Message, error)
}
