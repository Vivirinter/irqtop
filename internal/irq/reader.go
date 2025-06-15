package irq

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

const keySep = "\x00"

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

func ParseLines(lines []string, cpuIdx int) map[string]int {
	current := make(map[string]int)
	for idx, line := range lines {
		if idx == 0 {
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		key, total, ok := parseLine(fields, cpuIdx)
		if !ok {
			continue
		}
		current[key] = total
	}
	return current
}

func (r *Reader) Read() ([]Record, error) {
	file, err := os.Open("/proc/interrupts")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil, scanner.Err()
	}

	var lines []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	current := ParseLines(lines, r.cpuIdx)

	var out []Record
	for key, curr := range current {
		delta := curr - r.prev[key]
		if delta < 0 {
			delta = curr
		}
		if r.isFirstRun {
			delta = 0
		}
		parts := strings.Split(key, keySep)
		if len(parts) >= 2 {
			out = append(out, Record{IRQ: parts[0], Name: parts[1], Delta: delta})
		}
	}
	r.prev = current
	r.isFirstRun = false
	return out, nil
}

func parseLine(fields []string, cpuIdx int) (string, int, bool) {
	if len(fields) < 2 {
		return "", 0, false
	}

	irq := strings.TrimSuffix(fields[0], ":")

	var nums []int
	nameStartIdx := 1
	for i, f := range fields[1:] {
		if v, err := strconv.Atoi(f); err == nil {
			nums = append(nums, v)
			nameStartIdx = i + 2
		} else {
			break
		}
	}

	var total int
	if cpuIdx >= 0 {
		if cpuIdx < len(nums) {
			total = nums[cpuIdx]
		}
	} else {
		for _, v := range nums {
			total += v
		}
	}

	var name string
	if nameStartIdx < len(fields) {
		name = strings.Join(fields[nameStartIdx:], " ")
	}
	if name == "" {
		name = "unknown"
	}

	key := irq + keySep + name
	return key, total, true
}
