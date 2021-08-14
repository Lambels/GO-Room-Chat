package models

type Room struct {
	ID		int64	`json:"id" db:"ID"`

	Name	string	`json:"name" db:"Name"`
	
	OwnerID	int64	`json:"owner" db:"Owner"`
}