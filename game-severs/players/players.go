package players

import (
	"game-sever/board"

	"github.com/gorilla/websocket"
)

// type Players interface {
// 	Write(message string)
// }

type Player struct {
	Conn  *websocket.Conn
	Name  string
	Board *board.Board // Each player gets their own board.
	Input chan string
}

func NewPlayer(conn *websocket.Conn, board *board.Board, username string) *Player {
	return &Player{
		Conn:  conn,
		Name:  username,
		Board: board,
		Input: make(chan string),
	}
}

// write sends a message to the player.
func (p *Player) Write(message string) error {
	return p.Conn.WriteMessage(websocket.TextMessage, []byte(message))
}
