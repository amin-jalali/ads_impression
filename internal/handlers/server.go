package handlers

import (
	"go.uber.org/zap"
	"learning/internal/entities"
	"sync"
	"time"
)

type Server struct {
	mu          sync.Mutex
	campaigns   map[string]entities.Campaign
	impressions map[string]map[string]time.Time
	stats       map[string]entities.Stats
	Logger      *zap.Logger
}

//
//func NewServer() *Server {
//	return &Server{
//		campaigns:   make(map[string]entities.Campaign),
//		impressions: make(map[string]map[string]time.Time),
//		stats:       make(map[string]entities.Stats),
//		Logger:      logger.InitLogger(),
//	}
//}
//
//func (s *Server) SetupRoutes(mux *http.ServeMux) {
//	mux.HandleFunc("/api/v1/campaigns", s.CreateCampaignHandler)
//	mux.HandleFunc("/api/v1/impressions", s.TrackImpressionHandler)
//	mux.HandleFunc("/api/v1/campaigns/", s.GetCampaignStatsHandler)
//	mux.HandleFunc("/", NotFoundHandler)
//}
