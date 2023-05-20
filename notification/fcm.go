package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	env "go_stream_api/environment"
	db "go_stream_api/repository/database"
	"go_stream_api/repository/database/domain"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type fcmMessage struct {
	To   string      `json:"to"`
	Data interface{} `json:"data"`
}

type fcmData struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Body          string    `json:"body"`
	LatestEpisode string    `json:"latest_episode"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type fcmResponse struct {
	Success int `json:"success"`
	Failure int `json:"failure"`
}

var (
	wg sync.WaitGroup
)

func SendNotificationMessageToUsers(
	comparison domain.DataComparisonResult,
	anime *domain.Anime,
) error {
	if comparison != domain.NewEpisodeFound {
		return nil
	}

	usersToken, err := db.Conn.Pg.GetUsersTokenByAnimeID(anime.ID)
	if err != nil {
		return err
	}

	if len(usersToken) < 1 {
		return nil
	}

	for _, userToken := range usersToken {
		wg.Add(1)
		go handleNotificationTask(&wg, anime, userToken)
	}
	wg.Wait()

	err = db.Conn.Pg.UpdateBookmarkedLatestEpisode(anime.ID, anime.LatestEpisode)
	if err != nil {
		return err
	}

	return nil
}

func handleNotificationTask(wg *sync.WaitGroup, anime *domain.Anime, userToken string) {
	defer wg.Done()

	shouldWipeInactiveUserData, err := sendNotification(anime, userToken)
	if err != nil {
		log.Println(err)
		return
	}

	if shouldWipeInactiveUserData {
		err = wipeInactiveUserData(userToken)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func sendNotification(anime *domain.Anime, userToken string) (bool, error) {
	fcmMessage := fcmMessage{
		To: userToken,
		Data: fcmData{
			ID:            anime.ID,
			Title:         anime.Title,
			Body:          fmt.Sprintf("New episode: %s is available to watch", anime.LatestEpisode),
			LatestEpisode: anime.LatestEpisode,
			UpdatedAt:     anime.UpdatedAt,
		},
	}

	jsonBytes, err := json.Marshal(fcmMessage)
	if err != nil {
		return false, err
	}

	r, err := http.NewRequest("POST", env.FCMURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return false, err
	}

	r.Header.Set("Authorization", fmt.Sprintf("key=%s", env.FCMKey))
	r.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	jsonBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	response := fcmResponse{}
	err = json.Unmarshal(jsonBytes, &response)
	if err != nil {
		return false, err
	}

	return (response.Failure == 1), nil
}

func wipeInactiveUserData(userToken string) error {
	err := db.Conn.Pg.DeleteBookmarkByUserToken(userToken)
	if err != nil {
		return err
	}

	err = db.Conn.Pg.DeleteUser(userToken)
	if err != nil {
		return err
	}

	return nil
}
