package models

type User struct {
	ID            int64  `json:"id, omitempty" datastore:"-"`
	Sub           string `json:"sub" datastore:"sub"`
	Name          string `json:"name" datastore:"name"`
	GivenName     string `json:"given_name" datastore:"given_name"`
	FamilyName    string `json:"family_name" datastore:"family_name"`
	Profile       string `json:"profile" datastore:"profile"`
	Picture       string `json:"picture" datastore:"picture"`
	Email         string `json:"email" datastore:"email"`
	EmailVerified bool   `json:"email_verified" datastore:"email_verified"`
	Gender        string `json:"gender" datastore:"gender"`
}

// BookDatabase provides thread-safe access to a database of books.
type UserDatabase interface {
	// ListBooks returns a list of books, ordered by title.
	ListUsers() ([]*User, error)

	// GetBook retrieves a book by its ID.
	GetUser(id int64) (*User, error)

	// GetBook retrieves a book by its ID.
	GetUserByMail(mail string) (*User, error)

	// AddBook saves a given book, assigning it a new ID.
	AddUser(b *User) (id int64, err error)

	// DeleteBook removes a given book by its ID.
	DeleteUser(id int64) error

	// UpdateBook updates the entry for a given book.
	UpdateUser(b *User) error
}
