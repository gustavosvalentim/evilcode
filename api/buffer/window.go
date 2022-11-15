package buffer

import (
	"errors"
	"unicode"

	"github.com/gdamore/tcell/v2"
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
	newCx, newCy := w.cursor.X(), w.cursor.Y()
	if c := ev.Rune(); unicode.IsGraphic(c) {
		drawCharacter(w.cursor.X(), w.cursor.Y(), s, c)
		w.buf.Write(c)
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
		case tcell.KeyBackspace | tcell.KeyBackspace2:
			if w.cursor.X() == 0 && w.cursor.Y() == 0 {
				break
			}
			// Update buffer
			w.buf.Remove(NewLoc(w.cursor.X(), w.cursor.Y()), NewLoc(w.cursor.X(), w.cursor.Y()))
			newCx -= 1
		default:
		}
	}
	w.SetCursor(newCx, newCy)
	return nil
}

func (w *BufWindow) Update() {
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
	w.ShowCursor(w.cursor.X(), w.cursor.Y())
}
