package ui

import (
	"fmt"
	"sort"

	"github.com/Vivirinter/irqtop/internal/irq"
	"github.com/gdamore/tcell/v2"
)

const (
	SortByDelta = "delta"
	SortByName  = "name"
)

var pseudoIRQs = map[string]struct{}{
	"LOC": {}, "CAL": {}, "RES": {}, "TLB": {}, "ERR": {},
}

type View struct {
	screen       tcell.Screen
	topN         int
	redThresh    int
	yellowThresh int
	hidePseudo   bool
	sortBy       string
	alertMsg     string
	alertTicks   int
}

func NewView(s tcell.Screen, n, red, yellow int, hide bool, sortBy string) *View {
	return &View{screen: s, topN: n, redThresh: red, yellowThresh: yellow, hidePseudo: hide, sortBy: sortBy}
}

func (v *View) UpdateConfig(red, yellow int, sortBy string) {
	v.redThresh = red
	v.yellowThresh = yellow
	v.sortBy = sortBy
}

func (v *View) RenderHelp() {
	v.screen.Clear()
	lines := []string{
		"irqtop - IRQ activity monitor",
		"",
		"q / Ctrl-C : quit",
		"h         : toggle help",
		"s         : toggle sort (delta/name)",
		"r/R       : red threshold +100 / -100",
		"y/Y       : yellow threshold +10 / -10",
		"-cpu N    : show single CPU column",
		"-hide-pseudo : hide LOC/CAL/TLB/RES/ERR entries",
	}
	for i, l := range lines {
		v.printLine(i, l, tcell.StyleDefault)
	}
	v.screen.Show()
}

func (v *View) Render(records []irq.Record) {
	// optional filter pseudo-IRQs
	if v.hidePseudo {
		filtered := records[:0]
		for _, r := range records {
			if _, ok := pseudoIRQs[r.IRQ]; ok {
				continue
			}
			filtered = append(filtered, r)
		}
		records = filtered
	}

	if v.sortBy == SortByName {
		sort.Slice(records, func(i, j int) bool { return records[i].Delta > records[j].Delta })
		if len(records) > v.topN {
			records = records[:v.topN]
		}
		sort.Slice(records, func(i, j int) bool { return records[i].Name < records[j].Name })
	} else {
		sort.Slice(records, func(i, j int) bool { return records[i].Delta > records[j].Delta })
		if len(records) > v.topN {
			records = records[:v.topN]
		}
	}

	totalDelta := 0
	for _, r := range records {
		totalDelta += r.Delta
	}

	v.screen.Clear()
	header := "IRQ   Delta/s  Name"
	v.printLine(0, header, tcell.StyleDefault.Bold(true))

	for i, rec := range records {
		style := tcell.StyleDefault
		switch {
		case rec.Delta >= v.redThresh:
			style = style.Foreground(tcell.ColorRed)
		case rec.Delta >= v.yellowThresh:
			style = style.Foreground(tcell.ColorYellow)
		default:
			style = style.Foreground(tcell.ColorGreen)
		}
		line := fmt.Sprintf("%-5s %-8d %s", rec.IRQ, rec.Delta, rec.Name)
		v.printLine(i+1, line, style)
	}
	if v.alertTicks > 0 {
		v.alertTicks--
		v.printLine(v.topN+2, v.alertMsg, tcell.StyleDefault.Foreground(tcell.ColorRed))
	} else {
		status := fmt.Sprintf("sort:%s red:%d yellow:%d total:%d (r/R,y/Y adjust)", v.sortBy, v.redThresh, v.yellowThresh, totalDelta)
		v.printLine(v.topN+2, status, tcell.StyleDefault)
	}
	v.screen.Show()
}

func (v *View) printLine(row int, text string, style tcell.Style) {
	for col, r := range text {
		v.screen.SetContent(col, row, r, nil, style)
	}
}

func (v *View) SetAlert(msg string) {
	v.alertMsg = msg
	v.alertTicks = 3 // show for three refreshes
}
