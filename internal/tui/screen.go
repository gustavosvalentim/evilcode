package tui

import (
	"os"

	"github.com/gdamore/tcell/v2"
)

var screen tcell.Screen
var windows []Window
var focusedWindow *Window

func Init() error {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		return err
	}
	err = screen.Init()
	if err != nil {
		return err
	}

	screen.Show()

	screen.SetStyle(Style())
	screen.Clear()

	return nil
}

func Style() tcell.Style {
	return tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorReset)
}

func DrawCharacter(x, y int, c rune) {
	screen.SetContent(x, y, c, nil, tcell.StyleDefault)
}

func dispatchKeyEvents(event *tcell.EventKey) error {
	if err := (*focusedWindow).HandleKeyEvent(event); err != nil {
		return err
	}

	return nil
}

func HandleEvent() error {
	ev := screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventResize:
		screen.Sync()
	case *tcell.EventKey:
		if err := dispatchKeyEvents(ev); err != nil {
			Terminate()
		}
	default:
		break
	}

	return nil
}

func ScreenSize() (int, int) {
	return screen.Size()
}

func ShowCursor(x, y int) {
	screen.ShowCursor(x, y)
}

func Terminate() {
	screen.Fini()
	os.Exit(0)
}

func AddWindow(w Window) {
	focusedWindow = &w
	windows = append(windows, w)
}

func UpdateScreen() {
	for _, w := range windows {
		w.Display()
	}

	screen.Show()
}
