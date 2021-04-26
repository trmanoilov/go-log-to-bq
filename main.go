package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/joho/godotenv"

	Logs "nodeapps/logs/lib/logs"
)

/*
Main app logic.
*/
func main() {
	// Load env vars.
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file %v", err)
	}

	// Defining and parsing command line arguments.
	// More specifically "project" to map to app container.
	projectFlag := flag.String("project", "", "the project name to identify it later on")
	flag.Parse()

	// Invoking async channels to generate & process data.
	// The order is surprisingly important.
	dockerLogsGenerator := Logs.GenerateDockerLogs(*projectFlag)
	dockerUsageLogsParser := Logs.ParseDockerLogs(*projectFlag)
	accessLogsParser := Logs.ParseAccessLogs(*projectFlag)

	// If for whatever reason the channels have the data flow
	// stopped - then the value is returned and printed out.
	// Might also be forwarded to some sort of error log for better
	// overview in future.
	fmt.Println(<-dockerLogsGenerator)
	fmt.Println(<-accessLogsParser)
	fmt.Println(<-dockerUsageLogsParser)
}
