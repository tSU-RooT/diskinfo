# DiskInfo

![ss01](https://raw.githubusercontent.com/wiki/tSU-RooT/diskinfo/img/ss01.png)
## Introduction
DiskInfo is TUI S.M.A.R.T viewer for Unix systems.  
Clone of [CrazyDiskInfo](https://github.com/otakuto/CrazyDiskInfo)  
CrazyDiskInfo is written in C++ and using ncurses library.  
But this reimplementation is written in Go and using termbox-go instead ncurses.  

## Getting Started

### Required C library
libatasmart4  

When you have already installed golang environment  

```
$ sudo apt install libatasmart-dev
$ go get github.com/tSU-RooT/diskinfo
$ sudo diskinfo
```

### Just want binary?
Download archieve from release page.  
You must install `libatasmart4` before running.  

### How to operate

|Key|Command|
|:---|:---|
|<kbd>Ctrl-q</kbd>|Quit from application  |
|Arrow-Key|Switch target or Scroll screen |
