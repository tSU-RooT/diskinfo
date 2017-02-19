//+build !debug

package atasmart

/*
#cgo LDFLAGS: -latasmart
#include<stdlib.h>
#include<atasmart.h>

void appendAttribute(void*, SkSmartAttributeParsedData);

static void parse_attr_callback(SkDisk *d,
                                const SkSmartAttributeParsedData *a,
                                void *userdata) {
  SkSmartAttributeParsedData dat = *a;
  appendAttribute(userdata, dat);
}

static int attribute_wrapper(SkDisk *d, void * datap) {
  return sk_disk_smart_parse_attributes(d, parse_attr_callback, datap);
}

*/
import "C"
import (
	"unsafe"
)

const (
	DebugEnabled = false
)

func DiskOpen(path string) (*SkDisk, error) {
	var tmp *C.struct_SkDisk
	cn := C.CString(path)
	p1 := (**C.struct_SkDisk)(unsafe.Pointer(&tmp))
	check := (int)(C.sk_disk_open(cn, p1))
	if check != 0 {
		return nil, InternalError
	}
	return &SkDisk{internal: tmp, Path: path}, nil
}

func (sd *SkDisk) GetSize() (uint64, error) {
	var size uint64
	check := (int)(C.sk_disk_get_size((*C.struct_SkDisk)(sd.internal), (*C.uint64_t)(unsafe.Pointer(&size))))
	if check != 0 {
		return 0, InternalError
	}
	return size, nil
}

func (sd *SkDisk) IsSleepMode() (bool, error) {
	var b uint
	check := (int)(C.sk_disk_check_sleep_mode((*C.struct_SkDisk)(sd.internal), (*C.SkBool)(unsafe.Pointer(&b))))
	if check != 0 {
		return false, InternalError
	}
	if b == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (sd *SkDisk) ParseIdentify() error {
	var pd *C.SkIdentifyParsedData
	check := (int)(C.sk_disk_identify_parse((*C.struct_SkDisk)(sd.internal), (**C.SkIdentifyParsedData)(unsafe.Pointer(&pd))))
	if check != 0 {
		return InternalError
	}
	if pd != nil {
		dat := &SkDiskData{}
		dat.Serial = C.GoString(&(pd.serial[0]))
		dat.Firmware = C.GoString(&(pd.firmware[0]))
		dat.Model = C.GoString(&(pd.model[0]))
		sd.Data = dat
		return nil
	}

	return nil
}

func (sd *SkDisk) ReadSmartData() error {
	check := (int)(C.sk_disk_smart_read_data((*C.struct_SkDisk)(sd.internal)))
	if check != 0 {
		return InternalError
	}
	return nil
}

func (sd *SkDisk) GetSmartPowerCycle() (uint64, error) {
	var pcy uint64
	check := (int)(C.sk_disk_smart_get_power_cycle((*C.struct_SkDisk)(sd.internal), (*C.uint64_t)(unsafe.Pointer(&pcy))))
	if check != 0 {
		return 0, InternalError
	}
	return pcy, nil
}

func (sd *SkDisk) GetSmartPowerOn() (uint64, error) {
	var poweron uint64
	check := (int)(C.sk_disk_smart_get_power_on((*C.struct_SkDisk)(sd.internal), (*C.uint64_t)(unsafe.Pointer(&poweron))))
	if check != 0 {
		return 0, InternalError
	}
	return poweron, nil
}

func (sd *SkDisk) GetSmartTemperature() (uint64, error) {
	var tem uint64
	check := (int)(C.sk_disk_smart_get_temperature((*C.struct_SkDisk)(sd.internal), (*C.uint64_t)(unsafe.Pointer(&tem))))
	if check != 0 {
		return 0, InternalError
	}
	return tem, nil
}

func (sd *SkDisk) ParseSmartAttr() error {
	d := (*C.struct_SkDisk)(sd.internal)
	ars := make([]Attribute, 0)
	p := unsafe.Pointer(&ars)
	check := C.attribute_wrapper(d, p)
	if check != 0 {
		return InternalError
	}

	sd.Attributes = ([]Attribute)(*(*[]Attribute)(p))
	return nil
}

//export appendAttribute
func appendAttribute(skdp unsafe.Pointer, e C.SkSmartAttributeParsedData) {
	address := (*[]Attribute)(skdp)
	slice := *address
	a := Attribute{}
	a.ID = (uint8)(e.id)
	a.Name = C.GoString(e.name)
	a.Current = (uint8)(e.current_value)
	a.Worst = (uint8)(e.worst_value)
	a.Threshold = (uint8)(e.threshold)
	r := ([6]C.uint8_t)(e.raw)
	for i := (uint64)(0); i < 6; i++ {
		v := (uint64)(r[i])
		v = v << (8 * i)
		a.Raw += v
	}
	slice = append(slice, a)
	*address = slice
}

func (sd *SkDisk) Free() {
	C.sk_disk_free((*C.struct_SkDisk)(sd.internal))
}

func (sd *SkDisk) DiskDump() int {
	r := (int)(C.sk_disk_dump((*C.struct_SkDisk)(sd.internal)))
	return r
}
