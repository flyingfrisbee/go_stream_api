package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	env "go_stream_api/environment"
	db "go_stream_api/repository/database"
	"net/http"
	"strings"
)

type codesPayload struct {
	Codes string `json:"codes"`
}

func SendCodesNotificationToUsers(newCodes []string) error {
	// Reason for this is because i want user to copy the codes to the clipboard after clicking the notification
	concatenatedNewCodes := strings.Join(newCodes, ";")

	usersToken, err := db.Conn.Pg.GetAllUsersToken()
	if err != nil {
		return err
	}

	if len(usersToken) < 1 {
		return nil
	}

	for _, userToken := range usersToken {
		go sendCodesNotification(userToken, concatenatedNewCodes)
	}

	return nil
}

func sendCodesNotification(userToken, codes string) error {
	fcmMessage := fcmMessage{
		Message: message{
			Token: userToken,
			Data: codesPayload{
				Codes: codes,
			},
		},
	}

	jsonBytes, err := json.Marshal(fcmMessage)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(env.FCMURLFormat, env.FirebaseProjectID)
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}

	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", oauthToken))
	r.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
