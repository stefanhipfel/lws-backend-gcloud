package database

import (
	"context"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
)

// datastoreDB persists books to Cloud Datastore.
// https://cloud.google.com/datastore/docs/concepts/overview
type DatastoreDB struct {
	client *datastore.Client
	ctx    *context.Context
}

type StorageDB struct {
	bucket *storage.BucketHandle
	ctx    context.Context
}
