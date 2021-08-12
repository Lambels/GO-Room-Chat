package models

type Room struct {
	ID		int64	`json:"id" db:"ID"`

	Name	string	`json:"name" db:"Name"`
	
	Owner	Account	`json:"owner" db:"Owner"`
}