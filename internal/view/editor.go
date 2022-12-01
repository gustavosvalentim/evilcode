package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gustavosvalentim/evilcode/internal/buffer"
	"github.com/gustavosvalentim/evilcode/internal/tui"
)

type Editor struct {
	buffers    []*buffer.Buffer
	views      []View
	activeView int
}

func NewEditor(path string) *Editor {
	e := new(Editor)

	buf := buffer.NewBuffer(make([][]byte, 1))
	buf.Path = path
	buf.Cursors = append(buf.Cursors, buffer.NewCursor())

	view := new(BufView)
	view.SetBuffer(buf)

	e.views = append(e.views, view)
	e.buffers = append(e.buffers, buf)

	return e
}

func (e *Editor) ActiveView() View {
	return e.views[e.activeView]
}

func (e *Editor) ProcessKeyEvent(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyCtrlC:
		tui.Terminate()
	default:
		e.ActiveView().OnKeyEvent(ev)
	}
}

func (e *Editor) Draw() {
	for _, view := range e.views {
		view.Draw()
	}
}
