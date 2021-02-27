package websocket
type RegistrationResponse struct {
	Message    string `json:"message"`
	PlayerName string `json:"playerName"`
	Timestamp  int64  `json:"timestamp"`
	GameId     int64  `json:"gameId"`
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