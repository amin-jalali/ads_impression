package main

import (
	"learning/cmd/server"
	"learning/internal/logger"
	"net/http"
)

func main() {
	logger.InitLogger()
	defer logger.Sync()

	port, err := server.GetEnv("port")

	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	srv := &http.Server{Addr: ":" + port}
	err = server.Run(srv.ListenAndServe)
	if err != nil {
		logger.Log.Error(err.Error())
	}
}
