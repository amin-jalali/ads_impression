package main

import (
	"learning/cmd/server"
	"learning/internal/logger"
	"net/http"
)

func main() {
	logger.InitLogger()
	defer logger.Sync()

	srv := &http.Server{Addr: ":8081"}
	err := server.Run(srv.ListenAndServe)
	if err != nil {
		logger.Log.Error(err.Error())
	}
}
