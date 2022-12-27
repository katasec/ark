package logs

import (
	"fmt"
	"io"
	"os"
	"time"
)

type ExampleLogger struct {
	w io.Writer
}

func NewExampleLogger() *ExampleLogger {
	return &ExampleLogger{
		w: os.Stdout,
	}
}
func (e ExampleLogger) Write(data []byte) (int, error) {

	tmStamp := fmt.Sprint(time.Now().Format("2006-01-02 15:04:05"))

	data = []byte(fmt.Sprintf("%s  %s", tmStamp, string(data)))

	n, err := e.w.Write(data)
	if err != nil {
		return n, err
	}
	if n != len(data) {
		return n, io.ErrShortWrite
	}
	return len(data), nil
}
