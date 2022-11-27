package logging

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

var writer *bufio.Writer = newLogWriter()

func newLogWriter() *bufio.Writer {
	f, err := os.Create("log.txt")
	if err != nil {
		panic(err)
	}
	return bufio.NewWriter(f)
}

func Log(message string) {
	if _, err := writer.WriteString(fmt.Sprintf("[%s] %s\n", time.Now(), message)); err != nil {
		panic(err)
	}
	if err := writer.Flush(); err != nil {
		panic(err)
	}
}

func Logf(message string, params ...any) {
	Log(fmt.Sprintf(message, params...))
}
