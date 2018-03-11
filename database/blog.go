package database

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/appengine"

	"cloud.google.com/go/datastore"

	"github.com/stefanhipfel/lens-wide-shut/common"
	"github.com/stefanhipfel/lens-wide-shut/models"
)

// Ensure datastoreDB conforms to the BlockDatabase interface.
var _ models.BlogDatabase = &DatastoreDB{}

func NewBlogDatastoreDB(r *http.Request) (models.BlogDatabase, error) {
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

func (db *DatastoreDB) AddBlog(blogPost *models.BlogPost) (int64, error) {
	k := datastore.IncompleteKey("Blog", nil)

	blogPost.CreatedOn = time.Now()

	k, err := db.client.Put(*db.ctx, k, blogPost)

	defer db.client.Close()

	if err != nil {
		return 0, err
	}

	return k.ID, nil
}

func (db *DatastoreDB) ListBlogsCreatedBy(createdBy string) ([]*models.BlogPost, error) {
	blogPosts := make([]*models.BlogPost, 0)

	return blogPosts, nil
}

func (db *DatastoreDB) GetBlog(id int64) (*models.BlogPost, error) {
	blogPost := &models.BlogPost{}

	defer db.client.Close()

	k := datastore.IDKey("Blog", id, nil)
	err := db.client.Get(*db.ctx, k, blogPost)

	if err != nil {
		return blogPost, err
	}
	blogPost.ID = id

	return blogPost, nil
}

func (db *DatastoreDB) ListBlogs() ([]*models.BlogPost, error) {
	blogPosts := make([]*models.BlogPost, 0)

	q := datastore.NewQuery("Blog").
		Order("created_on")

	keys, err := db.client.GetAll(*db.ctx, q, &blogPosts)

	if err != nil {
		return nil, fmt.Errorf("datastoredb: could not list blogs: %v", err)
	}

	for i, k := range keys {
		blogPosts[i].ID = k.ID
	}

	defer db.client.Close()

	return blogPosts, nil
}

func (db *DatastoreDB) UpdateBlog(blog *models.BlogPost) error {
	return nil
}

func (db *DatastoreDB) DeleteBlog(i int64) error {
	return nil
}
