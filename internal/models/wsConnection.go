package models

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)


const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512
)

type WsConnection struct {
	// The websocket connection belonging to the user
	// used to communicate with the user
	Conn 	*websocket.Conn

	// The account with which the connestion is associatied with
	Acc  	Account

	// The channel from which messages are pumped to the client
	Send	chan Message	
}

// Listen registers the conn to the set register channel and starts the write and read pump respectively
func (conn WsConnection) Listen(register, unregister chan<- WsConnection, broadcast chan<- Message) {
	register <- conn
	go conn.writePump()
	go conn.readPump(unregister, broadcast)
}

// readPump pumps messages from the connection to the broadcast channel.
func (conn WsConnection) readPump(unregister chan<- WsConnection, broadcast chan<- Message) {
	defer func() {
		unregister <- conn
		conn.Conn.Close()
	}()

	conn.Conn.SetReadLimit(maxMessageSize)
	conn.Conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.Conn.SetPongHandler(func(appData string) (error) { conn.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		var msg Message
		err := conn.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("readPump: unexpected error: %v", err)
			}
			break
		}

		broadcast <- msg
	}
}

// writePump pumps messages from the send channel to the wbebsocket
func (conn WsConnection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func(){
		ticker.Stop()
		conn.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-conn.Send:
			conn.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			// User kicked or hub closed?
			if !ok {
				conn.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := conn.Conn.WriteJSON(msg)
			if err != nil {
				return
			}

		case <- ticker.C:
			conn.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

type connMap map[string]WsConnection

type CommunicationChannels struct {
	Register	chan<- WsConnection
	Unregister 	chan<- WsConnection
	Broadcast	chan<- Message
}

// Creates a new ConnMap and starts listening on
// the returned channels
//
// The returned channels are register, unregister, broadcast
func NewConnStore() (CommunicationChannels) {
	connMap := connMap{}

	register := make(chan WsConnection, 1)
	unregister := make(chan WsConnection, 1)
	broadcast := make(chan Message, 1)

	go connMap.listen(register, unregister, broadcast)
	
	return CommunicationChannels{
		Register: register,
		Unregister: unregister,
		Broadcast: broadcast,
	}
}

// Broadcasts the message to the connections and handles any error by closing the connection
func (c connMap) broadcast(msg Message) {
	for e, conn := range c {
		if e == msg.Sender.Email {
			continue
		}
		select {
		case conn.Send <- msg:

		default:
			close(conn.Send)
			delete(c, e)
		}
	}
}

// Listens for events:
//
// Register, unregister, broadcast
func (c connMap) listen(register, unregister <-chan WsConnection, broadcast <-chan Message) {
	for {
		select {
		case conn := <-register:
			c[conn.Acc.Email] = conn

		case conn := <-unregister:
			delete(c, conn.Acc.Email)
			close(conn.Send)

		case msg := <-broadcast:
			c.broadcast(msg)
		}
	}
}