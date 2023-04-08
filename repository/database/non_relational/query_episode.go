package nonrelational

import (
	"fmt"
	"go_stream_api/repository/database/domain"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type episodeRelatedQuery interface {
	// Insert last n episodes to mongoDB. If n < 1 then all episodes will be inserted
	InsertEpisodes(a *domain.Anime, eps []interface{}, n int) error
	GetEpisodes(animeID int) ([]domain.Episode, error)
	GetEpisodesCount(animeID int) (int, error)
}

type episodeCollections struct {
	conn mongoDBConn
}

func (ec *episodeCollections) InsertEpisodes(a *domain.Anime, eps []interface{}, n int) error {
	id := strconv.Itoa(a.ID)
	collection := ec.conn.client.Database("not_episodes").Collection(id)

	startIndex := 0
	if n >= 1 {
		startIndex = len(eps) - n
		if startIndex < 0 {
			return fmt.Errorf("starting index for episodes slice cannot be less than 0, got %d", startIndex)
		}
	}

	eps = eps[startIndex:]
	result, err := collection.InsertMany(ec.conn.ctx, eps)
	if err != nil {
		return err
	}

	failedInsertingDocuments := len(result.InsertedIDs) == 0
	if failedInsertingDocuments {
		return fmt.Errorf("failed inserting episodes to MongoDB for anime: %s", a.Title)
	}

	return nil
}

func (ec *episodeCollections) GetEpisodes(animeID int) ([]domain.Episode, error) {
	id := strconv.Itoa(animeID)
	collection := ec.conn.client.Database("not_episodes").Collection(id)

	// Sort episodes ascending
	opts := options.Find().SetSort(bson.D{{"_id", 1}})
	emptyFilter := bson.D{}

	cursor, err := collection.Find(ec.conn.ctx, emptyFilter, opts)
	if err != nil {
		return nil, err
	}

	var episodes []domain.Episode
	err = cursor.All(ec.conn.ctx, &episodes)
	if err != nil {
		return nil, err
	}

	episodesAreEmpty := len(episodes) == 0
	if episodesAreEmpty {
		return nil, fmt.Errorf("failed getting episodes from MongoDB for animeID: %d", animeID)
	}

	return episodes, nil
}
func (ec *episodeCollections) GetEpisodesCount(animeID int) (int, error) {
	id := strconv.Itoa(animeID)
	collection := ec.conn.client.Database("not_episodes").Collection(id)

	docsCount, err := collection.CountDocuments(ec.conn.ctx, bson.D{})
	if err != nil {
		return 0, err
	}

	return int(docsCount), nil
}
