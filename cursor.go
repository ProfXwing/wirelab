package main

import "redstone/blocks"

type Cursor struct {
	X, Y              int
	SelectedBlockType blocks.BlockType
	Direction         blocks.Direction
}

func NewCursor() *Cursor {
	return &Cursor{
		X:                 0,
		Y:                 0,
		SelectedBlockType: blocks.EmptyBlockType,
		Direction:         blocks.Right,
	}
}
