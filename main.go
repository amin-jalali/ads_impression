package main

import (
	"fmt"
	"learning/cmd/config"
	"learning/cmd/server"
	logger2 "learning/internal/logger"
	"net/http"
)

func main() {

	logger := logger2.InitLogger()
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error(err.Error())
	}

	// Read port from config and convert to string
	port := fmt.Sprintf(":%d", cfg.Server.Port)

	// Set up the server
	srv := server.SetupServer()

	// Start the server using config values
	err = server.Run(func() error {
		return http.ListenAndServe(port, srv)
	})
	if err != nil {
		logger.Fatal(err.Error())
	}
}
