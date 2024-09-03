package blocks

type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
	NonSurrounding
)
