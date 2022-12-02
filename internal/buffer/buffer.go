package buffer

import (
	"fmt"

	"github.com/gustavosvalentim/evilcode/internal/logging"
)

type Loc struct {
	X int
	Y int
}

func NewLoc(x, y int) *Loc {
	return &Loc{x, y}
}

type Buffer struct {
	Lines   [][]byte
	Cursors []*Cursor
	Path    string
}

func NewBuffer(content [][]byte) *Buffer {
	buf := new(Buffer)
	buf.Lines = content

	return buf
}

func (b *Buffer) NewLine(loc *Loc) {
	logging.Logf("[Buffer.NewLine] %d - %d", loc.X, loc.Y)
	if len(b.Lines[loc.Y]) == 0 || loc.X == len(b.Lines[loc.Y]) {
		b.Lines = append(b.Lines, make([]byte, 0))
	} else {
		b.splitLine(loc)
	}
	logging.Logf("NewLine line count %d", len(b.Lines))
}

func (b *Buffer) Delete(loc *Loc) error {
	if loc.X == 0 && loc.Y > 0 {
		b.joinLines(loc.Y-1, loc.Y)
	} else if loc.X > 0 {
		b.removeCharAtLoc(loc)
	}

	return nil
}

func (b *Buffer) Write(loc *Loc, c rune) {
	appendChar := append(b.Lines[loc.Y][:loc.X], byte(c))
	b.Lines[loc.Y] = append(appendChar, b.Lines[loc.Y][loc.X:]...)
}

func (b *Buffer) joinLines(l0, l1 int) {
	b.Lines[l0] = append(b.Lines[l0], b.Lines[l1]...)
	if l1 < len(b.Lines)-1 {
		b.Lines = append(b.Lines[:l0], b.Lines[l0+1:]...)
	} else {
		b.Lines = b.Lines[:len(b.Lines)-1]
	}
}

func (b *Buffer) removeCharAtLoc(loc *Loc) {
	logging.Log(fmt.Sprintf("[Buffer.removeCharAtLoc] %d %d", loc.X, loc.Y))
	b.Lines[loc.Y] = append(b.Lines[loc.Y][:loc.X-1], b.Lines[loc.Y][loc.X:]...)
}

func (b *Buffer) splitLine(loc *Loc) {
	logging.Log(fmt.Sprintf("[Buffer.splitLine] %d - %d", loc.X, loc.Y))
	newLine := b.Lines[loc.Y][loc.X:]
	b.Lines[loc.Y] = b.Lines[loc.Y][:loc.X]
	if loc.Y == len(b.Lines)-1 {
		b.Lines = append(b.Lines, newLine)
	} else {
		endLines := [][]byte{newLine}
		endLines = append(endLines, b.Lines[loc.Y+1:]...)
		b.Lines = append(b.Lines[:loc.Y+1], endLines...)
		logging.Log(fmt.Sprintf("[Buffer.splitLine] %s - %s", b.Lines[loc.Y], endLines))
	}
}

func (b *Buffer) Text() string {
	var content string
	for _, l := range b.Lines {
		content += fmt.Sprintf("%s\n", string(l))
	}
	return content
}
