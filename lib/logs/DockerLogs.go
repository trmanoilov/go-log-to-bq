package logs

import (
	"encoding/json"
	"fmt"
	"log"
	BQ "nodeapps/logs/lib/bigquery"
	Structs "nodeapps/logs/structs"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/hpcloud/tail"
)

/*
Parsing custom generated logs.
*/
func ParseDockerLogs(APP_ID string) chan int {
	r := make(chan int)
	go func() {

		// Tailing out custom logs with a follow option.
		t, err := tail.TailFile("/var/log/docker/"+APP_ID+".log", tail.Config{Follow: true})
		if err == nil {
			for line := range t.Lines {

				// Converting to string because the tail module does not return string :man-shrugging:
				stringified := fmt.Sprintf("%v", line)
				// Then split into fields.
				fields := strings.Fields(stringified)

				// Getting current timestamp.
				currentTime := time.Now()

				// Preparing usage log entry structure to ensure proper import.
				logEntry := Structs.UsageLogEntry{}

				// Parse the timestamp to adequate date format.
				logEntry.Timestamp = currentTime.UTC().Local().Format("2006-01-02 15:04:05")

				// Map to proper app that this process is assigned to.
				logEntry.App = APP_ID

				// Convert the CPU & Memory usage from string to float,
				// stripping the percentage and assigning to the struct.
				cpu, _ := strconv.ParseFloat(strings.Replace(fields[2], "%", "", -1), 64)
				logEntry.CPU = cpu

				memory, _ := strconv.ParseFloat(strings.Replace(fields[6], "%", "", -1), 64)
				logEntry.Memory = memory

				// Convert struct to JSON to be able to pass it as argument to the next fuction.
				logEntryJSON, err := json.Marshal(logEntry)
				if err != nil {
					// fmt.Println(err)
				} else {
					BQ.InsertDockerLog(logEntryJSON)
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
Log generator.

Since I couldn't use the built in logs and the already made Grafana
plugins that will output the data.. custom logs were needed.
*/
func GenerateDockerLogs(APP_ID string) chan string {
	r := make(chan string)
	go func() {

		// Invoke CMD command, that's executing a bash script, passing the
		// app ID as argument, as well as "&" to keep the process running.
		out, err := exec.Command("./dockerlog", string(APP_ID), "&").Output()

		if err != nil {
			log.Fatal(err)
		}

		r <- string(out)

		// Signal the function is done doing what it's doing.
		fmt.Println("Done ...")
	}()
	return r
}
