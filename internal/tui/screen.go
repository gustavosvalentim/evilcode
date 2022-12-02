package tui

import (
	"os"

	"github.com/gdamore/tcell/v2"
)

var Screen tcell.Screen

func Init() error {
	var err error
	Screen, err = tcell.NewScreen()
	if err != nil {
		return err
	}
	err = Screen.Init()
	if err != nil {
		return err
	}

	Screen.SetStyle(Style())

	return nil
}

func Style() tcell.Style {
	return tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorReset)
}

func DrawCharacter(x, y int, c rune) {
	Screen.SetContent(x, y, c, nil, Style())
}

func DrawLine(startX, y int, text string) {
	for i, c := range text {
		DrawCharacter(startX+i, y, c)
	}
}

func Terminate() {
	Screen.Fini()
	os.Exit(0)
}
