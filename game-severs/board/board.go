package board

// type Cell int

const (
	Empty int = iota
	Ship
	Hit
	Miss
)

type Boards interface {
	PlaceShip(x, y int)
	Fire(x, y int) string
}

type Board struct {
	Grid        [10][10]int
	ShipBattery int
}

func NewBoard() *Board {
	return &Board{
		Grid:        [10][10]int{},
		ShipBattery: 10,
	}
}

func (b *Board) PlaceShip(x, y int) {
	if x >= 0 && x < 10 && y >= 0 && y < 10 {
		b.Grid[x][y] = Ship
	}
}

func (b *Board) Fire(x, y int) string {
	if x < 0 || x > 10 || y < 0 || y > 10 {
		return "invalid coordinates"
	}

	switch b.Grid[x][y] {
	case Ship:
		b.Grid[x][y] = Hit
		return "HIT"
	case Empty:
		b.Grid[x][y] = Miss
		return "MISS"
	case Hit, Miss:
		return "Already fired in this position"
	default:
		return "Unknown"
	}
}

func (b *Board) ReduceBattery(damage int) bool {
	b.ShipBattery -= damage

	return b.ShipBattery <= 0
}

