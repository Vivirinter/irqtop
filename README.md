# irqtop

A lightweight terminal UI that shows the most active hardware interrupts on a Linux host in **real-time**.
Written in Go, based on `tcell`, intended for quick diagnostics and performance hunting.

---

## Features

* Live per-IRQ delta (interrupts/sec)
* Colour-coded table, thresholds configurable on the fly
* Sort by delta or alphabetically
* Option to hide pseudo-IRQs (LOC, CAL, TLB, RES, ERR)
* Sustained-leader alert with audible beep
* Designed to run either natively or in a privileged container

---

## Quick start

### Docker

```bash
docker run --rm -it --privileged --pid=host \
  -e TERM=xterm-256color \
  <dockerhub_user>/irqtop:latest \
  -interval 1s -n 15 -red 2000 -yellow 200
```

Why `--privileged --pid=host`?  Reading `/proc/interrupts` on the host is only
permitted to a privileged process that shares the PID namespace with the host
on common Docker setups (especially with SELinux or user-namespace remapping).

---

## Keyboard shortcuts

| Key | Action |
|-----|--------|
| q / Ctrl-C | Quit |
| h | Toggle help screen |
| s | Toggle sort (delta ↔ name) |
| r / R | Red threshold +100 / −100 |
| y / Y | Yellow threshold +10 / −10 |

---

## CLI flags

```
-cpu N         monitor only CPU N column (default: sum of all)
-hide-pseudo   hide LOC/CAL/TLB/RES/ERR rows
-interval 1s   refresh interval
-n 15          max rows to display
-red 2000      red threshold (irq/s)
-yellow 200    yellow threshold (irq/s)
-alert 5000    sustained leader alert threshold
-alertdur 3    intervals the leader must persist before alert
```

---

## Development

```bash
# Unit tests
go test ./...

# Static analysis
go vet ./...
golangci-lint run

# Build static binary
CGO_ENABLED=0 go build -ldflags="-s -w" -o irqtop ./cmd/irqtop
```

CI (GitHub Actions) runs on every push / pull-request:

1. Go tests, `go vet`, `golangci-lint` (matrix Go 1.22/1.23/1.24)
2. Multi-arch Docker build & push to Docker Hub (`latest` + commit SHA)

---

## License

MIT © 2025 Vivirinter


Real-time TUI showing the most active hardware interrupts on Linux.

```
irq   delta/s   name
 144     1200   timer
  10      550   AMDI0010:00
  37      210   xhci_hcd
```

## Build

```
go build ./cmd/irqtop
```

## Run

```
./irqtop -interval 1s -n 10       # press 'q' or Ctrl-C to quit
```

Colors:
* red   – ≥1000 irq/s
* yellow – 100-999 irq/s
* green  – <100 irq/s

Requires Go 1.24+ and `tcell/v2` (added via `go mod`).
