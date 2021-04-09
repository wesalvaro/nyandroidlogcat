package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

func createLogcatScanner() *bufio.Scanner {
	cmd := exec.Command("adb", "logcat", "--format=long")
	readCloser, _ := cmd.StdoutPipe()
	cmd.Start()
	return bufio.NewScanner(readCloser)
}

func nyan(scanner *bufio.Scanner, out chan *Entry) {
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
			e = newEntry(text)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading log:", err)
	}
}

func nyanForever(out chan *Entry) {
	one := make(chan *Entry)
	for {
		fmt.Println("Waiting for device...")
		go nyan(createLogcatScanner(), one)
		for e := range one {
			out <- e
		}
		one = make(chan *Entry)
	}
}

func getProfileName() string {
	if len(os.Args) == 2 {
		return os.Args[1]
	} else {
		return "default"
	}
}

const sampleProfile = `{
	"default": {
		"time-format": "15:04:05",
		"tag": {
			"filter": [],
			"ignore": [],
			"show": true
		},
		"message": {
			"filter": [],
			"highlight": [],
			"highlight-color": "yellow"
		},
		"level": {
			"bound": "Debug",
			"first": true,
			"color": true,
			"emoji": false,
			"long": false,
			"show": true
		}
	}
}`

func main() {
	home, _ := os.UserHomeDir()
	profileJson, _ := ioutil.ReadFile(path.Join(home, ".nyandroidlogcat.json"))
	configs, _ := NewPrinterConfigMap(profileJson)
	profile, ok := configs[getProfileName()]
	if !ok {
		configs, _ := NewPrinterConfigMap([]byte(sampleProfile))
		profile = configs["default"]
	}
	printer := newPrinter(profile)
	entries := make(chan *Entry)
	go nyanForever(entries)
	for e := range entries {
		printer.Print(e)
	}
}
