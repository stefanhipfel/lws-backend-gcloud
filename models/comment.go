package models

import "time"

type Comment struct {
	ID        int64     `json:"id, omitempty" datastore:"-"`
	Name      string    `json:"name" datastore:"name"`
	Website   string    `json:"website" datastore:"website"`
	Message   string    `json:"message" datastore:"message"`
	CreatedOn time.Time `json:"created_on,omitempty" datastore:"created_on"`
	CreatedBy string    `json:"user" datastore:"user"`
	BlogID    int64     `json:"blog_id" datastore:"blog_id"`
	Replies   []int64   `json:"replies" datastore:"replies"`
}

func (post *Comment) FormattedDate() string {
	return post.CreatedOn.Format(time.RFC822)
}

// BookDatabase provides thread-safe access to a database of books.
type CommentDatabase interface {
	// ListBooks returns a list of books, ordered by title.
	ListComments() ([]*Comment, error)

	// GetBook retrieves a book by its ID.
	GetComment(id int64) (*Comment, error)

	// GetBook retrieves a book by its ID.
	GetCommentsByBlogId(id int64) ([]*Comment, error)

	// AddBook saves a given book, assigning it a new ID.
	AddComment(b *Comment) (id int64, err error)

	// DeleteBook removes a given book by its ID.
	DeleteComment(id int64) error

	// UpdateBook updates the entry for a given book.
	UpdateComment(b *Comment) error
}
