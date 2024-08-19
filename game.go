package main

import (
	"github.com/gdamore/tcell/v2"
)

const GameWidth = 50
const GameHeight = 50

type Game struct {
	display_handler *DisplayHandler
	block_handler   *BlockHandler
	cursor          *Cursor
	custom_events   chan int
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

		dh.RedrawScreen()

	case tcell.KeyDown:
		if g.cursor.Y < GameHeight-1 {
			g.cursor.Y += 1
		}
		dh.RedrawScreen()

	case tcell.KeyLeft:
		if g.cursor.X > 0 {
			g.cursor.X -= 1
		}
		dh.RedrawScreen()

	case tcell.KeyRight:
		if g.cursor.X < GameWidth-1 {
			g.cursor.X += 1
		}
		dh.RedrawScreen()

	case tcell.KeyRune:
		switch ev.Rune() {
		case 'r':
			g.cursor.SelectedBlockType = Redstone
			dh.RedrawScreen()
		case 'b':
			g.cursor.SelectedBlockType = RedstoneBlock
			dh.RedrawScreen()
		case 'l':
			g.cursor.SelectedBlockType = RedstoneLamp
			dh.RedrawScreen()
		}

	case tcell.KeyEnter:
		if g.cursor.SelectedBlockType == EmptyCursor {
			return
		}

		g.block_handler.NewBlock(g.cursor)

		dh.RedrawScreen()
	}
}
