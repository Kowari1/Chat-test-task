package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Kowari1/TestTask/internal/mappers"
	"github.com/Kowari1/TestTask/internal/models"
	"github.com/Kowari1/TestTask/internal/services"
	"go.uber.org/zap"
)

type ChatService interface {
	CreateChat(title string) (*models.Chat, error)
	GetChatWithMessages(chatID uint, limit int) (*models.Chat, []models.Message, error)
	DeleteChat(chatID uint) error
}

type ChatHandler struct {
	chatService ChatService
	logger      *zap.Logger
}

type CreateChatRequest struct {
	Title string `json:"title"`
}

func NewChatHandler(chatService ChatService, logger *zap.Logger) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
		logger:      logger.With(zap.String("handler", "chat")),
	}
}

func (h *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Create chat request")

	var req CreateChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrInvalidJSON, h.logger)
		return
	}

	chat, err := h.chatService.CreateChat(req.Title)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, ErrInvalidInput, h.logger)
		} else {
			h.logger.Error("Failed to create chat", zap.Error(err))
			writeError(w, http.StatusInternalServerError, ErrInternal, h.logger)
		}
		return
	}

	response := struct {
		Chat any `json:"chat"`
	}{
		Chat: mappers.ToChatDTO(chat),
	}

	writeJSON(w, http.StatusCreated, response, h.logger)
}

func (h *ChatHandler) GetChat(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrInvalidChatID, h.logger)
		return
	}

	limit := services.DefaultMessageLimit
	if q := r.URL.Query().Get("limit"); q != "" {
		if l, err := strconv.Atoi(q); err == nil {
			limit = l
		}
	}

	chat, messages, err := h.chatService.GetChatWithMessages(uint(id), limit)
	if err != nil {
		if errors.Is(err, services.ErrChatNotFound) {
			writeError(w, http.StatusNotFound, ErrChatNotFound, h.logger)
		} else {
			h.logger.Error("Failed to get chat", zap.Error(err))
			writeError(w, http.StatusInternalServerError, ErrInternal, h.logger)
		}
		return
	}

	response := mappers.ToChatDTO(chat)

	writeJSON(w, http.StatusOK, map[string]any{
		"chat":     response,
		"messages": mappers.ToMessageDTOList(messages),
	}, h.logger)
}

func (h *ChatHandler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrInvalidChatID, h.logger)
		return
	}

	if err := h.chatService.DeleteChat(uint(id)); err != nil {
		h.logger.Error("Failed to delete chat", zap.Error(err))
		writeError(w, http.StatusInternalServerError, ErrInternal, h.logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
