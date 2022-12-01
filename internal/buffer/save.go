package buffer

import (
	"bufio"
	"os"

	"github.com/gustavosvalentim/evilcode/internal/logging"
)

func SaveToFile(b *Buffer) error {
	f, err := os.Create(b.Path)

	if err != nil {
		logging.Logf("[buffer.SaveToFile] %s", err.Error())
		return err
	}

	defer f.Close()

	writer := bufio.NewWriter(f)

	if _, err := writer.WriteString(b.Text()); err != nil {
		logging.Logf("[buffer.SaveToFile] %s", err.Error())
		return err
	}

	writer.Flush()

	return nil
}
