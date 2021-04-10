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

type cursedNyandraid struct {
	followBottom bool
	list         *widgets.List
	grid         *ui.Grid
	ring         *Entring
	entries      <-chan *lc.Entry
}

func NewCursedNyandraid() *cursedNyandraid {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	list := widgets.NewList()
	list.Title = getTitleString(initialFilterLevel, "")
	list.WrapText = true
	list.SelectedRowStyle = ui.NewStyle(ui.ColorWhite, ui.ColorBlue)

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(1.0,
			ui.NewCol(1.0, list),
		),
	)
	entring := newEntring(1_000)
	return &cursedNyandraid{
		followBottom: true,
		list:         list,
		grid:         grid,
		ring:         entring,
		entries:      lc.NyanForever(),
	}
}

func (n *cursedNyandraid) End() {
	ui.Close()
}

func (n *cursedNyandraid) start() {
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
			n.grid.SetRect(0, 0, width, height)
		case "<Backspace>":
			if len(filter) == 0 {
				continue
			}
			filter = filter[:len(filter)-1]
			n.list.Title = getTitleString(level, filter)
		case "<C-k>":
			n.followBottom = false
			n.list.ScrollPageUp()
		case "<C-j>":
			n.list.ScrollPageDown()
		case "<Home>":
			n.followBottom = false
			n.list.ScrollTop()
		case "<End>":
			n.followBottom = true
			n.list.ScrollBottom()
		case "<Up>", "<MouseWheelUp>":
			n.followBottom = false
			n.list.ScrollUp()
		case "<Down>", "<MouseWheelDown>":
			n.list.ScrollDown()
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

func (n *cursedNyandraid) render() {
	ui.Render(n.grid)
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
	curse := NewCursedNyandraid()
	defer curse.End()
	curse.start()
}
