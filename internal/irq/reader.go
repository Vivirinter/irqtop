package irq

import (
	"bufio"
	"os"
	"strings"
)

type Record struct {
	IRQ   string
	Name  string
	Delta int
}

type Reader struct {
	prev       map[string]int
	cpuIdx     int
	isFirstRun bool
}

func NewReader(cpuIdx int) *Reader {
	return &Reader{prev: make(map[string]int), cpuIdx: cpuIdx, isFirstRun: true}
}

func (r *Reader) Read() ([]Record, error) {
	file, err := os.Open("/proc/interrupts")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil, scanner.Err()
	}

	current := make(map[string]int)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		key, total, ok := parseLine(fields, r.cpuIdx)
		if !ok {
			continue
		}
		current[key] = total
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var out []Record
	for key, curr := range current {
		delta := curr - r.prev[key]
		if delta < 0 {
			delta = curr
		}
		if r.isFirstRun {
			delta = 0
		}
		parts := strings.Split(key, "\x00")
		if len(parts) >= 2 {
			out = append(out, Record{IRQ: parts[0], Name: parts[1], Delta: delta})
		}
	}
	r.prev = current
	r.isFirstRun = false
	return out, nil
}
