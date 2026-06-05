package main

import (
	"math"
	"strings"
	"time"
)

// braille dot bit per (row, col) inside a 2x4 cell
var brailleBits = [4][2]byte{
	{0x01, 0x08},
	{0x02, 0x10},
	{0x04, 0x20},
	{0x40, 0x80},
}

type canvas struct {
	w, h int
	dots []byte
}

func newCanvas(w, h int) *canvas {
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	return &canvas{w: w, h: h, dots: make([]byte, w*h)}
}

func (c *canvas) set(x, y int) {
	if x < 0 || y < 0 {
		return
	}
	cx, cy := x/2, y/4
	if cx >= c.w || cy >= c.h {
		return
	}
	c.dots[cy*c.w+cx] |= brailleBits[y%4][x%2]
}

func (c *canvas) line(x0, y0, x1, y1 int) {
	dx, dy := max(x1-x0, x0-x1), -max(y1-y0, y0-y1)
	sx, sy := 1, 1
	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}
	err := dx + dy
	for {
		c.set(x0, y0)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func (c *canvas) String() string {
	var b strings.Builder
	for row := 0; row < c.h; row++ {
		for col := 0; col < c.w; col++ {
			b.WriteRune(rune(0x2800 + int(c.dots[row*c.w+col])))
		}
		if row < c.h-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func renderClock(t time.Time, cellsW, cellsH int) string {
	d := max(12, min(cellsW*2, cellsH*4))
	c := newCanvas((d+1)/2, (d+3)/4)
	cx, cy := float64(d)/2, float64(d)/2
	r := cx - 1

	steps := int(2 * math.Pi * r)
	for i := 0; i < steps; i++ {
		a := float64(i) / float64(steps) * 2 * math.Pi
		c.set(int(math.Round(cx+r*math.Cos(a))), int(math.Round(cy+r*math.Sin(a))))
	}

	for i := 0; i < 12; i++ {
		a := float64(i) / 12 * 2 * math.Pi
		for rr := r * 0.86; rr <= r; rr += 0.5 {
			c.set(int(math.Round(cx+rr*math.Sin(a))), int(math.Round(cy-rr*math.Cos(a))))
		}
	}

	h := float64(t.Hour()%12) + float64(t.Minute())/60
	hand(c, cx, cy, h/12*2*math.Pi, r*0.5)
	hand(c, cx, cy, float64(t.Minute())/60*2*math.Pi, r*0.78)
	hand(c, cx, cy, float64(t.Second())/60*2*math.Pi, r*0.9)
	return c.String()
}

func hand(c *canvas, cx, cy, angle, length float64) {
	x := cx + length*math.Sin(angle)
	y := cy - length*math.Cos(angle)
	c.line(int(math.Round(cx)), int(math.Round(cy)), int(math.Round(x)), int(math.Round(y)))
}
