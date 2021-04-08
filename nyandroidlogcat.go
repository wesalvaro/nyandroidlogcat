package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

func createLogcatScanner() *bufio.Scanner {
	cmd := exec.Command("adb", "logcat", "--format=long")
	readCloser, _ := cmd.StdoutPipe()
	cmd.Start()
	return bufio.NewScanner(readCloser)
}

func nyan(c *Printer, scanner *bufio.Scanner) {
	e := (*Entry)(nil)
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) == 0 {
			if e != nil {
				c.Print(e)
				e = nil
			}
			continue
		}
		if e != nil {
			e.Message += text
		} else {
			e = newEntry(text)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading log:", err)
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
	nyan(newPrinter(profile), createLogcatScanner())
}
