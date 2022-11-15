package main

import (
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/gustavosvalentim/evilcode/api/buffer"
)

func main() {
	defaultStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err := s.Init(); err != nil {
		panic(err)
	}
	s.SetStyle(defaultStyle)
	s.Clear()
	buf := buffer.NewBuffer(make([]byte, 0), "teste.txt")
	bufWindow := buffer.NewBufWindow(s).
		SetCursor(0, 0).
		SetBuffer(buf)
	// mainloop
	for {
		s.Show()
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
			continue
		}
		bufWindow.Update()
	}
}
