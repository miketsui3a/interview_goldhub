package websocket

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"reflect"

	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

type GenericRequest struct {
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"-"`
}

type RegistrationRequest struct {
	Message    string `json:"message"`
	PlayerName string `json:"playerName"`
	Timestamp  int64  `json:"timestamp"`
}

type RegistrationResponse struct {
	Message    string `json:"message"`
	PlayerName string `json:"playerName"`
	Timestamp  int64  `json:"timestamp"`
	GameId     int64  `json:"gameId"`
}

type GuessRequest struct {
	Message   string `json:"message"`
	Guess     int    `json:"guess"`
	Timestamp int64  `json:"timestamp"`
	GameId    int64  `json:"gameId"`
}

type GuessResponse struct {
	Message     string `json:"message"`
	GuessResult int    `json:"guessResult"`
	Timestamp   int64  `json:"timestamp"`
	GameId      int64  `json:"gameId"`
}

type WinResponse struct {
	Message string `json:"message"`
	Answer  int    `json:"answer"`
	Winner  string `json:"winner"`
	GameId  int64  `json:"gameId"`
}

type StartResponse struct {
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
	GameId    int64  `json:"gameId"`
}

type ErrorResponse struct {
	Message   string `json:"message"`
	Reason    string `json:"reason"`
	Timestamp int64  `json:"timestamp"`
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

	return false
}

func checkHash(hashString string) string {
	registration := []byte("registration")
	guess := []byte("guess")

	registrationHash := sha256.Sum256(registration[:])
	guessHash := sha256.Sum256(guess[:])

	switch hashString {
	case hex.EncodeToString(registrationHash[:]):
		return "registration"
	case hex.EncodeToString(guessHash[:]):
		return "guess"
	default:
		return "?"
	}
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

		// message := Message{Type: messageType, Body: string(p)}

		var tmp map[string]interface{}
		json.Unmarshal([]byte(string(p)), &tmp)
		if inputValid(tmp, c) != true {
			continue
		}
		if tmp["message"] == nil {
			c.Conn.WriteJSON(ErrorResponse{Message: strToHashStr("error"), Reason: "I don't understand...", Timestamp: time.Now().UnixNano()})
			continue
		}
		switch checkHash(tmp["message"].(string)) {
		case "registration":
			{
				fmt.Println(tmp["playerName"])
				if _, ok := c.Pool.Clients[c]; ok {
					c.Conn.WriteJSON(ErrorResponse{Message: strToHashStr("error"), Reason: "you already in the game", Timestamp: time.Now().UnixNano()})
				} else {
					c.Pool.Clients[c] = tmp["playerName"].(string)
					c.Conn.WriteJSON(RegistrationResponse{Message: strToHashStr("registration"), PlayerName: tmp["playerName"].(string), Timestamp: int64(tmp["timestamp"].(float64)), GameId: c.Pool.GameId})

				}
				break
			}
		case "guess":
			{
				if _, ok := c.Pool.Clients[c]; !ok {
					c.Conn.WriteJSON(ErrorResponse{Message: strToHashStr("error"), Reason: "register to join the game", Timestamp: time.Now().UnixNano()})
					continue
				}

				fmt.Println(tmp["gameId"])
				if int64(tmp["gameId"].(float64)) != c.Pool.GameId {
					c.Conn.WriteJSON(ErrorResponse{Message: strToHashStr("error"), Reason: "no such gameId", Timestamp: time.Now().UnixNano()})
					continue
				}
				// num, _ := strconv.Atoi(string(tmp["guess"].(string)))
				if int(tmp["guess"].(float64)) == c.Pool.Number {
					fmt.Println("bingo")
					c.Pool.WinBroadcast <- WinResponse{Message: strToHashStr("win"), Answer: c.Pool.Number, Winner: c.Pool.Clients[c], GameId: c.Pool.GameId}
				} else if int(tmp["guess"].(float64)) < c.Pool.Number {
					fmt.Println("ok")
					c.Conn.WriteJSON(GuessResponse{Message: strToHashStr("guess"), GuessResult: 2, Timestamp: int64(tmp["timestamp"].(float64)), GameId: c.Pool.GameId})
					// c.Conn.WriteJSON(GuessResponse{Message: "guess", GuessResult: 2,GameId: c.Pool.GameId})
				} else if int(tmp["guess"].(float64)) > c.Pool.Number {
					c.Conn.WriteJSON(GuessResponse{Message: strToHashStr("guess"), GuessResult: 1, Timestamp: int64(tmp["timestamp"].(float64)), GameId: c.Pool.GameId})
					// c.Conn.WriteJSON(GuessResponse{Message: "guess", GuessResult: 1,GameId: c.Pool.GameId})
				}
				break
			}
		default:
			{
				c.Conn.WriteJSON(ErrorResponse{Message: strToHashStr("error"), Reason: "I don't understand...", Timestamp: time.Now().UnixNano()})
			}
		}

		num, err := strconv.Atoi(string(p))
		if c.Pool.Number == num {
			fmt.Println("ok")
			// c.Pool.Broadcast <- message
			newNum := rand.Intn(500)
			fmt.Println(newNum)
			c.Pool.Number = newNum
		} else {
			// c.Conn.WriteJSON(Message{Type: 1, Body: "DDDD"})
		}
		// fmt.Printf("Message Received: %+v\n", message)
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
