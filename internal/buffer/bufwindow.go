package buffer

import (
	"errors"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/gustavosvalentim/evilcode/internal/logging"
	"github.com/gustavosvalentim/evilcode/internal/tui"
)

type BufWindow struct {
	buf        *Buffer
	hasChanges bool
}

func NewBufWindow() *BufWindow {
	w := new(BufWindow)
	return w
}

func (w *BufWindow) SetBuffer(b *Buffer) *BufWindow {
	w.buf = b
	return w
}

func (w *BufWindow) HandleKeyEvent(ev *tcell.EventKey) error {
	key := ev.Key()
	logging.Log(fmt.Sprintf("[BufWindow.HandleKeyEvent] Key pressed: %s", ev.Name()))
	cur := w.buf.cursors[0].start
	if key == tcell.KeyRune {
		c := ev.Rune()
		w.buf.Write(c)
		w.hasChanges = true
	} else {
		switch ev.Key() {
		case tcell.KeyCtrlC:
			return errors.New("keyboard interrupt")
		case tcell.KeyCtrlS:
			Save(w.buf)
		case tcell.KeyEnter:
			w.buf.NewLine()
			w.hasChanges = true
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			err := w.buf.Delete()
			if err != nil {
				return err
			}
			w.hasChanges = true
		case tcell.KeyRight:
			if cur.x < len(w.buf.lines[cur.y]) {
				cur.x += 1
			}
		case tcell.KeyLeft:
			if cur.x > 0 {
				cur.x -= 1
			}
		case tcell.KeyUp:
			if cur.y > 0 {
				if prevLineLen := len(w.buf.lines[cur.y-1]); cur.x > prevLineLen {
					cur.x = prevLineLen
				}
				cur.y -= 1
			}
		case tcell.KeyDown:
			if cur.y < len(w.buf.lines[cur.y]) {
				if nextLineLen := len(w.buf.lines[cur.y+1]); cur.x > nextLineLen {
					cur.x = nextLineLen
				}
				cur.y += 1
			}
		default:
		}
	}
	return nil
}

func (w *BufWindow) Display() {
	cur := w.buf.cursors[0].start
	if w.hasChanges {
		width, _ := tui.ScreenSize()
		for y := 0; y < len(w.buf.lines); y++ {
			for x := 0; x < width; x++ {
				c := byte(0)
				if x <= len(w.buf.lines[y])-1 {
					c = w.buf.lines[y][x]
				}
				tui.DrawCharacter(x, y, rune(c))
			}
		}
		w.hasChanges = false
	}

	tui.ShowCursor(cur.x, cur.y)

	logging.Logf("[BufWindow.Display] Cursor location X: %d Y: %d", cur.x, cur.y)
}
