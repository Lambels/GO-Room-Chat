package models

type Account struct {
	// the account primary key
	ID       int64		`json:"id" db:"ID"`

	// email
	Email    string		`json:"email" db:"Email"`

	// username is used for message author
	Username string		`json:"username" db:"Username"`

	// account password encrypted with bcrypt package: https://pkg.go.dev/golang.org/x/crypto/bcrypt
	Password []byte		`json:"-" db:"Password"`
}

// Exists is a helper function used to check if the account exists
// when a query is made from the database
func (a *Account) Exists() bool {
	return !(a.Email == "")
}