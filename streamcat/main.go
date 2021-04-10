package main

import (
	"io/ioutil"
	"os"
	"path"

	lc "wesalvaro.com/nyandroidlogcat"
)

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
	configs, _ := lc.NewPrinterConfigMap(profileJson)
	profile, ok := configs[getProfileName()]
	if !ok {
		configs, _ := lc.NewPrinterConfigMap([]byte(sampleProfile))
		profile = configs["default"]
	}
	printer := lc.NewPrinter(profile)
	for e := range lc.NyanForever() {
		printer.Print(e)
	}
}
