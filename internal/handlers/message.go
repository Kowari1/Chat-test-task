package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Kowari1/TestTask/internal/services"
	"go.uber.org/zap"
)

type CreateMessageRequest struct {
	Text string `json:"text"`
}

type MessageHandler struct {
	messageService *services.MessageService
	logger         *zap.Logger
}

func NewMessageHandler(messageService *services.MessageService, logger *zap.Logger) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
		logger:         logger.With(zap.String("handler", "message")),
	}
}

func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Received create message request")

	chatIDStr := r.PathValue("id")
	chatIDUint64, err := strconv.ParseUint(chatIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrInvalidChatID, h.logger)
		return
	}
	chatID := uint(chatIDUint64)

	var req CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrInvalidJSON, h.logger)
		return
	}

	msg, err := h.messageService.CreateMessage(chatID, req.Text)
	if err != nil {
		if errors.Is(err, services.ErrChatNotFound) {
			writeError(w, http.StatusNotFound, ErrChatNotFound, h.logger)
		} else if errors.Is(err, services.ErrMessageInvalid) {
			writeError(w, http.StatusBadRequest, ErrInvalidInput, h.logger)
		} else {
			h.logger.Error("Unexpected error in CreateMessage", zap.Error(err))
			writeError(w, http.StatusInternalServerError, ErrInternal, h.logger)
		}
		return
	}

	writeJSON(w, http.StatusCreated, msg, h.logger)
}
