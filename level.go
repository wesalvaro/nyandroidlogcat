package nyandroidlogcat

import (
	"strings"

	"github.com/fatih/color"
)

type Level int

const (
	Debug Level = iota
	Verbose
	Info
	Warning
	Error
	Fatal
)

func (v Level) String() string {
	switch v {
	case Debug:
		return "Debug"
	case Verbose:
		return "Verbose"
	case Info:
		return "Info"
	case Warning:
		return "Warning"
	case Error:
		return "Error"
	case Fatal:
		return "Fatal"
	}
	panic("Unknown level")
}

func (v Level) Rune() rune {
	return []rune(v.String())[0]
}

func (v Level) Prev() Level {
	switch v {
	case Debug:
		return Debug
	case Verbose:
		return Debug
	case Info:
		return Verbose
	case Warning:
		return Info
	case Error:
		return Warning
	case Fatal:
		return Error
	}
	panic("Unknown level")
}

func (v Level) Next() Level {
	switch v {
	case Debug:
		return Verbose
	case Verbose:
		return Info
	case Info:
		return Warning
	case Warning:
		return Error
	case Error:
		return Fatal
	case Fatal:
		return Fatal
	}
	panic("Unknown level")
}

func (v Level) ColorString() string {
	switch v {
	case Debug:
		return "green"
	case Verbose:
		return "blue"
	case Info:
		return "white"
	case Warning:
		return "yellow"
	case Error:
		return "red"
	case Fatal:
		return "magenta"
	}
	panic("Unknown level")
}

func (v Level) Color() *color.Color {
	switch v {
	case Debug:
		return color.New(color.FgGreen)
	case Verbose:
		return color.New(color.FgBlue)
	case Info:
		return color.New(color.FgWhite)
	case Warning:
		return color.New(color.BgYellow, color.FgBlack)
	case Error:
		return color.New(color.BgRed, color.FgBlack)
	case Fatal:
		return color.New(color.BgHiWhite, color.FgBlack, color.Bold)
	}
	panic("Unknown level")
}

func (v Level) Emoji() string {
	switch v {
	case Debug:
		return `🪲`
	case Verbose:
		return `🔊`
	case Info:
		return `💁‍♂️`
	case Warning:
		return `⚠️ `
	case Error:
		return `⛔️`
	case Fatal:
		return `💀`
	}
	panic("Unknown level")
}

func strToLevel(lvl rune) Level {
	switch lvl {
	case 'D':
		return Debug
	case 'V':
		return Verbose
	case 'I':
		return Info
	case 'W':
		return Warning
	case 'E':
		return Error
	case 'F':
		return Fatal
	}
	panic("Unknown level")
}

func (v *Level) UnmarshalText(b []byte) error {
	*v = strToLevel([]rune(strings.ToUpper(string(b)))[0])
	return nil
}
