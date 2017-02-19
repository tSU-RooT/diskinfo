package atasmart

/*
#include<stdlib.h>
#include<atasmart.h>
*/
import "C"
import "errors"

var (
	InternalError = errors.New("atasmart internal Error")
)

type SkDisk struct {
	internal *C.struct_SkDisk

	Data       *SkDiskData
	Attributes []Attribute
	Path       string
}

type SkDiskData struct {
	Serial   string
	Firmware string
	Model    string
}

type Attribute struct {
	ID        uint8
	Name      string
	Current   uint8
	Worst     uint8
	Threshold uint8
	Raw       uint64
}
