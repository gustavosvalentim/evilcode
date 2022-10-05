package tui

import "github.com/rivo/tview"

type MainView struct {
	textEditor    *TextEditor
	filenameInput *tview.InputField
}

func NewMainView(filename, content string) *MainView {
	filenameInput := tview.NewInputField().
		SetText(filename).
		SetLabel("Filename")
	return &MainView{
		textEditor:    NewTextEditor(filename, content),
		filenameInput: filenameInput,
	}
}

func (view *MainView) View() *tview.Grid {
	return tview.NewGrid().
		SetRows(0, 1).
		AddItem(view.textEditor.Window(), 0, 0, 1, 2, 0, 0, true).
		AddItem(view.filenameInput, 1, 0, 1, 2, 0, 0, false)
}

func (view *MainView) GetFilename() string {
	return view.filenameInput.GetText()
}

func (view *MainView) GetContent() string {
	return view.textEditor.Content()
}

func (view *MainView) SetFilename(filename string) *MainView {
	view.textEditor.SetFilename(filename)
	return view
}
