package internal

import "github.com/gdamore/tcell/v2"

type Window interface {
	Display()
	HandleKeyEvent(ev *tcell.EventKey, s tcell.Screen) error
}
