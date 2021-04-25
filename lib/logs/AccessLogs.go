package logs

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/hpcloud/tail"

	BQ "nodeapps/logs/lib/bigquery"
	Slack "nodeapps/logs/lib/slack"
	Structs "nodeapps/logs/structs"
)

func ParseAccessLogs(APP_ID string) chan int {
	r := make(chan int)
	go func() {

		t, err := tail.TailFile("/var/log/nginx/"+APP_ID+".access.log", tail.Config{Follow: true})
		if err == nil {
			for line := range t.Lines {

				var dataJSON []string
				_ = json.Unmarshal([]byte(line.Text), &dataJSON)

				format := "2006-01-02 15:04:05"
				t, _ := time.Parse(time.RFC3339, dataJSON[1])
				timestamp := t.UTC().Local().Format(format)

				responseCode, err := strconv.Atoi(dataJSON[3])
				if err != nil {
					os.Exit(2)
				}

				request := dataJSON[2]

				logEntry := Structs.AccessLogEntry{}

				logEntry.Timestamp = timestamp
				logEntry.App = APP_ID
				logEntry.Request = request
				logEntry.Code = responseCode

				logEntryJSON, err := json.Marshal(logEntry)
				if err != nil {
					// fmt.Println(err)
				} else if request != "" {
					BQ.InsertAccessLog(logEntryJSON)
				}

				if !isOkCode(responseCode) {
					Slack.SendSlackNotification("Received a `" + strconv.Itoa(responseCode) + "` response code on: `" + APP_ID + "` app")
				}

			}
		}

		r <- 1
		fmt.Println("Done ...")
	}()
	return r
}

func isOkCode(responseCode int) bool {
	okResponseCodes := []int{200, 301}
	for _, code := range okResponseCodes {
		if responseCode == code {
			return true
		}
	}
	return false
}
