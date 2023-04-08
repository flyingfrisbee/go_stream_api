package relational

type userRelatedQuery interface {
	InsertUser(userToken string) error
	InsertBookmark(animeID int, userToken, latestEpisode string) error
	DeleteBookmark(animeID int, userToken string) error
}

type userTable struct {
	conn postgresConn
}

func (u *userTable) InsertUser(userToken string) error {
	_, err := u.conn.pool.Exec(u.conn.ctx, insertUserQuery, userToken)
	if err != nil {
		return err
	}

	return nil
}

func (u *userTable) InsertBookmark(animeID int, userToken, latestEpisode string) error {
	_, err := u.conn.pool.Exec(u.conn.ctx, insertBookmarkQuery, userToken, animeID, latestEpisode)
	if err != nil {
		return err
	}

	return nil
}

func (u *userTable) DeleteBookmark(animeID int, userToken string) error {
	_, err := u.conn.pool.Exec(u.conn.ctx, deleteBookmarkQuery, userToken, animeID)
	if err != nil {
		return err
	}

	return nil
}

func (u *userTable) GetUsersTokenByAnimeID(animeID int) ([]string, error) {
	rows, err := u.conn.pool.Query(u.conn.ctx, getUsersTokenByAnimeID, animeID)
	if err != nil {
		return nil, err
	}

	userTokens := []string{}
	for rows.Next() {
		var userToken string
		err = rows.Scan(&userToken)
		if err != nil {
			return nil, err
		}

		userTokens = append(userTokens, userToken)
	}

	return userTokens, nil
}

func (u *userTable) UpdateBookmarkedLatestEpisode(
	animeID int,
	newLatestEpisode string,
) error {
	_, err := u.conn.pool.Exec(
		u.conn.ctx, updateBookmarkedLatestEpisode, newLatestEpisode, animeID,
	)
	if err != nil {
		return err
	}

	return nil
}

var (
	insertUserQuery = `
	INSERT INTO stream_anime.user (user_token) VALUES ($1)
	ON CONFLICT DO NOTHING;`

	insertBookmarkQuery = `
	INSERT INTO stream_anime.user_anime_xref (user_id, anime_id, bookmarked_latest_episode)
	VALUES (
		(SELECT id FROM stream_anime.user WHERE user_token = $1),
		$2, $3
	);`

	deleteBookmarkQuery = `
	DELETE FROM stream_anime.user_anime_xref
	WHERE user_id = (SELECT id FROM stream_anime.user WHERE user_token = $1) AND anime_id = $2;`

	getUsersTokenByAnimeID = `
	SELECT u.user_token FROM stream_anime.user_anime_xref uax 
	JOIN stream_anime.user u
		ON uax.user_id = u.id 
	WHERE anime_id = $1 
		AND uax.bookmarked_latest_episode != (
			SELECT latest_episode FROM stream_anime.anime
			WHERE id = $1
		);`

	updateBookmarkedLatestEpisode = `
	UPDATE stream_anime.user_anime_xref 
	SET bookmarked_latest_episode = $1
	WHERE anime_id = $2;`
)
