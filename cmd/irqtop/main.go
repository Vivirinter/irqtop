package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"

	"github.com/Vivirinter/irqtop/internal/irq"
	"github.com/Vivirinter/irqtop/internal/ui"
)

const (
    adjustPercent = 10 // hotkey changes thresholds by ±10%
    minRed        = 1  // never let thresholds go below 1
    minYellow     = 1
)

type Config struct {
    Interval      time.Duration
    TopN          int
    CPUIdx        int
    RedThresh     int
    YellowThresh  int
    AlertThresh   int
    AlertDur      int
    HidePseudo    bool
    SortBy        string
}

func adjustThreshold(ptr *int, up bool, min int) {
	delta := *ptr * adjustPercent / 100
	if delta == 0 {
		delta = 1
	}
	if up {
		*ptr += delta
	} else {
		*ptr -= delta
		if *ptr < min {
			*ptr = min
		}
	}
}

func handleKey(e *tcell.EventKey, red, yellow *int, sortBy *string, view *ui.View, last *[]irq.Record) bool {
	if e.Key() == tcell.KeyCtrlC || e.Rune() == 'q' {
		return true
	}

	switch r := e.Rune(); r {
	case 'h':
		view.RenderHelp()
		return false
	case 's':
		if *sortBy == ui.SortByDelta {
			*sortBy = ui.SortByName
		} else {
			*sortBy = ui.SortByDelta
		}
	case 'r':
		adjustThreshold(red, true, minRed)
	case 'R':
		adjustThreshold(red, false, minRed)
	case 'y':
		adjustThreshold(yellow, true, minYellow)
	case 'Y':
		adjustThreshold(yellow, false, minYellow)
	default:
		return false
	}

	// apply changes to view and refresh immediately
	view.UpdateConfig(*red, *yellow, *sortBy)
	if *last != nil {
		view.Render(*last)
	}
	return false
}

func main() {
    cfg := &Config{}
    flag.DurationVar(&cfg.Interval, "interval", time.Second, "Refresh interval")
    flag.IntVar(&cfg.TopN, "n", 10, "Show top N IRQs")
    flag.IntVar(&cfg.CPUIdx, "cpu", -1, "CPU index to monitor (-1 = sum)")
    flag.IntVar(&cfg.RedThresh, "red", 1000, "Delta/s threshold for red color")
    flag.IntVar(&cfg.YellowThresh, "yellow", 100, "Delta/s threshold for yellow color")
    flag.IntVar(&cfg.AlertThresh, "alert", 5000, "Delta/s threshold for anomaly beep")
    flag.IntVar(&cfg.AlertDur, "alertdur", 3, "Number of consecutive intervals above alert to trigger beep")
    flag.BoolVar(&cfg.HidePseudo, "hide-pseudo", false, "Hide pseudo IRQs (LOC/CAL/TLB/RES/ERR)")
    flag.StringVar(&cfg.SortBy, "sort", ui.SortByDelta, "Sort by delta or name")
    flag.Parse()

    screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("failed to create screen: %v", err)
	}
	if err = screen.Init(); err != nil {
		log.Fatalf("failed to init screen: %v", err)
	}
	defer screen.Fini()

    reader := irq.NewReader(cfg.CPUIdx)
    view := ui.NewView(screen, cfg.TopN, cfg.RedThresh, cfg.YellowThresh, cfg.HidePseudo, cfg.SortBy)

    det := &Detector{}

    ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	_, _ = reader.Read()

    var lastRecords []irq.Record

    evCh := make(chan tcell.Event, 10)
	go func() {
		for {
			ev := screen.PollEvent()
			evCh <- ev
		}
	}()

	for {
		select {
		case <-ticker.C:
			deltas, err := reader.Read()
			if err != nil {
				log.Printf("read error: %v", err)
				continue
			}
			lastRecords = deltas
			if trigger, top := det.Update(deltas, cfg); trigger {
                if err := screen.Beep(); err != nil {
                    log.Printf("beep error: %v", err)
                }
                view.SetAlert(fmt.Sprintf("IRQ %s %d/s", top.IRQ, top.Delta))
            }
			view.Render(deltas)
		case ev := <-evCh:
			switch e := ev.(type) {
			case *tcell.EventKey:
	                if quit := handleKey(e, &cfg.RedThresh, &cfg.YellowThresh, &cfg.SortBy, view, &lastRecords); quit {
					return
				}
			case *tcell.EventResize:
				screen.Sync()
			}
		}
	}
}
