package services

import (
	"strings"

	"github.com/Kowari1/TestTask/internal/models"
	"github.com/Kowari1/TestTask/internal/storage"
	"go.uber.org/zap"
)

type MessageService struct {
	chatStore    storage.ChatStore
	messageStore storage.MessageStore
	logger       *zap.Logger
}

func NewMessageService(
	chatStore storage.ChatStore,
	messageStore storage.MessageStore,
	logger *zap.Logger,
) *MessageService {
	return &MessageService{
		chatStore:    chatStore,
		messageStore: messageStore,
		logger:       logger.With(zap.String("service", "message")),
	}
}

func (s *MessageService) CreateMessage(chatID uint, text string) (*models.Message, error) {
	text = strings.TrimSpace(text)
	if len(text) < MinMessageTextLength || len(text) > MaxMessageTextLength {
		s.logger.Warn("Invalid message text",
			zap.Uint("chat_id", chatID),
			zap.Int("length", len(text)),
		)
		return nil, ErrMessageInvalid
	}

	_, err := s.chatStore.GetChatByID(chatID)
	if err != nil {
		s.logger.Warn("Chat not found when sending message",
			zap.Uint("chat_id", chatID),
		)
		return nil, ErrChatNotFound
	}

	msg, err := s.messageStore.CreateMessage(chatID, text)
	if err != nil {
		s.logger.Error("Failed to create message in DB",
			zap.Uint("chat_id", chatID),
			zap.Error(err),
		)
		return nil, err
	}

	s.logger.Info("Message created",
		zap.Uint("message_id", msg.ID),
		zap.Uint("chat_id", chatID),
	)
	return msg, nil
}
