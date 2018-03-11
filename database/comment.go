package database

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/datastore"

	"github.com/stefanhipfel/lens-wide-shut/common"
	"github.com/stefanhipfel/lens-wide-shut/models"
	"google.golang.org/appengine"
)

// Ensure datastoreDB conforms to the BlockDatabase interface.
var _ models.CommentDatabase = &DatastoreDB{}

func NewCommentDatastoreDB(r *http.Request) (models.CommentDatabase, error) {
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

func (db *DatastoreDB) GetComment(i int64) (*models.Comment, error) {
	c := new(models.Comment)

	defer db.client.Close()

	k := datastore.IDKey("Blog", i, nil)
	err := db.client.Get(*db.ctx, k, c)

	if err != nil {
		return c, err
	}

	return &models.Comment{}, nil
}

func (db *DatastoreDB) AddComment(c *models.Comment) (int64, error) {
	k := datastore.IncompleteKey("Comment", nil)

	c.CreatedOn = time.Now()

	k, err := db.client.Put(*db.ctx, k, c)

	defer db.client.Close()

	if err != nil {
		return 0, err
	}

	return k.ID, nil
}

func (db *DatastoreDB) ListComments() ([]*models.Comment, error) {
	c := make([]*models.Comment, 0)

	q := datastore.NewQuery("Comment").
		Order("created_on")

	keys, err := db.client.GetAll(*db.ctx, q, &c)

	if err != nil {
		return nil, fmt.Errorf("datastoredb: could not list comments: %v", err)
	}

	for i, k := range keys {
		c[i].ID = k.ID
	}

	defer db.client.Close()

	return c, nil
}

func (db *DatastoreDB) GetCommentsByBlogId(i int64) ([]*models.Comment, error) {
	c := make([]*models.Comment, 0)

	q := datastore.NewQuery("Comment").Filter("blog_id =", i).Order("created_on")

	keys, err := db.client.GetAll(*db.ctx, q, &c)

	if err != nil {
		return nil, fmt.Errorf("datastoredb: could not list comments: %v", err)
	}

	for i, k := range keys {
		c[i].ID = k.ID
	}

	defer db.client.Close()

	return c, nil
}

func (db *DatastoreDB) DeleteComment(i int64) error {
	return nil
}

func (db *DatastoreDB) UpdateComment(*models.Comment) error {
	return nil
}
