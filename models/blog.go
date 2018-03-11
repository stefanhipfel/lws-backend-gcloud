package models

import (
	"html/template"
	"time"
)

type BlogPost struct {
	ID        int64         `json:"id, omitempty" datastore:"-"`
	Title     string        `json:"title" datastore:"title"`
	Text      template.HTML `json:"text" datastore:"text"`
	ImageUrl  string        `json:"imageUrl" datastore:"image_url"`
	CreatedOn time.Time     `json:"createdon" datastore:"created_on"`
	CreatedBy int64         `json:"created_by" datastore:"created_by"`
	User      *User          `json:"user, omitempty" datastore:"-"`
	Category  string        `json:"category" datastore:"category"`
}

func (post *BlogPost) FormattedDate() string {
	return post.CreatedOn.Format(time.RFC822)
}

// BookDatabase provides thread-safe access to a database of books.
type BlogDatabase interface {
	// ListBooks returns a list of books, ordered by title.
	ListBlogs() ([]*BlogPost, error)

	// ListBooksCreatedBy returns a list of books, ordered by title, filtered by
	// the user who created the book entry.
	ListBlogsCreatedBy(userID string) ([]*BlogPost, error)

	// GetBook retrieves a book by its ID.
	GetBlog(id int64) (*BlogPost, error)

	// AddBook saves a given book, assigning it a new ID.
	AddBlog(b *BlogPost) (id int64, err error)

	// DeleteBook removes a given book by its ID.
	DeleteBlog(id int64) error

	// UpdateBook updates the entry for a given book.
	UpdateBlog(b *BlogPost) error
}
