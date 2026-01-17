package services

import (
	"strings"

	"github.com/Kowari1/TestTask/internal/models"
	"github.com/Kowari1/TestTask/internal/storage"
	"go.uber.org/zap"
)

type ChatService struct {
	chatStore    storage.ChatStore
	messageStore storage.MessageStore
	logger       *zap.Logger
}

func NewChatService(
	chatStore storage.ChatStore,
	messageStore storage.MessageStore,
	logger *zap.Logger,
) *ChatService {
	return &ChatService{
		chatStore:    chatStore,
		messageStore: messageStore,
		logger:       logger.With(zap.String("service", "chat")),
	}
}

func (s *ChatService) CreateChat(title string) (*models.Chat, error) {
	title = strings.TrimSpace(title)
	if len(title) < MinChatTitleLength || len(title) > MaxChatTitleLength {
		s.logger.Warn("Invalid chat title", zap.String("title", title))
		return nil, ErrInvalidInput
	}

	chat, err := s.chatStore.CreateChat(title)
	if err != nil {
		s.logger.Error("Failed to create chat in DB", zap.String("title", title), zap.Error(err))
		return nil, err
	}

	s.logger.Info("Chat created", zap.Uint("chat_id", chat.ID), zap.String("title", chat.Title))
	return chat, nil
}

func (s *ChatService) GetChatWithMessages(chatID uint, limit int) (*models.Chat, []models.Message, error) {
	if limit <= 0 {
		limit = DefaultMessageLimit
	}
	if limit > MaxMessageLimit {
		limit = MaxMessageLimit
	}

	chat, err := s.chatStore.GetChatByID(chatID)
	if err != nil {
		s.logger.Warn("Chat not found", zap.Uint("chat_id", chatID))
		return nil, nil, ErrChatNotFound
	}

	messages, err := s.messageStore.GetLastMessageByChatID(chatID, limit)
	if err != nil {
		s.logger.Error("Failed to fetch messages", zap.Uint("chat_id", chatID), zap.Error(err))
		return nil, nil, err
	}

	s.logger.Debug("Fetched chat with messages",
		zap.Uint("chat_id", chatID),
		zap.Int("message_count", len(messages)),
		zap.Int("limit", limit),
	)
	return chat, messages, nil
}

func (s *ChatService) DeleteChat(chatID uint) error {
	err := s.chatStore.DeleteChat(chatID)
	if err != nil {
		s.logger.Error("Failed to delete chat", zap.Uint("chat_id", chatID), zap.Error(err))
		return err
	}
	s.logger.Info("Chat deleted", zap.Uint("chat_id", chatID))
	return nil
}
