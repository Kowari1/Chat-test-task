package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Kowari1/TestTask/internal/handlers"
	"github.com/Kowari1/TestTask/internal/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type mockChatService struct{}

func (m *mockChatService) CreateChat(title string) (*models.Chat, error) {
	return &models.Chat{
		ID:        1,
		Title:     title,
		CreatedAt: time.Date(2026, 1, 1, 1, 1, 1, 0, time.UTC),
	}, nil
}

func (m *mockChatService) GetChatWithMessages(uint, int) (*models.Chat, []models.Message, error) {
	return nil, nil, nil
}

func (m *mockChatService) DeleteChat(uint) error {
	return nil
}

func TestCreateChat_ReturnsDTO(t *testing.T) {
	logger := zap.NewNop()

	handler := handlers.NewChatHandler(&mockChatService{}, logger)

	body := `{"title":"My chat"}`
	req := httptest.NewRequest(http.MethodPost, "/chats", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	handler.CreateChat(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var resp struct {
		Chat struct {
			ID        uint   `json:"id"`
			Title     string `json:"title"`
			CreatedAt string `json:"created_at"`
		} `json:"chat"`
	}

	err := json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)

	assert.Equal(t, uint(1), resp.Chat.ID)
	assert.Equal(t, "My chat", resp.Chat.Title)
	assert.NotEmpty(t, resp.Chat.CreatedAt)
}
