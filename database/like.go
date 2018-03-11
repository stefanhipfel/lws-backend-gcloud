package database

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/stefanhipfel/lens-wide-shut/common"
	"github.com/stefanhipfel/lens-wide-shut/models"
	"google.golang.org/appengine"
)

// Ensure datastoreDB conforms to the BlockDatabase interface.
var _ models.LikeDatabase = &DatastoreDB{}

func NewLikeDatastoreDB(r *http.Request) (models.LikeDatabase, error) {
	client, err := common.CreateDatastoreClient(r)
	var ctx context.Context
	ctx = appengine.NewContext(r)
	// Verify that we can communicate and authenticate with the datastore service.
	t, err := client.NewTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("datastoredb: could not connect: %v", err)
	}
	if err := t.Rollback(); err != nil {
		return nil, fmt.Errorf("datastoredb: could not connect: %v", err)
	}
	return &DatastoreDB{
		client: client,
		ctx:    &ctx,
	}, nil
}

func (db *DatastoreDB) AddLike(l *models.Like) (int64, error) {
	k := datastore.IncompleteKey("Like", nil)

	k, err := db.client.Put(*db.ctx, k, l)

	defer db.client.Close()

	if err != nil {
		return 0, err
	}

	return k.ID, nil
}

func (db *DatastoreDB) GetLike(i int64) (*models.Like, error) {
	like := make([]*models.Like, 0)
	q := datastore.NewQuery("Like").Filter("blog_id =", i)

	k, err := db.client.GetAll(*db.ctx, q, &like)

	if err != nil {
		return &models.Like{}, fmt.Errorf("datastoredb: could not get like: %v", err)
	}
	like[0].ID = k[0].ID

	return like[0], nil
}

func (db *DatastoreDB) DeleteLike(i int64) error { return nil }

func (db *DatastoreDB) UpdateLike(l *models.Like) error { return nil }

func (db *DatastoreDB) Vote(i int64, e string) error {
	like := make([]*models.Like, 0)
	q := datastore.NewQuery("Like").Filter("blog_id =", i)
	k, err := db.client.GetAll(*db.ctx, q, &like)

	if err != nil {
		return fmt.Errorf("datastoredb: could not get like: %v", err)
	}
	like[0].ID = k[0].ID

	b := contains(like[0].Likes, e)

	if !b {
		like[0].Likes = append(like[0].Likes, e)
		if _, err := db.client.Put(*db.ctx, k[0], like[0]); err != nil {
			return fmt.Errorf("datastoredb: could not update like: %v", err)
		}
		return nil
	}

	return nil
}

func (db *DatastoreDB) UnVote(i int64, e string) error {
	return nil
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
