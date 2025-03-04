package entities

import (
	"sync"
	"time"
)

type Server struct {
	Mu          sync.Mutex
	Campaigns   map[string]Campaign
	Impressions map[string]map[string]time.Time
	Stats       map[string]Stats
}
