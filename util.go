package main

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

func drawText(str string, x int, y int, fg termbox.Attribute, bg termbox.Attribute) {
	i := 0
	for _, c := range str {
		termbox.SetCell(x+i, y, c, fg, bg)
		i += runewidth.RuneWidth(c)
	}
}

func drawFrame(x int, y int, width int, height int, fg termbox.Attribute, bg termbox.Attribute) {
	for py := y; py-y < height; py++ {
		if py == y || py-y+1 == height {
			for px := x; px-x < width; px++ {
				if px == x || px-x+1 == width {
					termbox.SetCell(px, py, '+', fg, bg)
				} else {
					termbox.SetCell(px, py, '-', fg, bg)
				}

			}
		} else {
			termbox.SetCell(x, py, '|', fg, bg)
			termbox.SetCell(x+width-1, py, '|', fg, bg)
		}
	}
}

func drawBgLine(x int, y int, width int, color termbox.Attribute) {
	for px := x; px-x < width; px++ {
		termbox.SetCell(px, y, ' ', termbox.ColorDefault, color)
	}
}

func getTermSize() (int, int) {
	return termWidth, termHeight
}

func setTermSize(w, h int) {
	termWidth, termHeight = w, h
}
