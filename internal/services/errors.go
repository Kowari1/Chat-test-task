package services

import "errors"

var (
	ErrInvalidInput   = errors.New("invalid input")
	ErrChatNotFound   = errors.New("chat not found")
	ErrMessageTooLong = errors.New("message too long")
	ErrMessageInvalid = errors.New("invalid message text")
)
