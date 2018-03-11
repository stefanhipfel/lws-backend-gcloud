package database

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/appengine"

	"cloud.google.com/go/datastore"

	"github.com/stefanhipfel/lens-wide-shut/common"
	"github.com/stefanhipfel/lens-wide-shut/models"
)

// Ensure datastoreDB conforms to the BlockDatabase interface.
var _ models.ProjectDatabase = &DatastoreDB{}

func NewProjectDatastoreDB(r *http.Request) (models.ProjectDatabase, error) {
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

func (db *DatastoreDB) ListProjects() ([]*models.Project, error) {
	projects := make([]*models.Project, 0)
	q := datastore.NewQuery("Project").
		Order("created_on")

	keys, err := db.client.GetAll(*db.ctx, q, &projects)

	if err != nil {
		return nil, fmt.Errorf("datastoredb: could not list blogs: %v", err)
	}

	for i, k := range keys {
		projects[i].ID = k.ID
	}

	defer db.client.Close()
	return projects, nil
}

func (db *DatastoreDB) GetProject(id int64) (*models.Project, error) {
	project := &models.Project{}

	defer db.client.Close()

	k := datastore.IDKey("Project", id, nil)
	err := db.client.Get(*db.ctx, k, project)

	if err != nil {
		return project, err
	}
	project.ID = id

	return project, nil
}
