package nyandroidlogcat

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func CreateLogcatScanner() *bufio.Scanner {
	cmd := exec.Command("adb", "logcat", "--format=long")
	readCloser, _ := cmd.StdoutPipe()
	cmd.Start()
	return bufio.NewScanner(readCloser)
}

func Nyan(scanner *bufio.Scanner) <- chan *Entry {
	out := make(chan *Entry)
	go func() {
		defer close(out)
		e := (*Entry)(nil)
		var mLines []string
		for scanner.Scan() {
			text := scanner.Text()
			if len(text) == 0 {
				if e != nil {
					e.Message = strings.Join(mLines, "\n")
					out <- e
					mLines = nil
					e = nil
				}
				continue
			}
			if e != nil {
				mLines = append(mLines, text)
			} else {
				e = NewEntry(text)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "error reading log:", err)
		}
	}()
	return out
}

func NyanForever() <- chan *Entry {
	out := make(chan *Entry)
	go func() {
		for {
			entries := Nyan(CreateLogcatScanner())
			fmt.Println("Waiting for device...")
			for e := range entries {
				out <- e
			}
		}
	}()
	return out
}
