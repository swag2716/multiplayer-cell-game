package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/websocket-multiplayer-game-go/models"
)

var clients = make(map[string]interface{})
var games = make(map[string]interface{})

func main() {
	upgrader := websocket.Upgrader{}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(*http.Request) bool {
			return true
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal("error in upgrading to websocket: ", err)
		}

		Connect(conn)

		var msg models.Connect

		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("error in reading message from client: ", err)
		}

		err = json.Unmarshal(p, &msg) //converts []byte to interface
		if err != nil {
			log.Fatal("error in unmarshalling message from client: ", err)
		}

		if msg.Method == "create" {
			Create(msg)
		}

	})
	http.ListenAndServe(":8080", nil)
}

func Connect(conn *websocket.Conn) {
	clientID := uuid.NewString()

	clients[clientID] = models.Client{
		ClientID:   clientID,
		Connection: conn,
	}
	payload := models.Connect{
		Method:   "connect",
		ClientID: clientID,
	}

	data, err := json.Marshal(payload) // converts interface to byte
	if err != nil {
		log.Fatal("error in marshalling the payload: ", err)
	}

	conn.WriteMessage(1, data)
}

func Create(msg models.Connect) {
	clientID := msg.ClientID
	gameID := uuid.NewString()
	games[gameID] = models.Game{
		GameId: gameID,
		Balls:  9,
	}

	payload := models.Create{
		Game:   games[gameID].(models.Game),
		Method: "create",
	}
	data, err := json.Marshal(payload)
	if err != nil {
		log.Fatal("error in marshalling the payload: ", err)
	}
	conn := clients[clientID].(models.Client).Connection

	conn.WriteMessage(1, data)
}
