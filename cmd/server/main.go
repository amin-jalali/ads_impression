package server

import (
	"learning/internal/handlers"
	"learning/internal/logger"
	"net/http"
)

var ListenAndServe = http.ListenAndServe
var SetupServer = setupServer

func setupServer() http.Handler {
	mux := http.NewServeMux()

	s := handlers.NewServer()

	mux.HandleFunc("/api/v1/campaigns", s.CreateCampaignHandler)
	mux.HandleFunc("/api/v1/impressions", s.TrackImpressionHandler)
	mux.HandleFunc("/api/v1/campaigns/", s.GetCampaignStatsHandler)
	mux.HandleFunc("/", handlers.NotFoundHandler)

	return mux
}

func Run(listenAndServe func() error) error {
	logger.InitLogger()
	defer logger.Sync()

	logger.Log.Info("Server started ...")

	err := listenAndServe()
	if err != nil {
		return err
	}
	return nil
}
