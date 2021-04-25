package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/joho/godotenv"

	Logs "nodeapps/logs/lib/logs"
)

func main() {
	projectFlag := flag.String("project", "", "the project name to identify it later on")

	// Load env vars.
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file %v", err)
	}

	flag.Parse()

	fmt.Println("Let's start ...")
	val2 := Logs.GenerateDockerLogs(*projectFlag)
	val := Logs.ParseDockerLogs(*projectFlag)
	val3 := Logs.ParseAccessLogs(*projectFlag)

	fmt.Println("Done are running ...")
	fmt.Println(<-val2)
	fmt.Println(<-val3)
	fmt.Println(<-val)
}
