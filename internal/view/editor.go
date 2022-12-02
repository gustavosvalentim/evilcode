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

	_, h := tui.Screen.Size()

	buf := buffer.NewBuffer(make([][]byte, 1))
	buf.Path = path
	buf.Cursors = append(buf.Cursors, buffer.NewCursor())

	bufView := new(BufView)
	bufView.SetBuffer(buf)
	bufView.SetOffset(0, 1)

	infoView := new(InfoView)
	infoView.SetBuffer(buf)
	infoView.SetOffset(0, h-1)

	e.AppendBuffer(buf)
	e.AppendView(bufView)
	e.AppendView(infoView)

	return e
}

func (e *Editor) AppendView(v View) {
	e.views = append(e.views, v)
}

func (e *Editor) AppendBuffer(b *buffer.Buffer) {
	e.buffers = append(e.buffers, b)
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
