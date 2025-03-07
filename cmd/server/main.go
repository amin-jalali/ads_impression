package server

import (
	"learning/internal/handlers"
	"learning/internal/logger"
	"learning/internal/repositories/memory"
	"net/http"
)

var SetupServer = setupServer

func setupServer() http.Handler {
	mux := http.NewServeMux()

	// Initialize shared in-memory server
	memServer := memory.NewServer()

	// Pass shared memory to repositories
	campaignRepo := memory.NewInMemoryCampaignRepository(memServer)
	impressionRepo := memory.NewInMemoryImpressionRepository(memServer)
	statsRepo := memory.NewInMemoryStatsRepository(memServer)

	campaignHandler := handlers.NewCampaignHandler(campaignRepo)
	impressionHandler := handlers.NewImpressionHandler(impressionRepo)
	statsHandler := handlers.NewStatsHandler(statsRepo)

	mux.HandleFunc("/api/v1/campaigns", campaignHandler.CreateCampaignHandler)
	mux.HandleFunc("/api/v1/impressions", impressionHandler.TrackImpressionHandler)
	mux.HandleFunc("/api/v1/campaigns/stats/", func(w http.ResponseWriter, r *http.Request) {
		statsHandler.GetCampaignStatsHandler(w, r)
	})
	mux.HandleFunc("/", handlers.NotFoundHandler)

	return mux
}

func Run(listenAndServe func() error) error {
	logger.InitLogger()
	defer logger.Sync()

	logger.Log.Info("Server started on port 8080...")

	err := listenAndServe()
	if err != nil {
		return err
	}
	return nil
}
