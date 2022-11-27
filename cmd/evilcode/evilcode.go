package main

import (
	"github.com/gustavosvalentim/evilcode/internal/buffer"
	"github.com/gustavosvalentim/evilcode/internal/tui"
)

func main() {
	var err error

	if err = tui.Init(); err != nil {
		tui.Terminate()
	}

	buf := buffer.NewBufferFromFile("teste.txt")
	bufWindow := buffer.NewBufWindow().
		SetBuffer(buf)

	tui.AddWindow(bufWindow)

	// mainloop
	for {
		if err = tui.HandleEvent(); err != nil {
			tui.Terminate()
		}
		tui.UpdateScreen()
	}
}
