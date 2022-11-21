package main

import (
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/gustavosvalentim/evilcode/internal"
	"github.com/gustavosvalentim/evilcode/internal/buffer"
)

func main() {
	var err error

	defaultStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	err = s.Init()
	if err != nil {
		panic(err)
	}
	s.SetStyle(defaultStyle)
	s.Clear()
	buf := buffer.NewBuffer(make([]byte, 0), "teste.txt")
	bufWindow := buffer.NewBufWindow(s).
		SetCursor(0, 0).
		SetBuffer(buf)
	visibleWindows := []internal.Window{
		bufWindow,
	}

	handleEvent := func() {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if err := bufWindow.HandleKeyEvent(ev, s); err != nil {
				s.Fini()
				os.Exit(0)
			}
		default:
			break
		}
	}

	// mainloop
	for {
		s.Show()
		handleEvent()
		for _, w := range visibleWindows {
			w.Display()
		}
	}
}
