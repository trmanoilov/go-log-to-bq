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

	// Prepare request body JSON, holding the message data.
	slackBody, _ := json.Marshal(Structs.SlackRequestBody{Text: msg})

	// Prepare an XHR request to the predefined webhook URL.
	req, err := http.NewRequest(http.MethodPost, WEBHOOK_URL, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	// Adjust request headers & set timeout. 10 sec comes by default from Slack.
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// If all is good - do nothing, otherwise alert.
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		log.Fatal("Non-ok response returned from Slack")
	}
	return nil
}
