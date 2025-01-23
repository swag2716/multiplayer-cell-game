package models

import "github.com/gorilla/websocket"

type Client struct {
	ClientID   string          `json:"clientID"`
	Connection *websocket.Conn `json:"connection"`
}

type Request struct {
	Method   string `json:"method"`
	ClientID string `json:"clientID"`
}

type Connect struct {
	Method   string `json:"method"`
	ClientID string `json:"clientID"`
}
type Create struct {
	Method   string `json:"method"`
	Game     Game   `json:"game"`
	ClientID string `json:"clientID"`
}

type Game struct {
	GameId string `json:"gameID"`
	Balls  int    `json:"balls"`
}
