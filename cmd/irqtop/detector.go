package main

import "github.com/Vivirinter/irqtop/internal/irq"

type Detector struct {
    leader irq.Record
    cnt    int
}

func (d *Detector) Update(recs []irq.Record, cfg *Config) (bool, irq.Record) {
    // find record with maximum delta
    var top irq.Record
    for _, r := range recs {
        if r.Delta > top.Delta {
            top = r
        }
    }
    if top.Delta < cfg.AlertThresh {
        d.cnt = 0
        return false, top
    }
    if top.IRQ == d.leader.IRQ {
        d.cnt++
    } else {
        d.leader = top
        d.cnt = 1
    }
    if d.cnt >= cfg.AlertDur {
        return true, top
    }
    return false, top
}
