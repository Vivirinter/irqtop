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

	t.Run("sum_mode", func(t *testing.T) {
		b0 := ParseLines(before, -1)
		b1 := ParseLines(after, -1)
		key := "25" + keySep + "nvme"

		expected := 9
		actual := b1[key] - b0[key]
		if actual != expected {
			t.Errorf("sum delta wrong: got %d, want %d", actual, expected)
		}
	})

	t.Run("cpu1_mode", func(t *testing.T) {
		b0cpu1 := ParseLines(before, 1)
		b1cpu1 := ParseLines(after, 1)
		key := "25" + keySep + "nvme"

		expected := 4
		actual := b1cpu1[key] - b0cpu1[key]
		if actual != expected {
			t.Errorf("cpu1 delta wrong: got %d, want %d", actual, expected)
		}
	})
}
