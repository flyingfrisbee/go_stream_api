package domain

import "time"

type SyncAnime struct {
	ID            int       `json:"id"`
	LatestEpisode string    `json:"latest_episode"`
	UpdatedAt     time.Time `json:"updated_at"`
}
