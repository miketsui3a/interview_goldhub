package websocket

import (
	"math/rand"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]string
	// Broadcast  chan Message
	WinBroadcast   chan WinResponse
	StartBroadcast chan StartResponse
	Number         int
	GameId         int64
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]string),
		// Broadcast:  make(chan Message),
		WinBroadcast:   make(chan WinResponse),
		StartBroadcast: make(chan StartResponse),
		Number:         rand.Intn(500),
		GameId:         0,
	}
}
