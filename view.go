package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"github.com/tSU-RooT/diskinfo/atasmart"
)

type view struct {
	disks      []*atasmart.SkDisk
	page       int
	scroll     int
	scrollable bool
	exitCh     chan struct{}
}

func initView() view {
	v := view{
		disks:  nil,
		page:   0,
		exitCh: make(chan struct{}, 5),
	}
	return v
}

func (v *view) paintScreen() {
	w, h := getTermSize()
	// draw frame
	drawFrame(0, 0, w, 4, termbox.ColorBlue, termbox.ColorDefault)
	drawFrame(0, 3, w, 5, termbox.ColorBlue, termbox.ColorDefault)
	drawFrame(0, 7, w, h-7, termbox.ColorBlue, termbox.ColorDefault)
	v.paintSelecter()
	v.paintData(v.disks[v.page])
	v.paintAttr(v.disks[v.page])
}

func (v *view) paintSelecter() {
	for i, d := range v.disks {
		x := 2 + i*12
		fg := termbox.ColorWhite
		bg := termbox.ColorDefault
		if v, err := d.GetSize(); err == nil {
			var unit float64
			unit = (float64)(v / 1024.0)
			c := 0
			u := []string{"KB", "MB", "GB", "TB", "PB"}
			for unit >= 1000 && c < len(u)-1 {
				unit /= 1024.0
				c++
			}
			t := fmt.Sprintf("%.2f[%s]", unit, u[c])
			drawText(t, x, 1, fg, bg)
		}
		if i == v.page {
			bg = termbox.ColorCyan
		}
		p := fmt.Sprintf("%- 10s", d.Path)
		drawText(p, x, 2, fg, bg)
	}

}

func (v *view) paintData(disk *atasmart.SkDisk) {
	y := 4
	text := fmt.Sprintf("Name: %s", disk.Data.Model)
	drawText(text, 2, y, termbox.ColorWhite, termbox.ColorDefault)
	y++
	text = fmt.Sprintf("Firmware: %s", disk.Data.Firmware)
	drawText(text, 2, y, termbox.ColorWhite, termbox.ColorDefault)
	y++
	text = fmt.Sprintf("Serial: %s", disk.Data.Serial)
	drawText(text, 2, y, termbox.ColorWhite, termbox.ColorDefault)

	y = 4
	if v, err := disk.GetSmartPowerCycle(); err == nil {
		text = fmt.Sprintf("Power On Count: %d", v)
		drawText(text, 42, y, termbox.ColorWhite, termbox.ColorDefault)
		y++
	}
	if v, err := disk.GetSmartTemperature(); err == nil {
		t := (float64)((v - 273150) / 1000.0)
		text = fmt.Sprintf("Temperature: %.1f", t)
		drawText(text, 42, y, termbox.ColorWhite, termbox.ColorDefault)
	}
}

func (v *view) paintAttr(disk *atasmart.SkDisk) {
	w, h := getTermSize()
	attrs := disk.Attributes
	y := 8
	drawBgLine(1, y, w-2, termbox.ColorGreen)
	drawText("Status", 1, y, termbox.ColorWhite, termbox.ColorGreen)
	drawText("ID", 9, y, termbox.ColorWhite, termbox.ColorGreen)
	drawText("Attribute", 13, y, termbox.ColorWhite, termbox.ColorGreen)
	drawText("Current", 41, y, termbox.ColorWhite, termbox.ColorGreen)
	drawText("Worst", 49, y, termbox.ColorWhite, termbox.ColorGreen)
	drawText("Threshold", 56, y, termbox.ColorWhite, termbox.ColorGreen)
	drawText("Raw", 66, y, termbox.ColorWhite, termbox.ColorGreen)
	y++
	i := v.scroll
	for {
		if i >= len(attrs) {
			break
		} else if y >= h-1 {
			// clear
			drawText("  ", 1, h-2, termbox.ColorBlack, termbox.ColorWhite)
			drawText("↓", 1, h-2, termbox.ColorBlack, termbox.ColorWhite)
			break
		}
		at := attrs[i]
		bg := termbox.ColorDefault
		st := "Good"
		if at.Threshold != 0 && at.Current < at.Threshold {
			bg = termbox.ColorRed
			st = "Bad"
		}
		drawBgLine(1, y, w-2, bg)
		drawText(st, 3, y, termbox.ColorWhite, bg)
		text := fmt.Sprintf("%02X", at.ID)
		drawText(text, 9, y, termbox.ColorWhite, bg)
		drawText(at.Name, 13, y, termbox.ColorWhite, bg)
		text = fmt.Sprintf("%03d", at.Current)
		drawText(text, 44, y, termbox.ColorWhite, bg)
		text = fmt.Sprintf("%03d", at.Worst)
		drawText(text, 50, y, termbox.ColorWhite, bg)
		text = fmt.Sprintf("% 3d", at.Threshold)
		drawText(text, 62, y, termbox.ColorWhite, bg)
		text = fmt.Sprintf("%012X", at.Raw)
		drawText(text, 66, y, termbox.ColorWhite, bg)
		i++
		y++
	}
	if v.scrollable && v.scroll > 0 {
		drawText("  ", 1, 9, termbox.ColorBlack, termbox.ColorWhite)
		drawText("↑", 1, 9, termbox.ColorBlack, termbox.ColorWhite)
	}
}

func (v *view) update() {
	v.paintScreen()
	termbox.Flush()
}

func (v *view) updateAttributes() {
	v.paintAttr(v.disks[v.page])
	termbox.Flush()
}

func (v *view) Loop() error {
	v.update()
	v.updateScrollable()
	evCh := make(chan termbox.Event)
	go func() {
		for {
			evCh <- termbox.PollEvent()
		}
	}()

loop:
	for {
		select {
		case ev := <-evCh:
			if ev.Type == termbox.EventResize {
				setTermSize(ev.Width, ev.Height)
				v.scroll = 0
				v.updateScrollable()
				termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
				v.update()
				continue
			}
			v.handleKey(ev.Key)
		case <-v.exitCh:
			break loop
		}

	}

	// terminate
	for _, d := range v.disks {
		d.Free()
	}
	v.disks = nil
	return nil
}

const occupiedspace = 10

func (v *view) updateScrollable() {
	_, h := getTermSize()
	v.scrollable = len(v.disks[v.page].Attributes) > (h - occupiedspace)
}

func (v *view) handleKey(key termbox.Key) {
	switch key {
	case termbox.KeyArrowLeft:
		if v.page >= 1 {
			v.page--
			v.scroll = 0
			v.update()
		}
	case termbox.KeyArrowRight:
		if v.page < len(v.disks)-1 {
			v.page++
			v.scroll = 0
			v.update()
		}
	case termbox.KeyArrowUp:
		if v.scrollable && v.scroll > 0 {
			v.scroll--
			v.updateAttributes()

		}
	case termbox.KeyArrowDown:
		_, h := getTermSize()
		if v.scrollable && v.scroll < len(v.disks[v.page].Attributes)-(h-occupiedspace) {
			v.scroll++
			v.updateAttributes()
		}
	case termbox.KeyCtrlQ:
		v.exitCh <- struct{}{}
	}
}
