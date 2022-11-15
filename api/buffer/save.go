package buffer

import (
	"bufio"
	"os"
)

func Save(b *Buffer) error {
	f, err := os.Create(b.Path)
	if err != nil {
		return err
	}
	defer func() error {
		if err := f.Close(); err != nil {
			return err
		}
		return nil
	}()
	w := bufio.NewWriter(f)
	if _, err := w.WriteString(b.Text()); err != nil {
		return err
	}
	w.Flush()
	b.UpdateModified(false)
	return nil
}
