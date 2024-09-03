package main

import (
	"log"
	"redstone/blocks"

	"github.com/gdamore/tcell/v2"
)

type DisplayHandler struct {
	screen tcell.Screen
	styles *GameStyles
	game   *Game
}

func NewDisplayHandler(game *Game) *DisplayHandler {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Error creating s: %v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("Error initializing s: %v", err)
	}

	styles := InitStyles()
	s.SetStyle(styles.DefaultStyle)

	display_handler := &DisplayHandler{
		screen: s,
		styles: styles,
		game:   game,
	}

	s.Clear()

	return display_handler
}

func (dh *DisplayHandler) DrawText(x1, y1, x2, y2 int, style tcell.Style, text string) {
	s := dh.screen

	row := y1
	col := x1
	for _, r := range text {
		s.SetContent(col, row, r, nil, style)
		col++
		if col > x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func (dh *DisplayHandler) DrawCursor(cursor *Cursor) {
	block_handler := dh.game.block_handler

	cursorBlock := block_handler.NewBlock(cursor, false)
	surroundingBlocks := block_handler.GetSurroundingBlocks(cursor.X, cursor.Y)

	char := cursorBlock.GetRune(surroundingBlocks)

	s := dh.screen
	s.SetContent(cursor.X, cursor.Y, char, nil, dh.styles.TextStyle)
}

func (dh *DisplayHandler) DrawBlock(block blocks.Block) {
	block_handler := dh.game.block_handler

	x, y := block.GetPosition()
	surroundingBlocks := block_handler.GetSurroundingBlocks(x, y)
	char := block.GetRune(surroundingBlocks)

	s := dh.screen

	s.SetContent(x, y, char, nil, dh.styles.DefaultStyle)
}

func (dh *DisplayHandler) RedrawScreen() {
	dh.Clear()

	for _, block := range dh.game.block_handler.Blocks {
		dh.DrawBlock(block)
	}
}

func (dh *DisplayHandler) Quit() {
	maybePanic := recover()
	dh.screen.Fini()
	if maybePanic != nil {
		panic(maybePanic)
	}
}

func (dh *DisplayHandler) Show() {
	dh.screen.Show()
}

func (dh *DisplayHandler) Clear() {
	dh.screen.Clear()
}
