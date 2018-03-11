package common

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/appengine"

	"cloud.google.com/go/datastore"
)

var client *datastore.Client

func GetDatastoreClient() *datastore.Client {
	return client
}

func CreateDatastoreClient(r *http.Request) (*datastore.Client, error) {
	var err error
	projID := os.Getenv("DATASTORE_PROJECT_ID")
	if projID == "" {
		log.Fatal(`You need to set the environment variable "DATASTORE_PROJECT_ID"`)
	}
	ctx := appengine.NewContext(r)
	client, err = datastore.NewClient(ctx, projID)

	if err != nil {
		fmt.Println(err.Error())
		return client, err
	}

	t, err := client.NewTransaction(ctx)
	if err != nil {
		fmt.Errorf("datastoredb: could not connect: %v", err)
		return client, err
	}
	if err := t.Rollback(); err != nil {
		fmt.Errorf("datastoredb: could not connect: %v", err)
		return client, err
	}
	return client, nil
}
