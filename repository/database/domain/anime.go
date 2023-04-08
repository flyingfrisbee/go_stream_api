package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
