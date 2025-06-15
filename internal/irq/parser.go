package irq

import (
	"strconv"
	"strings"
)

const keySep = "\x00"

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
