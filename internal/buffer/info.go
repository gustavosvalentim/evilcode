package buffer

type BufferInfo struct {
	Lines           int
	Columns         int
	CurrentLocation *Loc
}

func GetBufferInfo(buf *Buffer) *BufferInfo {
	info := new(BufferInfo)
	curLoc := buf.Cursors[0].Start

	info.Lines = len(buf.Lines)
	info.Columns = len(buf.Lines[curLoc.Y])
	info.CurrentLocation = curLoc

	return info
}
