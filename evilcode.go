package main // github.com/gustavosvalentim/evilcode

import (
	"bufio"
	"os"

	"github.com/gdamore/tcell/v2"
)

type row struct {
	chars    []byte
	modified bool
}

type editor struct {
	filename string
	rows     []row
	cx, cy   int
}

var screen tcell.Screen
var e editor

/** Terminal **/

func terminate() {
	screen.Fini()
	screen = nil
	os.Exit(0)
}

func initScreen() {
	var err error

	if screen, err = tcell.NewScreen(); err != nil {
		terminate()
	}

	if err = screen.Init(); err != nil {
		terminate()
	}
}

func updateScreen() {
	w, h := screen.Size()
	// Write through the entire screen
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			ch := rune(0)
			// Check if there is a character at the current position
			// If there is set `ch` to this character
			if y < len(e.rows) {
				r := e.rows[y]
				if x < len(r.chars) {
					ch = rune(r.chars[x])
				}
			}
			// Write `ch` to the screen
			screen.SetContent(x, y, ch, nil, tcell.StyleDefault)
		}
	}

	fixCursor()

	screen.ShowCursor(e.cx, e.cy)
	screen.Show()
}

/** Editor **/

func initEditor() {
	if len(os.Args) > 1 {
		e.filename = os.Args[1]
		f, err := os.Open(e.filename)
		if err != nil {
			terminate()
		}
		defer f.Close()
		reader := bufio.NewReader(f)

		for line, err := reader.ReadBytes('\n'); err == nil; line, err = reader.ReadBytes('\n') {
			e.rows = append(e.rows, row{line, false})
		}
	} else {
		e.rows = append(e.rows, row{make([]byte, 0), false})
	}
	e.cy = len(e.rows) - 1
	e.cx = len(e.rows[e.cy].chars)
}

func writeCharAt(r *row, ch rune, at int) {
	tail := make([]byte, len(r.chars[at:]))
	copy(tail, r.chars[at:])
	r.chars = append(append(r.chars[:at], byte(ch)), tail...)
	r.modified = true
	e.cx++
}

func newLine() {
	e.rows = append(e.rows, row{make([]byte, 0), true})
	e.cx = 0
	e.cy++
}

func splitRowAt(r *row, at int) {
	newRow := row{r.chars[at:], true}
	r.chars = r.chars[:at]
	r.modified = true
	tail := make([]row, len(e.rows[e.cy+1:]))
	copy(tail, e.rows[e.cy+1:])
	e.rows = append(append(e.rows[:e.cy+1], newRow), tail...)
	e.cx = 0
	e.cy++
}

func joinRows(brow, arow int) {
	rowToAppend := e.rows[arow]
	isize := len(e.rows[brow].chars)
	e.rows[brow].chars = append(e.rows[brow].chars, rowToAppend.chars...)
	e.rows = append(e.rows[:arow], e.rows[arow+1:]...)
	e.cy = brow
	e.cx = isize
}

func deleteCharAt(r *row, at int) {
	r.chars = append(r.chars[:at], r.chars[at+1:]...)
	e.cx--
}

func moveCursorLeft() {
	e.cx--
}

func moveCursorRight() {
	e.cx++
}

func moveCursorUp() {
	e.cy--
}

func moveCursorDown() {
	e.cy++
}

func fixCursor() {
	if e.cy < 0 {
		e.cy = 0
	}
	if maxY := len(e.rows) - 1; e.cy > maxY {
		e.cy = maxY
	}
	if e.cx < 0 {
		e.cx = 0
	}
	if maxX := len(e.rows[e.cy].chars); e.cx > maxX {
		e.cx = maxX
	}
}

func processKeyEvent(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyRune:
		writeCharAt(&e.rows[e.cy], ev.Rune(), e.cx)
	case tcell.KeyUp:
		moveCursorUp()
	case tcell.KeyDown:
		moveCursorDown()
	case tcell.KeyLeft:
		moveCursorLeft()
	case tcell.KeyRight:
		moveCursorRight()
	case tcell.KeyBackspace:
	case tcell.KeyBackspace2:
		if e.cx > 0 {
			deleteCharAt(&e.rows[e.cy], e.cx-1)
		} else if e.cy > 0 {
			joinRows(e.cy-1, e.cy)
		}
	case tcell.KeyEnter:
		if e.cx < len(e.rows[e.cy].chars) {
			splitRowAt(&e.rows[e.cy], e.cx)
		} else {
			newLine()
		}
	case tcell.KeyCtrlQ:
		terminate()
	default:
	}
}

/** Mainloop **/

func main() {
	initEditor()
	initScreen()

	for {
		ev := screen.PollEvent()
		switch evAsType := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			processKeyEvent(evAsType)
		default:
		}

		updateScreen()
	}
}
