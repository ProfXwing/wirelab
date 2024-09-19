package blocks

type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
	NonSurrounding
	NoDirection
)

func GetOppositeDirection(dir Direction) Direction {
	switch dir {
	case Left:
		return Right
	case Right:
		return Left
	case Up:
		return Down
	case Down:
		return Up
	}
	return NoDirection
}
