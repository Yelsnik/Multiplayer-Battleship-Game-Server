package server

import (
	"context"
	"encoding/json"
	"fmt"
	"game-sever/board"
	"game-sever/players"
	"game-sever/utils"
	"log"
	"net/http"

	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all connections (for simplicity). Adjust in production.
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Server struct {
	Port       string
	PlayerChan chan *players.Player
	Redis      *redis.Client
}

func NewServer(port string, redis *redis.Client) *Server {
	return &Server{
		Port:       port,
		PlayerChan: make(chan *players.Player, 10),
		Redis:      redis,
	}
}

func (s *Server) WsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP to WebSocket.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("Error reading auth message:", err)
		conn.Close()
		return
	}

	var user User
	if err := json.Unmarshal(msg, &user); err != nil {
		fmt.Println("Error parsing auth request:", err)
		conn.Close()
		return
	}

	// Create a new board for the player.
	b := board.NewBoard()

	// Handle incoming connections.
	s.HandleConnections(conn, b, s.PlayerChan, user.Username)
}

func (s *Server) Start() {
	// playersChan := make(chan *players.Player, 10)

	go func() {
		for {
			p1 := <-s.PlayerChan
			p2 := <-s.PlayerChan
			go s.RunGameSessions(p1, p2, s.PlayerChan)
		}
	}()

	// Start the WebSocket server.
	http.HandleFunc("/ws", s.WsHandler)
	fmt.Println("Battleship WebSocket server listening on", s.Port)
	if err := http.ListenAndServe(s.Port, nil); err != nil {
		log.Fatal("Error accepting connection:", err)
	}
}

func (s *Server) HandleConnections(conn *websocket.Conn, board *board.Board, playersChan chan *players.Player, username string) {
	var player *players.Player

	ctx := context.Background()

	// Check if the player already exists in Redis.
	player, err := s.GetPlayerFromRedis(ctx, username, conn)
	if err != nil {
		// If the player does not exist, create a new player.
		player = players.NewPlayer(conn, board, username)

		playerData := PlayerData{
			Username: player.Name,
			Board: BoardData{
				Grid: board.Grid,
				Ship: SData{
					X: board.Ship.X,
					Y: board.Ship.Y,
				},
				Battery: board.Battery,
			},
		}

		// Save the player data to Redis.
		data, err := json.Marshal(playerData)
		if err != nil {
			log.Fatalf("Error marshalling player data: %v", err)
		}

		if err := s.Redis.Set(ctx, player.Name, data, 0).Err(); err != nil {
			log.Fatalf("Error saving player data: %v", err)
		}
	}

	// The player's ship is randomly placed on the board.
	x := int(utils.RandomInt(0, 9))
	y := int(utils.RandomInt(0, 9))
	player.Board.PlaceShip(x, y)
	player.Board.Battery = 10
	player.Write(fmt.Sprintf("Your ship is placed at (%d, %d).", x, y))
	player.SendBoard()
	// fmt.Sprintln(player.Board.Battery)

	// Welcome message.
	player.Write("Welcome to the ultimate Battleship game server!!!...")

	// Read incoming messages from the player.
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Error reading WebSocket message:", err)
				break
			}
			player.Input <- string(msg)
		}
		close(player.Input)
	}()

	// Enqueue the player for matchmaking.
	playersChan <- player
}

func (s *Server) RunGameSessions(p1, p2 *players.Player, playersChan chan *players.Player) {
	// Prompt both players to accept the pairing.
	result, ok := promptPairing(p1, p2)
	if !ok {
		// If either player rejects, notify both and re-queue the ones that accept.
		p1.Write("Pairing rejected. Returning to lobby.")
		p2.Write("Pairing rejected. Returning to lobby.")

		if result == "player1 rejected" {
			playersChan <- p2
		} else if result == "player2 rejected" {
			playersChan <- p1
		} else {
			playersChan <- p1
			playersChan <- p2
		}
		return
	}

	// Notify both players that the game has started.
	p1.Write(fmt.Sprintf("Game started! You are Player 1. Your ship battery: %d", p1.Board.Battery))
	p2.Write(fmt.Sprintf("Game started! You are Player 2. Your ship battery: %d", p2.Board.Battery))

	// Start the game loop.
	currentPlayer := p1
	opponent := p2

	for {
		currentPlayer.Write("Your turn! Enter command (e.g., 'move 2 2' or'fire 2 2'):")
		opponent.Write("Waiting for opponent's move...")

		// Wait for the current player's input.
		cmd, ok := <-currentPlayer.Input
		if !ok {
			opponent.Write("Opponent disconnected. You win! noerr")
			return
		}

		// Validate the command.
		cmd = strings.TrimSpace(cmd)
		parts := strings.Split(cmd, " ")
		if len(parts) < 3 {
			currentPlayer.Write("Invalid command. Use 'fire x y' or 'move x y'. err")
			continue
		}

		// Extract the action and coordinates from the command.
		action := strings.ToLower(parts[0])
		x, err1 := strconv.Atoi(parts[1])
		y, err2 := strconv.Atoi(parts[2])
		if err1 != nil || err2 != nil {
			currentPlayer.Write("Coordinates must be integers. err")
			continue
		}

		// Process the move command.
		if action == "move" {
			// Use the opponent's board to process the shot.
			result := currentPlayer.Board.MoveShip(x, y)
			if result == fmt.Sprintf("Ship Moved to (%d, %d)", x, y) {
				currentPlayer.Write(fmt.Sprintf("You moved your ship to (%d, %d)! noerr", x, y))
				currentPlayer.SendBoard()
				opponent.Write("Opponent moved their ship! noerr")
			} else if result == "invalid coordinates" {
				currentPlayer.Write("Invalid coordinates. Coordinates must be within the grid. err")
				continue
			} else if result == "Position Occupied" {
				currentPlayer.Write("Invalid coordinates. Position already occupied. err")
				continue
			} else {
				currentPlayer.Write("Invalid coordinates. Your ship is already in this position. err")
				continue
			}

			// Swap turns after moving.
			currentPlayer, opponent = opponent, currentPlayer
			continue
		}

		// Process the fire command.
		if action == "fire" {
			// Convert coordinates from string to integer.
			x, err1 := strconv.Atoi(parts[1])
			y, err2 := strconv.Atoi(parts[2])
			if err1 != nil || err2 != nil {
				currentPlayer.Write("Coordinates must be integers. err")
				continue
			}

			// Use the opponent's board to process the shot.
			result := opponent.Board.Fire(x, y)
			if result == "HIT" {
				battery := opponent.Board.ReduceBattery(1)
				currentPlayer.Write(fmt.Sprintf("Your opponent's battery is now %d", opponent.Board.Battery))
				opponent.Write(fmt.Sprintf("Your battery is now %d", opponent.Board.Battery))

				// Check if the opponent's battery is depleted.
				if battery {
					currentPlayer.Write("You win🎉! Opponent's battery depleted. noerr")
					opponent.Write("You lose😢! Your battery is depleted. err")
					return
				}
			}

			currentPlayer.Write(fmt.Sprintf("Result of firing at (%d, %d): %s info", x, y, result))
			opponent.Write(fmt.Sprintf("Your board was fired at (%d, %d): %s info", x, y, result))
		}

		// Swap turns after each move.
		currentPlayer, opponent = opponent, currentPlayer
	}
}

func promptPairing(p1, p2 *players.Player) (string, bool) {
	// Send pairing invitation messages.
	p1.Write(fmt.Sprintf("You have been paired with %s. Do you accept? (yes/no)", p2.Name))
	p2.Write(fmt.Sprintf("You have been paired with %s. Do you accept? (yes/no)", p1.Name))

	// Wait for responses from both players.
	// In a real-world scenario, you might want to use a timeout here.
	response1 := strings.ToLower(strings.TrimSpace(<-p1.Input))
	response2 := strings.ToLower(strings.TrimSpace(<-p2.Input))

	// Check if either player rejected the pairing.
	if response1 == "no" || response1 == "n" {
		return "player1 rejected", false
	} else if response2 == "no" || response2 == "n" {
		return "player2 rejected", false
	} else {
		return "pairing accepted", true
	}
}
