package models

import (
	"time"
)

type Project struct {
	ID        int64     `json:"id, omitempty" datastore:"-"`
	Title     string    `json:"title" datastore:"title"`
	Text      string    `json:"text" datastore:"text"`
	Images    *[]Image  `json:"images" datastore:"-"`
	CoverURL  string    `json:"coverurl" datastore:"cover_url"`
	CreatedOn time.Time `json:"createdon" datastore:"created_on"`
	Category  string    `json:"category" datastore:"category"`
}

// ProjectDatabase provides thread-safe access to a database of Projects.
type ProjectDatabase interface {
	// ListProjects returns a list of projects, ordered by created date.
	ListProjects() ([]*Project, error)

	// GetProject retrieves a project by its ID.
	GetProject(id int64) (*Project, error)
}
