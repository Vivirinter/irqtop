package irq

import "testing"

func TestParseLines(t *testing.T) {
    before := []string{
        "            CPU0  CPU1",
        "25:          10    20   nvme",
        "85:           5     5   usb",
    }
    after := []string{
        "            CPU0  CPU1",
        "25:          15    24   nvme",
        "85:           7     5   usb",
    }

    b0 := ParseLines(before, -1)
    b1 := ParseLines(after, -1)
    deltaSum := make(map[string]int)
    for k, v := range b1 {
        deltaSum[k] = v - b0[k]
    }
    if deltaSum["25|nvme"] != 9 {
        t.Fatalf("sum mode wrong, got %d", deltaSum["25|nvme"])
    }

    b0cpu1 := ParseLines(before, 1)
    b1cpu1 := ParseLines(after, 1)
    if (b1cpu1["25|nvme"] - b0cpu1["25|nvme"]) != 4 {
        t.Fatalf("cpu1 delta wrong")
    }
}
