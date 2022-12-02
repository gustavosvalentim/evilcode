package view

import "github.com/gdamore/tcell/v2"

type View interface {
	SetOffset(x, y int)
	Draw()
	OnKeyEvent(ev *tcell.EventKey)
}
