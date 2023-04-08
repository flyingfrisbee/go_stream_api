package domain

import "time"

type RecentAnime struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	ImageURL      string    `json:"image_url"`
	LatestEpisode string    `json:"latest_episode"`
	UpdatedAt     time.Time `json:"updated_at"`
}
