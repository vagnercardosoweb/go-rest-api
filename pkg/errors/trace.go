package errors

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/vagnercardosoweb/go-rest-api/pkg/utils"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func GetStack(skip int) []string {
	var lines [][]byte
	buf := new(bytes.Buffer)
	var lastFile string

	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)

		if file != lastFile {
			if err := utils.ValidateSourceFilePath(file); err != nil {
				fmt.Fprintf(buf, "    %s: %s\n", function(pc), dunno)
				continue
			}

			// #nosec G304 - Path is validated by ValidateSourceFilePath above
			data, err := os.ReadFile(file)
			if err != nil {
				fmt.Fprintf(buf, "    %s: %s\n", function(pc), dunno)
				continue
			}

			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}

		fmt.Fprintf(buf, "    %s: %s\n", function(pc), source(lines, line))
	}

	return strings.Split(buf.String(), "\n")
}

func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed

	if n < 0 || n >= len(lines) {
		return dunno
	}

	return bytes.TrimSpace(lines[n])
}

func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}

	name := []byte(fn.Name())
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}

	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}

	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
