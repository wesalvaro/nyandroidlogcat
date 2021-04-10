package nyandroidlogcat

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"
)

var header = regexp.MustCompile(
	`^\[ (?P<month>\d{2})\-(?P<day>\d{2}) ` +
		`(?P<h>\d{2}):(?P<m>\d{2}):(?P<s>\d{2}).(?P<ss>\d+) ` +
		`(?P<pid>\s*\d+):(?P<tid>\s*\d+) (?P<lvl>[DVIEFW])/(?P<tag>.*?)\s+\]$`)
var now = time.Now()

type Entry struct {
	Time    time.Time
	Level   Level
	Tag     string
	Message string
}

func NewEntry(when time.Time, level Level, tag string, message string) *Entry {
	return &Entry{when, level, tag, message}
}

func NewEntryNow(level Level, tag string, message string) *Entry {
	return NewEntry(time.Now(), level, tag, message)
}

func NewEntryFromHeader(text string) *Entry {
	vv := match(header, text)
	if vv == nil {
		return nil
	}
	month, _ := strconv.Atoi(vv["month"])
	day, _ := strconv.Atoi(vv["day"])
	hour, _ := strconv.Atoi(vv["h"])
	minute, _ := strconv.Atoi(vv["m"])
	second, _ := strconv.Atoi(vv["s"])
	nsec, _ := strconv.Atoi(vv["ss"])
	t := time.Date(
		now.Year(),
		time.Month(month),
		day,
		hour, minute, second, nsec*int(math.Pow10(6)),
		now.Location())
	return NewEntry(t, strToLevel([]rune(vv["lvl"])[0]), vv["tag"], "")
}

func (e *Entry) TermUiString() string {
	return fmt.Sprintf(
		"[%c](fg:%s) [%s](fg:blue) [%s](fg:black,mod:bold) %s",
		e.Level.Rune(), e.Level.ColorString(), e.Time.Format("15:04:05"), e.Tag, e.Message)
}

func match(pattern *regexp.Regexp, text string) map[string]string {
	match := pattern.FindStringSubmatch(text)
	if len(match) == 0 {
		return nil
	}
	valueNames := pattern.SubexpNames()
	values := make(map[string]string)
	for i, val := range match {
		if i == 0 {
			continue
		}
		values[valueNames[i]] = val
	}
	return values
}
