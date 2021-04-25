package slack

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	Structs "nodeapps/logs/structs"
)

// SendSlackNotification will post to an 'Incoming Webook' url setup in Slack Apps. It accepts
// some text and the slack channel is saved within Slack.
func SendSlackNotification(msg string) error {
	WEBHOOK_URL := os.Getenv("WEBHOOK_URL")

	slackBody, _ := json.Marshal(Structs.SlackRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, WEBHOOK_URL, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		log.Fatal("Non-ok response returned from Slack")
	}
	return nil
}
