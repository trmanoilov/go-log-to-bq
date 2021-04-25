package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/hpcloud/tail"
	"github.com/joho/godotenv"
)

func DoneAsync() chan int {
	r := make(chan int)
	go func() {

		// t, err := tail.TailFile("/var/log/nginx/"+APP_ID+".access.log", tail.Config{Follow: true})
		t, err := tail.TailFile("./test.txt", tail.Config{Follow: true})
		if err == nil {
			for line := range t.Lines {

				fmt.Println(line)

				// var dataJSON []string
				// _ = json.Unmarshal([]byte(line.Text), &dataJSON)

				// format := "2006-01-02 15:04:05"
				// t, _ := time.Parse(time.RFC3339, dataJSON[1])
				// timestamp := t.UTC().Local().Format(format)

				// responseCode, err := strconv.Atoi(dataJSON[3])
				// if err != nil {
				// 	// handle error
				// 	// fmt.Println(err)
				// 	os.Exit(2)
				// }

				// logEntry := Structs.LogEntry{}

				// logEntry.Timestamp = timestamp
				// logEntry.App = APP_ID
				// logEntry.Request = dataJSON[2]
				// logEntry.Code = responseCode

				// logEntryJSON, err := json.Marshal(logEntry)
				// if err != nil {
				// 	// fmt.Println(err)
				// } else {
				// 	BQ.InsertAccessLog(logEntryJSON)
				// }

				// if !isOkCode(responseCode) {
				// 	Slack.SendSlackNotification("Received a `" + strconv.Itoa(responseCode) + "` response code on: `" + APP_ID + "` app")
				// 	// if slackErr != nil {
				// 	// log.Fatal(slackErr)
				// 	// }
				// }

			}
		}

		r <- 1
		fmt.Println("Done ...")
	}()
	return r
}

func main() {
	// projectFlag := flag.String("project", "", "the project name to identify it later on")

	// Load env vars.
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file %v", err)
	}

	flag.Parse()

	fmt.Println("Let's start ...")
	val := DoneAsync()
	// val2 := AL.ParseAccessLogs(*projectFlag)
	fmt.Println("Done are running ...")
	fmt.Println(<-val)
	// fmt.Println(<-val2)

}
