package buffer

type Cursor struct {
	Start *Loc
	End   *Loc
}

func NewCursor() *Cursor {
	start := NewLoc(0, 0)
	end := NewLoc(0, 0)
	return &Cursor{start, end}
}

func (c *Cursor) Move(x, y int) {
	c.Start.X = x
	c.Start.Y = y
	c.End.X = x
	c.End.Y = y
}

func (c *Cursor) MoveUp() {
	c.Move(c.Start.X, c.Start.Y-1)
}

func (c *Cursor) MoveDown() {
	c.Move(c.Start.X, c.Start.Y+1)
}

func (c *Cursor) MoveLeft() {
	c.Move(c.Start.X-1, c.Start.Y)
}

func (c *Cursor) MoveRight() {
	c.Move(c.Start.X+1, c.Start.Y)
}
