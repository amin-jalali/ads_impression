package server

import (
	"learning/handlers"
	"learning/internal/logger"
	"learning/internal/repositories/memory"
	"net/http"
)

var ListenAndServe = http.ListenAndServe

var SetupServer = setupServer

func setupServer() http.Handler {
	mux := http.NewServeMux()

	//s := memory.NewServer()

	campaignRepo := memory.NewInMemoryCampaignRepository()
	impressionRepo := memory.NewInMemoryImpressionRepository()
	statsRepo := memory.NewInMemoryStatsRepository()

	campaignHandler := handlers.NewCampaignHandler(campaignRepo)
	impressionHandler := handlers.NewImpressionHandler(impressionRepo)
	statsHandler := handlers.NewStatsHandler(statsRepo)

	mux.HandleFunc("/api/v1/campaigns", campaignHandler.CreateCampaignHandler)
	mux.HandleFunc("/api/v1/impressions", impressionHandler.TrackImpressionHandler)
	mux.HandleFunc("/api/v1/campaigns/", statsHandler.GetCampaignStatsHandler)
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
