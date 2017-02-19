package main

import (
	"errors"
	"fmt"
	"github.com/nsf/termbox-go"
	"github.com/tSU-RooT/diskinfo/atasmart"
	"io/ioutil"
	"os"
	"strings"
)

var (
	termWidth  int
	termHeight int
)

func main() {
	os.Exit(_main())
}

func _main() int {
	if os.Getuid() != 0 && !atasmart.DebugEnabled {
		fmt.Fprint(os.Stderr, "Root privileges are required\n")
		return 1
	}
	view := initView()
	if ds, err := setup(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	} else {
		view.disks = ds
	}

	if err := termbox.Init(); err != nil {
		fmt.Fprint(os.Stderr, "termbox init failed\n")
		return 1
	}
	defer func() {
		if err := recover(); err != nil {
			termbox.Close()
			panic(err)
		}
	}()
	if w, h := termbox.Size(); w < 80 {
		termbox.Close()
		fmt.Fprint(os.Stderr, "Width of terminal is small\n")
		return 1
	} else {
		setTermSize(w, h)
	}
	view.paintScreen()
	view.Loop()
	termbox.Close()
	return 0
}

func listDisks() ([]string, error) {
	// When Debug flag is enabled by
	// `$ go build -tags debug`
	if atasmart.DebugEnabled {
		return []string{"/dev/sda", "/dev/sdb"}, nil
	}

	finfo, err := ioutil.ReadDir("/sys/block")
	if err != nil {
		return nil, fmt.Errorf("ReadErr: /sys/block")
	}

	list := make([]string, 0, 3)
	for _, f := range finfo {
		n := f.Name()
		if strings.HasPrefix(n, "ram") ||
			strings.HasPrefix(n, "loop") {
			continue
		}
		list = append(list, "/dev/"+n)
	}
	if len(list) > 6 {
		list = list[:6]
	}
	return list, nil

}

func setup() ([]*atasmart.SkDisk, error) {
	ld, err := listDisks()
	if err != nil {
		return nil, err
	}

	available := make([]*atasmart.SkDisk, 0, 3)
	for _, dn := range ld {
		_ = dn
		var disk *atasmart.SkDisk
		var err error
		if disk, err = atasmart.DiskOpen(dn); err != nil {
			continue
		}
		if err = disk.ParseIdentify(); err != nil {
			continue
		}
		if err = disk.ReadSmartData(); err != nil {
			continue
		}
		if err = disk.ParseSmartAttr(); err != nil {
			continue
		}
		available = append(available, disk)
	}
	if len(available) == 0 {
		return nil, errors.New("No available disk")
	}
	return available, nil
}
