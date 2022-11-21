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

func (loc *Loc) IsEqual(other *Loc) bool {
	return loc.x == other.x && loc.y == other.y
}

type LineArray struct {
	lines [][]byte
}

func NewLineArray(data []byte) *LineArray {
	lines := strings.Split(string(data), "\n")
	bLines := make([][]byte, 0)
	for _, l := range lines {
		bLines = append(bLines, []byte(l))
	}
	return &LineArray{bLines}
}

func (l *LineArray) Line(y int) []byte {
	return l.lines[y]
}

func (l *LineArray) NewLine(from, to *Loc) {
	fromRow := l.lines[from.y]
	rowSize := len(fromRow)
	if rowSize == 0 || from.x == rowSize-1 {
		LogLocation("NewLine", from.y, from.x, to.y, to.x)
		l.lines = append(l.lines, make([]byte, 0))
	} else {
		LogLocation("NewLine splitLine", from.y, from.x, to.y, to.x)
		l.splitLine(from)
	}
	logging.Log(fmt.Sprintf("NewLine line count %d", len(l.lines)))
}

func (l *LineArray) Write(c rune) {
	lnum := len(l.lines) - 1
	l.lines[lnum] = append(l.lines[lnum], byte(c))
}

func (l *LineArray) joinLines(l0, l1 int) {
	l.lines[l0] = append(l.lines[l0], l.lines[l1]...)
	if l1 < len(l.lines)-1 {
		l.lines = append(l.lines[:l0], l.lines[l0+1:]...)
	} else {
		l.lines = l.lines[:len(l.lines)-1]
	}
}

func (l *LineArray) removeCharAtLoc(loc *Loc) {
	logging.Log(fmt.Sprintf("[removeCharAtLoc] %d %d", loc.x, loc.y))
	if loc.x == len(l.lines[loc.y]) {
		l.lines[loc.y] = l.lines[loc.y][:loc.x-1]
	} else {
		l.lines[loc.y] = append(l.lines[loc.y][:loc.x], l.lines[loc.y][loc.x:]...)
	}
}

func (l *LineArray) Remove(start, end *Loc) {
	if start.y == end.y {
		y := end.y
		x := end.x
		if y > 0 && x == 0 {
			l.joinLines(y-1, y)
		} else {
			l.removeCharAtLoc(end)
		}
	} else {
		// TODO: multi line selection
	}
}

func (l *LineArray) splitLine(loc *Loc) {
	newLine := l.lines[loc.y][loc.x:]
	if loc.y == len(l.lines)-1 {
		l.lines = append(l.lines, newLine)
	} else {
		endLines := [][]byte{newLine}
		endLines = append(endLines, l.lines[loc.y:]...)
		l.lines = append(l.lines[:loc.y], endLines...)
	}
}

type Buffer struct {
	*LineArray
	Path       string
	modified   bool
	modifiedAt time.Time
}

func NewBuffer(data []byte, path string) *Buffer {
	buf := &Buffer{
		NewLineArray(data),
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
	buf := NewBuffer(make([]byte, 0), path)
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
	buf.LineArray = NewLineArray(data)
	return buf
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

func (b *Buffer) Size(y int) int {
	return len(b.lines[y])
}

func LogLocation(name string, fromRow, fromColumn, toRow, toColumn int) {
	logging.Log(fmt.Sprintf("%s (%d - %d / %d - %d)", name, fromRow, fromColumn, toRow, toColumn))
}
