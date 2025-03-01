package server

import (
	"context"
	"encoding/json"
	"game-sever/board"
	"game-sever/players"

	"github.com/gorilla/websocket"
)

type Cell int

type User struct {
	Username string
}

type PlayerData struct {
	Username string    `json:"username"`
	Board    BoardData `json:"board"`
}

type BoardData struct {
	Grid        [10][10]int `json:"grid"`
	ShipBattery int         `json:"ship_battery"`
}

func (s *Server) GetPlayerFromRedis(ctx context.Context, username string, conn *websocket.Conn) (*players.Player, error) {
	result, err := s.Redis.Get(ctx, username).Result()
	if err != nil {
		return nil, err
	}

	var playerData PlayerData
	if err := json.Unmarshal([]byte(result), &playerData); err != nil {
		return nil, err
	}

	player := &players.Player{
		Conn:  conn,
		Name:  playerData.Username,
		Board: NewBoardFromData(&playerData.Board),
		Input: make(chan string),
	}

	return player, nil
}

func (s *Server) SavePlayerToRedis() {}

// NewBoardFromData creates a new Board instance from BoardData.
func NewBoardFromData(data *BoardData) *board.Board {
	var b board.Board
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			// Cast the integer value to Cell.
			b.Grid[i][j] = data.Grid[i][j]
		}
	}
	b.ShipBattery = data.ShipBattery
	return &b
}
