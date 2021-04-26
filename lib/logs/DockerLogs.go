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

func ParseDockerLogs(APP_ID string) chan int {
	r := make(chan int)
	go func() {

		t, err := tail.TailFile("/var/log/docker/"+APP_ID+".log", tail.Config{Follow: true})
		if err == nil {
			for line := range t.Lines {

				stringified := fmt.Sprintf("%v", line)
				fields := strings.Fields(stringified)

				currentTime := time.Now()

				logEntry := Structs.UsageLogEntry{}

				logEntry.Timestamp = currentTime.UTC().Local().Format("2006-01-02 15:04:05")
				logEntry.App = APP_ID

				cpu, _ := strconv.ParseFloat(strings.Replace(fields[2], "%", "", -1), 64)
				logEntry.CPU = cpu

				memory, _ := strconv.ParseFloat(strings.Replace(fields[6], "%", "", -1), 64)
				logEntry.Memory = memory

				logEntryJSON, err := json.Marshal(logEntry)
				if err != nil {
					// fmt.Println(err)
				} else {
					BQ.InsertDockerLog(logEntryJSON)
				}

			}
		}

		r <- 1
		fmt.Println("Done ...")
	}()
	return r
}

func GenerateDockerLogs(APP_ID string) chan string {
	r := make(chan string)
	go func() {

		out, err := exec.Command("./dockerlog", string(APP_ID), "&").Output()

		if err != nil {
			log.Fatal(err)
		}

		r <- string(out)

		fmt.Println("Done ...")
	}()
	return r
}
