package websocket

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func containsValue(m map[*Client]string, v string) bool {
	for _, x := range m {
		if x == v {
			return true
		}
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

func gameLogic(c *Client,input map[string]interface{}){
	switch checkHash(input["message"].(string)) {
	case "registration":
		{

			if containsValue(c.Pool.Clients,input["playerName"].(string)){
				c.Conn.WriteJSON(ErrorResponse{Message: strToHashStr("error"), Reason: "this name has been registered", Timestamp: time.Now().UnixNano()})
				return
			}

			if _, ok := c.Pool.Clients[c]; ok {
				c.Conn.WriteJSON(ErrorResponse{Message: strToHashStr("error"), Reason: "you already in the game", Timestamp: time.Now().UnixNano()})
			} else {
				c.Pool.Clients[c] = input["playerName"].(string)
				c.Conn.WriteJSON(RegistrationResponse{Message: strToHashStr("registration"), PlayerName: input["playerName"].(string), Timestamp: int64(input["timestamp"].(float64)), GameId: c.Pool.GameId})

			}
			return
		}
	case "guess":
		{
			if _, ok := c.Pool.Clients[c]; !ok {
				c.Conn.WriteJSON(ErrorResponse{Message: strToHashStr("error"), Reason: "register to join the game", Timestamp: time.Now().UnixNano()})
				return
			}

			fmt.Println(input["gameId"])
			if int64(input["gameId"].(float64)) != c.Pool.GameId {
				c.Conn.WriteJSON(ErrorResponse{Message: strToHashStr("error"), Reason: "no such gameId", Timestamp: time.Now().UnixNano()})
				return
			}
			if int(input["guess"].(float64)) == c.Pool.Number {
				c.Pool.WinBroadcast <- WinResponse{Message: strToHashStr("win"), Answer: c.Pool.Number, Winner: c.Pool.Clients[c], GameId: c.Pool.GameId}
			} else if int(input["guess"].(float64)) < c.Pool.Number {
				c.Conn.WriteJSON(GuessResponse{Message: strToHashStr("guess"), GuessResult: 2, Timestamp: int64(input["timestamp"].(float64)), GameId: c.Pool.GameId})
			} else if int(input["guess"].(float64)) > c.Pool.Number {
				c.Conn.WriteJSON(GuessResponse{Message: strToHashStr("guess"), GuessResult: 1, Timestamp: int64(input["timestamp"].(float64)), GameId: c.Pool.GameId})
			}
			return
		}
	default:
		{
			c.Conn.WriteJSON(ErrorResponse{Message: strToHashStr("error"), Reason: "I don't understand...", Timestamp: time.Now().UnixNano()})
			return
		}
	}
}