// 60-second, 60+fps capability demo for modernretro in 80x24.
//
//	go run ./modernretro
//
// Five scenes drive truecolor + half-block rendering through synchronized output
// (DEC 2026), with a live FPS counter. Ctrl+C exits cleanly.
package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	tm "github.com/angch/termdemo/internal/term"
)

const (
	cols      = 80
	rows      = 24
	targetFPS = 60
)

type scene interface {
	name() string
	duration() float64
	render(b *strings.Builder, st float64, frame int, dt float64)
}

func main() {
	tm.EnterFullscreen()
	defer tm.LeaveFullscreen()

	scenes := []scene{
		&intro{},
		newPlasma(),
		newFire(),
		newMatrix(),
		newStarfield(),
		newSpectrum(),
		&outro{},
	}

	var total float64
	for _, s := range scenes {
		total += s.duration()
	}

	var buf strings.Builder
	buf.Grow(cols * rows * 64)

	frameDur := time.Second / targetFPS
	start := time.Now()
	nextFrame := start
	lastTick := start

	var totalFrames int
	fpsTimer := start
	fpsFrames := 0
	var lastFPS float64

	for {
		now := time.Now()
		elapsed := now.Sub(start).Seconds()
		if elapsed >= total {
			break
		}
		dt := now.Sub(lastTick).Seconds()
		lastTick = now

		acc := 0.0
		var cur scene
		var st float64
		for _, s := range scenes {
			if elapsed < acc+s.duration() {
				cur = s
				st = elapsed - acc
				break
			}
			acc += s.duration()
		}
		if cur == nil {
			break
		}

		buf.Reset()
		buf.WriteString(tm.SyncStart)
		buf.WriteString("\x1b[2J\x1b[H")
		cur.render(&buf, st, totalFrames, dt)
		statusBar(&buf, cur.name(), elapsed, total, lastFPS, totalFrames)
		buf.WriteString(tm.SyncEnd)
		os.Stdout.WriteString(buf.String())

		totalFrames++
		fpsFrames++
		if d := time.Since(fpsTimer); d > 500*time.Millisecond {
			lastFPS = float64(fpsFrames) / d.Seconds()
			fpsFrames = 0
			fpsTimer = time.Now()
		}

		nextFrame = nextFrame.Add(frameDur)
		if d := time.Until(nextFrame); d > 0 {
			time.Sleep(d)
		} else if d < -3*frameDur {
			nextFrame = time.Now()
		}
	}

	// Final summary frame
	elapsed := time.Since(start).Seconds()
	avg := float64(totalFrames) / elapsed
	buf.Reset()
	buf.WriteString(tm.SyncStart + "\x1b[2J\x1b[H" + tm.Reset)
	finalSummary(&buf, totalFrames, elapsed, avg)
	buf.WriteString(tm.SyncEnd)
	os.Stdout.WriteString(buf.String())
	time.Sleep(2 * time.Second)
}

func statusBar(b *strings.Builder, name string, elapsed, total, fps float64, frame int) {
	b.WriteString(tm.MoveTo(24, 1))
	b.WriteString(tm.BG(15, 18, 32))
	b.WriteString(tm.FG(180, 210, 255))
	left := fmt.Sprintf(" %s", name)
	right := fmt.Sprintf("%5.1f fps · %4.1fs/%2.0fs · f%05d ", fps, elapsed, total, frame)
	gap := 80 - len([]rune(left)) - len([]rune(right))
	if gap < 1 {
		gap = 1
	}
	b.WriteString(left)
	b.WriteString(strings.Repeat(" ", gap))
	b.WriteString(right)
	b.WriteString(tm.Reset)
}

func finalSummary(b *strings.Builder, frames int, elapsed, avg float64) {
	lines := []string{
		"",
		"",
		"",
		"",
		"   ╔══════════════════════════════════════════════════════════════════╗",
		"   ║                                                                  ║",
		"   ║           M O D E R N R E T R O · S P E E D · T E S T            ║",
		"   ║                                                                  ║",
		"   ╠══════════════════════════════════════════════════════════════════╣",
		"   ║                                                                  ║",
		fmt.Sprintf("   ║     total frames     : %-41d ║", frames),
		fmt.Sprintf("   ║     elapsed seconds  : %-41.2f ║", elapsed),
		fmt.Sprintf("   ║     average fps      : %-41.2f ║", avg),
		fmt.Sprintf("   ║     cells per frame  : %-41d ║", 80*24),
		fmt.Sprintf("   ║     cells per second : %-41.0f ║", avg*80*24),
		"   ║                                                                  ║",
		"   ╚══════════════════════════════════════════════════════════════════╝",
		"",
	}
	for i, line := range lines {
		b.WriteString(tm.MoveTo(3+i, 1))
		// rainbow gradient the title
		if i == 6 {
			runes := []rune(line)
			for j, r := range runes {
				if j > 3 && j < len(runes)-4 {
					h := math.Mod(float64(j)*5+elapsed*30, 360)
					rr, gg, bb := tm.HSV(h, 0.7, 1.0)
					fmt.Fprintf(b, "%s%s%c", tm.Bold, tm.FG(rr, gg, bb), r)
				} else {
					b.WriteString(tm.FG(150, 200, 255))
					b.WriteRune(r)
				}

			}
			b.WriteString(tm.Reset)
		} else {
			b.WriteString(tm.FG(150, 200, 255))
			b.WriteString(line)
			b.WriteString(tm.Reset)
		}
	}
}

// ───── INTRO ────────────────────────────────────────────────────────────

type intro struct{}

func (intro) name() string      { return "INTRO" }
func (intro) duration() float64 { return 4 }
func (intro) render(b *strings.Builder, st float64, frame int, dt float64) {
	banner := []string{
		"███▄ ▄███▓  ▒█████   ▓█████▄  ▓█████   ██▀███   ███▄    █ ",
		"▓██▒▀█▀ ██▒▒██▒  ██▒ ▒██▀ ██▌ ▒██▀_   ▓██ ▒ ██▒ ██ ▀█   █ ",
		"▓██    ▓██_▒██░  ██▒ ░██   █▌ ░███    ▓██ ░▄█ ▒▓██  ▀█ ██▒",
		"▒██    ▒██ ▒██   ██░ ░▓█▄   ▌ ▒▓█  ▄  ▒██▀▀█▄  ▓██▒  ▐▌██▒",
		"▒██▒   ░██▒░ ████▓▒░ ░▒████▓  ░▒████▒ ░██▓ ▒██▒▒██░   ▓██_",
		"░ ▒░   ░  ░░ ▒░▒░▒░   ▒▒▓  ▒  ░░ ▒░ ░ ░ ▒▓ ░▒▓░░ ▒░   ▒ ▒ ",
		"                                                          ",
		"      ██▀███   ▓█████  ▄▄▄█████▓  ██▀███   ▒█████         ",
		"     ▓██ ▒ ██▒ ▓█   ▀  ▓  ██▒ ▓▒ ▓██ ▒ ██▒ ▒██▒  ██▒      ",
		"     ▓██ ░▄█ ▒ ▒███    ▒ ▓██░ ▒░ ▓██ ░▄█ ▒ ▒██░  ██▒      ",
		"     ▒██▀▀█▄   ▒▓█  ▄  ░ ▓██▓ ░  ▒██▀▀█▄   ▒██   ██░      ",
		"     ░██▓ ▒██▒ ░▒████▒   ▒██▒ ░  ░██▓ ▒██▒ ░ ████▓▒░      ",
		"     ░ ▒▓ ░▒▓░ ░░ ▒░ ░   ▒ ░░    ░ ▒▓ ░▒▓░ ░ ▒░▒░▒░       ",
	}
	startY := 4
	for i, line := range banner {
		b.WriteString(tm.MoveTo(startY+i, (80-len([]rune(line)))/2+1))
		for j, r := range []rune(line) {
			h := math.Mod(float64(j+i*2)*6+st*180, 360)
			rr, gg, bb := tm.HSV(h, 0.8, 1.0)
			fmt.Fprintf(b, "%s%s%c", tm.Bold, tm.FG(rr, gg, bb), r)
		}
	}
	b.WriteString(tm.Reset)
	sub := "— CAPABILITY DEMO · 80×24 · 60 fps target —"
	b.WriteString(tm.MoveTo(startY+16, (80-len([]rune(sub)))/2+1))
	alpha := math.Min(1, st)
	c := int(220 * alpha)
	fmt.Fprintf(b, "%s%s%s", tm.Italic, tm.FG(c, c, c), sub)
	b.WriteString(tm.Reset)
}

// ───── PLASMA ───────────────────────────────────────────────────────────

type plasma struct{}

func newPlasma() *plasma          { return &plasma{} }
func (*plasma) name() string      { return "PLASMA · 80×48 half-block · truecolor BG+FG" }
func (*plasma) duration() float64 { return 12 }

func plasmaColor(x, y int, t float64) (int, int, int) {
	fx, fy := float64(x), float64(y)
	v := math.Sin(fx*0.15 + t*2)
	v += math.Sin(fy*0.18 + t*1.7)
	v += math.Sin((fx+fy)*0.10 + t*1.3)
	cx, cy := fx-40, fy-24
	v += math.Sin(math.Sqrt(cx*cx+cy*cy)*0.15 + t*1.9)
	v /= 4
	h := (v+1)*180 + t*30
	return tm.HSV(h, 0.85, 0.95)
}

func (*plasma) render(b *strings.Builder, st float64, frame int, dt float64) {
	var lfr, lfg, lfb, lbr, lbg, lbb int = -1, -1, -1, -1, -1, -1
	for cy := range 24 {
		py1 := cy * 2
		py2 := cy*2 + 1
		for x := range 80 {
			r1, g1, b1 := plasmaColor(x, py1, st)
			r2, g2, b2 := plasmaColor(x, py2, st)
			if r1 != lfr || g1 != lfg || b1 != lfb || r2 != lbr || g2 != lbg || b2 != lbb {
				fmt.Fprintf(b, "\x1b[38;2;%d;%d;%d;48;2;%d;%d;%dm", r1, g1, b1, r2, g2, b2)
				lfr, lfg, lfb, lbr, lbg, lbb = r1, g1, b1, r2, g2, b2
			}
			b.WriteRune('▀')
		}
		b.WriteString(tm.Reset)
		lfr, lfg, lfb, lbr, lbg, lbb = -1, -1, -1, -1, -1, -1
		if cy < 23 {
			b.WriteByte('\n')
		}
	}
}

// ───── DOOM FIRE ────────────────────────────────────────────────────────

type fire struct {
	buf [48][80]uint8
}

var firePalette [37][3]int

func init() {
	for i := 0; i < 37; i++ {
		f := float64(i) / 36
		switch {
		case f < 0.25:
			tt := f / 0.25
			firePalette[i] = [3]int{int(tt * 180), 0, int(tt * 20)}
		case f < 0.5:
			tt := (f - 0.25) / 0.25
			firePalette[i] = [3]int{180 + int(75*tt), int(100 * tt), int(20 * (1 - tt))}
		case f < 0.75:
			tt := (f - 0.5) / 0.25
			firePalette[i] = [3]int{255, 100 + int(155*tt), 0}
		default:
			tt := (f - 0.75) / 0.25
			firePalette[i] = [3]int{255, 255, int(255 * tt)}
		}
	}
}

func newFire() *fire {
	f := &fire{}
	for x := range 80 {
		f.buf[47][x] = 36
	}
	return f
}

func (*fire) name() string      { return "DOOM FIRE · cellular automaton · half-block" }
func (*fire) duration() float64 { return 12 }

func (f *fire) render(b *strings.Builder, st float64, frame int, dt float64) {
	for x := range 80 {
		f.buf[47][x] = 36
	}
	for y := range 47 {
		for x := range 80 {
			src := f.buf[y+1][x]
			r := rand.Intn(3)
			off := rand.Intn(3) - 1
			xx := max(x+off, 0)
			if xx >= 80 {
				xx = 79
			}
			if int(src)-r < 0 {
				f.buf[y][xx] = 0
			} else {
				f.buf[y][xx] = src - uint8(r)
			}
		}
	}
	var lfr, lfg, lfb, lbr, lbg, lbb int = -1, -1, -1, -1, -1, -1
	for cy := range 24 {
		for x := range 80 {
			tc := firePalette[f.buf[cy*2][x]]
			bc := firePalette[f.buf[cy*2+1][x]]
			if tc[0] != lfr || tc[1] != lfg || tc[2] != lfb || bc[0] != lbr || bc[1] != lbg || bc[2] != lbb {
				fmt.Fprintf(b, "\x1b[38;2;%d;%d;%d;48;2;%d;%d;%dm",
					tc[0], tc[1], tc[2], bc[0], bc[1], bc[2])
				lfr, lfg, lfb = tc[0], tc[1], tc[2]
				lbr, lbg, lbb = bc[0], bc[1], bc[2]
			}
			b.WriteRune('▀')
		}
		b.WriteString(tm.Reset)
		lfr, lfg, lfb, lbr, lbg, lbb = -1, -1, -1, -1, -1, -1
		if cy < 23 {
			b.WriteByte('\n')
		}
	}
}

// ───── MATRIX RAIN ──────────────────────────────────────────────────────

type matrix struct {
	heads [80]float64
	speed [80]float64
	cells [24][80]rune
	age   [24][80]int
}

var matrixChars = []rune("ｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎ0123456789ABCDEF<>/\\|=+-*[]{}!?")

func newMatrix() *matrix {
	m := &matrix{}
	for x := 0; x < 80; x++ {
		m.heads[x] = -rand.Float64() * 30
		m.speed[x] = 8 + rand.Float64()*20
	}
	for y := 0; y < 24; y++ {
		for x := 0; x < 80; x++ {
			m.age[y][x] = 9999
			m.cells[y][x] = ' '
		}
	}
	return m
}

func (*matrix) name() string      { return "MATRIX RAIN · katakana + ascii · per-cell fade" }
func (*matrix) duration() float64 { return 12 }

func (m *matrix) render(b *strings.Builder, st float64, frame int, dt float64) {
	if dt <= 0 || dt > 0.5 {
		dt = 1.0 / 60
	}
	for y := 0; y < 24; y++ {
		for x := 0; x < 80; x++ {
			if m.age[y][x] < 9999 {
				m.age[y][x]++
			}
		}
	}
	for x := 0; x < 80; x++ {
		prev := m.heads[x]
		m.heads[x] += m.speed[x] * dt
		// fill every cell the head passed through this frame
		for yi := int(math.Max(0, prev)); yi <= int(m.heads[x]); yi++ {
			if yi >= 0 && yi < 24 {
				m.cells[yi][x] = matrixChars[rand.Intn(len(matrixChars))]
				m.age[yi][x] = 0
			}
		}
		if int(m.heads[x]) >= 28 {
			m.heads[x] = -rand.Float64() * 30
			m.speed[x] = 8 + rand.Float64()*22
		}
		// occasionally mutate trailing chars
		if rand.Intn(8) == 0 {
			yy := rand.Intn(24)
			if m.age[yy][x] < 40 {
				m.cells[yy][x] = matrixChars[rand.Intn(len(matrixChars))]
			}
		}
	}
	// solid black bg for the whole scene to overwrite any leftover pixels
	b.WriteString(tm.BG(0, 0, 0))
	var lr, lg, lbl int = -1, -1, -1
	for y := range 24 {
		for x := range 80 {
			a := m.age[y][x]
			var r, g, bl int
			ch := m.cells[y][x]
			if a >= 60 || ch == ' ' {
				if lr != 0 || lg != 0 || lbl != 0 {
					b.WriteString(tm.FG(0, 0, 0))
					lr, lg, lbl = 0, 0, 0
				}
				b.WriteByte(' ')
				continue
			}
			if a == 0 {
				r, g, bl = 230, 255, 230
			} else {
				f := 1.0 - float64(a)/60
				g = 80 + int(175*f)
				r = int(25 * f)
				bl = int(50 * f)
			}
			if r != lr || g != lg || bl != lbl {
				b.WriteString(tm.FG(r, g, bl))
				lr, lg, lbl = r, g, bl
			}
			b.WriteRune(ch)
		}
		b.WriteString(tm.Reset)
		b.WriteString(tm.BG(0, 0, 0))
		lr, lg, lbl = -1, -1, -1
		if y < 23 {
			b.WriteByte('\n')
		}
	}
	b.WriteString(tm.Reset)
}

// ───── STARFIELD ────────────────────────────────────────────────────────

type starfield struct {
	x, y, z [240]float64
}

func newStarfield() *starfield {
	s := &starfield{}
	for i := range len(s.x) {
		s.x[i] = (rand.Float64()*2 - 1) * 80
		s.y[i] = (rand.Float64()*2 - 1) * 48
		s.z[i] = rand.Float64()*99 + 1
	}
	return s
}

func (*starfield) name() string      { return "STARFIELD · 240 stars · 3D perspective" }
func (*starfield) duration() float64 { return 12 }

func (s *starfield) render(b *strings.Builder, st float64, frame int, dt float64) {
	if dt <= 0 || dt > 0.5 {
		dt = 1.0 / 60
	}
	var grid [48][80]uint8
	speed := 35.0 + 25*math.Sin(st*0.4)
	for i := range len(s.x) {
		s.z[i] -= dt * speed
		if s.z[i] <= 0.5 {
			s.x[i] = (rand.Float64()*2 - 1) * 80
			s.y[i] = (rand.Float64()*2 - 1) * 48
			s.z[i] = 100
		}
		px := int(s.x[i]/s.z[i]*60) + 40
		py := int(s.y[i]/s.z[i]*36) + 24
		if px >= 0 && px < 80 && py >= 0 && py < 48 {
			br := max(255-int(s.z[i]*2.4), 0)
			if grid[py][px] < uint8(br) {
				grid[py][px] = uint8(br)
			}
			// streak trailing pixel for fast stars
			if s.z[i] < 30 && py+1 < 48 {
				dim := br / 3
				if grid[py+1][px] < uint8(dim) {
					grid[py+1][px] = uint8(dim)
				}
			}
		}
	}
	var lr, lg, lbl, lbr2, lbg2, lbb2 int = -1, -1, -1, -1, -1, -1
	for cy := range 24 {
		for x := range 80 {
			ti := int(grid[cy*2][x])
			bi := int(grid[cy*2+1][x])
			// slight blue tint
			tr, tg, tbl := ti, ti, clamp_u8(ti+20)
			br2, bg2, bb2 := bi, bi, clamp_u8(bi+20)
			if tr != lr || tg != lg || tbl != lbl || br2 != lbr2 || bg2 != lbg2 || bb2 != lbb2 {
				fmt.Fprintf(b, "\x1b[38;2;%d;%d;%d;48;2;%d;%d;%dm", tr, tg, tbl, br2, bg2, bb2)
				lr, lg, lbl = tr, tg, tbl
				lbr2, lbg2, lbb2 = br2, bg2, bb2
			}
			b.WriteRune('▀')
		}
		b.WriteString(tm.Reset)
		lr, lg, lbl, lbr2, lbg2, lbb2 = -1, -1, -1, -1, -1, -1
		if cy < 23 {
			b.WriteByte('\n')
		}
	}
}

func clamp_u8(v int) int {
	return min(max(v, 0), 255)
}

// ───── SPECTRUM ─────────────────────────────────────────────────────────

type spectrum struct{}

func newSpectrum() *spectrum        { return &spectrum{} }
func (*spectrum) name() string      { return "SPECTRUM · bars + rainbow text + marquee" }
func (*spectrum) duration() float64 { return 12 }

func (*spectrum) render(b *strings.Builder, st float64, frame int, dt float64) {
	// Top title (row 2), rainbow
	title := "─── M O D E R N R E T R O · 80×24 · 24-bit · 60+ fps ───"
	b.WriteString(tm.MoveTo(2, (80-len([]rune(title)))/2+1))
	for i, r := range title {
		h := math.Mod(float64(i)*10+st*120, 360)
		rr, gg, bb := tm.HSV(h, 0.85, 1.0)
		fmt.Fprintf(b, "%s%s%c", tm.Bold, tm.FG(rr, gg, bb), r)
	}
	b.WriteString(tm.Reset)

	// Spectrum bars: 80 cols, height 0..16, base at row 20
	heights := make([]int, 80)
	for x := range 80 {
		v := math.Sin(float64(x)*0.15+st*4) * 0.4
		v += math.Sin(float64(x)*0.08+st*2.3) * 0.4
		v += math.Sin(float64(x)*0.30+st*5.5) * 0.2
		v = (v + 1) / 2
		heights[x] = int(v * 16)
	}
	for row := 5; row <= 20; row++ {
		b.WriteString(tm.MoveTo(row, 1))
		for x := range 80 {
			barTop := 20 - heights[x]
			if row >= barTop {
				frac := float64(20-row) / 16
				h := 240 - frac*240
				rr, gg, bb := tm.HSV(h, 0.9, 0.85+0.15*frac)
				fmt.Fprintf(b, "%s█", tm.FG(rr, gg, bb))
			} else {
				b.WriteByte(' ')
			}
		}
		b.WriteString(tm.Reset)
	}

	// Marquee row 23
	msg := "  ★  modernretro · synchronized output · OSC 8 hyperlinks · kitty graphics · ligatures · unicode 15 · GPU acceleration · 24-bit truecolor · "
	full := strings.Repeat(msg, 4)
	runes := []rune(full)
	offset := int(st*30) % len(runes)
	view := make([]rune, 0, 80)
	for i := range 80 {
		view = append(view, runes[(offset+i)%len(runes)])
	}
	b.WriteString(tm.MoveTo(23, 1))
	for i, r := range view {
		h := math.Mod(float64(i)*4+st*60, 360)
		rr, gg, bb := tm.HSV(h, 0.7, 1.0)
		fmt.Fprintf(b, "%s%c", tm.FG(rr, gg, bb), r)
	}
	b.WriteString(tm.Reset)
}

// ───── OUTRO ────────────────────────────────────────────────────────────

type outro struct{}

func (outro) name() string      { return "OUTRO" }
func (outro) duration() float64 { return 4 }
func (outro) render(b *strings.Builder, st float64, frame int, dt float64) {
	// message
	msg := "thank you · 80×24 · 60 fps · modernretro"
	b.WriteString(tm.MoveTo(14, (80-len([]rune(msg)))/2+1))
	for i, r := range msg {
		h := math.Mod(float64(i)*8+st*180, 360)
		rr, gg, bb := tm.HSV(h, 0.7, 1.0)
		fmt.Fprintf(b, "%s%s%c", tm.Bold, tm.FG(rr, gg, bb), r)
	}
	b.WriteString(tm.Reset)
}
