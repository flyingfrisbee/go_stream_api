package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	env "go_stream_api/environment"
	db "go_stream_api/repository/database"
	"go_stream_api/repository/database/domain"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"golang.org/x/oauth2/google"
)

type fcmMessage struct {
	Message message `json:"message"`
}

type message struct {
	Token string      `json:"token"`
	Data  interface{} `json:"data"`
}

type fcmData struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Body          string    `json:"body"`
	LatestEpisode string    `json:"latest_episode"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// when success sending message: status_code 200
// {
// 	"name": "projects/{project_id}/messages/{some_generated_code}"
// }

// user doesn't exist: status_code 404
// {
// 	"error": {
// 		"code": 404,
// 		"message": "Requested entity was not found.",
// 		"status": "NOT_FOUND",
// 		"details": [
// 			{
// 				"@type": "type.googleapis.com/google.firebase.fcm.v1.FcmError",
// 				"errorCode": "UNREGISTERED"
// 			}
// 		]
// 	}
// }

var (
	wg         sync.WaitGroup
	oauthToken string
)

func StartOAuthTokenGenerator() {
	token, err := generateOAuthToken()
	if err != nil {
		log.Fatal(err)
	}
	oauthToken = token

	for {
		select {
		// Expiration of OAuth token is 60 minutes, hence the 55 minutes hardcode
		case <-time.After(55 * time.Minute):
			token, err := generateOAuthToken()
			if err != nil {
				log.Fatal(err)
			}
			oauthToken = token
		}
	}
}

func generateOAuthToken() (string, error) {
	cred, err := google.FindDefaultCredentials(context.Background(), "https://www.googleapis.com/auth/firebase.messaging")
	if err != nil {
		return "", err
	}

	token, err := cred.TokenSource.Token()
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

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
		Message: message{
			Token: userToken,
			Data: fcmData{
				ID:            strconv.Itoa(anime.ID),
				Title:         anime.Title,
				Body:          fmt.Sprintf("New episode: %s is available to watch", anime.LatestEpisode),
				LatestEpisode: anime.LatestEpisode,
				UpdatedAt:     anime.UpdatedAt,
			},
		},
	}

	jsonBytes, err := json.Marshal(fcmMessage)
	if err != nil {
		return false, err
	}

	url := fmt.Sprintf(env.FCMURLFormat, env.FirebaseProjectID)
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return false, err
	}

	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", oauthToken))
	r.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 404:
		return true, nil
	case 401:
		return false, fmt.Errorf("cannot authorize oauth credentials")
	}

	return false, nil
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
