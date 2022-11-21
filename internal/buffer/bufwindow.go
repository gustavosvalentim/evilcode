package buffer

import (
	"errors"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/gustavosvalentim/evilcode/internal/logging"
)

type BufWindow struct {
	tcell.Screen

	cursor *Loc
	buf    *Buffer
}

func NewBufWindow(s tcell.Screen) *BufWindow {
	w := new(BufWindow)
	w.Screen = s
	return w
}

func (w *BufWindow) SetCursor(x, y int) *BufWindow {
	w.cursor = NewLoc(x, y)
	return w
}

func (w *BufWindow) SetBuffer(b *Buffer) *BufWindow {
	w.buf = b
	for y, l := range b.lines {
		for x, bc := range l {
			drawCharacter(x, y, w.Screen, rune(bc))
		}
	}
	return w
}

func drawCharacter(x, y int, s tcell.Screen, c rune) {
	s.SetContent(x, y, c, nil, tcell.StyleDefault)
}

func (w *BufWindow) HandleKeyEvent(ev *tcell.EventKey, s tcell.Screen) error {
	newCx, newCy := w.cursor.x, w.cursor.y
	key := ev.Key()
	logging.Log(fmt.Sprintf("[BufWindow.HandleKeyEvent] Key pressed: %s", string(key)))
	if key == tcell.KeyRune {
		c := ev.Rune()
		drawCharacter(w.cursor.x, w.cursor.y, s, c)
		w.buf.Write(c)
		w.buf.UpdateModified(true)
		newCx += 1
	} else {
		switch ev.Key() {
		case tcell.KeyCtrlC:
			return errors.New("keyboard interrupt")
		case tcell.KeyCtrlS:
			Save(w.buf)
		case tcell.KeyEnter:
			newCy += 1
			newCx = 0
			loc := NewLoc(w.cursor.x, w.cursor.y)
			w.buf.NewLine(loc, loc)
			w.buf.UpdateModified(true)
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if w.cursor.x == 0 && w.cursor.y == 0 {
				break
			}
			// Update buffer
			w.buf.Remove(NewLoc(w.cursor.x, w.cursor.y), NewLoc(w.cursor.x, w.cursor.y))
			w.buf.UpdateModified(true)
			newCx -= 1
		default:
		}
	}
	w.SetCursor(newCx, newCy)
	return nil
}

func (w *BufWindow) Display() {
	width, _ := w.Screen.Size()
	for y := 0; y < len(w.buf.lines); y++ {
		for x := 0; x < width; x++ {
			c := byte(0)
			if x <= len(w.buf.lines[y])-1 {
				c = w.buf.lines[y][x]
			}
			drawCharacter(x, y, w.Screen, rune(c))
		}
	}

	lastRowNum := len(w.buf.lines) - 1

	if w.cursor.y > lastRowNum {
		w.SetCursor(len(w.buf.lines[lastRowNum]), lastRowNum)
	}

	w.ShowCursor(w.cursor.x, w.cursor.y)

	logging.Log(fmt.Sprintf("[BufWindow.Display] Cursor location X: %d Y: %d", w.cursor.x, w.cursor.y))
}
