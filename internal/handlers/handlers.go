package handlers

import (
	"GO-Group-Chat/internal/config"
	"GO-Group-Chat/internal/helpers"
	"GO-Group-Chat/internal/middleware"
	"GO-Group-Chat/internal/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

var (
	app *config.AppConfig

	upgrader = websocket.Upgrader {
		WriteBufferSize: 1024,
		ReadBufferSize: 1024,
	}

	indexChannels	models.CommunicationChannels
)

func init() {
	indexChannels = models.NewConnStore()
}

func InitHandlers(a *config.AppConfig) {
	app = a
}

// Handles "/" and renders to different templates depending if the user
// sending the request is auth
//
// and accepts only the GET verb
func Index(w http.ResponseWriter, r *http.Request) {
	if !helpers.IsAuth(r) {
		err := helpers.RenderTemplate(w, "index.html", nil)
		if err != nil {
			http.Error(w, "Couldn't render template", http.StatusInternalServerError)
		}
		return
	}

	u, err := helpers.GetUser(r)
	if err != nil {
		return
	}

	err = helpers.RenderTemplate(w, "index.html", u)
	if err != nil {
		http.Error(w, "Couldn't render template", http.StatusInternalServerError)
	}
}

// Handels "/index/ws" and connects to the index room used to feed users on the index page updates such as
// new rooms
//
// IndexWS needs to be protected by onlyAuth middleware
func IndexWS(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	// browser doesent support ws ?
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	user := r.Context().Value(middleware.AccountKey{})

	conn := models.WsConnection {
		Conn: 	ws,
		Acc: 	user.(models.Account),
		Send: 	make(chan models.Message),
	}

	go conn.Listen(indexChannels.Register, indexChannels.Unregister, indexChannels.Broadcast)
}

// Signup handles "/signup" and accepts POST and GET verbs
//
// signup needs to be protected by onlyUnAuth middleware
func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()

		email := r.FormValue("email")
		usn := r.FormValue("username")
		password := r.FormValue("password")

		switch err := helpers.Signup(w, r, email, usn, []byte(password)); err {
		case nil:
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		
		case helpers.ErrUserAlreadyExists:
			http.Redirect(w, r, "/auth-err/user-exists", http.StatusSeeOther)
			return
		
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	helpers.RenderTemplate(w, "signup.html", nil)
}

// Login handles "/login" and accepts POST and GET verbs
//
// login needs to be protected by onlyUnAuth middleware
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()

		email := r.FormValue("email")
		password := r.FormValue("password")

		switch err := helpers.Login(w, r, email, []byte(password)); err {
		case nil:
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		
		case helpers.ErrUserNotFound:
			http.Redirect(w, r, "/auth-err/user-not-found", http.StatusSeeOther)
			return

		case bcrypt.ErrMismatchedHashAndPassword:
			http.Redirect(w, r, "/auth-err/wrong-pass", http.StatusSeeOther)
			return

		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	helpers.RenderTemplate(w, "login.html", nil)
}

// Logout handles "/logout" and accepts the GET verb
//
// Logout needs to be protected by onlyAUth middleware
func Logout(w http.ResponseWriter, r *http.Request) {
	err := helpers.Logout(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// CreateRoom handles "/rooms/new" anb accepts POST and GET verbs
//
// CreateRoom needs to be protected by onlyAuth middleware
func CreateRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()

		roomName := r.FormValue("roomName")
		user := r.Context().Value(middleware.AccountKey{})

		res, err := app.DB.SQL.Exec("INSERT INTO rooms (name, owner) VALUES (?, ?)", roomName, user.(models.Account).ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		id, _ := res.LastInsertId()

		message := models.Message{
			Type: 		models.TYPE_NEW_ROOM,
			Sender:		user.(models.Account),
			Data: 		make(map[string]interface{}),
		}

		message.Data["authorName"] = user.(models.Account).Username
		message.Data["authorID"] = user.(models.Account).ID
		message.Data["roomName"] = roomName
		message.Data["roomId"] = id
		
		indexChannels.Broadcast <- message

		http.Redirect(w, r, fmt.Sprintf("/room/%v", id), http.StatusSeeOther)
		return
	}

	helpers.RenderTemplate(w, "newRoom.html", nil)
}

func JoinRoom(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	roomID := params["pk"]

	user := r.Context().Value(middleware.AccountKey{}).(models.Account)

	id, _ := strconv.Atoi(roomID) // Irelevant error as roomId goes through regexp
	room, err := helpers.GetRoom(int64(id))
	switch err {
	case nil:
		ctx := struct {
			Room 	models.Room
			User	models.Account
		} {
			Room: room,
			User: user,
		}
		helpers.RenderTemplate(w, "room.html", ctx)
		return
	
	case helpers.ErrRoomNotFound:
		w.WriteHeader(http.StatusNotFound)
		http.Redirect(w, r, "/room-err/not-found", http.StatusSeeOther)
		return

	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	
}

func JoinRoomWS(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	roomID := params["pk"]

	id, _ := strconv.Atoi(roomID) // Irelevant error as roomId goes through regexp

	room, _ := helpers.GetRoom(int64(id))

	roomChannels, ok := app.ActiveConnections[room.ID]
	// First Connection?
	if !ok {
		roomChannels = models.NewConnStore()
		app.ActiveConnections[room.ID] = roomChannels
	}

	acc := r.Context().Value(middleware.AccountKey{}).(models.Account)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	ws := models.WsConnection {
		Conn: conn,
		Acc: acc,
		Send: make(chan models.Message),
	}

	go ws.Listen(roomChannels.Register, roomChannels.Unregister, roomChannels.Broadcast)

	msg := models.Message {
		Type: models.TYPE_NEW_CONNECTION,
		Sender: acc,
		Data: make(map[string]interface{}),
	}

	msg.Data["authorName"] = acc.Username
	msg.Data["auhtorID"] = acc.ID
	
	roomChannels.Broadcast <- msg
}