package main

import (
	"fmt"
	"os"

	"github.com/hpcloud/tail"
)

// main
func main() {

	err := os.Truncate("./test.txt", 0)
	if err != nil {
		fmt.Println(err)
	}

	t, err := tail.TailFile("./test.txt", tail.Config{Follow: true})
	if err == nil {
		for line := range t.Lines {
			fmt.Println(line.Text)
		}
	}

}
