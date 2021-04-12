package main

import (
	"fmt"
	"log"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	lc "wesalvaro.com/nyandroidlogcat"
)

const initialFilterLevel = lc.Warning

type cursedNyandroid struct {
	followBottom bool
	list         *widgets.List
	ring         *Entring
	entries      <-chan *lc.Entry
}

func NewCursedNyandroid() *cursedNyandroid {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	w, h := ui.TerminalDimensions()

	list := widgets.NewList()
	list.Title = getTitleString(initialFilterLevel, "")
	list.WrapText = true
	list.SetRect(0, 0, w, h)
	list.SelectedRowStyle = ui.NewStyle(ui.ColorWhite, ui.ColorBlue)

	ui.Render(list)
	return &cursedNyandroid{
		followBottom: true,
		list:         list,
		ring:         newEntring(1_000),
		entries:      lc.NyanForever(),
	}
}

func (n *cursedNyandroid) End() {
	ui.Close()
}

func (n *cursedNyandroid) start() {
	level := initialFilterLevel
	filter := ""
	go func() {
		for e := range n.entries {
			if e.Level < level {
				continue
			}
			if len(filter) > 0 && !strings.Contains(e.Message, filter) {
				continue
			}
			n.ring.Append(e)
			n.list.Rows = n.ring.ToList()
			if n.followBottom {
				n.list.ScrollBottom()
			}
			n.render()
		}
	}()
	for e := range ui.PollEvents() {
		switch e.ID {
		case "<Resize>":
			payload := e.Payload.(ui.Resize)
			width, height := payload.Width, payload.Height
			n.list.SetRect(0, 0, width, height)
		case "<Backspace>":
			if len(filter) == 0 {
				continue
			}
			filter = filter[:len(filter)-1]
			n.list.Title = getTitleString(level, filter)
		case "<C-k>":
			n.list.ScrollPageUp()
			n.followBottom = false
		case "<C-j>":
			n.list.ScrollPageDown()
			n.followBottom = n.isAtBottom()
		case "<Home>":
			n.list.ScrollTop()
			n.followBottom = false
		case "<End>":
			n.list.ScrollBottom()
			n.followBottom = true
		case "<Up>", "<MouseWheelUp>":
			n.list.ScrollUp()
			n.followBottom = false
		case "<Down>", "<MouseWheelDown>":
			n.list.ScrollDown()
			n.followBottom = n.isAtBottom()
		case "<Escape>", "<C-c>":
			return
		case "<Right>":
			level = level.Next()
			n.list.Title = getTitleString(level, filter)
		case "<Left>":
			level = level.Prev()
			n.list.Title = getTitleString(level, filter)
		case "<Space>":
			filter = filter + " "
			n.list.Title = getTitleString(level, filter)
		}
		if (e.ID[0] >= 'A' && e.ID[0] <= 'z') || (e.ID[0] >= '0' && e.ID[0] <= '9') {
			filter = filter + string(e.ID[0])
			n.list.Title = getTitleString(level, filter)
		}
		n.render()
	}
}

func (n *cursedNyandroid) isAtBottom() bool {
	return n.list.SelectedRow == (len(n.list.Rows) - 1)
}

func (n *cursedNyandroid) render() {
	ui.Render(n.list)
}

func getTitleString(level lc.Level, filter string) string {
	if len(filter) == 0 {
		return fmt.Sprintf(
			" Level: %s (type to filter messages, 0-5 to filter levels, esc to exit) ",
			level.String())
	} else {
		return fmt.Sprintf(" Level: %s | *%s* ", level.String(), filter)
	}
}

func main() {
	curse := NewCursedNyandroid()
	defer curse.End()
	curse.start()
}
