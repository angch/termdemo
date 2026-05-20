package term

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	Reset       = "\x1b[0m"
	Clear       = "\x1b[2J\x1b[H"
	HideCursor  = "\x1b[?25l"
	ShowCursor  = "\x1b[?25h"
	AltScreen   = "\x1b[?1049h"
	MainScreen  = "\x1b[?1049l"
	SyncStart   = "\x1b[?2026h" // ghostty supports DEC 2026
	SyncEnd     = "\x1b[?2026l"
	Home        = "\x1b[H"
	Bold        = "\x1b[1m"
	Dim         = "\x1b[2m"
	Italic      = "\x1b[3m"
	Underline   = "\x1b[4m"
	Blink       = "\x1b[5m"
	Reverse     = "\x1b[7m"
	Strike      = "\x1b[9m"
	UnderDouble = "\x1b[4:2m"
	UnderCurly  = "\x1b[4:3m"
	UnderDotted = "\x1b[4:4m"
	UnderDashed = "\x1b[4:5m"
)

func MoveTo(row, col int) string { return fmt.Sprintf("\x1b[%d;%dH", row, col) }
func FG(r, g, b int) string      { return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b) }
func BG(r, g, b int) string      { return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b) }
func FG256(c int) string         { return fmt.Sprintf("\x1b[38;5;%dm", c) }
func BG256(c int) string         { return fmt.Sprintf("\x1b[48;5;%dm", c) }

// HSV: h in [0,360), s,v in [0,1]. Returns r,g,b in [0,255].
func HSV(h, s, v float64) (int, int, int) {
	h = math.Mod(math.Mod(h, 360)+360, 360)
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := v - c
	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}
	return clamp_u8(r + m), clamp_u8(g + m), clamp_u8(b + m)
}

// clamp_u8 converts f in [0,1] to int in [0,255], clamping out-of-range values.
func clamp_u8(f float64) int {
	// can't use min or max, because different types.
	if f <= 0 {
		return 0
	}
	if f >= 1 {
		return 255
	}
	return int(f * 255)
}

// Lerp linearly blends a→b by t in [0,1].
func Lerp(a, b, t float64) float64 { return a + (b-a)*t }

// EnterFullscreen swaps to alt screen, hides cursor, installs SIGINT handler that
// restores the terminal and exits cleanly.
func EnterFullscreen() {
	os.Stdout.WriteString(AltScreen + HideCursor + Clear)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		LeaveFullscreen()
		os.Exit(0)
	}()
}

func LeaveFullscreen() {
	os.Stdout.WriteString(Reset + ShowCursor + MainScreen)
}

// PadCenter centers s within width n, padded with spaces.
func PadCenter(s string, n int) string {
	visible := visibleLen(s)
	if visible >= n {
		return s
	}
	left := (n - visible) / 2
	right := n - visible - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

// visibleLen counts runes excluding ANSI CSI sequences. Approximate but works for our use.
func visibleLen(s string) int {
	n := 0
	inEsc := false
	for _, r := range s {
		if inEsc {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEsc = false
			}
			continue
		}
		if r == 0x1b {
			inEsc = true
			continue
		}
		n++
	}
	return n
}
