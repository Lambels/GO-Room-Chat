package models

const (
	TYPE_NEW_ROOM = iota
	TYPE_NEW_CONNECTION
	TYPE_MESSAGE
)

type Message struct {
	// Type of message used by the reciever to know how to process the message
	// Ex:
	//
	// TYPE_NEW_ROOM
	Type 		int								`json:"type"`

	// Sender is the account which send the message
	Sender		Account							`json:"account"`

	// The data is map and its processed by the reciever in different ways
	// depending on what type the message is
	Data		map[string]interface{}			`json:"data"`
}