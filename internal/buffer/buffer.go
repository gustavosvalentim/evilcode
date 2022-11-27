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

type Cursor struct {
	start *Loc
	end   *Loc
}

func (c *Cursor) IsSelection() bool {
	return c.start.x == c.end.x && c.start.y == c.end.y
}

type Buffer struct {
	Path       string
	lines      [][]byte
	cursors    []*Cursor
	modified   bool
	modifiedAt time.Time
}

func NewBuffer(path string) *Buffer {
	loc := NewLoc(0, 0)
	cur := &Cursor{loc, loc}
	buf := &Buffer{
		path,
		make([][]byte, 1),
		[]*Cursor{cur},
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

func (b *Buffer) NewLine() {
	cur := b.cursors[0].start
	x, y := cur.x, cur.y
	logging.Logf("[Buffer.NewLine] %d - %d", x, y)
	fromRow := b.lines[y]
	rowSize := len(fromRow)
	if rowSize == 0 || x == rowSize {
		b.lines = append(b.lines, make([]byte, 0))
	} else {
		b.splitLine(cur)
	}
	cur.x = 0
	cur.y += 1
	b.UpdateModified(true)
	logging.Logf("NewLine line count %d", len(b.lines))
}

func (b *Buffer) Delete() error {
	cur := b.cursors[0].start
	x, y := cur.x, cur.y
	if x == 0 && y > 0 {
		b.joinLines(y-1, y)
		cur.x = len(b.lines[y-1])
		cur.y -= 1
	} else if x > 0 {
		b.removeCharAtLoc(cur)
		cur.x -= 1
	}

	b.UpdateModified(true)

	return nil
}

func (b *Buffer) Write(c rune) {
	cur := b.cursors[0].start
	appendChar := append(b.lines[cur.y][:cur.x], byte(c))
	b.lines[cur.y] = append(appendChar, b.lines[cur.y][cur.x:]...)
	cur.x += 1
	b.UpdateModified(true)
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
	logging.Log(fmt.Sprintf("[Buffer.removeCharAtLoc] %d %d", loc.x, loc.y))
	b.lines[loc.y] = append(b.lines[loc.y][:loc.x-1], b.lines[loc.y][loc.x:]...)
}

func (b *Buffer) splitLine(loc *Loc) {
	logging.Log(fmt.Sprintf("[Buffer.splitLine] %d - %d", loc.x, loc.y))
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
	var content string
	for _, l := range b.lines {
		content += fmt.Sprintf("%s\n", string(l))
	}
	return content
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
