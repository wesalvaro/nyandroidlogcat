package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

type TagConfig struct {
	Filter []string `json:"filter"`
	Ignore []string `json:"ignore"`
	Show   bool     `json:"show"`
}

type MsgConfig struct {
	Filter         []string     `json:"filter"`
	Highlight      []string     `json:"highlight"`
	HighlightColor *Highlighter `json:"highlight-color"`
}

type LvlConfig struct {
	Bound Level `json:"bound"`
	Emoji bool  `json:"emoji"`
	Long  bool  `json:"long"`
	First bool  `json:"first"`
	Color bool  `json:"color"`
	Show  bool  `json:"show"`
}

type PrinterConfig struct {
	TimeFormat string    `json:"time-format"`
	Tag        TagConfig `json:"tag"`
	Msg        MsgConfig `json:"message"`
	Lvl        LvlConfig `json:"level"`
}

func NewPrinterConfig(blob []byte) (*PrinterConfig, error) {
	var config PrinterConfig
	if err := json.Unmarshal(blob, &config); err != nil {
		return nil, fmt.Errorf("could not parse config: %s", err)
	}
	return &config, nil
}

func NewPrinterConfigMap(blob []byte) (map[string]*PrinterConfig, error) {
	var configs map[string]*PrinterConfig
	if err := json.Unmarshal(blob, &configs); err != nil {
		return nil, fmt.Errorf("could not parse config: %s", err)
	}
	return configs, nil
}

type Highlighter struct {
	Color *color.Color
}

func (h *Highlighter) UnmarshalText(text []byte) error {
	c := (*color.Color)(nil)
	switch strings.ToLower(string(text)) {
	case "black":
		c = color.New(color.BgBlack)
	case "red":
		c = color.New(color.BgRed)
	case "green":
		c = color.New(color.BgGreen)
	case "yellow":
		c = color.New(color.BgYellow)
	case "blue":
		c = color.New(color.BgBlue)
	case "magenta":
		c = color.New(color.BgMagenta)
	case "cyan":
		c = color.New(color.BgCyan)
	case "white":
		c = color.New(color.BgWhite)
	}
	if c == nil {
		return fmt.Errorf("unknown color: %s", text)
	}
	h.Color = c.Add(color.FgBlack)
	return nil
}

type Printer struct {
	TimeFormat        string
	TagFilter         []*regexp.Regexp
	TagIgnore         []*regexp.Regexp
	MsgFilter         []*regexp.Regexp
	MsgHighlight      []*regexp.Regexp
	MsgHighlightColor *color.Color
	ShowLevel         bool
	LevelBound        Level
	LevelEmoji        bool
	LevelLong         bool
	LevelFirst        bool
	LevelColor        bool
	ShowTag           bool
}

func newPrinter(config *PrinterConfig) *Printer {
	tagFilter := make([]*regexp.Regexp, len(config.Tag.Filter))
	for i, f := range config.Tag.Filter {
		tagFilter[i] = regexp.MustCompile(f)
	}
	tagIgnore := make([]*regexp.Regexp, len(config.Tag.Ignore))
	for i, f := range config.Tag.Ignore {
		tagIgnore[i] = regexp.MustCompile(f)
	}
	msgFilter := make([]*regexp.Regexp, len(config.Msg.Filter))
	for i, f := range config.Msg.Filter {
		msgFilter[i] = regexp.MustCompile(f)
	}
	msgHighlight := make([]*regexp.Regexp, len(config.Msg.Highlight))
	for i, f := range config.Msg.Highlight {
		msgHighlight[i] = regexp.MustCompile(f)
	}
	highlighter := color.New(color.FgBlack, color.BgYellow)
	if config.Msg.HighlightColor != nil {
		highlighter = config.Msg.HighlightColor.Color
	}
	return &Printer{
		TimeFormat:        config.TimeFormat,
		ShowTag:           config.Tag.Show,
		TagFilter:         tagFilter,
		TagIgnore:         tagIgnore,
		MsgFilter:         msgFilter,
		MsgHighlight:      msgHighlight,
		MsgHighlightColor: highlighter,
		ShowLevel:         config.Tag.Show,
		LevelBound:        config.Lvl.Bound,
		LevelEmoji:        config.Lvl.Emoji,
		LevelLong:         config.Lvl.Long,
		LevelFirst:        config.Lvl.First,
		LevelColor:        config.Lvl.Color,
	}
}

func matchesFilter(text string, filters []*regexp.Regexp, empty bool) bool {
	if len(filters) == 0 {
		return empty
	}
	for _, f := range filters {
		if f.MatchString(text) {
			return true
		}
	}
	return false
}

func highlightText(text string, highlights []*regexp.Regexp, highlighter *color.Color) string {
	for _, h := range highlights {
		text = h.ReplaceAllString(text, highlighter.Sprint("${0}"))
	}
	return text
}

func (c *Printer) Print(e *Entry) {
	if e.Level < c.LevelBound {
		return
	}
	if matchesFilter(e.Tag, c.TagIgnore, false) {
		return
	}
	if !matchesFilter(e.Message, c.MsgFilter, true) {
		return
	}
	if !matchesFilter(e.Tag, c.TagFilter, true) {
		return
	}

	if c.LevelFirst {
		c.printLevel(e)
	}
	if len(c.TimeFormat) > 0 {
		fmt.Print(color.BlueString(e.Time.Format(c.TimeFormat + " ")))
	}
	if !c.LevelFirst {
		c.printLevel(e)
	}
	if c.ShowTag {
		color.New(color.FgHiBlack).Printf("%s ", e.Tag)
	}

	fmt.Println(highlightText(e.Message, c.MsgHighlight, c.MsgHighlightColor))
}

func (c *Printer) printLevel(e *Entry) {
	if !c.ShowLevel {
		return
	}
	if c.LevelEmoji {
		if c.LevelFirst {
			fmt.Printf("%s  ", e.Level.Emoji())
		} else {
			fmt.Printf(" %s  ", e.Level.Emoji())
		}
	} else {
		if c.LevelColor {
			if c.LevelLong {
				e.Level.Color().Printf("%s ", e.Level.String())
			} else {
				e.Level.Color().Printf("%c ", e.Level.Rune())
			}
		} else {
			if c.LevelLong {
				fmt.Printf("%s ", e.Level.String())
			} else {
				fmt.Printf("%c ", e.Level.Rune())
			}
		}
	}
}
