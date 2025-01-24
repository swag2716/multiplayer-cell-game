package models

import "github.com/gorilla/websocket"

type Client struct {
	ClientID   string          `json:"clientID"`
	Connection *websocket.Conn `json:"connection"`
	Color      string          `json:"color"`
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
	GameID  string   `json:"gameID"`
	Balls   int      `json:"balls"`
	Clients []Client `json:"clients"`
	State   []string `json:"state"`
}

type Play struct {
	GameID   string `json:"gameID"`
	BallId   int    `json:"ballId"`
	ClientID string `json:"clientID"`
	Color    string `json:"color"`
}

type Join struct {
	ClientID string `json:"clientID"`
	GameID   string `json:"gameID"`
}
