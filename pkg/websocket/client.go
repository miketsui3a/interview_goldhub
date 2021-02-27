package websocket

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"reflect"

	// "strconv"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}




func strToHashStr(input string) string {
	b := []byte(input)
	inputHash := sha256.Sum256(b[:])
	return hex.EncodeToString(inputHash[:])
}

func inputValid(input map[string]interface{}, c *Client) bool {
	if input["message"] == nil {
		c.Conn.WriteJSON(ErrorResponse{Message: strToHashStr("error"), Reason: "I don't understand...", Timestamp: time.Now().UnixNano()})
		return false
	}

	messageType := reflect.TypeOf(input["message"])
	if k := messageType.Kind(); k != reflect.String {
		c.Conn.WriteJSON(ErrorResponse{Message: strToHashStr("error"), Reason: "message should be string", Timestamp: time.Now().UnixNano()})
		return false
	}

	if input["timestamp"] == nil {
		c.Conn.WriteJSON(ErrorResponse{Message: strToHashStr("error"), Reason: "missing timestamp...", Timestamp: time.Now().UnixNano()})
		return false
	}

	timestampType := reflect.TypeOf(input["timestamp"])
	if k := timestampType.Kind(); k != reflect.Float64 {
		fmt.Println(reflect.TypeOf(input["timestamp"]))
		c.Conn.WriteJSON(ErrorResponse{Message: strToHashStr("error"), Reason: "timestamp should be integer", Timestamp: time.Now().UnixNano()})
		return false
	}

	if input["message"] == strToHashStr("registration") {
		if input["playerName"] != nil {
			playerNameType := reflect.TypeOf(input["playerName"])
			if k := playerNameType.Kind(); k != reflect.String {
				c.Conn.WriteJSON(ErrorResponse{Message: strToHashStr("error"), Reason: "playerName should be string", Timestamp: time.Now().UnixNano()})
				return false
			}

			return true
		}
		c.Conn.WriteJSON(ErrorResponse{Message: strToHashStr("error"), Reason: "missing playerName...", Timestamp: time.Now().UnixNano()})
	}

	if input["message"] == strToHashStr("guess") {
		if input["guess"] == nil {
			c.Conn.WriteJSON(ErrorResponse{Message: strToHashStr("error"), Reason: "missing guess...", Timestamp: time.Now().UnixNano()})
			return false
		}

		guessType := reflect.TypeOf(input["guess"])
		if k := guessType.Kind(); k != reflect.Float64 {
			c.Conn.WriteJSON(ErrorResponse{Message: strToHashStr("error"), Reason: "guess should be int", Timestamp: time.Now().UnixNano()})
			return false
		}

		if input["gameId"] == nil {
			c.Conn.WriteJSON(ErrorResponse{Message: strToHashStr("error"), Reason: "missing gameId...", Timestamp: time.Now().UnixNano()})
			return false
		}

		gameIdType := reflect.TypeOf(input["gameId"])
		if k := gameIdType.Kind(); k != reflect.Float64 {
			c.Conn.WriteJSON(ErrorResponse{Message: strToHashStr("error"), Reason: "gameId should be int", Timestamp: time.Now().UnixNano()})
			return false
		}

		return true
	}

	c.Conn.WriteJSON(ErrorResponse{Message: strToHashStr("error"), Reason: "I don't understand...", Timestamp: time.Now().UnixNano()})
	return false
}


func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var tmp map[string]interface{}
		json.Unmarshal([]byte(string(p)), &tmp)
		if inputValid(tmp, c) != true {
			continue
		}
		if tmp["message"] == nil {
			c.Conn.WriteJSON(ErrorResponse{Message: strToHashStr("error"), Reason: "I don't understand...", Timestamp: time.Now().UnixNano()})
			continue
		}
		gameLogic(c, tmp)

	}
}

func (pool *Pool) Start() {
	for {
		select {

		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			break
		case message := <-pool.WinBroadcast:
			fmt.Println("Broadcast win message and next game message")
			pool.GameId++
			pool.Number = rand.Intn(500)
			for client, _ := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
				if err := client.Conn.WriteJSON(StartResponse{Message: strToHashStr("gameStart"), GameId: pool.GameId, Timestamp: time.Now().UnixNano()}); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
