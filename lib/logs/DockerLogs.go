package logs

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/hpcloud/tail"
)

func ParseDockerLogs(APP_ID string) chan int {
	r := make(chan int)
	go func() {

		t, err := tail.TailFile("/var/log/docker/"+APP_ID+".log", tail.Config{Follow: true})
		if err == nil {
			for line := range t.Lines {

				stringified := fmt.Sprintf("%v", line)
				words := strings.Fields(stringified)

				fmt.Println(words[1], words[2], words[6], len(words)) // [one two three four] 4

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
