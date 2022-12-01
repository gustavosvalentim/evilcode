package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gustavosvalentim/evilcode/internal/tui"
	"github.com/gustavosvalentim/evilcode/internal/view"
)

func main() {
	var err error

	if err = tui.Init(); err != nil {
		tui.Terminate()
	}

	editor := view.NewEditor("teste.txt")

	// mainloop
	for {
		ev := tui.Screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			tui.Screen.Sync()
		case *tcell.EventKey:
			editor.ProcessKeyEvent(ev)
		default:
			break
		}

		editor.Draw()

		tui.Screen.Show()
	}
}
