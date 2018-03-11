package models

type Like struct {
	ID     int64    `json:"id, omitempty" datastore:"-"`
	BlogID int64    `json:"blog_id" datastore:"blog_id"`
	Likes  []string `json:"likes" datastore:"likes"`
}

// LikeDatabase provides thread-safe access to a database of books.
type LikeDatabase interface {
	// ListBooks returns a list of books, ordered by title.
	GetLike(id int64) (*Like, error)

	Vote(id int64, email string) error

	UnVote(id int64, email string) error

	// AddBook saves a given book, assigning it a new ID.
	AddLike(l *Like) (id int64, err error)

	// DeleteBook removes a given book by its ID.
	DeleteLike(id int64) error

	// UpdateBook updates the entry for a given book.
	UpdateLike(l *Like) error
}
