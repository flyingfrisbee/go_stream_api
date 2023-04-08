package relational

import (
	"fmt"
	"go_stream_api/repository/database/domain"

	"github.com/jackc/pgx/v5"
)

type animeRelatedQuery interface {
	UpsertAnime(*domain.Anime) error
	// Insert last n episodes to postgres. If n < 1 then all episodes will be inserted
	InsertEpisodesToPostgres(anime *domain.Anime, n int) (int, error)
	GetRecentAnimes(page int) ([]domain.RecentAnime, error)
	GetLatestEpisode(animeID int) (*string, error)
	GetEpisodesCount(anime *domain.Anime) (int, error)
}

type animeTable struct {
	conn postgresConn
}

func (a *animeTable) UpsertAnime(anime *domain.Anime) error {
	_, err := a.conn.pool.Exec(
		a.conn.ctx, upsertAnimeQuery, anime.ID, anime.Title,
		anime.Type, anime.Summary, anime.Genre, anime.AiringYear,
		anime.Status, anime.ImageURL, anime.LatestEpisode, anime.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (a *animeTable) InsertEpisodesToPostgres(anime *domain.Anime, n int) (int, error) {
	values := [][]interface{}{}

	startIndex := 0
	if n >= 1 {
		startIndex = len(anime.Episodes) - n
		if startIndex < 0 {
			return 0, fmt.Errorf("starting index for episodes slice cannot be less than 0, got %d", startIndex)
		}
	}

	for _, episode := range anime.Episodes[startIndex:] {
		row := []interface{}{anime.ID, episode.Text, episode.Endpoint}
		values = append(values, row)
	}

	insertedRow, err := a.conn.pool.CopyFrom(
		a.conn.ctx, pgx.Identifier{"stream_anime", "episode"},
		[]string{"anime_id", "text", "endpoint"}, pgx.CopyFromRows(values),
	)
	if err != nil {
		return 0, err
	}

	return int(insertedRow), nil
}

func (a *animeTable) GetRecentAnimes(page int) ([]domain.RecentAnime, error) {
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * 20
	rows, err := a.conn.pool.Query(a.conn.ctx, recentAnimeQuery, offset)
	if err != nil {
		return nil, err
	}

	result := []domain.RecentAnime{}
	for rows.Next() {
		var anime domain.RecentAnime

		err = rows.Scan(&anime.ID, &anime.Title, &anime.ImageURL, &anime.LatestEpisode, &anime.UpdatedAt)
		if err != nil {
			return nil, err
		}

		result = append(result, anime)
	}

	return result, nil
}

func (a *animeTable) GetLatestEpisode(animeID int) (*string, error) {
	var latestEpisode string
	row := a.conn.pool.QueryRow(a.conn.ctx, latestEpisodeQuery, animeID)
	err := row.Scan(&latestEpisode)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &latestEpisode, nil
}

func (a *animeTable) GetEpisodesCount(anime *domain.Anime) (int, error) {
	epsCount := 0
	row := a.conn.pool.QueryRow(a.conn.ctx, episodesCountQuery, anime.ID)
	err := row.Scan(&epsCount)
	if err != nil {
		return epsCount, err
	}

	return epsCount, nil
}

var (
	upsertAnimeQuery = `
	INSERT INTO stream_anime.anime (
		id, title, type, summary, genre, airing_year, 
		status, image_url, latest_episode, updated_at
	)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
	ON CONFLICT (id) 
	DO UPDATE SET
		title = EXCLUDED.title,
		type = EXCLUDED.type,
		summary = EXCLUDED.summary,
		genre = EXCLUDED.genre,
		airing_year = EXCLUDED.airing_year,
		status = EXCLUDED.status,
		image_url = EXCLUDED.image_url,
		latest_episode = EXCLUDED.latest_episode,
		updated_at = EXCLUDED.updated_at;`

	recentAnimeQuery = `
	SELECT id, title, image_url, latest_episode, updated_at
	FROM stream_anime.anime
	ORDER BY updated_at DESC
	LIMIT 20
	OFFSET $1;`

	latestEpisodeQuery = `
	SELECT latest_episode
	FROM stream_anime.anime
	WHERE id = $1;`

	episodesCountQuery = `
	SELECT COUNT(e.id) FROM stream_anime.anime a
	JOIN stream_anime.episode e
		ON a.id = e.anime_id
	WHERE a.id = $1;`
)
