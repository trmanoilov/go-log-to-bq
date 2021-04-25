package BQ

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/option"

	Structs "nodeapps/logs/structs"
)

func InsertAccessLog(logEntryJSON []byte) error {
	ctx := context.Background()
	datasetID := os.Getenv("GOOGLE_CLOUD_DATASET_ID")
	tableID := os.Getenv("GOOGLE_CLOUD_TICKER_TABLE")
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	credentialsJSON := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	client, err := bigquery.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsJSON))
	fmt.Printf("bigquery.credentialsJSON: %s", credentialsJSON)
	fmt.Printf("bigquery.NewClient: %v", client)

	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	// Prepare proper structure and parse JSON.
	tickerData := Structs.AccessLogEntry{}
	decodeError := json.Unmarshal(logEntryJSON, &tickerData)
	if decodeError != nil {
		log.Fatal(decodeError)
	}
	fmt.Println(tickerData)

	// Insert to dataset.
	inserter := client.Dataset(datasetID).Table(tableID).Inserter()
	if err := inserter.Put(ctx, tickerData); err != nil {
		log.Fatal(err)
		return err
	}
	return nil

}
