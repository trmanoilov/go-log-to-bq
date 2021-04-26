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

/*
Send access logs data to BigQuery to be able to parse and show in Grafana.
*/
func InsertAccessLog(logEntryJSON []byte) error {
	// Create contenxt and parse environment variables.
	ctx := context.Background()
	datasetID := os.Getenv("GOOGLE_CLOUD_DATASET_ID")
	tableID := os.Getenv("GOOGLE_CLOUD_ACCESS_TABLE")
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	credentialsJSON := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	// Then utilize all of the above to instantiate a BigQuery client.
	client, err := bigquery.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsJSON))

	// Obsviously if things don't work - we shall not proceed.
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

	// Insert to dataset.
	inserter := client.Dataset(datasetID).Table(tableID).Inserter()
	if err := inserter.Put(ctx, tickerData); err != nil {
		log.Fatal(err)
		return err
	}
	return nil

}

/*
Send usage logs data to BigQuery to be able to parse and show in Grafana.
*/
func InsertDockerLog(logEntryJSON []byte) error {
	// Create contenxt and parse environment variables.
	ctx := context.Background()
	datasetID := os.Getenv("GOOGLE_CLOUD_DATASET_ID")
	tableID := os.Getenv("GOOGLE_CLOUD_USAGE_TABLE")
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	credentialsJSON := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	// Then utilize all of the above to instantiate a BigQuery client.
	client, err := bigquery.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsJSON))

	// Obsviously if things don't work - we shall not proceed.
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	// Prepare proper structure and parse JSON.
	tickerData := Structs.UsageLogEntry{}
	decodeError := json.Unmarshal(logEntryJSON, &tickerData)
	if decodeError != nil {
		log.Fatal(decodeError)
	}

	// Insert to dataset.
	inserter := client.Dataset(datasetID).Table(tableID).Inserter()
	if err := inserter.Put(ctx, tickerData); err != nil {
		log.Fatal(err)
		return err
	}
	return nil

}
