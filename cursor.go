package main

type Cursor struct {
	X, Y              int
	SelectedBlockType BlockType
}

func NewCursor() *Cursor {
	return &Cursor{
		X:                 0,
		Y:                 0,
		SelectedBlockType: EmptyCursor,
	}
}
