package buffer

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gustavosvalentim/evilcode/internal/logging"
)

var bufferList []*Buffer = make([]*Buffer, 0)

func GetBufferList() []*Buffer {
	return bufferList
}

type Loc struct {
	x int
	y int
}

func NewLoc(x, y int) *Loc {
	return &Loc{x, y}
}

func (loc *Loc) X() int {
	return loc.x
}

func (loc *Loc) Y() int {
	return loc.y
}

type Buffer struct {
	lines      [][]byte
	Path       string
	modified   bool
	modifiedAt time.Time
}

func NewBuffer(path string) *Buffer {
	buf := &Buffer{
		make([][]byte, 1),
		path,
		false,
		time.Now(),
	}
	bufferList = append(bufferList, buf)
	return buf
}

func NewBufferFromFile(path string) *Buffer {
	if buf := GetOpenBuffer(path); buf != nil {
		return buf
	}
	buf := NewBuffer(path)
	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
		return buf
	}
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	for _, l := range strings.Split(string(data), "/n") {
		buf.lines = append(buf.lines, []byte(l))
	}
	return buf
}

func (b *Buffer) NewLine(from, to *Loc) {
	fromRow := b.lines[from.y]
	rowSize := len(fromRow)
	if rowSize == 0 || from.x == rowSize {
		LogLocation("NewLine", from.y, from.x, to.y, to.x)
		b.lines = append(b.lines, make([]byte, 0))
	} else {
		LogLocation("NewLine splitLine", from.y, from.x, to.y, to.x)
		b.splitLine(from)
	}
	logging.Log(fmt.Sprintf("NewLine line count %d", len(b.lines)))
}

func (b *Buffer) Delete(start, end *Loc) {
	if start.y == end.y {
		y := end.y
		x := end.x
		if y > 0 && x == 0 {
			b.joinLines(y-1, y)
		} else {
			b.removeCharAtLoc(end)
		}
	} else {
		// TODO: multi line selection
	}
}

func (b *Buffer) Write(c rune) {
	lnum := len(b.lines) - 1
	b.lines[lnum] = append(b.lines[lnum], byte(c))
}

func (b *Buffer) joinLines(l0, l1 int) {
	b.lines[l0] = append(b.lines[l0], b.lines[l1]...)
	if l1 < len(b.lines)-1 {
		b.lines = append(b.lines[:l0], b.lines[l0+1:]...)
	} else {
		b.lines = b.lines[:len(b.lines)-1]
	}
}

func (b *Buffer) removeCharAtLoc(loc *Loc) {
	logging.Log(fmt.Sprintf("[removeCharAtLoc] %d %d", loc.x, loc.y))
	if loc.x == len(b.lines[loc.y]) {
		b.lines[loc.y] = b.lines[loc.y][:loc.x-1]
	} else {
		b.lines[loc.y] = append(b.lines[loc.y][:loc.x], b.lines[loc.y][loc.x+1:]...)
	}
}

func (b *Buffer) splitLine(loc *Loc) {
	newLine := b.lines[loc.y][loc.x:]
	b.lines[loc.y] = b.lines[loc.y][:loc.x]
	if loc.y == len(b.lines)-1 {
		b.lines = append(b.lines, newLine)
	} else {
		endLines := [][]byte{newLine}
		endLines = append(endLines, b.lines[loc.y+1:]...)
		b.lines = append(b.lines[:loc.y+1], endLines...)
		logging.Log(fmt.Sprintf("[Buffer.splitLine] %s - %s", b.lines[loc.y], endLines))
	}
}

func (b *Buffer) UpdateModified(modified bool) {
	b.modified = modified
	if modified {
		b.modifiedAt = time.Now()
	}
}

func (b *Buffer) Text() string {
	tLines := make([]string, 0)
	for _, l := range b.lines {
		tLines = append(tLines, string(l))
	}
	return strings.Join(tLines, "\n")
}

func GetOpenBuffer(path string) *Buffer {
	for _, b := range bufferList {
		if b.Path == path {
			return b
		}
	}
	return nil
}

func LogLocation(name string, fromRow, fromColumn, toRow, toColumn int) {
	logging.Log(fmt.Sprintf("%s (%d - %d / %d - %d)", name, fromRow, fromColumn, toRow, toColumn))
}
