package database

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/stefanhipfel/lens-wide-shut/common"
	"github.com/stefanhipfel/lens-wide-shut/models"
	"google.golang.org/api/iterator"
	"google.golang.org/appengine"
)

// Ensure datastoreDB conforms to the BlockDatabase interface.
var _ models.UserDatabase = &DatastoreDB{}

func NewUserDatastoreDB(r *http.Request) (models.UserDatabase, error) {
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

func (db *DatastoreDB) ListUsers() ([]*models.User, error) {
	users := make([]*models.User, 0)
	defer db.client.Close()

	q := datastore.NewQuery("User").
		Order("name")

	keys, err := db.client.GetAll(*db.ctx, q, &users)

	if err != nil {
		return nil, fmt.Errorf("datastoredb: could not list blogs: %v", err)
	}

	for i, k := range keys {
		users[i].ID = k.ID
	}

	return users, nil
}

func (db *DatastoreDB) AddUser(user *models.User) (int64, error) {
	k := datastore.IncompleteKey("User", nil)

	defer db.client.Close()

	k, err := db.client.Put(*db.ctx, k, user)

	if err != nil {
		return 0, err
	}

	return k.ID, nil
}

func (db *DatastoreDB) GetUser(id int64) (*models.User, error) {
	user := &models.User{}

	k := datastore.IDKey("User", id, nil)

	err := db.client.Get(*db.ctx, k, user)

	if err != nil {
		return user, err
	}

	user.ID = id

	return user, nil
}

func (db *DatastoreDB) GetUserByMail(mail string) (*models.User, error) {
	defer db.client.Close()

	q := datastore.NewQuery("User").Filter("email =", mail)

	it := db.client.Run(*db.ctx, q)

	user := &models.User{}

	for {
		if _, err := it.Next(&user); err == iterator.Done {
			break
		} else if err != nil {
			return user, err
		}
	}

	return user, nil
}

func (db *DatastoreDB) DeleteUser(i int64) error {
	return nil
}

func (db *DatastoreDB) UpdateUser(user *models.User) error {
	return nil
}
