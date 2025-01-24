package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/websocket-multiplayer-game-go/models"
)

var clients = make(map[string]interface{})
var games = make(map[string]interface{})
var color = make(map[int]string)

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

		color[0] = "red"
		color[1] = "blue"
		color[2] = "yellow"

		Connect(conn)
		go func(conn *websocket.Conn) {
			for {
				_, p, err := conn.ReadMessage()
				if err != nil {
					log.Fatal("error in reading message from client: ", err)
				}
				var msg models.Connect

				err = json.Unmarshal(p, &msg) //converts []byte to interface
				if err != nil {
					log.Fatal("error in unmarshalling message from client: ", err)
				}

				if msg.Method == "create" {
					Create(msg)
				}

				if msg.Method == "join" {
					fmt.Println("recived join game")
					var joinMsg models.Join
					err := json.Unmarshal(p, &joinMsg) //converts []byte to interface
					if err != nil {
						log.Fatal("error in unmarshalling message from client: ", err)
					}
					Join(joinMsg)
				}

				if msg.Method == "play" {
					var playMsg models.Play
					err := json.Unmarshal(p, &playMsg) //converts []byte to interface
					if err != nil {
						log.Fatal("error in unmarshalling message from client: ", err)
					}

					Play(playMsg)
				}
			}
		}(conn)

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

	var allClients []models.Client
	games[gameID] = models.Game{
		GameID:  gameID,
		Balls:   9,
		Clients: allClients,
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

func Join(msg models.Join) {
	clientID := msg.ClientID
	gameID := msg.GameID

	fmt.Println("game id", gameID)

	conn := clients[clientID].(models.Client).Connection

	game := games[gameID].(models.Game)

	clientsInGame := len(game.Clients)
	if clientsInGame >= 3 {
		return
	}

	var client models.Client

	client.ClientID = clientID
	client.Color = color[clientsInGame]

	game.Clients = append(games[gameID].(models.Game).Clients, client)

	if len(game.Clients) == 3 {
		fmt.Println("sending game state")
		SendGameStatePeriodically(gameID)
	}

	games[gameID] = game
	payload := map[string]interface{}{
		"method": "join",
		"game":   game,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Fatal("error in marshalling the payload: ", err)
	}

	for _, client := range game.Clients {
		fmt.Println("joined game client id:", client.ClientID)
		conn = clients[client.ClientID].(models.Client).Connection
		conn.WriteMessage(1, data)
	}

}

func Play(msg models.Play) {
	// clientID := msg.ClientID
	gameID := msg.GameID
	ballID := msg.BallId
	color := msg.Color
	game := games[gameID].(models.Game)
	state := game.State
	if state == nil || len(state) == 0 {
		state = make([]string, game.Balls+1)
	}
	fmt.Println("state, ballID, color", state, ballID, color)
	state[ballID] = color
	game.State = state
	games[gameID] = game
	fmt.Println("state1, ballID1, color1", state, ballID, color)
}

func SendGameStatePeriodically(gameID string) {

	go func() {
		log.Println("Starting periodic game state updates for gameID:", gameID)
		for i := 0; i < 30; i++ {
			SendGameState(gameID)
			time.Sleep(1 * time.Second)
		}
	}()
}

func SendGameState(gameID string) {
	game := games[gameID].(models.Game)
	state := game.State

	fmt.Println("sending state in every 500 ms")

	payload := map[string]interface{}{
		"method": "state",
		"state":  state,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Fatal("error in marshalling the payload: ", err)
	}

	for _, client := range game.Clients {
		conn := clients[client.ClientID].(models.Client).Connection
		conn.WriteMessage(1, data)
	}
}

func Stop(gameID string) {
	// payload :=
}
