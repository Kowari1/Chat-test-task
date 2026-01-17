package mappers

import (
	"github.com/Kowari1/TestTask/internal/dto"
	"github.com/Kowari1/TestTask/internal/models"
)

func ToChatDTO(chat *models.Chat) dto.ChatResponse {
	return dto.ChatResponse{
		ID:        chat.ID,
		Title:     chat.Title,
		CreatedAt: chat.CreatedAt,
	}
}

func ToMessageDTO(msg *models.Message) dto.MessageResponse {
	return dto.MessageResponse{
		ID:        msg.ID,
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt,
	}
}

func ToMessageDTOList(messages []models.Message) []dto.MessageResponse {
	result := make([]dto.MessageResponse, 0, len(messages))
	for _, m := range messages {
		msg := m
		result = append(result, ToMessageDTO(&msg))
	}
	return result
}
