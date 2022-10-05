package tui

import (
	"github.com/rivo/tview"
)

type TextEditor struct {
	filename string
	content  string
	window   *tview.TextArea
}

func NewTextEditor(filename, content string) *TextEditor {
	window := tview.NewTextArea().
		SetText(content, false)
	window.SetTitle(filename).SetBorder(true)
	return &TextEditor{
		filename: filename,
		content:  content,
		window:   window,
	}
}

func (editor *TextEditor) Window() *tview.TextArea {
	return editor.window
}

func (editor *TextEditor) Filename() string {
	return editor.filename
}

func (editor *TextEditor) Content() string {
	return editor.window.GetText()
}

func (editor *TextEditor) SetFilename(filename string) *TextEditor {
	editor.filename = filename
	editor.window.SetTitle(filename)
	return editor
}
