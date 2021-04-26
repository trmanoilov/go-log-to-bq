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

/*
Parsing NGINX access logs.
*/
func ParseAccessLogs(APP_ID string) chan int {
	r := make(chan int)
	go func() {

		// Tailing NGINX logs with a follow option.
		t, err := tail.TailFile("/var/log/nginx/"+APP_ID+".access.log", tail.Config{Follow: true})
		if err == nil {
			for line := range t.Lines {

				// Preparing a new object to hold parsed JSON data and parse the data.
				var dataJSON []string
				_ = json.Unmarshal([]byte(line.Text), &dataJSON)

				// Prepare format and convert the log timestamp to the selected format.
				format := "2006-01-02 15:04:05"
				t, _ := time.Parse(time.RFC3339, dataJSON[1])
				timestamp := t.UTC().Local().Format(format)

				// Parse the server response code from string to int.
				responseCode, err := strconv.Atoi(dataJSON[3])
				if err != nil {
					os.Exit(2)
				}

				// Prepare & assign the generated values.
				logEntry := Structs.AccessLogEntry{}

				request := dataJSON[2]

				logEntry.Timestamp = timestamp
				logEntry.App = APP_ID
				logEntry.Request = request
				logEntry.Code = responseCode

				// Convert struct to JSON to be able to pass it as argument to the next fuction.
				logEntryJSON, err := json.Marshal(logEntry)
				if err != nil {
					// fmt.Println(err)
				} else {
					BQ.InsertAccessLog(logEntryJSON)
				}

				// If the response code does not qualify as an "OK" one -
				// send a message to a dedicated Slack channel.
				if !isOkCode(responseCode) {
					Slack.SendSlackNotification("Received a `" + strconv.Itoa(responseCode) + "` response code on: `" + APP_ID + "` app")
				}

			}
		}

		r <- 1
		// Signal the function is done doing what it's doing.
		fmt.Println("Done ...")
	}()
	return r
}

/*
Verify whether the response code is part of the predefined set of allowed ones.
*/
func isOkCode(responseCode int) bool {
	okResponseCodes := []int{200, 301}
	for _, code := range okResponseCodes {
		if responseCode == code {
			return true
		}
	}
	return false
}
