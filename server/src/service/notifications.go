package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"strconv"

	"net/http"
	"os"
	"server/src/models"
)


func sendMentionNotif(newMentionNotification models.MentionNotificationRequest, token string) error {

	if os.Getenv("ENVIROMENT") == "test" {
		slog.Info("Notification sent to ", "user_id", newMentionNotification.UserId)
		return nil
	}

	url := "http://" + os.Getenv("NOTIF_HOST") + "/notification/mention"

	marshalledData, _ := json.Marshal(newMentionNotification)

	req, err := http.NewRequest("POST", url, bytes.NewReader(marshalledData))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	if err != nil {
		return errors.New("error creating request")
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return errors.New("error sending request, " + err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("error sending request, status code: "  + strconv.Itoa(resp.StatusCode))
	}

	slog.Info("Notification sent to ", "user_id", newMentionNotification.UserId)

	return nil
}