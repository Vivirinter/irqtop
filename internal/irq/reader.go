package irq

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Record struct {
	IRQ   string
	Name  string
	Delta int
}

type Reader struct {
	prev   map[string]int
	cpuIdx int
}

func NewReader(cpuIdx int) *Reader {
	return &Reader{prev: make(map[string]int), cpuIdx: cpuIdx}
}

// Read parses /proc/interrupts and returns delta since previous call.
func (r *Reader) Read() ([]Record, error) {
	file, err := os.Open("/proc/interrupts")
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

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
		irq := strings.TrimSuffix(fields[0], ":")

		var nums []int
		nameStartIdx := 1
		for i, f := range fields[1:] {
			if v, err := strconv.Atoi(f); err == nil {
				nums = append(nums, v)
				nameStartIdx = i + 2
				break
			}
		}
		var total int
		if r.cpuIdx >= 0 {
			if r.cpuIdx < len(nums) {
				total = nums[r.cpuIdx]
			} else {
				total = 0
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

		key := irq + "\x00" + name
		current[key] = total
	}

	var out []Record
	for key, curr := range current {
		delta := curr - r.prev[key]
		// Handle counter resets by taking current value as delta
		if delta < 0 {
			delta = curr
		}
		parts := strings.Split(key, "\x00")
		if len(parts) >= 2 {
			out = append(out, Record{IRQ: parts[0], Name: parts[1], Delta: delta})
		}
	}
	r.prev = current
	return out, nil
}
