package main

import "github.com/gdamore/tcell/v2"

type GameStyles struct {
	DefaultStyle tcell.Style
	TextStyle    tcell.Style
}

func InitStyles() *GameStyles {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	pieceStyle := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)

	return &GameStyles{DefaultStyle: defStyle, TextStyle: pieceStyle}
}
