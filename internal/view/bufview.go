package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gustavosvalentim/evilcode/internal/buffer"
	"github.com/gustavosvalentim/evilcode/internal/tui"
)

type BufView struct {
	buf     *buffer.Buffer
	offsetX int
	offsetY int
}

func (view *BufView) SetBuffer(buf *buffer.Buffer) {
	view.buf = buf
}

func (view *BufView) SetOffset(x, y int) {
	view.offsetX = x
	view.offsetY = y
}

func (view *BufView) Draw() {
	width, height := tui.Screen.Size()

	for y := 0; y < height-view.offsetY; y++ {
		for x := 0; x < width; x++ {
			c := byte(0)

			if y < len(view.buf.Lines) {
				row := view.buf.Lines[y]

				if x < len(row) {
					c = row[x]
				}
			}

			tui.DrawCharacter(x+view.offsetX, y, rune(c))
		}
	}

	for _, c := range view.buf.Cursors {
		tui.Screen.ShowCursor(c.Start.X, c.Start.Y)
	}
}

func (view *BufView) OnKeyEvent(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyRune:
		for _, c := range view.buf.Cursors {
			view.buf.Write(c.Start, ev.Rune())
			c.MoveRight()
		}
	case tcell.KeyBackspace:
		for _, c := range view.buf.Cursors {
			if c.Start.X == 0 && c.Start.Y == 0 {
				break
			}

			view.buf.Delete(c.Start)

			if c.Start.X == 0 && c.Start.Y > 0 {
				c.Start.X = len(view.buf.Lines[c.Start.Y])
				c.MoveUp()
			} else {
				c.MoveLeft()
			}
		}
	case tcell.KeyBackspace2:
		for _, c := range view.buf.Cursors {
			if c.Start.X == 0 && c.Start.Y == 0 {
				break
			}

			view.buf.Delete(c.Start)

			if c.Start.X == 0 && c.Start.Y > 0 {
				c.Start.X = len(view.buf.Lines[c.Start.Y])
				c.MoveUp()
			} else {
				c.MoveLeft()
			}
		}
	case tcell.KeyEnter:
		for _, c := range view.buf.Cursors {
			view.buf.NewLine(c.Start)
			c.Start.X = 0
			c.MoveDown()
		}
	case tcell.KeyUp:
		for _, c := range view.buf.Cursors {
			c.MoveUp()
		}
	case tcell.KeyDown:
		for _, c := range view.buf.Cursors {
			c.MoveDown()
		}
	case tcell.KeyLeft:
		for _, c := range view.buf.Cursors {
			c.MoveLeft()
		}
	case tcell.KeyRight:
		for _, c := range view.buf.Cursors {
			c.MoveRight()
		}
	case tcell.KeyCtrlS:
		buffer.SaveToFile(view.buf)
	}
}
