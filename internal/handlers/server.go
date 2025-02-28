package handlers

import (
	"go.uber.org/zap"
	"learning/internal/logger"
	"sync"
	"time"
)

type Campaign struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	StartTime time.Time `json:"start_time"`
}

type Impression struct {
	CampaignID string    `json:"campaign_id"`
	Timestamp  time.Time `json:"timestamp"`
	UserID     string    `json:"user_id"`
	AdID       string    `json:"ad_id"`
}

type Stats struct {
	CampaignID string `json:"campaign_id"`
	LastHour   int64  `json:"last_hour"`
	LastDay    int64  `json:"last_day"`
	TotalCount int64  `json:"total"`
}

type Server struct {
	mu          sync.Mutex
	campaigns   map[string]Campaign
	impressions map[string]map[string]time.Time
	stats       map[string]Stats
	Logger      *zap.Logger
}

func NewServer() *Server {
	logger.InitLogger()
	defer logger.Sync()

	return &Server{
		campaigns:   make(map[string]Campaign),
		impressions: make(map[string]map[string]time.Time),
		stats:       make(map[string]Stats),
		Logger:      logger.Log,
	}
}
