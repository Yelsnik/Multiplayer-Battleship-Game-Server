package board

import (
	"fmt"
	"game-sever/utils"
)

// type Cell int

// const (
// 	Empty int = iota
// 	Ship
// 	Hit
// 	Miss
// )

type S struct {
	X       int
	Y       int
}

type Boards interface {
	PlaceShip(x, y int)
	Fire(x, y int) string
}

type Board struct {
	Grid [10][10]int
	Ship S
	Battery int
}

func NewBoard() *Board {
	return &Board{
		Grid: [10][10]int{},
		Ship: S{
			X:       0,
			Y:       0,
		},
		Battery: 10,
	}
}

func (b *Board) PlaceShip(x, y int) {
	if x >= 0 && x < 10 && y >= 0 && y < 10 {
		// // b.Grid[x][y] = Ship
		// b.Ship.X = x
		// b.Ship.Y = y
		// Keep finding an empty spot if the chosen position is occupied
		for b.Grid[x][y] != 0 { // 0 means Empty
			x = int(utils.RandomInt(0, 9))
			y = int(utils.RandomInt(0, 9))
		}

		// Mark the ship's position in the grid
		b.Grid[x][y] = 1 // 1 represents Ship

		// Store ship details
		b.Ship.X = x
		b.Ship.Y = y
	}

}

func (b *Board) Fire(x, y int) string {
	if x < 0 || x > 10 || y < 0 || y > 10 {
		return "invalid coordinates"
	}

	fmt.Println(b.Ship.X, b.Ship.Y)

	if x == b.Ship.X && y == b.Ship.Y {
		// b.Grid[x][y] = Hit
		return "HIT"
	} else {
		// b.Grid[x][y] = Miss
		return "MISS"
	}

	// switch b.Grid[x][y] {
	// case Ship:
	// 	b.Grid[x][y] = Hit
	// 	return "HIT"
	// case Empty:
	// 	b.Grid[x][y] = Miss
	// 	return "MISS"
	// case Hit, Miss:
	// 	return "Already fired in this position"
	// default:
	// 	return "Unknown"
	// }

}

func (b *Board) ReduceBattery(damage int) bool {
	fmt.Println(b.Battery)
	b.Battery -= damage

	return b.Battery <= 0
}

func (b *Board) ToDisplayString() string {
	var output string

	for i := range b.Grid {
		for j := range b.Grid[i] {
			if i == b.Ship.X && j == b.Ship.Y {
				output += " [s] " // Ship's position
			} else {
				output += " [] " // Empty space
			}
		}
		output += "\n"
	}

	// for i := 0; i < 10; i++ {
	// 	for j := 0; j < 10; j++ {
	// 		if i == b.Ship.X && j == b.Ship.Y {
	// 			output += " S " // Ship's position
	// 		} else {
	// 			output += " . " // Empty space
	// 		}
	// 	}
	// 	output += "\n"
	// }
	return output
}

func (b *Board) MoveShip(x, y int) string {
	// Check if the new position is within the grid
	if x < 0 || x > 10 || y < 0 || y > 10 {
		return "invalid coordinates"
	}

	// Check if the ship is already in the new position
	if x == b.Ship.X && y == b.Ship.Y {
		return "Already in this position"
	}

	// check if the new position is occupied
	if b.Grid[x][y] != 0 {
		return "Position Occupied"
	}

	// Mark the ship's position in the grid
	b.Grid[x][y] = 1 // 1 represents Ship

	// Store ship details
	b.Ship.X = x
	b.Ship.Y = y

	return fmt.Sprintf("Ship Moved to (%d, %d)", x, y)
}