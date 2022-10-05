package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gustavosvalentim/evilcode/api/commands"
	"github.com/gustavosvalentim/evilcode/pkg/fsstorage"
	"github.com/gustavosvalentim/evilcode/tui"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	pages := tview.NewPages()
	mainView := tui.NewMainView("", "")
	fileStorageRepo := fsstorage.NewFileSystemRepo()
	pages.AddAndSwitchToPage("main", mainView.View(), true)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch key := event.Key(); key {
		case tcell.KeyCtrlC:
			app.Stop()
		case tcell.KeyCtrlS:
			saveCommand, err := commands.NewSaveCommand(mainView.GetFilename(), mainView.GetContent())
			if err != nil {
				panic(err)
			}
			saveCommandHandler := commands.NewSaveCommandHandler(fileStorageRepo)
			if err := saveCommandHandler.Handle(saveCommand); err != nil {
				panic(err)
			}
			mainView.SetFilename(mainView.GetFilename())
		default:
		}
		return event
	})
	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
