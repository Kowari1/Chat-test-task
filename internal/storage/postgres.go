package storage

import (
	"fmt"

	"github.com/Kowari1/TestTask/internal/config"
	"github.com/Kowari1/TestTask/internal/models"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresStorage struct {
	db *gorm.DB
}

func NewPostgresStorage(cfg *config.Config, logger *zap.Logger) (*PostgresStorage, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.DBConnMaxLifetime)

	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) CreateChat(title string) (*models.Chat, error) {
	chat := &models.Chat{Title: title}
	err := s.db.Create(chat).Error

	return chat, err
}

func (s *PostgresStorage) GetChatByID(id uint) (*models.Chat, error) {
	var chat models.Chat
	err := s.db.First(&chat, id).Error
	return &chat, err
}

func (s *PostgresStorage) DeleteChat(id uint) error {
	return s.db.Delete(&models.Chat{}, id).Error
}

func (s *PostgresStorage) CreateMessage(chatID uint, text string) (*models.Message, error) {
	msg := &models.Message{ChatID: chatID, Text: text}
	err := s.db.Create(msg).Error
	return msg, err
}

func (s *PostgresStorage) GetLastMessageByChatID(chatID uint, limit int) ([]models.Message, error) {
	var messages []models.Message
	err := s.db.
		Where("chat_id = ?", chatID).
		Limit(limit).
		Find(&messages).Error

	return messages, err
}
