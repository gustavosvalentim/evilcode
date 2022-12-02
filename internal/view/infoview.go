package view

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/gustavosvalentim/evilcode/internal/buffer"
	"github.com/gustavosvalentim/evilcode/internal/tui"
)

type InfoView struct {
	buf     *buffer.Buffer
	info    *buffer.BufferInfo
	offsetX int
	offsetY int
}

func (view *InfoView) SetBuffer(buf *buffer.Buffer) {
	view.buf = buf
}

func (view *InfoView) SetOffset(x, y int) {
	view.offsetX = x
	view.offsetY = y
}

func (view *InfoView) Draw() {
	view.update()

	loc := view.info.CurrentLocation
	text := fmt.Sprintf("L %d/%d C %d/%d", loc.Y+1, view.info.Lines, loc.X, view.info.Columns)
	tui.DrawLine(view.offsetX, view.offsetY, text)
}

func (view *InfoView) update() {
	view.info = buffer.GetBufferInfo(view.buf)
}

func (view *InfoView) OnKeyEvent(ev *tcell.EventKey) {}
