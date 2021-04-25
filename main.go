package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/hpcloud/tail"
	"github.com/joho/godotenv"

	BQ "nodeapps/logs/lib/bigquery"
	Structs "nodeapps/logs/structs"
)

// // main
// func main() {

// 	// construct `go version` command
// 	cmd := exec.Command("go", "run", "reader/reader.go")

// 	// configure `Stdout` and `Stderr`
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stdout

// 	// run command
// 	if err := cmd.Run(); err != nil {
// 		fmt.Println("Error:", err)
// 	}

// 	// construct `go version` command
// 	cmd2 := exec.Command("go", "run", "reader/reader.go")

// 	// configure `Stdout` and `Stderr`
// 	cmd2.Stdout = os.Stdout
// 	cmd2.Stderr = os.Stdout

// 	// run command
// 	if err := cmd2.Run(); err != nil {
// 		fmt.Println("Error:", err)
// 	}
// }

func DoneAsync() chan int {
	r := make(chan int)
	fmt.Println("Warming up ...")
	go func() {
		time.Sleep(3 * time.Second)

		t, err := tail.TailFile("./test.txt", tail.Config{Follow: true})
		if err == nil {
			for line := range t.Lines {
				fmt.Println(line.Text)
			}
		}

		r <- 1
		fmt.Println("Done ...")
	}()
	return r
}

func DoneAsync2() chan int {
	APP_ID := os.Getenv("NODE_APP_ID")

	r := make(chan int)
	go func() {

		t, err := tail.TailFile("/var/log/nginx/"+APP_ID+".access.log", tail.Config{Follow: true})
		if err == nil {
			for line := range t.Lines {
				// fmt.Println(line.Text)

				var dataJSON []string
				_ = json.Unmarshal([]byte(line.Text), &dataJSON)

				format := "2006-01-02 15:04:05"
				t, _ := time.Parse(time.RFC3339, dataJSON[1])
				timestamp := t.UTC().Local().Format(format)

				responseCode, err := strconv.Atoi(dataJSON[3])
				if err != nil {
					// handle error
					fmt.Println(err)
					os.Exit(2)
				}

				logEntry := Structs.LogEntry{}

				logEntry.Timestamp = timestamp
				logEntry.App = APP_ID
				logEntry.Request = dataJSON[2]
				logEntry.Code = responseCode

				logEntryJSON, err := json.Marshal(logEntry)
				if err != nil {
					fmt.Println(err)
				} else {
					BQ.InsertAccessLog(logEntryJSON)
				}

			}
		}

		r <- 1
		fmt.Println("Done ...")
	}()
	return r
}

func main() {
	projectFlag := flag.String("project", "", "the project name to identify it later on")

	// Load env vars.
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file %v", err)
	}

	flag.Parse()
	fmt.Println("word:", *projectFlag)

	fmt.Println("Let's start ...")
	// val := DoneAsync()
	val2 := DoneAsync2()
	fmt.Println("Done are running ...")
	// fmt.Println(<-val)
	fmt.Println(<-val2)

}
