package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"video-agent-go/config"
	"video-agent-go/handler"
	"video-agent-go/model"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	// Initialize configuration
	config.Init()

	// Initialize database
	model.InitDB()

	// Create Hertz server
	h := server.Default(
		server.WithHostPorts(fmt.Sprintf("%s:%d",
			config.AppConfig.Server.Host,
			config.AppConfig.Server.Port)),
	)

	// Register routes
	handler.RegisterRoutes(h)

	// Setup graceful shutdown
	setupGracefulShutdown(h)

	log.Printf("🚀 Video Agent server starting on %s:%d",
		config.AppConfig.Server.Host,
		config.AppConfig.Server.Port)

	// Start server
	h.Spin()
}

func setupGracefulShutdown(h *server.Hertz) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down server...")

		// Close database connection
		if model.DB != nil {
			model.DB.Close()
		}

		// Shutdown server
		if err := h.Shutdown(); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}

		os.Exit(0)
	}()
}
