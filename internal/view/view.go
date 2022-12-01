package view

import "github.com/gdamore/tcell/v2"

type View interface {
	Draw()
	OnKeyEvent(ev *tcell.EventKey)
}
