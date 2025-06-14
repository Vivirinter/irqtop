package irq

import (
    "bufio"
    "os"
    "strconv"
    "strings"
)

type Record struct {
    IRQ     string
    Name    string
    Delta   int
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
        irq := strings.TrimSuffix(fields[0], ":")

        var nums []int
        for _, f := range fields[1:] {
            if v, err := strconv.Atoi(f); err == nil {
                nums = append(nums, v)
            } else {
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
        name := fields[len(fields)-1]
        current[irq+"|"+name] = total
    }
    var out []Record
    for key, curr := range current {
        delta := curr - r.prev[key]
        if delta < 0 {
            delta = 0
        }
        parts := strings.Split(key, "|")
        out = append(out, Record{IRQ: parts[0], Name: parts[1], Delta: delta})
    }
    r.prev = current
    return out, nil
}
