package main

import (
	"learning/cmd/server"
	"net/http"
)

func main() {
	//lgr := logger.InitLogger()

	port, err := server.GetEnv("port")

	if err != nil {
		//logger.Log.Error(err.Error())
		return
	}

	srv := &http.Server{Addr: ":" + port}
	err = server.Run(srv.ListenAndServe)
	if err != nil {
		//logger.Log.Error(err.Error())
	}
}
