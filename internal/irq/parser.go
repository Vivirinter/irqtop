package irq

import (
    "strconv"
    "strings"
)

func ParseLines(lines []string, cpuIdx int) map[string]int {
    current := make(map[string]int)
    if len(lines) == 0 {
        return current
    }
    for i, line := range lines {
        if i == 0 {
            continue
        }
        line = strings.TrimSpace(line)
        if line == "" {
            continue
        }
        fields := strings.Fields(line)
        if len(fields) < 2 {
            continue
        }
        irq := strings.TrimSuffix(fields[0], ":")

        // collect numeric CPU columns until non-numeric field encountered
        var nums []int
        for _, f := range fields[1:] {
            if v, err := strconv.Atoi(f); err == nil {
                nums = append(nums, v)
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
        name := fields[len(fields)-1]
        current[irq+"|"+name] = total
    }

    return current
}
