package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Kowari1/TestTask/internal/config"
	"github.com/Kowari1/TestTask/internal/handlers"
	"github.com/Kowari1/TestTask/internal/services"
	"github.com/Kowari1/TestTask/internal/storage"
	"go.uber.org/zap"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to config file (not used, env only)")
	flag.Parse()

	cfg := config.Load()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	store, err := storage.NewPostgresStorage(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	chatService := services.NewChatService(store, store, logger)
	messageService := services.NewMessageService(store, store, logger)

	chatHandler := handlers.NewChatHandler(chatService, logger)
	messageHandler := handlers.NewMessageHandler(messageService, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /chats/", chatHandler.CreateChat)
	mux.HandleFunc("GET /chats/{id}", chatHandler.GetChat)
	mux.HandleFunc("DELETE /chats/{id}", chatHandler.DeleteChat)
	mux.HandleFunc("POST /chats/{id}/messages/", messageHandler.CreateMessage)

	server := &http.Server{
		Addr:         cfg.ServerPort,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		logger.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Fatal("Server forced to shutdown", zap.Error(err))
		}

		close(done)
	}()

	logger.Info("Server started", zap.String("addr", cfg.ServerPort))

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Server failed", zap.Error(err))
	}

	<-done
	logger.Info("Server stopped")
}
