package buffer

// Initial size of the buffer
const chunk_size = 256 * 1024

type Buffer interface {
	// Write r at position in the buffer
	WriteAt(pos int, r []byte) error

	// Delete character at position from the buffer
	DeleteAt(pos int) error

	// Returns the rune at a given position
	RuneAt(pos int) (rune, error)

	// Return the character index based on
	// row and col
	Index(row, col int) int

	// Returns the size of the buffer
	Size() int
}

type buf struct {
	data []byte
}

func NewBuffer() Buffer {
	return &buf{make([]byte, chunk_size)}
}

func (b *buf) WriteAt(pos int, r []byte) error {
	if pos == len(b.data) {
		copy(b.data[pos:pos+len(r)], r)
	} else {
		copy(b.data[pos+len(r):cap(b.data)], b.data[pos:len(b.data)])
		copy(b.data[pos:pos+len(r)], r)
	}
	return nil
}

func (b *buf) DeleteAt(pos int) error {
	b.data = append(b.data[:pos], b.data[pos+1:]...)
	return nil
}

func (b *buf) RuneAt(pos int) (rune, error) {
	return rune(b.data[pos]), nil
}

func (b *buf) Index(row, col int) int {
	var i int
	for i = 0; row > 0 || i < col; i++ {
		if b.data[i] == '\n' {
			row--
		}
	}
	return i
}

func (b *buf) Size() int {
	return len(b.data)
}
