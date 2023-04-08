package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataComparisonResult int

const (
	// New anime, save all episodes
	EntryNotFound DataComparisonResult = iota
	// New episode update on existing anime, save the recent episode(s) only
	NewEpisodeFound
	// No new episode on existing anime, don't do anything
	NoChangesFound
)

type Anime struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Type           string    `json:"type"`
	Summary        string    `json:"summary"`
	Genre          string    `json:"genre"`
	AiringYear     string    `json:"airing_year"`
	Status         string    `json:"status"`
	ImageURL       string    `json:"image_url"`
	LatestEpisode  string    `json:"latest_episode"`
	Episodes       []Episode `json:"episodes"`
	UpdatedAt      time.Time `json:"updated_at"`
	StreamEndpoint string    `json:"-"`
	DetailEndpoint string    `json:"-"`
}

type Episode struct {
	ID       primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Text     string             `json:"text" bson:"text,omitempty"`
	Endpoint string             `json:"endpoint" bson:"endpoint,omitempty"`
}

func (a *Anime) GetEpisodesAsSliceInterface() []interface{} {
	epsLength := len(a.Episodes)
	result := make([]interface{}, epsLength)
	for i := 0; i < epsLength; i++ {
		result[i] = a.Episodes[i]
	}

	return result
}
