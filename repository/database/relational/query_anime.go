package relational

import (
	"fmt"
	"go_stream_api/repository/database/domain"
	"strings"

	"github.com/jackc/pgx/v5"
)

type animeRelatedQuery interface {
	UpsertAnime(*domain.Anime) error
	InsertEpisodes(anime *domain.Anime) error
	GetRecentAnimes(page int) ([]domain.RecentAnime, error)
	GetLatestEpisode(animeID int) (*string, error)
	GetEpisodesCount(anime *domain.Anime) (int, error)
	GetEpisodes(animeID int) ([]domain.Episode, error)
	GetAnimeDetail(animeID int) (domain.Anime, error)
	GetSyncAnime(animeIDs []int) ([]domain.SyncAnime, error)
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

func (a *animeTable) InsertEpisodes(anime *domain.Anime) error {
	episodesAsAny := []interface{}{}

	var sb strings.Builder
	counter := 1
	for _, episode := range anime.Episodes {
		episodesAsAny = append(episodesAsAny, anime.ID, episode.Text, episode.Endpoint)

		valuesText := fmt.Sprintf(`($%d, $%d, $%d),`, counter, counter+1, counter+2)
		sb.WriteString(valuesText)
		counter += 3
	}

	valuesSQL := sb.String()
	// Remove last comma
	valuesSQL = valuesSQL[:len(valuesSQL)-1]

	insertEpisodesQuery := fmt.Sprintf(insertEpisodesFormat, valuesSQL)

	_, err := a.conn.pool.Exec(a.conn.ctx, insertEpisodesQuery, episodesAsAny...)
	if err != nil {
		return err
	}

	return nil
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

func (a *animeTable) GetEpisodes(animeID int) ([]domain.Episode, error) {
	rows, err := a.conn.pool.Query(a.conn.ctx, episodesQuery, animeID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	eps := []domain.Episode{}
	for rows.Next() {
		var ep domain.Episode
		err = rows.Scan(&ep.Text, &ep.Endpoint)
		if err != nil {
			return nil, err
		}
		eps = append(eps, ep)
	}

	return eps, nil
}

func (a *animeTable) GetAnimeDetail(animeID int) (domain.Anime, error) {
	var anime domain.Anime
	row := a.conn.pool.QueryRow(a.conn.ctx, animeDetailQuery, animeID)
	err := row.Scan(
		&anime.ID, &anime.Title, &anime.Type, &anime.Summary,
		&anime.Genre, &anime.AiringYear, &anime.Status,
		&anime.ImageURL, &anime.LatestEpisode, &anime.UpdatedAt,
	)
	if err != nil {
		return anime, err
	}

	return anime, nil
}

func (a *animeTable) GetSyncAnime(animeIDs []int) ([]domain.SyncAnime, error) {
	length := len(animeIDs)
	if length < 1 {
		return nil, fmt.Errorf("anime ids length cannot be less than 1")
	}

	ids := make([]interface{}, length)
	var sb strings.Builder
	for i := 1; i <= length; i++ {
		// Can only use spread operator in Query func if type is []any
		idx := i - 1
		ids[idx] = animeIDs[idx]

		// Write to strings.Builder
		if i == length {
			sb.WriteString(fmt.Sprintf("$%d", i))
			continue
		}
		sb.WriteString(fmt.Sprintf("$%d,", i))
	}

	finalQuery := fmt.Sprintf(syncAnimeQuery, sb.String())
	rows, err := a.conn.pool.Query(a.conn.ctx, finalQuery, ids...)
	if err != nil {
		return nil, err
	}

	result := []domain.SyncAnime{}
	for rows.Next() {
		var sa domain.SyncAnime
		err = rows.Scan(&sa.ID, &sa.LatestEpisode, &sa.UpdatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, sa)
	}

	return result, nil
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

	insertEpisodesFormat = `
	INSERT INTO stream_anime.episode (anime_id, text, endpoint)
	VALUES %s
	ON CONFLICT (endpoint) DO NOTHING;`

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

	episodesQuery = `
	SELECT text, endpoint from stream_anime.episode
	WHERE anime_id = $1
	ORDER BY id ASC;`

	animeDetailQuery = `
	SELECT * FROM stream_anime.anime
	WHERE id = $1;`

	syncAnimeQuery = `
	SELECT id, latest_episode, updated_at FROM stream_anime.anime
	WHERE id IN (%s)
	ORDER BY updated_at DESC;`
)
