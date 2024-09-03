package main

import (
	"redstone/blocks"

	"github.com/gdamore/tcell/v2"
)

const GameWidth = 50
const GameHeight = 50

type Game struct {
	display_handler *DisplayHandler
	block_handler   *BlockHandler
	cursor          *Cursor
}

func NewGame() *Game {
	cursor := NewCursor()

	g := &Game{
		display_handler: nil,
		block_handler:   nil,
		cursor:          cursor,
	}

	g.display_handler = NewDisplayHandler(g)
	g.block_handler = NewBlockHandler(g)

	return g
}

func (g *Game) Run() {
	defer g.display_handler.Quit()

	for {
		g.display_handler.DrawCursor(g.cursor)
		g.display_handler.Show()

		if !g.HandleEvents() {
			break
		}
	}
}

func (g *Game) HandleEvents() bool {
	ev := g.display_handler.screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyCtrlC {
			return false
		} else {
			g.HandleInput(ev)
		}
	}

	return true
}

func (g *Game) HandleInput(ev *tcell.EventKey) {
	dh := g.display_handler

	switch ev.Key() {
	case tcell.KeyUp:
		if g.cursor.Y > 0 {
			g.cursor.Y -= 1
		}

	case tcell.KeyDown:
		if g.cursor.Y < GameHeight-1 {
			g.cursor.Y += 1
		}

	case tcell.KeyLeft:
		if g.cursor.X > 0 {
			g.cursor.X -= 1
		}

	case tcell.KeyRight:
		if g.cursor.X < GameWidth-1 {
			g.cursor.X += 1
		}

	case tcell.KeyEnter:
		if g.cursor.SelectedBlockType == blocks.EmptyBlockType {
			x, y := g.cursor.X, g.cursor.Y
			selectedBlock := g.block_handler.GetBlock(x, y)

			if selectedBlock == nil {
				return
			}

			if selectedBlock.GetBlockType() == blocks.LeverType {
				selectedBlock.SetPowered(!selectedBlock.IsPowered())

				g.block_handler.UpdateSurroundingBlocks(x, y)
			}
		} else {
			g.block_handler.NewBlock(g.cursor, true)
		}

	case tcell.KeyEscape:
		g.cursor.SelectedBlockType = blocks.EmptyBlockType

	case tcell.KeyRune:
		switch ev.Rune() {
		case 'w':
			g.cursor.SelectedBlockType = blocks.WireType
		case 'p':
			g.cursor.SelectedBlockType = blocks.PoweredBlockType
		case 'l':
			g.cursor.SelectedBlockType = blocks.WiredLampType
		case 't':
			g.cursor.SelectedBlockType = blocks.LeverType
		case 'i':
			g.cursor.SelectedBlockType = blocks.InverterType
		case 'r':
			switch g.cursor.Direction {
			case blocks.Right:
				g.cursor.Direction = blocks.Down
			case blocks.Down:
				g.cursor.Direction = blocks.Left
			case blocks.Left:
				g.cursor.Direction = blocks.Up
			case blocks.Up:
				g.cursor.Direction = blocks.Right
			}

		default:
			return
		}
	default:
		return
	}

	// redraw screen, assuming we did something
	dh.RedrawScreen()
}
