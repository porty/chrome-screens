package chrome

import (
	"io"
	"strings"
)

type PrefixWriter struct {
	writer io.Writer
	prefix string
}

func (w *PrefixWriter) Write(p []byte) (n int, err error) {
	lines := strings.Split(string(p), "\n")
	for _, line := range lines {
		if _, err := w.writer.Write([]byte(w.prefix + line + "\n")); err != nil {
			return 0, err
		}
	}
	return len(p), nil
}
