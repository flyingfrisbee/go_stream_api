package relational

import (
	"go_stream_api/repository/database/domain"
)

type animeRelatedQuery interface {
	UpsertAnime(*domain.Anime) error
	GetRecentAnimes(page int) ([]domain.RecentAnime, error)
}

type animeTable struct {
	conn *postgresConn
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

// func (conn *rdbmsConnImpl) bulkInsert(
// 	schemaName, tableName string,
// 	columns []string,
// 	values [][]interface{},
// ) (int, error) {
// 	rowsInsertedCount, err := conn.pool.CopyFrom(
// 		conn.ctx, pgx.Identifier{schemaName, tableName},
// 		columns, pgx.CopyFromRows(values),
// 	)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return int(rowsInsertedCount), nil
// }

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
)
